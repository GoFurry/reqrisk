package policy

import "math"

const defaultMaxScore = 100

// Config controls scoring and aggregation behavior.
type Config struct {
	MaxScore int

	Bands LevelBands

	Entropy     EntropyPolicy
	Complexity  ComplexityPolicy
	Charset     CharsetPolicy
	Fingerprint FingerprintPolicy
}

// LevelBands maps aggregate scores into risk levels.
type LevelBands struct {
	Medium   int
	High     int
	Critical int
}

// EntropyPolicy controls entropy assessment thresholds.
type EntropyPolicy struct {
	Weight    float64
	MinSample int
	Elevated  float64
	High      float64
	Extreme   float64
}

// ComplexityPolicy controls structural complexity assessment thresholds.
type ComplexityPolicy struct {
	Weight         float64
	DepthHigh      int
	FieldCountHigh int
	ArrayHigh      int
}

// CharsetPolicy controls text anomaly assessment thresholds.
type CharsetPolicy struct {
	Weight             float64
	NonASCIIHigh       float64
	ControlHigh        float64
	ZeroWidthHigh      int
	ReplacementHigh    int
	MixedScriptsWeight float64
}

// FingerprintPolicy controls fingerprint rarity and volatility thresholds.
type FingerprintPolicy struct {
	Weight         float64
	RarityHigh     float64
	VolatilityHigh float64
}

func DefaultConfig() Config {
	return NormalizeConfig(Config{
		MaxScore: defaultMaxScore,
		Bands: LevelBands{
			Medium:   25,
			High:     50,
			Critical: 75,
		},
		Entropy: EntropyPolicy{
			Weight:    1,
			MinSample: 256,
			Elevated:  7.2,
			High:      7.6,
			Extreme:   7.85,
		},
		Complexity: ComplexityPolicy{
			Weight:         1,
			DepthHigh:      6,
			FieldCountHigh: 40,
			ArrayHigh:      20,
		},
		Charset: CharsetPolicy{
			Weight:             1,
			NonASCIIHigh:       0.85,
			ControlHigh:        0.02,
			ZeroWidthHigh:      2,
			ReplacementHigh:    2,
			MixedScriptsWeight: 20,
		},
		Fingerprint: FingerprintPolicy{
			Weight:         1,
			RarityHigh:     0.80,
			VolatilityHigh: 0.60,
		},
	})
}

func NormalizeConfig(cfg Config) Config {
	if cfg.MaxScore <= 0 {
		cfg.MaxScore = defaultMaxScore
	}

	cfg.Bands.Medium = clampInt(cfg.Bands.Medium, 0, cfg.MaxScore)
	cfg.Bands.High = clampInt(cfg.Bands.High, cfg.Bands.Medium, cfg.MaxScore)
	cfg.Bands.Critical = clampInt(cfg.Bands.Critical, cfg.Bands.High, cfg.MaxScore)

	cfg.Entropy.Weight = maxFloat64(0, cfg.Entropy.Weight)
	cfg.Entropy.MinSample = maxInt(0, cfg.Entropy.MinSample)
	cfg.Entropy.Elevated = clampFloat64(cfg.Entropy.Elevated, 0, 8)
	cfg.Entropy.High = clampFloat64(maxFloat64(cfg.Entropy.High, cfg.Entropy.Elevated), 0, 8)
	cfg.Entropy.Extreme = clampFloat64(maxFloat64(cfg.Entropy.Extreme, cfg.Entropy.High), 0, 8)

	cfg.Complexity.Weight = maxFloat64(0, cfg.Complexity.Weight)
	cfg.Complexity.DepthHigh = maxInt(0, cfg.Complexity.DepthHigh)
	cfg.Complexity.FieldCountHigh = maxInt(0, cfg.Complexity.FieldCountHigh)
	cfg.Complexity.ArrayHigh = maxInt(0, cfg.Complexity.ArrayHigh)

	cfg.Charset.Weight = maxFloat64(0, cfg.Charset.Weight)
	cfg.Charset.NonASCIIHigh = clampFloat64(cfg.Charset.NonASCIIHigh, 0, 1)
	cfg.Charset.ControlHigh = clampFloat64(cfg.Charset.ControlHigh, 0, 1)
	cfg.Charset.ZeroWidthHigh = maxInt(0, cfg.Charset.ZeroWidthHigh)
	cfg.Charset.ReplacementHigh = maxInt(0, cfg.Charset.ReplacementHigh)
	cfg.Charset.MixedScriptsWeight = maxFloat64(0, cfg.Charset.MixedScriptsWeight)

	cfg.Fingerprint.Weight = maxFloat64(0, cfg.Fingerprint.Weight)
	cfg.Fingerprint.RarityHigh = clampFloat64(cfg.Fingerprint.RarityHigh, 0, 1)
	cfg.Fingerprint.VolatilityHigh = clampFloat64(cfg.Fingerprint.VolatilityHigh, 0, 1)

	return cfg
}

func defaultBandsForMaxScore(maxScore int) LevelBands {
	if maxScore <= 0 {
		maxScore = defaultMaxScore
	}

	medium := int(math.Round(float64(maxScore) * 0.25))
	high := int(math.Round(float64(maxScore) * 0.50))
	critical := int(math.Round(float64(maxScore) * 0.75))

	return LevelBands{
		Medium:   clampInt(medium, 0, maxScore),
		High:     clampInt(high, medium, maxScore),
		Critical: clampInt(critical, high, maxScore),
	}
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
