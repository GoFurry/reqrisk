// Package reqrisk provides explainable risk assessment for precomputed request features.
//
// The package is intentionally narrow: it does not parse raw HTTP requests or extract
// low-level features. Instead, it accepts already-derived features and turns them into
// signal reports, an aggregate score, a risk level, findings, tags, and suggested actions.
// The aggregate score is clamped to the configured max score, each signal report carries
// an explainable weighted contribution, and suggestions are lightweight deduplicated hints.
package reqrisk
