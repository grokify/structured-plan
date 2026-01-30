package prd

// AssumptionsConstraints contains assumptions and constraints.
type AssumptionsConstraints struct {
	Assumptions  []Assumption `json:"assumptions"`
	Constraints  []Constraint `json:"constraints"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

// Assumption represents a condition assumed to be true.
type Assumption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Rationale   string `json:"rationale,omitempty"`
	Risk        string `json:"risk,omitempty"` // What happens if assumption is wrong
	Validated   bool   `json:"validated,omitempty"`
}

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
type Constraint struct {
	ID          string         `json:"id"`
	Type        ConstraintType `json:"type"`
	Description string         `json:"description"`
	Impact      string         `json:"impact,omitempty"`
	Mitigation  string         `json:"mitigation,omitempty"`
}

// Dependency represents an external dependency.
type Dependency struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type,omitempty"` // API, Service, Team, Vendor
	Owner       string `json:"owner,omitempty"`
	Status      string `json:"status,omitempty"` // Available, Pending, Blocked
	DueDate     string `json:"due_date,omitempty"`
}

// TechnicalArchitecture contains technical design information.
type TechnicalArchitecture struct {
	Overview          string          `json:"overview"`
	SystemDiagram     string          `json:"system_diagram,omitempty"` // URL or path to diagram
	DataModel         string          `json:"data_model,omitempty"`     // URL or path to ERD
	IntegrationPoints []Integration   `json:"integration_points,omitempty"`
	TechnologyStack   TechnologyStack `json:"technology_stack,omitempty"`
	SecurityDesign    string          `json:"security_design,omitempty"`
	ScalabilityDesign string          `json:"scalability_design,omitempty"`
}

// Integration represents an external integration point.
type Integration struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"` // REST API, GraphQL, Event, Database
	Description   string `json:"description"`
	Protocol      string `json:"protocol,omitempty"`
	AuthMethod    string `json:"auth_method,omitempty"`
	DataFormat    string `json:"data_format,omitempty"` // JSON, XML, Protobuf
	RateLimit     string `json:"rate_limit,omitempty"`
	Documentation string `json:"documentation,omitempty"` // URL to docs
}

// TechnologyStack defines the technology choices.
type TechnologyStack struct {
	Frontend       []Technology `json:"frontend,omitempty"`
	Backend        []Technology `json:"backend,omitempty"`
	Database       []Technology `json:"database,omitempty"`
	Infrastructure []Technology `json:"infrastructure,omitempty"`
	DevOps         []Technology `json:"devops,omitempty"`
	Monitoring     []Technology `json:"monitoring,omitempty"`
}

// Technology represents a technology choice.
type Technology struct {
	Name         string   `json:"name"`
	Version      string   `json:"version,omitempty"`
	Purpose      string   `json:"purpose,omitempty"`
	Rationale    string   `json:"rationale,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"` // Considered alternatives
}

// UXRequirements contains UX/UI requirements.
type UXRequirements struct {
	DesignPrinciples []string          `json:"design_principles,omitempty"`
	Wireframes       []Wireframe       `json:"wireframes,omitempty"`
	InteractionFlows []InteractionFlow `json:"interaction_flows,omitempty"`
	Accessibility    AccessibilitySpec `json:"accessibility,omitempty"`
	BrandGuidelines  string            `json:"brand_guidelines,omitempty"` // URL or path
	DesignSystem     string            `json:"design_system,omitempty"`    // URL or path
}

// Wireframe represents a wireframe or mockup.
type Wireframe struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`              // Link to wireframe
	Status      string `json:"status,omitempty"` // Draft, Approved
}

// InteractionFlow represents a user interaction flow.
type InteractionFlow struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	DiagramURL  string   `json:"diagram_url,omitempty"`
}

// AccessibilitySpec defines accessibility requirements.
type AccessibilitySpec struct {
	Standard        string   `json:"standard"` // WCAG 2.1 AA
	Requirements    []string `json:"requirements,omitempty"`
	TestingApproach string   `json:"testing_approach,omitempty"`
}

// RiskProbability represents risk probability levels.
type RiskProbability string

const (
	RiskProbabilityLow    RiskProbability = "low"
	RiskProbabilityMedium RiskProbability = "medium"
	RiskProbabilityHigh   RiskProbability = "high"
)

// RiskImpact represents risk impact levels.
type RiskImpact string

const (
	RiskImpactLow      RiskImpact = "low"
	RiskImpactMedium   RiskImpact = "medium"
	RiskImpactHigh     RiskImpact = "high"
	RiskImpactCritical RiskImpact = "critical"
)

// RiskStatus represents risk status.
type RiskStatus string

const (
	RiskStatusOpen      RiskStatus = "open"
	RiskStatusMitigated RiskStatus = "mitigated"
	RiskStatusAccepted  RiskStatus = "accepted"
	RiskStatusClosed    RiskStatus = "closed"
)

// Risk represents a project risk.
type Risk struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Probability RiskProbability `json:"probability"`
	Impact      RiskImpact      `json:"impact"`
	Mitigation  string          `json:"mitigation"`
	Owner       string          `json:"owner,omitempty"`
	Status      RiskStatus      `json:"status,omitempty"`
	DueDate     string          `json:"due_date,omitempty"`
	Tags        []string        `json:"tags,omitempty"` // For filtering by topic/domain
	Notes       string          `json:"notes,omitempty"`

	// AppendixRefs references appendices with additional details for this risk.
	AppendixRefs []string `json:"appendix_refs,omitempty"`
}

// GlossaryTerm defines a glossary entry.
type GlossaryTerm struct {
	Term       string   `json:"term"`
	Definition string   `json:"definition"`
	Acronym    string   `json:"acronym,omitempty"`
	Context    string   `json:"context,omitempty"`
	Related    []string `json:"related,omitempty"` // Related terms
}

// CustomSection allows project-specific sections.
type CustomSection struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Content     any    `json:"content"`          // Flexible content structure
	Schema      string `json:"schema,omitempty"` // Optional JSON schema for validation
}
