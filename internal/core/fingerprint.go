package core

import (
	"github.com/GoFurry/reqrisk/internal/model"
	"github.com/GoFurry/reqrisk/internal/policy"
)

type fingerprintAnalyzer struct{}

// AssessFingerprint evaluates a single fingerprint feature with the default policy.
func AssessFingerprint(feature model.FingerprintFeature) model.SignalReport {
	cfg := singleSignalConfig("fingerprint")
	return assessFingerprint(feature, cfg)
}

func (fingerprintAnalyzer) key() string {
	return "fingerprint"
}

func (fingerprintAnalyzer) enabled(cfg policy.Config) bool {
	return cfg.Fingerprint.Weight > 0
}

func (fingerprintAnalyzer) assess(features model.Features, cfg policy.Config) model.SignalReport {
	if features.Fingerprint == nil {
		return model.SignalReport{}
	}
	return assessFingerprint(*features.Fingerprint, cfg)
}

func assessFingerprint(feature model.FingerprintFeature, cfg policy.Config) model.SignalReport {
	findings := make([]assessedFinding, 0, 2)

	appendFinding := func(key string, severity model.Severity, summary string, localScore int, tag string, evidence []model.Evidence) {
		findings = append(findings, assessedFinding{
			finding: model.Finding{
				Key:      key,
				Signal:   "fingerprint",
				Severity: severity,
				Summary:  summary,
				Evidence: evidence,
			},
			localScore: localScore,
			tags:       []string{tag},
		})
	}

	if tier := tieredFloat(feature.Rarity, cfg.Fingerprint.RarityHigh, secondaryFloatThreshold(cfg.Fingerprint.RarityHigh, 0.15, 1)); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "rarity", Value: formatFloat(feature.Rarity)},
			{Key: "threshold", Value: formatFloat(cfg.Fingerprint.RarityHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("fingerprint_rare_extreme", model.SeverityHigh, "Fingerprint rarity is far above the configured baseline.", 50, "fingerprint_rare_extreme", evidence)
		} else {
			appendFinding("fingerprint_rare_high", model.SeverityMedium, "Fingerprint rarity is above the configured baseline.", 35, "fingerprint_rare_high", evidence)
		}
	}

	if tier := tieredFloat(feature.Volatility, cfg.Fingerprint.VolatilityHigh, secondaryFloatThreshold(cfg.Fingerprint.VolatilityHigh, 0.20, 1)); tier > 0 {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "volatility", Value: formatFloat(feature.Volatility)},
			{Key: "threshold", Value: formatFloat(cfg.Fingerprint.VolatilityHigh)},
		}, feature.Source)
		if tier == 2 {
			appendFinding("fingerprint_volatile_extreme", model.SeverityHigh, "Fingerprint volatility is far above the configured baseline.", 35, "fingerprint_volatile_extreme", evidence)
		} else {
			appendFinding("fingerprint_volatile_high", model.SeverityMedium, "Fingerprint volatility is above the configured baseline.", 25, "fingerprint_volatile_high", evidence)
		}
	}

	score := 0
	for _, finding := range findings {
		score += finding.localScore
	}

	return finalizeSignal("fingerprint", score, cfg.Fingerprint.Weight, cfg, findings, nil)
}
