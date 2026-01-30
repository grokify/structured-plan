package mrd

import (
	"testing"
)

func TestFilterByTags_EmptyTags(t *testing.T) {
	doc := Document{
		MarketRequirements: []MarketRequirement{
			{ID: "mr1", Title: "Req 1", Tags: []string{"data-management"}},
		},
	}
	filtered := doc.FilterByTags()
	if len(filtered.MarketRequirements) != 1 {
		t.Errorf("expected 1 market requirement, got %d", len(filtered.MarketRequirements))
	}
}

func TestFilterByTags_SingleTag(t *testing.T) {
	doc := Document{
		TargetMarket: TargetMarket{
			PrimarySegments: []MarketSegment{
				{ID: "s1", Name: "Segment 1", Tags: []string{"enterprise"}},
				{ID: "s2", Name: "Segment 2", Tags: []string{"smb"}},
			},
			BuyerPersonas: []BuyerPersona{
				{ID: "bp1", Name: "CTO", Tags: []string{"enterprise"}},
				{ID: "bp2", Name: "Developer", Tags: []string{"developer"}},
			},
		},
		MarketRequirements: []MarketRequirement{
			{ID: "mr1", Title: "Req 1", Tags: []string{"enterprise"}},
			{ID: "mr2", Title: "Req 2", Tags: []string{"smb"}},
		},
	}

	filtered := doc.FilterByTags("enterprise")

	if len(filtered.TargetMarket.PrimarySegments) != 1 {
		t.Errorf("expected 1 primary segment, got %d", len(filtered.TargetMarket.PrimarySegments))
	}
	if len(filtered.TargetMarket.BuyerPersonas) != 1 {
		t.Errorf("expected 1 buyer persona, got %d", len(filtered.TargetMarket.BuyerPersonas))
	}
	if len(filtered.MarketRequirements) != 1 {
		t.Errorf("expected 1 market requirement, got %d", len(filtered.MarketRequirements))
	}
}

func TestFilterByTags_MultipleTags_ORLogic(t *testing.T) {
	doc := Document{
		CompetitiveLandscape: CompetitiveLandscape{
			Competitors: []Competitor{
				{ID: "c1", Name: "Competitor 1", Tags: []string{"direct"}},
				{ID: "c2", Name: "Competitor 2", Tags: []string{"indirect"}},
				{ID: "c3", Name: "Competitor 3", Tags: []string{"substitute"}},
			},
		},
	}

	// OR logic: should return competitors with direct OR indirect
	filtered := doc.FilterByTags("direct", "indirect")

	if len(filtered.CompetitiveLandscape.Competitors) != 2 {
		t.Errorf("expected 2 competitors (direct OR indirect), got %d", len(filtered.CompetitiveLandscape.Competitors))
	}
}

func TestFilterByTags_Risks(t *testing.T) {
	doc := Document{
		Risks: []Risk{
			{ID: "r1", Description: "Risk 1", Tags: []string{"regulatory"}},
			{ID: "r2", Description: "Risk 2", Tags: []string{"competitive"}},
			{ID: "r3", Description: "Risk 3", Tags: []string{"market", "regulatory"}},
		},
	}

	filtered := doc.FilterByTags("regulatory")

	if len(filtered.Risks) != 2 {
		t.Errorf("expected 2 risks with regulatory tag, got %d", len(filtered.Risks))
	}
}
