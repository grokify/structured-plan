# TODO: Rename to structured-plan

This checklist tracks renaming `structured-requirements` to `structured-plan` and restructuring packages.

## Overview

**Current:** `github.com/grokify/structured-requirements`
**Target:** `github.com/grokify/structured-plan`
**CLI:** `splan` (not `splanning`)

### Ecosystem

```
structured-plan       →  Planning artifacts (V2MOM, OKR, PRD, MRD, TRD, Roadmap, Status Reports)
structured-metrics    →  DMAIC metrics (separate repo, internal implementation)
structured-evaluation →  Quality assessment (existing)
structured-changelog  →  Release management (existing)
```

### Why "plan" over "planning"?

- **"Plan" = artifact** (noun) - what you produce
- **"Planning" = activity** (verb) - what you do
- The system is about **artifacts**, not meetings
- **Naming convention**: `structured-*` repos use nouns (`changelog`, `evaluation`, `goals`)
- **Lifecycle symmetry**: `structured-plan` → `structured-changelog` → `structured-evaluation`
- **CLI ergonomics**: `splan` is shorter and cleaner than `splanning`

### Why consolidate into one repo?

These aren't separate tools - they're a **cascading planning system** used together at multiple levels:

```
Organization
├── V2MOM (company vision/values/methods)
├── OKRs (company objectives)
└── Roadmap
    ├── Deliverable: Auth Initiative PRD ─────────┐
    ├── Deliverable: Platform PRD                 │
    └── Deliverable: Analytics PRD                │
                                                  ▼
        PRD (project level) ◄─────────────────────┘
        ├── V2MOM (project methods, cascaded from org)
        ├── OKRs (project KRs, aligned to company OKRs)
        └── Roadmap
            ├── Deliverable: Dashboard Feature
            ├── Deliverable: Reports API
            └── Deliverable: Export Tool
```

**Same types at every level, different content granularity:**

| Level | V2MOM Methods | OKR Objectives | Roadmap Deliverables |
|-------|---------------|----------------|---------------------|
| **Org** | Strategic initiatives | Company OKRs | PRDs, initiatives |
| **Dept** | Department goals | Dept OKRs | Team PRDs |
| **Team** | Team methods | Team OKRs | Features |
| **PRD** | Project methods | Project KRs | Features, APIs, docs |

V2MOM is explicitly a **cascading** framework (Salesforce uses Company→Dept→Team→Individual).
OKRs cascade similarly (Company→Team→Individual alignment).
Roadmaps nest (Portfolio→Product→Release).

**Benefits of single repo:**
- One set of types used at all levels
- Cross-level linking (PRD refs in Portfolio Roadmap, OKR alignment)
- Evolve together as a unified system
- Shared rendering, validation, filtering
- Single versioning and release

## Phase 1: Restructure Packages (before rename)

### 1.1 Create requirements/ directory and move PRD/MRD/TRD

```
Current:                          Target:
structured-requirements/          structured-plan/
├── prd/                          ├── requirements/
├── mrd/                          │   ├── prd/
├── trd/                          │   ├── mrd/
├── common/                       │   └── trd/
└── schema/                       ├── common/
                                  └── schema/
```

Tasks:
- [x] Create `requirements/` directory
- [x] Move `prd/` → `requirements/prd/`
- [x] Move `mrd/` → `requirements/mrd/`
- [x] Move `trd/` → `requirements/trd/`
- [x] Update all internal imports
- [x] Update tests
- [x] Verify build passes

### 1.2 Extract roadmap from PRD to top-level

The roadmap package uses a **unified structure** that works at any level - same types, different content granularity:

| Level | Phases | Deliverables | Rollout | OKRs |
|-------|--------|--------------|---------|------|
| **Product** | Phase 1, Phase 2 | Features, APIs, Docs | Customer groups | Product KRs |
| **Portfolio** | Q1, Q2, H1, H2 | PRDs (or summaries) | Combined rollouts | Org-level KRs |

```
Portfolio Roadmap Example:
├── Q1 2025
│   ├── Deliverables: Auth PRD, API PRD
│   ├── Rollout: Internal beta
│   └── KR: 10% market share
├── Q2 2025
│   ├── Deliverables: Dashboard PRD, Analytics PRD
│   ├── Rollout: Early adopters
│   └── KR: 15% market share
└── ...
```

