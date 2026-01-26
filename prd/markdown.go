package prd

import (
	"fmt"
	"strings"
)

// MarkdownOptions configures markdown generation.
type MarkdownOptions struct {
	// IncludeFrontmatter adds YAML frontmatter for Pandoc
	IncludeFrontmatter bool
	// Margin sets the page margin (e.g., "2cm")
	Margin string
	// MainFont sets the main font family
	MainFont string
	// SansFont sets the sans-serif font family
	SansFont string
	// MonoFont sets the monospace font family
	MonoFont string
	// FontFamily sets the LaTeX font family (e.g., "helvet")
	FontFamily string
}

// DefaultMarkdownOptions returns sensible defaults for markdown generation.
func DefaultMarkdownOptions() MarkdownOptions {
	return MarkdownOptions{
		IncludeFrontmatter: true,
		Margin:             "2cm",
		MainFont:           "Helvetica",
		SansFont:           "Helvetica",
		MonoFont:           "Courier New",
		FontFamily:         "helvet",
	}
}

// ToMarkdown converts a PRD Document to markdown format.
func (d *Document) ToMarkdown(opts MarkdownOptions) string {
	var sb strings.Builder

	// YAML Frontmatter
	if opts.IncludeFrontmatter {
		sb.WriteString(d.generateFrontmatter(opts))
	}

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", d.Metadata.Title))

	// Metadata table
	sb.WriteString(d.generateMetadataTable())

	// Executive Summary
	sb.WriteString(d.generateExecutiveSummary())

	// Objectives
	sb.WriteString(d.generateObjectives())

	// Personas
	sb.WriteString(d.generatePersonas())

	// User Stories
	sb.WriteString(d.generateUserStories())

	// Requirements
	sb.WriteString(d.generateRequirements())

	// Roadmap
	sb.WriteString(d.generateRoadmap())

	// Optional sections
	if d.TechArchitecture != nil {
		sb.WriteString(d.generateTechArchitecture())
	}

	if d.Assumptions != nil {
		sb.WriteString(d.generateAssumptions())
	}

	if len(d.OutOfScope) > 0 {
		sb.WriteString(d.generateOutOfScope())
	}

	if len(d.Risks) > 0 {
		sb.WriteString(d.generateRisks())
	}

	if len(d.Glossary) > 0 {
		sb.WriteString(d.generateGlossary())
	}

	// Custom sections
	if len(d.CustomSections) > 0 {
		sb.WriteString(d.generateCustomSections())
	}

	// Footer
	sb.WriteString("\n---\n\n*Generated from structured PRD JSON format*\n")

	return sb.String()
}

func (d *Document) generateFrontmatter(opts MarkdownOptions) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", d.Metadata.Title))

	// Authors
	if len(d.Metadata.Authors) > 0 {
		names := make([]string, len(d.Metadata.Authors))
		for i, a := range d.Metadata.Authors {
			names[i] = a.Name
		}
		sb.WriteString(fmt.Sprintf("author: %q\n", strings.Join(names, ", ")))
	}

	// Date
	if !d.Metadata.UpdatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("date: %q\n", d.Metadata.UpdatedAt.Format("2006-01-02")))
	} else if !d.Metadata.CreatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("date: %q\n", d.Metadata.CreatedAt.Format("2006-01-02")))
	}

	sb.WriteString(fmt.Sprintf("version: %q\n", d.Metadata.Version))
	sb.WriteString(fmt.Sprintf("status: %q\n", d.Metadata.Status))

	// Pandoc/LaTeX settings
	if opts.Margin != "" {
		sb.WriteString(fmt.Sprintf("geometry: margin=%s\n", opts.Margin))
	}
	if opts.MainFont != "" {
		sb.WriteString(fmt.Sprintf("mainfont: %q\n", opts.MainFont))
	}
	if opts.SansFont != "" {
		sb.WriteString(fmt.Sprintf("sansfont: %q\n", opts.SansFont))
	}
	if opts.MonoFont != "" {
		sb.WriteString(fmt.Sprintf("monofont: %q\n", opts.MonoFont))
	}
	if opts.FontFamily != "" {
		sb.WriteString(fmt.Sprintf("fontfamily: %s\n", opts.FontFamily))
	}

	sb.WriteString("header-includes:\n")
	sb.WriteString("  - \\renewcommand{\\familydefault}{\\sfdefault}\n")
	sb.WriteString("---\n\n")

	return sb.String()
}

