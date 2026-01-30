package prd

import (
	"strings"
	"testing"
)

func TestGenerateSixPagerView(t *testing.T) {
	// Create a comprehensive test PRD
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-001",
			Title:   "Test Product",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "John Doe", Role: "Product Manager"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Users struggle with complex workflows",
			ProposedSolution: "A streamlined workflow automation tool",
			ExpectedOutcomes: []string{"Reduced time to complete tasks", "Improved user satisfaction"},
		},
		Problem: &ProblemDefinition{
			Statement:  "Users spend too much time on manual tasks",
			UserImpact: "Reduced productivity and job satisfaction",
			Evidence: []Evidence{
				{Type: EvidenceInterview, Source: "User interviews", Summary: "80% of users reported frustration", Strength: StrengthHigh},
			},
		},
		Objectives: Objectives{
			OKRs: []OKR{
				{
					Objective: Objective{ID: "O-1", Description: "Increase user retention by 20%"},
					KeyResults: []KeyResult{
						{ID: "KR-1", Description: "Task Completion Time", Target: "< 5 minutes", Baseline: "10 minutes", MeasurementMethod: "Analytics"},
						{ID: "KR-2", Description: "User Satisfaction", Target: "NPS > 50", Baseline: "NPS 30"},
					},
				},
				{
					Objective: Objective{ID: "O-2", Description: "Reduce task completion time by 50%"},
					KeyResults: []KeyResult{
						{ID: "KR-3", Description: "Time reduction", Target: "50%"},
					},
				},
			},
		},
		Personas: []Persona{
			{
				ID:         "PER-1",
				Name:       "Busy Professional",
				Role:       "Project Manager",
				IsPrimary:  true,
				Goals:      []string{"Complete tasks quickly", "Track project progress"},
				PainPoints: []string{"Too many manual steps", "Lack of visibility"},
			},
		},
		Requirements: Requirements{
			Functional: []FunctionalRequirement{
				{ID: "FR-1", Description: "Automate repetitive tasks", Priority: MoSCoWMust},
				{ID: "FR-2", Description: "Provide real-time notifications", Priority: MoSCoWShould},
			},
		},
		Roadmap: Roadmap{
			Phases: []Phase{
				{
					ID:     "phase-1",
					Name:   "MVP",
					Status: PhaseStatusPlanned,
					Goals:  []string{"Launch core automation features"},
					Deliverables: []Deliverable{
						{ID: "D-1", Title: "Task Automation Engine"},
						{ID: "D-2", Title: "Basic Dashboard"},
					},
				},
			},
		},
		Risks: []Risk{
			{
				ID:          "R-1",
				Description: "Technical complexity may delay launch",
				Impact:      RiskImpactHigh,
				Mitigation:  "Start with MVP scope and iterate",
			},
		},
		OutOfScope: []string{"Mobile app", "Third-party integrations"},
		Market: &MarketDefinition{
			Alternatives: []Alternative{
				{ID: "ALT-1", Name: "Manual Process", Type: AlternativeWorkaround, Weaknesses: []string{"Time consuming", "Error prone"}},
			},
			Differentiation: []string{"AI-powered automation", "Intuitive interface"},
		},
	}

	view := GenerateSixPagerView(doc)

	// Test metadata
	if view.Title != "Test Product" {
		t.Errorf("Expected title 'Test Product', got '%s'", view.Title)
	}
	if view.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", view.Version)
	}
	if view.Author != "John Doe" {
		t.Errorf("Expected author 'John Doe', got '%s'", view.Author)
	}

	// Test Press Release section
	if !strings.Contains(view.PressRelease.Headline, "Test Product") {
		t.Errorf("Press release headline should contain product name")
	}
	if view.PressRelease.ProblemSolved == "" {
		t.Error("Press release should have problem solved")
	}
	if len(view.PressRelease.Benefits) == 0 {
		t.Error("Press release should have benefits")
	}

	// Test FAQ section
	if len(view.FAQ.CustomerFAQs) == 0 {
		t.Error("FAQ should have customer FAQs")
	}
	if len(view.FAQ.InternalFAQs) == 0 {
		t.Error("FAQ should have internal FAQs")
	}

	// Test Customer Problem section
	if view.CustomerProblem.Statement == "" {
		t.Error("Customer problem should have statement")
	}
	if len(view.CustomerProblem.Personas) == 0 {
		t.Error("Customer problem should have personas")
	}
	if len(view.CustomerProblem.CurrentAlternatives) == 0 {
		t.Error("Customer problem should have alternatives")
	}
	if len(view.CustomerProblem.Evidence) == 0 {
		t.Error("Customer problem should have evidence")
	}

	// Test Solution section
	if view.Solution.Overview == "" {
		t.Error("Solution should have overview")
	}
	if len(view.Solution.KeyFeatures) == 0 {
		t.Error("Solution should have key features")
	}
	if len(view.Solution.Differentiators) == 0 {
		t.Error("Solution should have differentiators")
	}
	if len(view.Solution.Scope.OutOfScope) == 0 {
		t.Error("Solution should have out of scope items")
	}

	// Test Success Metrics section
	if view.SuccessMetrics.PrimaryMetric.Name == "" {
		t.Error("Success metrics should have primary metric")
	}
	if view.SuccessMetrics.PrimaryMetric.Target != "< 5 minutes" {
		t.Errorf("Expected target '< 5 minutes', got '%s'", view.SuccessMetrics.PrimaryMetric.Target)
	}
	if len(view.SuccessMetrics.SecondaryMetrics) == 0 {
		t.Error("Success metrics should have secondary metrics")
	}
	if len(view.SuccessMetrics.BusinessGoals) == 0 {
		t.Error("Success metrics should have business goals")
	}

	// Test Timeline section
	if len(view.Timeline.Phases) == 0 {
		t.Error("Timeline should have phases")
	}
	if view.Timeline.Phases[0].Name != "MVP" {
		t.Errorf("Expected phase name 'MVP', got '%s'", view.Timeline.Phases[0].Name)
	}
	if len(view.Timeline.Risks) == 0 {
		t.Error("Timeline should have risks")
	}
}

