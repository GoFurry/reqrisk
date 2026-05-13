# Benchmark Baseline

This document captures a lightweight benchmark baseline for `reqrisk`.
It is a reference point for future policy or analyzer changes, not a universal performance promise.

## Environment

- Date: `2026-04-08`
- OS: `windows`
- Arch: `amd64`
- CPU: `AMD Ryzen 7 5800H with Radeon Graphics`
- Package: `github.com/gofurry/reqrisk`

## Commands

```bash
go test ./...
go test -run ^$ -bench . ./...
```

## Test Status

- `go test ./...`: passed

## Benchmark Results

| Benchmark | Result |
| --- | --- |
| `BenchmarkEvaluateFull-16` | `9668 ns/op`, `8704 B/op`, `92 allocs/op` |
| `BenchmarkAssessEntropy-16` | `587.2 ns/op`, `400 B/op`, `11 allocs/op` |
| `BenchmarkAssessComplexity-16` | `1373 ns/op`, `1232 B/op`, `19 allocs/op` |
| `BenchmarkAssessCharset-16` | `1920 ns/op`, `2029 B/op`, `31 allocs/op` |
| `BenchmarkAssessFingerprint-16` | `974.9 ns/op`, `845 B/op`, `19 allocs/op` |

## Notes

- `Evaluate` remains the heaviest public entry point, which is expected because it runs all enabled analyzers and assembles the aggregate result.
- Single-signal helpers are small and stable enough to serve as cheap building blocks for callers that already have feature extraction upstream.
- Re-run this baseline whenever scoring rules, preset profiles, or signal allocation logic changes.
