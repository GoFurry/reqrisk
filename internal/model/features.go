package model

// EntropyFeature is an input signal for entropy-based risk assessment.
type EntropyFeature struct {
	Value      float64
	SampleSize int
	Source     string
}

// ComplexityFeature is an input signal for structure-based risk assessment.
type ComplexityFeature struct {
	Depth          int
	FieldCount     int
	MaxArrayLength int
	Source         string
}

// CharsetFeature is an input signal for text anomaly assessment.
type CharsetFeature struct {
	NonASCIIRatio    float64
	ControlRatio     float64
	ZeroWidthCount   int
	MixedScripts     bool
	ReplacementCount int
	Source           string
}

// FingerprintFeature is an input signal for fingerprint rarity and volatility assessment.
type FingerprintFeature struct {
	Rarity     float64
	Volatility float64
	Source     string
}

// Features is the aggregate input for full request risk assessment.
type Features struct {
	Entropy     *EntropyFeature
	Complexity  *ComplexityFeature
	Charset     *CharsetFeature
	Fingerprint *FingerprintFeature

	Meta Meta
}

// Meta carries lightweight caller context that may be useful for display or explanation.
type Meta struct {
	Target string
	Route  string
	Method string
}
