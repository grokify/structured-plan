package common

// CustomSection allows project-specific sections.
// Used across PRD, MRD, and TRD documents.
type CustomSection struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Content     any    `json:"content"`          // Flexible content structure
	Schema      string `json:"schema,omitempty"` // Optional JSON schema for validation
}
