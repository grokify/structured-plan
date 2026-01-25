package prd

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PMView represents the Product Manager view of a PRD.
type PMView struct {
	Title          string           `json:"title"`
	Status         string           `json:"status"`
	Owner          string           `json:"owner"`
	Version        string           `json:"version"`
	ProblemSummary string           `json:"problem_summary"`
	Personas       []PersonaSummary `json:"personas"`
	Goals          []string         `json:"goals"`
	NonGoals       []string         `json:"non_goals"`
	Solution       SolutionSummary  `json:"solution"`
	Requirements   RequirementsList `json:"requirements"`
	Metrics        MetricsSummary   `json:"metrics"`
	Risks          []RiskSummary    `json:"risks"`
	OpenQuestions  []string         `json:"open_questions,omitempty"`
}

// PersonaSummary is a condensed persona view.
type PersonaSummary struct {
	Name       string   `json:"name"`
	Role       string   `json:"role"`
	PainPoints []string `json:"pain_points"`
	IsPrimary  bool     `json:"is_primary"`
}

// SolutionSummary is a condensed solution view.
type SolutionSummary struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Rationale   string   `json:"rationale"`
	Tradeoffs   []string `json:"tradeoffs"`
}

// RequirementsList groups requirements by priority.
type RequirementsList struct {
	Must   []string `json:"must"`
	Should []string `json:"should"`
	Could  []string `json:"could"`
}

// MetricsSummary shows key metrics.
type MetricsSummary struct {
	Primary    string   `json:"primary,omitempty"`
	Supporting []string `json:"supporting,omitempty"`
	Guardrails []string `json:"guardrails,omitempty"`
}

// RiskSummary is a condensed risk view.
type RiskSummary struct {
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Mitigation  string `json:"mitigation"`
}

// GeneratePMView creates a PM-friendly view of the Document.
func GeneratePMView(doc *Document) *PMView {
	view := &PMView{
		Title:   doc.Metadata.Title,
		Status:  string(doc.Metadata.Status),
		Version: doc.Metadata.Version,
	}

	// Get owner from authors
	if len(doc.Metadata.Authors) > 0 {
		view.Owner = doc.Metadata.Authors[0].Name
	}

	// Problem summary
	if doc.Problem != nil && doc.Problem.Statement != "" {
		view.ProblemSummary = doc.Problem.Statement
		if doc.Problem.UserImpact != "" {
			view.ProblemSummary += "\n\nImpact: " + doc.Problem.UserImpact
		}
	} else {
		view.ProblemSummary = doc.ExecutiveSummary.ProblemStatement
	}

	// Personas
	for _, persona := range doc.Personas {
		view.Personas = append(view.Personas, PersonaSummary{
			Name:       persona.Name,
			Role:       persona.Role,
			PainPoints: persona.PainPoints,
			IsPrimary:  persona.IsPrimary,
		})
	}

	// Goals from objectives
	for _, obj := range doc.Objectives.BusinessObjectives {
		view.Goals = append(view.Goals, obj.Description)
	}
	for _, obj := range doc.Objectives.ProductGoals {
		view.Goals = append(view.Goals, obj.Description)
	}

	// Non-goals from out of scope
	view.NonGoals = doc.OutOfScope

	// Solution
	if doc.Solution != nil && doc.Solution.SelectedSolutionID != "" {
		if selected := doc.Solution.SelectedSolution(); selected != nil {
			view.Solution = SolutionSummary{
				Name:        selected.Name,
				Description: selected.Description,
				Rationale:   doc.Solution.SolutionRationale,
				Tradeoffs:   selected.Tradeoffs,
			}
		}
	} else if doc.ExecutiveSummary.ProposedSolution != "" {
		view.Solution = SolutionSummary{
			Name:        "Proposed Solution",
			Description: doc.ExecutiveSummary.ProposedSolution,
		}
	}

	// Requirements by priority
	for _, req := range doc.Requirements.Functional {
		switch req.Priority {
		case MoSCoWMust:
			view.Requirements.Must = append(view.Requirements.Must, req.Description)
		case MoSCoWShould:
			view.Requirements.Should = append(view.Requirements.Should, req.Description)
		case MoSCoWCould:
			view.Requirements.Could = append(view.Requirements.Could, req.Description)
		}
	}

	// Metrics
	if len(doc.Objectives.SuccessMetrics) > 0 {
		// First metric as primary
		view.Metrics.Primary = fmt.Sprintf("%s: %s", doc.Objectives.SuccessMetrics[0].Name, doc.Objectives.SuccessMetrics[0].Target)
		// Rest as supporting
		for i := 1; i < len(doc.Objectives.SuccessMetrics); i++ {
			m := doc.Objectives.SuccessMetrics[i]
			view.Metrics.Supporting = append(view.Metrics.Supporting, m.Name)
		}
	}

	// Risks
	for _, r := range doc.Risks {
		view.Risks = append(view.Risks, RiskSummary{
			Description: r.Description,
			Impact:      string(r.Impact),
			Mitigation:  r.Mitigation,
		})
	}

	return view
}

