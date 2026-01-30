package mrd

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestDocumentParsing tests JSON unmarshaling of MRD documents.
func TestDocumentParsing(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(t *testing.T, doc Document)
	}{
		{
			name: "minimal valid document",
			json: `{
				"metadata": {
					"id": "mrd-001",
					"title": "Test MRD",
					"version": "1.0.0",
					"status": "draft",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Test Author"}]
				},
				"executiveSummary": {
					"marketOpportunity": "Big opportunity",
					"proposedOffering": "Great product",
					"keyFindings": ["Finding 1"]
				},
				"marketOverview": {
					"tam": {"value": "$10B"},
					"sam": {"value": "$5B"},
					"som": {"value": "$1B"}
				},
				"targetMarket": {
					"primarySegments": [{"id": "seg1", "name": "Enterprise", "description": "Large companies"}]
				},
				"competitiveLandscape": {
					"overview": "Competitive market",
					"competitors": [{"id": "c1", "name": "Competitor 1", "strengths": ["Strong brand"], "weaknesses": ["High price"]}]
				},
				"marketRequirements": [{"id": "mr1", "title": "Feature X", "description": "Need feature X", "priority": "must"}],
				"positioning": {
					"statement": "For enterprises who need X",
					"targetAudience": "Enterprise IT",
					"category": "SaaS",
					"keyBenefits": ["Benefit 1"],
					"differentiators": ["Diff 1"]
				},
				"successMetrics": [{"id": "sm1", "name": "Revenue", "description": "ARR", "metric": "ARR", "target": "$1M"}]
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if doc.Metadata.ID != "mrd-001" {
					t.Errorf("expected ID 'mrd-001', got %q", doc.Metadata.ID)
				}
				if doc.Metadata.Title != "Test MRD" {
					t.Errorf("expected Title 'Test MRD', got %q", doc.Metadata.Title)
				}
				if doc.MarketOverview.TAM.Value != "$10B" {
					t.Errorf("expected TAM '$10B', got %q", doc.MarketOverview.TAM.Value)
				}
				if len(doc.TargetMarket.PrimarySegments) != 1 {
					t.Errorf("expected 1 primary segment, got %d", len(doc.TargetMarket.PrimarySegments))
				}
				if len(doc.CompetitiveLandscape.Competitors) != 1 {
					t.Errorf("expected 1 competitor, got %d", len(doc.CompetitiveLandscape.Competitors))
				}
			},
		},
		{
			name: "document with buyer personas",
			json: `{
				"metadata": {
					"id": "mrd-002",
					"title": "Test",
					"version": "1.0.0",
					"status": "draft",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author"}]
				},
				"executiveSummary": {"marketOpportunity": "M", "proposedOffering": "P", "keyFindings": []},
				"marketOverview": {"tam": {"value": "$1B"}, "sam": {"value": "$500M"}, "som": {"value": "$100M"}},
				"targetMarket": {
					"primarySegments": [{"id": "s1", "name": "Segment", "description": "Desc"}],
					"buyerPersonas": [
						{
							"id": "bp1",
							"name": "CTO",
							"title": "Chief Technology Officer",
							"description": "Tech leader",
							"buyingRole": "Decision Maker",
							"budgetAuthority": true,
							"painPoints": ["Pain 1"],
							"goals": ["Goal 1"]
						}
					]
				},
				"competitiveLandscape": {"overview": "O", "competitors": [{"id": "c1", "name": "C1", "strengths": ["S"], "weaknesses": ["W"]}]},
				"marketRequirements": [{"id": "mr1", "title": "T", "description": "D", "priority": "must"}],
				"positioning": {"statement": "S", "targetAudience": "T", "category": "C", "keyBenefits": ["B"], "differentiators": ["D"]},
				"successMetrics": []
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if len(doc.TargetMarket.BuyerPersonas) != 1 {
					t.Errorf("expected 1 buyer persona, got %d", len(doc.TargetMarket.BuyerPersonas))
				}
				bp := doc.TargetMarket.BuyerPersonas[0]
				if bp.BuyingRole != "Decision Maker" {
					t.Errorf("expected buying role 'Decision Maker', got %q", bp.BuyingRole)
				}
				if !bp.BudgetAuthority {
					t.Error("expected budget authority to be true")
				}
			},
		},
		{
			name: "document with go-to-market",
			json: `{
				"metadata": {
					"id": "mrd-003",
					"title": "Test",
					"version": "1.0.0",
					"status": "draft",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author"}]
				},
				"executiveSummary": {"marketOpportunity": "M", "proposedOffering": "P", "keyFindings": []},
				"marketOverview": {"tam": {"value": "$1B"}, "sam": {"value": "$500M"}, "som": {"value": "$100M"}},
				"targetMarket": {"primarySegments": [{"id": "s1", "name": "S", "description": "D"}]},
				"competitiveLandscape": {"overview": "O", "competitors": [{"id": "c1", "name": "C", "strengths": ["S"], "weaknesses": ["W"]}]},
				"marketRequirements": [{"id": "mr1", "title": "T", "description": "D", "priority": "must"}],
				"positioning": {"statement": "S", "targetAudience": "T", "category": "C", "keyBenefits": ["B"], "differentiators": ["D"]},
				"goToMarket": {
					"launchStrategy": "Phased rollout",
					"pricingStrategy": {
						"model": "Subscription",
						"tiers": [
							{"name": "Free", "price": "$0", "features": ["Basic"]},
							{"name": "Pro", "price": "$99/mo", "features": ["Advanced"]}
						]
					},
					"distributionChannels": ["Direct", "Partners"]
				},
				"successMetrics": []
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if doc.GoToMarket == nil {
					t.Fatal("expected go_to_market to be present")
				}
				if doc.GoToMarket.LaunchStrategy != "Phased rollout" {
					t.Errorf("expected launch strategy 'Phased rollout', got %q", doc.GoToMarket.LaunchStrategy)
				}
				if doc.GoToMarket.PricingStrategy == nil {
					t.Fatal("expected pricing_strategy to be present")
				}
				if len(doc.GoToMarket.PricingStrategy.Tiers) != 2 {
					t.Errorf("expected 2 pricing tiers, got %d", len(doc.GoToMarket.PricingStrategy.Tiers))
				}
			},
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc Document
			err := json.Unmarshal([]byte(tt.json), &doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, doc)
			}
		})
	}
}

// TestDocumentMarshaling tests JSON marshaling of MRD documents.
func TestDocumentMarshaling(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:        "mrd-test",
			Title:     "Test Document",
			Version:   "1.0.0",
			Status:    StatusDraft,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Authors:   []Person{{Name: "Test Author"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			MarketOpportunity: "Big market",
			ProposedOffering:  "Great product",
			KeyFindings:       []string{"Finding 1"},
		},
		MarketOverview: MarketOverview{
			TAM: MarketSize{Value: "$10B"},
			SAM: MarketSize{Value: "$5B"},
			SOM: MarketSize{Value: "$1B"},
		},
		TargetMarket: TargetMarket{
			PrimarySegments: []MarketSegment{{ID: "s1", Name: "Enterprise", Description: "Large cos"}},
		},
		CompetitiveLandscape: CompetitiveLandscape{
			Overview:    "Competitive",
			Competitors: []Competitor{{ID: "c1", Name: "Competitor", Strengths: []string{"S"}, Weaknesses: []string{"W"}}},
		},
		MarketRequirements: []MarketRequirement{{ID: "mr1", Title: "Feature", Description: "Desc", Priority: PriorityMust}},
		Positioning: Positioning{
			Statement:       "Position statement",
			TargetAudience:  "Enterprise",
			Category:        "SaaS",
			KeyBenefits:     []string{"Benefit"},
			Differentiators: []string{"Diff"},
		},
		SuccessMetrics: []SuccessMetric{{ID: "sm1", Name: "Revenue", Metric: "ARR", Target: "$1M"}},
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var doc2 Document
	if err := json.Unmarshal(data, &doc2); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if doc2.Metadata.ID != doc.Metadata.ID {
		t.Errorf("round-trip ID mismatch: got %q, want %q", doc2.Metadata.ID, doc.Metadata.ID)
	}
	if doc2.MarketOverview.TAM.Value != doc.MarketOverview.TAM.Value {
		t.Errorf("round-trip TAM mismatch: got %q, want %q", doc2.MarketOverview.TAM.Value, doc.MarketOverview.TAM.Value)
	}
}

// TestStatusConstants verifies status constant values.
func TestStatusConstants(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{StatusDraft, "draft"},
		{StatusInReview, "in_review"},
		{StatusApproved, "approved"},
		{StatusDeprecated, "deprecated"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("Status = %q, want %q", tt.status, tt.want)
			}
		})
	}
}

// TestPriorityConstants verifies priority constant values.
func TestPriorityConstants(t *testing.T) {
	tests := []struct {
		priority Priority
		want     string
	}{
		{PriorityMust, "must"},
		{PriorityShould, "should"},
		{PriorityCould, "could"},
		{PriorityWont, "wont"},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			if string(tt.priority) != tt.want {
				t.Errorf("Priority = %q, want %q", tt.priority, tt.want)
			}
		})
	}
}

// TestMarkdownGeneration tests the ToMarkdown method.
func TestMarkdownGeneration(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:        "mrd-test",
			Title:     "Test Market Requirements",
			Version:   "1.0.0",
			Status:    StatusDraft,
			CreatedAt: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			Authors:   []Person{{Name: "Jane Doe", Role: "Product Marketing"}},
			Tags:      []string{"market", "analysis"},
		},
		ExecutiveSummary: ExecutiveSummary{
			MarketOpportunity: "Growing market opportunity in AI.",
			ProposedOffering:  "AI-powered analytics platform.",
			KeyFindings:       []string{"Market growing 40% YoY", "No dominant player"},
			Recommendation:    "Launch MVP in Q2",
		},
		MarketOverview: MarketOverview{
			TAM:        MarketSize{Value: "$50B", Year: 2025, Source: "Gartner"},
			SAM:        MarketSize{Value: "$10B"},
			SOM:        MarketSize{Value: "$500M"},
			GrowthRate: "40% CAGR",
		},
		TargetMarket: TargetMarket{
			PrimarySegments: []MarketSegment{
				{ID: "s1", Name: "Enterprise", Description: "Fortune 500 companies"},
			},
			BuyerPersonas: []BuyerPersona{
				{ID: "bp1", Name: "CTO", Title: "Chief Technology Officer", BuyingRole: "Decision Maker"},
			},
		},
		CompetitiveLandscape: CompetitiveLandscape{
			Overview: "Fragmented market with emerging players.",
			Competitors: []Competitor{
				{ID: "c1", Name: "Competitor A", Strengths: []string{"Brand"}, Weaknesses: []string{"Price"}},
			},
		},
		MarketRequirements: []MarketRequirement{
			{ID: "mr1", Title: "Real-time Analytics", Description: "Need real-time data", Priority: PriorityMust},
		},
		Positioning: Positioning{
			Statement:       "For enterprise teams who need insights",
			TargetAudience:  "Enterprise IT teams",
			Category:        "Analytics",
			KeyBenefits:     []string{"Fast insights", "Easy to use"},
			Differentiators: []string{"AI-powered"},
		},
		SuccessMetrics: []SuccessMetric{
			{ID: "sm1", Name: "Revenue", Metric: "ARR", Target: "$10M", Timeframe: "Year 1"},
		},
		Glossary: []GlossaryTerm{
			{Term: "TAM", Definition: "Total Addressable Market"},
		},
	}

	t.Run("with frontmatter", func(t *testing.T) {
		opts := MarkdownOptions{
			IncludeFrontmatter: true,
			Margin:             "2cm",
			MainFont:           "Helvetica",
		}
		md := doc.ToMarkdown(opts)

		// Check frontmatter
		if !strings.HasPrefix(md, "---\n") {
			t.Error("expected markdown to start with YAML frontmatter")
		}
		if !strings.Contains(md, `title: "Test Market Requirements"`) {
			t.Error("expected frontmatter to contain title")
		}

		// Check content sections
		if !strings.Contains(md, "# Test Market Requirements") {
			t.Error("expected markdown to contain document title as H1")
		}
		if !strings.Contains(md, "Executive Summary") {
			t.Error("expected markdown to contain Executive Summary section")
		}
		if !strings.Contains(md, "Market Overview") {
			t.Error("expected markdown to contain Market Overview section")
		}
		if !strings.Contains(md, "$50B") {
			t.Error("expected markdown to contain TAM value")
		}
		if !strings.Contains(md, "Target Market") {
			t.Error("expected markdown to contain Target Market section")
		}
		if !strings.Contains(md, "Competitive Landscape") {
			t.Error("expected markdown to contain Competitive Landscape section")
		}
		if !strings.Contains(md, "Market Requirements") {
			t.Error("expected markdown to contain Market Requirements section")
		}
		if !strings.Contains(md, "Positioning") {
			t.Error("expected markdown to contain Positioning section")
		}
		if !strings.Contains(md, "Success Metrics") {
			t.Error("expected markdown to contain Success Metrics section")
		}
		if !strings.Contains(md, "Glossary") {
			t.Error("expected markdown to contain Glossary section")
		}
	})

	t.Run("without frontmatter", func(t *testing.T) {
		opts := MarkdownOptions{
			IncludeFrontmatter: false,
		}
		md := doc.ToMarkdown(opts)

		if strings.HasPrefix(md, "---\n") {
			t.Error("expected markdown to NOT start with YAML frontmatter")
		}
		if !strings.HasPrefix(md, "# Test Market Requirements") {
			t.Error("expected markdown to start with document title")
		}
	})
}

// TestValidation tests document validation logic.
func TestValidation(t *testing.T) {
	tests := []struct {
		name       string
		doc        Document
		wantErrors []string
	}{
		{
			name: "valid document",
			doc: Document{
				Metadata: Metadata{
					ID:      "mrd-001",
					Title:   "Valid MRD",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{
					MarketOpportunity: "Opportunity",
					ProposedOffering:  "Offering",
				},
				MarketOverview: MarketOverview{
					TAM: MarketSize{Value: "$1B"},
				},
				TargetMarket: TargetMarket{
					PrimarySegments: []MarketSegment{{ID: "s1", Name: "Seg"}},
				},
				CompetitiveLandscape: CompetitiveLandscape{
					Competitors: []Competitor{{ID: "c1", Name: "C", Strengths: []string{"S"}, Weaknesses: []string{"W"}}},
				},
				MarketRequirements: []MarketRequirement{{ID: "mr1", Title: "R", Description: "D"}},
				Positioning: Positioning{
					Statement: "Statement",
				},
			},
			wantErrors: nil,
		},
		{
			name: "missing metadata.id",
			doc: Document{
				Metadata: Metadata{
					Title:   "Test",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary:     ExecutiveSummary{MarketOpportunity: "M", ProposedOffering: "O"},
				MarketOverview:       MarketOverview{TAM: MarketSize{Value: "$1B"}},
				TargetMarket:         TargetMarket{PrimarySegments: []MarketSegment{{ID: "s1"}}},
				CompetitiveLandscape: CompetitiveLandscape{Competitors: []Competitor{{ID: "c1", Strengths: []string{"S"}, Weaknesses: []string{"W"}}}},
				MarketRequirements:   []MarketRequirement{{ID: "mr1"}},
				Positioning:          Positioning{Statement: "S"},
			},
			wantErrors: []string{"metadata.id is required"},
		},
		{
			name: "missing TAM",
			doc: Document{
				Metadata: Metadata{
					ID:      "mrd-001",
					Title:   "Test",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary:     ExecutiveSummary{MarketOpportunity: "M", ProposedOffering: "O"},
				MarketOverview:       MarketOverview{},
				TargetMarket:         TargetMarket{PrimarySegments: []MarketSegment{{ID: "s1"}}},
				CompetitiveLandscape: CompetitiveLandscape{Competitors: []Competitor{{ID: "c1", Strengths: []string{"S"}, Weaknesses: []string{"W"}}}},
				MarketRequirements:   []MarketRequirement{{ID: "mr1"}},
				Positioning:          Positioning{Statement: "S"},
			},
			wantErrors: []string{"market_overview.tam.value is required"},
		},
		{
			name: "multiple errors",
			doc: Document{
				Metadata:       Metadata{Version: "1.0.0"},
				MarketOverview: MarketOverview{},
			},
			wantErrors: []string{
				"metadata.id is required",
				"metadata.title is required",
				"metadata.authors is required",
				"executive_summary.market_opportunity is required",
				"executive_summary.proposed_offering is required",
				"market_overview.tam.value is required",
				"target_market.primary_segments is required",
				"competitive_landscape.competitors is required",
				"market_requirements is required",
				"positioning.statement is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validateDocument(tt.doc)

			if len(tt.wantErrors) == 0 {
				if len(errors) > 0 {
					t.Errorf("expected no errors, got %v", errors)
				}
				return
			}

			if len(errors) != len(tt.wantErrors) {
				t.Errorf("expected %d errors, got %d: %v", len(tt.wantErrors), len(errors), errors)
				return
			}

			for i, wantErr := range tt.wantErrors {
				if errors[i] != wantErr {
					t.Errorf("error[%d] = %q, want %q", i, errors[i], wantErr)
				}
			}
		})
	}
}

// validateDocument mirrors CLI validation logic.
func validateDocument(doc Document) []string {
	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required")
	}
	if doc.ExecutiveSummary.MarketOpportunity == "" {
		errors = append(errors, "executive_summary.market_opportunity is required")
	}
	if doc.ExecutiveSummary.ProposedOffering == "" {
		errors = append(errors, "executive_summary.proposed_offering is required")
	}
	if doc.MarketOverview.TAM.Value == "" {
		errors = append(errors, "market_overview.tam.value is required")
	}
	if len(doc.TargetMarket.PrimarySegments) == 0 {
		errors = append(errors, "target_market.primary_segments is required")
	}
	if len(doc.CompetitiveLandscape.Competitors) == 0 {
		errors = append(errors, "competitive_landscape.competitors is required")
	}
	if len(doc.MarketRequirements) == 0 {
		errors = append(errors, "market_requirements is required")
	}
	if doc.Positioning.Statement == "" {
		errors = append(errors, "positioning.statement is required")
	}

	return errors
}
