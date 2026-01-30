package prd

import (
	"fmt"
	"strings"
	"time"
)

// SixPagerView represents the Amazon-style 6-pager document format.
// This is a deterministic transformation from PRD data.
type SixPagerView struct {
	// Metadata
	Title   string `json:"title"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	PRDID   string `json:"prdId"`

	// The six sections
	PressRelease    PressReleaseSection    `json:"pressRelease"`
	FAQ             FAQSection             `json:"faq"`
	CustomerProblem CustomerProblemSection `json:"customerProblem"`
	Solution        SolutionSection        `json:"solution"`
	SuccessMetrics  SuccessMetricsSection  `json:"successMetrics"`
	Timeline        TimelineSection        `json:"timeline"`
}

// PressReleaseSection is a future-dated announcement of the product.
type PressReleaseSection struct {
	Headline      string   `json:"headline"`
	Subheadline   string   `json:"subheadline"`
	Summary       string   `json:"summary"`
	ProblemSolved string   `json:"problemSolved"`
	Solution      string   `json:"solution"`
	Quote         Quote    `json:"quote,omitempty"`
	CustomerQuote Quote    `json:"customerQuote,omitempty"`
	CallToAction  string   `json:"callToAction"`
	Benefits      []string `json:"benefits"`
}

// Quote represents a quote in the press release.
type Quote struct {
	Speaker string `json:"speaker"`
	Role    string `json:"role,omitempty"`
	Text    string `json:"text"`
}

// FAQSection contains anticipated questions and answers.
type FAQSection struct {
	CustomerFAQs  []FAQ `json:"customerFaqs"`
	InternalFAQs  []FAQ `json:"internalFaqs"`
	TechnicalFAQs []FAQ `json:"technicalFaqs,omitempty"`
}

// FAQ represents a question and answer pair.
type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// CustomerProblemSection describes the user pain points.
type CustomerProblemSection struct {
	Statement           string                `json:"statement"`
	Impact              string                `json:"impact"`
	Personas            []PersonaSnapshot     `json:"personas"`
	CurrentAlternatives []AlternativeSnapshot `json:"currentAlternatives"`
	Evidence            []EvidenceSnapshot    `json:"evidence,omitempty"`
}

// PersonaSnapshot is a brief persona summary for the 6-pager.
type PersonaSnapshot struct {
	Name       string   `json:"name"`
	Role       string   `json:"role"`
	PainPoints []string `json:"painPoints"`
}

// AlternativeSnapshot describes current alternatives/workarounds.
type AlternativeSnapshot struct {
	Name       string   `json:"name"`
	Weaknesses []string `json:"weaknesses"`
}

// EvidenceSnapshot is a brief evidence summary.
type EvidenceSnapshot struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Strength    string `json:"strength"`
}

// SolutionSection describes the proposed solution.
type SolutionSection struct {
	Overview        string            `json:"overview"`
	HowItWorks      string            `json:"howItWorks"`
	KeyFeatures     []FeatureSnapshot `json:"keyFeatures"`
	Differentiators []string          `json:"differentiators"`
	Scope           ScopeSnapshot     `json:"scope"`
}

// FeatureSnapshot is a brief feature summary.
type FeatureSnapshot struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// ScopeSnapshot defines what's in and out of scope.
type ScopeSnapshot struct {
	InScope    []string `json:"inScope"`
	OutOfScope []string `json:"outOfScope"`
}

// SuccessMetricsSection defines how success will be measured.
type SuccessMetricsSection struct {
	PrimaryMetric    MetricSnapshot   `json:"primaryMetric"`
	SecondaryMetrics []MetricSnapshot `json:"secondaryMetrics"`
	Guardrails       []MetricSnapshot `json:"guardrails,omitempty"`
	BusinessGoals    []string         `json:"businessGoals"`
}

// MetricSnapshot is a brief metric summary.
type MetricSnapshot struct {
	Name        string `json:"name"`
	Target      string `json:"target"`
	Baseline    string `json:"baseline,omitempty"`
	Measurement string `json:"measurement,omitempty"`
}

// TimelineSection describes the roadmap and resources.
type TimelineSection struct {
	Phases       []PhaseSnapshot `json:"phases"`
	Dependencies []string        `json:"dependencies"`
	Risks        []RiskSnapshot  `json:"risks"`
	TeamNeeds    string          `json:"teamNeeds,omitempty"`
}

// PhaseSnapshot is a brief phase summary.
type PhaseSnapshot struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Deliverables []string `json:"deliverables"`
	Status       string   `json:"status,omitempty"`
}

// RiskSnapshot is a brief risk summary.
type RiskSnapshot struct {
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Mitigation  string `json:"mitigation"`
}

// GenerateSixPagerView creates an Amazon-style 6-pager view from a PRD.
func GenerateSixPagerView(doc *Document) *SixPagerView {
	view := &SixPagerView{
		Title:   doc.Metadata.Title,
		Version: doc.Metadata.Version,
		PRDID:   doc.Metadata.ID,
		Date:    time.Now().Format("January 2, 2006"),
	}

	if len(doc.Metadata.Authors) > 0 {
		view.Author = doc.Metadata.Authors[0].Name
	}

	// Generate each section
	view.PressRelease = generatePressRelease(doc)
	view.FAQ = generateFAQ(doc)
	view.CustomerProblem = generateCustomerProblem(doc)
	view.Solution = generateSolution(doc)
	view.SuccessMetrics = generateSuccessMetrics(doc)
	view.Timeline = generateTimeline(doc)

	return view
}

func generatePressRelease(doc *Document) PressReleaseSection {
	pr := PressReleaseSection{}

	// Headline from title
	pr.Headline = fmt.Sprintf("Introducing %s", doc.Metadata.Title)

	// Subheadline from problem + solution
	if doc.Problem != nil && doc.Problem.Statement != "" {
		pr.Subheadline = summarizeSentence(doc.Problem.Statement, 100)
	} else if doc.ExecutiveSummary.ProblemStatement != "" {
		pr.Subheadline = summarizeSentence(doc.ExecutiveSummary.ProblemStatement, 100)
	}

	// Summary paragraph
	if doc.ExecutiveSummary.ProposedSolution != "" {
		pr.Summary = doc.ExecutiveSummary.ProposedSolution
	}

	// Problem solved
	if doc.Problem != nil && doc.Problem.Statement != "" {
		pr.ProblemSolved = doc.Problem.Statement
		if doc.Problem.UserImpact != "" {
			pr.ProblemSolved += " " + doc.Problem.UserImpact
		}
	} else {
		pr.ProblemSolved = doc.ExecutiveSummary.ProblemStatement
	}

	// Solution description
	if doc.Solution != nil {
		if selected := doc.Solution.SelectedSolution(); selected != nil {
			pr.Solution = selected.Description
		}
	}
	if pr.Solution == "" {
		pr.Solution = doc.ExecutiveSummary.ProposedSolution
	}

	// Internal quote from author
	if len(doc.Metadata.Authors) > 0 {
		author := doc.Metadata.Authors[0]
		pr.Quote = Quote{
			Speaker: author.Name,
			Role:    author.Role,
			Text:    fmt.Sprintf("We're excited to deliver %s to help our users solve their most pressing challenges.", doc.Metadata.Title),
		}
	}

	// Customer quote from primary persona
	for _, persona := range doc.Personas {
		if persona.IsPrimary && len(persona.Goals) > 0 {
			pr.CustomerQuote = Quote{
				Speaker: persona.Name,
				Role:    persona.Role,
				Text:    fmt.Sprintf("This solution helps me %s - exactly what I've been looking for.", strings.ToLower(persona.Goals[0])),
			}
			break
		}
	}

	// Benefits from OKR objectives
	for _, okr := range doc.Objectives.OKRs {
		pr.Benefits = append(pr.Benefits, okr.Objective.Description)
	}

	// Call to action
	pr.CallToAction = fmt.Sprintf("Learn more about %s and get started today.", doc.Metadata.Title)

	return pr
}

func generateFAQ(doc *Document) FAQSection {
	faq := FAQSection{}

	// Customer FAQs from personas and user stories
	if len(doc.Personas) > 0 {
		for _, persona := range doc.Personas {
			if len(persona.PainPoints) > 0 {
				faq.CustomerFAQs = append(faq.CustomerFAQs, FAQ{
					Question: fmt.Sprintf("How does this help someone in a %s role?", persona.Role),
					Answer:   fmt.Sprintf("For %s, this solution addresses: %s", persona.Role, strings.Join(persona.PainPoints, ", ")),
				})
			}
		}
	}

	// What is this?
	if doc.ExecutiveSummary.ProposedSolution != "" {
		faq.CustomerFAQs = append(faq.CustomerFAQs, FAQ{
			Question: fmt.Sprintf("What is %s?", doc.Metadata.Title),
			Answer:   doc.ExecutiveSummary.ProposedSolution,
		})
	}

	// Internal FAQs from assumptions and constraints
	if doc.Assumptions != nil {
		for _, assumption := range doc.Assumptions.Assumptions {
			faq.InternalFAQs = append(faq.InternalFAQs, FAQ{
				Question: fmt.Sprintf("Why do we assume %s?", summarizeSentence(assumption.Description, 50)),
				Answer:   assumption.Description,
			})
		}

		for _, constraint := range doc.Assumptions.Constraints {
			faq.InternalFAQs = append(faq.InternalFAQs, FAQ{
				Question: fmt.Sprintf("What about the constraint: %s?", summarizeSentence(constraint.Description, 50)),
				Answer:   constraint.Description,
			})
		}
	}

	// Why not alternatives?
	if doc.Market != nil {
		for _, alt := range doc.Market.Alternatives {
			if alt.WhyNotChosen != "" {
				faq.InternalFAQs = append(faq.InternalFAQs, FAQ{
					Question: fmt.Sprintf("Why not use %s instead?", alt.Name),
					Answer:   alt.WhyNotChosen,
				})
			}
		}
	}

	// Out of scope as FAQ
	for _, item := range doc.OutOfScope {
		faq.InternalFAQs = append(faq.InternalFAQs, FAQ{
			Question: fmt.Sprintf("Why isn't %s included?", summarizeSentence(item, 50)),
			Answer:   fmt.Sprintf("%s is out of scope for this release.", item),
		})
	}

	// Risk FAQs
	for _, risk := range doc.Risks {
		if risk.Impact == RiskImpactHigh || risk.Impact == RiskImpactCritical {
			faq.InternalFAQs = append(faq.InternalFAQs, FAQ{
				Question: fmt.Sprintf("What if %s?", summarizeSentence(risk.Description, 50)),
				Answer:   fmt.Sprintf("We've identified this risk and plan to mitigate it by: %s", risk.Mitigation),
			})
		}
	}

	// Technical FAQs from architecture
	if doc.TechArchitecture != nil {
		if doc.TechArchitecture.Overview != "" {
			faq.TechnicalFAQs = append(faq.TechnicalFAQs, FAQ{
				Question: "How does the system work at a high level?",
				Answer:   doc.TechArchitecture.Overview,
			})
		}

		for _, integration := range doc.TechArchitecture.IntegrationPoints {
			faq.TechnicalFAQs = append(faq.TechnicalFAQs, FAQ{
				Question: fmt.Sprintf("How does it integrate with %s?", integration.Name),
				Answer:   integration.Description,
			})
		}
	}

	return faq
}

func generateCustomerProblem(doc *Document) CustomerProblemSection {
	cp := CustomerProblemSection{}

	// Statement
	if doc.Problem != nil && doc.Problem.Statement != "" {
		cp.Statement = doc.Problem.Statement
		cp.Impact = doc.Problem.UserImpact
	} else {
		cp.Statement = doc.ExecutiveSummary.ProblemStatement
	}

	// Personas
	for _, persona := range doc.Personas {
		cp.Personas = append(cp.Personas, PersonaSnapshot{
			Name:       persona.Name,
			Role:       persona.Role,
			PainPoints: persona.PainPoints,
		})
	}

	// Current alternatives
	if doc.Market != nil {
		for _, alt := range doc.Market.Alternatives {
			cp.CurrentAlternatives = append(cp.CurrentAlternatives, AlternativeSnapshot{
				Name:       alt.Name,
				Weaknesses: alt.Weaknesses,
			})
		}
	}

	// Evidence
	if doc.Problem != nil {
		for _, evidence := range doc.Problem.Evidence {
			desc := evidence.Summary
			if desc == "" {
				desc = evidence.Source
			}
			cp.Evidence = append(cp.Evidence, EvidenceSnapshot{
				Type:        string(evidence.Type),
				Description: desc,
				Strength:    string(evidence.Strength),
			})
		}
	}

	return cp
}

func generateSolution(doc *Document) SolutionSection {
	sol := SolutionSection{}

	// Overview
	if doc.Solution != nil {
		if selected := doc.Solution.SelectedSolution(); selected != nil {
			sol.Overview = selected.Description
			// HowItWorks from benefits if available
			if len(selected.Benefits) > 0 {
				sol.HowItWorks = strings.Join(selected.Benefits, ". ")
			}
		}
	}
	if sol.Overview == "" {
		sol.Overview = doc.ExecutiveSummary.ProposedSolution
	}

	// Key features from Must-have requirements
	for _, req := range doc.Requirements.Functional {
		if req.Priority == MoSCoWMust {
			sol.KeyFeatures = append(sol.KeyFeatures, FeatureSnapshot{
				Name:        req.ID,
				Description: req.Description,
				Priority:    string(req.Priority),
			})
		}
	}

	// Differentiators
	if doc.Market != nil {
		sol.Differentiators = doc.Market.Differentiation
	}

	// Scope from OKR objectives
	for _, okr := range doc.Objectives.OKRs {
		sol.Scope.InScope = append(sol.Scope.InScope, okr.Objective.Description)
	}
	sol.Scope.OutOfScope = doc.OutOfScope

	return sol
}

func generateSuccessMetrics(doc *Document) SuccessMetricsSection {
	sm := SuccessMetricsSection{}

	// Metrics from OKR Key Results
	isFirst := true
	for _, okr := range doc.Objectives.OKRs {
		// Add objective description as business goal
		sm.BusinessGoals = append(sm.BusinessGoals, okr.Objective.Description)

		// Add key results as metrics
		for _, kr := range okr.KeyResults {
			metric := MetricSnapshot{
				Name:        kr.Description,
				Target:      kr.Target,
				Baseline:    kr.Baseline,
				Measurement: kr.MeasurementMethod,
			}
			if isFirst {
				sm.PrimaryMetric = metric
				isFirst = false
			} else {
				sm.SecondaryMetrics = append(sm.SecondaryMetrics, metric)
			}
		}
	}

	return sm
}

func generateTimeline(doc *Document) TimelineSection {
	tl := TimelineSection{}

	// Phases from roadmap
	for _, phase := range doc.Roadmap.Phases {
		var deliverables []string
		for _, d := range phase.Deliverables {
			deliverables = append(deliverables, d.Title)
		}

		// Use Goals joined as description if available
		description := ""
		if len(phase.Goals) > 0 {
			description = strings.Join(phase.Goals, "; ")
		}

		tl.Phases = append(tl.Phases, PhaseSnapshot{
			Name:         phase.Name,
			Description:  description,
			Deliverables: deliverables,
			Status:       string(phase.Status),
		})
	}

	// Dependencies from assumptions
	if doc.Assumptions != nil {
		for _, dep := range doc.Assumptions.Dependencies {
			tl.Dependencies = append(tl.Dependencies, dep.Description)
		}
	}

	// Technical dependencies
	if doc.TechArchitecture != nil {
		stack := doc.TechArchitecture.TechnologyStack
		for _, tech := range stack.Backend {
			tl.Dependencies = append(tl.Dependencies, fmt.Sprintf("Backend: %s", tech.Name))
		}
		for _, tech := range stack.Database {
			tl.Dependencies = append(tl.Dependencies, fmt.Sprintf("Database: %s", tech.Name))
		}
		for _, tech := range stack.Infrastructure {
			tl.Dependencies = append(tl.Dependencies, fmt.Sprintf("Infrastructure: %s", tech.Name))
		}
	}

	// Top risks
	for _, risk := range doc.Risks {
		if risk.Impact == RiskImpactHigh || risk.Impact == RiskImpactCritical || risk.Impact == RiskImpactMedium {
			tl.Risks = append(tl.Risks, RiskSnapshot{
				Description: risk.Description,
				Impact:      string(risk.Impact),
				Mitigation:  risk.Mitigation,
			})
		}
	}

	return tl
}

// Helper function to summarize/truncate a sentence
func summarizeSentence(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	// Find last space before maxLen
	idx := strings.LastIndex(s[:maxLen], " ")
	if idx == -1 {
		idx = maxLen
	}
	return s[:idx] + "..."
}

// RenderSixPagerMarkdown generates markdown output for the 6-pager view.
func RenderSixPagerMarkdown(view *SixPagerView) string {
	var sb strings.Builder

	// Title page
	sb.WriteString(fmt.Sprintf("# %s\n\n", view.Title))
	sb.WriteString(fmt.Sprintf("**Version:** %s | **Author:** %s | **Date:** %s\n\n", view.Version, view.Author, view.Date))
	sb.WriteString("---\n\n")

	// Section 1: Press Release
	sb.WriteString("## 1. Press Release\n\n")
	sb.WriteString(fmt.Sprintf("### %s\n\n", view.PressRelease.Headline))
	if view.PressRelease.Subheadline != "" {
		sb.WriteString(fmt.Sprintf("*%s*\n\n", view.PressRelease.Subheadline))
	}
	if view.PressRelease.Summary != "" {
		sb.WriteString(view.PressRelease.Summary + "\n\n")
	}

	sb.WriteString("**The Problem:**\n\n")
	sb.WriteString(view.PressRelease.ProblemSolved + "\n\n")

	sb.WriteString("**The Solution:**\n\n")
	sb.WriteString(view.PressRelease.Solution + "\n\n")

	if view.PressRelease.Quote.Text != "" {
		sb.WriteString(fmt.Sprintf("> \"%s\"\n>\n> — %s", view.PressRelease.Quote.Text, view.PressRelease.Quote.Speaker))
		if view.PressRelease.Quote.Role != "" {
			sb.WriteString(fmt.Sprintf(", %s", view.PressRelease.Quote.Role))
		}
		sb.WriteString("\n\n")
	}

	if view.PressRelease.CustomerQuote.Text != "" {
		sb.WriteString(fmt.Sprintf("> \"%s\"\n>\n> — %s", view.PressRelease.CustomerQuote.Text, view.PressRelease.CustomerQuote.Speaker))
		if view.PressRelease.CustomerQuote.Role != "" {
			sb.WriteString(fmt.Sprintf(", %s", view.PressRelease.CustomerQuote.Role))
		}
		sb.WriteString("\n\n")
	}

	if len(view.PressRelease.Benefits) > 0 {
		sb.WriteString("**Key Benefits:**\n\n")
		for _, b := range view.PressRelease.Benefits {
			sb.WriteString(fmt.Sprintf("- %s\n", b))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Section 2: FAQ
	sb.WriteString("## 2. Frequently Asked Questions\n\n")

	if len(view.FAQ.CustomerFAQs) > 0 {
		sb.WriteString("### Customer Questions\n\n")
		for _, faq := range view.FAQ.CustomerFAQs {
			sb.WriteString(fmt.Sprintf("**Q: %s**\n\n", faq.Question))
			sb.WriteString(fmt.Sprintf("A: %s\n\n", faq.Answer))
		}
	}

	if len(view.FAQ.InternalFAQs) > 0 {
		sb.WriteString("### Internal Questions\n\n")
		for _, faq := range view.FAQ.InternalFAQs {
			sb.WriteString(fmt.Sprintf("**Q: %s**\n\n", faq.Question))
			sb.WriteString(fmt.Sprintf("A: %s\n\n", faq.Answer))
		}
	}

	if len(view.FAQ.TechnicalFAQs) > 0 {
		sb.WriteString("### Technical Questions\n\n")
		for _, faq := range view.FAQ.TechnicalFAQs {
			sb.WriteString(fmt.Sprintf("**Q: %s**\n\n", faq.Question))
			sb.WriteString(fmt.Sprintf("A: %s\n\n", faq.Answer))
		}
	}

	sb.WriteString("---\n\n")

	// Section 3: Customer Problem
	sb.WriteString("## 3. Customer Problem\n\n")

	sb.WriteString("### Problem Statement\n\n")
	sb.WriteString(view.CustomerProblem.Statement + "\n\n")

	if view.CustomerProblem.Impact != "" {
		sb.WriteString("### Impact\n\n")
		sb.WriteString(view.CustomerProblem.Impact + "\n\n")
	}

	if len(view.CustomerProblem.Personas) > 0 {
		sb.WriteString("### Who Is Affected\n\n")
		for _, p := range view.CustomerProblem.Personas {
			sb.WriteString(fmt.Sprintf("**%s** (%s)\n\n", p.Name, p.Role))
			if len(p.PainPoints) > 0 {
				sb.WriteString("Pain Points:\n\n")
				for _, pp := range p.PainPoints {
					sb.WriteString(fmt.Sprintf("- %s\n", pp))
				}
				sb.WriteString("\n")
			}
		}
	}

	if len(view.CustomerProblem.CurrentAlternatives) > 0 {
		sb.WriteString("### Current Alternatives\n\n")
		for _, alt := range view.CustomerProblem.CurrentAlternatives {
			sb.WriteString(fmt.Sprintf("**%s**\n\n", alt.Name))
			if len(alt.Weaknesses) > 0 {
				sb.WriteString("Weaknesses:\n\n")
				for _, w := range alt.Weaknesses {
					sb.WriteString(fmt.Sprintf("- %s\n", w))
				}
				sb.WriteString("\n")
			}
		}
	}

	if len(view.CustomerProblem.Evidence) > 0 {
		sb.WriteString("### Evidence\n\n")
		for _, e := range view.CustomerProblem.Evidence {
			sb.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", e.Type, e.Strength, e.Description))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Section 4: Solution
	sb.WriteString("## 4. Solution\n\n")

	sb.WriteString("### Overview\n\n")
	sb.WriteString(view.Solution.Overview + "\n\n")

	if view.Solution.HowItWorks != "" {
		sb.WriteString("### How It Works\n\n")
		sb.WriteString(view.Solution.HowItWorks + "\n\n")
	}

	if len(view.Solution.KeyFeatures) > 0 {
		sb.WriteString("### Key Features\n\n")
		for _, f := range view.Solution.KeyFeatures {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", f.Name, f.Description))
		}
		sb.WriteString("\n")
	}

	if len(view.Solution.Differentiators) > 0 {
		sb.WriteString("### Differentiators\n\n")
		for _, d := range view.Solution.Differentiators {
			sb.WriteString(fmt.Sprintf("- %s\n", d))
		}
		sb.WriteString("\n")
	}

	if len(view.Solution.Scope.InScope) > 0 || len(view.Solution.Scope.OutOfScope) > 0 {
		sb.WriteString("### Scope\n\n")
		if len(view.Solution.Scope.InScope) > 0 {
			sb.WriteString("**In Scope:**\n\n")
			for _, s := range view.Solution.Scope.InScope {
				sb.WriteString(fmt.Sprintf("- %s\n", s))
			}
			sb.WriteString("\n")
		}
		if len(view.Solution.Scope.OutOfScope) > 0 {
			sb.WriteString("**Out of Scope:**\n\n")
			for _, s := range view.Solution.Scope.OutOfScope {
				sb.WriteString(fmt.Sprintf("- %s\n", s))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("---\n\n")

	// Section 5: Success Metrics
	sb.WriteString("## 5. Success Metrics\n\n")

	if view.SuccessMetrics.PrimaryMetric.Name != "" {
		sb.WriteString("### Primary Metric\n\n")
		m := view.SuccessMetrics.PrimaryMetric
		sb.WriteString(fmt.Sprintf("**%s**\n\n", m.Name))
		sb.WriteString(fmt.Sprintf("- Target: %s\n", m.Target))
		if m.Baseline != "" {
			sb.WriteString(fmt.Sprintf("- Baseline: %s\n", m.Baseline))
		}
		if m.Measurement != "" {
			sb.WriteString(fmt.Sprintf("- Measurement: %s\n", m.Measurement))
		}
		sb.WriteString("\n")
	}

	if len(view.SuccessMetrics.SecondaryMetrics) > 0 {
		sb.WriteString("### Secondary Metrics\n\n")
		for _, m := range view.SuccessMetrics.SecondaryMetrics {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", m.Name, m.Target))
		}
		sb.WriteString("\n")
	}

	if len(view.SuccessMetrics.BusinessGoals) > 0 {
		sb.WriteString("### Business Goals\n\n")
		for _, g := range view.SuccessMetrics.BusinessGoals {
			sb.WriteString(fmt.Sprintf("- %s\n", g))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Section 6: Timeline & Resources
	sb.WriteString("## 6. Timeline & Resources\n\n")

	if len(view.Timeline.Phases) > 0 {
		sb.WriteString("### Phases\n\n")
		for _, p := range view.Timeline.Phases {
			status := ""
			if p.Status != "" {
				status = fmt.Sprintf(" [%s]", p.Status)
			}
			sb.WriteString(fmt.Sprintf("#### %s%s\n\n", p.Name, status))
			if p.Description != "" {
				sb.WriteString(p.Description + "\n\n")
			}
			if len(p.Deliverables) > 0 {
				sb.WriteString("Deliverables:\n\n")
				for _, d := range p.Deliverables {
					sb.WriteString(fmt.Sprintf("- %s\n", d))
				}
				sb.WriteString("\n")
			}
		}
	}

	if len(view.Timeline.Dependencies) > 0 {
		sb.WriteString("### Dependencies\n\n")
		for _, d := range view.Timeline.Dependencies {
			sb.WriteString(fmt.Sprintf("- %s\n", d))
		}
		sb.WriteString("\n")
	}

	if len(view.Timeline.Risks) > 0 {
		sb.WriteString("### Risks\n\n")
		for _, r := range view.Timeline.Risks {
			sb.WriteString(fmt.Sprintf("- **%s** (%s impact)\n", r.Description, r.Impact))
			sb.WriteString(fmt.Sprintf("  - Mitigation: %s\n", r.Mitigation))
		}
		sb.WriteString("\n")
	}

	if view.Timeline.TeamNeeds != "" {
		sb.WriteString("### Team Needs\n\n")
		sb.WriteString(view.Timeline.TeamNeeds + "\n\n")
	}

	return sb.String()
}
