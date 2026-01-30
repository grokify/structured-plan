package prd

// CurrentState documents the existing state before the proposed solution.
type CurrentState struct {
	// Overview provides a high-level summary of the current state.
	Overview string `json:"overview"`

	// Approaches describes current approaches/solutions in use.
	Approaches []CurrentApproach `json:"approaches,omitempty"`

	// Problems lists specific problems with the current state.
	Problems []CurrentProblem `json:"problems,omitempty"`

	// TargetState describes the desired future state.
	TargetState string `json:"targetState"`

	// Metrics provides baseline metrics for comparison.
	Metrics []BaselineMetric `json:"metrics,omitempty"`

	// Diagrams provides links to architecture or flow diagrams.
	Diagrams []DiagramRef `json:"diagrams,omitempty"`
}

// CurrentApproach describes an existing approach or solution.
type CurrentApproach struct {
	// ID is the unique identifier for this approach.
	ID string `json:"id,omitempty"`

	// Name is the identifier for this approach.
	Name string `json:"name"`

	// Description explains how this approach works.
	Description string `json:"description"`

	// Problems lists issues with this approach.
	Problems []string `json:"problems,omitempty"`

	// Usage indicates adoption level (e.g., "80% of customers").
	Usage string `json:"usage,omitempty"`

	// Owner is the team/person responsible for this approach.
	Owner string `json:"owner,omitempty"`
}

// CurrentProblem describes a specific problem with the current state.
type CurrentProblem struct {
	// ID is the unique identifier for this problem.
	ID string `json:"id,omitempty"`

	// Description of the problem.
	Description string `json:"description"`

	// Impact on users or business.
	Impact string `json:"impact,omitempty"`

	// Frequency of occurrence.
	Frequency string `json:"frequency,omitempty"`

	// AffectedUsers describes who is impacted.
	AffectedUsers string `json:"affectedUsers,omitempty"`

	// RelatedIDs links to related requirements or risks.
	RelatedIDs []string `json:"relatedIds,omitempty"`
}

// BaselineMetric provides current state metrics for comparison.
type BaselineMetric struct {
	// ID is the unique identifier for this metric.
	ID string `json:"id,omitempty"`

	// Name of the metric.
	Name string `json:"name"`

	// CurrentValue is the baseline value.
	CurrentValue string `json:"currentValue"`

	// TargetValue is the desired value after implementation.
	TargetValue string `json:"targetValue,omitempty"`

	// MeasurementMethod describes how this is measured.
	MeasurementMethod string `json:"measurementMethod,omitempty"`

	// Source is where the current value was obtained.
	Source string `json:"source,omitempty"`
}

// DiagramRef references a diagram or visual.
type DiagramRef struct {
	// Title is the diagram title.
	Title string `json:"title"`

	// URL is the link to the diagram.
	URL string `json:"url"`

	// Description provides context.
	Description string `json:"description,omitempty"`

	// Type is the diagram type (e.g., "architecture", "flow", "sequence").
	Type string `json:"type,omitempty"`
}
