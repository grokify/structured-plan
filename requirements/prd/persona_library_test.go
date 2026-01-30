package prd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPersonaLibrary(t *testing.T) {
	lib := NewPersonaLibrary()

	if lib.SchemaVersion != "1.0" {
		t.Errorf("expected schema version 1.0, got %s", lib.SchemaVersion)
	}
	if len(lib.Personas) != 0 {
		t.Errorf("expected 0 personas, got %d", len(lib.Personas))
	}
	if lib.Metadata.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestPersonaLibraryAddAndGet(t *testing.T) {
	lib := NewPersonaLibrary()

	persona := Persona{
		Name:       "Developer Dan",
		Role:       "Backend Developer",
		PainPoints: []string{"Slow builds", "Complex configs"},
	}

	id, err := lib.Add(persona)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Verify ID was generated with correct format
	if id != "PER-001" {
		t.Errorf("expected ID PER-001, got %s", id)
	}

	// Get by ID
	retrieved := lib.Get(id)
	if retrieved == nil {
		t.Fatal("Get returned nil")
	}

	if retrieved.Name != "Developer Dan" {
		t.Errorf("expected name 'Developer Dan', got %s", retrieved.Name)
	}
	if retrieved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestPersonaLibraryAddWithID(t *testing.T) {
	lib := NewPersonaLibrary()

	persona := Persona{
		ID:   "CUSTOM-001",
		Name: "Custom Persona",
		Role: "Tester",
	}

	id, err := lib.Add(persona)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if id != "CUSTOM-001" {
		t.Errorf("expected ID CUSTOM-001, got %s", id)
	}

	// Try to add duplicate ID
	_, err = lib.Add(persona)
	if err != ErrPersonaAlreadyExists {
		t.Errorf("expected ErrPersonaAlreadyExists, got %v", err)
	}
}

func TestPersonaLibraryAddNoName(t *testing.T) {
	lib := NewPersonaLibrary()

	persona := Persona{
		Role: "No Name",
	}

	_, err := lib.Add(persona)
	if err != ErrPersonaNameRequired {
		t.Errorf("expected ErrPersonaNameRequired, got %v", err)
	}
}

func TestPersonaLibraryGetByName(t *testing.T) {
	lib := NewPersonaLibrary()

	_, _ = lib.Add(Persona{Name: "Developer Dan", Role: "Developer"})
	_, _ = lib.Add(Persona{Name: "Manager Maria", Role: "Manager"})

	found := lib.GetByName("Manager Maria")
	if found == nil {
		t.Fatal("GetByName returned nil")
	}
	if found.Role != "Manager" {
		t.Errorf("expected role 'Manager', got %s", found.Role)
	}

	notFound := lib.GetByName("Unknown")
	if notFound != nil {
		t.Error("expected nil for unknown name")
	}
}

func TestPersonaLibraryList(t *testing.T) {
	lib := NewPersonaLibrary()

	_, _ = lib.Add(Persona{Name: "Persona 1", Role: "Role 1"})
	_, _ = lib.Add(Persona{Name: "Persona 2", Role: "Role 2"})
	_, _ = lib.Add(Persona{Name: "Persona 3", Role: "Role 3"})

	personas := lib.List()
	if len(personas) != 3 {
		t.Errorf("expected 3 personas, got %d", len(personas))
	}

	if lib.Count() != 3 {
		t.Errorf("Count() expected 3, got %d", lib.Count())
	}
}

func TestPersonaLibraryListByTag(t *testing.T) {
	lib := NewPersonaLibrary()

	id1, _ := lib.Add(Persona{Name: "Dev 1", Role: "Developer"})
	id2, _ := lib.Add(Persona{Name: "Dev 2", Role: "Developer"})
	_, _ = lib.Add(Persona{Name: "Manager", Role: "Manager"})

	// Add tags
	lib.Get(id1).Tags = []string{"engineering", "backend"}
	lib.Get(id2).Tags = []string{"engineering", "frontend"}

	engineering := lib.ListByTag("engineering")
	if len(engineering) != 2 {
		t.Errorf("expected 2 engineering personas, got %d", len(engineering))
	}

	backend := lib.ListByTag("backend")
	if len(backend) != 1 {
		t.Errorf("expected 1 backend persona, got %d", len(backend))
	}

	unknown := lib.ListByTag("unknown")
	if len(unknown) != 0 {
		t.Errorf("expected 0 unknown personas, got %d", len(unknown))
	}
}

func TestPersonaLibraryUpdate(t *testing.T) {
	lib := NewPersonaLibrary()

	id, _ := lib.Add(Persona{Name: "Original", Role: "Developer"})

	updated := Persona{
		ID:   id,
		Name: "Updated",
		Role: "Senior Developer",
	}

	err := lib.Update(updated)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved := lib.Get(id)
	if retrieved.Name != "Updated" {
		t.Errorf("expected name 'Updated', got %s", retrieved.Name)
	}
	if retrieved.Role != "Senior Developer" {
		t.Errorf("expected role 'Senior Developer', got %s", retrieved.Role)
	}
	if retrieved.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestPersonaLibraryUpdateNotFound(t *testing.T) {
	lib := NewPersonaLibrary()

	err := lib.Update(Persona{ID: "NONEXISTENT", Name: "Test"})
	if err != ErrPersonaNotFound {
		t.Errorf("expected ErrPersonaNotFound, got %v", err)
	}
}

func TestPersonaLibraryUpdateNoID(t *testing.T) {
	lib := NewPersonaLibrary()

	err := lib.Update(Persona{Name: "No ID"})
	if err != ErrInvalidPersonaID {
		t.Errorf("expected ErrInvalidPersonaID, got %v", err)
	}
}

func TestPersonaLibraryRemove(t *testing.T) {
	lib := NewPersonaLibrary()

	id, _ := lib.Add(Persona{Name: "To Remove", Role: "Test"})

	if lib.Count() != 1 {
		t.Errorf("expected 1 persona before remove, got %d", lib.Count())
	}

	err := lib.Remove(id)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if lib.Count() != 0 {
		t.Errorf("expected 0 personas after remove, got %d", lib.Count())
	}

	if lib.Get(id) != nil {
		t.Error("expected nil after remove")
	}
}

func TestPersonaLibraryRemoveNotFound(t *testing.T) {
	lib := NewPersonaLibrary()

	err := lib.Remove("NONEXISTENT")
	if err != ErrPersonaNotFound {
		t.Errorf("expected ErrPersonaNotFound, got %v", err)
	}
}

func TestPersonaLibrarySaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "personas.json")

	// Create and save library
	lib := NewPersonaLibrary()
	lib.Metadata.Name = "Test Library"
	_, _ = lib.Add(Persona{Name: "Developer Dan", Role: "Developer", PainPoints: []string{"Slow builds"}})
	_, _ = lib.Add(Persona{Name: "Manager Maria", Role: "Manager"})

	err := lib.Save(path)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("library file was not created")
	}

	// Load library
	loaded, err := LoadPersonaLibrary(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Metadata.Name != "Test Library" {
		t.Errorf("expected name 'Test Library', got %s", loaded.Metadata.Name)
	}
	if loaded.Count() != 2 {
		t.Errorf("expected 2 personas, got %d", loaded.Count())
	}

	dan := loaded.GetByName("Developer Dan")
	if dan == nil {
		t.Fatal("expected to find Developer Dan")
	}
	if len(dan.PainPoints) != 1 || dan.PainPoints[0] != "Slow builds" {
		t.Errorf("expected pain point 'Slow builds', got %v", dan.PainPoints)
	}
}

func TestPersonaLibraryLoadNonExistent(t *testing.T) {
	_, err := LoadPersonaLibrary("/nonexistent/path/personas.json")
	if err == nil {
		t.Error("expected error loading non-existent file")
	}
}

func TestPersonaLibraryImportTo(t *testing.T) {
	lib := NewPersonaLibrary()
	id, _ := lib.Add(Persona{
		Name:       "Developer Dan",
		Role:       "Developer",
		PainPoints: []string{"Slow builds"},
	})

	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})

	// Import persona
	err := lib.ImportTo(doc, id, true)
	if err != nil {
		t.Fatalf("ImportTo failed: %v", err)
	}

	if len(doc.Personas) != 1 {
		t.Errorf("expected 1 persona in doc, got %d", len(doc.Personas))
	}

	persona := doc.Personas[0]
	if persona.Name != "Developer Dan" {
		t.Errorf("expected name 'Developer Dan', got %s", persona.Name)
	}
	if !persona.IsPrimary {
		t.Error("expected persona to be primary")
	}
	if persona.LibraryRef != id {
		t.Errorf("expected LibraryRef %s, got %s", id, persona.LibraryRef)
	}
}

