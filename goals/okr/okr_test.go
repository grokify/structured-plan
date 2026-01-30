package okr

import (
	"testing"
)

func TestNewOKRDocument(t *testing.T) {
	doc := New("OKR-2025-001", "Q1 OKRs", "Alice Smith")

	if doc.Metadata == nil {
		t.Fatal("Expected metadata to be initialized")
	}
	if doc.Metadata.ID != "OKR-2025-001" {
		t.Errorf("Expected ID 'OKR-2025-001', got '%s'", doc.Metadata.ID)
	}
	if doc.Metadata.Name != "Q1 OKRs" {
		t.Errorf("Expected name 'Q1 OKRs', got '%s'", doc.Metadata.Name)
	}
	if doc.Metadata.Owner != "Alice Smith" {
		t.Errorf("Expected owner 'Alice Smith', got '%s'", doc.Metadata.Owner)
	}
	if doc.Metadata.Status != StatusDraft {
		t.Errorf("Expected status 'Draft', got '%s'", doc.Metadata.Status)
	}
	if len(doc.Objectives) != 0 {
		t.Errorf("Expected empty objectives, got %d", len(doc.Objectives))
	}
}

func TestCalculateProgress(t *testing.T) {
	obj := Objective{
		Title: "Improve customer satisfaction",
		KeyResults: []KeyResult{
			{Title: "Increase NPS", Score: 0.8},
			{Title: "Reduce churn", Score: 0.6},
			{Title: "Improve response time", Score: 0.9},
		},
	}

	// Expected: (0.8 + 0.6 + 0.9) / 3 = 0.7666...
	progress := obj.CalculateProgress()
	expected := 0.7666666666666667

	if progress < expected-0.01 || progress > expected+0.01 {
		t.Errorf("Expected progress ~%.2f, got %.4f", expected, progress)
	}
}

func TestCalculateProgressEmpty(t *testing.T) {
	obj := Objective{
		Title:      "Empty objective",
		KeyResults: []KeyResult{},
	}

	progress := obj.CalculateProgress()
	if progress != 0 {
		t.Errorf("Expected 0 progress for empty key results, got %.2f", progress)
	}
}

func TestDocumentOverallProgress(t *testing.T) {
	doc := &OKRDocument{
		Objectives: []Objective{
			{
				Title: "Objective 1",
				KeyResults: []KeyResult{
					{Title: "KR 1.1", Score: 0.8},
					{Title: "KR 1.2", Score: 0.6},
				},
			},
			{
				Title: "Objective 2",
				KeyResults: []KeyResult{
					{Title: "KR 2.1", Score: 1.0},
					{Title: "KR 2.2", Score: 0.4},
				},
			},
		},
	}

	// Obj1: (0.8 + 0.6) / 2 = 0.7
	// Obj2: (1.0 + 0.4) / 2 = 0.7
	// Overall: (0.7 + 0.7) / 2 = 0.7
	progress := doc.CalculateOverallProgress()
	expected := 0.7

	if progress < expected-0.01 || progress > expected+0.01 {
		t.Errorf("Expected overall progress ~%.2f, got %.4f", expected, progress)
	}
}

func TestAllKeyResults(t *testing.T) {
	doc := &OKRDocument{
		Objectives: []Objective{
			{
				Title: "Objective 1",
				KeyResults: []KeyResult{
					{Title: "KR 1.1"},
					{Title: "KR 1.2"},
				},
			},
			{
				Title: "Objective 2",
				KeyResults: []KeyResult{
					{Title: "KR 2.1"},
				},
			},
		},
	}

	krs := doc.AllKeyResults()
	if len(krs) != 3 {
		t.Errorf("Expected 3 key results, got %d", len(krs))
	}
}

func TestAllRisks(t *testing.T) {
	doc := &OKRDocument{
		Risks: []Risk{
			{Title: "Global Risk"},
		},
		Objectives: []Objective{
			{
				Title: "Objective 1",
				Risks: []Risk{
					{Title: "Objective-specific Risk"},
				},
			},
		},
	}

	risks := doc.AllRisks()
	if len(risks) != 2 {
		t.Errorf("Expected 2 risks, got %d", len(risks))
	}
}

func TestScoreGrade(t *testing.T) {
	tests := []struct {
		score    float64
		expected string
	}{
		{1.0, "A"},
		{0.9, "A"},
		{0.85, "B"},
		{0.7, "B"},
		{0.5, "C"},
		{0.4, "C"},
		{0.3, "D"},
		{0.2, "D"},
		{0.1, "F"},
		{0.0, "F"},
	}

	for _, tt := range tests {
		grade := ScoreGrade(tt.score)
		if grade != tt.expected {
			t.Errorf("ScoreGrade(%.2f) = %s, expected %s", tt.score, grade, tt.expected)
		}
	}
}

