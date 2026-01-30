// Package v2mom provides types and utilities for V2MOM strategic planning documents.
// V2MOM (Vision, Values, Methods, Obstacles, Measures) is a framework created by
// Marc Benioff at Salesforce for organizational alignment.
//
// This package supports both traditional flat V2MOM structure and OKR-aligned
// nested structure where Measures (Key Results) are nested under Methods (Objectives).
package v2mom

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Structure constants define V2MOM organizational styles.
const (
	StructureFlat   = "flat"   // Traditional V2MOM: measures/obstacles at V2MOM level only
	StructureNested = "nested" // OKR-aligned: measures under Methods, global obstacles allowed
	StructureHybrid = "hybrid" // Both levels allowed (default)
)

// Terminology constants define display label modes.
const (
	TerminologyV2MOM  = "v2mom"  // Methods/Measures/Obstacles
	TerminologyOKR    = "okr"    // Objectives/Key Results/Risks
	TerminologyHybrid = "hybrid" // Methods (Objectives)/Measures (Key Results)/Obstacles
)

// Status constants for document lifecycle.
const (
	StatusDraft     = "Draft"
	StatusInReview  = "In Review"
	StatusApproved  = "Approved"
	StatusActive    = "Active"
	StatusCompleted = "Completed"
	StatusArchived  = "Archived"
)

// Priority constants.
const (
	PriorityP0 = "P0"
	PriorityP1 = "P1"
	PriorityP2 = "P2"
	PriorityP3 = "P3"
)

// V2MOM represents a complete V2MOM strategic planning document.
// It supports both traditional flat structure and OKR-aligned nested structure.
type V2MOM struct {
	Schema   string    `json:"$schema,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	Vision   string    `json:"vision"`
	Values   []Value   `json:"values"`
	Methods  []Method  `json:"methods"`
	// Global obstacles (traditional V2MOM or cross-cutting in nested mode)
	Obstacles []Obstacle `json:"obstacles,omitempty"`
	// Global measures (traditional V2MOM only; use Method.Measures for OKR alignment)
	Measures []Measure `json:"measures,omitempty"`
	// Projects for roadmap visualization
	Projects []Project `json:"projects,omitempty"`
}

// Metadata contains document metadata and configuration.
type Metadata struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Author     string    `json:"author,omitempty"`
	Team       string    `json:"team,omitempty"`
	FiscalYear string    `json:"fiscalYear,omitempty"` // e.g., "FY2025"
	Quarter    string    `json:"quarter,omitempty"`    // Q1, Q2, Q3, Q4, H1, H2, Annual
	Version    string    `json:"version,omitempty"`
	Status     string    `json:"status,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	ParentID   string    `json:"parentId,omitempty"` // For cascading V2MOMs

	// Structure defines the V2MOM organizational style.
	// - "flat": Traditional V2MOM (measures/obstacles at V2MOM level only)
	// - "nested": OKR-aligned (measures under Methods, global obstacles allowed)
	// - "hybrid": Both levels allowed (default)
	Structure string `json:"structure,omitempty"`

	// Terminology defines display labels for rendering.
	// - "v2mom": Methods/Measures/Obstacles (default)
	// - "okr": Objectives/Key Results/Risks
	// - "hybrid": Methods (Objectives)/Measures (Key Results)/Obstacles
	Terminology string `json:"terminology,omitempty"`
}

// Value represents a guiding principle that supports the vision.
type Value struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"` // 1 = highest priority
}

// Method represents an action or objective to achieve the vision.
// In OKR terminology, this corresponds to an Objective.
type Method struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority,omitempty"` // P0, P1, P2, P3
	Status      string `json:"status,omitempty"`   // Not Started, Planning, In Progress, At Risk, Completed, Cancelled
	Owner       string `json:"owner,omitempty"`
	StartDate   string `json:"startDate,omitempty"` // ISO 8601 date
	EndDate     string `json:"endDate,omitempty"`   // ISO 8601 date

	// Nested measures (OKR Key Results) - used in nested/hybrid mode
	Measures []Measure `json:"measures,omitempty"`
	// Method-specific obstacles - used in nested/hybrid mode
	Obstacles []Obstacle `json:"obstacles,omitempty"`
	// Linked project IDs
	Projects []string `json:"projects,omitempty"`
}

