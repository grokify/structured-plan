// Package prd provides data types for structured Product Requirements Documents.
package prd

import "time"

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

	// Reviews contains review outcomes and quality assessments.
	Reviews *ReviewsDefinition `json:"reviews,omitempty"`

	// RevisionHistory tracks changes to the PRD over time.
	RevisionHistory []RevisionRecord `json:"revision_history,omitempty"`

	// Goals contains alignment with strategic goals (V2MOM, OKR).
	Goals *GoalsAlignment `json:"goals,omitempty"`
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
}

// Person represents an individual contributor.
type Person struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
}

// Approver represents a person with approval authority.
type Approver struct {
	Person
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	Approved   bool       `json:"approved"`
	Comments   string     `json:"comments,omitempty"`
}

// ExecutiveSummary provides high-level product overview.
type ExecutiveSummary struct {
	ProblemStatement string   `json:"problem_statement"`
	ProposedSolution string   `json:"proposed_solution"`
	ExpectedOutcomes []string `json:"expected_outcomes"`
	TargetAudience   string   `json:"target_audience,omitempty"`
	ValueProposition string   `json:"value_proposition,omitempty"`
}

// Objectives defines business and product goals.
type Objectives struct {
	BusinessObjectives []Objective     `json:"business_objectives"`
	ProductGoals       []Objective     `json:"product_goals"`
	SuccessMetrics     []SuccessMetric `json:"success_metrics"`
}

// Objective represents a single business or product objective.
type Objective struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Rationale   string `json:"rationale,omitempty"`
	AlignedWith string `json:"aligned_with,omitempty"` // Parent strategy/OKR
}

// SuccessMetric defines how success is measured.
type SuccessMetric struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Metric            string `json:"metric"` // What is measured
	Target            string `json:"target"` // Target value
	CurrentBaseline   string `json:"current_baseline,omitempty"`
	MeasurementMethod string `json:"measurement_method,omitempty"`
}