```
Target:
structured-plan/
├── requirements/
│   └── prd/           # PRD embeds roadmap.Roadmap
├── roadmap/
│   ├── roadmap.go     # Unified Roadmap type (works at any level)
│   ├── phase.go       # Phase type
│   ├── deliverable.go # Deliverable types
│   ├── swimlane.go    # Swimlane table generation
│   ├── legend.go      # Status icons and legend
│   └── table.go       # Table rendering utilities
└── ...
```

#### Unified Roadmap Type

One `Roadmap` type that scales from product to portfolio:

```go
type Roadmap struct {
    Metadata    Metadata    `json:"metadata"`
    Phases      []Phase     `json:"phases"`
    OKRs        []OKR       `json:"okrs,omitempty"`        // Optional OKR alignment
}

type Phase struct {
    ID           string        `json:"id"`
    Name         string        `json:"name"`
    Type         PhaseType     `json:"type,omitempty"`      // quarter, half, sprint, custom
    StartDate    string        `json:"start_date,omitempty"`
    EndDate      string        `json:"end_date,omitempty"`
    Goals        []string      `json:"goals,omitempty"`
    Deliverables []Deliverable `json:"deliverables,omitempty"`
    Rollouts     []Rollout     `json:"rollouts,omitempty"`  // Customer rollout groups
    Tags         []string      `json:"tags,omitempty"`
}

type Deliverable struct {
    ID          string           `json:"id"`
    Title       string           `json:"title"`
    Description string           `json:"description,omitempty"`
    Type        DeliverableType  `json:"type"`    // feature, prd, milestone, etc.
    Status      DeliverableStatus `json:"status,omitempty"`
    PRDRef      string           `json:"prd_ref,omitempty"`  // Reference to PRD (for portfolio level)
    Tags        []string         `json:"tags,omitempty"`
}
```

**At Product level**: Deliverables are features, APIs, docs
**At Portfolio level**: Deliverables are PRDs (with `prd_ref` and optional summary)

Tasks:
- [x] Identify roadmap types to extract from `prd/`:
  - [x] `Phase`, `Deliverable`, `DeliverableType`, `DeliverableStatus`
  - [x] `RoadmapTableOptions`, swimlane table generation
  - [x] `StatusLegend()`, status icons
  - [ ] `Rollout` types (DeliverableRollout) - not yet implemented
- [x] Create `roadmap/` package with unified types
- [ ] Add `DeliverableType` value for `prd` (portfolio-level deliverables)
- [ ] Add `PRDRef` field to Deliverable for portfolio→PRD linking
- [x] Move extracted types to `roadmap/`
- [x] Update `prd.Roadmap` to embed/alias `roadmap.Roadmap`
- [x] Update imports throughout
- [x] Verify build and tests pass

### 1.3 Update CLI structure

```
Current:                          Target:
srequirements prd ...             splan requirements prd ...
srequirements mrd ...             splan requirements mrd ...
srequirements trd ...             splan requirements trd ...
```

Tasks:
- [x] Create `cmd/splan/` with new CLI structure
- [x] Update command hierarchy:
  - [x] Add `requirements` parent command with `req` alias
  - [x] Move `prd`, `mrd`, `trd` under `requirements`
- [ ] Update binary name in goreleaser (no goreleaser config found)
- [x] Keep `srequirements` for backwards compatibility

## Phase 2: Rename Repository

### 2.1 GitHub rename

- [ ] Rename repo: `structured-requirements` → `structured-plan`
- [ ] GitHub auto-creates redirect from old name

### 2.2 Update go.mod

```go
// Before
module github.com/grokify/structured-requirements

// After
module github.com/grokify/structured-plan
```

Tasks:
- [ ] Update `go.mod` module path
- [ ] Update all import paths in all `.go` files
- [ ] Update `go.sum`
- [ ] Run `go mod tidy`
- [ ] Verify build passes

### 2.3 Update documentation

- [ ] Update README.md
- [ ] Update CHANGELOG.md
- [ ] Update any hardcoded repo references
- [ ] Update examples/

## Phase 3: Consolidate OKR and V2MOM

Move OKR and V2MOM from `structured-goals` into `structured-plan`. Single source of truth.

### 3.1 Target structure

