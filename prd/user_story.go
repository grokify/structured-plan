package prd

// Priority represents requirement priority levels.
type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// MoSCoW represents the MoSCoW prioritization method.
type MoSCoW string

const (
	MoSCoWMust   MoSCoW = "must"
	MoSCoWShould MoSCoW = "should"
	MoSCoWCould  MoSCoW = "could"
	MoSCoWWont   MoSCoW = "wont"
)

// UserStory represents a user story with acceptance criteria.
type UserStory struct {
	ID                 string                `json:"id"`
	PersonaID          string                `json:"persona_id"` // Reference to persona
	Title              string                `json:"title"`
	Story              string                `json:"story"` // "As a [persona], I want [goal] so that [reason]"
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptance_criteria"`
	Priority           Priority              `json:"priority"`
	PhaseID            string                `json:"phase_id"` // Reference to roadmap phase
	StoryPoints        *int                  `json:"story_points,omitempty"`
	Dependencies       []string              `json:"dependencies,omitempty"` // Dependent story IDs
	Epic               string                `json:"epic,omitempty"`         // Parent epic
	Labels             []string              `json:"labels,omitempty"`
	Notes              string                `json:"notes,omitempty"`
}

// AcceptanceCriterion defines a testable condition for a user story.
type AcceptanceCriterion struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Given       string `json:"given,omitempty"` // Precondition
	When        string `json:"when,omitempty"`  // Action
	Then        string `json:"then,omitempty"`  // Expected result
}