func (d *Document) generateMetadataTable() string {
	var sb strings.Builder
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|-------|-------|\n")
	sb.WriteString(fmt.Sprintf("| **ID** | %s |\n", d.Metadata.ID))
	sb.WriteString(fmt.Sprintf("| **Version** | %s |\n", d.Metadata.Version))
	sb.WriteString(fmt.Sprintf("| **Status** | %s |\n", d.Metadata.Status))

	if !d.Metadata.CreatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("| **Created** | %s |\n", d.Metadata.CreatedAt.Format("2006-01-02")))
	}
	if !d.Metadata.UpdatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("| **Updated** | %s |\n", d.Metadata.UpdatedAt.Format("2006-01-02")))
	}

	if len(d.Metadata.Authors) > 0 {
		names := make([]string, len(d.Metadata.Authors))
		for i, a := range d.Metadata.Authors {
			names[i] = a.Name
		}
		sb.WriteString(fmt.Sprintf("| **Author(s)** | %s |\n", strings.Join(names, ", ")))
	}

	if len(d.Metadata.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("| **Tags** | %s |\n", strings.Join(d.Metadata.Tags, ", ")))
	}

	sb.WriteString("\n---\n\n")
	return sb.String()
}

func (d *Document) generateExecutiveSummary() string {
	var sb strings.Builder
	sb.WriteString("## 1. Executive Summary\n\n")

	sb.WriteString("### 1.1 Problem Statement\n\n")
	sb.WriteString(d.ExecutiveSummary.ProblemStatement + "\n\n")

	sb.WriteString("### 1.2 Proposed Solution\n\n")
	sb.WriteString(d.ExecutiveSummary.ProposedSolution + "\n\n")

	if len(d.ExecutiveSummary.ExpectedOutcomes) > 0 {
		sb.WriteString("### 1.3 Expected Outcomes\n\n")
		for _, outcome := range d.ExecutiveSummary.ExpectedOutcomes {
			sb.WriteString(fmt.Sprintf("- %s\n", outcome))
		}
		sb.WriteString("\n")
	}

	if d.ExecutiveSummary.TargetAudience != "" {
		sb.WriteString("### 1.4 Target Audience\n\n")
		sb.WriteString(d.ExecutiveSummary.TargetAudience + "\n\n")
	}

	if d.ExecutiveSummary.ValueProposition != "" {
		sb.WriteString("### 1.5 Value Proposition\n\n")
		sb.WriteString(d.ExecutiveSummary.ValueProposition + "\n\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateObjectives() string {
	var sb strings.Builder
	sb.WriteString("## 2. Objectives and Goals\n\n")

	// Business Objectives
	if len(d.Objectives.BusinessObjectives) > 0 {
		sb.WriteString("### 2.1 Business Objectives\n\n")
		sb.WriteString("| ID | Objective | Rationale | Aligned With |\n")
		sb.WriteString("|----|-----------|-----------|---------------|\n")
		for _, obj := range d.Objectives.BusinessObjectives {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				obj.ID, obj.Description, obj.Rationale, obj.AlignedWith))
		}
		sb.WriteString("\n")
	}

	// Product Goals
	if len(d.Objectives.ProductGoals) > 0 {
		sb.WriteString("### 2.2 Product Goals\n\n")
		sb.WriteString("| ID | Goal | Rationale |\n")
		sb.WriteString("|----|------|----------|\n")
		for _, goal := range d.Objectives.ProductGoals {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				goal.ID, goal.Description, goal.Rationale))
		}
		sb.WriteString("\n")
	}

	// Success Metrics
	if len(d.Objectives.SuccessMetrics) > 0 {
		sb.WriteString("### 2.3 Success Metrics\n\n")
		sb.WriteString("| ID | Metric | Target | Measurement Method |\n")
		sb.WriteString("|----|--------|--------|-------------------|\n")
		for _, m := range d.Objectives.SuccessMetrics {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				m.ID, m.Name, m.Target, m.MeasurementMethod))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generatePersonas() string {
	var sb strings.Builder
	sb.WriteString("## 3. Personas\n\n")

	for i, p := range d.Personas {
		primary := ""
		if p.IsPrimary {
			primary = " (Primary)"
		}
		sb.WriteString(fmt.Sprintf("### 3.%d %s%s\n\n", i+1, p.Name, primary))

		sb.WriteString("| Attribute | Description |\n")
		sb.WriteString("|-----------|-------------|\n")
		sb.WriteString(fmt.Sprintf("| **Role** | %s |\n", p.Role))
		sb.WriteString(fmt.Sprintf("| **Description** | %s |\n", p.Description))
		if p.TechnicalProficiency != "" {
			sb.WriteString(fmt.Sprintf("| **Technical Proficiency** | %s |\n", p.TechnicalProficiency))
		}
		sb.WriteString("\n")

		if len(p.Goals) > 0 {
			sb.WriteString("**Goals:**\n\n")
			for _, g := range p.Goals {
				sb.WriteString(fmt.Sprintf("- %s\n", g))
			}
			sb.WriteString("\n")
		}

		if len(p.PainPoints) > 0 {
			sb.WriteString("**Pain Points:**\n\n")
			for _, pp := range p.PainPoints {
				sb.WriteString(fmt.Sprintf("- %s\n", pp))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateUserStories() string {
	var sb strings.Builder
	sb.WriteString("## 4. User Stories\n\n")

	// Group by persona
	personaStories := make(map[string][]UserStory)
	for _, us := range d.UserStories {
		personaStories[us.PersonaID] = append(personaStories[us.PersonaID], us)
	}

	sectionNum := 1
	for _, p := range d.Personas {
		stories, ok := personaStories[p.ID]
		if !ok || len(stories) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("### 4.%d %s Stories\n\n", sectionNum, p.Name))
		sb.WriteString("| ID | Story | Priority | Phase |\n")
		sb.WriteString("|----|-------|----------|-------|\n")
		for _, us := range stories {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				us.ID, us.Story(), us.Priority, us.PhaseID))
		}
		sb.WriteString("\n")
		sectionNum++
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateRequirements() string {
	var sb strings.Builder

	// Functional Requirements
	sb.WriteString("## 5. Functional Requirements\n\n")

	// Group by category
	categories := make(map[string][]FunctionalRequirement)
	for _, fr := range d.Requirements.Functional {
		categories[fr.Category] = append(categories[fr.Category], fr)
	}

	sectionNum := 1
	for cat, reqs := range categories {
		sb.WriteString(fmt.Sprintf("### 5.%d %s\n\n", sectionNum, cat))
		sb.WriteString("| ID | Title | Description | Priority | Phase |\n")
		sb.WriteString("|----|-------|-------------|----------|-------|\n")
		for _, r := range reqs {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				r.ID, r.Title, truncate(r.Description, 60), r.Priority, r.PhaseID))
		}
		sb.WriteString("\n")
		sectionNum++
	}

	// Non-Functional Requirements
	sb.WriteString("## 6. Non-Functional Requirements\n\n")

	// Group by category
	nfrCategories := make(map[NFRCategory][]NonFunctionalRequirement)
	for _, nfr := range d.Requirements.NonFunctional {
		nfrCategories[nfr.Category] = append(nfrCategories[nfr.Category], nfr)
	}

	sectionNum = 1
	categoryNames := map[NFRCategory]string{
		NFRPerformance:      "Performance",
		NFRScalability:      "Scalability",
		NFRReliability:      "Reliability",
		NFRAvailability:     "Availability",
		NFRSecurity:         "Security",
		NFRMultiTenancy:     "Multi-Tenancy",
		NFRObservability:    "Observability",
		NFRMaintainability:  "Maintainability",
		NFRUsability:        "Usability",
		NFRCompatibility:    "Compatibility",
		NFRCompliance:       "Compliance",
		NFRDisasterRecovery: "Disaster Recovery",
		NFRCostEfficiency:   "Cost Efficiency",
	}

	for cat, reqs := range nfrCategories {
		catName := categoryNames[cat]
		if catName == "" {
			catName = string(cat)
		}
		sb.WriteString(fmt.Sprintf("### 6.%d %s\n\n", sectionNum, catName))
		sb.WriteString("| ID | Title | Target | Priority | Phase |\n")
		sb.WriteString("|----|-------|--------|----------|-------|\n")
		for _, r := range reqs {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				r.ID, r.Title, r.Target, r.Priority, r.PhaseID))
		}
		sb.WriteString("\n")
		sectionNum++
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateRoadmap() string {
	var sb strings.Builder
	sb.WriteString("## 7. Roadmap\n\n")

	for _, phase := range d.Roadmap.Phases {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", phase.ID, phase.Name))

		sb.WriteString(fmt.Sprintf("**Type:** %s\n\n", phase.Type))

		if len(phase.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("**Dependencies:** %s\n\n", strings.Join(phase.Dependencies, ", ")))
		}

		if len(phase.Goals) > 0 {
			sb.WriteString("**Goals:**\n\n")
			for _, g := range phase.Goals {
				sb.WriteString(fmt.Sprintf("- %s\n", g))
			}
			sb.WriteString("\n")
		}

		if len(phase.Deliverables) > 0 {
			sb.WriteString("**Deliverables:**\n\n")
			sb.WriteString("| ID | Title | Type | Status |\n")
			sb.WriteString("|----|-------|------|--------|\n")
			for _, del := range phase.Deliverables {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					del.ID, del.Title, del.Type, del.Status))
			}
			sb.WriteString("\n")
		}

		if len(phase.SuccessCriteria) > 0 {
			sb.WriteString("**Success Criteria:**\n\n")
			for _, sc := range phase.SuccessCriteria {
				sb.WriteString(fmt.Sprintf("- %s\n", sc))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

func (d *Document) generateTechArchitecture() string {
	var sb strings.Builder
	sb.WriteString("## 8. Technical Architecture\n\n")

	if d.TechArchitecture.Overview != "" {
		sb.WriteString("### 8.1 Overview\n\n")
		sb.WriteString(d.TechArchitecture.Overview + "\n\n")
	}

	if len(d.TechArchitecture.IntegrationPoints) > 0 {
		sb.WriteString("### 8.2 Integration Points\n\n")
		sb.WriteString("| ID | Name | Type | Description | Auth Method |\n")
		sb.WriteString("|----|------|------|-------------|-------------|\n")
		for _, ip := range d.TechArchitecture.IntegrationPoints {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				ip.ID, ip.Name, ip.Type, ip.Description, ip.AuthMethod))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateAssumptions() string {
	var sb strings.Builder
	sb.WriteString("## 9. Assumptions and Constraints\n\n")

	if len(d.Assumptions.Assumptions) > 0 {
		sb.WriteString("### 9.1 Assumptions\n\n")
		sb.WriteString("| ID | Assumption | Risk if Invalid |\n")
		sb.WriteString("|----|------------|------------------|\n")
		for _, a := range d.Assumptions.Assumptions {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				a.ID, a.Description, a.Risk))
		}
		sb.WriteString("\n")
	}

	if len(d.Assumptions.Constraints) > 0 {
		sb.WriteString("### 9.2 Constraints\n\n")
		sb.WriteString("| ID | Type | Constraint | Impact | Mitigation |\n")
		sb.WriteString("|----|------|------------|--------|------------|\n")
		for _, c := range d.Assumptions.Constraints {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				c.ID, c.Type, c.Description, c.Impact, c.Mitigation))
		}
		sb.WriteString("\n")
	}

	if len(d.Assumptions.Dependencies) > 0 {
		sb.WriteString("### 9.3 Dependencies\n\n")
		sb.WriteString("| ID | Name | Type | Status |\n")
		sb.WriteString("|----|------|------|--------|\n")
		for _, dep := range d.Assumptions.Dependencies {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				dep.ID, dep.Name, dep.Type, dep.Status))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateOutOfScope() string {
	var sb strings.Builder
	sb.WriteString("## 10. Out of Scope\n\n")

	for _, item := range d.OutOfScope {
		sb.WriteString(fmt.Sprintf("- %s\n", item))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateRisks() string {
	var sb strings.Builder
	sb.WriteString("## 11. Risk Assessment\n\n")

	sb.WriteString("| ID | Risk | Probability | Impact | Mitigation | Status |\n")
	sb.WriteString("|----|------|-------------|--------|------------|--------|\n")
	for _, r := range d.Risks {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			r.ID, truncate(r.Description, 50), r.Probability, r.Impact, truncate(r.Mitigation, 40), r.Status))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateGlossary() string {
	var sb strings.Builder
	sb.WriteString("## 12. Glossary\n\n")

	sb.WriteString("| Term | Definition |\n")
	sb.WriteString("|------|------------|\n")
	for _, term := range d.Glossary {
		var name string
		if term.Acronym != "" {
			name = fmt.Sprintf("**%s** (%s)", term.Term, term.Acronym)
		} else {
			name = fmt.Sprintf("**%s**", term.Term)
		}
		sb.WriteString(fmt.Sprintf("| %s | %s |\n", name, term.Definition))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateCustomSections() string {
	var sb strings.Builder

	sectionNum := 13
	for _, cs := range d.CustomSections {
		sb.WriteString(fmt.Sprintf("## %d. %s\n\n", sectionNum, cs.Title))
		if cs.Description != "" {
			sb.WriteString(cs.Description + "\n\n")
		}
		// Content is interface{}, so we just note it exists
		sb.WriteString("*See JSON source for detailed content.*\n\n")
		sb.WriteString("---\n\n")
		sectionNum++
	}

	return sb.String()
}

// truncate shortens a string to maxLen, adding "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
