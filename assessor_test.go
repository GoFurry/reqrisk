package reqrisk

import "testing"

func TestEvaluateAggregatesSignalsIntoCriticalResult(t *testing.T) {
	assessor := New()

	result := assessor.Evaluate(Features{
		Entropy: &EntropyFeature{
			Value:      7.9,
			SampleSize: 1024,
			Source:     "body",
		},
		Complexity: &ComplexityFeature{
			Depth:          8,
			FieldCount:     96,
			MaxArrayLength: 42,
			Source:         "json",
		},
		Charset: &CharsetFeature{
			NonASCIIRatio:    0.92,
			ControlRatio:     0.04,
			ZeroWidthCount:   3,
			MixedScripts:     true,
			ReplacementCount: 2,
			Source:           "body",
		},
		Fingerprint: &FingerprintFeature{
			Rarity:     0.96,
			Volatility: 0.85,
			Source:     "request-shape",
		},
		Meta: Meta{
			Target: "body",
			Route:  "/api/comment",
			Method: "POST",
		},
	})

	if result.Score != 94 {
		t.Fatalf("unexpected score: got %d want 94", result.Score)
	}

	if result.Level != LevelCritical {
		t.Fatalf("unexpected level: got %q want %q", result.Level, LevelCritical)
	}

	if len(result.Signals) != 4 {
		t.Fatalf("unexpected signal count: got %d want 4", len(result.Signals))
	}

	if len(result.Findings) != 11 {
		t.Fatalf("unexpected flattened finding count: got %d want 11", len(result.Findings))
	}

	if len(result.Tags) == 0 {
		t.Fatal("expected aggregate tags")
	}

	if len(result.Suggestions) != 2 || result.Suggestions[0] != SuggestChallenge || result.Suggestions[1] != SuggestBlock {
		t.Fatalf("unexpected suggestions: got %v want [%q %q]", result.Suggestions, SuggestChallenge, SuggestBlock)
	}
}

func TestEvaluateUsesOnlyProvidedFeatures(t *testing.T) {
	assessor := New()

	result := assessor.Evaluate(Features{
		Entropy: &EntropyFeature{
			Value:      7.7,
			SampleSize: 800,
			Source:     "body",
		},
	})

	if len(result.Signals) != 1 {
		t.Fatalf("unexpected signal count: got %d want 1", len(result.Signals))
	}

	if result.Signals[0].Key != "entropy" {
		t.Fatalf("unexpected signal key: got %q want %q", result.Signals[0].Key, "entropy")
	}

	if result.Score != 18 {
		t.Fatalf("unexpected score: got %d want 18", result.Score)
	}

	if result.Level != LevelLow {
		t.Fatalf("unexpected level: got %q want %q", result.Level, LevelLow)
	}
}

func TestNilAssessorFallsBackToDefault(t *testing.T) {
	var assessor *Assessor

	result := assessor.Evaluate(Features{
		Fingerprint: &FingerprintFeature{
			Rarity:     0.90,
			Volatility: 0.10,
		},
	})

	if len(result.Signals) != 1 {
		t.Fatalf("unexpected signal count: got %d want 1", len(result.Signals))
	}

	if result.Signals[0].Key != "fingerprint" {
		t.Fatalf("unexpected signal key: got %q", result.Signals[0].Key)
	}
}
