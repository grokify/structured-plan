package common

// GlossaryTerm defines a glossary entry.
// Used across PRD, MRD, and TRD documents.
type GlossaryTerm struct {
	Term       string   `json:"term"`
	Definition string   `json:"definition"`
	Acronym    string   `json:"acronym,omitempty"`
	Context    string   `json:"context,omitempty"`
	Related    []string `json:"related,omitempty"` // Related terms
}