func TestPersonaLibraryImportToUpdatesExisting(t *testing.T) {
	lib := NewPersonaLibrary()
	id, _ := lib.Add(Persona{
		Name: "Developer Dan",
		Role: "Developer",
	})

	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})

	// Import twice
	_ = lib.ImportTo(doc, id, false)
	_ = lib.ImportTo(doc, id, true)

	// Should still have only 1 persona
	if len(doc.Personas) != 1 {
		t.Errorf("expected 1 persona after double import, got %d", len(doc.Personas))
	}

	// Should be updated to primary
	if !doc.Personas[0].IsPrimary {
		t.Error("expected persona to be primary after second import")
	}
}

func TestPersonaLibraryImportToNotFound(t *testing.T) {
	lib := NewPersonaLibrary()
	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})

	err := lib.ImportTo(doc, "NONEXISTENT", true)
	if err != ErrPersonaNotFound {
		t.Errorf("expected ErrPersonaNotFound, got %v", err)
	}
}

func TestPersonaLibraryImportToTracksUsage(t *testing.T) {
	lib := NewPersonaLibrary()
	id, _ := lib.Add(Persona{Name: "Developer Dan", Role: "Developer"})

	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})
	_ = lib.ImportTo(doc, id, true)

	persona := lib.Get(id)
	if len(persona.UsedInPRDs) != 1 {
		t.Errorf("expected 1 PRD in usage, got %d", len(persona.UsedInPRDs))
	}
	if persona.UsedInPRDs[0] != "PRD-2026-001" {
		t.Errorf("expected PRD-2026-001 in usage, got %s", persona.UsedInPRDs[0])
	}
}

