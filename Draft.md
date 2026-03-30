# reqrisk Draft

## Overview

`reqrisk` is a small, dependency-free Go core library for turning request-related features into explainable risk results.

It does **not** parse HTTP requests, compute raw request features, or make blocking decisions by itself. Its responsibility is focused and narrow:

**feature -> signal -> score -> explanation**

This makes `reqrisk` suitable as a downstream package for tools like `web-profiler`, while also remaining usable as a standalone library with manually supplied feature inputs.

---

## Design Goals

- Zero third-party runtime dependencies
- Small and stable public API
- Explainable results instead of opaque scores
- Framework-agnostic
- Not coupled to any specific middleware or request profiler
- Suitable for embedding in services, gateways, analyzers, or higher-level middleware

---

## Non-Goals

`reqrisk` does **not** aim to:

- parse raw `*http.Request`
- capture or mutate request bodies
- compute entropy, fingerprint, or structure metrics from raw payloads
- act as a WAF
- automatically block traffic
- manage remote rules or dynamic policy distribution
- persist risk history or maintain stateful scoring in V1

---

## Core Responsibility

The library accepts already-derived request features and transforms them into:

- per-signal analysis
- a total risk score
- a risk level
- findings with human-readable explanations
- lightweight tags
- suggested actions

---

## V1 Architecture

```text
reqrisk/
  doc.go
  assessor.go
  option.go
  config.go

  features.go
  result.go
  finding.go
  signal.go
  level.go
  suggestion.go
  severity.go

  entropy.go
  complexity.go
  charset.go
  fingerprint.go

  internal/
    aggregate/
      aggregate.go
      bands.go
      suggest.go
    mathx/
      clamp.go
      scale.go
```

---

## Public API Strategy

V1 exposes only two main entry styles.

### 1. Full assessment

Used when the caller already has a set of features and wants a full risk result.

```go
assessor := reqrisk.New()

result := assessor.Evaluate(reqrisk.Features{
    Entropy: &reqrisk.EntropyFeature{
        Value:      7.9,
        SampleSize: 1024,
        Source:     "body",
    },
    Complexity: &reqrisk.ComplexityFeature{
        Depth:          7,
        FieldCount:     42,
        MaxArrayLength: 18,
        Source:         "json",
    },
    Charset: &reqrisk.CharsetFeature{
        NonASCIIRatio:  0.82,
        ZeroWidthCount: 3,
        MixedScripts:   true,
        Source:         "body",
    },
    Fingerprint: &reqrisk.FingerprintFeature{
        Rarity:     0.91,
        Volatility: 0.63,
        Source:     "request-shape",
    },
})
```

### 2. Single-signal analysis

Used when the caller only wants one focused analysis.

```go
report := reqrisk.AssessEntropy(reqrisk.EntropyFeature{
    Value:      7.9,
    SampleSize: 1024,
    Source:     "body",
})
```

Equivalent functions exist for:

- `AssessEntropy`
- `AssessComplexity`
- `AssessCharset`
- `AssessFingerprint`

---

## Core Data Model

### Features

`Features` is the main aggregate input type.

All major feature fields are optional and represented as pointers.

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

`Meta` carries lightweight context useful for explanation or display, but does not directly drive scoring in V1.

```go
type Meta struct {
    Target string // body/query/header/path/form
    Route  string // optional
    Method string // optional
}
```

### Feature Types

#### EntropyFeature

```go
type EntropyFeature struct {
    Value      float64
    SampleSize int
    Source     string
}
```

#### ComplexityFeature

```go
type ComplexityFeature struct {
    Depth          int
    FieldCount     int
    MaxArrayLength int
    Source         string
}
```

#### CharsetFeature

```go
type CharsetFeature struct {
    NonASCIIRatio   float64
    ControlRatio    float64
    ZeroWidthCount  int
    MixedScripts    bool
    ReplacementCount int
    Source          string
}
```

#### FingerprintFeature

```go
type FingerprintFeature struct {
    Rarity     float64 // 0~1, higher means rarer
    Volatility float64 // 0~1, higher means less stable
    Source     string
}
```

---

## Result Model

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

## Levels and Suggestions

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

Recommended default bands:

- `0~24` -> `low`
- `25~49` -> `medium`
- `50~74` -> `high`
- `75~100` -> `critical`

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

## Assessor

The core orchestration type is `Assessor`.

```go
type Assessor struct {
    cfg Config
}
```

### Construction

```go
func New(opts ...Option) *Assessor
```

### Evaluation

```go
func (a *Assessor) Evaluate(f Features) Result
```

Responsibilities:

1. run enabled signal analyzers
2. aggregate signal scores
3. compute total level
4. merge tags and findings
5. generate suggestions
6. return final `Result`

---

## Configuration Model

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

### Option

```go
type Option func(*Config)
```

---

## Scoring Strategy

V1 uses weighted heuristic scoring.

### Flow

1. each signal computes a local `SignalReport`
2. each local score is weighted into a contribution
3. contributions are aggregated into a total score
4. total score is clamped into `0~100`
5. score maps to a `Level`
6. suggestions are derived from score and notable findings

This keeps the core model understandable and avoids introducing a rule engine too early.

---

## Analyzer Extensibility

Even in V1, the design should leave room for internal extensibility.

```go
type Analyzer interface {
    Key() string
    Assess(Features) SignalReport
}
```

Default analyzers in V1:

- entropy
- complexity
- charset
- fingerprint

This allows future growth without bloating the public API.

---

## V1 Scope

### Included

- unified `Assessor`
- aggregate `Features`
- four built-in feature types
- four single-signal assessment functions
- explainable `Result`
- configurable scoring policies
- zero third-party runtime dependencies

### Deferred

- HTTP adapters
- `web-profiler` adapter package
- rule DSL
- remote config
- Redis or DB state
- dynamic learning
- plugin ecosystem
- blocking middleware

---

## Positioning

`reqrisk` is not a WAF and not a middleware.

It is a small Go core library for **explainable request risk assessment**, designed to sit downstream of feature extraction and upstream of policy decisions.
