package prd

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestDocumentParsing tests JSON unmarshaling of PRD documents.
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
					"id": "prd-001",
					"title": "Test PRD",
					"version": "1.0.0",
					"status": "draft",
					"created_at": "2025-01-01T00:00:00Z",
					"updated_at": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Test Author"}]
				},
				"executive_summary": {
					"problem_statement": "Test problem",
					"proposed_solution": "Test solution",
					"expected_outcomes": ["Outcome 1"]
				},
				"objectives": {
					"business_objectives": [],
					"product_goals": [],
					"success_metrics": []
				},
				"personas": [{"id": "p1", "name": "User", "role": "End User", "description": "Test user"}],
				"user_stories": [{"id": "us1", "persona_id": "p1", "title": "Test Story", "story": "As a user..."}],
				"requirements": {"functional": [], "non_functional": []},
				"roadmap": {"phases": [{"id": "phase1", "name": "Phase 1"}]}
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if doc.Metadata.ID != "prd-001" {
					t.Errorf("expected ID 'prd-001', got %q", doc.Metadata.ID)
				}
				if doc.Metadata.Title != "Test PRD" {
					t.Errorf("expected Title 'Test PRD', got %q", doc.Metadata.Title)
				}
				if doc.Metadata.Status != StatusDraft {
					t.Errorf("expected Status 'draft', got %q", doc.Metadata.Status)
				}
				if len(doc.Personas) != 1 {
					t.Errorf("expected 1 persona, got %d", len(doc.Personas))
				}
				if len(doc.UserStories) != 1 {
					t.Errorf("expected 1 user story, got %d", len(doc.UserStories))
				}
			},
		},
		{
			name: "document with all status types",
			json: `{
				"metadata": {
					"id": "prd-002",
					"title": "Test",
					"version": "1.0.0",
					"status": "approved",
					"created_at": "2025-01-01T00:00:00Z",
					"updated_at": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author"}]
				},
				"executive_summary": {"problem_statement": "P", "proposed_solution": "S", "expected_outcomes": []},
				"objectives": {"business_objectives": [], "product_goals": [], "success_metrics": []},
				"personas": [],
				"user_stories": [],
				"requirements": {"functional": [], "non_functional": []},
				"roadmap": {"phases": []}
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if doc.Metadata.Status != StatusApproved {
					t.Errorf("expected Status 'approved', got %q", doc.Metadata.Status)
				}
			},
		},
		{
			name: "document with optional sections",
			json: `{
				"metadata": {
					"id": "prd-003",
					"title": "Test",
					"version": "1.0.0",
					"status": "draft",
					"created_at": "2025-01-01T00:00:00Z",
					"updated_at": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author", "email": "author@example.com", "role": "PM"}],
					"tags": ["tag1", "tag2"]
				},
				"executive_summary": {
					"problem_statement": "Problem",
					"proposed_solution": "Solution",
					"expected_outcomes": ["O1", "O2"],
					"target_audience": "Developers",
					"value_proposition": "Value prop"
				},
				"objectives": {"business_objectives": [], "product_goals": [], "success_metrics": []},
				"personas": [],
				"user_stories": [],
				"requirements": {"functional": [], "non_functional": []},
				"roadmap": {"phases": []},
				"out_of_scope": ["Item 1", "Item 2"],
				"glossary": [{"term": "API", "definition": "Application Programming Interface"}],
				"risks": [{"id": "r1", "description": "Risk 1", "probability": "low", "impact": "high", "mitigation": "Mitigate"}]
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if len(doc.Metadata.Tags) != 2 {
					t.Errorf("expected 2 tags, got %d", len(doc.Metadata.Tags))
				}
				if doc.Metadata.Authors[0].Email != "author@example.com" {
					t.Errorf("expected email 'author@example.com', got %q", doc.Metadata.Authors[0].Email)
				}
				if len(doc.OutOfScope) != 2 {
					t.Errorf("expected 2 out of scope items, got %d", len(doc.OutOfScope))
				}
				if len(doc.Glossary) != 1 {
					t.Errorf("expected 1 glossary term, got %d", len(doc.Glossary))
				}
				if len(doc.Risks) != 1 {
					t.Errorf("expected 1 risk, got %d", len(doc.Risks))
				}
			},
		},
		{
			name: "document with requirements",
			json: `{
				"metadata": {
					"id": "prd-004",
					"title": "Test",
					"version": "1.0.0",
					"status": "draft",
					"created_at": "2025-01-01T00:00:00Z",
					"updated_at": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author"}]
				},
				"executive_summary": {"problem_statement": "P", "proposed_solution": "S", "expected_outcomes": []},
				"objectives": {"business_objectives": [], "product_goals": [], "success_metrics": []},
				"personas": [],
				"user_stories": [],
				"requirements": {
					"functional": [
						{"id": "FR-1", "title": "Feature 1", "description": "Desc", "category": "Core", "priority": "must"}
					],
					"non_functional": [
						{"id": "NFR-1", "category": "performance", "title": "Perf Req", "description": "Fast", "metric": "latency", "target": "<100ms"}
					]
				},
				"roadmap": {"phases": []}
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if len(doc.Requirements.Functional) != 1 {
					t.Errorf("expected 1 functional requirement, got %d", len(doc.Requirements.Functional))
				}
				if len(doc.Requirements.NonFunctional) != 1 {
					t.Errorf("expected 1 non-functional requirement, got %d", len(doc.Requirements.NonFunctional))
				}
				if doc.Requirements.Functional[0].Priority != MoSCoWMust {
					t.Errorf("expected priority 'must', got %q", doc.Requirements.Functional[0].Priority)
				}
			},
		},
		{
			name:    "invalid json",
			json:    `{invalid json}`,
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

// TestDocumentMarshaling tests JSON marshaling of PRD documents.
func TestDocumentMarshaling(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:        "prd-test",
			Title:     "Test Document",
			Version:   "1.0.0",
			Status:    StatusDraft,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Authors:   []Person{{Name: "Test Author"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Test problem",
			ProposedSolution: "Test solution",
			ExpectedOutcomes: []string{"Outcome 1"},
		},
		Personas:    []Persona{{ID: "p1", Name: "User", Role: "End User", Description: "Test"}},
		UserStories: []UserStory{{ID: "us1", PersonaID: "p1", Title: "Story", Story: "As a user..."}},
		Roadmap:     Roadmap{Phases: []Phase{{ID: "ph1", Name: "Phase 1"}}},
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Unmarshal back and verify
	var doc2 Document
	if err := json.Unmarshal(data, &doc2); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if doc2.Metadata.ID != doc.Metadata.ID {
		t.Errorf("round-trip ID mismatch: got %q, want %q", doc2.Metadata.ID, doc.Metadata.ID)
	}
	if doc2.Metadata.Title != doc.Metadata.Title {
		t.Errorf("round-trip Title mismatch: got %q, want %q", doc2.Metadata.Title, doc.Metadata.Title)
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

// TestMoSCoWConstants verifies MoSCoW priority constant values.
func TestMoSCoWConstants(t *testing.T) {
	tests := []struct {
		priority MoSCoW
		want     string
	}{
		{MoSCoWMust, "must"},
		{MoSCoWShould, "should"},
		{MoSCoWCould, "could"},
		{MoSCoWWont, "wont"},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			if string(tt.priority) != tt.want {
				t.Errorf("MoSCoW = %q, want %q", tt.priority, tt.want)
			}
		})
	}
}

// TestPriorityConstants verifies Priority constant values for user stories.
func TestPriorityConstants(t *testing.T) {
	tests := []struct {
		priority Priority
		want     string
	}{
		{PriorityCritical, "critical"},
		{PriorityHigh, "high"},
		{PriorityMedium, "medium"},
		{PriorityLow, "low"},
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
			ID:        "prd-test",
			Title:     "Test Product Requirements",
			Version:   "1.0.0",
			Status:    StatusDraft,
			CreatedAt: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			Authors:   []Person{{Name: "John Doe", Role: "Product Manager"}},
			Tags:      []string{"test", "example"},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Users need a better solution.",
			ProposedSolution: "Build an awesome product.",
			ExpectedOutcomes: []string{"Increased satisfaction", "Higher retention"},
			TargetAudience:   "Enterprise customers",
			ValueProposition: "Save time and money",
		},
		Objectives: Objectives{
			BusinessObjectives: []Objective{
				{ID: "bo-1", Description: "Increase revenue by 20%"},
			},
			ProductGoals: []Objective{
				{ID: "pg-1", Description: "Launch MVP"},
			},
			SuccessMetrics: []SuccessMetric{
				{ID: "sm-1", Name: "NPS", Description: "Net Promoter Score", Metric: "NPS", Target: "> 50"},
			},
		},
		Personas: []Persona{
			{ID: "p1", Name: "Developer", Role: "Software Engineer", Description: "Builds applications"},
		},
		UserStories: []UserStory{
			{ID: "us1", PersonaID: "p1", Title: "Quick Setup", Story: "As a developer, I want quick setup", Priority: "high"},
		},
		Requirements: Requirements{
			Functional: []FunctionalRequirement{
				{ID: "FR-1", Title: "User Login", Description: "Support OAuth login", Priority: MoSCoWMust},
			},
			NonFunctional: []NonFunctionalRequirement{
				{ID: "NFR-1", Category: "performance", Title: "Response Time", Description: "Fast responses", Target: "<100ms"},
			},
		},
		Roadmap: Roadmap{
			Phases: []Phase{
				{ID: "ph1", Name: "MVP", Goals: []string{"Core features"}},
			},
		},
		Glossary: []GlossaryTerm{
			{Term: "MVP", Definition: "Minimum Viable Product"},
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
		if !strings.Contains(md, `title: "Test Product Requirements"`) {
			t.Error("expected frontmatter to contain title")
		}
		if !strings.Contains(md, `geometry: margin=2cm`) {
			t.Error("expected frontmatter to contain margin")
		}

		// Check content sections
		if !strings.Contains(md, "# Test Product Requirements") {
			t.Error("expected markdown to contain document title as H1")
		}
		if !strings.Contains(md, "## 1. Executive Summary") {
			t.Error("expected markdown to contain Executive Summary section")
		}
		if !strings.Contains(md, "Users need a better solution.") {
			t.Error("expected markdown to contain problem statement")
		}
		if !strings.Contains(md, "## 2. Objectives") {
			t.Error("expected markdown to contain Objectives section")
		}
		if !strings.Contains(md, "## 3. Personas") {
			t.Error("expected markdown to contain Personas section")
		}
		if !strings.Contains(md, "Developer") {
			t.Error("expected markdown to contain persona name")
		}
		if !strings.Contains(md, "## 4. User Stories") {
			t.Error("expected markdown to contain User Stories section")
		}
		if !strings.Contains(md, "## 5. Functional Requirements") {
			t.Error("expected markdown to contain Functional Requirements section")
		}
		if !strings.Contains(md, "## 7. Roadmap") {
			t.Error("expected markdown to contain Roadmap section")
		}
		if !strings.Contains(md, "Glossary") {
			t.Error("expected markdown to contain Glossary section")
		}
		if !strings.Contains(md, "MVP") && !strings.Contains(md, "Minimum Viable Product") {
			t.Error("expected markdown to contain glossary terms")
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
		if !strings.HasPrefix(md, "# Test Product Requirements") {
			t.Error("expected markdown to start with document title")
		}
	})
}

// TestMarkdownGenerationEmptyDocument tests markdown generation with minimal document.
func TestMarkdownGenerationEmptyDocument(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:        "prd-empty",
			Title:     "Empty Document",
			Version:   "0.1.0",
			Status:    StatusDraft,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Authors:   []Person{{Name: "Author"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Problem",
			ProposedSolution: "Solution",
		},
		Roadmap: Roadmap{
			Phases: []Phase{{ID: "ph1", Name: "Phase 1"}},
		},
	}

	opts := MarkdownOptions{IncludeFrontmatter: true}
	md := doc.ToMarkdown(opts)

	// Should not panic and should produce valid output
	if md == "" {
		t.Error("expected non-empty markdown output")
	}
	if !strings.Contains(md, "# Empty Document") {
		t.Error("expected markdown to contain document title")
	}
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
					ID:      "prd-001",
					Title:   "Valid PRD",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{
					ProblemStatement: "Problem",
					ProposedSolution: "Solution",
				},
				Personas:    []Persona{{ID: "p1", Name: "User"}},
				UserStories: []UserStory{{ID: "us1", Story: "Story"}},
				Roadmap:     Roadmap{Phases: []Phase{{ID: "ph1", Name: "Phase"}}},
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
				ExecutiveSummary: ExecutiveSummary{ProblemStatement: "P", ProposedSolution: "S"},
				Personas:         []Persona{{ID: "p1"}},
				UserStories:      []UserStory{{ID: "us1"}},
				Roadmap:          Roadmap{Phases: []Phase{{ID: "ph1"}}},
			},
			wantErrors: []string{"metadata.id is required"},
		},
		{
			name: "missing metadata.title",
			doc: Document{
				Metadata: Metadata{
					ID:      "prd-001",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{ProblemStatement: "P", ProposedSolution: "S"},
				Personas:         []Persona{{ID: "p1"}},
				UserStories:      []UserStory{{ID: "us1"}},
				Roadmap:          Roadmap{Phases: []Phase{{ID: "ph1"}}},
			},
			wantErrors: []string{"metadata.title is required"},
		},
		{
			name: "missing authors",
			doc: Document{
				Metadata: Metadata{
					ID:      "prd-001",
					Title:   "Test",
					Version: "1.0.0",
				},
				ExecutiveSummary: ExecutiveSummary{ProblemStatement: "P", ProposedSolution: "S"},
				Personas:         []Persona{{ID: "p1"}},
				UserStories:      []UserStory{{ID: "us1"}},
				Roadmap:          Roadmap{Phases: []Phase{{ID: "ph1"}}},
			},
			wantErrors: []string{"metadata.authors is required"},
		},
		{
			name: "missing problem statement",
			doc: Document{
				Metadata: Metadata{
					ID:      "prd-001",
					Title:   "Test",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{ProposedSolution: "S"},
				Personas:         []Persona{{ID: "p1"}},
				UserStories:      []UserStory{{ID: "us1"}},
				Roadmap:          Roadmap{Phases: []Phase{{ID: "ph1"}}},
			},
			wantErrors: []string{"executive_summary.problem_statement is required"},
		},
		{
			name: "missing personas",
			doc: Document{
				Metadata: Metadata{
					ID:      "prd-001",
					Title:   "Test",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{ProblemStatement: "P", ProposedSolution: "S"},
				UserStories:      []UserStory{{ID: "us1"}},
				Roadmap:          Roadmap{Phases: []Phase{{ID: "ph1"}}},
			},
			wantErrors: []string{"personas is required"},
		},
		{
			name: "multiple errors",
			doc: Document{
				Metadata: Metadata{
					Version: "1.0.0",
				},
				ExecutiveSummary: ExecutiveSummary{},
				Roadmap:          Roadmap{},
			},
			wantErrors: []string{
				"metadata.id is required",
				"metadata.title is required",
				"metadata.authors is required",
				"executive_summary.problem_statement is required",
				"executive_summary.proposed_solution is required",
				"personas is required",
				"user_stories is required",
				"roadmap.phases is required",
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

// validateDocument is a helper function that mirrors CLI validation logic.
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
	if doc.ExecutiveSummary.ProblemStatement == "" {
		errors = append(errors, "executive_summary.problem_statement is required")
	}
	if doc.ExecutiveSummary.ProposedSolution == "" {
		errors = append(errors, "executive_summary.proposed_solution is required")
	}
	if len(doc.Personas) == 0 {
		errors = append(errors, "personas is required")
	}
	if len(doc.UserStories) == 0 {
		errors = append(errors, "user_stories is required")
	}
	if len(doc.Roadmap.Phases) == 0 {
		errors = append(errors, "roadmap.phases is required")
	}

	return errors
}
