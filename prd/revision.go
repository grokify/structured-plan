package prd

import "time"

// RevisionRecord documents a revision to the PRD.
type RevisionRecord struct {
	// Version is the version number after this revision.
	Version string `json:"version"`

	// Changes lists what changed in this revision.
	Changes []string `json:"changes"`

	// Trigger indicates what triggered this revision.
	Trigger RevisionTriggerType `json:"trigger"`

	// Date is when this revision was made.
	Date time.Time `json:"date"`

	// Author is who made this revision.
	Author string `json:"author,omitempty"`
}

// RevisionTriggerType indicates what triggered a revision.
type RevisionTriggerType string

const (
	// TriggerInitial is for the initial PRD creation.
	TriggerInitial RevisionTriggerType = "initial"

	// TriggerReview is for revisions from review feedback.
	TriggerReview RevisionTriggerType = "review"

	// TriggerScore is for revisions from scoring feedback.
	TriggerScore RevisionTriggerType = "score"

	// TriggerHuman is for revisions from human feedback.
	TriggerHuman RevisionTriggerType = "human"
)
