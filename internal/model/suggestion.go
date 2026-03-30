package model

// Suggestion is a lightweight action hint derived from a result.
type Suggestion string

const (
	SuggestObserve   Suggestion = "observe"
	SuggestReview    Suggestion = "review"
	SuggestChallenge Suggestion = "challenge"
	SuggestBlock     Suggestion = "block"
)
