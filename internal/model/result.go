package model

// Result is the aggregate output returned by Assessor.Evaluate.
type Result struct {
	Score       int
	Level       Level
	Signals     []SignalReport
	Findings    []Finding
	Tags        []string
	Suggestions []Suggestion
}

// SignalReport is the explainable result for a single signal.
type SignalReport struct {
	Key          string
	Score        int
	Contribution int
	Findings     []Finding
	Tags         []string
}
