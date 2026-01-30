// Package okr provides types and utilities for OKR (Objectives and Key Results)
// goal-setting documents.
//
// OKR is a framework popularized by Intel and Google for setting and
// communicating goals and results. Each Objective has associated Key Results
// that define how success is measured.
//
// Key characteristics of OKRs:
//   - Objectives are qualitative, inspirational goals
//   - Key Results are quantitative, measurable outcomes
//   - Progress is scored 0.0-1.0, where 0.7 is typically considered success
//   - OKRs are typically set quarterly with annual themes
//
// This package consolidates OKR types from both standalone OKR documents
// and PRD-embedded OKRs, providing a single source of truth.
package okr

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Status constants for OKR lifecycle.
const (
	StatusDraft     = "Draft"
	StatusActive    = "Active"
	StatusCompleted = "Completed"
	StatusCancelled = "Cancelled"
)

// Confidence constants for Key Result confidence levels.
const (
	ConfidenceLow    = "Low"    // 0-30% likely to achieve
	ConfidenceMedium = "Medium" // 30-70% likely to achieve
	ConfidenceHigh   = "High"   // 70-100% likely to achieve
)

// ScoreThresholds for OKR evaluation.
const (
	ScoreExcellent = 1.0 // Fully achieved
	ScoreGood      = 0.7 // Typical success threshold
	ScoreOK        = 0.4 // Partial achievement
	ScoreFailed    = 0.0 // Not achieved
)

// OKRDocument represents a complete OKR document containing objectives.
// Used for standalone OKR files (team/company OKRs).
type OKRDocument struct {
	Schema     string      `json:"$schema,omitempty"`
	Metadata   *Metadata   `json:"metadata,omitempty"`
	Theme      string      `json:"theme,omitempty"`     // Annual or quarterly theme
	Objectives []Objective `json:"objectives"`          // The OKRs
	Risks      []Risk      `json:"risks,omitempty"`     // Cross-cutting risks
	Alignment  *Alignment  `json:"alignment,omitempty"` // Links to parent/company OKRs
}

