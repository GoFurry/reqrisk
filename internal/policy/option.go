package policy

// Option mutates Config during assessor construction.
type Option func(*Config)

// WithBands overrides aggregate level bands.
func WithBands(b LevelBands) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		cfg.Bands = b
		cfg.bandsCustomized = true
	}
}

// WithEntropyPolicy overrides the entropy policy.
// A fully zero policy disables the signal; otherwise zero-valued fields keep the current defaults.
func WithEntropyPolicy(p EntropyPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		if isZeroEntropyPolicy(p) {
			cfg.Entropy = EntropyPolicy{}
			return
		}
		cfg.Entropy = mergeEntropyPolicy(cfg.Entropy, p)
	}
}

// WithComplexityPolicy overrides the complexity policy.
// A fully zero policy disables the signal; otherwise zero-valued fields keep the current defaults.
func WithComplexityPolicy(p ComplexityPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		if isZeroComplexityPolicy(p) {
			cfg.Complexity = ComplexityPolicy{}
			return
		}
		cfg.Complexity = mergeComplexityPolicy(cfg.Complexity, p)
	}
}

// WithCharsetPolicy overrides the charset policy.
// A fully zero policy disables the signal; otherwise zero-valued fields keep the current defaults.
func WithCharsetPolicy(p CharsetPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		if isZeroCharsetPolicy(p) {
			cfg.Charset = CharsetPolicy{}
			return
		}
		cfg.Charset = mergeCharsetPolicy(cfg.Charset, p)
	}
}

// WithFingerprintPolicy overrides the fingerprint policy.
// A fully zero policy disables the signal; otherwise zero-valued fields keep the current defaults.
func WithFingerprintPolicy(p FingerprintPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		if isZeroFingerprintPolicy(p) {
			cfg.Fingerprint = FingerprintPolicy{}
			return
		}
		cfg.Fingerprint = mergeFingerprintPolicy(cfg.Fingerprint, p)
	}
}

// WithMaxScore overrides the aggregate score ceiling.
func WithMaxScore(n int) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		cfg.MaxScore = n
	}
}

func isZeroEntropyPolicy(p EntropyPolicy) bool {
	return p == (EntropyPolicy{})
}

func mergeEntropyPolicy(base, override EntropyPolicy) EntropyPolicy {
	if override.Weight != 0 {
		base.Weight = override.Weight
	}
	if override.MinSample != 0 {
		base.MinSample = override.MinSample
	}
	if override.Elevated != 0 {
		base.Elevated = override.Elevated
	}
	if override.High != 0 {
		base.High = override.High
	}
	if override.Extreme != 0 {
		base.Extreme = override.Extreme
	}
	return base
}

func isZeroComplexityPolicy(p ComplexityPolicy) bool {
	return p == (ComplexityPolicy{})
}

func mergeComplexityPolicy(base, override ComplexityPolicy) ComplexityPolicy {
	if override.Weight != 0 {
		base.Weight = override.Weight
	}
	if override.DepthHigh != 0 {
		base.DepthHigh = override.DepthHigh
	}
	if override.FieldCountHigh != 0 {
		base.FieldCountHigh = override.FieldCountHigh
	}
	if override.ArrayHigh != 0 {
		base.ArrayHigh = override.ArrayHigh
	}
	return base
}

func isZeroCharsetPolicy(p CharsetPolicy) bool {
	return p == (CharsetPolicy{})
}

func mergeCharsetPolicy(base, override CharsetPolicy) CharsetPolicy {
	if override.Weight != 0 {
		base.Weight = override.Weight
	}
	if override.NonASCIIHigh != 0 {
		base.NonASCIIHigh = override.NonASCIIHigh
	}
	if override.ControlHigh != 0 {
		base.ControlHigh = override.ControlHigh
	}
	if override.ZeroWidthHigh != 0 {
		base.ZeroWidthHigh = override.ZeroWidthHigh
	}
	if override.ReplacementHigh != 0 {
		base.ReplacementHigh = override.ReplacementHigh
	}
	if override.MixedScriptsWeight != 0 {
		base.MixedScriptsWeight = override.MixedScriptsWeight
	}
	return base
}

func isZeroFingerprintPolicy(p FingerprintPolicy) bool {
	return p == (FingerprintPolicy{})
}

func mergeFingerprintPolicy(base, override FingerprintPolicy) FingerprintPolicy {
	if override.Weight != 0 {
		base.Weight = override.Weight
	}
	if override.RarityHigh != 0 {
		base.RarityHigh = override.RarityHigh
	}
	if override.VolatilityHigh != 0 {
		base.VolatilityHigh = override.VolatilityHigh
	}
	return base
}
