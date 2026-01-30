package prd

import "time"

// DecisionsDefinition contains decision records for the PRD.
type DecisionsDefinition struct {
	// Records are the decision records.
	Records []DecisionRecord `json:"records,omitempty"`
}

// DecisionRecord documents a decision made during PRD development.
type DecisionRecord struct {
	// ID is the unique identifier for this decision.
	ID string `json:"id"`

	// Decision is the decision that was made.
	Decision string `json:"decision"`

	// Rationale explains why this decision was made.
	Rationale string `json:"rationale,omitempty"`

	// AlternativesConsidered lists other options that were evaluated.
	AlternativesConsidered []string `json:"alternatives_considered,omitempty"`

	// MadeBy is the person or group who made the decision.
	MadeBy string `json:"made_by,omitempty"`

	// Date is when the decision was made.
	Date time.Time `json:"date,omitempty"`

	// Status is the current status of the decision.
	Status DecisionStatus `json:"status,omitempty"`

	// RelatedIDs are IDs of related items (requirements, risks, etc.).
	RelatedIDs []string `json:"related_ids,omitempty"`
}

// DecisionStatus represents the status of a decision.
type DecisionStatus string

const (
	// DecisionProposed means the decision is proposed but not finalized.
	DecisionProposed DecisionStatus = "proposed"

	// DecisionAccepted means the decision has been accepted.
	DecisionAccepted DecisionStatus = "accepted"

	// DecisionSuperseded means the decision has been replaced by another.
	DecisionSuperseded DecisionStatus = "superseded"

	// DecisionDeprecated means the decision is no longer relevant.
	DecisionDeprecated DecisionStatus = "deprecated"
)
