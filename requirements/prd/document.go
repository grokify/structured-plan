// Package prd provides data types for structured Product Requirements Documents.
package prd

import (
	"time"

	"github.com/grokify/structured-plan/common"
	"github.com/grokify/structured-plan/goals/okr"
)

// Person is an alias for common.Person for backwards compatibility.
type Person = common.Person

// Approver is an alias for common.Approver for backwards compatibility.
type Approver = common.Approver

// OKR type aliases from goals/okr for backward compatibility.
// These allow existing PRD code to continue using prd.OKR, prd.Objective, etc.
type (
	// OKR represents an Objective with its Key Results.
	OKR = okr.OKR

	// Objective represents a qualitative, inspirational goal.
	Objective = okr.Objective

	// KeyResult represents a measurable outcome that indicates objective achievement.
	KeyResult = okr.KeyResult

	// PhaseTarget represents a Key Result target for a specific roadmap phase.
	PhaseTarget = okr.PhaseTarget
)

// Common type aliases for backward compatibility.
// These allow existing PRD code to continue using prd.Status, prd.Risk, etc.
type (
	// Status represents the document lifecycle status.
	Status = common.Status

	// Risk represents a project risk.
	Risk = common.Risk

	// RiskProbability represents risk probability levels.
	RiskProbability = common.RiskProbability

	// RiskImpact represents risk impact levels.
	RiskImpact = common.RiskImpact

	// RiskStatus represents risk status.
	RiskStatus = common.RiskStatus

	// Assumption represents a condition assumed to be true.
	Assumption = common.Assumption

	// Constraint represents a limitation on the project.
	Constraint = common.Constraint

	// ConstraintType represents types of constraints.
	ConstraintType = common.ConstraintType

	// GlossaryTerm defines a glossary entry.
	GlossaryTerm = common.GlossaryTerm

	// CustomSection allows project-specific sections.
	CustomSection = common.CustomSection

	// OpenItem represents a pending decision or question.
	OpenItem = common.OpenItem

	// Option represents one possible choice for an open item.
	Option = common.Option

	// OpenItemStatus represents the status of an open item.
	OpenItemStatus = common.OpenItemStatus

	// OpenItemResolution documents how an open item was resolved.
	OpenItemResolution = common.OpenItemResolution

	// EffortLevel represents effort estimates.
	EffortLevel = common.EffortLevel

	// RiskLevel represents risk levels for options.
	RiskLevel = common.RiskLevel

	// DecisionRecord documents a completed decision.
	DecisionRecord = common.DecisionRecord

	// DecisionStatus represents the status of a decision.
	DecisionStatus = common.DecisionStatus

	// Priority represents priority levels.
	Priority = common.Priority

	// NonGoal represents an explicit out-of-scope item.
	NonGoal = common.NonGoal
)

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

// Status constants re-exported from common for backward compatibility.
const (
	StatusDraft      = common.StatusDraft
	StatusInReview   = common.StatusInReview
	StatusApproved   = common.StatusApproved
	StatusDeprecated = common.StatusDeprecated
)

// Risk constants re-exported from common for backward compatibility.
const (
	RiskProbabilityLow    = common.RiskProbabilityLow
	RiskProbabilityMedium = common.RiskProbabilityMedium
	RiskProbabilityHigh   = common.RiskProbabilityHigh

	RiskImpactLow      = common.RiskImpactLow
	RiskImpactMedium   = common.RiskImpactMedium
	RiskImpactHigh     = common.RiskImpactHigh
	RiskImpactCritical = common.RiskImpactCritical

	RiskStatusOpen      = common.RiskStatusOpen
	RiskStatusMitigated = common.RiskStatusMitigated
	RiskStatusAccepted  = common.RiskStatusAccepted
	RiskStatusClosed    = common.RiskStatusClosed
)

// Constraint type constants re-exported from common for backward compatibility.
const (
	ConstraintTechnical  = common.ConstraintTechnical
	ConstraintBudget     = common.ConstraintBudget
	ConstraintTimeline   = common.ConstraintTimeline
	ConstraintRegulatory = common.ConstraintRegulatory
	ConstraintResource   = common.ConstraintResource
	ConstraintLegal      = common.ConstraintLegal
)

// OpenItem status constants re-exported from common for backward compatibility.
const (
	OpenItemStatusOpen         = common.OpenItemStatusOpen
	OpenItemStatusInDiscussion = common.OpenItemStatusInDiscussion
	OpenItemStatusBlocked      = common.OpenItemStatusBlocked
	OpenItemStatusResolved     = common.OpenItemStatusResolved
	OpenItemStatusDeferred     = common.OpenItemStatusDeferred
)

// Effort level constants re-exported from common for backward compatibility.
const (
	EffortLow    = common.EffortLow
	EffortMedium = common.EffortMedium
	EffortHigh   = common.EffortHigh
)

// Risk level constants re-exported from common for backward compatibility.
const (
	RiskLevelLow    = common.RiskLevelLow
	RiskLevelMedium = common.RiskLevelMedium
	RiskLevelHigh   = common.RiskLevelHigh
)

// Decision status constants re-exported from common for backward compatibility.
const (
	DecisionProposed   = common.DecisionProposed
	DecisionAccepted   = common.DecisionAccepted
	DecisionSuperseded = common.DecisionSuperseded
	DecisionDeprecated = common.DecisionDeprecated
)

// Priority constants re-exported from common for backward compatibility.
const (
	PriorityCritical = common.PriorityCritical
	PriorityHigh     = common.PriorityHigh
	PriorityMedium   = common.PriorityMedium
	PriorityLow      = common.PriorityLow
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
