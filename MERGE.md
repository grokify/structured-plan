# Merge Plan: agent-team-prd → structured-prd

This document outlines the plan to merge PRD capabilities from `agentplexus/agent-team-prd` into `grokify/structured-prd`, creating a unified PRD framework.

## Overview

**Goal**: Combine the best of both projects:

- **structured-prd**: Rich data model (60+ types), comprehensive markdown export, completeness analysis
- **agent-team-prd**: Quality scoring (10-category rubric), view generation (PM/Executive), review workflow, validation

**Result**: A single, enterprise-grade PRD library with both structural completeness checking AND quality-based scoring.

## Phase 1: Add Missing Types

### 1.1 Market Definition (from agent-team-prd)

Add to `prd/market.go`:

```go
type MarketDefinition struct {
    Alternatives    []Alternative `json:"alternatives,omitempty"`
    Differentiation []string      `json:"differentiation,omitempty"`
    MarketRisks     []string      `json:"market_risks,omitempty"`
}

type Alternative struct {
    ID           string          `json:"id"`
    Name         string          `json:"name"`
    Type         AlternativeType `json:"type"`
    Description  string          `json:"description,omitempty"`
    Strengths    []string        `json:"strengths,omitempty"`
    Weaknesses   []string        `json:"weaknesses,omitempty"`
    WhyNotChosen string          `json:"why_not_chosen,omitempty"`
}

type AlternativeType string

const (
    AlternativeCompetitor   AlternativeType = "competitor"
    AlternativeWorkaround   AlternativeType = "workaround"
    AlternativeDoNothing    AlternativeType = "do_nothing"
    AlternativeInternalTool AlternativeType = "internal_tool"
)
```

### 1.2 Solution Definition (from agent-team-prd)

Add to `prd/solution.go`:

```go
type SolutionDefinition struct {
    SolutionOptions    []SolutionOption `json:"solution_options,omitempty"`
    SelectedSolutionID string           `json:"selected_solution_id,omitempty"`
    SolutionRationale  string           `json:"solution_rationale,omitempty"`
    Confidence         float64          `json:"confidence,omitempty"`
}

type SolutionOption struct {
    ID               string   `json:"id"`
    Name             string   `json:"name"`
    Description      string   `json:"description,omitempty"`
    ProblemsAddressed []string `json:"problems_addressed,omitempty"`
    Benefits         []string `json:"benefits,omitempty"`
    Tradeoffs        []string `json:"tradeoffs,omitempty"`
    Risks            []string `json:"risks,omitempty"`
    EstimatedEffort  string   `json:"estimated_effort,omitempty"`
}
```

### 1.3 Decisions Definition (from agent-team-prd)

Add to `prd/decisions.go`:

```go
type DecisionsDefinition struct {
    Records []DecisionRecord `json:"records,omitempty"`
}

type DecisionRecord struct {
    ID            string    `json:"id"`
    Decision      string    `json:"decision"`
    Rationale     string    `json:"rationale,omitempty"`
    Alternatives  []string  `json:"alternatives_considered,omitempty"`
    MadeBy        string    `json:"made_by,omitempty"`
    Date          time.Time `json:"date,omitempty"`
    Status        string    `json:"status,omitempty"`
    RelatedIDs    []string  `json:"related_ids,omitempty"`
}
```

### 1.4 Reviews Definition (from agent-team-prd)

Add to `prd/reviews.go`:

