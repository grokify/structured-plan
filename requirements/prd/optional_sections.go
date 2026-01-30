package prd

// AssumptionsConstraints contains assumptions and constraints.
type AssumptionsConstraints struct {
	Assumptions  []Assumption `json:"assumptions"`
	Constraints  []Constraint `json:"constraints"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

// Dependency represents an external dependency.
type Dependency struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type,omitempty"` // API, Service, Team, Vendor
	Owner       string `json:"owner,omitempty"`
	Status      string `json:"status,omitempty"` // Available, Pending, Blocked
	DueDate     string `json:"dueDate,omitempty"`
}

// TechnicalArchitecture contains technical design information.
type TechnicalArchitecture struct {
	Overview          string          `json:"overview"`
	SystemDiagram     string          `json:"systemDiagram,omitempty"` // URL or path to diagram
	DataModel         string          `json:"dataModel,omitempty"`     // URL or path to ERD
	IntegrationPoints []Integration   `json:"integrationPoints,omitempty"`
	TechnologyStack   TechnologyStack `json:"technologyStack,omitempty"`
	SecurityDesign    string          `json:"securityDesign,omitempty"`
	ScalabilityDesign string          `json:"scalabilityDesign,omitempty"`
}

// Integration represents an external integration point.
type Integration struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"` // REST API, GraphQL, Event, Database
	Description   string `json:"description"`
	Protocol      string `json:"protocol,omitempty"`
	AuthMethod    string `json:"authMethod,omitempty"`
	DataFormat    string `json:"dataFormat,omitempty"` // JSON, XML, Protobuf
	RateLimit     string `json:"rateLimit,omitempty"`
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
	DesignPrinciples []string          `json:"designPrinciples,omitempty"`
	Wireframes       []Wireframe       `json:"wireframes,omitempty"`
	InteractionFlows []InteractionFlow `json:"interactionFlows,omitempty"`
	Accessibility    AccessibilitySpec `json:"accessibility,omitempty"`
	BrandGuidelines  string            `json:"brandGuidelines,omitempty"` // URL or path
	DesignSystem     string            `json:"designSystem,omitempty"`    // URL or path
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
	DiagramURL  string   `json:"diagramUrl,omitempty"`
}

// AccessibilitySpec defines accessibility requirements.
type AccessibilitySpec struct {
	Standard        string   `json:"standard"` // WCAG 2.1 AA
	Requirements    []string `json:"requirements,omitempty"`
	TestingApproach string   `json:"testingApproach,omitempty"`
}