func TestPersonaLibraryExportFrom(t *testing.T) {
	lib := NewPersonaLibrary()
	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})

	// Add persona to doc
	doc.Personas = append(doc.Personas, Persona{
		ID:         "DOC-PER-001",
		Name:       "New Persona",
		Role:       "Tester",
		PainPoints: []string{"Manual testing"},
	})

	// Export to library
	id, err := lib.ExportFrom(doc, "DOC-PER-001")
	if err != nil {
		t.Fatalf("ExportFrom failed: %v", err)
	}

	if id != "DOC-PER-001" {
		t.Errorf("expected ID DOC-PER-001, got %s", id)
	}

	if lib.Count() != 1 {
		t.Errorf("expected 1 persona in library, got %d", lib.Count())
	}

	persona := lib.Get(id)
	if persona.Name != "New Persona" {
		t.Errorf("expected name 'New Persona', got %s", persona.Name)
	}

	// Check that LibraryRef was set in doc
	if doc.Personas[0].LibraryRef != id {
		t.Errorf("expected LibraryRef to be set, got %s", doc.Personas[0].LibraryRef)
	}
}

func TestPersonaLibraryExportFromUpdatesExisting(t *testing.T) {
	lib := NewPersonaLibrary()
	id, _ := lib.Add(Persona{
		Name: "Original Name",
		Role: "Original Role",
	})

	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})
	doc.Personas = append(doc.Personas, Persona{
		ID:   id,
		Name: "Updated Name",
		Role: "Updated Role",
	})

	// Export should update existing
	exportedID, err := lib.ExportFrom(doc, id)
	if err != nil {
		t.Fatalf("ExportFrom failed: %v", err)
	}

	if exportedID != id {
		t.Errorf("expected same ID, got %s", exportedID)
	}

	if lib.Count() != 1 {
		t.Errorf("expected 1 persona (not duplicated), got %d", lib.Count())
	}

	persona := lib.Get(id)
	if persona.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got %s", persona.Name)
	}
}