func TestGenerateSixPagerViewMinimal(t *testing.T) {
	// Test with minimal PRD
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-002",
			Title:   "Minimal Product",
			Version: "0.1.0",
			Status:  StatusDraft,
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Basic problem",
			ProposedSolution: "Basic solution",
		},
	}

	view := GenerateSixPagerView(doc)

	if view.Title != "Minimal Product" {
		t.Errorf("Expected title 'Minimal Product', got '%s'", view.Title)
	}
	if view.PressRelease.ProblemSolved != "Basic problem" {
		t.Errorf("Expected problem 'Basic problem', got '%s'", view.PressRelease.ProblemSolved)
	}
	if view.Solution.Overview != "Basic solution" {
		t.Errorf("Expected solution 'Basic solution', got '%s'", view.Solution.Overview)
	}
}

func TestRenderSixPagerMarkdown(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-003",
			Title:   "Markdown Test",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "Jane Smith"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Test problem statement",
			ProposedSolution: "Test solution description",
		},
		Objectives: Objectives{
			OKRs: []OKR{
				{
					Objective: Objective{ID: "O-1", Description: "Test business goal"},
					KeyResults: []KeyResult{
						{ID: "KR-1", Description: "Test Metric", Target: "100%"},
					},
				},
			},
		},
		Personas: []Persona{
			{ID: "P-1", Name: "Test User", Role: "Developer", IsPrimary: true, PainPoints: []string{"Pain point 1"}},
		},
		Roadmap: Roadmap{
			Phases: []Phase{
				{ID: "P1", Name: "Phase 1", Goals: []string{"Goal 1"}, Deliverables: []Deliverable{{ID: "D1", Title: "Deliverable 1"}}},
			},
		},
		Risks: []Risk{
			{ID: "R1", Description: "Test risk", Impact: RiskImpactHigh, Mitigation: "Test mitigation"},
		},
	}

	view := GenerateSixPagerView(doc)
	markdown := RenderSixPagerMarkdown(view)

	// Test markdown structure
	expectedSections := []string{
		"# Markdown Test",
		"## 1. Press Release",
		"## 2. Frequently Asked Questions",
		"## 3. Customer Problem",
		"## 4. Solution",
		"## 5. Success Metrics",
		"## 6. Timeline & Resources",
	}

	for _, section := range expectedSections {
		if !strings.Contains(markdown, section) {
			t.Errorf("Markdown should contain section: %s", section)
		}
	}

	// Test content
	if !strings.Contains(markdown, "Jane Smith") {
		t.Error("Markdown should contain author name")
	}
	if !strings.Contains(markdown, "Test problem statement") {
		t.Error("Markdown should contain problem statement")
	}
	if !strings.Contains(markdown, "Test Metric") {
		t.Error("Markdown should contain metric name")
	}
	if !strings.Contains(markdown, "Phase 1") {
		t.Error("Markdown should contain phase name")
	}
}

func TestSummarizeSentence(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"Short text", 20, "Short text"},
		{"This is a longer text that should be truncated", 20, "This is a longer..."},
		{"NoSpacesHere", 5, "NoSpa..."},
		{"", 10, ""},
		{"  Trimmed  ", 20, "Trimmed"},
	}

	for _, tt := range tests {
		got := summarizeSentence(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("summarizeSentence(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}
