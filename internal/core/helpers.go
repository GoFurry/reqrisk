package core

import (
	"strconv"
	"strings"

	"github.com/gofurry/reqrisk/internal/model"
)

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func dedupeStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
}

func formatInt(value int) string {
	return strconv.Itoa(value)
}

func appendSourceEvidence(evidence []model.Evidence, source string) []model.Evidence {
	source = strings.TrimSpace(source)
	if source == "" {
		return evidence
	}
	return append(evidence, model.Evidence{Key: "source", Value: source})
}

func severityRank(severity model.Severity) int {
	switch severity {
	case model.SeverityCritical:
		return 5
	case model.SeverityHigh:
		return 4
	case model.SeverityMedium:
		return 3
	case model.SeverityLow:
		return 2
	case model.SeverityInfo:
		return 1
	default:
		return 0
	}
}