```
structured-plan/
├── goals/
│   ├── okr/              # OKR types (from structured-goals + PRD enhancements)
│   │   ├── okr.go
│   │   ├── progress.go   # Score calculation, grading
│   │   └── render/marp/  # Marp slides
│   └── v2mom/            # V2MOM types (from structured-goals)
│       ├── v2mom.go
│       └── render/marp/
├── roadmap/              # Roadmap types (from PRD, made standalone)
│   ├── roadmap.go
│   ├── phase.go
│   ├── deliverable.go
│   └── render/
├── requirements/
│   ├── prd/              # Imports from goals/okr, roadmap/
│   ├── mrd/
│   └── trd/
└── common/
```

**structured-goals disposition:**
- Re-export from structured-plan for backwards compatibility, OR
- Deprecate with migration guide pointing to structured-plan

### 3.2 Analysis: structured-goals vs PRD OKR Types

**PRD OKR Types (prd/document.go):**
| Type | Key Fields |
|------|------------|
| `OKR` | Objective + KeyResults |
| `Objective` | ID, Description, Rationale, AlignedWith, Category, Owner, Timeframe, Tags |
| `KeyResult` | ID, Description, Metric, Baseline, Target, Current, Unit, MeasurementMethod, Owner, Confidence, PhaseTargets, Tags |
| `PhaseTarget` | PhaseID, Target, Status, Actual, Notes |

**structured-goals OKR Types (okr/okr.go):**
| Type | Key Fields |
|------|------------|
| `OKRDocument` | Metadata, Theme, Objectives, Risks, Alignment |
| `Objective` | ID, Title, Description, Owner, Timeframe, Status, KeyResults, **Progress**, Risks, ParentID, AlignedWith[] |
| `KeyResult` | ID, Title, Description, Owner, Metric, Baseline, Target, Current, Unit, **Score**, Confidence, Status, **DueDate** |
| `Risk` | ID, Title, Description, Impact, Likelihood, Mitigation, Status |

**Features in structured-goals NOT in PRD:**
1. **Score (0.0-1.0)** - KeyResult actual achievement score
2. **Progress calculation** - `Objective.CalculateProgress()` averaging KR scores
3. **Score grading** - `ScoreGrade()` returns A-F, `ScoreDescription()` returns text
4. **DueDate** - KeyResult deadline tracking
5. **Richer Alignment** - Array of aligned objective IDs vs single string
6. **Document-level Risks** - Cross-cutting risks at OKRDocument level
7. **Theme** - Annual/quarterly theme for context

**V2MOM Unique Features:**
- Vision/Values - Strategic context PRD lacks
- Projects with ExternalLinks (Jira, Aha, Productboard, Confluence)
- Structure modes (flat/nested/hybrid)
- Terminology modes (v2mom/okr/hybrid for display)

### 3.3 Framework-agnostic Goals wrapper

Support organizations using either V2MOM or OKR with a discriminated union:

```go
// goals/goals.go

type GoalFramework string

const (
    FrameworkOKR   GoalFramework = "okr"
    FrameworkV2MOM GoalFramework = "v2mom"
)

// Goals supports either OKR or V2MOM framework.
type Goals struct {
    Framework GoalFramework `json:"framework"`       // "okr" or "v2mom"
    OKR       *OKRSet       `json:"okr,omitempty"`   // When framework = "okr"
    V2MOM     *V2MOM        `json:"v2mom,omitempty"` // When framework = "v2mom"
}
```

**Abstraction layer for rendering** (roadmap swimlanes, Marp slides):

```go
// goals/abstract.go

// GoalItem = Objective (OKR) or Method (V2MOM)
type GoalItem struct {
    ID, Title, Description, Owner, Status string
    Tags []string
}

// ResultItem = Key Result (OKR) or Measure (V2MOM)
type ResultItem struct {
    ID, Title, Metric, Baseline, Target, Current, Status, PhaseID string
}

// Abstraction methods on Goals
func (g *Goals) GoalItems() []GoalItem      // Objectives or Methods
func (g *Goals) ResultItems() []ResultItem  // Key Results or Measures
func (g *Goals) GoalLabel() string          // "Objectives" or "Methods"
func (g *Goals) ResultLabel() string        // "Key Results" or "Measures"
```

