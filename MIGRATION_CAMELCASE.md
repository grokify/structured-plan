# JSON Naming Convention Migration Plan

## Overview

This document outlines the plan to standardize all JSON field names to **camelCase** for consistency with JSON/JavaScript conventions and modern API standards (OpenAPI, GraphQL, Google APIs).

## Current State

| Style | Count | Status |
|-------|-------|--------|
| camelCase | 19 | ✅ Keep |
| snake_case | 252 | ❌ Migrate |
| lowercase | 150+ | ✅ Keep (single words) |

## Migration Scope

### Files Requiring Changes (27 files)

```
common/
├── decision.go
├── nongoals.go
├── person.go
└── risk.go

goals/okr/
└── okr.go

requirements/mrd/
└── document.go

requirements/prd/
├── appendices.go
├── completeness.go
├── current_state.go
├── document.go
├── goals.go
├── market.go
├── optional_sections.go
├── persona.go
├── persona_library.go
├── prfaq.go
├── problem.go
├── requirements.go
├── reviews.go
├── scoring.go
├── security_model.go
├── six_pager.go
├── solution.go
├── user_story.go
└── views.go

requirements/trd/
└── document.go

roadmap/
└── roadmap.go
```

## Complete Tag Migration Map

### A-C

| snake_case | camelCase |
|------------|-----------|
| `acceptance_criteria` | `acceptanceCriteria` |
| `access_control` | `accessControl` |
| `as_a` | `asA` |
| `at_rest` | `atRest` |
| `audit_logging` | `auditLogging` |
| `below_threshold` | `belowThreshold` |
| `budget_authority` | `budgetAuthority` |
| `business_goals` | `businessGoals` |
| `buying_role` | `buyingRole` |
| `call_to_action` | `callToAction` |
| `category_scores` | `categoryScores` |
| `collection_interval` | `collectionInterval` |
| `competitive_landscape` | `competitiveLandscape` |
| `confidence_level` | `confidenceLevel` |
| `correlation_id` | `correlationId` |
| `created_at` | `createdAt` |
| `current_alternatives` | `currentAlternatives` |
| `current_value` | `currentValue` |
| `customer_faqs` | `customerFaqs` |
| `customer_problem` | `customerProblem` |
| `custom_sections` | `customSections` |

### D-I

| snake_case | camelCase |
|------------|-----------|
| `data_segregation` | `dataSegregation` |
| `due_date` | `dueDate` |
| `executive_summary` | `executiveSummary` |
| `expected_outcomes` | `expectedOutcomes` |
| `export_format` | `exportFormat` |
| `go_to_market` | `goToMarket` |
| `how_it_works` | `howItWorks` |
| `i_want` | `iWant` |
| `in_scope` | `inScope` |
| `in_transit` | `inTransit` |
| `internal_faqs` | `internalFaqs` |
| `is_primary` | `isPrimary` |
| `isolation_model` | `isolationModel` |
| `issue_id` | `issueId` |

### K-M

| snake_case | camelCase |
|------------|-----------|
| `key_benefits` | `keyBenefits` |
| `key_features` | `keyFeatures` |
| `key_findings` | `keyFindings` |
| `key_management` | `keyManagement` |
| `key_results` | `keyResults` |
| `key_threats` | `keyThreats` |
| `log_levels` | `logLevels` |
| `market_awareness` | `marketAwareness` |
| `market_opportunity` | `marketOpportunity` |
| `market_overview` | `marketOverview` |
| `market_requirements` | `marketRequirements` |
| `max_points` | `maxPoints` |
| `max_score` | `maxScore` |
| `metrics_quality` | `metricsQuality` |

### N-P

