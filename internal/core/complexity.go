package core

import (
	"github.com/GoFurry/reqrisk/internal/model"
	"github.com/GoFurry/reqrisk/internal/policy"
)

type complexityAnalyzer struct{}

// AssessComplexity evaluates a single complexity feature with the default policy.
func AssessComplexity(feature model.ComplexityFeature) model.SignalReport {
	cfg := singleSignalConfig("complexity")
	return assessComplexity(feature, cfg)
}

func (complexityAnalyzer) key() string {
	return "complexity"
}

func (complexityAnalyzer) enabled(cfg policy.Config) bool {
	return cfg.Complexity.Weight > 0
}

func (complexityAnalyzer) assess(features model.Features, cfg policy.Config) model.SignalReport {
	if features.Complexity == nil {
		return model.SignalReport{}
	}
	return assessComplexity(*features.Complexity, cfg)
}

func assessComplexity(feature model.ComplexityFeature, cfg policy.Config) model.SignalReport {
	findings := make([]assessedFinding, 0, 3)

	appendFinding := func(key string, severity model.Severity, summary string, localScore int, tag string, evidence []model.Evidence) {
		findings = append(findings, assessedFinding{
			finding: model.Finding{
				Key:      key,
				Signal:   "complexity",
				Severity: severity,
				Summary:  summary,
				Evidence: evidence,
			},
			localScore: localScore,
			tags:       []string{tag},
		})
	}

	if tier := tieredInt(feature.Depth, cfg.Complexity.DepthHigh); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "depth", Value: formatInt(feature.Depth)},
			{Key: "threshold", Value: formatInt(cfg.Complexity.DepthHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("complexity_depth_extreme", model.SeverityHigh, "Structural depth is far above the configured baseline.", 45, "complexity_depth_extreme", evidence)
		} else {
			appendFinding("complexity_depth_high", model.SeverityMedium, "Structural depth is above the configured baseline.", 30, "complexity_depth_high", evidence)
		}
	}

	if tier := tieredInt(feature.FieldCount, cfg.Complexity.FieldCountHigh); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "field_count", Value: formatInt(feature.FieldCount)},
			{Key: "threshold", Value: formatInt(cfg.Complexity.FieldCountHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("complexity_fields_extreme", model.SeverityHigh, "Field volume is far above the configured baseline.", 45, "complexity_fields_extreme", evidence)
		} else {
			appendFinding("complexity_fields_high", model.SeverityMedium, "Field volume is above the configured baseline.", 30, "complexity_fields_high", evidence)
		}
	}

	if tier := tieredInt(feature.MaxArrayLength, cfg.Complexity.ArrayHigh); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "max_array_length", Value: formatInt(feature.MaxArrayLength)},
			{Key: "threshold", Value: formatInt(cfg.Complexity.ArrayHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("complexity_array_span_extreme", model.SeverityHigh, "Maximum array span is far above the configured baseline.", 30, "complexity_array_span_extreme", evidence)
		} else {
			appendFinding("complexity_array_span_high", model.SeverityMedium, "Maximum array span is above the configured baseline.", 20, "complexity_array_span_high", evidence)
		}
	}

	score := 0
	for _, finding := range findings {
		score += finding.localScore
	}

	return finalizeSignal("complexity", score, cfg.Complexity.Weight, cfg, findings, nil)
}

func tieredInt(value, threshold int) int {
	if threshold <= 0 || value < threshold {
		return 0
	}
	if value >= threshold*2 {
		return 2
	}
	return 1
}