// ExecView represents the Executive view of a PRD.
type ExecView struct {
	Header                ExecHeader   `json:"header"`
	Strengths             []string     `json:"strengths"`
	Blockers              []string     `json:"blockers"`
	RequiredActions       []ExecAction `json:"required_actions"`
	TopRisks              []ExecRisk   `json:"top_risks"`
	RecommendationSummary string       `json:"recommendation_summary"`
}

// ExecHeader contains high-level decision info.
type ExecHeader struct {
	PRDID           string  `json:"prd_id"`
	Title           string  `json:"title"`
	OverallDecision string  `json:"overall_decision"`
	ConfidenceLevel string  `json:"confidence_level"`
	OverallScore    float64 `json:"overall_score"`
}

// ExecAction represents a required action.
type ExecAction struct {
	Action    string `json:"action"`
	Owner     string `json:"owner"`
	Priority  string `json:"priority"`
	Rationale string `json:"rationale"`
}

// ExecRisk is a high-level risk summary.
type ExecRisk struct {
	Risk       string `json:"risk"`
	Impact     string `json:"impact"`
	Mitigation string `json:"mitigation"`
	Confidence string `json:"confidence"`
}

// GenerateExecView creates an executive-friendly view of the Document.
func GenerateExecView(doc *Document, scores *ScoringResult) *ExecView {
	view := &ExecView{
		Header: ExecHeader{
			PRDID: doc.Metadata.ID,
			Title: doc.Metadata.Title,
		},
	}

	// Decision and confidence from scores
	if scores != nil {
		view.Header.OverallScore = scores.WeightedScore
		view.Header.OverallDecision = decisionToExec(scores.Decision)
		view.Header.ConfidenceLevel = scoreToConfidence(scores)
		view.Blockers = scores.Blockers
	} else {
		view.Header.OverallDecision = "Pending Review"
		view.Header.ConfidenceLevel = "Unknown"
	}

	// Strengths (top 3 scoring categories)
	view.Strengths = extractStrengths(doc, scores)

	// Required actions from revision triggers
	if scores != nil {
		for _, trigger := range scores.RevisionTriggers {
			if trigger.Severity == "blocker" || trigger.Severity == "major" {
				view.RequiredActions = append(view.RequiredActions, ExecAction{
					Action:    trigger.Description,
					Owner:     trigger.RecommendedOwner,
					Priority:  trigger.Severity,
					Rationale: fmt.Sprintf("Scoring identified issue in %s", trigger.Category),
				})
			}
		}
	}

	// Top risks (max 3)
	count := 0
	for _, r := range doc.Risks {
		if count >= 3 {
			break
		}
		if r.Impact == RiskImpactHigh || r.Impact == RiskImpactMedium || r.Impact == RiskImpactCritical {
			view.TopRisks = append(view.TopRisks, ExecRisk{
				Risk:       r.Description,
				Impact:     string(r.Impact),
				Mitigation: r.Mitigation,
				Confidence: "Medium",
			})
			count++
		}
	}

	// Generate recommendation summary
	view.RecommendationSummary = generateExecSummary(doc, scores)

	return view
}

func decisionToExec(decision string) string {
	switch decision {
	case "approve":
		return "Proceed"
	case "revise":
		return "Proceed with Revisions"
	case "reject":
		return "Do Not Proceed"
	case "human_review":
		return "Requires Leadership Review"
	default:
		return "Pending"
	}
}

func scoreToConfidence(scores *ScoringResult) string {
	if len(scores.Blockers) > 0 {
		return "Low"
	}
	if scores.WeightedScore >= 8.0 {
		return "High"
	}
	if scores.WeightedScore >= 6.5 {
		return "Medium"
	}
	return "Low"
}

