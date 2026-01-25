// Command srequirements converts structured requirements documents (PRD, MRD, TRD) to markdown.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-requirements/mrd"
	"github.com/grokify/structured-requirements/prd"
	"github.com/grokify/structured-requirements/schema"
	"github.com/grokify/structured-requirements/trd"
)

var version = "0.4.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "srequirements",
	Short: "Structured requirements document generator and validator",
	Long: `srequirements is a CLI tool for working with structured requirements documents.

Supported document types:
  - PRD (Product Requirements Document)
  - MRD (Market Requirements Document)
  - TRD (Technical Requirements Document)

It can convert requirements JSON files to markdown with Pandoc-compatible YAML
frontmatter, and validate files against their respective schemas.

Example usage:
  srequirements prd generate myproduct.prd.json
  srequirements mrd validate market-analysis.mrd.json
  srequirements trd generate architecture.trd.json -o output.md`,
	Version: version,
}

// Shared generate flags
type generateFlags struct {
	output        string
	margin        string
	mainFont      string
	sansFont      string
	monoFont      string
	fontFamily    string
	noFrontmatter bool
}

func init() {
	// Add document type commands
	rootCmd.AddCommand(prdCmd)
	rootCmd.AddCommand(mrdCmd)
	rootCmd.AddCommand(trdCmd)
	rootCmd.AddCommand(schemaCmd)
}

// ============================================================================
// PRD Commands
// ============================================================================

var prdCmd = &cobra.Command{
	Use:   "prd",
	Short: "Work with Product Requirements Documents",
	Long:  `Commands for generating and validating Product Requirements Documents (PRD).`,
}

var prdGenerateFlags generateFlags

var prdGenerateCmd = &cobra.Command{
	Use:   "generate <input.json>",
	Short: "Convert PRD JSON to markdown",
	Long: `Generate markdown from a Product Requirements Document (PRD).

The output includes YAML frontmatter compatible with Pandoc for PDF generation.
By default, the output file has the same name as the input with a .md extension.`,
	Example: `  srequirements prd generate myproduct.prd.json
  srequirements prd generate myproduct.json -o output.md
  srequirements prd generate myproduct.json --no-frontmatter`,
	Args: cobra.ExactArgs(1),
	RunE: runPRDGenerate,
}

