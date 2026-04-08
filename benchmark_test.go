package reqrisk

import "testing"

var (
	benchmarkResult       Result
	benchmarkSignalReport SignalReport
)

func BenchmarkEvaluateFull(b *testing.B) {
	assessor := New()
	features := benchmarkFeatures()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkResult = assessor.Evaluate(features)
	}
}

func BenchmarkAssessEntropy(b *testing.B) {
	feature := EntropyFeature{
		Value:      7.9,
		SampleSize: 1024,
		Source:     "body",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSignalReport = AssessEntropy(feature)
	}
}

func BenchmarkAssessComplexity(b *testing.B) {
	feature := ComplexityFeature{
		Depth:          8,
		FieldCount:     96,
		MaxArrayLength: 42,
		Source:         "json",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSignalReport = AssessComplexity(feature)
	}
}

func BenchmarkAssessCharset(b *testing.B) {
	feature := CharsetFeature{
		NonASCIIRatio:    0.92,
		ControlRatio:     0.04,
		ZeroWidthCount:   3,
		MixedScripts:     true,
		ReplacementCount: 2,
		Source:           "body",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSignalReport = AssessCharset(feature)
	}
}

func BenchmarkAssessFingerprint(b *testing.B) {
	feature := FingerprintFeature{
		Rarity:     0.96,
		Volatility: 0.85,
		Source:     "request-shape",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSignalReport = AssessFingerprint(feature)
	}
}

func benchmarkFeatures() Features {
	return Features{
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
	}
}