// Obstacle represents a challenge or risk that could prevent success.
type Obstacle struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Severity    string `json:"severity,omitempty"`   // Low, Medium, High, Critical
	Likelihood  string `json:"likelihood,omitempty"` // Low, Medium, High
	Mitigation  string `json:"mitigation,omitempty"`
	Status      string `json:"status,omitempty"` // Identified, Mitigating, Resolved, Accepted
}

// Measure represents a success metric or key result.
// In OKR terminology, this corresponds to a Key Result.
type Measure struct {
	ID          string  `json:"id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Baseline    string  `json:"baseline,omitempty"` // Starting value
	Target      string  `json:"target,omitempty"`   // Target value
	Current     string  `json:"current,omitempty"`  // Current value
	Unit        string  `json:"unit,omitempty"`     // Unit of measurement
	Progress    float64 `json:"progress,omitempty"` // 0.0-1.0 (OKR scoring)
	Timeline    string  `json:"timeline,omitempty"` // Target timeline
	Status      string  `json:"status,omitempty"`   // On Track, At Risk, Behind, Achieved, Missed
}

// Project represents a roadmap project linked to methods.
type Project struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Category      string            `json:"category,omitempty"`
	MethodID      string            `json:"methodId,omitempty"`
	Priority      string            `json:"priority,omitempty"` // P0, P1, P2, P3
	Status        string            `json:"status,omitempty"`   // Proposed, Approved, In Progress, Completed, Cancelled
	StartDate     string            `json:"startDate,omitempty"`
	EndDate       string            `json:"endDate,omitempty"`
	Quarter       string            `json:"quarter,omitempty"`
	Dependencies  []string          `json:"dependencies,omitempty"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty"` // jira, aha, productboard, confluence URLs
}

// ReadFile reads a V2MOM from a JSON file.
func ReadFile(filepath string) (*V2MOM, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return Parse(data)
}

// Parse parses V2MOM JSON data.
func Parse(data []byte) (*V2MOM, error) {
	var v V2MOM
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &v, nil
}

// JSON returns the V2MOM as formatted JSON.
func (v *V2MOM) JSON() ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// WriteFile writes the V2MOM to a JSON file.
func (v *V2MOM) WriteFile(filepath string) error {
	data, err := v.JSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// AllMeasures returns all measures (global + nested), flattened.
func (v *V2MOM) AllMeasures() []Measure {
	all := make([]Measure, 0, len(v.Measures))
	all = append(all, v.Measures...)
	for _, m := range v.Methods {
		all = append(all, m.Measures...)
	}
	return all
}

// AllObstacles returns all obstacles (global + nested), flattened.
func (v *V2MOM) AllObstacles() []Obstacle {
	all := make([]Obstacle, 0, len(v.Obstacles))
	all = append(all, v.Obstacles...)
	for _, m := range v.Methods {
		all = append(all, m.Obstacles...)
	}
	return all
}

// InferStructure detects the structure based on data placement.
func (v *V2MOM) InferStructure() string {
	hasGlobalMeasures := len(v.Measures) > 0
	hasNestedMeasures := false
	for _, m := range v.Methods {
		if len(m.Measures) > 0 {
			hasNestedMeasures = true
			break
		}
	}

	switch {
	case hasGlobalMeasures && !hasNestedMeasures:
		return StructureFlat
	case !hasGlobalMeasures && hasNestedMeasures:
		return StructureNested
	case hasGlobalMeasures && hasNestedMeasures:
		return StructureHybrid
	default:
		return "" // No measures anywhere
	}
}

// HasNestedStructure returns true if any method has nested measures or obstacles.
func (v *V2MOM) HasNestedStructure() bool {
	for _, m := range v.Methods {
		if len(m.Measures) > 0 || len(m.Obstacles) > 0 {
			return true
		}
	}
	return false
}

// GetStructure returns the structure mode, using metadata if set or inferring from data.
func (v *V2MOM) GetStructure() string {
	if v.Metadata != nil && v.Metadata.Structure != "" {
		return v.Metadata.Structure
	}
	return v.InferStructure()
}

// GetTerminology returns the terminology mode, defaulting to "v2mom".
func (v *V2MOM) GetTerminology() string {
	if v.Metadata != nil && v.Metadata.Terminology != "" {
		return v.Metadata.Terminology
	}
	return TerminologyV2MOM
}