var prdValidateCmd = &cobra.Command{
	Use:     "validate <input.json>",
	Short:   "Validate PRD structure",
	Long:    `Validate a Product Requirements Document by parsing it and checking required fields.`,
	Example: `  srequirements prd validate myproduct.prd.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    runPRDValidate,
}

var prdCheckFlags struct {
	json bool
}

var prdCheckCmd = &cobra.Command{
	Use:   "check <input.json>",
	Short: "Check PRD completeness",
	Long: `Analyze a Product Requirements Document for completeness and quality.

This command evaluates each section of the PRD and provides:
  - Overall completeness score (0-100%) and letter grade
  - Per-section scores for required and optional sections
  - Specific recommendations for improvement

The check examines:
  - Required sections: metadata, executive summary, objectives, personas,
    user stories, requirements, and roadmap
  - Optional sections: assumptions, out of scope, technical architecture,
    UX requirements, risks, and glossary
  - Quality indicators: depth of content, cross-references between sections,
    acceptance criteria coverage, and NFR category coverage`,
	Example: `  srequirements prd check myproduct.prd.json
  srequirements prd check myproduct.prd.json --json`,
	Args: cobra.ExactArgs(1),
	RunE: runPRDCheck,
}

func init() {
	// PRD generate flags
	prdGenerateCmd.Flags().StringVarP(&prdGenerateFlags.output, "output", "o", "", "Output markdown file path (default: input with .md extension)")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.margin, "margin", "2cm", "Page margin for Pandoc")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.mainFont, "mainfont", "Helvetica", "Main font family")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.sansFont, "sansfont", "Helvetica", "Sans-serif font family")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.monoFont, "monofont", "Courier New", "Monospace font family")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.fontFamily, "fontfamily", "helvet", "LaTeX font family")
	prdGenerateCmd.Flags().BoolVar(&prdGenerateFlags.noFrontmatter, "no-frontmatter", false, "Disable YAML frontmatter generation")

	prdCmd.AddCommand(prdGenerateCmd)
	prdCmd.AddCommand(prdValidateCmd)
	prdCmd.AddCommand(prdCheckCmd)

	// PRD check flags
	prdCheckCmd.Flags().BoolVar(&prdCheckFlags.json, "json", false, "Output report as JSON")
}

func runPRDGenerate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Determine output file
	output := prdGenerateFlags.output
	if output == "" {
		output = deriveOutputPath(inputFile)
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	opts := prd.MarkdownOptions{
		IncludeFrontmatter: !prdGenerateFlags.noFrontmatter,
		Margin:             prdGenerateFlags.margin,
		MainFont:           prdGenerateFlags.mainFont,
		SansFont:           prdGenerateFlags.sansFont,
		MonoFont:           prdGenerateFlags.monoFont,
		FontFamily:         prdGenerateFlags.fontFamily,
	}
	markdown := doc.ToMarkdown(opts)

	if err := os.WriteFile(output, []byte(markdown), 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}

func runPRDValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required (at least one author)")
	}
	if doc.ExecutiveSummary.ProblemStatement == "" {
		errors = append(errors, "executive_summary.problem_statement is required")
	}
	if doc.ExecutiveSummary.ProposedSolution == "" {
		errors = append(errors, "executive_summary.proposed_solution is required")
	}
	if len(doc.Personas) == 0 {
		errors = append(errors, "personas is required (at least one persona)")
	}
	if len(doc.UserStories) == 0 {
		errors = append(errors, "user_stories is required (at least one user story)")
	}
	if len(doc.Roadmap.Phases) == 0 {
		errors = append(errors, "roadmap.phases is required (at least one phase)")
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(errors))
	}

	fmt.Printf("Valid PRD: %s\n", inputFile)
	fmt.Printf("  Title: %s\n", doc.Metadata.Title)
	fmt.Printf("  Version: %s\n", doc.Metadata.Version)
	fmt.Printf("  Personas: %d\n", len(doc.Personas))
	fmt.Printf("  User Stories: %d\n", len(doc.UserStories))
	fmt.Printf("  Functional Requirements: %d\n", len(doc.Requirements.Functional))
	fmt.Printf("  Non-Functional Requirements: %d\n", len(doc.Requirements.NonFunctional))
	fmt.Printf("  Phases: %d\n", len(doc.Roadmap.Phases))

	return nil
}

func runPRDCheck(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	report := doc.CheckCompleteness()

	if prdCheckFlags.json {
		output, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling report: %w", err)
		}
		fmt.Println(string(output))
	} else {
		fmt.Print(report.FormatReport())
	}

	// Return non-zero exit code if PRD has critical issues
	if report.Grade == "F" {
		return fmt.Errorf("PRD completeness check failed (Grade: F)")
	}

	return nil
}

// ============================================================================
// MRD Commands
// ============================================================================

var mrdCmd = &cobra.Command{
	Use:   "mrd",
	Short: "Work with Market Requirements Documents",
	Long:  `Commands for generating and validating Market Requirements Documents (MRD).`,
}

var mrdGenerateFlags generateFlags

var mrdGenerateCmd = &cobra.Command{
	Use:   "generate <input.json>",
	Short: "Convert MRD JSON to markdown",
	Long: `Generate markdown from a Market Requirements Document (MRD).

The output includes YAML frontmatter compatible with Pandoc for PDF generation.
By default, the output file has the same name as the input with a .md extension.`,
	Example: `  srequirements mrd generate market-analysis.mrd.json
  srequirements mrd generate market.json -o output.md
  srequirements mrd generate market.json --no-frontmatter`,
	Args: cobra.ExactArgs(1),
	RunE: runMRDGenerate,
}

var mrdValidateCmd = &cobra.Command{
	Use:     "validate <input.json>",
	Short:   "Validate MRD structure",
	Long:    `Validate a Market Requirements Document by parsing it and checking required fields.`,
	Example: `  srequirements mrd validate market-analysis.mrd.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    runMRDValidate,
}

func init() {
	// MRD generate flags
	mrdGenerateCmd.Flags().StringVarP(&mrdGenerateFlags.output, "output", "o", "", "Output markdown file path (default: input with .md extension)")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.margin, "margin", "2cm", "Page margin for Pandoc")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.mainFont, "mainfont", "Helvetica", "Main font family")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.sansFont, "sansfont", "Helvetica", "Sans-serif font family")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.monoFont, "monofont", "Courier New", "Monospace font family")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.fontFamily, "fontfamily", "helvet", "LaTeX font family")
	mrdGenerateCmd.Flags().BoolVar(&mrdGenerateFlags.noFrontmatter, "no-frontmatter", false, "Disable YAML frontmatter generation")

	mrdCmd.AddCommand(mrdGenerateCmd)
	mrdCmd.AddCommand(mrdValidateCmd)
}

