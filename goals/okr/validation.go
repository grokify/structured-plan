package okr

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation issue.
type ValidationError struct {
	Path    string // JSON path to the problematic field
	Message string
	IsError bool // true for errors, false for warnings
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationOptions configures validation behavior.
type ValidationOptions struct {
	RequireKeyResults   bool // Require at least one key result per objective
	RequireScores       bool // Require scores to be set on key results
	MinKeyResultsPerObj int  // Minimum key results per objective (default: 1)
	MaxKeyResultsPerObj int  // Maximum key results per objective (default: 5, 0 = no limit)
	MaxObjectives       int  // Maximum objectives per document (default: 5, 0 = no limit)
	RequireTargets      bool // Require target values on key results
	ValidateScoreRange  bool // Ensure scores are in 0.0-1.0 range
	RequireTimeframe    bool // Require timeframe on active objectives
}

// DefaultValidationOptions returns sensible defaults.
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireKeyResults:   true,
		RequireScores:       false,
		MinKeyResultsPerObj: 1,
		MaxKeyResultsPerObj: 5,
		MaxObjectives:       5,
		RequireTargets:      false,
		ValidateScoreRange:  true,
		RequireTimeframe:    false,
	}
}

// StrictValidationOptions returns strict validation settings.
func StrictValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireKeyResults:   true,
		RequireScores:       true,
		MinKeyResultsPerObj: 2,
		MaxKeyResultsPerObj: 5,
		MaxObjectives:       5,
		RequireTargets:      true,
		ValidateScoreRange:  true,
		RequireTimeframe:    true,
	}
}

// Validate checks the OKR document for issues.
func (doc *OKRDocument) Validate(opts *ValidationOptions) []ValidationError {
	if opts == nil {
		opts = DefaultValidationOptions()
	}

	var errs []ValidationError

	// Validate objectives exist
	if len(doc.Objectives) == 0 {
		errs = append(errs, ValidationError{
			Path:    "objectives",
			Message: "at least one objective is required",
			IsError: true,
		})
	}

	// Check max objectives
	if opts.MaxObjectives > 0 && len(doc.Objectives) > opts.MaxObjectives {
		errs = append(errs, ValidationError{
			Path:    "objectives",
			Message: fmt.Sprintf("too many objectives: %d (max: %d)", len(doc.Objectives), opts.MaxObjectives),
			IsError: false, // Warning - OKR best practice suggests 3-5 objectives
		})
	}

	// Validate each objective
	for i, obj := range doc.Objectives {
		objPath := fmt.Sprintf("objectives[%d]", i)
		errs = append(errs, validateObjective(obj, objPath, opts)...)
	}

	// Validate metadata if present
	if doc.Metadata != nil {
		errs = append(errs, validateMetadata(doc.Metadata)...)
	}

	return errs
}

func validateObjective(obj Objective, path string, opts *ValidationOptions) []ValidationError {
	var errs []ValidationError

	// Title is required
	if strings.TrimSpace(obj.Title) == "" {
		errs = append(errs, ValidationError{
			Path:    path + ".title",
			Message: "objective title is required",
			IsError: true,
		})
	}

	// Validate key results exist
	if opts.RequireKeyResults && len(obj.KeyResults) == 0 {
		errs = append(errs, ValidationError{
			Path:    path + ".keyResults",
			Message: "at least one key result is required",
			IsError: true,
		})
	}

	// Check minimum key results
	if opts.MinKeyResultsPerObj > 0 && len(obj.KeyResults) < opts.MinKeyResultsPerObj {
		errs = append(errs, ValidationError{
			Path:    path + ".keyResults",
			Message: fmt.Sprintf("too few key results: %d (min: %d)", len(obj.KeyResults), opts.MinKeyResultsPerObj),
			IsError: false, // Warning
		})
	}

	// Check maximum key results
	if opts.MaxKeyResultsPerObj > 0 && len(obj.KeyResults) > opts.MaxKeyResultsPerObj {
		errs = append(errs, ValidationError{
			Path:    path + ".keyResults",
			Message: fmt.Sprintf("too many key results: %d (max: %d)", len(obj.KeyResults), opts.MaxKeyResultsPerObj),
			IsError: false, // Warning - OKR best practice suggests 3-5 KRs per objective
		})
	}

	// Validate each key result
	for j, kr := range obj.KeyResults {
		krPath := fmt.Sprintf("%s.keyResults[%d]", path, j)
		errs = append(errs, validateKeyResult(kr, krPath, opts)...)
	}

	// Validate progress is in range
	if obj.Progress < 0 || obj.Progress > 1 {
		errs = append(errs, ValidationError{
			Path:    path + ".progress",
			Message: fmt.Sprintf("progress must be between 0.0 and 1.0, got %.2f", obj.Progress),
			IsError: true,
		})
	}

	// Validate timeframe for active objectives
	if opts.RequireTimeframe && strings.TrimSpace(obj.Timeframe) == "" {
		if obj.Status == "" || obj.Status == StatusActive {
			errs = append(errs, ValidationError{
				Path:    path + ".timeframe",
				Message: "timeframe is required for active objectives (e.g., Q2 2026, H1 2026)",
				IsError: false, // Warning
			})
		}
	}

	return errs
}

