package prd

import (
	"testing"
)

func TestScore(t *testing.T) {
	tests := []struct {
		name             string
		doc              *Document
		minScore         float64
		maxScore         float64
		expectBlockers   bool
		expectedDecision string
	}{
		{
			name:             "empty document scores low",
			doc:              &Document{},
			minScore:         0,
			maxScore:         3,
			expectBlockers:   true,
			expectedDecision: "reject",
		},
		{
			name: "document with problem statement",
			doc: &Document{
				ExecutiveSummary: ExecutiveSummary{
					ProblemStatement: "Users struggle to find relevant content",
					ProposedSolution: "Implement smart search with ML recommendations",
				},
				Personas: []Persona{
					{
						ID:         "P-001",
						Name:       "Content Consumer",
						Role:       "End User",
						IsPrimary:  true,
						Goals:      []string{"Find content quickly"},
						PainPoints: []string{"Current search is slow"},
					},
				},
				Objectives: Objectives{
					BusinessObjectives: []Objective{
						{ID: "BO-001", Description: "Increase user engagement"},
					},
					ProductGoals: []Objective{
						{ID: "PG-001", Description: "Improve search accuracy"},
					},
					SuccessMetrics: []SuccessMetric{
						{ID: "SM-001", Name: "Search CTR", Target: "30%"},
					},
				},
			},
			minScore:         2,
			maxScore:         5,
			expectBlockers:   true, // partial document will have blockers in missing categories
			expectedDecision: "reject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Score(tt.doc)

			if result.WeightedScore < tt.minScore || result.WeightedScore > tt.maxScore {
				t.Errorf("WeightedScore = %.2f, want between %.2f and %.2f",
					result.WeightedScore, tt.minScore, tt.maxScore)
			}

			if tt.expectBlockers && len(result.Blockers) == 0 {
				t.Error("Expected blockers but got none")
			}

			if !tt.expectBlockers && len(result.Blockers) > 0 {
				t.Errorf("Expected no blockers but got %d: %v", len(result.Blockers), result.Blockers)
			}

			if tt.expectedDecision != "" && result.Decision != tt.expectedDecision {
				t.Errorf("Decision = %s, want %s", result.Decision, tt.expectedDecision)
			}

			// All results should have category scores
			if len(result.CategoryScores) != 10 {
				t.Errorf("Expected 10 category scores, got %d", len(result.CategoryScores))
			}

			// Summary should not be empty
			if result.Summary == "" {
				t.Error("Summary should not be empty")
			}
		})
	}
}

func TestDefaultWeights(t *testing.T) {
	weights := DefaultWeights()

	if len(weights) != 10 {
		t.Errorf("Expected 10 weights, got %d", len(weights))
	}

	var totalWeight float64
	for _, w := range weights {
		if w.Weight <= 0 || w.Weight > 1 {
			t.Errorf("Weight for %s should be between 0 and 1, got %f", w.Category, w.Weight)
		}
		totalWeight += w.Weight
	}

	// Weights should sum to 1.0 (with small tolerance for floating point)
	if totalWeight < 0.99 || totalWeight > 1.01 {
		t.Errorf("Total weights should sum to 1.0, got %f", totalWeight)
	}
}

func TestScoreProblemDefinition(t *testing.T) {
	tests := []struct {
		name     string
		doc      *Document
		minScore float64
		maxScore float64
	}{
		{
			name:     "empty document",
			doc:      &Document{},
			minScore: 0,
			maxScore: 1,
		},
		{
			name: "with problem statement only",
			doc: &Document{
				ExecutiveSummary: ExecutiveSummary{
					ProblemStatement: "Users cannot find relevant content",
				},
			},
			minScore: 2.5,
			maxScore: 4,
		},
		{
			name: "with full problem definition",
			doc: &Document{
				ExecutiveSummary: ExecutiveSummary{
					ProblemStatement: "Users cannot find relevant content",
				},
				Problem: &ProblemDefinition{
					Statement:  "Users cannot find relevant content efficiently",
					UserImpact: "Users spend 30% more time searching than necessary",
					Evidence: []Evidence{
						{
							Type:     EvidenceInterview,
							Summary:  "User interviews showed frustration",
							Strength: StrengthHigh,
						},
					},
					Confidence: 0.8,
					RootCauses: []string{"Poor search algorithm", "No personalization"},
				},
			},
			minScore: 7,
			maxScore: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scoreProblemDefinition(tt.doc)

			if score.Score < tt.minScore || score.Score > tt.maxScore {
				t.Errorf("Score = %.2f, want between %.2f and %.2f",
					score.Score, tt.minScore, tt.maxScore)
			}

			if score.Category != "problem_definition" {
				t.Errorf("Category = %s, want problem_definition", score.Category)
			}
		})
	}
}

func TestGenerateJustification(t *testing.T) {
	tests := []struct {
		category string
		score    float64
		contains string
	}{
		{"problem_definition", 9.0, "strong"},
		{"user_understanding", 7.0, "adequate"},
		{"market_awareness", 5.0, "gaps"},
		{"solution_fit", 2.5, "weak"},
		{"scope_discipline", 0.5, "missing"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			result := generateJustification(tt.category, tt.score)

			if result == "" {
				t.Error("Justification should not be empty")
			}
		})
	}
}

func TestCategoryToOwner(t *testing.T) {
	tests := []struct {
		category string
		expected string
	}{
		{"problem_definition", "problem-discovery"},
		{"user_understanding", "user-research"},
		{"market_awareness", "market-intel"},
		{"solution_fit", "solution-ideation"},
		{"requirements_quality", "requirements"},
		{"unknown_category", "prd-lead"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			result := categoryToOwner(tt.category)
			if result != tt.expected {
				t.Errorf("Owner for %s = %s, want %s", tt.category, result, tt.expected)
			}
		})
	}
}

func TestMinFloat(t *testing.T) {
	tests := []struct {
		a, b, expected float64
	}{
		{5.0, 10.0, 5.0},
		{10.0, 5.0, 5.0},
		{7.5, 7.5, 7.5},
	}

	for _, tt := range tests {
		result := minFloat(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("minFloat(%f, %f) = %f, want %f", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestHasTechnologyStack(t *testing.T) {
	tests := []struct {
		name     string
		stack    TechnologyStack
		expected bool
	}{
		{
			name:     "empty stack",
			stack:    TechnologyStack{},
			expected: false,
		},
		{
			name: "with frontend",
			stack: TechnologyStack{
				Frontend: []Technology{{Name: "React"}},
			},
			expected: true,
		},
		{
			name: "with backend",
			stack: TechnologyStack{
				Backend: []Technology{{Name: "Go"}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasTechnologyStack(tt.stack)
			if result != tt.expected {
				t.Errorf("hasTechnologyStack() = %v, want %v", result, tt.expected)
			}
		})
	}
}
