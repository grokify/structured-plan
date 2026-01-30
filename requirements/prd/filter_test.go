package prd

import (
	"testing"
)

func TestFilterByTags_EmptyTags(t *testing.T) {
	doc := Document{
		Personas: []Persona{
			{ID: "p1", Name: "Alice", Tags: []string{"data-management"}},
		},
	}
	filtered := doc.FilterByTags()
	if len(filtered.Personas) != 1 {
		t.Errorf("expected 1 persona, got %d", len(filtered.Personas))
	}
}

func TestFilterByTags_SingleTag(t *testing.T) {
	doc := Document{
		Personas: []Persona{
			{ID: "p1", Name: "Alice", Tags: []string{"data-management"}},
			{ID: "p2", Name: "Bob", Tags: []string{"security"}},
			{ID: "p3", Name: "Carol", Tags: []string{"data-management", "privacy"}},
		},
		UserStories: []UserStory{
			{ID: "us1", Title: "Story 1", Tags: []string{"retention"}},
			{ID: "us2", Title: "Story 2", Tags: []string{"archiving"}},
		},
		Requirements: Requirements{
			Functional: []FunctionalRequirement{
				{ID: "fr1", Title: "FR 1", Tags: []string{"data-management"}},
				{ID: "fr2", Title: "FR 2", Tags: []string{"security"}},
			},
		},
	}

	filtered := doc.FilterByTags("data-management")

	if len(filtered.Personas) != 2 {
		t.Errorf("expected 2 personas with data-management tag, got %d", len(filtered.Personas))
	}
	if len(filtered.UserStories) != 0 {
		t.Errorf("expected 0 user stories, got %d", len(filtered.UserStories))
	}
	if len(filtered.Requirements.Functional) != 1 {
		t.Errorf("expected 1 functional requirement, got %d", len(filtered.Requirements.Functional))
	}
}

func TestFilterByTags_MultipleTags_ORLogic(t *testing.T) {
	doc := Document{
		Personas: []Persona{
			{ID: "p1", Name: "Alice", Tags: []string{"data-management"}},
			{ID: "p2", Name: "Bob", Tags: []string{"security"}},
			{ID: "p3", Name: "Carol", Tags: []string{"compliance"}},
		},
	}

	// OR logic: should return personas with data-management OR security
	filtered := doc.FilterByTags("data-management", "security")

	if len(filtered.Personas) != 2 {
		t.Errorf("expected 2 personas (data-management OR security), got %d", len(filtered.Personas))
	}
}

func TestFilterByTags_Roadmap(t *testing.T) {
	doc := Document{
		Roadmap: Roadmap{
			Phases: []Phase{
				{
					ID:   "phase1",
					Name: "MVP",
					Tags: []string{"archiving"},
					Deliverables: []Deliverable{
						{ID: "d1", Title: "Deliverable 1", Tags: []string{"retention"}},
					},
				},
				{
					ID:   "phase2",
					Name: "Phase 2",
					Tags: []string{"security"},
					Deliverables: []Deliverable{
						{ID: "d2", Title: "Deliverable 2", Tags: []string{"ha-dr"}},
					},
				},
				{
					ID:   "phase3",
					Name: "Phase 3",
					Tags: []string{},
					Deliverables: []Deliverable{
						{ID: "d3", Title: "Deliverable 3", Tags: []string{"archiving"}},
					},
				},
			},
		},
	}

	filtered := doc.FilterByTags("archiving")

	if len(filtered.Roadmap.Phases) != 2 {
		t.Errorf("expected 2 phases (one with tag, one with deliverable tag), got %d", len(filtered.Roadmap.Phases))
	}
}

func TestFilterByTags_Risks(t *testing.T) {
	doc := Document{
		Risks: []Risk{
			{ID: "r1", Description: "Risk 1", Tags: []string{"privacy"}},
			{ID: "r2", Description: "Risk 2", Tags: []string{"security"}},
			{ID: "r3", Description: "Risk 3", Tags: []string{"compliance", "privacy"}},
		},
	}

	filtered := doc.FilterByTags("privacy")

	if len(filtered.Risks) != 2 {
		t.Errorf("expected 2 risks with privacy tag, got %d", len(filtered.Risks))
	}
}

func TestHasAnyTag(t *testing.T) {
	tests := []struct {
		name       string
		entityTags []string
		filterTags []string
		want       bool
	}{
		{"no match", []string{"a", "b"}, []string{"c", "d"}, false},
		{"single match", []string{"a", "b"}, []string{"b", "c"}, true},
		{"multiple matches", []string{"a", "b", "c"}, []string{"a", "c"}, true},
		{"empty entity tags", []string{}, []string{"a"}, false},
		{"empty filter tags", []string{"a"}, []string{}, false},
		{"both empty", []string{}, []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasAnyTag(tt.entityTags, tt.filterTags)
			if got != tt.want {
				t.Errorf("hasAnyTag(%v, %v) = %v, want %v", tt.entityTags, tt.filterTags, got, tt.want)
			}
		})
	}
}
