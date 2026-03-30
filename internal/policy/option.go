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
	}
}

// WithEntropyPolicy overrides the entropy policy.
func WithEntropyPolicy(p EntropyPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		cfg.Entropy = p
	}
}

// WithComplexityPolicy overrides the complexity policy.
func WithComplexityPolicy(p ComplexityPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		cfg.Complexity = p
	}
}

// WithCharsetPolicy overrides the charset policy.
func WithCharsetPolicy(p CharsetPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		cfg.Charset = p
	}
}

// WithFingerprintPolicy overrides the fingerprint policy.
func WithFingerprintPolicy(p FingerprintPolicy) Option {
	return func(cfg *Config) {
		if cfg == nil {
			return
		}
		cfg.Fingerprint = p
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
