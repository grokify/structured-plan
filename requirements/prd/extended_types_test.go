package prd

import (
	"encoding/json"
	"testing"
)

func TestProblemDefinitionJSON(t *testing.T) {
	problem := ProblemDefinition{
		ID:         "PROB-001",
		Statement:  "Users struggle to find relevant content",
		UserImpact: "30% longer search times than necessary",
		Evidence: []Evidence{
			{
				Type:     EvidenceInterview,
				Summary:  "Interview findings",
				Source:   "Q4 2024 User Study",
				Strength: StrengthHigh,
			},
			{
				Type:     EvidenceAnalytics,
				Summary:  "Search abandonment data",
				Source:   "Mixpanel",
				Strength: StrengthMedium,
			},
		},
		RootCauses: []string{"Poor algorithm", "No personalization"},
		Confidence: 0.85,
	}

	// Test marshaling
	data, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Test unmarshaling
	var loaded ProblemDefinition
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if loaded.ID != problem.ID {
		t.Errorf("ID = %s, want %s", loaded.ID, problem.ID)
	}

	if len(loaded.Evidence) != 2 {
		t.Errorf("Evidence count = %d, want 2", len(loaded.Evidence))
	}

	if loaded.Confidence != 0.85 {
		t.Errorf("Confidence = %f, want 0.85", loaded.Confidence)
	}
}

func TestMarketDefinitionJSON(t *testing.T) {
	market := MarketDefinition{
		Alternatives: []Alternative{
			{
				ID:          "ALT-001",
				Name:        "Competitor A",
				Type:        AlternativeCompetitor,
				Description: "Leading competitor",
				Strengths:   []string{"Market share", "Brand recognition"},
				Weaknesses:  []string{"Expensive", "Slow support"},
			},
			{
				ID:          "ALT-002",
				Name:        "Do Nothing",
				Type:        AlternativeDoNothing,
				Description: "Status quo option",
			},
		},
		Differentiation: []string{"AI-powered", "Better UX", "Lower cost"},
		MarketRisks:     []string{"Market saturation", "Economic downturn"},
	}

	data, err := json.Marshal(market)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var loaded MarketDefinition
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(loaded.Alternatives) != 2 {
		t.Errorf("Alternatives count = %d, want 2", len(loaded.Alternatives))
	}

	if loaded.Alternatives[0].Type != AlternativeCompetitor {
		t.Errorf("First alternative type = %s, want competitor", loaded.Alternatives[0].Type)
	}

	if len(loaded.Differentiation) != 3 {
		t.Errorf("Differentiation count = %d, want 3", len(loaded.Differentiation))
	}
}

func TestSolutionDefinitionJSON(t *testing.T) {
	solution := SolutionDefinition{
		SolutionOptions: []SolutionOption{
			{
				ID:                "SOL-001",
				Name:              "ML-Powered Search",
				Description:       "Implement ML recommendations",
				ProblemsAddressed: []string{"PROB-001"},
				Tradeoffs:         []string{"Higher complexity", "More compute cost"},
				EstimatedEffort:   "large",
			},
			{
				ID:                "SOL-002",
				Name:              "Enhanced Filtering",
				Description:       "Add advanced filters",
				ProblemsAddressed: []string{"PROB-001"},
				EstimatedEffort:   "medium",
			},
		},
		SelectedSolutionID: "SOL-001",
		SolutionRationale:  "Best long-term value despite higher complexity",
		Confidence:         0.85,
	}

	data, err := json.Marshal(solution)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var loaded SolutionDefinition
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(loaded.SolutionOptions) != 2 {
		t.Errorf("SolutionOptions count = %d, want 2", len(loaded.SolutionOptions))
	}

	if loaded.SelectedSolutionID != "SOL-001" {
		t.Errorf("SelectedSolutionID = %s, want SOL-001", loaded.SelectedSolutionID)
	}
}

func TestSolutionDefinitionSelectedSolution(t *testing.T) {
	solution := SolutionDefinition{
		SolutionOptions: []SolutionOption{
			{ID: "SOL-001", Name: "Option 1"},
			{ID: "SOL-002", Name: "Option 2"},
		},
		SelectedSolutionID: "SOL-002",
	}

	selected := solution.SelectedSolution()
	if selected == nil {
		t.Fatal("SelectedSolution returned nil")
	}

	if selected.Name != "Option 2" {
		t.Errorf("SelectedSolution name = %s, want Option 2", selected.Name)
	}
}

func TestSolutionDefinitionSelectedSolutionNotFound(t *testing.T) {
	solution := SolutionDefinition{
		SolutionOptions: []SolutionOption{
			{ID: "SOL-001", Name: "Option 1"},
		},
		SelectedSolutionID: "SOL-999",
	}

	selected := solution.SelectedSolution()
	if selected != nil {
		t.Error("SelectedSolution should return nil when ID not found")
	}
}

