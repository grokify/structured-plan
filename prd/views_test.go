package prd

import (
	"strings"
	"testing"
)

func TestGeneratePMView(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			Title:   "Test PRD",
			Status:  StatusDraft,
			Version: "1.0.0",
			Authors: []Person{{Name: "John Doe"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Users struggle with content discovery",
			ProposedSolution: "Implement AI-powered search",
		},
		Personas: []Persona{
			{
				ID:         "P-001",
				Name:       "Power User",
				Role:       "Content Creator",
				IsPrimary:  true,
				PainPoints: []string{"Slow search", "Poor results"},
			},
		},
		Objectives: Objectives{
			BusinessObjectives: []Objective{
				{ID: "BO-001", Description: "Increase engagement"},
			},
			ProductGoals: []Objective{
				{ID: "PG-001", Description: "Improve search accuracy"},
			},
			SuccessMetrics: []SuccessMetric{
				{ID: "SM-001", Name: "Search CTR", Target: "30%"},
				{ID: "SM-002", Name: "User satisfaction", Target: "4.5/5"},
			},
		},
		Requirements: Requirements{
			Functional: []FunctionalRequirement{
				{ID: "FR-001", Description: "Full-text search", Priority: MoSCoWMust},
				{ID: "FR-002", Description: "Faceted filtering", Priority: MoSCoWShould},
				{ID: "FR-003", Description: "Search history", Priority: MoSCoWCould},
			},
		},
		OutOfScope: []string{"Mobile app", "Voice search"},
		Risks: []Risk{
			{
				ID:          "R-001",
				Description: "ML model accuracy",
				Impact:      RiskImpactHigh,
				Mitigation:  "Iterative testing",
			},
		},
	}

	view := GeneratePMView(doc)

	// Test basic fields
	if view.Title != "Test PRD" {
		t.Errorf("Title = %s, want Test PRD", view.Title)
	}

	if view.Owner != "John Doe" {
		t.Errorf("Owner = %s, want John Doe", view.Owner)
	}

	if view.Status != "draft" {
		t.Errorf("Status = %s, want draft", view.Status)
	}

	// Test personas
	if len(view.Personas) != 1 {
		t.Errorf("Personas count = %d, want 1", len(view.Personas))
	}

	if !view.Personas[0].IsPrimary {
		t.Error("First persona should be primary")
	}

	// Test goals
	if len(view.Goals) != 2 {
		t.Errorf("Goals count = %d, want 2", len(view.Goals))
	}

	// Test non-goals
	if len(view.NonGoals) != 2 {
		t.Errorf("NonGoals count = %d, want 2", len(view.NonGoals))
	}

	// Test requirements by priority
	if len(view.Requirements.Must) != 1 {
		t.Errorf("Must requirements = %d, want 1", len(view.Requirements.Must))
	}
	if len(view.Requirements.Should) != 1 {
		t.Errorf("Should requirements = %d, want 1", len(view.Requirements.Should))
	}
	if len(view.Requirements.Could) != 1 {
		t.Errorf("Could requirements = %d, want 1", len(view.Requirements.Could))
	}

	// Test metrics
	if view.Metrics.Primary == "" {
		t.Error("Primary metric should not be empty")
	}

	// Test risks
	if len(view.Risks) != 1 {
		t.Errorf("Risks count = %d, want 1", len(view.Risks))
	}
}

func TestGenerateExecView(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:    "PRD-2025-001",
			Title: "AI Search Feature",
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Users need better search",
		},
		Risks: []Risk{
			{Description: "Technical risk", Impact: RiskImpactHigh, Mitigation: "Testing"},
			{Description: "Market risk", Impact: RiskImpactMedium, Mitigation: "Research"},
		},
	}

	scores := &ScoringResult{
		WeightedScore: 7.5,
		Decision:      "revise",
		CategoryScores: []CategoryScore{
			{Category: "problem_definition", Score: 8.5},
			{Category: "user_understanding", Score: 7.0},
		},
		Blockers: []string{},
		RevisionTriggers: []RevisionTrigger{
			{IssueID: "REV-1", Category: "market_awareness", Severity: "major", Description: "Add market analysis", RecommendedOwner: "market-intel"},
		},
	}

	view := GenerateExecView(doc, scores)

	// Test header
	if view.Header.PRDID != "PRD-2025-001" {
		t.Errorf("PRDID = %s, want PRD-2025-001", view.Header.PRDID)
	}

	if view.Header.OverallScore != 7.5 {
		t.Errorf("OverallScore = %f, want 7.5", view.Header.OverallScore)
	}

	if view.Header.OverallDecision != "Proceed with Revisions" {
		t.Errorf("OverallDecision = %s, want 'Proceed with Revisions'", view.Header.OverallDecision)
	}

	// Test required actions
	if len(view.RequiredActions) != 1 {
		t.Errorf("RequiredActions count = %d, want 1", len(view.RequiredActions))
	}

	// Test top risks (should include high and medium impact)
	if len(view.TopRisks) != 2 {
		t.Errorf("TopRisks count = %d, want 2", len(view.TopRisks))
	}

	// Test recommendation summary
	if view.RecommendationSummary == "" {
		t.Error("RecommendationSummary should not be empty")
	}
}

