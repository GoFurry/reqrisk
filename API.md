# reqrisk API Draft

## Package

```go
package reqrisk
```

---

## Core Constructor

```go
func New(opts ...Option) *Assessor
```

Creates a new assessor with default configuration, optionally overridden by options.

---

## Core Types

### Assessor

```go
type Assessor struct {
    // contains filtered or unexported fields
}
```

#### Methods

```go
func (a *Assessor) Evaluate(f Features) Result
```

Evaluates the provided features and returns a complete explainable risk result.

---

## Aggregate Input

### Features

```go
type Features struct {
    Entropy     *EntropyFeature
    Complexity  *ComplexityFeature
    Charset     *CharsetFeature
    Fingerprint *FingerprintFeature

    Meta Meta
}
```

### Meta

```go
type Meta struct {
    Target string
    Route  string
    Method string
}
```

---

## Feature Types

### EntropyFeature

```go
type EntropyFeature struct {
    Value      float64
    SampleSize int
    Source     string
}
```

### ComplexityFeature

```go
type ComplexityFeature struct {
    Depth          int
    FieldCount     int
    MaxArrayLength int
    Source         string
}
```

### CharsetFeature

```go
type CharsetFeature struct {
    NonASCIIRatio    float64
    ControlRatio     float64
    ZeroWidthCount   int
    MixedScripts     bool
    ReplacementCount int
    Source           string
}
```

### FingerprintFeature

```go
type FingerprintFeature struct {
    Rarity     float64
    Volatility float64
    Source     string
}
```

---

## Single-Signal Assessment

### Entropy

```go
func AssessEntropy(f EntropyFeature) SignalReport
```

### Complexity

```go
func AssessComplexity(f ComplexityFeature) SignalReport
```

### Charset

```go
func AssessCharset(f CharsetFeature) SignalReport
```

### Fingerprint

```go
func AssessFingerprint(f FingerprintFeature) SignalReport
```

These functions provide focused analysis for callers that do not need full aggregation.

---

## Result Types

### Result

```go
type Result struct {
    Score       int
    Level       Level
    Signals     []SignalReport
    Findings    []Finding
    Tags        []string
    Suggestions []Suggestion
}
```

### SignalReport

```go
type SignalReport struct {
    Key          string
    Score        int
    Contribution int
    Findings     []Finding
    Tags         []string
}
```

### Finding

```go
type Finding struct {
    Key          string
    Signal       string
    Severity     Severity
    Contribution int
    Summary      string
    Evidence     []Evidence
}
```

### Evidence

```go
type Evidence struct {
    Key   string
    Value string
}
```

---

## Level and Severity

### Level

```go
type Level string

const (
    LevelLow      Level = "low"
    LevelMedium   Level = "medium"
    LevelHigh     Level = "high"
    LevelCritical Level = "critical"
)
```

### Severity

```go
type Severity string

const (
    SeverityInfo     Severity = "info"
    SeverityLow      Severity = "low"
    SeverityMedium   Severity = "medium"
    SeverityHigh     Severity = "high"
    SeverityCritical Severity = "critical"
)
```

---

## Suggestions

### Suggestion

```go
type Suggestion string

const (
    SuggestObserve   Suggestion = "observe"
    SuggestReview    Suggestion = "review"
    SuggestChallenge Suggestion = "challenge"
    SuggestBlock     Suggestion = "block"
)
```

---

## Configuration

### Config

```go
type Config struct {
    MaxScore int

    Bands LevelBands

    Entropy     EntropyPolicy
    Complexity  ComplexityPolicy
    Charset     CharsetPolicy
    Fingerprint FingerprintPolicy
}
```

### LevelBands

```go
type LevelBands struct {
    Medium   int
    High     int
    Critical int
}
```

### EntropyPolicy

```go
type EntropyPolicy struct {
    Weight    float64
    MinSample int
    Elevated  float64
    High      float64
    Extreme   float64
}
```

### ComplexityPolicy

```go
type ComplexityPolicy struct {
    Weight         float64
    DepthHigh      int
    FieldCountHigh int
    ArrayHigh      int
}
```

### CharsetPolicy

```go
type CharsetPolicy struct {
    Weight             float64
    NonASCIIHigh       float64
    ControlHigh        float64
    ZeroWidthHigh      int
    ReplacementHigh    int
    MixedScriptsWeight float64
}
```

### FingerprintPolicy

```go
type FingerprintPolicy struct {
    Weight         float64
    RarityHigh     float64
    VolatilityHigh float64
}
```

---

## Options

```go
type Option func(*Config)
```

Planned examples:

```go
func WithBands(b LevelBands) Option
func WithEntropyPolicy(p EntropyPolicy) Option
func WithComplexityPolicy(p ComplexityPolicy) Option
func WithCharsetPolicy(p CharsetPolicy) Option
func WithFingerprintPolicy(p FingerprintPolicy) Option
func WithMaxScore(n int) Option
```

---

## Internal Extension Contract

This contract may remain internal in V1, but the architecture assumes analyzer-style composition.

```go
type Analyzer interface {
    Key() string
    Assess(Features) SignalReport
}
```

---

## Example

```go
assessor := reqrisk.New()

result := assessor.Evaluate(reqrisk.Features{
    Entropy: &reqrisk.EntropyFeature{
        Value:      7.8,
        SampleSize: 1200,
        Source:     "body",
    },
    Charset: &reqrisk.CharsetFeature{
        NonASCIIRatio:  0.91,
        ZeroWidthCount: 2,
        MixedScripts:   true,
        Source:         "body",
    },
    Meta: reqrisk.Meta{
        Target: "body",
        Route:  "/api/comment",
        Method: "POST",
    },
})
```

---

## Versioning Notes

The V1 API should favor:

- stable core types
- explicit fields over hidden magic
- predictable scoring behavior
- explainable outputs over compact but opaque abstractions
