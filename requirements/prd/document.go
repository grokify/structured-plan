// Package prd provides data types for structured Product Requirements Documents.
package prd

import (
	"time"

	"github.com/grokify/structured-plan/common"
)

// Person is an alias for common.Person for backwards compatibility.
type Person = common.Person

// Approver is an alias for common.Approver for backwards compatibility.
type Approver = common.Approver

// Document represents a complete Product Requirements Document.
type Document struct {
	Metadata         Metadata         `json:"metadata"`
	ExecutiveSummary ExecutiveSummary `json:"executive_summary"`
	Objectives       Objectives       `json:"objectives"`
	Personas         []Persona        `json:"personas"`
	UserStories      []UserStory      `json:"user_stories"`
	Requirements     Requirements     `json:"requirements"`
	Roadmap          Roadmap          `json:"roadmap"`

	// Optional sections
	Assumptions      *AssumptionsConstraints `json:"assumptions,omitempty"`
	OutOfScope       []string                `json:"out_of_scope,omitempty"`
	TechArchitecture *TechnicalArchitecture  `json:"technical_architecture,omitempty"`
	UXRequirements   *UXRequirements         `json:"ux_requirements,omitempty"`
	Risks            []Risk                  `json:"risks,omitempty"`
	Glossary         []GlossaryTerm          `json:"glossary,omitempty"`

	// Custom sections for project-specific needs
	CustomSections []CustomSection `json:"custom_sections,omitempty"`

	// Extended sections (from agent-team-prd merge)
	// These provide additional structure for problem definition, market analysis,
	// solution evaluation, decision tracking, and quality reviews.

	// Problem provides detailed problem definition with evidence.
	Problem *ProblemDefinition `json:"problem,omitempty"`

	// Market contains market analysis and competitive landscape.
	Market *MarketDefinition `json:"market,omitempty"`

	// Solution contains solution options and selection rationale.
	Solution *SolutionDefinition `json:"solution,omitempty"`

	// Decisions contains decision records for the PRD.
	Decisions *DecisionsDefinition `json:"decisions,omitempty"`

	// OpenItems contains pending decisions that need resolution.
	OpenItems []OpenItem `json:"open_items,omitempty"`

	// Reviews contains review outcomes and quality assessments.
	Reviews *ReviewsDefinition `json:"reviews,omitempty"`

	// RevisionHistory tracks changes to the PRD over time.
	RevisionHistory []RevisionRecord `json:"revision_history,omitempty"`

	// Goals contains alignment with strategic goals (V2MOM, OKR).
	Goals *GoalsAlignment `json:"goals,omitempty"`

	// CurrentState documents the existing state before the proposed solution.
	CurrentState *CurrentState `json:"current_state,omitempty"`

	// SecurityModel documents security architecture and threat model.
	// This section is strongly recommended for all PRDs.
	SecurityModel *SecurityModel `json:"security_model,omitempty"`

	// Appendices contains supplementary information and domain-specific data.
	Appendices []Appendix `json:"appendices,omitempty"`
}

// Status represents the document lifecycle status.
type Status string

const (
	StatusDraft      Status = "draft"
	StatusInReview   Status = "in_review"
	StatusApproved   Status = "approved"
	StatusDeprecated Status = "deprecated"
)

// Metadata contains document metadata.
type Metadata struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Version   string     `json:"version"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Authors   []Person   `json:"authors"`
	Reviewers []Person   `json:"reviewers,omitempty"`
	Approvers []Approver `json:"approvers,omitempty"`
	Tags      []string   `json:"tags,omitempty"`

	// SemanticVersioning indicates the Version field follows Semantic Versioning (semver.org).
	SemanticVersioning bool `json:"semantic_versioning,omitempty"`
}

// ExecutiveSummary provides high-level product overview.
type ExecutiveSummary struct {
	ProblemStatement string   `json:"problem_statement"`
	ProposedSolution string   `json:"proposed_solution"`
	ExpectedOutcomes []string `json:"expected_outcomes"`
	TargetAudience   string   `json:"target_audience,omitempty"`
	ValueProposition string   `json:"value_proposition,omitempty"`
}

// Objectives defines business and product goals using OKR structure.
type Objectives struct {
	// OKRs contains Objectives and Key Results in nested OKR format.
	OKRs []OKR `json:"okrs"`
}

// OKR represents an Objective with its Key Results.
// Following the OKR methodology used at Google, Intel, Intuit, and others.
type OKR struct {
	Objective  Objective   `json:"objective"`
	KeyResults []KeyResult `json:"key_results"` // Must have 1+ Key Results
}

// Objective represents a qualitative, inspirational goal.
type Objective struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Rationale   string   `json:"rationale,omitempty"`
	AlignedWith string   `json:"aligned_with,omitempty"` // Parent strategy/OKR
	Category    string   `json:"category,omitempty"`     // Business, Product, Team, etc.
	Owner       string   `json:"owner,omitempty"`        // Person or team responsible
	Timeframe   string   `json:"timeframe,omitempty"`    // e.g., "Q1 2026", "H1 2026", "FY2026"
	Tags        []string `json:"tags,omitempty"`         // For filtering by topic/domain
}

// KeyResult represents a measurable outcome that indicates objective achievement.
type KeyResult struct {
	ID                string        `json:"id"`
	Description       string        `json:"description"`
	Metric            string        `json:"metric"`                       // What is measured (e.g., "Monthly Active Users")
	Baseline          string        `json:"baseline,omitempty"`           // Starting value
	Target            string        `json:"target"`                       // Target value to achieve
	Current           string        `json:"current,omitempty"`            // Current value (for tracking)
	Unit              string        `json:"unit,omitempty"`               // e.g., "%", "users", "$"
	MeasurementMethod string        `json:"measurement_method,omitempty"` // How it's measured
	Owner             string        `json:"owner,omitempty"`              // Person or team responsible
	Confidence        float64       `json:"confidence,omitempty"`         // 0.0-1.0 confidence score
	PhaseTargets      []PhaseTarget `json:"phase_targets,omitempty"`      // Per-phase targets for roadmap alignment
	Tags              []string      `json:"tags,omitempty"`               // For filtering by topic/domain
}

// PhaseTarget represents a Key Result target for a specific roadmap phase.
type PhaseTarget struct {
	PhaseID string `json:"phase_id"`         // Reference to roadmap phase
	Target  string `json:"target"`           // Target value for this phase
	Status  string `json:"status,omitempty"` // not_started, in_progress, achieved, missed
	Actual  string `json:"actual,omitempty"` // Actual value achieved
	Notes   string `json:"notes,omitempty"`  // Commentary on progress
}
