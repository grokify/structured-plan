package v2mom

import "fmt"

// ValidationError represents a validation error with path and severity.
type ValidationError struct {
	Path     string // JSON path to the error (e.g., "methods[0].measures")
	Message  string // Human-readable error message
	Severity string // "error" or "warning"
}

func (e ValidationError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s: %s", e.Path, e.Message)
	}
	return e.Message
}

// ValidationOptions configures validation behavior.
type ValidationOptions struct {
	// Structure enforcement mode: "flat", "nested", "hybrid", or "" (no enforcement)
	Structure string

	// RequireMethodMeasures requires each method to have at least one measure
	RequireMethodMeasures bool

	// MaxMeasuresPerMethod sets OKR best practice limit (3-5); 0 = unlimited
	MaxMeasuresPerMethod int

	// RequireGlobalObstacles requires at least one top-level obstacle
	RequireGlobalObstacles bool

	// WarnEmptyMethods generates warnings (not errors) for methods without measures
	WarnEmptyMethods bool
}

// DefaultValidationOptions returns sensible defaults for validation.
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		Structure:             "",    // No enforcement
		RequireMethodMeasures: false, // Don't require
		MaxMeasuresPerMethod:  0,     // Unlimited
		WarnEmptyMethods:      true,  // Warn but don't fail
	}
}

// OKRValidationOptions returns options suitable for OKR-aligned V2MOMs.
func OKRValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		Structure:             StructureNested,
		RequireMethodMeasures: true,
		MaxMeasuresPerMethod:  5, // OKR best practice: 3-5 key results
		WarnEmptyMethods:      false,
	}
}

// Validate validates the V2MOM against the provided options.
// Returns a slice of validation errors (may include warnings).
func (v *V2MOM) Validate(opts *ValidationOptions) []ValidationError {
	if opts == nil {
		opts = DefaultValidationOptions()
	}

	var errs []ValidationError

	// Basic required field validation
	if v.Vision == "" {
		errs = append(errs, ValidationError{
			Path:     "vision",
			Message:  "vision is required",
			Severity: "error",
		})
	}

	if len(v.Values) == 0 {
		errs = append(errs, ValidationError{
			Path:     "values",
			Message:  "at least one value is required",
			Severity: "error",
		})
	}

	if len(v.Methods) == 0 {
		errs = append(errs, ValidationError{
			Path:     "methods",
			Message:  "at least one method is required",
			Severity: "error",
		})
	}

	// Validate individual methods
	for i, m := range v.Methods {
		if m.Name == "" {
			errs = append(errs, ValidationError{
				Path:     fmt.Sprintf("methods[%d].name", i),
				Message:  "method name is required",
				Severity: "error",
			})
		}
	}

	// Validate individual values
	for i, val := range v.Values {
		if val.Name == "" {
			errs = append(errs, ValidationError{
				Path:     fmt.Sprintf("values[%d].name", i),
				Message:  "value name is required",
				Severity: "error",
			})
		}
	}

	// Structure-specific validation
	structure := opts.Structure
	if structure == "" && v.Metadata != nil && v.Metadata.Structure != "" {
		structure = v.Metadata.Structure
	}

	switch structure {
	case StructureFlat:
		errs = append(errs, v.validateFlat()...)
	case StructureNested:
		errs = append(errs, v.validateNested()...)
	case StructureHybrid, "":
		errs = append(errs, v.validateHybrid(opts)...)
	}

	// OKR best practice: 3-5 key results per objective
	if opts.MaxMeasuresPerMethod > 0 {
		for i, m := range v.Methods {
			if len(m.Measures) > opts.MaxMeasuresPerMethod {
				errs = append(errs, ValidationError{
					Path: fmt.Sprintf("methods[%d].measures", i),
					Message: fmt.Sprintf("method %q has %d measures (max %d recommended)",
						m.Name, len(m.Measures), opts.MaxMeasuresPerMethod),
					Severity: "warning",
				})
			}
		}
	}

	// Warn about methods without measures
	if opts.WarnEmptyMethods && !opts.RequireMethodMeasures {
		for i, m := range v.Methods {
			if len(m.Measures) == 0 && len(v.Measures) == 0 {
				errs = append(errs, ValidationError{
					Path:     fmt.Sprintf("methods[%d]", i),
					Message:  fmt.Sprintf("method %q has no measures", m.Name),
					Severity: "warning",
				})
			}
		}
	}

	// Require global obstacles
	if opts.RequireGlobalObstacles && len(v.Obstacles) == 0 {
		errs = append(errs, ValidationError{
			Path:     "obstacles",
			Message:  "at least one global obstacle is required",
			Severity: "error",
		})
	}

	return errs
}

// validateFlat validates flat (traditional V2MOM) structure.
func (v *V2MOM) validateFlat() []ValidationError {
	var errs []ValidationError

	// Measures must be at V2MOM level only
	for i, m := range v.Methods {
		if len(m.Measures) > 0 {
			errs = append(errs, ValidationError{
				Path:     fmt.Sprintf("methods[%d].measures", i),
				Message:  "flat mode: measures must be at V2MOM level, not under methods",
				Severity: "error",
			})
		}
		if len(m.Obstacles) > 0 {
			errs = append(errs, ValidationError{
				Path:     fmt.Sprintf("methods[%d].obstacles", i),
				Message:  "flat mode: obstacles must be at V2MOM level, not under methods",
				Severity: "error",
			})
		}
	}

	// Require global measures
	if len(v.Measures) == 0 {
		errs = append(errs, ValidationError{
			Path:     "measures",
			Message:  "flat mode: V2MOM-level measures required",
			Severity: "error",
		})
	}

	return errs
}

// validateNested validates nested (OKR-aligned) structure.
func (v *V2MOM) validateNested() []ValidationError {
	var errs []ValidationError

	// Global measures forbidden (global obstacles allowed for cross-cutting risks)
	if len(v.Measures) > 0 {
		errs = append(errs, ValidationError{
			Path:     "measures",
			Message:  "nested mode: measures must be under methods, not at V2MOM level",
			Severity: "error",
		})
	}

	// Each method must have measures
	for i, m := range v.Methods {
		if len(m.Measures) == 0 {
			errs = append(errs, ValidationError{
				Path:     fmt.Sprintf("methods[%d].measures", i),
				Message:  fmt.Sprintf("nested mode: method %q must have at least one measure", m.Name),
				Severity: "error",
			})
		}
	}

	return errs
}

// validateHybrid validates hybrid structure (both levels allowed).
func (v *V2MOM) validateHybrid(opts *ValidationOptions) []ValidationError {
	var errs []ValidationError

	// In hybrid mode, require method measures if option is set
	if opts.RequireMethodMeasures {
		for i, m := range v.Methods {
			if len(m.Measures) == 0 {
				errs = append(errs, ValidationError{
					Path:     fmt.Sprintf("methods[%d].measures", i),
					Message:  fmt.Sprintf("method %q has no measures", m.Name),
					Severity: "error",
				})
			}
		}
	}

	return errs
}

// HasErrors returns true if there are any errors (not just warnings).
func HasErrors(errs []ValidationError) bool {
	for _, e := range errs {
		if e.Severity != "warning" {
			return true
		}
	}
	return false
}

// Errors filters to only errors (excludes warnings).
func Errors(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if e.Severity != "warning" {
			result = append(result, e)
		}
	}
	return result
}

// Warnings filters to only warnings.
func Warnings(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if e.Severity == "warning" {
			result = append(result, e)
		}
	}
	return result
}