// Metadata contains document metadata.
type Metadata struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Owner      string    `json:"owner,omitempty"`
	Team       string    `json:"team,omitempty"`
	Period     string    `json:"period,omitempty"`     // e.g., "2025-Q1", "FY2025"
	PeriodType string    `json:"periodType,omitempty"` // "quarter", "half", "annual"
	Version    string    `json:"version,omitempty"`
	Status     string    `json:"status,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
}

// Objective represents an inspirational, qualitative goal.
// Merged from standalone OKR and PRD objective types.
type Objective struct {
	ID          string      `json:"id,omitempty"`
	Title       string      `json:"title"`                 // Short display title
	Description string      `json:"description,omitempty"` // Detailed description
	Rationale   string      `json:"rationale,omitempty"`   // Why this objective matters (from PRD)
	Category    string      `json:"category,omitempty"`    // Business, Product, Team, etc. (from PRD)
	Owner       string      `json:"owner,omitempty"`       // Person or team responsible
	Timeframe   string      `json:"timeframe,omitempty"`   // Target period (e.g., "Q2 2026")
	Status      string      `json:"status,omitempty"`      // Draft, Active, Completed, Cancelled
	KeyResults  []KeyResult `json:"keyResults"`            // Must have 1+ Key Results
	Progress    float64     `json:"progress,omitempty"`    // Calculated from key results (0.0-1.0)
	Risks       []Risk      `json:"risks,omitempty"`       // Objective-specific risks
	ParentID    string      `json:"parentId,omitempty"`    // Link to parent/company objective
	AlignedWith []string    `json:"alignedWith,omitempty"` // IDs of objectives this supports
	Tags        []string    `json:"tags,omitempty"`        // For filtering by topic/domain (from PRD)
}

// KeyResult represents a measurable outcome for an Objective.
// Merged from standalone OKR and PRD key result types.
type KeyResult struct {
	ID                string        `json:"id,omitempty"`
	Title             string        `json:"title"`                       // Short display title
	Description       string        `json:"description,omitempty"`       // Detailed description
	Owner             string        `json:"owner,omitempty"`             // Person or team responsible
	Metric            string        `json:"metric,omitempty"`            // What is being measured
	Baseline          string        `json:"baseline,omitempty"`          // Starting value
	Target            string        `json:"target,omitempty"`            // Target value to achieve
	Current           string        `json:"current,omitempty"`           // Current value
	Unit              string        `json:"unit,omitempty"`              // Unit of measurement
	MeasurementMethod string        `json:"measurementMethod,omitempty"` // How it's measured (from PRD)
	Score             float64       `json:"score,omitempty"`             // 0.0-1.0 achievement score
	Confidence        string        `json:"confidence,omitempty"`        // Low, Medium, High
	Status            string        `json:"status,omitempty"`            // On Track, At Risk, Behind, Achieved
	DueDate           string        `json:"dueDate,omitempty"`           // ISO 8601 date
	PhaseTargets      []PhaseTarget `json:"phaseTargets,omitempty"`      // Per-phase targets for roadmap alignment (from PRD)
	Tags              []string      `json:"tags,omitempty"`              // For filtering by topic/domain (from PRD)
}

// PhaseTarget represents a Key Result target for a specific roadmap phase.
// This enables alignment between OKRs and roadmap phases.
type PhaseTarget struct {
	PhaseID string `json:"phaseId"`          // Reference to roadmap phase
	Target  string `json:"target"`           // Target value for this phase
	Status  string `json:"status,omitempty"` // not_started, in_progress, achieved, missed
	Actual  string `json:"actual,omitempty"` // Actual value achieved
	Notes   string `json:"notes,omitempty"`  // Commentary on progress
}

// Risk represents a challenge or risk to achieving objectives.
type Risk struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Impact      string `json:"impact,omitempty"`     // Low, Medium, High, Critical
	Likelihood  string `json:"likelihood,omitempty"` // Low, Medium, High
	Mitigation  string `json:"mitigation,omitempty"`
	Status      string `json:"status,omitempty"` // Identified, Mitigating, Resolved, Accepted
}

// Alignment represents how OKRs align with parent/company objectives.
type Alignment struct {
	ParentOKRID   string   `json:"parentOkrId,omitempty"`   // Parent OKR document ID
	CompanyOKRIDs []string `json:"companyOkrIds,omitempty"` // Company-level objective IDs this supports
}

// DefaultFilename is the standard OKR filename.
const DefaultFilename = "okr.json"

// ReadFile reads an OKR document from a JSON file.
func ReadFile(filepath string) (*OKRDocument, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return Parse(data)
}

// Parse parses OKR JSON data.
func Parse(data []byte) (*OKRDocument, error) {
	var doc OKRDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &doc, nil
}

// JSON returns the OKR document as formatted JSON.
func (doc *OKRDocument) JSON() ([]byte, error) {
	return json.MarshalIndent(doc, "", "  ")
}

// WriteFile writes the OKR document to a JSON file.
func (doc *OKRDocument) WriteFile(filepath string) error {
	data, err := doc.JSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// CalculateProgress calculates the overall progress of an Objective
// based on its Key Results. Uses average scoring by default.
func (o *Objective) CalculateProgress() float64 {
	if len(o.KeyResults) == 0 {
		return 0
	}
	var total float64
	for _, kr := range o.KeyResults {
		total += kr.Score
	}
	return total / float64(len(o.KeyResults))
}

// UpdateProgress recalculates the progress for an objective.
func (o *Objective) UpdateProgress() {
	o.Progress = o.CalculateProgress()
}

// CalculateOverallProgress calculates the overall OKR document progress.
func (doc *OKRDocument) CalculateOverallProgress() float64 {
	if len(doc.Objectives) == 0 {
		return 0
	}
	var total float64
	for _, obj := range doc.Objectives {
		total += obj.CalculateProgress()
	}
	return total / float64(len(doc.Objectives))
}

// AllKeyResults returns all key results from all objectives, flattened.
func (doc *OKRDocument) AllKeyResults() []KeyResult {
	var all []KeyResult
	for _, obj := range doc.Objectives {
		all = append(all, obj.KeyResults...)
	}
	return all
}

// AllRisks returns all risks (global + objective-specific), flattened.
func (doc *OKRDocument) AllRisks() []Risk {
	all := make([]Risk, 0, len(doc.Risks))
	all = append(all, doc.Risks...)
	for _, obj := range doc.Objectives {
		all = append(all, obj.Risks...)
	}
	return all
}

// ScoreGrade returns a letter grade for a score.
func ScoreGrade(score float64) string {
	switch {
	case score >= 0.9:
		return "A"
	case score >= 0.7:
		return "B"
	case score >= 0.4:
		return "C"
	case score >= 0.2:
		return "D"
	default:
		return "F"
	}
}

// ScoreDescription returns a description for a score.
func ScoreDescription(score float64) string {
	switch {
	case score >= 0.9:
		return "Exceeded expectations"
	case score >= 0.7:
		return "Achieved target"
	case score >= 0.4:
		return "Partial achievement"
	case score >= 0.2:
		return "Below expectations"
	default:
		return "Not achieved"
	}
}

// New creates a new OKR document with required fields initialized.
func New(id, name, owner string) *OKRDocument {
	now := time.Now()
	return &OKRDocument{
		Metadata: &Metadata{
			ID:        id,
			Name:      name,
			Owner:     owner,
			Status:    StatusDraft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Objectives: []Objective{},
	}
}

// GenerateID generates an OKR ID based on the current date.
// Format: OKR-YYYY-DDD where DDD is the day of year.
func GenerateID() string {
	now := time.Now()
	return fmt.Sprintf("OKR-%d-%03d", now.Year(), now.YearDay())
}

// OKRSet represents a set of OKRs within a PRD or other document.
// This is the embedded form (vs standalone OKRDocument).
type OKRSet struct {
	OKRs []OKR `json:"okrs"`
}

// OKR represents an Objective with its Key Results in nested form.
// This format is commonly used in PRDs for cleaner nesting.
type OKR struct {
	Objective  Objective   `json:"objective"`
	KeyResults []KeyResult `json:"key_results"` // Alternative to Objective.KeyResults
}

// ToObjectives converts an OKRSet to a flat list of Objectives.
// Merges KeyResults from OKR struct into Objective.KeyResults.
func (s *OKRSet) ToObjectives() []Objective {
	objectives := make([]Objective, len(s.OKRs))
	for i, okr := range s.OKRs {
		obj := okr.Objective
		// Merge KeyResults if Objective.KeyResults is empty
		if len(obj.KeyResults) == 0 {
			obj.KeyResults = okr.KeyResults
		}
		objectives[i] = obj
	}
	return objectives
}

// FromObjectives creates an OKRSet from a list of Objectives.
func FromObjectives(objectives []Objective) *OKRSet {
	okrs := make([]OKR, len(objectives))
	for i, obj := range objectives {
		okrs[i] = OKR{
			Objective:  obj,
			KeyResults: obj.KeyResults,
		}
	}
	return &OKRSet{OKRs: okrs}
}
