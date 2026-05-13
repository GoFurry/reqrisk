package core

import (
	"github.com/gofurry/reqrisk/internal/model"
	"github.com/gofurry/reqrisk/internal/policy"
)

// Assessor evaluates request-related features into explainable risk results.
type Assessor struct {
	cfg       policy.Config
	analyzers []analyzer
}

// New constructs an assessor with the default configuration plus any supplied options.
func New(opts ...policy.Option) *Assessor {
	cfg := policy.DefaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	cfg = policy.NormalizeConfig(cfg)

	return &Assessor{
		cfg:       cfg,
		analyzers: defaultAnalyzers(),
	}
}

// Evaluate runs the built-in analyzers and returns an aggregate explainable result.
func (a *Assessor) Evaluate(features model.Features) model.Result {
	if a == nil {
		return New().Evaluate(features)
	}

	result := model.Result{
		Signals:  make([]model.SignalReport, 0, len(a.analyzers)),
		Findings: make([]model.Finding, 0, 8),
	}

	for _, analyzer := range a.analyzers {
		if !analyzer.enabled(a.cfg) {
			continue
		}

		report := analyzer.assess(features, a.cfg)
		if report.Key == "" {
			continue
		}

		result.Signals = append(result.Signals, report)
		result.Score += report.Contribution
		result.Findings = append(result.Findings, report.Findings...)
		result.Tags = append(result.Tags, report.Tags...)
	}

	result.Score = clampInt(result.Score, 0, a.cfg.MaxScore)
	result.Level = levelFromScore(result.Score, a.cfg.Bands)
	result.Tags = dedupeStrings(result.Tags)
	result.Suggestions = deriveSuggestions(result.Level, result.Findings)

	if len(result.Signals) == 0 {
		result.Suggestions = []model.Suggestion{model.SuggestObserve}
	}

	return result
}

func levelFromScore(score int, bands policy.LevelBands) model.Level {
	switch {
	case score >= bands.Critical:
		return model.LevelCritical
	case score >= bands.High:
		return model.LevelHigh
	case score >= bands.Medium:
		return model.LevelMedium
	default:
		return model.LevelLow
	}
}

func deriveSuggestions(level model.Level, findings []model.Finding) []model.Suggestion {
	suggestions := make([]model.Suggestion, 0, 3)

	switch level {
	case model.LevelCritical:
		suggestions = append(suggestions, model.SuggestChallenge, model.SuggestBlock)
	case model.LevelHigh:
		suggestions = append(suggestions, model.SuggestReview, model.SuggestChallenge)
	case model.LevelMedium:
		suggestions = append(suggestions, model.SuggestObserve, model.SuggestReview)
	default:
		suggestions = append(suggestions, model.SuggestObserve)
	}

	maxSeverity := model.SeverityInfo
	for _, finding := range findings {
		if severityRank(finding.Severity) > severityRank(maxSeverity) {
			maxSeverity = finding.Severity
		}
	}

	switch {
	case severityRank(maxSeverity) >= severityRank(model.SeverityCritical):
		suggestions = append(suggestions, model.SuggestBlock)
	case severityRank(maxSeverity) >= severityRank(model.SeverityHigh):
		suggestions = append(suggestions, model.SuggestChallenge)
	case severityRank(maxSeverity) >= severityRank(model.SeverityMedium):
		suggestions = append(suggestions, model.SuggestReview)
	}

	return dedupeSuggestions(suggestions)
}

func dedupeSuggestions(values []model.Suggestion) []model.Suggestion {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[model.Suggestion]struct{}, len(values))
	out := make([]model.Suggestion, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
