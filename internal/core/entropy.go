package core

import (
	"github.com/GoFurry/reqrisk/internal/model"
	"github.com/GoFurry/reqrisk/internal/policy"
)

type entropyAnalyzer struct{}

// AssessEntropy evaluates a single entropy feature with the default policy.
func AssessEntropy(feature model.EntropyFeature) model.SignalReport {
	cfg := singleSignalConfig("entropy")
	return assessEntropy(feature, cfg)
}

func (entropyAnalyzer) key() string {
	return "entropy"
}

func (entropyAnalyzer) enabled(cfg policy.Config) bool {
	return cfg.Entropy.Weight > 0
}

func (entropyAnalyzer) assess(features model.Features, cfg policy.Config) model.SignalReport {
	if features.Entropy == nil {
		return model.SignalReport{}
	}
	return assessEntropy(*features.Entropy, cfg)
}

func assessEntropy(feature model.EntropyFeature, cfg policy.Config) model.SignalReport {
	findings := make([]assessedFinding, 0, 1)

	if feature.SampleSize >= cfg.Entropy.MinSample {
		evidence := appendSourceEvidence([]model.Evidence{
			{Key: "value", Value: formatFloat(feature.Value)},
			{Key: "sample_size", Value: formatInt(feature.SampleSize)},
		}, feature.Source)

		switch {
		case feature.Value >= cfg.Entropy.Extreme && cfg.Entropy.Extreme > 0:
			findings = append(findings, assessedFinding{
				finding: model.Finding{
					Key:      "entropy_extreme",
					Signal:   "entropy",
					Severity: model.SeverityCritical,
					Summary:  "Entropy is extremely high for the observed sample.",
					Evidence: evidence,
				},
				localScore: 90,
				tags:       []string{"entropy_extreme"},
			})
		case feature.Value >= cfg.Entropy.High && cfg.Entropy.High > 0:
			findings = append(findings, assessedFinding{
				finding: model.Finding{
					Key:      "entropy_high",
					Signal:   "entropy",
					Severity: model.SeverityHigh,
					Summary:  "Entropy is unusually high for the observed sample.",
					Evidence: evidence,
				},
				localScore: 70,
				tags:       []string{"entropy_high"},
			})
		case feature.Value >= cfg.Entropy.Elevated && cfg.Entropy.Elevated > 0:
			findings = append(findings, assessedFinding{
				finding: model.Finding{
					Key:      "entropy_elevated",
					Signal:   "entropy",
					Severity: model.SeverityMedium,
					Summary:  "Entropy is elevated relative to the configured baseline.",
					Evidence: evidence,
				},
				localScore: 40,
				tags:       []string{"entropy_elevated"},
			})
		}
	}

	score := 0
	for _, finding := range findings {
		score += finding.localScore
	}

	return finalizeSignal("entropy", score, cfg.Entropy.Weight, cfg, findings, nil)
}
