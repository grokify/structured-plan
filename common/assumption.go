package common

// Assumption represents a condition assumed to be true.
// Used across PRD, MRD, and TRD documents.
type Assumption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Rationale   string `json:"rationale,omitempty"`
	Risk        string `json:"risk,omitempty"` // What happens if assumption is wrong
	Validated   bool   `json:"validated,omitempty"`
}
