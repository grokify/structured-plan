package prd

// ReviewsDefinition contains review outcomes and quality assessments.
type ReviewsDefinition struct {
	// ReviewBoardSummary is a summary from the review board.
	ReviewBoardSummary string `json:"review_board_summary,omitempty"`

	// QualityScores contains scores across quality dimensions.
	QualityScores *QualityScores `json:"quality_scores,omitempty"`

	// Decision is the review decision.
	Decision ReviewDecision `json:"decision,omitempty"`

	// Blockers are issues that block approval.
	Blockers []Blocker `json:"blockers,omitempty"`

	// RevisionTriggers are issues requiring revision.
	RevisionTriggers []RevisionTrigger `json:"revision_triggers,omitempty"`
}

// QualityScores contains scores across the 10 quality dimensions.
type QualityScores struct {
	ProblemDefinition    float64 `json:"problem_definition"`
	UserUnderstanding    float64 `json:"user_understanding"`
	MarketAwareness      float64 `json:"market_awareness"`
	SolutionFit          float64 `json:"solution_fit"`
	ScopeDiscipline      float64 `json:"scope_discipline"`
	RequirementsQuality  float64 `json:"requirements_quality"`
	UXCoverage           float64 `json:"ux_coverage"`
	TechnicalFeasibility float64 `json:"technical_feasibility"`
	MetricsQuality       float64 `json:"metrics_quality"`
	RiskManagement       float64 `json:"risk_management"`
	OverallScore         float64 `json:"overall_score"`
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
	IssueID string `json:"issue_id"`

	// Category is the scoring category related to this issue.
	Category string `json:"category"`

	// Severity indicates how severe the issue is (blocker, major, minor).
	Severity string `json:"severity"`

	// Description describes the issue.
	Description string `json:"description"`

	// RecommendedOwner suggests who should address this issue.
	RecommendedOwner string `json:"recommended_owner,omitempty"`
}
