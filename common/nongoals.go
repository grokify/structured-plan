package common

// NonGoal represents an explicit out-of-scope item with rationale.
// More structured than a simple []string for OutOfScope.
type NonGoal struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Rationale   string   `json:"rationale,omitempty"`   // Why it's out of scope
	FuturePhase string   `json:"futurePhase,omitempty"` // "Phase 2", "v2.0", etc.
	Tags        []string `json:"tags,omitempty"`
}
