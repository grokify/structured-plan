# Release Notes - v0.5.0

**Release Date:** 2026-01-30

## Overview

This release marks the rename from `structured-requirements` to `structured-plan` and introduces a unified planning system. The repository now provides a complete set of planning document types (PRD, MRD, TRD) with framework-agnostic goal support for both OKR and V2MOM methodologies.

## Highlights

- **Repository Rename**: `structured-requirements` → `structured-plan`
- **Framework-Agnostic Goals**: New `goals/` package supporting both OKR and V2MOM
- **Common Types Extraction**: Shared types moved to `common/` package for consistency
- **Dynamic Roadmap Labels**: Roadmap tables automatically use correct terminology (Objectives/Methods)

## Breaking Changes

### Module Path Change

The Go module path has changed:

```go
// Before
import "github.com/grokify/structured-requirements/prd"

// After
import "github.com/grokify/structured-plan/requirements/prd"
```

GitHub automatically redirects the old repository URL, but you should update your imports.

### Package Structure

| Before | After |
|--------|-------|
| `prd/` | `requirements/prd/` |
| `mrd/` | `requirements/mrd/` |
| `trd/` | `requirements/trd/` |
| (new) | `goals/` |
| (new) | `goals/okr/` |
| (new) | `goals/v2mom/` |
| (new) | `roadmap/` |

## New Features

### Framework-Agnostic Goals

The new `goals/` package provides a unified interface for both OKR and V2MOM:

```go
import "github.com/grokify/structured-plan/goals"

// Create OKR-based goals
g := goals.NewOKR(okrSet)

// Create V2MOM-based goals
g := goals.NewV2MOM(v2mom)

// Framework-agnostic access
items := g.GoalItems()     // Returns Objectives (OKR) or Methods (V2MOM)
results := g.ResultItems() // Returns Key Results (OKR) or Measures (V2MOM)

// Dynamic labels
g.GoalLabel()   // "Objectives" or "Methods"
g.ResultLabel() // "Key Results" or "Measures"
```

### PRD ProductGoals Field

PRDs now support the framework-agnostic Goals wrapper:

```go
// New field
doc.ProductGoals = goals.NewOKR(okrSet)  // or goals.NewV2MOM(v2mom)

// Helper method (prefers ProductGoals, falls back to legacy Objectives)
goals := doc.GetProductGoals()
```

### Dynamic Roadmap Labels

Roadmap swimlane tables now use framework-appropriate labels:

```go
// Generates table with "Objectives" / "Key Results" for OKR
// or "Methods" / "Measures" for V2MOM
table := doc.ToSwimlaneTableWithGoals(opts)
```

### Common Types Package

Shared types extracted to `common/` for consistency across PRD/MRD/TRD:

- `Status` - Document lifecycle status
- `Risk`, `RiskProbability`, `RiskImpact`, `RiskStatus`
- `Assumption`, `Constraint`, `ConstraintType`
- `GlossaryTerm`, `CustomSection`
- `OpenItem`, `Option`, `OpenItemResolution`
- `DecisionRecord`, `DecisionStatus`
- `Priority`, `EffortLevel`, `RiskLevel`
- `NonGoal` - Structured out-of-scope items

## Migration Guide

### Update Imports

```bash
# Find and replace in your codebase
find . -name "*.go" -exec sed -i '' \
  's|github.com/grokify/structured-requirements|github.com/grokify/structured-plan|g' {} +
```

### Update Package Paths

```go
// Before
import "github.com/grokify/structured-requirements/prd"

// After
import "github.com/grokify/structured-plan/requirements/prd"
```

### Using Framework-Agnostic Goals

To migrate from OKR-only to framework-agnostic goals:

```go
// Before: OKR-only
for _, okr := range doc.Objectives.OKRs {
    // process objectives
}

// After: Framework-agnostic
goals := doc.GetProductGoals()
for _, item := range goals.GoalItems() {
    // process goals (works with OKR or V2MOM)
}
```

## Ecosystem

The structured-plan ecosystem now includes:

| Repository | Purpose |
|------------|---------|
| `structured-plan` | Planning documents (PRD, MRD, TRD, OKR, V2MOM, Roadmap) |
| `structured-tasks` | AI agent task tracking (renamed from structured-roadmap) |
| `structured-changelog` | Release management |
| `structured-evaluation` | Quality assessment |

## Installation

```bash
go get github.com/grokify/structured-plan@v0.5.0
```

Or install the CLI:

```bash
go install github.com/grokify/structured-plan/cmd/splan@v0.5.0
```

## Full Changelog

See [CHANGELOG.md](CHANGELOG.md) for complete details.