**Roadmap swimlane uses abstraction:**
```go
// Works for both V2MOM and OKR
sb.WriteString(fmt.Sprintf("| **%s** |", goals.GoalLabel()))
for _, item := range goals.GoalItems() { ... }
```

### 3.4 Recommendation: Consolidate into structured-plan

**Goal:** Single source of truth for OKR types - no duplication between repos.

**Approach:** Move OKR/V2MOM into `structured-plan/goals/` and have PRD use those types directly.

```
structured-plan/
├── goals/
│   ├── okr/
│   │   ├── okr.go         # OKR types (from structured-goals)
│   │   ├── progress.go    # Score calculation, grading
│   │   └── render/        # Marp renderer
│   └── v2mom/
│       ├── v2mom.go       # V2MOM types (from structured-goals)
│       └── render/        # Marp renderer
├── requirements/
│   └── prd/
│       ├── document.go    # Uses goals/okr types
│       └── ...
└── ...
```

**Type Reconciliation:**

| structured-goals | PRD current | Consolidated |
|------------------|-------------|--------------|
| `okr.Objective` | `prd.Objective` | `okr.Objective` (add PRD fields) |
| `okr.KeyResult` | `prd.KeyResult` | `okr.KeyResult` (add PhaseTargets) |
| N/A | `prd.PhaseTarget` | `okr.PhaseTarget` (move to okr) |
| `okr.Risk` | `prd.Risk` | Keep both (different contexts) |

**Merged KeyResult type:**
```go
type KeyResult struct {
    ID                string        `json:"id,omitempty"`
    Title             string        `json:"title"`              // from structured-goals
    Description       string        `json:"description,omitempty"`
    Owner             string        `json:"owner,omitempty"`
    Metric            string        `json:"metric,omitempty"`
    Baseline          string        `json:"baseline,omitempty"`
    Target            string        `json:"target"`
    Current           string        `json:"current,omitempty"`
    Unit              string        `json:"unit,omitempty"`
    MeasurementMethod string        `json:"measurement_method,omitempty"` // from PRD
    Score             float64       `json:"score,omitempty"`              // from structured-goals
    Confidence        string        `json:"confidence,omitempty"`         // Low/Medium/High
    Status            string        `json:"status,omitempty"`             // On Track, At Risk, etc.
    DueDate           string        `json:"due_date,omitempty"`           // from structured-goals
    PhaseTargets      []PhaseTarget `json:"phase_targets,omitempty"`      // from PRD
    Tags              []string      `json:"tags,omitempty"`               // from PRD
}
```

**Merged Objective type:**
```go
type Objective struct {
    ID          string      `json:"id,omitempty"`
    Title       string      `json:"title"`                    // from structured-goals (maps to Description)
    Description string      `json:"description,omitempty"`
    Rationale   string      `json:"rationale,omitempty"`      // from PRD
    Owner       string      `json:"owner,omitempty"`
    Timeframe   string      `json:"timeframe,omitempty"`
    Status      string      `json:"status,omitempty"`         // from structured-goals
    Category    string      `json:"category,omitempty"`       // from PRD
    KeyResults  []KeyResult `json:"key_results"`
    Progress    float64     `json:"progress,omitempty"`       // from structured-goals
    Risks       []Risk      `json:"risks,omitempty"`          // from structured-goals
    AlignedWith []string    `json:"aligned_with,omitempty"`   // array from structured-goals
    ParentID    string      `json:"parent_id,omitempty"`      // from structured-goals
    Tags        []string    `json:"tags,omitempty"`           // from PRD
}
```

### 3.5 Migration Tasks

- [x] Analyze structured-goals OKR features vs PRD
- [x] Analyze structured-goals V2MOM features
- [x] Decide on integration approach: Consolidate into structured-plan
- [x] Design framework-agnostic Goals wrapper (V2MOM or OKR)
- [x] Create `goals/` package structure:
  ```
  goals/
  ├── goals.go       # Goals wrapper with Framework discriminator (TODO)
  ├── abstract.go    # GoalItem, ResultItem abstractions (TODO)
  ├── okr/
  │   ├── okr.go     # OKR types (merged from structured-goals)
  │   ├── validation.go
  │   └── render/marp/
  └── v2mom/
      ├── v2mom.go   # V2MOM types (from structured-goals)
      ├── validation.go
      ├── terminology.go
      └── render/marp/
  ```