func runMRDGenerate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Determine output file
	output := mrdGenerateFlags.output
	if output == "" {
		output = deriveOutputPath(inputFile)
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc mrd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	opts := mrd.MarkdownOptions{
		IncludeFrontmatter: !mrdGenerateFlags.noFrontmatter,
		Margin:             mrdGenerateFlags.margin,
		MainFont:           mrdGenerateFlags.mainFont,
		SansFont:           mrdGenerateFlags.sansFont,
		MonoFont:           mrdGenerateFlags.monoFont,
		FontFamily:         mrdGenerateFlags.fontFamily,
	}
	markdown := doc.ToMarkdown(opts)

	if err := os.WriteFile(output, []byte(markdown), 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}

func runMRDValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc mrd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required (at least one author)")
	}
	if doc.ExecutiveSummary.MarketOpportunity == "" {
		errors = append(errors, "executive_summary.market_opportunity is required")
	}
	if doc.ExecutiveSummary.ProposedOffering == "" {
		errors = append(errors, "executive_summary.proposed_offering is required")
	}
	if doc.MarketOverview.TAM.Value == "" {
		errors = append(errors, "market_overview.tam.value is required")
	}
	if len(doc.TargetMarket.PrimarySegments) == 0 {
		errors = append(errors, "target_market.primary_segments is required (at least one segment)")
	}
	if len(doc.CompetitiveLandscape.Competitors) == 0 {
		errors = append(errors, "competitive_landscape.competitors is required (at least one competitor)")
	}
	if len(doc.MarketRequirements) == 0 {
		errors = append(errors, "market_requirements is required (at least one requirement)")
	}
	if doc.Positioning.Statement == "" {
		errors = append(errors, "positioning.statement is required")
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(errors))
	}

	fmt.Printf("Valid MRD: %s\n", inputFile)
	fmt.Printf("  Title: %s\n", doc.Metadata.Title)
	fmt.Printf("  Version: %s\n", doc.Metadata.Version)
	fmt.Printf("  TAM: %s\n", doc.MarketOverview.TAM.Value)
	fmt.Printf("  Primary Segments: %d\n", len(doc.TargetMarket.PrimarySegments))
	fmt.Printf("  Buyer Personas: %d\n", len(doc.TargetMarket.BuyerPersonas))
	fmt.Printf("  Competitors: %d\n", len(doc.CompetitiveLandscape.Competitors))
	fmt.Printf("  Market Requirements: %d\n", len(doc.MarketRequirements))
	fmt.Printf("  Success Metrics: %d\n", len(doc.SuccessMetrics))

	return nil
}

// ============================================================================
// TRD Commands
// ============================================================================

var trdCmd = &cobra.Command{
	Use:   "trd",
	Short: "Work with Technical Requirements Documents",
	Long:  `Commands for generating and validating Technical Requirements Documents (TRD).`,
}

var trdGenerateFlags generateFlags

var trdGenerateCmd = &cobra.Command{
	Use:   "generate <input.json>",
	Short: "Convert TRD JSON to markdown",
	Long: `Generate markdown from a Technical Requirements Document (TRD).

The output includes YAML frontmatter compatible with Pandoc for PDF generation.
By default, the output file has the same name as the input with a .md extension.`,
	Example: `  srequirements trd generate architecture.trd.json
  srequirements trd generate tech-spec.json -o output.md
  srequirements trd generate tech-spec.json --no-frontmatter`,
	Args: cobra.ExactArgs(1),
	RunE: runTRDGenerate,
}