func validateKeyResult(kr KeyResult, path string, opts *ValidationOptions) []ValidationError {
	var errs []ValidationError

	// Title is required
	if strings.TrimSpace(kr.Title) == "" {
		errs = append(errs, ValidationError{
			Path:    path + ".title",
			Message: "key result title is required",
			IsError: true,
		})
	}

	// Validate score if required or if present
	if opts.RequireScores && kr.Score == 0 {
		errs = append(errs, ValidationError{
			Path:    path + ".score",
			Message: "score is required",
			IsError: false, // Warning - score might legitimately be 0
		})
	}

	if opts.ValidateScoreRange && (kr.Score < 0 || kr.Score > 1) {
		errs = append(errs, ValidationError{
			Path:    path + ".score",
			Message: fmt.Sprintf("score must be between 0.0 and 1.0, got %.2f", kr.Score),
			IsError: true,
		})
	}

	// Validate target if required
	if opts.RequireTargets && strings.TrimSpace(kr.Target) == "" {
		errs = append(errs, ValidationError{
			Path:    path + ".target",
			Message: "target is required for measurable key results",
			IsError: false, // Warning
		})
	}

	// Validate confidence value
	if kr.Confidence != "" && kr.Confidence != ConfidenceLow &&
		kr.Confidence != ConfidenceMedium && kr.Confidence != ConfidenceHigh {
		errs = append(errs, ValidationError{
			Path:    path + ".confidence",
			Message: fmt.Sprintf("invalid confidence value: %s (expected: Low, Medium, High)", kr.Confidence),
			IsError: true,
		})
	}

	return errs
}

func validateMetadata(meta *Metadata) []ValidationError {
	var errs []ValidationError

	// Name should be present
	if strings.TrimSpace(meta.Name) == "" {
		errs = append(errs, ValidationError{
			Path:    "metadata.name",
			Message: "OKR document name is recommended",
			IsError: false, // Warning
		})
	}

	// Status should be valid if present
	if meta.Status != "" &&
		meta.Status != StatusDraft &&
		meta.Status != StatusActive &&
		meta.Status != StatusCompleted &&
		meta.Status != StatusCancelled {
		errs = append(errs, ValidationError{
			Path:    "metadata.status",
			Message: fmt.Sprintf("invalid status: %s (expected: Draft, Active, Completed, Cancelled)", meta.Status),
			IsError: true,
		})
	}

	return errs
}

// Errors returns only error-level validation results.
func Errors(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if e.IsError {
			result = append(result, e)
		}
	}
	return result
}

// Warnings returns only warning-level validation results.
func Warnings(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if !e.IsError {
			result = append(result, e)
		}
	}
	return result
}

// IsValid returns true if there are no error-level validation issues.
func IsValid(errs []ValidationError) bool {
	return len(Errors(errs)) == 0
}
