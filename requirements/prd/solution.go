package prd

// SolutionDefinition contains solution options and selection rationale.
type SolutionDefinition struct {
	// SolutionOptions are the possible solutions considered.
	SolutionOptions []SolutionOption `json:"solution_options,omitempty"`

	// SelectedSolutionID is the ID of the chosen solution.
	SelectedSolutionID string `json:"selected_solution_id,omitempty"`

	// SolutionRationale explains why the selected solution was chosen.
	SolutionRationale string `json:"solution_rationale,omitempty"`

	// Confidence is the confidence level in the solution (0.0-1.0).
	Confidence float64 `json:"confidence,omitempty"`
}

// SolutionOption represents a possible solution approach.
type SolutionOption struct {
	// ID is the unique identifier for this solution option.
	ID string `json:"id"`

	// Name is the name of this solution option.
	Name string `json:"name"`

	// Description provides details about the solution.
	Description string `json:"description,omitempty"`

	// ProblemsAddressed lists problem IDs this solution addresses.
	ProblemsAddressed []string `json:"problems_addressed,omitempty"`

	// Benefits are advantages of this solution.
	Benefits []string `json:"benefits,omitempty"`

	// Tradeoffs are compromises or downsides of this solution.
	Tradeoffs []string `json:"tradeoffs,omitempty"`

	// Risks are potential risks of this solution.
	Risks []string `json:"risks,omitempty"`

	// EstimatedEffort is a high-level effort estimate.
	EstimatedEffort string `json:"estimated_effort,omitempty"`
}

// SelectedSolution returns the selected solution option, or nil if none selected.
func (s *SolutionDefinition) SelectedSolution() *SolutionOption {
	if s == nil || s.SelectedSolutionID == "" {
		return nil
	}
	for i := range s.SolutionOptions {
		if s.SolutionOptions[i].ID == s.SelectedSolutionID {
			return &s.SolutionOptions[i]
		}
	}
	return nil
}