var trdValidateCmd = &cobra.Command{
	Use:     "validate <input.json>",
	Short:   "Validate TRD structure",
	Long:    `Validate a Technical Requirements Document by parsing it and checking required fields.`,
	Example: `  srequirements trd validate architecture.trd.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    runTRDValidate,
}

func init() {
	// TRD generate flags
	trdGenerateCmd.Flags().StringVarP(&trdGenerateFlags.output, "output", "o", "", "Output markdown file path (default: input with .md extension)")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.margin, "margin", "2cm", "Page margin for Pandoc")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.mainFont, "mainfont", "Helvetica", "Main font family")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.sansFont, "sansfont", "Helvetica", "Sans-serif font family")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.monoFont, "monofont", "Courier New", "Monospace font family")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.fontFamily, "fontfamily", "helvet", "LaTeX font family")
	trdGenerateCmd.Flags().BoolVar(&trdGenerateFlags.noFrontmatter, "no-frontmatter", false, "Disable YAML frontmatter generation")

	trdCmd.AddCommand(trdGenerateCmd)
	trdCmd.AddCommand(trdValidateCmd)
}

func runTRDGenerate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Determine output file
	output := trdGenerateFlags.output
	if output == "" {
		output = deriveOutputPath(inputFile)
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc trd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	opts := trd.MarkdownOptions{
		IncludeFrontmatter: !trdGenerateFlags.noFrontmatter,
		Margin:             trdGenerateFlags.margin,
		MainFont:           trdGenerateFlags.mainFont,
		SansFont:           trdGenerateFlags.sansFont,
		MonoFont:           trdGenerateFlags.monoFont,
		FontFamily:         trdGenerateFlags.fontFamily,
	}
	markdown := doc.ToMarkdown(opts)

	if err := os.WriteFile(output, []byte(markdown), 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}

func runTRDValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc trd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required (at least one author)")
	}
	if doc.ExecutiveSummary.Purpose == "" {
		errors = append(errors, "executive_summary.purpose is required")
	}
	if doc.ExecutiveSummary.Scope == "" {
		errors = append(errors, "executive_summary.scope is required")
	}
	if doc.Architecture.Overview == "" {
		errors = append(errors, "architecture.overview is required")
	}
	if len(doc.Architecture.Components) == 0 {
		errors = append(errors, "architecture.components is required (at least one component)")
	}
	if doc.SecurityDesign.Overview == "" {
		errors = append(errors, "security_design.overview is required")
	}
	if len(doc.Performance.Requirements) == 0 {
		errors = append(errors, "performance.requirements is required (at least one requirement)")
	}
	if len(doc.Deployment.Environments) == 0 {
		errors = append(errors, "deployment.environments is required (at least one environment)")
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(errors))
	}

	fmt.Printf("Valid TRD: %s\n", inputFile)
	fmt.Printf("  Title: %s\n", doc.Metadata.Title)
	fmt.Printf("  Version: %s\n", doc.Metadata.Version)
	fmt.Printf("  Components: %d\n", len(doc.Architecture.Components))
	fmt.Printf("  APIs: %d\n", len(doc.APISpecifications))
	fmt.Printf("  Performance Requirements: %d\n", len(doc.Performance.Requirements))
	fmt.Printf("  Environments: %d\n", len(doc.Deployment.Environments))
	fmt.Printf("  Integrations: %d\n", len(doc.Integration))

	return nil
}

// ============================================================================
// Utility Functions
// ============================================================================

func deriveOutputPath(inputFile string) string {
	ext := filepath.Ext(inputFile)
	if ext == ".json" {
		base := strings.TrimSuffix(inputFile, ext)
		return base + ".md"
	}
	return inputFile + ".md"
}

// ============================================================================
// Schema Commands
// ============================================================================

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate JSON Schema files from Go types",
	Long: `Commands for generating JSON Schema files from Go type definitions.

This implements the Go-first approach where Go structs are the source of truth
and JSON Schema is generated from them.`,
}

var schemaGenerateFlags struct {
	output  string
	docType string
}

var schemaGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate JSON Schema from Go types",
	Long: `Generate JSON Schema files from Go type definitions.

By default, generates all schema files (PRD, MRD, TRD) to the schema/ directory.
Use --type to generate a specific document type's schema.`,
	Example: `  srequirements schema generate
  srequirements schema generate -o ./schema/
  srequirements schema generate --type prd -o prd.schema.json`,
	RunE: runSchemaGenerate,
}

func init() {
	schemaGenerateCmd.Flags().StringVarP(&schemaGenerateFlags.output, "output", "o", ".", "Output directory or file path")
	schemaGenerateCmd.Flags().StringVarP(&schemaGenerateFlags.docType, "type", "t", "all", "Document type to generate (prd, mrd, trd, or all)")

	schemaCmd.AddCommand(schemaGenerateCmd)
}

func runSchemaGenerate(cmd *cobra.Command, args []string) error {
	gen := schema.NewGenerator()
	output := schemaGenerateFlags.output
	docType := strings.ToLower(schemaGenerateFlags.docType)

	switch docType {
	case "prd":
		// Single PRD schema
		path := output
		if isDir(output) {
			path = filepath.Join(output, "prd.schema.json")
		}
		if err := gen.WritePRDSchema(path); err != nil {
			return fmt.Errorf("generating PRD schema: %w", err)
		}
		fmt.Printf("Generated: %s\n", path)

	case "all":
		// All schemas to directory
		dir := output
		if !isDir(dir) {
			dir = filepath.Dir(output)
		}
		if err := gen.GenerateAll(dir); err != nil {
			return fmt.Errorf("generating schemas: %w", err)
		}
		fmt.Printf("Generated schemas in: %s\n", dir)

	case "mrd", "trd":
		return fmt.Errorf("schema generation for %s is not yet implemented", docType)

	default:
		return fmt.Errorf("unknown document type: %s (expected prd, mrd, trd, or all)", docType)
	}

	return nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