- [ ] Implement Goals wrapper:
  - [ ] GoalFramework discriminator (okr, v2mom)
  - [ ] GoalItems() abstraction (Objectives or Methods)
  - [ ] ResultItems() abstraction (Key Results or Measures)
  - [ ] GoalLabel()/ResultLabel() for rendering
- [ ] Merge OKR types (combine best of both):
  - [ ] Merge Objective types
  - [ ] Merge KeyResult types (add PhaseTargets)
  - [ ] Move PhaseTarget to okr package
  - [ ] Add progress calculation utilities
  - [ ] Add score grading utilities
- [ ] Move Marp renderers from structured-goals
- [ ] Update PRD:
  - [ ] Replace `Objectives` field with `Goals` field
  - [ ] Remove inline OKR types (use goals/okr)
  - [ ] Remove `GoalsAlignment` wrapper
- [ ] Update roadmap swimlane to use Goals abstraction
- [ ] Update structured-goals to re-export from structured-plan (or deprecate)
- [ ] Update CLI: `splan goals okr`, `splan goals v2mom`

## Phase 4: Status Reports (Weekly/Quarterly)

Add support for operational status reporting within the planning lifecycle.

### 4.1 Rationale

Status reports are **plan tracking**, not separate artifacts:

| Cadence | Purpose | Relationship to Plan |
|---------|---------|---------------------|
| Weekly | Track delivery, slips, blockers | "Are we on track with the plan?" |
| Quarterly | Evaluate OKR/objective completion | "Did we hit planned objectives?" |

### 4.2 Target structure

```
structured-plan/
├── reports/
│   ├── weekly/
│   │   ├── weekly.go        # WeeklyStatus type
│   │   ├── slide.go         # Exec slide generation
│   │   └── render/marp/     # Marp slides
│   └── quarterly/
│       ├── quarterly.go     # QuarterlyReview type
│       ├── okr_summary.go   # OKR scoring summary
│       └── render/marp/     # Marp slides
└── ...
```

### 4.3 Weekly Status Type

```go
type WeeklyStatus struct {
    Metadata      Metadata          `json:"metadata"`
    Period        WeekPeriod        `json:"period"`           // Week start/end
    Highlights    []Highlight       `json:"highlights"`       // Key wins
    Deliverables  []DeliverableStatus `json:"deliverables"`   // Target vs actual dates
    Blockers      []Blocker         `json:"blockers"`
    Risks         []RiskUpdate      `json:"risks,omitempty"`
    NextWeek      []PlannedItem     `json:"next_week,omitempty"`
    Metrics       []MetricSnapshot  `json:"metrics,omitempty"` // Refs to structured-metrics
}

type DeliverableStatus struct {
    ID              string `json:"id"`
    Title           string `json:"title"`
    TargetDate      string `json:"target_date"`      // Original target
    CurrentDate     string `json:"current_date"`     // Current expected date
    DeliveredDate   string `json:"delivered_date,omitempty"`
    SlipDays        int    `json:"slip_days,omitempty"`
    Status          string `json:"status"`           // on_track, at_risk, blocked, complete
    StatusReason    string `json:"status_reason,omitempty"`
    PRDRef          string `json:"prd_ref,omitempty"`
}

type Blocker struct {
    ID          string `json:"id"`
    Description string `json:"description"`
    Owner       string `json:"owner,omitempty"`
    DependsOn   string `json:"depends_on,omitempty"` // External team/resource
    Since       string `json:"since,omitempty"`      // How long blocked
    Resolution  string `json:"resolution,omitempty"`
}
```

### 4.4 Quarterly Review Type

```go
type QuarterlyReview struct {
    Metadata       Metadata           `json:"metadata"`
    Period         QuarterPeriod      `json:"period"`           // Q1 2026, etc.
    OKRSummary     []ObjectiveResult  `json:"okr_summary"`      // OKR scoring
    Achievements   []Achievement      `json:"achievements"`
    Misses         []Miss             `json:"misses,omitempty"`
    Learnings      []Learning         `json:"learnings,omitempty"`
    NextQuarter    []PlannedObjective `json:"next_quarter,omitempty"`
}

type ObjectiveResult struct {
    ObjectiveID   string          `json:"objective_id"`
    Title         string          `json:"title"`
    Score         float64         `json:"score"`          // 0.0-1.0
    Grade         string          `json:"grade"`          // A, B, C, D, F
    KeyResults    []KeyResultScore `json:"key_results"`
    Commentary    string          `json:"commentary,omitempty"`
}
```