func TestGenerateExecViewWithoutScores(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:    "PRD-2025-001",
			Title: "Test Feature",
		},
	}

	view := GenerateExecView(doc, nil)

	if view.Header.OverallDecision != "Pending Review" {
		t.Errorf("OverallDecision = %s, want 'Pending Review'", view.Header.OverallDecision)
	}

	if view.Header.ConfidenceLevel != "Unknown" {
		t.Errorf("ConfidenceLevel = %s, want 'Unknown'", view.Header.ConfidenceLevel)
	}
}

func TestDecisionToExec(t *testing.T) {
	tests := []struct {
		decision string
		expected string
	}{
		{"approve", "Proceed"},
		{"revise", "Proceed with Revisions"},
		{"reject", "Do Not Proceed"},
		{"human_review", "Requires Leadership Review"},
		{"unknown", "Pending"},
	}

	for _, tt := range tests {
		t.Run(tt.decision, func(t *testing.T) {
			result := decisionToExec(tt.decision)
			if result != tt.expected {
				t.Errorf("decisionToExec(%s) = %s, want %s", tt.decision, result, tt.expected)
			}
		})
	}
}

func TestScoreToConfidence(t *testing.T) {
	tests := []struct {
		name     string
		scores   *ScoringResult
		expected string
	}{
		{
			name:     "with blockers",
			scores:   &ScoringResult{WeightedScore: 9.0, Blockers: []string{"blocker1"}},
			expected: "Low",
		},
		{
			name:     "high score no blockers",
			scores:   &ScoringResult{WeightedScore: 8.5, Blockers: []string{}},
			expected: "High",
		},
		{
			name:     "medium score",
			scores:   &ScoringResult{WeightedScore: 7.0, Blockers: []string{}},
			expected: "Medium",
		},
		{
			name:     "low score",
			scores:   &ScoringResult{WeightedScore: 5.0, Blockers: []string{}},
			expected: "Low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scoreToConfidence(tt.scores)
			if result != tt.expected {
				t.Errorf("scoreToConfidence() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestRenderPMMarkdown(t *testing.T) {
	view := &PMView{
		Title:          "Test PRD",
		Status:         "draft",
		Owner:          "John Doe",
		Version:        "1.0.0",
		ProblemSummary: "Users need better search",
		Personas: []PersonaSummary{
			{Name: "User", Role: "End User", IsPrimary: true},
		},
		Goals:    []string{"Improve search"},
		NonGoals: []string{"Mobile app"},
		Solution: SolutionSummary{
			Name:        "AI Search",
			Description: "ML-powered search",
		},
		Requirements: RequirementsList{
			Must:   []string{"Full-text search"},
			Should: []string{"Faceted filtering"},
		},
		Metrics: MetricsSummary{
			Primary: "Search CTR: 30%",
		},
		Risks: []RiskSummary{
			{Description: "Tech risk", Impact: "high", Mitigation: "Testing"},
		},
	}

	md := RenderPMMarkdown(view)

	// Check essential sections are present
	requiredSections := []string{
		"# Test PRD",
		"## Problem",
		"## Target Users",
		"## Goals",
		"## Non-Goals",
		"## Solution",
		"## Requirements",
		"## Success Metrics",
		"## Key Risks",
	}

	for _, section := range requiredSections {
		if !strings.Contains(md, section) {
			t.Errorf("Markdown missing section: %s", section)
		}
	}
}

func TestRenderExecMarkdown(t *testing.T) {
	view := &ExecView{
		Header: ExecHeader{
			PRDID:           "PRD-001",
			Title:           "Test Feature",
			OverallDecision: "Proceed",
			ConfidenceLevel: "High",
			OverallScore:    8.5,
		},
		Strengths: []string{"Clear problem", "Good metrics"},
		Blockers:  []string{},
		RequiredActions: []ExecAction{
			{Action: "Add market analysis", Owner: "PM", Priority: "major"},
		},
		TopRisks: []ExecRisk{
			{Risk: "Tech risk", Impact: "high", Mitigation: "Testing"},
		},
		RecommendationSummary: "Ready to proceed with minor revisions.",
	}

	md := RenderExecMarkdown(view)

	// Check essential sections
	requiredSections := []string{
		"# Executive Summary",
		"## Decision",
		"## What's Working",
		"## Required Actions",
		"## Top Risks",
		"## Recommendation",
	}

	for _, section := range requiredSections {
		if !strings.Contains(md, section) {
			t.Errorf("Markdown missing section: %s", section)
		}
	}

	// Check score formatting
	if !strings.Contains(md, "8.5/10") {
		t.Error("Markdown should contain formatted score")
	}
}

func TestToJSON(t *testing.T) {
	view := &PMView{
		Title:  "Test",
		Status: "draft",
	}

	json, err := ToJSON(view)
	if err != nil {
		t.Errorf("ToJSON failed: %v", err)
	}

	if !strings.Contains(json, "\"title\": \"Test\"") {
		t.Error("JSON should contain title field")
	}
}