```go
type ReviewsDefinition struct {
    ReviewBoardSummary string          `json:"review_board_summary,omitempty"`
    QualityScores      *QualityScores  `json:"quality_scores,omitempty"`
    Decision           ReviewDecision  `json:"decision,omitempty"`
    Blockers           []Blocker       `json:"blockers,omitempty"`
    RevisionTriggers   []RevisionTrigger `json:"revision_triggers,omitempty"`
}

type QualityScores struct {
    ProblemDefinition    float64 `json:"problem_definition"`
    UserUnderstanding    float64 `json:"user_understanding"`
    MarketAwareness      float64 `json:"market_awareness"`
    SolutionFit          float64 `json:"solution_fit"`
    ScopeDiscipline      float64 `json:"scope_discipline"`
    RequirementsQuality  float64 `json:"requirements_quality"`
    UXCoverage           float64 `json:"ux_coverage"`
    TechnicalFeasibility float64 `json:"technical_feasibility"`
    MetricsQuality       float64 `json:"metrics_quality"`
    RiskManagement       float64 `json:"risk_management"`
    OverallScore         float64 `json:"overall_score"`
}

type ReviewDecision string

const (
    ReviewApprove     ReviewDecision = "approve"
    ReviewRevise      ReviewDecision = "revise"
    ReviewReject      ReviewDecision = "reject"
    ReviewHumanReview ReviewDecision = "human_review"
)

type Blocker struct {
    ID          string `json:"id"`
    Category    string `json:"category"`
    Description string `json:"description"`
}

type RevisionTrigger struct {
    IssueID          string `json:"issue_id"`
    Category         string `json:"category"`
    Severity         string `json:"severity"`
    Description      string `json:"description"`
    RecommendedOwner string `json:"recommended_owner,omitempty"`
}
```

### 1.5 Revision History (from agent-team-prd)

Add to `prd/revision.go`:

```go
type RevisionRecord struct {
    Version string              `json:"version"`
    Changes []string            `json:"changes"`
    Trigger RevisionTriggerType `json:"trigger"`
    Date    time.Time           `json:"date"`
}

type RevisionTriggerType string

const (
    TriggerInitial RevisionTriggerType = "initial"
    TriggerReview  RevisionTriggerType = "review"
    TriggerScore   RevisionTriggerType = "score"
    TriggerHuman   RevisionTriggerType = "human"
)
```

### 1.6 Enhanced Problem Definition

Add confidence and evidence to existing types in `prd/document.go`:

```go
// Add to ExecutiveSummary or create ProblemDefinition
type ProblemDefinition struct {
    Statement   string     `json:"statement"`
    UserImpact  string     `json:"user_impact,omitempty"`
    Evidence    []Evidence `json:"evidence,omitempty"`
    Confidence  float64    `json:"confidence,omitempty"`
    RootCauses  []string   `json:"root_causes,omitempty"`
}

type Evidence struct {
    Type       EvidenceType `json:"type"`
    Source     string       `json:"source"`
    Summary    string       `json:"summary,omitempty"`
    SampleSize int          `json:"sample_size,omitempty"`
    Strength   Strength     `json:"strength,omitempty"`
}

type EvidenceType string

const (
    EvidenceInterview      EvidenceType = "interview"
    EvidenceSurvey         EvidenceType = "survey"
    EvidenceAnalytics      EvidenceType = "analytics"
    EvidenceSupportTicket  EvidenceType = "support_ticket"
    EvidenceMarketResearch EvidenceType = "market_research"
    EvidenceAssumption     EvidenceType = "assumption"
)

type Strength string

const (
    StrengthLow    Strength = "low"
    StrengthMedium Strength = "medium"
    StrengthHigh   Strength = "high"
)
```

### 1.7 Update Document Struct

Modify `prd/document.go` to include new fields:

```go
type Document struct {
    // Existing required fields
    Metadata         Metadata         `json:"metadata"`
    ExecutiveSummary ExecutiveSummary `json:"executive_summary"`
    Objectives       Objectives       `json:"objectives"`
    Personas         []Persona        `json:"personas"`
    UserStories      []UserStory      `json:"user_stories"`
    Requirements     Requirements     `json:"requirements"`
    Roadmap          Roadmap          `json:"roadmap"`

    // Existing optional fields
    Assumptions      *AssumptionsConstraints `json:"assumptions_constraints,omitempty"`
    OutOfScope       []string                `json:"out_of_scope,omitempty"`
    TechArchitecture *TechnicalArchitecture  `json:"technical_architecture,omitempty"`
    UXRequirements   *UXRequirements         `json:"ux_requirements,omitempty"`
    Risks            []Risk                  `json:"risks,omitempty"`
    Glossary         []GlossaryTerm          `json:"glossary,omitempty"`
    CustomSections   []CustomSection         `json:"custom_sections,omitempty"`

    // NEW: From agent-team-prd
    Problem          *ProblemDefinition   `json:"problem,omitempty"`
    Market           *MarketDefinition    `json:"market,omitempty"`
    Solution         *SolutionDefinition  `json:"solution,omitempty"`
    Decisions        *DecisionsDefinition `json:"decisions,omitempty"`
    Reviews          *ReviewsDefinition   `json:"reviews,omitempty"`
    RevisionHistory  []RevisionRecord     `json:"revision_history,omitempty"`
}
```

