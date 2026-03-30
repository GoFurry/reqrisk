package model

// Finding describes one explainable observation within a signal report.
type Finding struct {
	Key          string
	Signal       string
	Severity     Severity
	Contribution int
	Summary      string
	Evidence     []Evidence
}

// Evidence is a small key/value proof attached to a finding.
type Evidence struct {
	Key   string
	Value string
}
