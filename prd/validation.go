package prd

import (
	"fmt"
	"regexp"
)

// ValidationResult contains validation errors and warnings.
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Errors   []ValidationError   `json:"errors,omitempty"`
	Warnings []ValidationWarning `json:"warnings,omitempty"`
}

// ValidationError represents a validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationWarning represents a non-blocking issue.
type ValidationWarning struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate checks the Document for structural and content issues.
func Validate(doc *Document) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Required metadata fields
	if doc.Metadata.ID == "" {
		result.addError("metadata.id", "Document ID is required")
	}

	if doc.Metadata.Title == "" {
		result.addError("metadata.title", "Title is required")
	} else if len(doc.Metadata.Title) < 5 {
		result.addError("metadata.title", "Title must be at least 5 characters")
	}

	if len(doc.Metadata.Authors) == 0 {
		result.addWarning("metadata.authors", "No authors specified")
	}

	if doc.Metadata.Status == "" {
		result.addError("metadata.status", "Status is required")
	}

	// Executive summary
	if doc.ExecutiveSummary.ProblemStatement == "" {
		result.addWarning("executive_summary.problem_statement", "Problem statement is empty")
	}

	if doc.ExecutiveSummary.ProposedSolution == "" {
		result.addWarning("executive_summary.proposed_solution", "Proposed solution is empty")
	}

	// Objectives
	if len(doc.Objectives.BusinessObjectives) == 0 && len(doc.Objectives.ProductGoals) == 0 {
		result.addWarning("objectives", "No business objectives or product goals defined")
	}

	if len(doc.Objectives.SuccessMetrics) == 0 {
		result.addWarning("objectives.success_metrics", "No success metrics defined")
	}

	// Personas
	if len(doc.Personas) == 0 {
		result.addWarning("personas", "No personas defined")
	}

	// User stories
	if len(doc.UserStories) == 0 {
		result.addWarning("user_stories", "No user stories defined")
	}

	// Validate IDs are unique
	result.validateIDs(doc)

	// Validate traceability
	result.validateTraceability(doc)

	return result
}

