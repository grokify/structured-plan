package prd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// DefaultPersonaLibraryFilename is the standard filename for persona libraries.
const DefaultPersonaLibraryFilename = "personas.json"

// PersonaLibrary manages reusable personas across documents.
type PersonaLibrary struct {
	SchemaVersion string           `json:"schema_version"`
	Personas      []LibraryPersona `json:"personas"`
	Metadata      LibraryMetadata  `json:"metadata,omitempty"`
}

// LibraryPersona extends Persona with library-specific metadata.
type LibraryPersona struct {
	Persona
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	UsedInPRDs  []string  `json:"used_in_prds,omitempty"`        // Track which PRDs use this persona
	Tags        []string  `json:"tags,omitempty"`                // For organization/filtering
	Description string    `json:"library_description,omitempty"` // Additional context for library
}

// LibraryMetadata contains metadata about the persona library itself.
type LibraryMetadata struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// Common errors for persona library operations.
var (
	ErrPersonaNotFound      = errors.New("persona not found in library")
	ErrPersonaAlreadyExists = errors.New("persona with this ID already exists")
	ErrPersonaNameRequired  = errors.New("persona name is required")
	ErrInvalidPersonaID     = errors.New("invalid persona ID")
)

// NewPersonaLibrary creates a new empty persona library.
func NewPersonaLibrary() *PersonaLibrary {
	now := time.Now()
	return &PersonaLibrary{
		SchemaVersion: "1.0",
		Personas:      []LibraryPersona{},
		Metadata: LibraryMetadata{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// LoadPersonaLibrary reads a persona library from a JSON file.
func LoadPersonaLibrary(path string) (*PersonaLibrary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read persona library: %w", err)
	}

	var lib PersonaLibrary
	if err := json.Unmarshal(data, &lib); err != nil {
		return nil, fmt.Errorf("failed to parse persona library: %w", err)
	}

	return &lib, nil
}

// Save writes the persona library to a JSON file.
func (lib *PersonaLibrary) Save(path string) error {
	lib.Metadata.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(lib, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal persona library: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write persona library: %w", err)
	}

	return nil
}

// Get retrieves a persona by ID.
func (lib *PersonaLibrary) Get(id string) *LibraryPersona {
	for i := range lib.Personas {
		if lib.Personas[i].ID == id {
			return &lib.Personas[i]
		}
	}
	return nil
}

// GetByName retrieves a persona by name (case-sensitive).
func (lib *PersonaLibrary) GetByName(name string) *LibraryPersona {
	for i := range lib.Personas {
		if lib.Personas[i].Name == name {
			return &lib.Personas[i]
		}
	}
	return nil
}

// List returns all personas in the library.
func (lib *PersonaLibrary) List() []LibraryPersona {
	return lib.Personas
}

// ListByTag returns personas matching a specific tag.
func (lib *PersonaLibrary) ListByTag(tag string) []LibraryPersona {
	var result []LibraryPersona
	for _, p := range lib.Personas {
		for _, t := range p.Tags {
			if t == tag {
				result = append(result, p)
				break
			}
		}
	}
	return result
}

// Add adds a new persona to the library.
// Generates an ID if not provided. Returns the ID.
func (lib *PersonaLibrary) Add(p Persona) (string, error) {
	if p.Name == "" {
		return "", ErrPersonaNameRequired
	}

	// Generate ID if not provided
	if p.ID == "" {
		p.ID = lib.generatePersonaID()
	} else {
		// Check for duplicate ID
		if lib.Get(p.ID) != nil {
			return "", ErrPersonaAlreadyExists
		}
	}

	now := time.Now()
	libPersona := LibraryPersona{
		Persona:   p,
		CreatedAt: now,
		UpdatedAt: now,
	}

	lib.Personas = append(lib.Personas, libPersona)
	lib.Metadata.UpdatedAt = now

	return p.ID, nil
}

// Update updates an existing persona in the library.
func (lib *PersonaLibrary) Update(p Persona) error {
	if p.ID == "" {
		return ErrInvalidPersonaID
	}

	for i := range lib.Personas {
		if lib.Personas[i].ID == p.ID {
			// Preserve library metadata
			lib.Personas[i].Persona = p
			lib.Personas[i].UpdatedAt = time.Now()
			lib.Metadata.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrPersonaNotFound
}

// Remove removes a persona from the library by ID.
func (lib *PersonaLibrary) Remove(id string) error {
	for i := range lib.Personas {
		if lib.Personas[i].ID == id {
			lib.Personas = append(lib.Personas[:i], lib.Personas[i+1:]...)
			lib.Metadata.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrPersonaNotFound
}

// ImportTo imports a persona from the library into a document.
// The persona is copied (not referenced) so the document is self-contained.
// Sets LibraryRef to track the origin.
func (lib *PersonaLibrary) ImportTo(doc *Document, personaID string, isPrimary bool) error {
	libPersona := lib.Get(personaID)
	if libPersona == nil {
		return ErrPersonaNotFound
	}

	// Create a copy of the persona for the document
	docPersona := libPersona.Persona
	docPersona.IsPrimary = isPrimary
	docPersona.LibraryRef = personaID

	// Check if persona already exists in document
	for i, existing := range doc.Personas {
		if existing.ID == personaID || existing.LibraryRef == personaID {
			// Update existing
			doc.Personas[i] = docPersona
			return nil
		}
	}

	// If this is primary, unset primary on others
	if isPrimary {
		for i := range doc.Personas {
			doc.Personas[i].IsPrimary = false
		}
	}

	// Add new persona
	doc.Personas = append(doc.Personas, docPersona)

	// Track usage in library
	lib.trackUsage(personaID, doc.Metadata.ID)

	return nil
}

// ExportFrom exports a persona from a document to the library.
// If the persona already exists in the library (by ID), it's updated.
// If not, it's added as a new persona.
// Returns the library persona ID.
func (lib *PersonaLibrary) ExportFrom(doc *Document, personaID string) (string, error) {
	// Find persona in document
	var docPersona *Persona
	for i := range doc.Personas {
		if doc.Personas[i].ID == personaID {
			docPersona = &doc.Personas[i]
			break
		}
	}

	if docPersona == nil {
		return "", fmt.Errorf("persona %s not found in document", personaID)
	}

	// Check if persona already exists in library
	if existing := lib.Get(personaID); existing != nil {
		// Update existing
		if err := lib.Update(*docPersona); err != nil {
			return "", err
		}
		lib.trackUsage(personaID, doc.Metadata.ID)
		return personaID, nil
	}

	// Check if there's a library ref
	if docPersona.LibraryRef != "" {
		if existing := lib.Get(docPersona.LibraryRef); existing != nil {
			// Update the original library persona
			existingID := docPersona.LibraryRef
			updated := *docPersona
			updated.ID = existingID
			if err := lib.Update(updated); err != nil {
				return "", err
			}
			lib.trackUsage(existingID, doc.Metadata.ID)
			return existingID, nil
		}
	}

	// Add as new persona
	id, err := lib.Add(*docPersona)
	if err != nil {
		return "", err
	}

	// Update document persona with library ref
	for i := range doc.Personas {
		if doc.Personas[i].ID == personaID {
			doc.Personas[i].LibraryRef = id
			break
		}
	}

	lib.trackUsage(id, doc.Metadata.ID)
	return id, nil
}

// ExportAllFrom exports all personas from a document to the library.
// Returns the count of added and updated personas.
func (lib *PersonaLibrary) ExportAllFrom(doc *Document) (added int, updated int, errs []error) {
	for i := range doc.Personas {
		persona := &doc.Personas[i]
		wasNew := lib.Get(persona.ID) == nil && (persona.LibraryRef == "" || lib.Get(persona.LibraryRef) == nil)

		_, err := lib.ExportFrom(doc, persona.ID)
		if err != nil {
			errs = append(errs, fmt.Errorf("persona %s: %w", persona.ID, err))
			continue
		}

		if wasNew {
			added++
		} else {
			updated++
		}
	}
	return added, updated, errs
}

// SyncFromLibrary updates all personas in a document with their latest library versions.
// Only updates personas that have a LibraryRef set.
// Returns the count of updated personas.
func (lib *PersonaLibrary) SyncFromLibrary(doc *Document) int {
	updated := 0
	for i := range doc.Personas {
		ref := doc.Personas[i].LibraryRef
		if ref == "" {
			continue
		}

		libPersona := lib.Get(ref)
		if libPersona == nil {
			continue
		}

		// Preserve document-specific fields
		isPrimary := doc.Personas[i].IsPrimary

		// Update from library
		doc.Personas[i] = libPersona.Persona
		doc.Personas[i].IsPrimary = isPrimary
		doc.Personas[i].LibraryRef = ref
		updated++
	}
	return updated
}

// generatePersonaID generates a unique persona ID.
func (lib *PersonaLibrary) generatePersonaID() string {
	maxNum := 0
	for _, p := range lib.Personas {
		var num int
		if _, err := fmt.Sscanf(p.ID, "PER-%d", &num); err == nil {
			if num > maxNum {
				maxNum = num
			}
		}
	}
	return fmt.Sprintf("PER-%03d", maxNum+1)
}

// trackUsage adds a PRD ID to the persona's usage list if not already present.
func (lib *PersonaLibrary) trackUsage(personaID, prdID string) {
	if prdID == "" {
		return
	}

	persona := lib.Get(personaID)
	if persona == nil {
		return
	}

	// Check if already tracked
	for _, id := range persona.UsedInPRDs {
		if id == prdID {
			return
		}
	}

	persona.UsedInPRDs = append(persona.UsedInPRDs, prdID)
}

// Count returns the number of personas in the library.
func (lib *PersonaLibrary) Count() int {
	return len(lib.Personas)
}

// ToPersona converts a LibraryPersona to a basic Persona (drops library metadata).
func (lp *LibraryPersona) ToPersona() Persona {
	return lp.Persona
}