func extractStrengths(doc *Document, scores *ScoringResult) []string {
	var strengths []string

	// If we have scores, use top categories
	if scores != nil {
		type catScore struct {
			name  string
			score float64
		}
		var topCats []catScore
		for _, cs := range scores.CategoryScores {
			if cs.Score >= 7.0 {
				topCats = append(topCats, catScore{cs.Category, cs.Score})
			}
		}
		// Sort by score descending (simple bubble sort for small list)
		for i := 0; i < len(topCats)-1; i++ {
			for j := i + 1; j < len(topCats); j++ {
				if topCats[j].score > topCats[i].score {
					topCats[i], topCats[j] = topCats[j], topCats[i]
				}
			}
		}
		// Take top 3
		for i := 0; i < 3 && i < len(topCats); i++ {
			strengths = append(strengths, categoryToStrength(topCats[i].name, doc))
		}
	}

	// If we don't have 3 strengths, add generic ones
	if len(strengths) < 3 {
		if doc.ExecutiveSummary.ProblemStatement != "" {
			strengths = append(strengths, "Clear problem definition")
		}
		if doc.Solution != nil && doc.Solution.SelectedSolutionID != "" {
			strengths = append(strengths, "Solution selected with documented rationale")
		}
		if len(doc.OutOfScope) > 0 {
			strengths = append(strengths, "Well-defined scope with explicit non-goals")
		}
		if len(doc.Personas) >= 2 {
			strengths = append(strengths, "Multiple user personas defined")
		}
	}

	if len(strengths) > 3 {
		strengths = strengths[:3]
	}

	return strengths
}

func categoryToStrength(category string, doc *Document) string {
	switch category {
	case "problem_definition":
		return "Clear articulation of the core user problem"
	case "user_understanding":
		return "Strong user understanding with validated personas"
	case "market_awareness":
		return "Good competitive awareness and differentiation strategy"
	case "solution_fit":
		return "Solution is well-aligned with identified problems"
	case "scope_discipline":
		return "Clear scope with well-defined goals and non-goals"
	case "requirements_quality":
		return "Requirements are clear, prioritized, and testable"
	case "ux_coverage":
		return "User journeys well-documented with edge cases"
	case "technical_feasibility":
		return "Technical approach is realistic and well-assessed"
	case "metrics_quality":
		return "Strong success metrics with clear targets"
	case "risk_management":
		return "Proactive risk identification with mitigations"
	default:
		return fmt.Sprintf("Strong %s", strings.ReplaceAll(category, "_", " "))
	}
}

func generateExecSummary(doc *Document, scores *ScoringResult) string {
	var parts []string

	// Opening - what is this about
	parts = append(parts, fmt.Sprintf("This initiative (%s) ", doc.Metadata.Title))

	// Problem validity
	if doc.ExecutiveSummary.ProblemStatement != "" {
		parts = append(parts, "addresses a defined user problem")
		if doc.Problem != nil && len(doc.Problem.Evidence) > 0 {
			parts = append(parts, " backed by evidence")
		}
	} else {
		parts = append(parts, "lacks a clear problem definition")
	}

	// Decision
	if scores != nil {
		switch scores.Decision {
		case "approve":
			parts = append(parts, ". The PRD is ready for approval and development can proceed.")
		case "revise":
			parts = append(parts, fmt.Sprintf(". The PRD requires targeted revisions (%d issues identified) before approval.", len(scores.RevisionTriggers)))
		case "reject":
			parts = append(parts, fmt.Sprintf(". The PRD has blocking issues (%d blockers) that must be resolved before proceeding.", len(scores.Blockers)))
		case "human_review":
			parts = append(parts, ". The PRD requires leadership review due to significant gaps.")
		}
	}

	return strings.Join(parts, "")
}

