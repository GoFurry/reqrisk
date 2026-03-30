package reqrisk

import "testing"

func TestAssessEntropyUsesDefaultThresholds(t *testing.T) {
	report := AssessEntropy(EntropyFeature{
		Value:      7.9,
		SampleSize: 1024,
		Source:     "body",
	})

	if report.Key != "entropy" {
		t.Fatalf("unexpected key: got %q want %q", report.Key, "entropy")
	}

	if report.Score != 90 {
		t.Fatalf("unexpected score: got %d want 90", report.Score)
	}

	if report.Contribution != 23 {
		t.Fatalf("unexpected contribution: got %d want 23", report.Contribution)
	}

	if len(report.Findings) != 1 {
		t.Fatalf("unexpected finding count: got %d want 1", len(report.Findings))
	}

	if report.Findings[0].Severity != SeverityCritical {
		t.Fatalf("unexpected severity: got %q want %q", report.Findings[0].Severity, SeverityCritical)
	}
}

func TestAssessEntropyIgnoresSmallSamples(t *testing.T) {
	report := AssessEntropy(EntropyFeature{
		Value:      7.9,
		SampleSize: 64,
		Source:     "body",
	})

	if report.Score != 0 || report.Contribution != 0 {
		t.Fatalf("unexpected report for small sample: %+v", report)
	}

	if len(report.Findings) != 0 {
		t.Fatalf("expected no findings for small sample, got %+v", report.Findings)
	}
}

func TestAssessComplexityProducesMultipleFindings(t *testing.T) {
	report := AssessComplexity(ComplexityFeature{
		Depth:          7,
		FieldCount:     80,
		MaxArrayLength: 25,
		Source:         "json",
	})

	if report.Score != 95 {
		t.Fatalf("unexpected score: got %d want 95", report.Score)
	}

	if report.Contribution != 24 {
		t.Fatalf("unexpected contribution: got %d want 24", report.Contribution)
	}

	if len(report.Findings) != 3 {
		t.Fatalf("unexpected finding count: got %d want 3", len(report.Findings))
	}
}

func TestAssessCharsetProducesExplainableTags(t *testing.T) {
	report := AssessCharset(CharsetFeature{
		NonASCIIRatio:    0.91,
		ControlRatio:     0.01,
		ZeroWidthCount:   4,
		MixedScripts:     true,
		ReplacementCount: 0,
		Source:           "body",
	})

	if report.Score != 70 {
		t.Fatalf("unexpected score: got %d want 70", report.Score)
	}

	if len(report.Tags) < 3 {
		t.Fatalf("expected multiple tags, got %v", report.Tags)
	}
}

func TestAssessFingerprintUsesRarityAndVolatility(t *testing.T) {
	report := AssessFingerprint(FingerprintFeature{
		Rarity:     0.97,
		Volatility: 0.61,
		Source:     "request-shape",
	})

	if report.Score != 75 {
		t.Fatalf("unexpected score: got %d want 75", report.Score)
	}

	if len(report.Findings) != 2 {
		t.Fatalf("unexpected finding count: got %d want 2", len(report.Findings))
	}

	if report.Findings[0].Signal != "fingerprint" || report.Findings[1].Signal != "fingerprint" {
		t.Fatalf("unexpected finding signals: %+v", report.Findings)
	}
}
