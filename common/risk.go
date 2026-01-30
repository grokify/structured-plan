package common

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
// Used across PRD, MRD, and TRD documents.
type Risk struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Probability RiskProbability `json:"probability"`
	Impact      RiskImpact      `json:"impact"`
	Mitigation  string          `json:"mitigation"`
	Owner       string          `json:"owner,omitempty"`
	Status      RiskStatus      `json:"status,omitempty"`
	Category    string          `json:"category,omitempty"` // Market, Competitive, Technical, etc.
	DueDate     string          `json:"dueDate,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Notes       string          `json:"notes,omitempty"`

	// AppendixRefs references appendices with additional details for this risk.
	AppendixRefs []string `json:"appendixRefs,omitempty"`
}