### 4.5 Tasks

- [ ] Create `reports/` package structure
- [ ] Create `reports/weekly/weekly.go` with WeeklyStatus type
- [ ] Create `reports/weekly/slide.go` for exec slide generation
- [ ] Create `reports/quarterly/quarterly.go` with QuarterlyReview type
- [ ] Create `reports/quarterly/okr_summary.go` for OKR scoring
- [ ] Add Marp renderers for both report types
- [ ] Add CLI: `splan reports weekly`, `splan reports quarterly`
- [ ] Integration: weekly status can pull from PRD/Roadmap deliverables
- [ ] Integration: quarterly review can pull from OKR scores

## Phase 5: Evaluate structured-roadmap

Separate from above - `structured-roadmap` is project task tracking, not product planning.

- [ ] Consider renaming to `simple-todo` or `structured-tasks`
- [ ] Clarify in README: project tasks vs product roadmap
- [ ] Keep separate from structured-plan

## File Change Summary

### Moves
| From | To |
|------|-----|
| `prd/` | `requirements/prd/` |
| `mrd/` | `requirements/mrd/` |
| `trd/` | `requirements/trd/` |
| `prd/roadmap.go` (types) | `roadmap/` |
| `prd/roadmap_table.go` | `roadmap/` |
| `cmd/srequirements/` | `cmd/splan/` |

### Import Path Changes
| From | To |
|------|-----|
| `structured-requirements/prd` | `structured-plan/requirements/prd` |
| `structured-requirements/mrd` | `structured-plan/requirements/mrd` |
| `structured-requirements/trd` | `structured-plan/requirements/trd` |
| `structured-requirements/common` | `structured-plan/common` |

## Phase 6: Extract Common Types

Move duplicated types from PRD/MRD/TRD to `common/` package.

### 6.1 Types to extract

Currently duplicated across packages:

| Type | PRD | MRD | TRD | Action |
|------|-----|-----|-----|--------|
| `Assumption` | ✅ | ✅ | ✅ | → common/ |
| `Constraint` | ✅ | ❌ | ✅ | → common/ |
| `Risk` | ✅ | ✅ | ✅ | → common/ |
| `GlossaryTerm` | ✅ | ✅ | ✅ | → common/ |
| `CustomSection` | ✅ | ✅ | ✅ | → common/ |
| `Status` | ✅ | ✅ | ✅ | → common/ |
| `OpenItem` | ✅ | ❌ | ❌ | → common/ (useful for all) |
| `DecisionRecord` | ✅ | ❌ | ❌ | → common/ (useful for all) |

### 6.2 Target common/ structure

```
common/
├── person.go        # Person, Approver (existing)
├── status.go        # Status constants (draft, in_review, approved, deprecated)
├── assumption.go    # Assumption
├── constraint.go    # Constraint, ConstraintType
├── risk.go          # Risk, RiskProbability, RiskImpact, RiskStatus
├── decision.go      # OpenItem, Option, DecisionRecord, Resolution
├── glossary.go      # GlossaryTerm
├── custom.go        # CustomSection
└── nongoals.go      # NonGoal (structured, not just []string)
```

### 6.3 Unified types

**OpenItem** (decisions pending):
```go
type OpenItem struct {
    ID           string        `json:"id"`
    Title        string        `json:"title"`
    Description  string        `json:"description,omitempty"`
    Context      string        `json:"context,omitempty"`
    Options      []Option      `json:"options,omitempty"`
    Status       OpenItemStatus `json:"status,omitempty"`      // open, in_discussion, resolved, deferred
    Priority     Priority      `json:"priority,omitempty"`
    Owner        string        `json:"owner,omitempty"`
    Stakeholders []string      `json:"stakeholders,omitempty"`
    DueDate      *time.Time    `json:"due_date,omitempty"`
    Resolution   *Resolution   `json:"resolution,omitempty"`
    RelatedIDs   []string      `json:"related_ids,omitempty"`
    Tags         []string      `json:"tags,omitempty"`
}

type Option struct {
    ID                      string `json:"id"`
    Title                   string `json:"title"`
    Description             string `json:"description,omitempty"`
    Pros                    []string `json:"pros,omitempty"`
    Cons                    []string `json:"cons,omitempty"`
    Effort                  string `json:"effort,omitempty"`  // low, medium, high
    Risk                    string `json:"risk,omitempty"`    // low, medium, high
    Cost                    string `json:"cost,omitempty"`
    Timeline                string `json:"timeline,omitempty"`
    Recommended             bool   `json:"recommended,omitempty"`
    RecommendationRationale string `json:"recommendation_rationale,omitempty"`
}
```

