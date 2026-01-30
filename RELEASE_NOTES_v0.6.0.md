# Release Notes - v0.6.0

**Release Date:** 2026-01-30

## Overview

This release standardizes all JSON field names to **camelCase** for consistency with JSON/JavaScript conventions and modern API standards (OpenAPI, GraphQL, Google APIs).

## Breaking Changes

### JSON Field Names Migrated to camelCase

All 252 JSON field names have been migrated from snake_case to camelCase. This is a **breaking change** that requires updating existing JSON files.

**Examples of changes:**

| Before (snake_case) | After (camelCase) |
|---------------------|-------------------|
| `created_at` | `createdAt` |
| `updated_at` | `updatedAt` |
| `user_stories` | `userStories` |
| `executive_summary` | `executiveSummary` |
| `problem_statement` | `problemStatement` |
| `proposed_solution` | `proposedSolution` |
| `expected_outcomes` | `expectedOutcomes` |
| `success_metrics` | `successMetrics` |
| `key_results` | `keyResults` |
| `non_functional` | `nonFunctional` |
| `out_of_scope` | `outOfScope` |
| `target_audience` | `targetAudience` |

See [MIGRATION_CAMELCASE.md](MIGRATION_CAMELCASE.md) for the complete field mapping.

## Migration Guide

### Automatic Migration Script

A jq-based migration script is provided to convert existing JSON files:

```bash
# Convert a single file (in place)
./scripts/migrate_json_to_camelcase.sh myproduct.prd.json

# Convert to a new file
./scripts/migrate_json_to_camelcase.sh myproduct.prd.json myproduct-new.prd.json

# Batch convert all PRD files
for f in *.prd.json; do
  ./scripts/migrate_json_to_camelcase.sh "$f"
done
```

**Requirements:** `jq` (install with `brew install jq`)

### Manual Migration

If you prefer manual migration, update JSON keys according to the pattern:

```json
// Before
{
  "executive_summary": {
    "problem_statement": "...",
    "proposed_solution": "..."
  },
  "user_stories": [...]
}

// After
{
  "executiveSummary": {
    "problemStatement": "...",
    "proposedSolution": "..."
  },
  "userStories": [...]
}
```

## Rationale

The migration to camelCase provides:

1. **JSON Convention Alignment**: camelCase is the de facto standard for JSON (derived from JavaScript)
2. **API Compatibility**: Matches OpenAPI, GraphQL, and Google API style guides
3. **Smaller Payloads**: No underscores means slightly smaller JSON files
4. **IDE Support**: Better autocomplete in JavaScript/TypeScript editors

## New Files

- `MIGRATION_CAMELCASE.md` - Complete migration documentation with field mapping
- `scripts/migrate_json_to_camelcase.sh` - jq-based migration script for user JSON files
- `scripts/migrate_to_camelcase.sh` - Migration script for Go source files (for reference)

## Installation

### Homebrew (macOS/Linux)

```bash
brew upgrade grokify/tap/splan
```

### Go Install

```bash
go install github.com/grokify/structured-plan/cmd/splan@v0.6.0
```

### Go Module

```bash
go get github.com/grokify/structured-plan@v0.6.0
```

## Full Changelog

See [CHANGELOG.md](CHANGELOG.md) for complete details.
