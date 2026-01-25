package prd

import "time"

// Roadmap contains the product roadmap with phases.
type Roadmap struct {
	Phases []Phase `json:"phases"`
}

// PhaseType represents the type of roadmap phase.
type PhaseType string

const (
	PhaseTypeGeneric   PhaseType = "generic"   // Phase 1, 2, 3
	PhaseTypeQuarter   PhaseType = "quarter"   // Q1 2026, Q2 2026
	PhaseTypeMonth     PhaseType = "month"     // January 2026
	PhaseTypeSprint    PhaseType = "sprint"    // Sprint 1, Sprint 2
	PhaseTypeMilestone PhaseType = "milestone" // MVP, GA, etc.
)

// Phase represents a roadmap phase.
type Phase struct {
	ID              string        `json:"id"`   // e.g., "phase-1", "q1-2026"
	Name            string        `json:"name"` // e.g., "MVP", "Q1 2026"
	Type            PhaseType     `json:"type"`
	StartDate       *time.Time    `json:"start_date,omitempty"`
	EndDate         *time.Time    `json:"end_date,omitempty"`
	Goals           []string      `json:"goals"`
	Deliverables    []Deliverable `json:"deliverables"`
	SuccessCriteria []string      `json:"success_criteria"`
	Dependencies    []string      `json:"dependencies,omitempty"` // Dependent phase IDs
	Risks           []Risk        `json:"risks,omitempty"`
	Status          PhaseStatus   `json:"status,omitempty"`
	Progress        *int          `json:"progress,omitempty"` // 0-100 percentage
	Notes           string        `json:"notes,omitempty"`
}

// PhaseStatus represents the current status of a phase.
type PhaseStatus string

const (
	PhaseStatusPlanned    PhaseStatus = "planned"
	PhaseStatusInProgress PhaseStatus = "in_progress"
	PhaseStatusCompleted  PhaseStatus = "completed"
	PhaseStatusDelayed    PhaseStatus = "delayed"
	PhaseStatusCancelled  PhaseStatus = "cancelled"
)

// Deliverable represents a phase deliverable.
type Deliverable struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Type        DeliverableType   `json:"type"`
	Status      DeliverableStatus `json:"status,omitempty"`
}

// DeliverableType represents types of deliverables.
type DeliverableType string

const (
	DeliverableFeature        DeliverableType = "feature"
	DeliverableDocumentation  DeliverableType = "documentation"
	DeliverableInfrastructure DeliverableType = "infrastructure"
	DeliverableIntegration    DeliverableType = "integration"
	DeliverableMilestone      DeliverableType = "milestone"
)

// DeliverableStatus represents the status of a deliverable.
type DeliverableStatus string

const (
	DeliverableNotStarted DeliverableStatus = "not_started"
	DeliverableInProgress DeliverableStatus = "in_progress"
	DeliverableCompleted  DeliverableStatus = "completed"
	DeliverableBlocked    DeliverableStatus = "blocked"
)