**NonGoal** (explicit out-of-scope with rationale):
```go
type NonGoal struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    Rationale   string `json:"rationale,omitempty"`  // Why it's out of scope
    FuturePhase string `json:"future_phase,omitempty"` // "Phase 2", "v2.0", etc.
    Tags        []string `json:"tags,omitempty"`
}
```

### 6.4 Tasks

- [ ] Create `common/status.go` with Status type and constants
- [ ] Create `common/assumption.go` with unified Assumption type
- [ ] Create `common/constraint.go` with unified Constraint type
- [ ] Create `common/risk.go` with unified Risk type
- [ ] Create `common/decision.go` with OpenItem, Option, DecisionRecord
- [ ] Create `common/glossary.go` with GlossaryTerm
- [ ] Create `common/custom.go` with CustomSection
- [ ] Create `common/nongoals.go` with NonGoal type
- [ ] Update PRD to use common types (type aliases for backwards compat)
- [ ] Update MRD to use common types
- [ ] Update TRD to use common types
- [ ] Add OpenItems to MRD and TRD documents
- [ ] Add NonGoals to MRD (PRD has OutOfScope, TRD has in ExecutiveSummary)

---

## Phase 7: Documentation and Presentation

### 7.1 Marp Presentation

Created: `docs/presentations/structured-plan-overview.md`

Covers:
- The cascading planning model
- V2MOM/OKR/Roadmap at multiple levels
- How PRD fits in
- Unified type system
- Benefits and CLI

Tasks:
- [x] Create initial presentation structure
- [ ] Add screenshots when renderers are working
- [ ] Update CLI examples when implemented
- [ ] Add architecture diagrams
- [ ] Review and polish after Phase 1-3 complete

### 7.2 Documentation Updates

- [ ] Update README.md with new structure
- [ ] Create migration guide from structured-goals
- [ ] Document cascading alignment patterns
- [ ] Add examples for each level (org, team, PRD)

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2025-01-28 | Rename to structured-plan | Consolidate planning docs |
| 2025-01-28 | Move prd/mrd/trd under requirements/ | Clearer hierarchy |
| 2025-01-28 | Extract roadmap to top-level | Shared by multiple doc types, used at org/portfolio/product levels |
| 2025-01-29 | Consolidate OKR/V2MOM into structured-plan | Cascading planning system - same types at all levels (org→dept→team→PRD) |
| 2025-01-29 | structured-plan owns all planning types | Single source of truth; structured-goals re-exports or deprecated |
| 2025-01-29 | Merge OKR implementations | PhaseTargets from PRD + Score/Progress from structured-goals |
| 2025-01-29 | Extract common types | Assumption, Constraint, Risk, OpenItem, Decision, NonGoal shared across PRD/MRD/TRD |
| 2025-01-29 | Framework-agnostic Goals wrapper | Support V2MOM or OKR via discriminated union with abstraction layer for rendering |
| 2025-01-30 | Name: structured-plan (not structured-orgops) | "Plan" is a noun matching naming convention; status reports are plan tracking |
| 2025-01-30 | Keep metrics in separate structured-metrics | DMAIC methodology is different domain; metrics referenced by ID from plans |
| 2025-01-30 | Add weekly/quarterly status reports | Status reports are plan tracking (weekly = on track?, quarterly = objectives hit?) |
| 2025-01-30 | V2MOM: Add Assumptions field | Maps to: Methods→Objectives, Measures→KRs, Obstacles→Constraints, NEW: Assumptions |
| TBD | structured-roadmap → simple-todo | Different purpose (project tasks vs product roadmap) |
