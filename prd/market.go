package prd

// MarketDefinition contains market analysis and competitive landscape.
type MarketDefinition struct {
	// Alternatives are competing products, workarounds, or alternative approaches.
	Alternatives []Alternative `json:"alternatives,omitempty"`

	// Differentiation describes how this solution differs from alternatives.
	Differentiation []string `json:"differentiation,omitempty"`

	// MarketRisks are risks related to market conditions or competition.
	MarketRisks []string `json:"market_risks,omitempty"`
}

// Alternative represents a competing product or alternative approach.
type Alternative struct {
	// ID is the unique identifier for this alternative.
	ID string `json:"id"`

	// Name is the name of the alternative.
	Name string `json:"name"`

	// Type categorizes the alternative.
	Type AlternativeType `json:"type"`

	// Description provides details about the alternative.
	Description string `json:"description,omitempty"`

	// Strengths are advantages of this alternative.
	Strengths []string `json:"strengths,omitempty"`

	// Weaknesses are disadvantages of this alternative.
	Weaknesses []string `json:"weaknesses,omitempty"`

	// WhyNotChosen explains why this alternative was not selected.
	WhyNotChosen string `json:"why_not_chosen,omitempty"`
}

// AlternativeType categorizes alternatives.
type AlternativeType string

const (
	// AlternativeCompetitor is a competing product.
	AlternativeCompetitor AlternativeType = "competitor"

	// AlternativeWorkaround is a manual or existing workaround.
	AlternativeWorkaround AlternativeType = "workaround"

	// AlternativeDoNothing represents the option to not address the problem.
	AlternativeDoNothing AlternativeType = "do_nothing"

	// AlternativeInternalTool is an existing internal solution.
	AlternativeInternalTool AlternativeType = "internal_tool"
)
