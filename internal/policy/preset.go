package policy

import "math"

// Preset selects a ready-made scoring profile for the assessor.
type Preset string

const (
	// PresetBalanced matches the default baseline and is suitable for general use.
	PresetBalanced Preset = "balanced"
	// PresetSensitive lowers thresholds so borderline signals surface earlier.
	PresetSensitive Preset = "sensitive"
	// PresetConservative raises thresholds so only stronger signals contribute.
	PresetConservative Preset = "conservative"
)

// WithPreset applies a lightweight scoring preset.
//
// Presets are intentionally opinionated starting points, not a replacement for
// explicit policy overrides. Later options can still adjust individual fields.
func WithPreset(p Preset) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		applyPreset(cfg, p)
	}
}

func applyPreset(cfg *Config, preset Preset) {
	switch preset {
	case PresetSensitive:
		applySensitivePreset(cfg)
	case PresetConservative:
		applyConservativePreset(cfg)
	default:
		applyBalancedPreset(cfg)
	}
}

func applyBalancedPreset(cfg *Config) {
	if cfg == nil {
		return
	}

	cfg.Bands = defaultBandsForMaxScore(cfg.MaxScore)
	cfg.bandsCustomized = false
	cfg.Entropy = EntropyPolicy{
		Weight:    1,
		MinSample: 256,
		Elevated:  7.2,
		High:      7.6,
		Extreme:   7.85,
	}
	cfg.Complexity = ComplexityPolicy{
		Weight:         1,
		DepthHigh:      6,
		FieldCountHigh: 40,
		ArrayHigh:      20,
	}
	cfg.Charset = CharsetPolicy{
		Weight:             1,
		NonASCIIHigh:       0.85,
		ControlHigh:        0.02,
		ZeroWidthHigh:      2,
		ReplacementHigh:    2,
		MixedScriptsWeight: 20,
	}
	cfg.Fingerprint = FingerprintPolicy{
		Weight:         1,
		RarityHigh:     0.80,
		VolatilityHigh: 0.60,
	}
}

func applySensitivePreset(cfg *Config) {
	if cfg == nil {
		return
	}

	cfg.Bands = presetBandsForMaxScore(cfg.MaxScore, 0.20, 0.40, 0.60)
	cfg.bandsCustomized = true
	cfg.Entropy = EntropyPolicy{
		Weight:    1,
		MinSample: 128,
		Elevated:  7.0,
		High:      7.35,
		Extreme:   7.65,
	}
	cfg.Complexity = ComplexityPolicy{
		Weight:         1,
		DepthHigh:      4,
		FieldCountHigh: 24,
		ArrayHigh:      12,
	}
	cfg.Charset = CharsetPolicy{
		Weight:             1,
		NonASCIIHigh:       0.75,
		ControlHigh:        0.01,
		ZeroWidthHigh:      1,
		ReplacementHigh:    1,
		MixedScriptsWeight: 30,
	}
	cfg.Fingerprint = FingerprintPolicy{
		Weight:         1,
		RarityHigh:     0.70,
		VolatilityHigh: 0.45,
	}
}

func applyConservativePreset(cfg *Config) {
	if cfg == nil {
		return
	}

	cfg.Bands = presetBandsForMaxScore(cfg.MaxScore, 0.30, 0.60, 0.85)
	cfg.bandsCustomized = true
	cfg.Entropy = EntropyPolicy{
		Weight:    1,
		MinSample: 512,
		Elevated:  7.4,
		High:      7.7,
		Extreme:   7.9,
	}
	cfg.Complexity = ComplexityPolicy{
		Weight:         1,
		DepthHigh:      8,
		FieldCountHigh: 60,
		ArrayHigh:      30,
	}
	cfg.Charset = CharsetPolicy{
		Weight:             1,
		NonASCIIHigh:       0.92,
		ControlHigh:        0.04,
		ZeroWidthHigh:      3,
		ReplacementHigh:    3,
		MixedScriptsWeight: 12,
	}
	cfg.Fingerprint = FingerprintPolicy{
		Weight:         1,
		RarityHigh:     0.88,
		VolatilityHigh: 0.75,
	}
}

func presetBandsForMaxScore(maxScore int, mediumRatio, highRatio, criticalRatio float64) LevelBands {
	if maxScore <= 0 {
		maxScore = defaultMaxScore
	}

	medium := int(math.Round(float64(maxScore) * mediumRatio))
	high := int(math.Round(float64(maxScore) * highRatio))
	critical := int(math.Round(float64(maxScore) * criticalRatio))

	return LevelBands{
		Medium:   clampInt(medium, 0, maxScore),
		High:     clampInt(high, medium, maxScore),
		Critical: clampInt(critical, high, maxScore),
	}
}
