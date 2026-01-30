package prd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	doc := New("PRD-001", "Test PRD", Person{Name: "John Doe", Email: "john@example.com"})

	if doc.Metadata.ID != "PRD-001" {
		t.Errorf("ID = %s, want PRD-001", doc.Metadata.ID)
	}

	if doc.Metadata.Title != "Test PRD" {
		t.Errorf("Title = %s, want Test PRD", doc.Metadata.Title)
	}

	if doc.Metadata.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", doc.Metadata.Version)
	}

	if doc.Metadata.Status != StatusDraft {
		t.Errorf("Status = %s, want draft", doc.Metadata.Status)
	}

	if len(doc.Metadata.Authors) != 1 {
		t.Errorf("Authors count = %d, want 1", len(doc.Metadata.Authors))
	}

	if doc.Metadata.Authors[0].Name != "John Doe" {
		t.Errorf("Author name = %s, want John Doe", doc.Metadata.Authors[0].Name)
	}

	// Check revision history
	if len(doc.RevisionHistory) != 1 {
		t.Errorf("RevisionHistory count = %d, want 1", len(doc.RevisionHistory))
	}

	if doc.RevisionHistory[0].Version != "1.0.0" {
		t.Errorf("Initial revision version = %s, want 1.0.0", doc.RevisionHistory[0].Version)
	}

	if doc.RevisionHistory[0].Trigger != TriggerInitial {
		t.Errorf("Initial revision trigger = %s, want initial", doc.RevisionHistory[0].Trigger)
	}
}

func TestNewWithoutAuthors(t *testing.T) {
	doc := New("PRD-002", "Test PRD Without Authors")

	if len(doc.Metadata.Authors) != 0 {
		t.Errorf("Authors count = %d, want 0", len(doc.Metadata.Authors))
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateID()

	if !strings.HasPrefix(id, "PRD-") {
		t.Errorf("ID should start with 'PRD-', got %s", id)
	}

	// ID format should be PRD-YYYY-DDD
	parts := strings.Split(id, "-")
	if len(parts) != 3 {
		t.Errorf("ID should have format PRD-YYYY-DDD, got %s", id)
	}
}

func TestGenerateIDWithPrefix(t *testing.T) {
	id := GenerateIDWithPrefix("FEAT")

	if !strings.HasPrefix(id, "FEAT-") {
		t.Errorf("ID should start with 'FEAT-', got %s", id)
	}
}

func TestAddRevision(t *testing.T) {
	doc := New("PRD-001", "Test PRD")
	initialVersion := doc.Metadata.Version

	doc.AddRevision(
		[]string{"Added user stories", "Updated metrics"},
		TriggerReview,
		"Jane Doe",
	)

	// Version should be incremented
	if doc.Metadata.Version == initialVersion {
		t.Error("Version should be incremented")
	}

	// Revision history should have 2 entries
	if len(doc.RevisionHistory) != 2 {
		t.Errorf("RevisionHistory count = %d, want 2", len(doc.RevisionHistory))
	}

	// Check latest revision
	latest := doc.RevisionHistory[len(doc.RevisionHistory)-1]
	if latest.Author != "Jane Doe" {
		t.Errorf("Revision author = %s, want Jane Doe", latest.Author)
	}

	if latest.Trigger != TriggerReview {
		t.Errorf("Revision trigger = %s, want review", latest.Trigger)
	}

	if len(latest.Changes) != 2 {
		t.Errorf("Revision changes count = %d, want 2", len(latest.Changes))
	}
}

func TestIncrementVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.0", "1.0.1"},
		{"1.2.3", "1.2.4"},
		{"0.0.0", "0.0.1"},
		{"invalid", "invalid.1"},
		{"1.0", "1.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := incrementVersion(tt.input)
			if result != tt.expected {
				t.Errorf("incrementVersion(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temp directory for the test
	tempDir, err := os.MkdirTemp("", "prd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a document
	doc := New("PRD-TEST", "Test Document", Person{Name: "Test Author"})
	doc.ExecutiveSummary.ProblemStatement = "Test problem"
	doc.ExecutiveSummary.ProposedSolution = "Test solution"

	// Save it
	filePath := filepath.Join(tempDir, "test.json")
	err = Save(doc, filePath)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("File was not created")
	}

	// Load it back
	loaded, err := Load(filePath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify contents
	if loaded.Metadata.ID != doc.Metadata.ID {
		t.Errorf("Loaded ID = %s, want %s", loaded.Metadata.ID, doc.Metadata.ID)
	}

	if loaded.Metadata.Title != doc.Metadata.Title {
		t.Errorf("Loaded Title = %s, want %s", loaded.Metadata.Title, doc.Metadata.Title)
	}

	if loaded.ExecutiveSummary.ProblemStatement != doc.ExecutiveSummary.ProblemStatement {
		t.Errorf("Loaded ProblemStatement = %s, want %s",
			loaded.ExecutiveSummary.ProblemStatement, doc.ExecutiveSummary.ProblemStatement)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/file.json")
	if err == nil {
		t.Error("Load should fail for non-existent file")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	// Create temp file with invalid JSON
	tempFile, err := os.CreateTemp("", "invalid-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("{ invalid json }")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	_, err = Load(tempFile.Name())
	if err == nil {
		t.Error("Load should fail for invalid JSON")
	}
}

func TestDocumentToJSON(t *testing.T) {
	doc := New("PRD-001", "Test PRD")

	jsonStr, err := doc.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if !strings.Contains(jsonStr, "PRD-001") {
		t.Error("JSON should contain document ID")
	}

	if !strings.Contains(jsonStr, "Test PRD") {
		t.Error("JSON should contain document title")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "prd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	doc := New("PRD-001", "Test")
	nestedPath := filepath.Join(tempDir, "subdir", "nested", "test.json")

	err = Save(doc, nestedPath)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
		t.Error("File was not created in nested directory")
	}
}
