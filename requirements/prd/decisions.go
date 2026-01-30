package prd

// DecisionsDefinition contains decision records for the PRD.
type DecisionsDefinition struct {
	// Records are the decision records.
	Records []DecisionRecord `json:"records,omitempty"`
}

// Note: DecisionRecord, DecisionStatus types are now defined in common/
// and aliased in document.go.