func TestDecisionsDefinitionJSON(t *testing.T) {
	decisions := DecisionsDefinition{
		Records: []DecisionRecord{
			{
				ID:                     "DEC-001",
				Decision:               "Use ML for Search",
				Status:                 DecisionAccepted,
				Rationale:              "Best long-term solution",
				AlternativesConsidered: []string{"Rule-based", "Third-party"},
				MadeBy:                 "Product Team",
			},
		},
	}

	data, err := json.Marshal(decisions)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var loaded DecisionsDefinition
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(loaded.Records) != 1 {
		t.Errorf("Records count = %d, want 1", len(loaded.Records))
	}

	if loaded.Records[0].Status != DecisionAccepted {
		t.Errorf("Status = %s, want accepted", loaded.Records[0].Status)
	}
}

func TestReviewsDefinitionJSON(t *testing.T) {
	reviews := ReviewsDefinition{
		QualityScores: &QualityScores{
			ProblemDefinition:    8.5,
			SolutionFit:          7.5,
			UserUnderstanding:    8.0,
			MarketAwareness:      6.5,
			ScopeDiscipline:      7.0,
			RequirementsQuality:  7.5,
			MetricsQuality:       8.0,
			UXCoverage:           6.0,
			TechnicalFeasibility: 7.5,
			RiskManagement:       7.0,
		},
		Decision: ReviewRevise,
		Blockers: []Blocker{
			{
				ID:          "BLK-001",
				Category:    "market_awareness",
				Description: "Missing competitive analysis",
			},
		},
		RevisionTriggers: []RevisionTrigger{
			{
				IssueID:          "REV-001",
				Category:         "market_awareness",
				Severity:         "major",
				Description:      "Add competitor analysis",
				RecommendedOwner: "market-intel",
			},
		},
	}

	data, err := json.Marshal(reviews)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var loaded ReviewsDefinition
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if loaded.QualityScores.ProblemDefinition != 8.5 {
		t.Errorf("ProblemDefinition score = %f, want 8.5", loaded.QualityScores.ProblemDefinition)
	}

	if loaded.Decision != ReviewRevise {
		t.Errorf("Decision = %s, want revise", loaded.Decision)
	}

	if len(loaded.Blockers) != 1 {
		t.Errorf("Blockers count = %d, want 1", len(loaded.Blockers))
	}
}

func TestRevisionRecordJSON(t *testing.T) {
	record := RevisionRecord{
		Version: "1.1.0",
		Changes: []string{"Added market analysis", "Updated personas"},
		Trigger: TriggerReview,
		Author:  "Jane Doe",
	}

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var loaded RevisionRecord
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if loaded.Version != "1.1.0" {
		t.Errorf("Version = %s, want 1.1.0", loaded.Version)
	}

	if loaded.Trigger != TriggerReview {
		t.Errorf("Trigger = %s, want review", loaded.Trigger)
	}

	if len(loaded.Changes) != 2 {
		t.Errorf("Changes count = %d, want 2", len(loaded.Changes))
	}
}

func TestEvidenceTypeConstants(t *testing.T) {
	tests := []struct {
		value    EvidenceType
		expected string
	}{
		{EvidenceInterview, "interview"},
		{EvidenceSurvey, "survey"},
		{EvidenceAnalytics, "analytics"},
		{EvidenceSupportTicket, "support_ticket"},
		{EvidenceMarketResearch, "market_research"},
		{EvidenceAssumption, "assumption"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("EvidenceType %v = %s, want %s", tt.value, string(tt.value), tt.expected)
		}
	}
}

func TestEvidenceStrengthConstants(t *testing.T) {
	tests := []struct {
		value    EvidenceStrength
		expected string
	}{
		{StrengthHigh, "high"},
		{StrengthMedium, "medium"},
		{StrengthLow, "low"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("EvidenceStrength %v = %s, want %s", tt.value, string(tt.value), tt.expected)
		}
	}
}

func TestAlternativeTypeConstants(t *testing.T) {
	tests := []struct {
		value    AlternativeType
		expected string
	}{
		{AlternativeCompetitor, "competitor"},
		{AlternativeWorkaround, "workaround"},
		{AlternativeDoNothing, "do_nothing"},
		{AlternativeInternalTool, "internal_tool"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("AlternativeType %v = %s, want %s", tt.value, string(tt.value), tt.expected)
		}
	}
}

func TestDecisionStatusConstants(t *testing.T) {
	tests := []struct {
		value    DecisionStatus
		expected string
	}{
		{DecisionProposed, "proposed"},
		{DecisionAccepted, "accepted"},
		{DecisionSuperseded, "superseded"},
		{DecisionDeprecated, "deprecated"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("DecisionStatus %v = %s, want %s", tt.value, string(tt.value), tt.expected)
		}
	}
}

func TestReviewDecisionConstants(t *testing.T) {
	tests := []struct {
		value    ReviewDecision
		expected string
	}{
		{ReviewApprove, "approve"},
		{ReviewRevise, "revise"},
		{ReviewReject, "reject"},
		{ReviewHumanReview, "human_review"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("ReviewDecision %v = %s, want %s", tt.value, string(tt.value), tt.expected)
		}
	}
}

func TestRevisionTriggerTypeConstants(t *testing.T) {
	tests := []struct {
		value    RevisionTriggerType
		expected string
	}{
		{TriggerInitial, "initial"},
		{TriggerReview, "review"},
		{TriggerScore, "score"},
		{TriggerHuman, "human"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("RevisionTriggerType %v = %s, want %s", tt.value, string(tt.value), tt.expected)
		}
	}
}