## Phase 2: Add File I/O Functions

### 2.1 Create `prd/io.go`

```go
package prd

// Load reads a Document from a JSON file.
func Load(path string) (*Document, error)

// Save writes a Document to a JSON file.
func Save(doc *Document, path string) error

// New creates a new Document with required fields initialized.
func New(id, title, owner string) *Document

// GenerateID generates a PRD ID based on the current date.
func GenerateID() string
```

## Phase 3: Add Validation Framework

### 3.1 Create `prd/validation.go`

```go
package prd

type ValidationResult struct {
    Valid    bool
    Errors   []ValidationError
    Warnings []ValidationWarning
}

type ValidationError struct {
    Field   string
    Message string
}

type ValidationWarning struct {
    Field   string
    Message string
}

// Validate checks the Document for structural and content issues.
func Validate(doc *Document) *ValidationResult

// validateIDs checks for duplicate and malformed IDs
func (r *ValidationResult) validateIDs(doc *Document)

// validateTraceability checks cross-references between sections
func (r *ValidationResult) validateTraceability(doc *Document)
```

## Phase 4: Add Quality Scoring

### 4.1 Create `prd/scoring.go`

Port the 10-category scoring system, adapted for Document type:

```go
package prd

// Category weights (must sum to 1.0)
var DefaultWeights = []CategoryWeight{
    {Category: "problem_definition", Weight: 0.20},
    {Category: "solution_fit", Weight: 0.15},
    {Category: "user_understanding", Weight: 0.10},
    {Category: "market_awareness", Weight: 0.10},
    {Category: "scope_discipline", Weight: 0.10},
    {Category: "requirements_quality", Weight: 0.10},
    {Category: "metrics_quality", Weight: 0.10},
    {Category: "ux_coverage", Weight: 0.05},
    {Category: "technical_feasibility", Weight: 0.05},
    {Category: "risk_management", Weight: 0.05},
}

// Thresholds
const (
    ThresholdApprove = 8.0
    ThresholdRevise  = 6.5
    ThresholdBlocker = 3.0
)

type ScoringResult struct {
    CategoryScores   []CategoryScore
    WeightedScore    float64
    Decision         string
    Blockers         []string
    RevisionTriggers []RevisionTrigger
    Summary          string
}

type CategoryScore struct {
    Category       string
    Weight         float64
    Score          float64
    MaxScore       float64
    Justification  string
    Evidence       string
    BelowThreshold bool
}

// Score evaluates a Document and returns scoring results.
func Score(doc *Document) *ScoringResult
```

### 4.2 Scoring Category Mappings

| Category | agent-team-prd Field | structured-prd Field |
|----------|---------------------|----------------------|
| problem_definition | `PRD.Problem.PrimaryProblem` | `Document.Problem` OR `Document.ExecutiveSummary` |
| user_understanding | `PRD.Users.Personas` | `Document.Personas` |
| market_awareness | `PRD.Market` | `Document.Market` (NEW) |
| solution_fit | `PRD.Solution` | `Document.Solution` (NEW) |
| scope_discipline | `PRD.GoalsAndNonGoals` | `Document.Objectives` + `Document.OutOfScope` |
| requirements_quality | `PRD.Requirements` | `Document.Requirements` |
| ux_coverage | `PRD.UX` | `Document.UXRequirements` |
| technical_feasibility | `PRD.Technical` | `Document.TechArchitecture` |
| metrics_quality | `PRD.Metrics` | `Document.Objectives.SuccessMetrics` |
| risk_management | `PRD.RisksAndAssumption` | `Document.Risks` + `Document.Assumptions` |

## Phase 5: Add View Generation

### 5.1 Create `prd/views.go`

