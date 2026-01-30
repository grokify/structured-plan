package common

// Status represents the document lifecycle status.
// Used across PRD, MRD, and TRD documents.
type Status string

const (
	StatusDraft      Status = "draft"
	StatusInReview   Status = "in_review"
	StatusApproved   Status = "approved"
	StatusDeprecated Status = "deprecated"
)

// Priority represents priority levels.
type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)
