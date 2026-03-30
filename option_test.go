package reqrisk

import (
	"testing"

	internalpolicy "github.com/GoFurry/reqrisk/internal/policy"
)

func TestOptionsCanFocusScoringOnSingleSignal(t *testing.T) {
	assessor := New(
		WithMaxScore(80),
		WithBands(LevelBands{Medium: 20, High: 40, Critical: 60}),
		WithEntropyPolicy(EntropyPolicy{
			Weight:    1,
			MinSample: 128,
			Elevated:  7.0,
			High:      7.5,
			Extreme:   7.8,
		}),
		WithComplexityPolicy(ComplexityPolicy{}),
		WithCharsetPolicy(CharsetPolicy{}),
		WithFingerprintPolicy(FingerprintPolicy{}),
	)

	result := assessor.Evaluate(Features{
		Entropy: &EntropyFeature{
			Value:      7.9,
			SampleSize: 1024,
		},
	})

	if result.Score != 72 {
		t.Fatalf("unexpected score: got %d want 72", result.Score)
	}

	if result.Level != LevelCritical {
		t.Fatalf("unexpected level: got %q want %q", result.Level, LevelCritical)
	}

	if len(result.Signals) != 1 {
		t.Fatalf("unexpected signal count: got %d want 1", len(result.Signals))
	}
}

func TestNormalizeConfigClampsInvalidValues(t *testing.T) {
	cfg := internalpolicy.NormalizeConfig(Config{
		MaxScore: -10,
		Bands: LevelBands{
			Medium:   -1,
			High:     200,
			Critical: 300,
		},
		Entropy: EntropyPolicy{
			Weight:    -1,
			MinSample: -10,
			Elevated:  9,
			High:      -1,
			Extreme:   20,
		},
		Complexity: ComplexityPolicy{
			Weight:         -1,
			DepthHigh:      -1,
			FieldCountHigh: -10,
			ArrayHigh:      -5,
		},
		Charset: CharsetPolicy{
			Weight:             -1,
			NonASCIIHigh:       2,
			ControlHigh:        -1,
			ZeroWidthHigh:      -5,
			ReplacementHigh:    -5,
			MixedScriptsWeight: -3,
		},
		Fingerprint: FingerprintPolicy{
			Weight:         -1,
			RarityHigh:     2,
			VolatilityHigh: -1,
		},
	})

	if cfg.MaxScore != 100 {
		t.Fatalf("unexpected max score: got %d want %d", cfg.MaxScore, 100)
	}

	if cfg.Bands.Medium != 0 || cfg.Bands.High != 100 || cfg.Bands.Critical != 100 {
		t.Fatalf("unexpected normalized bands: %+v", cfg.Bands)
	}

	if cfg.Entropy.Weight != 0 || cfg.Entropy.MinSample != 0 || cfg.Entropy.Elevated != 8 || cfg.Entropy.High != 8 || cfg.Entropy.Extreme != 8 {
		t.Fatalf("unexpected normalized entropy policy: %+v", cfg.Entropy)
	}

	if cfg.Charset.NonASCIIHigh != 1 || cfg.Charset.ControlHigh != 0 || cfg.Fingerprint.RarityHigh != 1 || cfg.Fingerprint.VolatilityHigh != 0 {
		t.Fatalf("unexpected normalized ratio thresholds: charset=%+v fingerprint=%+v", cfg.Charset, cfg.Fingerprint)
	}
}
