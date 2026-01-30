package prd

// ReviewsDefinition contains review outcomes and quality assessments.
type ReviewsDefinition struct {
	// ReviewBoardSummary is a summary from the review board.
	ReviewBoardSummary string `json:"reviewBoardSummary,omitempty"`

	// QualityScores contains scores across quality dimensions.
	QualityScores *QualityScores `json:"qualityScores,omitempty"`

	// Decision is the review decision.
	Decision ReviewDecision `json:"decision,omitempty"`

	// Blockers are issues that block approval.
	Blockers []Blocker `json:"blockers,omitempty"`

	// RevisionTriggers are issues requiring revision.
	RevisionTriggers []RevisionTrigger `json:"revisionTriggers,omitempty"`
}

// QualityScores contains scores across the 10 quality dimensions.
type QualityScores struct {
	ProblemDefinition    float64 `json:"problemDefinition"`
	UserUnderstanding    float64 `json:"userUnderstanding"`
	MarketAwareness      float64 `json:"marketAwareness"`
	SolutionFit          float64 `json:"solutionFit"`
	ScopeDiscipline      float64 `json:"scopeDiscipline"`
	RequirementsQuality  float64 `json:"requirementsQuality"`
	UXCoverage           float64 `json:"uxCoverage"`
	TechnicalFeasibility float64 `json:"technicalFeasibility"`
	MetricsQuality       float64 `json:"metricsQuality"`
	RiskManagement       float64 `json:"riskManagement"`
	OverallScore         float64 `json:"overallScore"`
}

// ReviewDecision represents the outcome of a review.
type ReviewDecision string

const (
	// ReviewApprove means the PRD is approved.
	ReviewApprove ReviewDecision = "approve"

	// ReviewRevise means the PRD needs targeted revisions.
	ReviewRevise ReviewDecision = "revise"

	// ReviewReject means the PRD has blocking issues.
	ReviewReject ReviewDecision = "reject"

	// ReviewHumanReview means the PRD requires human review.
	ReviewHumanReview ReviewDecision = "human_review"
)

// Blocker represents an issue that blocks PRD approval.
type Blocker struct {
	// ID is the unique identifier for this blocker.
	ID string `json:"id"`

	// Category is the scoring category related to this blocker.
	Category string `json:"category"`

	// Description describes the blocking issue.
	Description string `json:"description"`
}

// RevisionTrigger represents an issue that requires revision.
type RevisionTrigger struct {
	// IssueID is the unique identifier for this issue.
	IssueID string `json:"issueId"`

	// Category is the scoring category related to this issue.
	Category string `json:"category"`

	// Severity indicates how severe the issue is (blocker, major, minor).
	Severity string `json:"severity"`

	// Description describes the issue.
	Description string `json:"description"`

	// RecommendedOwner suggests who should address this issue.
	RecommendedOwner string `json:"recommendedOwner,omitempty"`
}
