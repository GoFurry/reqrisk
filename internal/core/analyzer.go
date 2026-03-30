package core

import (
	"math"
	"sort"

	"github.com/GoFurry/reqrisk/internal/model"
	"github.com/GoFurry/reqrisk/internal/policy"
)

type analyzer interface {
	key() string
	enabled(policy.Config) bool
	assess(model.Features, policy.Config) model.SignalReport
}

type assessedFinding struct {
	finding    model.Finding
	localScore int
	tags       []string
}

func defaultAnalyzers() []analyzer {
	return []analyzer{
		entropyAnalyzer{},
		complexityAnalyzer{},
		charsetAnalyzer{},
		fingerprintAnalyzer{},
	}
}

func finalizeSignal(key string, score int, weight float64, cfg policy.Config, findings []assessedFinding, tags []string) model.SignalReport {
	score = clampInt(score, 0, 100)
	contribution := scoreToContribution(score, weight, cfg)

	report := model.SignalReport{
		Key:          key,
		Score:        score,
		Contribution: contribution,
	}

	if len(findings) == 0 {
		report.Tags = dedupeStrings(tags)
		return report
	}

	localScores := make([]int, 0, len(findings))
	report.Findings = make([]model.Finding, 0, len(findings))
	for _, finding := range findings {
		localScores = append(localScores, maxInt(0, finding.localScore))
		tags = append(tags, finding.tags...)
		report.Findings = append(report.Findings, finding.finding)
	}

	allocated := allocateContributions(localScores, contribution)
	for i := range report.Findings {
		report.Findings[i].Contribution = allocated[i]
	}

	report.Tags = dedupeStrings(tags)
	return report
}

func scoreToContribution(score int, weight float64, cfg policy.Config) int {
	if score <= 0 || weight <= 0 || cfg.MaxScore <= 0 {
		return 0
	}

	totalWeight := totalWeight(cfg)
	if totalWeight <= 0 {
		return 0
	}

	exact := (float64(score) / 100.0) * (weight / totalWeight) * float64(cfg.MaxScore)
	return clampInt(int(math.Round(exact)), 0, cfg.MaxScore)
}

func totalWeight(cfg policy.Config) float64 {
	total := 0.0
	total += maxFloat64(0, cfg.Entropy.Weight)
	total += maxFloat64(0, cfg.Complexity.Weight)
	total += maxFloat64(0, cfg.Charset.Weight)
	total += maxFloat64(0, cfg.Fingerprint.Weight)
	return total
}

func allocateContributions(localScores []int, total int) []int {
	allocated := make([]int, len(localScores))
	if total <= 0 || len(localScores) == 0 {
		return allocated
	}

	sumLocal := 0
	for _, score := range localScores {
		sumLocal += maxInt(0, score)
	}
	if sumLocal == 0 {
		return allocated
	}

	type remainder struct {
		index int
		value float64
	}

	remainders := make([]remainder, 0, len(localScores))
	used := 0
	for i, score := range localScores {
		exact := float64(total) * float64(maxInt(0, score)) / float64(sumLocal)
		base := int(math.Floor(exact))
		allocated[i] = base
		used += base
		remainders = append(remainders, remainder{
			index: i,
			value: exact - float64(base),
		})
	}

	sort.SliceStable(remainders, func(i, j int) bool {
		return remainders[i].value > remainders[j].value
	})

	for i := 0; i < total-used && i < len(remainders); i++ {
		allocated[remainders[i].index]++
	}

	return allocated
}
