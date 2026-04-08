package reqrisk

import "testing"

func TestPresetSensitiveSurfacesBorderlineSignals(t *testing.T) {
	features := Features{
		Entropy: &EntropyFeature{
			Value:      7.1,
			SampleSize: 512,
			Source:     "body",
		},
		Complexity: &ComplexityFeature{
			Depth:          5,
			FieldCount:     30,
			MaxArrayLength: 15,
			Source:         "json",
		},
		Charset: &CharsetFeature{
			NonASCIIRatio:    0.80,
			ControlRatio:     0.01,
			ZeroWidthCount:   1,
			MixedScripts:     false,
			ReplacementCount: 1,
			Source:           "body",
		},
		Fingerprint: &FingerprintFeature{
			Rarity:     0.75,
			Volatility: 0.50,
			Source:     "request-shape",
		},
	}

	defaultResult := New().Evaluate(features)
	sensitiveResult := New(WithPreset(PresetSensitive)).Evaluate(features)

	if defaultResult.Score != 0 {
		t.Fatalf("expected default preset to stay below threshold for borderline features, got %d", defaultResult.Score)
	}

	if sensitiveResult.Score <= defaultResult.Score {
		t.Fatalf("sensitive preset should produce a higher score, got default=%d sensitive=%d", defaultResult.Score, sensitiveResult.Score)
	}

	if len(sensitiveResult.Signals) != 4 {
		t.Fatalf("expected all four signals to contribute under sensitive preset, got %d", len(sensitiveResult.Signals))
	}
}

func TestEvaluateWithoutSignalsSuggestsObserve(t *testing.T) {
	result := New().Evaluate(Features{})

	if len(result.Signals) != 0 {
		t.Fatalf("expected no signals, got %d", len(result.Signals))
	}

	if len(result.Suggestions) != 1 || result.Suggestions[0] != SuggestObserve {
		t.Fatalf("unexpected suggestions for empty features: %v", result.Suggestions)
	}
}
