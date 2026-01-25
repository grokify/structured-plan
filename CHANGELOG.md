# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Extended PRD Types** - New data structures for comprehensive PRD management:
  - `ProblemDefinition` with evidence tracking and confidence scoring (prd/problem.go)
  - `MarketDefinition` with alternatives and differentiation analysis (prd/market.go)
  - `SolutionDefinition` with multiple options and selection rationale (prd/solution.go)
  - `DecisionsDefinition` for tracking architectural and product decisions (prd/decisions.go)
  - `ReviewsDefinition` for quality assessments and review outcomes (prd/reviews.go)
  - `RevisionRecord` for document history tracking (prd/revision.go)

- **Quality Scoring System** - 10-category scoring framework (prd/scoring.go):
  - Problem Definition (20% weight)
  - Solution Fit (15% weight)
  - User Understanding (10% weight)
  - Market Awareness (10% weight)
  - Scope Discipline (10% weight)
  - Requirements Quality (10% weight)
  - Metrics Quality (10% weight)
  - UX Coverage (5% weight)
  - Technical Feasibility (5% weight)
  - Risk Management (5% weight)

- **View Generation** - Stakeholder-specific document views (prd/views.go):
  - `GeneratePMView()` - Product Manager focused view
  - `GenerateExecView()` - Executive summary view
  - `RenderPMMarkdown()` - Markdown rendering for PM view
  - `RenderExecMarkdown()` - Markdown rendering for exec view

- **File I/O Operations** (prd/io.go):
  - `Load()` - Read PRD from JSON file
  - `Save()` - Write PRD to JSON file
  - `New()` - Create new PRD with initialized fields
  - `GenerateID()` - Generate date-based PRD ID
  - `AddRevision()` - Add revision record with version increment

- **Enhanced Validation** (prd/validation.go):
  - ID uniqueness checking across all sections
  - Cross-reference validation (traceability)
  - Required field validation with specific error messages
  - Warning generation for missing optional content

- **Comprehensive Tests**:
  - scoring_test.go - Scoring system tests
  - views_test.go - View generation tests
  - io_test.go - File I/O tests
  - extended_types_test.go - Type serialization tests

### Changed

- Updated `Document` struct in prd/document.go to include new extended sections:
  - `Problem *ProblemDefinition`
  - `Market *MarketDefinition`
  - `Solution *SolutionDefinition`
  - `Decisions *DecisionsDefinition`
  - `Reviews *ReviewsDefinition`
  - `RevisionHistory []RevisionRecord`

### Notes

This release merges capabilities from the agent-team-prd project, bringing:

- Evidence-based problem definition tracking
- Multi-option solution evaluation
- Structured decision recording (ADR-style)
- 10-category quality scoring with weighted evaluation
- Automatic revision history tracking
- Stakeholder-specific document views
