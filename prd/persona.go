package prd

// TechnicalProficiency represents a user's technical skill level.
type TechnicalProficiency string

const (
	ProficiencyLow    TechnicalProficiency = "low"
	ProficiencyMedium TechnicalProficiency = "medium"
	ProficiencyHigh   TechnicalProficiency = "high"
	ProficiencyExpert TechnicalProficiency = "expert"
)

// Persona represents a user persona for the product.
type Persona struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`                // e.g., "Developer Dan"
	Role                 string               `json:"role"`                // Job title
	Description          string               `json:"description"`         // Background and context
	Goals                []string             `json:"goals"`               // What they want to achieve
	PainPoints           []string             `json:"pain_points"`         // Current frustrations
	Behaviors            []string             `json:"behaviors,omitempty"` // Typical patterns
	TechnicalProficiency TechnicalProficiency `json:"technical_proficiency,omitempty"`
	Demographics         *Demographics        `json:"demographics,omitempty"`
	Motivations          []string             `json:"motivations,omitempty"`
	Frustrations         []string             `json:"frustrations,omitempty"`
	PreferredChannels    []string             `json:"preferred_channels,omitempty"` // How they prefer to interact
	Quote                string               `json:"quote,omitempty"`              // Representative quote
	ImageURL             string               `json:"image_url,omitempty"`
	IsPrimary            bool                 `json:"is_primary,omitempty"`  // Is this the primary persona?
	LibraryRef           string               `json:"library_ref,omitempty"` // Reference to persona in library (for tracking origin)
}

// Demographics contains optional demographic information.
type Demographics struct {
	AgeRange    string `json:"age_range,omitempty"`
	Location    string `json:"location,omitempty"`
	Industry    string `json:"industry,omitempty"`
	CompanySize string `json:"company_size,omitempty"`
	Experience  string `json:"experience,omitempty"` // Years of experience
}