func (r *ValidationResult) addError(field, message string) {
	r.Valid = false
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

func (r *ValidationResult) addWarning(field, message string) {
	r.Warnings = append(r.Warnings, ValidationWarning{Field: field, Message: message})
}

// validateIDs checks for duplicate and malformed IDs.
func (r *ValidationResult) validateIDs(doc *Document) {
	ids := make(map[string]string) // id -> location

	checkID := func(id, location string) {
		if id == "" {
			return
		}
		if existingLoc, exists := ids[id]; exists {
			r.addError(location, fmt.Sprintf("Duplicate ID '%s' (also at %s)", id, existingLoc))
		}
		ids[id] = location
	}

	// Check objective IDs
	for i, obj := range doc.Objectives.BusinessObjectives {
		checkID(obj.ID, fmt.Sprintf("objectives.business_objectives[%d].id", i))
	}
	for i, obj := range doc.Objectives.ProductGoals {
		checkID(obj.ID, fmt.Sprintf("objectives.product_goals[%d].id", i))
	}
	for i, m := range doc.Objectives.SuccessMetrics {
		checkID(m.ID, fmt.Sprintf("objectives.success_metrics[%d].id", i))
	}

	// Check persona IDs
	for i, p := range doc.Personas {
		checkID(p.ID, fmt.Sprintf("personas[%d].id", i))
	}

	// Check user story IDs
	for i, s := range doc.UserStories {
		checkID(s.ID, fmt.Sprintf("user_stories[%d].id", i))
	}

	// Check functional requirement IDs
	for i, req := range doc.Requirements.Functional {
		checkID(req.ID, fmt.Sprintf("requirements.functional[%d].id", i))
	}

	// Check NFR IDs
	for i, nfr := range doc.Requirements.NonFunctional {
		checkID(nfr.ID, fmt.Sprintf("requirements.non_functional[%d].id", i))
	}

	// Check roadmap phase IDs
	for i, phase := range doc.Roadmap.Phases {
		checkID(phase.ID, fmt.Sprintf("roadmap.phases[%d].id", i))
	}

	// Check problem ID (if present)
	if doc.Problem != nil && doc.Problem.ID != "" {
		checkID(doc.Problem.ID, "problem.id")
	}

	// Check market alternative IDs
	if doc.Market != nil {
		for i, alt := range doc.Market.Alternatives {
			checkID(alt.ID, fmt.Sprintf("market.alternatives[%d].id", i))
		}
	}

	// Check solution option IDs
	if doc.Solution != nil {
		for i, opt := range doc.Solution.SolutionOptions {
			checkID(opt.ID, fmt.Sprintf("solution.solution_options[%d].id", i))
		}
	}

	// Check decision IDs
	if doc.Decisions != nil {
		for i, d := range doc.Decisions.Records {
			checkID(d.ID, fmt.Sprintf("decisions.records[%d].id", i))
		}
	}
}

// validateTraceability checks cross-references between sections.
func (r *ValidationResult) validateTraceability(doc *Document) {
	// Build sets of defined IDs
	definedIDs := make(map[string]bool)

	// Collect all defined IDs
	for _, obj := range doc.Objectives.BusinessObjectives {
		if obj.ID != "" {
			definedIDs[obj.ID] = true
		}
	}
	for _, obj := range doc.Objectives.ProductGoals {
		if obj.ID != "" {
			definedIDs[obj.ID] = true
		}
	}
	for _, m := range doc.Objectives.SuccessMetrics {
		if m.ID != "" {
			definedIDs[m.ID] = true
		}
	}
	for _, p := range doc.Personas {
		if p.ID != "" {
			definedIDs[p.ID] = true
		}
	}
	for _, s := range doc.UserStories {
		if s.ID != "" {
			definedIDs[s.ID] = true
		}
	}
	for _, phase := range doc.Roadmap.Phases {
		if phase.ID != "" {
			definedIDs[phase.ID] = true
		}
	}
	if doc.Problem != nil && doc.Problem.ID != "" {
		definedIDs[doc.Problem.ID] = true
	}
	if doc.Solution != nil {
		for _, opt := range doc.Solution.SolutionOptions {
			if opt.ID != "" {
				definedIDs[opt.ID] = true
			}
		}
	}

	// Check user story references to personas
	for i, story := range doc.UserStories {
		if story.PersonaID != "" && !definedIDs[story.PersonaID] {
			r.addWarning(
				fmt.Sprintf("user_stories[%d].persona_id", i),
				fmt.Sprintf("Reference to undefined persona: %s", story.PersonaID),
			)
		}
		if story.PhaseID != "" && !definedIDs[story.PhaseID] {
			r.addWarning(
				fmt.Sprintf("user_stories[%d].phase_id", i),
				fmt.Sprintf("Reference to undefined phase: %s", story.PhaseID),
			)
		}
	}

	// Check functional requirement references
	for i, req := range doc.Requirements.Functional {
		for _, storyID := range req.UserStoryIDs {
			if storyID != "" && !definedIDs[storyID] {
				r.addWarning(
					fmt.Sprintf("requirements.functional[%d].user_story_ids", i),
					fmt.Sprintf("Reference to undefined user story: %s", storyID),
				)
			}
		}
		if req.PhaseID != "" && !definedIDs[req.PhaseID] {
			r.addWarning(
				fmt.Sprintf("requirements.functional[%d].phase_id", i),
				fmt.Sprintf("Reference to undefined phase: %s", req.PhaseID),
			)
		}
	}

	// Check solution problem references
	if doc.Solution != nil {
		for i, opt := range doc.Solution.SolutionOptions {
			for _, probID := range opt.ProblemsAddressed {
				if probID != "" && !definedIDs[probID] {
					r.addWarning(
						fmt.Sprintf("solution.solution_options[%d].problems_addressed", i),
						fmt.Sprintf("Reference to undefined problem: %s", probID),
					)
				}
			}
		}

		// Check selected solution exists
		if doc.Solution.SelectedSolutionID != "" {
			found := false
			for _, opt := range doc.Solution.SolutionOptions {
				if opt.ID == doc.Solution.SelectedSolutionID {
					found = true
					break
				}
			}
			if !found {
				r.addError(
					"solution.selected_solution_id",
					fmt.Sprintf("Selected solution '%s' not found in solution options", doc.Solution.SelectedSolutionID),
				)
			}
		}
	}
}

// isValidIDFormat checks if an ID matches the expected pattern.
func isValidIDFormat(id, prefix string) bool {
	pattern := regexp.MustCompile(fmt.Sprintf(`^%s-\d+$`, prefix))
	return pattern.MatchString(id)
}
