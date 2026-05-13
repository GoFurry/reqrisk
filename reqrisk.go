package reqrisk

import (
	"github.com/gofurry/reqrisk/internal/core"
	"github.com/gofurry/reqrisk/internal/model"
	"github.com/gofurry/reqrisk/internal/policy"
)

type (
	Assessor           = core.Assessor
	Features           = model.Features
	Meta               = model.Meta
	EntropyFeature     = model.EntropyFeature
	ComplexityFeature  = model.ComplexityFeature
	CharsetFeature     = model.CharsetFeature
	FingerprintFeature = model.FingerprintFeature
	Result             = model.Result
	SignalReport       = model.SignalReport
	Finding            = model.Finding
	Evidence           = model.Evidence
	Level              = model.Level
	Severity           = model.Severity
	Suggestion         = model.Suggestion
	Config             = policy.Config
	LevelBands         = policy.LevelBands
	EntropyPolicy      = policy.EntropyPolicy
	ComplexityPolicy   = policy.ComplexityPolicy
	CharsetPolicy      = policy.CharsetPolicy
	FingerprintPolicy  = policy.FingerprintPolicy
	Preset             = policy.Preset
	Option             = policy.Option
)

const (
	PresetBalanced     = policy.PresetBalanced
	PresetSensitive    = policy.PresetSensitive
	PresetConservative = policy.PresetConservative
)

const (
	LevelLow      = model.LevelLow
	LevelMedium   = model.LevelMedium
	LevelHigh     = model.LevelHigh
	LevelCritical = model.LevelCritical

	SeverityInfo     = model.SeverityInfo
	SeverityLow      = model.SeverityLow
	SeverityMedium   = model.SeverityMedium
	SeverityHigh     = model.SeverityHigh
	SeverityCritical = model.SeverityCritical

	SuggestObserve   = model.SuggestObserve
	SuggestReview    = model.SuggestReview
	SuggestChallenge = model.SuggestChallenge
	SuggestBlock     = model.SuggestBlock
)

func New(opts ...Option) *Assessor { return core.New(opts...) }

func AssessEntropy(feature EntropyFeature) SignalReport { return core.AssessEntropy(feature) }

func AssessComplexity(feature ComplexityFeature) SignalReport { return core.AssessComplexity(feature) }

func AssessCharset(feature CharsetFeature) SignalReport { return core.AssessCharset(feature) }

func AssessFingerprint(feature FingerprintFeature) SignalReport {
	return core.AssessFingerprint(feature)
}

func WithBands(b LevelBands) Option { return policy.WithBands(b) }

func WithEntropyPolicy(p EntropyPolicy) Option { return policy.WithEntropyPolicy(p) }

func WithComplexityPolicy(p ComplexityPolicy) Option { return policy.WithComplexityPolicy(p) }

func WithCharsetPolicy(p CharsetPolicy) Option { return policy.WithCharsetPolicy(p) }

func WithFingerprintPolicy(p FingerprintPolicy) Option { return policy.WithFingerprintPolicy(p) }

func WithPreset(p Preset) Option { return policy.WithPreset(p) }

func WithMaxScore(n int) Option { return policy.WithMaxScore(n) }
