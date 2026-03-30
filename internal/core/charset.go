package core

import (
	"github.com/GoFurry/reqrisk/internal/model"
	"github.com/GoFurry/reqrisk/internal/policy"
)

type charsetAnalyzer struct{}

// AssessCharset evaluates a single charset feature with the default policy.
func AssessCharset(feature model.CharsetFeature) model.SignalReport {
	cfg := singleSignalConfig("charset")
	return assessCharset(feature, cfg)
}

func (charsetAnalyzer) key() string {
	return "charset"
}

func (charsetAnalyzer) enabled(cfg policy.Config) bool {
	return cfg.Charset.Weight > 0
}

func (charsetAnalyzer) assess(features model.Features, cfg policy.Config) model.SignalReport {
	if features.Charset == nil {
		return model.SignalReport{}
	}
	return assessCharset(*features.Charset, cfg)
}

func assessCharset(feature model.CharsetFeature, cfg policy.Config) model.SignalReport {
	findings := make([]assessedFinding, 0, 5)

	appendFinding := func(key string, severity model.Severity, summary string, localScore int, tag string, evidence []model.Evidence) {
		findings = append(findings, assessedFinding{
			finding: model.Finding{
				Key:      key,
				Signal:   "charset",
				Severity: severity,
				Summary:  summary,
				Evidence: evidence,
			},
			localScore: localScore,
			tags:       []string{tag},
		})
	}

	if tier := tieredFloat(feature.NonASCIIRatio, cfg.Charset.NonASCIIHigh, secondaryFloatThreshold(cfg.Charset.NonASCIIHigh, 0.12, 1)); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "non_ascii_ratio", Value: formatFloat(feature.NonASCIIRatio)},
			{Key: "threshold", Value: formatFloat(cfg.Charset.NonASCIIHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("charset_non_ascii_heavy", model.SeverityHigh, "Non-ASCII usage is unusually heavy for the configured baseline.", 30, "charset_non_ascii_heavy", evidence)
		} else {
			appendFinding("charset_non_ascii_elevated", model.SeverityMedium, "Non-ASCII usage is elevated for the configured baseline.", 20, "charset_non_ascii_elevated", evidence)
		}
	}

	if tier := tieredFloat(feature.ControlRatio, cfg.Charset.ControlHigh, secondaryFloatThreshold(cfg.Charset.ControlHigh, cfg.Charset.ControlHigh, 1)); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "control_ratio", Value: formatFloat(feature.ControlRatio)},
			{Key: "threshold", Value: formatFloat(cfg.Charset.ControlHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("charset_control_chars_extreme", model.SeverityCritical, "Control-character density is far above the configured baseline.", 45, "charset_control_chars_extreme", evidence)
		} else {
			appendFinding("charset_control_chars_high", model.SeverityHigh, "Control-character density is above the configured baseline.", 30, "charset_control_chars_high", evidence)
		}
	}

	if tier := tieredInt(feature.ZeroWidthCount, cfg.Charset.ZeroWidthHigh); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "zero_width_count", Value: formatInt(feature.ZeroWidthCount)},
			{Key: "threshold", Value: formatInt(cfg.Charset.ZeroWidthHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("charset_zero_width_heavy", model.SeverityHigh, "Zero-width characters appear frequently.", 30, "charset_zero_width_heavy", evidence)
		} else {
			appendFinding("charset_zero_width_present", model.SeverityMedium, "Zero-width characters are present above the configured baseline.", 20, "charset_zero_width_present", evidence)
		}
	}

	if tier := tieredInt(feature.ReplacementCount, cfg.Charset.ReplacementHigh); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "replacement_count", Value: formatInt(feature.ReplacementCount)},
			{Key: "threshold", Value: formatInt(cfg.Charset.ReplacementHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("charset_replacement_chars_extreme", model.SeverityCritical, "Replacement-character count is far above the configured baseline.", 35, "charset_replacement_chars_extreme", evidence)
		} else {
			appendFinding("charset_replacement_chars_high", model.SeverityHigh, "Replacement-character count is above the configured baseline.", 25, "charset_replacement_chars_high", evidence)
		}
	}

	if feature.MixedScripts && cfg.Charset.MixedScriptsWeight > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "mixed_scripts", Value: "true"},
		}, feature.Source)
		appendFinding("charset_mixed_scripts", model.SeverityMedium, "Mixed scripts were observed in the analyzed text.", clampInt(int(cfg.Charset.MixedScriptsWeight), 0, 100), "charset_mixed_scripts", evidence)
	}

	score := 0
	for _, finding := range findings {
		score += finding.localScore
	}

	return finalizeSignal("charset", score, cfg.Charset.Weight, cfg, findings, nil)
}

func tieredFloat(value, threshold, extreme float64) int {
	if threshold <= 0 || value < threshold {
		return 0
	}
	if extreme < threshold {
		extreme = threshold
	}
	if value >= extreme {
		return 2
	}
	return 1
}

func secondaryFloatThreshold(threshold, delta, max float64) float64 {
	if threshold <= 0 {
		return 0
	}
	value := threshold + delta
	if value > max {
		return max
	}
	return value
}