// RenderPMMarkdown generates markdown output for PM view.
func RenderPMMarkdown(view *PMView) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", view.Title))
	sb.WriteString(fmt.Sprintf("**Status:** %s | **Owner:** %s | **Version:** %s\n\n", view.Status, view.Owner, view.Version))

	sb.WriteString("## Problem\n\n")
	sb.WriteString(view.ProblemSummary + "\n\n")

	if len(view.Personas) > 0 {
		sb.WriteString("## Target Users\n\n")
		for _, p := range view.Personas {
			primary := ""
			if p.IsPrimary {
				primary = " (Primary)"
			}
			sb.WriteString(fmt.Sprintf("### %s%s\n", p.Name, primary))
			sb.WriteString(fmt.Sprintf("**Role:** %s\n\n", p.Role))
			if len(p.PainPoints) > 0 {
				sb.WriteString("**Pain Points:**\n\n")
				for _, pp := range p.PainPoints {
					sb.WriteString(fmt.Sprintf("- %s\n", pp))
				}
				sb.WriteString("\n")
			}
		}
	}

	if len(view.Goals) > 0 {
		sb.WriteString("## Goals\n\n")
		for _, g := range view.Goals {
			sb.WriteString(fmt.Sprintf("- %s\n", g))
		}
		sb.WriteString("\n")
	}

	if len(view.NonGoals) > 0 {
		sb.WriteString("## Non-Goals\n\n")
		for _, ng := range view.NonGoals {
			sb.WriteString(fmt.Sprintf("- %s\n", ng))
		}
		sb.WriteString("\n")
	}

	if view.Solution.Name != "" {
		sb.WriteString("## Solution\n\n")
		sb.WriteString(fmt.Sprintf("### %s\n\n", view.Solution.Name))
		sb.WriteString(view.Solution.Description + "\n\n")
		if view.Solution.Rationale != "" {
			sb.WriteString(fmt.Sprintf("**Rationale:** %s\n\n", view.Solution.Rationale))
		}
	}

	if len(view.Requirements.Must) > 0 || len(view.Requirements.Should) > 0 {
		sb.WriteString("## Requirements\n\n")
		if len(view.Requirements.Must) > 0 {
			sb.WriteString("### Must Have\n\n")
			for _, r := range view.Requirements.Must {
				sb.WriteString(fmt.Sprintf("- %s\n", r))
			}
			sb.WriteString("\n")
		}
		if len(view.Requirements.Should) > 0 {
			sb.WriteString("### Should Have\n\n")
			for _, r := range view.Requirements.Should {
				sb.WriteString(fmt.Sprintf("- %s\n", r))
			}
			sb.WriteString("\n")
		}
		if len(view.Requirements.Could) > 0 {
			sb.WriteString("### Could Have\n\n")
			for _, r := range view.Requirements.Could {
				sb.WriteString(fmt.Sprintf("- %s\n", r))
			}
			sb.WriteString("\n")
		}
	}

	if view.Metrics.Primary != "" {
		sb.WriteString("## Success Metrics\n\n")
		sb.WriteString(fmt.Sprintf("**Primary:** %s\n\n", view.Metrics.Primary))
	}

	if len(view.Risks) > 0 {
		sb.WriteString("## Key Risks\n\n")
		for _, r := range view.Risks {
			sb.WriteString(fmt.Sprintf("- **%s** (%s impact): %s\n", r.Description, r.Impact, r.Mitigation))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderExecMarkdown generates markdown output for exec view.
func RenderExecMarkdown(view *ExecView) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Executive Summary: %s\n\n", view.Header.Title))
	sb.WriteString(fmt.Sprintf("**PRD ID:** %s\n\n", view.Header.PRDID))

	// Decision box
	sb.WriteString("## Decision\n\n")
	sb.WriteString("| Decision | Confidence | Score |\n")
	sb.WriteString("|----------|------------|-------|\n")
	sb.WriteString(fmt.Sprintf("| **%s** | %s | %.1f/10 |\n\n", view.Header.OverallDecision, view.Header.ConfidenceLevel, view.Header.OverallScore))

	if len(view.Strengths) > 0 {
		sb.WriteString("## What's Working\n\n")
		for _, s := range view.Strengths {
			sb.WriteString(fmt.Sprintf("- %s\n", s))
		}
		sb.WriteString("\n")
	}

	if len(view.Blockers) > 0 {
		sb.WriteString("## Blocking Issues\n\n")
		for _, b := range view.Blockers {
			sb.WriteString(fmt.Sprintf("- %s\n", b))
		}
		sb.WriteString("\n")
	}

	if len(view.RequiredActions) > 0 {
		sb.WriteString("## Required Actions\n\n")
		for i, a := range view.RequiredActions {
			sb.WriteString(fmt.Sprintf("%d. **%s** (%s)\n   - Owner: %s\n   - Rationale: %s\n\n", i+1, a.Action, a.Priority, a.Owner, a.Rationale))
		}
	}

	if len(view.TopRisks) > 0 {
		sb.WriteString("## Top Risks\n\n")
		for _, r := range view.TopRisks {
			sb.WriteString(fmt.Sprintf("- **%s** (%s impact)\n  - Mitigation: %s\n\n", r.Risk, r.Impact, r.Mitigation))
		}
	}

	sb.WriteString("## Recommendation\n\n")
	sb.WriteString(view.RecommendationSummary + "\n")

	return sb.String()
}

// ToJSON converts a view to JSON string.
func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
