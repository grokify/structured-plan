package common

// ConstraintType represents types of constraints.
type ConstraintType string

const (
	ConstraintTechnical  ConstraintType = "technical"
	ConstraintBudget     ConstraintType = "budget"
	ConstraintTimeline   ConstraintType = "timeline"
	ConstraintRegulatory ConstraintType = "regulatory"
	ConstraintResource   ConstraintType = "resource"
	ConstraintLegal      ConstraintType = "legal"
)

// Constraint represents a limitation on the project.
// Used across PRD and TRD documents.
type Constraint struct {
	ID          string         `json:"id"`
	Type        ConstraintType `json:"type"`
	Description string         `json:"description"`
	Impact      string         `json:"impact,omitempty"`
	Mitigation  string         `json:"mitigation,omitempty"`
	Rationale   string         `json:"rationale,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
}