func TestScoreDescription(t *testing.T) {
	tests := []struct {
		score    float64
		expected string
	}{
		{1.0, "Exceeded expectations"},
		{0.7, "Achieved target"},
		{0.5, "Partial achievement"},
		{0.2, "Below expectations"},
		{0.0, "Not achieved"},
	}

	for _, tt := range tests {
		desc := ScoreDescription(tt.score)
		if desc != tt.expected {
			t.Errorf("ScoreDescription(%.2f) = %s, expected %s", tt.score, desc, tt.expected)
		}
	}
}

func TestValidateBasic(t *testing.T) {
	// Valid minimal document
	doc := &OKRDocument{
		Objectives: []Objective{
			{
				Title: "Improve product quality",
				KeyResults: []KeyResult{
					{Title: "Reduce bugs by 50%", Target: "50%"},
				},
			},
		},
	}

	errs := doc.Validate(nil)
	errorCount := len(Errors(errs))
	if errorCount != 0 {
		t.Errorf("Expected no errors for valid document, got %d: %v", errorCount, errs)
	}
}

func TestValidateMissingObjectives(t *testing.T) {
	doc := &OKRDocument{
		Objectives: []Objective{},
	}

	errs := doc.Validate(nil)
	if len(Errors(errs)) == 0 {
		t.Error("Expected error for missing objectives")
	}
}

func TestValidateMissingTitle(t *testing.T) {
	doc := &OKRDocument{
		Objectives: []Objective{
			{
				Title: "", // Missing
				KeyResults: []KeyResult{
					{Title: "KR 1"},
				},
			},
		},
	}

	errs := doc.Validate(nil)
	hasError := false
	for _, e := range errs {
		if e.Path == "objectives[0].title" && e.IsError {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("Expected error for missing objective title")
	}
}

func TestValidateInvalidScore(t *testing.T) {
	doc := &OKRDocument{
		Objectives: []Objective{
			{
				Title: "Test objective",
				KeyResults: []KeyResult{
					{Title: "KR 1", Score: 1.5}, // Invalid: > 1.0
				},
			},
		},
	}

	opts := &ValidationOptions{
		ValidateScoreRange: true,
	}
	errs := doc.Validate(opts)
	hasError := false
	for _, e := range errs {
		if e.Path == "objectives[0].keyResults[0].score" && e.IsError {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("Expected error for invalid score > 1.0")
	}
}

func TestValidateTooManyObjectives(t *testing.T) {
	doc := &OKRDocument{
		Objectives: []Objective{
			{Title: "Obj 1", KeyResults: []KeyResult{{Title: "KR 1"}}},
			{Title: "Obj 2", KeyResults: []KeyResult{{Title: "KR 2"}}},
			{Title: "Obj 3", KeyResults: []KeyResult{{Title: "KR 3"}}},
			{Title: "Obj 4", KeyResults: []KeyResult{{Title: "KR 4"}}},
			{Title: "Obj 5", KeyResults: []KeyResult{{Title: "KR 5"}}},
			{Title: "Obj 6", KeyResults: []KeyResult{{Title: "KR 6"}}},
		},
	}

	opts := DefaultValidationOptions() // Max 5 objectives
	errs := doc.Validate(opts)

	// Should have a warning about too many objectives
	hasWarning := false
	for _, e := range errs {
		if e.Path == "objectives" && !e.IsError {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Error("Expected warning for too many objectives")
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateID()
	if len(id) < 12 { // "OKR-YYYY-DDD"
		t.Errorf("Generated ID too short: %s", id)
	}
	if id[:4] != "OKR-" {
		t.Errorf("Generated ID should start with 'OKR-': %s", id)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	original := &OKRDocument{
		Metadata: &Metadata{
			ID:     "OKR-TEST",
			Name:   "Test OKRs",
			Owner:  "Test Owner",
			Status: StatusActive,
		},
		Theme: "Growth",
		Objectives: []Objective{
			{
				ID:    "O1",
				Title: "Increase revenue",
				KeyResults: []KeyResult{
					{ID: "KR1", Title: "Grow ARR by 20%", Target: "$1M", Score: 0.75},
				},
			},
		},
	}

	// Convert to JSON
	data, err := original.JSON()
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Parse back
	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify
	if parsed.Metadata.ID != original.Metadata.ID {
		t.Errorf("ID mismatch: got %s, want %s", parsed.Metadata.ID, original.Metadata.ID)
	}
	if parsed.Theme != original.Theme {
		t.Errorf("Theme mismatch: got %s, want %s", parsed.Theme, original.Theme)
	}
	if len(parsed.Objectives) != len(original.Objectives) {
		t.Errorf("Objectives count mismatch: got %d, want %d", len(parsed.Objectives), len(original.Objectives))
	}
	if parsed.Objectives[0].KeyResults[0].Score != original.Objectives[0].KeyResults[0].Score {
		t.Errorf("Score mismatch: got %.2f, want %.2f",
			parsed.Objectives[0].KeyResults[0].Score,
			original.Objectives[0].KeyResults[0].Score)
	}
}
