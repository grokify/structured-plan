package prd

// Appendix represents a single appendix section.
// Content can be provided via ContentString, ContentTable, or both.
// When both are set, ContentString is rendered before ContentTable.
type Appendix struct {
	// ID is the unique identifier for this appendix.
	ID string `json:"id"`

	// Title is the appendix title.
	Title string `json:"title"`

	// Description provides context for this appendix.
	Description string `json:"description,omitempty"`

	// Type indicates the primary content type (hint for rendering).
	Type AppendixType `json:"type"`

	// ContentString is Markdown text content.
	// Rendered before ContentTable if both are set.
	ContentString string `json:"contentString,omitempty"`

	// ContentTable is structured table data.
	// Rendered after ContentString if both are set.
	ContentTable *AppendixTable `json:"contentTable,omitempty"`

	// Schema is the standard schema type (for validation and rendering hints).
	Schema AppendixSchema `json:"schema,omitempty"`

	// Tags for filtering and categorization.
	Tags []string `json:"tags,omitempty"`

	// ReferencedBy lists IDs of items that reference this appendix.
	// This is typically computed, not manually set.
	ReferencedBy []string `json:"referencedBy,omitempty"`
}

// AppendixTable represents tabular data.
type AppendixTable struct {
	// Headers are column headers.
	Headers []string `json:"headers,omitempty"`

	// Rows are table rows.
	Rows [][]string `json:"rows"`

	// Caption provides optional table description/footer.
	Caption string `json:"caption,omitempty"`
}

// AppendixType indicates the primary content type.
type AppendixType string

const (
	// AppendixTypeTable uses ContentTable for structured data.
	AppendixTypeTable AppendixType = "table"

	// AppendixTypeText uses ContentString for free-form text/markdown.
	AppendixTypeText AppendixType = "text"

	// AppendixTypeCode uses ContentString for code blocks.
	AppendixTypeCode AppendixType = "code"

	// AppendixTypeReference uses ContentString for reference links.
	AppendixTypeReference AppendixType = "reference"

	// AppendixTypeDiagram uses ContentString for diagram references or embedded diagrams.
	AppendixTypeDiagram AppendixType = "diagram"
)

// AppendixSchema identifies the schema type for appendices.
// This field can be used for validation hints or future plugin-based rendering.
// For now, use "custom" or any descriptive string for your use case.
type AppendixSchema string

const (
	// AppendixSchemaCustom is a custom/unstructured appendix (default).
	AppendixSchemaCustom AppendixSchema = "custom"
)

// NewTableAppendix creates a new table appendix with custom headers.
func NewTableAppendix(id, title, description string, headers []string) Appendix {
	return Appendix{
		ID:          id,
		Title:       title,
		Description: description,
		Type:        AppendixTypeTable,
		Schema:      AppendixSchemaCustom,
		ContentTable: &AppendixTable{
			Headers: headers,
			Rows:    [][]string{},
		},
	}
}

// NewTextAppendix creates a new text/markdown appendix.
func NewTextAppendix(id, title, description, content string) Appendix {
	return Appendix{
		ID:            id,
		Title:         title,
		Description:   description,
		Type:          AppendixTypeText,
		ContentString: content,
	}
}

// NewCodeAppendix creates a new code appendix.
func NewCodeAppendix(id, title, description, code string) Appendix {
	return Appendix{
		ID:            id,
		Title:         title,
		Description:   description,
		Type:          AppendixTypeCode,
		ContentString: code,
	}
}

// NewReferenceAppendix creates a new reference appendix for links and citations.
func NewReferenceAppendix(id, title, description, references string) Appendix {
	return Appendix{
		ID:            id,
		Title:         title,
		Description:   description,
		Type:          AppendixTypeReference,
		ContentString: references,
	}
}
