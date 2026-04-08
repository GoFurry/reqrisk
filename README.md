# reqrisk

**[中文文档](docs/README_zh.md) | English**

Explainable risk assessment for HTTP request features in Go.

`reqrisk` is a small Go core library that turns precomputed request features into explainable risk results.

It is designed for cases where feature extraction already exists upstream, and a downstream component needs to answer questions like:

- How risky does this request look?
- Which signals contributed to that result?
- What findings should be surfaced to operators or middleware?
- Should the caller observe, review, challenge, or block?

---

## What reqrisk does

`reqrisk` accepts already-derived features such as:

- entropy
- structural complexity
- charset anomalies
- fingerprint rarity or volatility

It then converts them into:

- per-signal analysis
- a total score
- a risk level
- findings with evidence
- lightweight tags
- suggested actions

---

## What reqrisk does not do

`reqrisk` is intentionally narrow.

It does **not**:

- parse raw `*http.Request`
- capture request bodies
- compute entropy from raw payloads
- generate request fingerprints
- act as a WAF
- automatically mutate or block traffic

That makes it a good fit as a downstream core library, including use alongside tools like `web-profiler`.

---

## Design goals

- zero third-party runtime dependencies
- small and stable API
- explainable results
- framework-agnostic design
- not coupled to any single middleware or profiler

---

## Installation

```bash
go get github.com/GoFurry/reqrisk
```

---

## Quick example

A runnable example lives at [`example/main.go`](example/main.go).

The current benchmark baseline is tracked in [`docs/benchmark_baseline.md`](docs/benchmark_baseline.md).

```go
package main

import (
    "fmt"

    "github.com/GoFurry/reqrisk"
)

func main() {
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
        Meta: reqrisk.Meta{
            Target: "body",
            Route:  "/api/comment",
            Method: "POST",
        },
    })

    fmt.Println(result.Score)
    fmt.Println(result.Level)
    fmt.Println(result.Tags)
}
```

---

## Core concepts

### Features

The main input is `Features`.

```go
type Features struct {
    Entropy     *EntropyFeature
    Complexity  *ComplexityFeature
    Charset     *CharsetFeature
    Fingerprint *FingerprintFeature

    Meta Meta
}
```

Each feature is optional, so callers can provide only what they have.

### Result

The main output is `Result`.

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

This keeps scoring explainable instead of opaque.

---

## Built-in signals in V1

- entropy
- complexity
- charset
- fingerprint

Each signal can also be assessed individually:

```go
entropyReport := reqrisk.AssessEntropy(reqrisk.EntropyFeature{
    Value:      7.7,
    SampleSize: 900,
    Source:     "body",
})
```

---

## Risk levels

`reqrisk` maps total score into four default levels:

- `low`
- `medium`
- `high`
- `critical`

Suggested actions are derived from score and findings:

- `observe`
- `review`
- `challenge`
- `block`

---

## Presets

`reqrisk` ships with a few lightweight presets for common scoring styles:

- `PresetBalanced` keeps the default baseline
- `PresetSensitive` lowers thresholds so borderline signals surface earlier
- `PresetConservative` raises thresholds so only stronger signals contribute

Example:

```go
assessor := reqrisk.New(reqrisk.WithPreset(reqrisk.PresetSensitive))
```

Presets are a starting point only. You can still layer explicit policy overrides after them when you need to tune a specific field.

---

## Project scope

V1 focuses on the core scoring model only.

Included:

- unified assessor
- independent signal analysis
- configurable scoring policies
- explainable findings and evidence
- zero third-party runtime dependencies

Deferred:

- HTTP adapters
- web-profiler adapter package
- remote config
- rule DSL
- stateful history or storage backends

---

## Repository layout

- root package: stable public API for callers
- `internal/model`: input and output data structures
- `internal/policy`: scoring policies, options, and config normalization
- `internal/core`: assessment engine and built-in signal implementations

The public import path stays at the root package on purpose so callers keep a simple integration point, while the implementation is split internally by responsibility for long-term maintenance.

---

## Status

`reqrisk` is now implemented as a small, stable core library for explainable request risk assessment.

The first version is intentionally narrow and focused.

---

## License

MIT
