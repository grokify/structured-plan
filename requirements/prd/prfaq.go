package prd

import (
	"fmt"
	"strings"
	"time"
)

// PRFAQView represents the Amazon-style PR/FAQ document format.
// This is a lighter-weight alternative to the full 6-pager, focusing
// on the press release and frequently asked questions.
type PRFAQView struct {
	// Metadata
	Title   string `json:"title"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	PRDID   string `json:"prd_id"`

	// The two main sections
	PressRelease PressReleaseSection `json:"press_release"`
	FAQ          FAQSection          `json:"faq"`
}

// GeneratePRFAQView creates an Amazon-style PR/FAQ view from a PRD.
// This reuses the press release and FAQ generation logic from the 6-pager.
func GeneratePRFAQView(doc *Document) *PRFAQView {
	view := &PRFAQView{
		Title:   doc.Metadata.Title,
		Version: doc.Metadata.Version,
		PRDID:   doc.Metadata.ID,
		Date:    time.Now().Format("January 2, 2006"),
	}

	if len(doc.Metadata.Authors) > 0 {
		view.Author = doc.Metadata.Authors[0].Name
	}

	// Reuse the generation functions from 6-pager
	view.PressRelease = generatePressRelease(doc)
	view.FAQ = generateFAQ(doc)

	return view
}

// RenderPRFAQMarkdown generates markdown output for the PR/FAQ view.
func RenderPRFAQMarkdown(view *PRFAQView) string {
	var sb strings.Builder

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", view.Title))
	sb.WriteString(fmt.Sprintf("**Version:** %s | **Author:** %s | **Date:** %s\n\n", view.Version, view.Author, view.Date))
	sb.WriteString("---\n\n")

	// Press Release Section
	sb.WriteString("## Press Release\n\n")
	sb.WriteString(fmt.Sprintf("### %s\n\n", view.PressRelease.Headline))

	if view.PressRelease.Subheadline != "" {
		sb.WriteString(fmt.Sprintf("*%s*\n\n", view.PressRelease.Subheadline))
	}

	if view.PressRelease.Summary != "" {
		sb.WriteString(view.PressRelease.Summary + "\n\n")
	}

	sb.WriteString("#### The Problem\n\n")
	sb.WriteString(view.PressRelease.ProblemSolved + "\n\n")

	sb.WriteString("#### The Solution\n\n")
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
		sb.WriteString("#### Key Benefits\n\n")
		for _, b := range view.PressRelease.Benefits {
			sb.WriteString(fmt.Sprintf("- %s\n", b))
		}
		sb.WriteString("\n")
	}

	if view.PressRelease.CallToAction != "" {
		sb.WriteString(fmt.Sprintf("**%s**\n\n", view.PressRelease.CallToAction))
	}

	sb.WriteString("---\n\n")

	// FAQ Section
	sb.WriteString("## Frequently Asked Questions\n\n")

	if len(view.FAQ.CustomerFAQs) > 0 {
		sb.WriteString("### External FAQ\n\n")
		sb.WriteString("*Questions customers and users are likely to ask.*\n\n")
		for _, faq := range view.FAQ.CustomerFAQs {
			sb.WriteString(fmt.Sprintf("**Q: %s**\n\n", faq.Question))
			sb.WriteString(fmt.Sprintf("A: %s\n\n", faq.Answer))
		}
	}

	if len(view.FAQ.InternalFAQs) > 0 {
		sb.WriteString("### Internal FAQ\n\n")
		sb.WriteString("*Questions from stakeholders, leadership, and team members.*\n\n")
		for _, faq := range view.FAQ.InternalFAQs {
			sb.WriteString(fmt.Sprintf("**Q: %s**\n\n", faq.Question))
			sb.WriteString(fmt.Sprintf("A: %s\n\n", faq.Answer))
		}
	}

	if len(view.FAQ.TechnicalFAQs) > 0 {
		sb.WriteString("### Technical FAQ\n\n")
		sb.WriteString("*Questions about implementation and architecture.*\n\n")
		for _, faq := range view.FAQ.TechnicalFAQs {
			sb.WriteString(fmt.Sprintf("**Q: %s**\n\n", faq.Question))
			sb.WriteString(fmt.Sprintf("A: %s\n\n", faq.Answer))
		}
	}

	return sb.String()
}
