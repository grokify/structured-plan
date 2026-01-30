package prd

import "time"

// OpenItem represents a pending decision or question that needs resolution.
// Unlike DecisionRecord (for completed decisions), OpenItem tracks items
// that are still under consideration with options and tradeoffs.
type OpenItem struct {
	// ID is the unique identifier for this open item.
	ID string `json:"id"`

	// Title is a brief summary of the decision needed.
	Title string `json:"title"`

	// Description provides detailed context about what needs to be decided.
	Description string `json:"description,omitempty"`

	// Context explains the background and why this decision is needed.
	Context string `json:"context,omitempty"`

	// Options are the available choices with their tradeoffs.
	Options []Option `json:"options,omitempty"`

	// Status is the current status of this open item.
	Status OpenItemStatus `json:"status,omitempty"`

	// Priority indicates how urgent this decision is.
	Priority Priority `json:"priority,omitempty"`

	// Owner is the person or group responsible for making this decision.
	Owner string `json:"owner,omitempty"`

	// Stakeholders are people who should be consulted.
	Stakeholders []string `json:"stakeholders,omitempty"`

	// DueDate is when this decision needs to be made.
	DueDate *time.Time `json:"due_date,omitempty"`

	// CreatedAt is when this open item was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// Resolution documents the final decision once made.
	Resolution *OpenItemResolution `json:"resolution,omitempty"`

	// RelatedIDs links to related requirements, risks, or other items.
	RelatedIDs []string `json:"related_ids,omitempty"`

	// Tags for filtering by topic/domain.
	Tags []string `json:"tags,omitempty"`
}

// Option represents one possible choice for an open item decision.
type Option struct {
	// ID is the unique identifier for this option.
	ID string `json:"id"`

	// Title is a brief name for this option.
	Title string `json:"title"`

	// Description explains this option in detail.
	Description string `json:"description,omitempty"`

	// Pros lists the benefits and advantages of this option.
	Pros []string `json:"pros,omitempty"`

	// Cons lists the drawbacks and disadvantages of this option.
	Cons []string `json:"cons,omitempty"`

	// Effort estimates the implementation effort.
	Effort EffortLevel `json:"effort,omitempty"`

	// Risk estimates the risk level of this option.
	Risk RiskLevel `json:"risk,omitempty"`

	// Cost provides cost estimate or impact.
	Cost string `json:"cost,omitempty"`

	// Timeline provides time estimate or impact.
	Timeline string `json:"timeline,omitempty"`

	// Recommended indicates if this is the recommended option.
	Recommended bool `json:"recommended,omitempty"`

	// RecommendationRationale explains why this option is recommended (if applicable).
	RecommendationRationale string `json:"recommendation_rationale,omitempty"`
}

// OpenItemStatus represents the status of an open item.
type OpenItemStatus string

const (
	// OpenItemStatusOpen means the item is awaiting decision.
	OpenItemStatusOpen OpenItemStatus = "open"

	// OpenItemStatusInDiscussion means the item is being actively discussed.
	OpenItemStatusInDiscussion OpenItemStatus = "in_discussion"

	// OpenItemStatusBlocked means the item is blocked on something else.
	OpenItemStatusBlocked OpenItemStatus = "blocked"

	// OpenItemStatusResolved means a decision has been made.
	OpenItemStatusResolved OpenItemStatus = "resolved"

	// OpenItemStatusDeferred means the decision has been postponed.
	OpenItemStatusDeferred OpenItemStatus = "deferred"
)

// EffortLevel represents effort estimates.
type EffortLevel string

const (
	EffortLow    EffortLevel = "low"
	EffortMedium EffortLevel = "medium"
	EffortHigh   EffortLevel = "high"
)

// RiskLevel represents risk levels for options.
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

// OpenItemResolution documents how an open item was resolved.
type OpenItemResolution struct {
	// ChosenOptionID is the ID of the option that was selected.
	ChosenOptionID string `json:"chosen_option_id,omitempty"`

	// Decision summarizes the final decision.
	Decision string `json:"decision"`

	// Rationale explains why this decision was made.
	Rationale string `json:"rationale,omitempty"`

	// DecidedBy is who made the final decision.
	DecidedBy string `json:"decided_by,omitempty"`

	// DecidedAt is when the decision was made.
	DecidedAt *time.Time `json:"decided_at,omitempty"`

	// Notes captures any additional context.
	Notes string `json:"notes,omitempty"`
}