| snake_case | camelCase |
|------------|-----------|
| `non_functional` | `nonFunctional` |
| `non_goals` | `nonGoals` |
| `optional_complete` | `optionalComplete` |
| `optional_total` | `optionalTotal` |
| `out_of_scope` | `outOfScope` |
| `overall_decision` | `overallDecision` |
| `overall_score` | `overallScore` |
| `pain_points` | `painPoints` |
| `persona_id` | `personaId` |
| `phase_id` | `phaseId` |
| `prd_id` | `prdId` |
| `press_release` | `pressRelease` |
| `primary_metric` | `primaryMetric` |
| `primary_segments` | `primarySegments` |
| `problem_definition` | `problemDefinition` |
| `problem_solved` | `problemSolved` |
| `problem_statement` | `problemStatement` |
| `problem_summary` | `problemSummary` |
| `propagation_format` | `propagationFormat` |
| `proposed_offering` | `proposedOffering` |
| `proposed_solution` | `proposedSolution` |

### R-S

| snake_case | camelCase |
|------------|-----------|
| `recommendation_summary` | `recommendationSummary` |
| `required_actions` | `requiredActions` |
| `required_complete` | `requiredComplete` |
| `required_total` | `requiredTotal` |
| `requirements_quality` | `requirementsQuality` |
| `retention_period` | `retentionPeriod` |
| `revision_triggers` | `revisionTriggers` |
| `risk_management` | `riskManagement` |
| `sampling_rate` | `samplingRate` |
| `schema_version` | `schemaVersion` |
| `scope_discipline` | `scopeDiscipline` |
| `secondary_metrics` | `secondaryMetrics` |
| `security_design` | `securityDesign` |
| `slo_target` | `sloTarget` |
| `so_that` | `soThat` |
| `solution_fit` | `solutionFit` |
| `start_date` | `startDate` |
| `success_criteria` | `successCriteria` |
| `success_metrics` | `successMetrics` |

### T-Z

| snake_case | camelCase |
|------------|-----------|
| `target_audience` | `targetAudience` |
| `target_market` | `targetMarket` |
| `target_state` | `targetState` |
| `target_uptime` | `targetUptime` |
| `technical_approach` | `technicalApproach` |
| `technical_feasibility` | `technicalFeasibility` |
| `technology_stack` | `technologyStack` |
| `threat_actors` | `threatActors` |
| `threat_model` | `threatModel` |
| `top_risks` | `topRisks` |
| `updated_at` | `updatedAt` |
| `user_stories` | `userStories` |
| `user_story_ids` | `userStoryIds` |
| `user_understanding` | `userUnderstanding` |
| `ux_coverage` | `uxCoverage` |
| `weighted_score` | `weightedScore` |

## Implementation Plan

### Phase 1: Run Migration Script

The migration script is located at `scripts/migrate_to_camelcase.sh`:

```bash
# From repository root
./scripts/migrate_to_camelcase.sh
```

This script updates all 252 snake_case JSON tags in Go source files to camelCase.

### Phase 2: Update Example Files

Update all example JSON files in `examples/` directory to use camelCase.

### Phase 3: Update JSON Schema

Regenerate JSON schemas to reflect new field names.

### Phase 4: Update Documentation

Update README.md and any documentation showing JSON examples.

### Phase 5: Testing

1. Run `go test ./...` to ensure all tests pass
2. Verify example files parse correctly
3. Generate markdown from updated examples

## Breaking Change Notice

This is a **breaking change** for existing JSON files. Users must update their JSON files to use camelCase field names.

### Migration Script for User JSON Files

A jq-based migration script is provided at `scripts/migrate_json_to_camelcase.sh`:

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

Requires `jq` (install with `brew install jq`).

## Version Impact

- **Target Version**: v0.6.0
- **Type**: Breaking change
- **Migration Required**: Yes, for all existing JSON files

## Checklist

- [x] Run migration script on Go source files
- [x] Update example JSON files
- [x] Update test files with camelCase JSON
- [x] Regenerate JSON schemas (PRD only - MRD/TRD not yet implemented)
- [x] Update README with new JSON examples
- [x] Add migration script to repository
- [x] Update CHANGELOG.json
- [x] Create RELEASE_NOTES_v0.6.0.md documenting breaking change
- [x] Run full test suite
- [x] Verify golangci-lint passes
