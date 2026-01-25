package prd

import (
	"strings"
	"testing"
)

func TestGeneratePRFAQView(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-001",
			Title:   "Amazing Product",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "Alice Smith", Role: "Product Lead"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Users struggle with tedious manual processes",
			ProposedSolution: "An intelligent automation platform that learns user patterns",
			ExpectedOutcomes: []string{"50% time savings", "Reduced errors"},
		},
		Objectives: Objectives{
			BusinessObjectives: []Objective{
				{ID: "BO-1", Description: "Capture 10% market share in year one"},
			},
		},
		Personas: []Persona{
			{
				ID:         "PER-1",
				Name:       "Operations Manager",
				Role:       "Ops Manager",
				IsPrimary:  true,
				Goals:      []string{"Streamline operations"},
				PainPoints: []string{"Manual data entry", "Lack of visibility"},
			},
		},
		Assumptions: &AssumptionsConstraints{
			Assumptions: []Assumption{
				{ID: "A-1", Description: "Users have basic technical skills"},
			},
			Constraints: []Constraint{
				{ID: "C-1", Description: "Must integrate with existing systems"},
			},
		},
		Risks: []Risk{
			{ID: "R-1", Description: "Adoption resistance", Impact: RiskImpactHigh, Mitigation: "Change management program"},
		},
		OutOfScope: []string{"Mobile app support", "Offline mode"},
	}

	view := GeneratePRFAQView(doc)

	// Test metadata
	if view.Title != "Amazing Product" {
		t.Errorf("Expected title 'Amazing Product', got '%s'", view.Title)
	}
	if view.Author != "Alice Smith" {
		t.Errorf("Expected author 'Alice Smith', got '%s'", view.Author)
	}
	if view.PRDID != "PRD-001" {
		t.Errorf("Expected PRDID 'PRD-001', got '%s'", view.PRDID)
	}

	// Test Press Release
	if !strings.Contains(view.PressRelease.Headline, "Amazing Product") {
		t.Error("Press release headline should contain product name")
	}
	if view.PressRelease.ProblemSolved == "" {
		t.Error("Press release should have problem solved")
	}
	if view.PressRelease.Solution == "" {
		t.Error("Press release should have solution")
	}
	if view.PressRelease.Quote.Speaker != "Alice Smith" {
		t.Errorf("Expected quote speaker 'Alice Smith', got '%s'", view.PressRelease.Quote.Speaker)
	}

	// Test FAQ
	if len(view.FAQ.CustomerFAQs) == 0 {
		t.Error("Should have customer FAQs")
	}
	if len(view.FAQ.InternalFAQs) == 0 {
		t.Error("Should have internal FAQs")
	}

	// Check for specific FAQ content
	hasPersonaFAQ := false
	for _, faq := range view.FAQ.CustomerFAQs {
		if strings.Contains(faq.Question, "Ops Manager") {
			hasPersonaFAQ = true
			break
		}
	}
	if !hasPersonaFAQ {
		t.Error("Should have FAQ about persona role")
	}

	hasOutOfScopeFAQ := false
	for _, faq := range view.FAQ.InternalFAQs {
		if strings.Contains(faq.Question, "Mobile app") || strings.Contains(faq.Answer, "Mobile app") {
			hasOutOfScopeFAQ = true
			break
		}
	}
	if !hasOutOfScopeFAQ {
		t.Error("Should have FAQ about out of scope items")
	}
}

func TestGeneratePRFAQViewMinimal(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-002",
			Title:   "Simple Product",
			Version: "0.1.0",
			Status:  StatusDraft,
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Basic problem",
			ProposedSolution: "Basic solution",
		},
	}

	view := GeneratePRFAQView(doc)

	if view.Title != "Simple Product" {
		t.Errorf("Expected title 'Simple Product', got '%s'", view.Title)
	}
	if view.PressRelease.ProblemSolved != "Basic problem" {
		t.Errorf("Expected problem 'Basic problem', got '%s'", view.PressRelease.ProblemSolved)
	}
}

func TestRenderPRFAQMarkdown(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-003",
			Title:   "Markdown PR/FAQ Test",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "Bob Jones"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Test problem for PR/FAQ",
			ProposedSolution: "Test solution for PR/FAQ",
		},
		Objectives: Objectives{
			BusinessObjectives: []Objective{
				{ID: "BO-1", Description: "Test business goal"},
			},
		},
		Personas: []Persona{
			{ID: "P-1", Name: "Test User", Role: "Tester", IsPrimary: true, PainPoints: []string{"Testing pain"}},
		},
		Assumptions: &AssumptionsConstraints{
			Assumptions: []Assumption{{ID: "A-1", Description: "Test assumption"}},
		},
	}

	view := GeneratePRFAQView(doc)
	markdown := RenderPRFAQMarkdown(view)

	// Test structure
	expectedSections := []string{
		"# Markdown PR/FAQ Test",
		"## Press Release",
		"## Frequently Asked Questions",
		"### External FAQ",
		"### Internal FAQ",
	}

	for _, section := range expectedSections {
		if !strings.Contains(markdown, section) {
			t.Errorf("Markdown should contain section: %s", section)
		}
	}

	// Test content
	if !strings.Contains(markdown, "Bob Jones") {
		t.Error("Markdown should contain author name")
	}
	if !strings.Contains(markdown, "Test problem for PR/FAQ") {
		t.Error("Markdown should contain problem statement")
	}
	if !strings.Contains(markdown, "Test solution for PR/FAQ") {
		t.Error("Markdown should contain solution")
	}
	if !strings.Contains(markdown, "Q:") && !strings.Contains(markdown, "A:") {
		t.Error("Markdown should contain Q&A format")
	}
}

func TestPRFAQViewComparedToSixPager(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-004",
			Title:   "Comparison Test",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "Test Author"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Shared problem",
			ProposedSolution: "Shared solution",
		},
		Personas: []Persona{
			{ID: "P-1", Name: "Shared Persona", Role: "User", IsPrimary: true, PainPoints: []string{"Pain"}},
		},
	}

	prfaq := GeneratePRFAQView(doc)
	sixPager := GenerateSixPagerView(doc)

	// PR/FAQ should have same press release as 6-pager
	if prfaq.PressRelease.Headline != sixPager.PressRelease.Headline {
		t.Error("PR/FAQ and 6-pager should have same press release headline")
	}
	if prfaq.PressRelease.ProblemSolved != sixPager.PressRelease.ProblemSolved {
		t.Error("PR/FAQ and 6-pager should have same problem solved")
	}

	// PR/FAQ should have same FAQ as 6-pager
	if len(prfaq.FAQ.CustomerFAQs) != len(sixPager.FAQ.CustomerFAQs) {
		t.Error("PR/FAQ and 6-pager should have same customer FAQs")
	}
	if len(prfaq.FAQ.InternalFAQs) != len(sixPager.FAQ.InternalFAQs) {
		t.Error("PR/FAQ and 6-pager should have same internal FAQs")
	}
}