```go
package prd

// PMView represents the Product Manager view of a PRD.
type PMView struct {
    Title          string
    Status         string
    Owner          string
    Version        string
    ProblemSummary string
    Personas       []PersonaSummary
    Goals          []string
    NonGoals       []string
    Solution       SolutionSummary
    Requirements   RequirementsList
    Metrics        MetricsSummary
    Risks          []RiskSummary
    OpenQuestions  []string
}

// ExecView represents the Executive view of a PRD.
type ExecView struct {
    Header                ExecHeader
    Strengths             []string
    Blockers              []string
    RequiredActions       []ExecAction
    TopRisks              []ExecRisk
    RecommendationSummary string
}

// GeneratePMView creates a PM-friendly view of the Document.
func GeneratePMView(doc *Document) *PMView

// GenerateExecView creates an executive-friendly view of the Document.
func GenerateExecView(doc *Document, scores *ScoringResult) *ExecView

// RenderPMMarkdown generates markdown output for PM view.
func RenderPMMarkdown(view *PMView) string

// RenderExecMarkdown generates markdown output for exec view.
func RenderExecMarkdown(view *ExecView) string
```

## Phase 6: Update Tests

### 6.1 New Test Files

- `prd/io_test.go` - File I/O tests
- `prd/validation_test.go` - Validation tests
- `prd/scoring_test.go` - Scoring tests
- `prd/views_test.go` - View generation tests
- `prd/market_test.go` - Market types tests
- `prd/solution_test.go` - Solution types tests

### 6.2 Update Existing Tests

- `prd/document_test.go` - Add tests for new Document fields
- `prd/completeness_test.go` - Ensure compatibility with new fields

## Phase 7: Update Documentation

### 7.1 Update README.md

- Add sections for scoring, validation, and views
- Document new types and their usage
- Add examples for quality scoring workflow

### 7.2 Update CHANGELOG

- Document all new features
- Note migration path from agent-team-prd

## File Summary

### New Files

| File | Purpose | Lines (est.) |
|------|---------|--------------|
| `prd/market.go` | Market definition types | ~60 |
| `prd/solution.go` | Solution definition types | ~50 |
| `prd/decisions.go` | Decision record types | ~40 |
| `prd/reviews.go` | Review and scoring result types | ~80 |
| `prd/revision.go` | Revision history types | ~30 |
| `prd/io.go` | File I/O functions | ~100 |
| `prd/validation.go` | Validation framework | ~250 |
| `prd/scoring.go` | Quality scoring (10 categories) | ~700 |
| `prd/views.go` | PM and Executive views | ~400 |

### Modified Files

| File | Changes |
|------|---------|
| `prd/document.go` | Add Problem, Market, Solution, Decisions, Reviews, RevisionHistory fields |
| `prd/completeness.go` | Update to handle new optional fields |
| `prd/markdown.go` | Add rendering for new sections |

## Migration Notes

### For agent-team-prd Users

1. Replace `github.com/agentplexus/agent-team-prd/pkg/prd` with `github.com/grokify/structured-plan/prd`
2. Replace `github.com/agentplexus/agent-team-prd/pkg/scoring` with `github.com/grokify/structured-plan/prd` (scoring functions now in prd package)
3. Replace `github.com/agentplexus/agent-team-prd/pkg/views` with `github.com/grokify/structured-plan/prd` (view functions now in prd package)
4. Update type references:
   - `prd.PRD` → `prd.Document`
   - `scoring.Score()` → `prd.Score()`
   - `views.GeneratePMView()` → `prd.GeneratePMView()`

### Backward Compatibility

- All new fields are optional (`omitempty`)
- Existing structured-prd Documents remain valid
- Completeness analysis works with or without new fields
- Scoring gracefully handles missing sections (scores 0 for missing)

## Timeline

| Phase | Effort | Dependencies |
|-------|--------|--------------|
| Phase 1: Types | 1 hour | None |
| Phase 2: File I/O | 30 min | Phase 1 |
| Phase 3: Validation | 1 hour | Phase 1 |
| Phase 4: Scoring | 1.5 hours | Phase 1 |
| Phase 5: Views | 1 hour | Phase 4 |
| Phase 6: Tests | 1 hour | All phases |
| Phase 7: Docs | 30 min | All phases |

**Total estimated effort: ~6.5 hours**