func TestPersonaLibraryExportFromNotFound(t *testing.T) {
	lib := NewPersonaLibrary()
	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})

	_, err := lib.ExportFrom(doc, "NONEXISTENT")
	if err == nil {
		t.Error("expected error for non-existent persona")
	}
}

func TestPersonaLibraryExportAllFrom(t *testing.T) {
	lib := NewPersonaLibrary()

	// Pre-add one persona
	existingID, _ := lib.Add(Persona{Name: "Existing", Role: "Existing"})

	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})
	doc.Personas = append(doc.Personas,
		Persona{ID: existingID, Name: "Existing Updated", Role: "Updated"},
		Persona{ID: "NEW-001", Name: "New Persona", Role: "New"},
	)

	added, updated, errs := lib.ExportAllFrom(doc)

	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
	if added != 1 {
		t.Errorf("expected 1 added, got %d", added)
	}
	if updated != 1 {
		t.Errorf("expected 1 updated, got %d", updated)
	}
	if lib.Count() != 2 {
		t.Errorf("expected 2 personas in library, got %d", lib.Count())
	}
}

func TestPersonaLibrarySyncFromLibrary(t *testing.T) {
	lib := NewPersonaLibrary()
	id, _ := lib.Add(Persona{
		Name:       "Developer Dan",
		Role:       "Developer",
		PainPoints: []string{"Updated pain point"},
	})

	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})
	doc.Personas = append(doc.Personas, Persona{
		ID:         "DOC-001",
		Name:       "Old Name",
		Role:       "Old Role",
		LibraryRef: id,
		IsPrimary:  true,
	})

	updated := lib.SyncFromLibrary(doc)

	if updated != 1 {
		t.Errorf("expected 1 updated, got %d", updated)
	}

	persona := doc.Personas[0]
	if persona.Name != "Developer Dan" {
		t.Errorf("expected name 'Developer Dan', got %s", persona.Name)
	}
	if !persona.IsPrimary {
		t.Error("expected IsPrimary to be preserved")
	}
	if persona.LibraryRef != id {
		t.Errorf("expected LibraryRef to be preserved, got %s", persona.LibraryRef)
	}
}

func TestPersonaLibrarySyncSkipsUnreferenced(t *testing.T) {
	lib := NewPersonaLibrary()
	_, _ = lib.Add(Persona{Name: "Library Persona", Role: "Test"})

	doc := New("PRD-2026-001", "Test PRD", Person{Name: "Owner"})
	doc.Personas = append(doc.Personas, Persona{
		ID:   "DOC-001",
		Name: "Doc Persona",
		Role: "Original",
		// No LibraryRef
	})

	updated := lib.SyncFromLibrary(doc)

	if updated != 0 {
		t.Errorf("expected 0 updated (no LibraryRef), got %d", updated)
	}

	// Should be unchanged
	if doc.Personas[0].Name != "Doc Persona" {
		t.Error("persona without LibraryRef should not be modified")
	}
}

func TestLibraryPersonaToPersona(t *testing.T) {
	lib := NewPersonaLibrary()
	id, _ := lib.Add(Persona{
		Name:       "Developer Dan",
		Role:       "Developer",
		PainPoints: []string{"Slow builds"},
	})

	libPersona := lib.Get(id)
	persona := libPersona.ToPersona()

	if persona.Name != "Developer Dan" {
		t.Errorf("expected name 'Developer Dan', got %s", persona.Name)
	}
	if persona.ID != id {
		t.Errorf("expected ID %s, got %s", id, persona.ID)
	}
}

func TestPersonaLibraryGenerateIDIncrementing(t *testing.T) {
	lib := NewPersonaLibrary()

	id1, _ := lib.Add(Persona{Name: "P1", Role: "R1"})
	id2, _ := lib.Add(Persona{Name: "P2", Role: "R2"})
	id3, _ := lib.Add(Persona{Name: "P3", Role: "R3"})

	if id1 != "PER-001" {
		t.Errorf("expected PER-001, got %s", id1)
	}
	if id2 != "PER-002" {
		t.Errorf("expected PER-002, got %s", id2)
	}
	if id3 != "PER-003" {
		t.Errorf("expected PER-003, got %s", id3)
	}
}
