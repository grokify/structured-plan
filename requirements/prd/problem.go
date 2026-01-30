package prd

// ProblemDefinition contains the problem statement with evidence.
type ProblemDefinition struct {
	// ID is the unique identifier for this problem.
	ID string `json:"id,omitempty"`

	// Statement is the problem statement.
	Statement string `json:"statement"`

	// UserImpact describes how users are affected by this problem.
	UserImpact string `json:"userImpact,omitempty"`

	// Evidence supports the existence and severity of the problem.
	Evidence []Evidence `json:"evidence,omitempty"`

	// Confidence is the confidence level in the problem definition (0.0-1.0).
	Confidence float64 `json:"confidence,omitempty"`

	// RootCauses are the underlying causes of the problem.
	RootCauses []string `json:"rootCauses,omitempty"`

	// AffectedSegments are user segments affected by this problem.
	AffectedSegments []string `json:"affectedSegments,omitempty"`

	// SecondaryProblems are related or secondary problems.
	SecondaryProblems []ProblemDefinition `json:"secondaryProblems,omitempty"`
}

// Evidence supports a problem statement or claim.
type Evidence struct {
	// Type categorizes the evidence source.
	Type EvidenceType `json:"type"`

	// Source identifies where the evidence came from.
	Source string `json:"source"`

	// Summary describes what the evidence shows.
	Summary string `json:"summary,omitempty"`

	// SampleSize is the number of data points (for quantitative evidence).
	SampleSize int `json:"sampleSize,omitempty"`

	// Strength indicates how strong the evidence is.
	Strength EvidenceStrength `json:"strength,omitempty"`

	// Date is when the evidence was collected.
	Date string `json:"date,omitempty"`
}

// EvidenceType categorizes evidence sources.
type EvidenceType string

const (
	// EvidenceInterview is from user interviews.
	EvidenceInterview EvidenceType = "interview"

	// EvidenceSurvey is from user surveys.
	EvidenceSurvey EvidenceType = "survey"

	// EvidenceAnalytics is from product analytics.
	EvidenceAnalytics EvidenceType = "analytics"

	// EvidenceSupportTicket is from support tickets.
	EvidenceSupportTicket EvidenceType = "support_ticket"

	// EvidenceMarketResearch is from market research.
	EvidenceMarketResearch EvidenceType = "market_research"

	// EvidenceAssumption is an assumption (not validated).
	EvidenceAssumption EvidenceType = "assumption"
)

// EvidenceStrength indicates how strong evidence is.
type EvidenceStrength string

const (
	// StrengthLow indicates weak evidence.
	StrengthLow EvidenceStrength = "low"

	// StrengthMedium indicates moderate evidence.
	StrengthMedium EvidenceStrength = "medium"

	// StrengthHigh indicates strong evidence.
	StrengthHigh EvidenceStrength = "high"
)
