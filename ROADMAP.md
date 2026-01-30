# Roadmap

This document outlines the development roadmap for structured-plan.

## Current Version: v0.5.0

The v0.5.0 release establishes structured-plan as a unified planning system with:

- Framework-agnostic goals (OKR and V2MOM support)
- Common types extracted for consistency
- Restructured package layout

## Upcoming Releases

### v0.6.0 - Status Reports

Add support for operational status reporting within the planning lifecycle.

**Planned Features:**

- Weekly status reports (`reports/weekly/`)
  - WeeklyStatus type with highlights, deliverables, blockers
  - DeliverableStatus with slip tracking
  - Exec slide generation
  - Marp renderer

- Quarterly reviews (`reports/quarterly/`)
  - QuarterlyReview type with OKR scoring
  - ObjectiveResult with score and grade
  - Achievements and learnings tracking
  - Marp renderer

- CLI commands: `splan reports weekly`, `splan reports quarterly`

- Integration with PRD/Roadmap deliverables and OKR scores

### v0.7.0 - Enhanced CLI

Expand CLI capabilities for the new structure.

**Planned Features:**

- `splan goals okr` - OKR document management
- `splan goals v2mom` - V2MOM document management
- `splan roadmap` - Standalone roadmap management
- Enhanced filtering across all document types

### v0.8.0 - Portfolio Roadmaps

Add support for portfolio-level planning.

**Planned Features:**

- Portfolio roadmap type (phases contain PRD references)
- Cross-PRD deliverable tracking
- Organization-level OKR alignment
- Cascading V2MOM support

## Completed Milestones

### v0.5.0 (2026-01-30)

- [x] Rename repository to structured-plan
- [x] Restructure packages (requirements/, goals/, roadmap/)
- [x] Consolidate OKR types into goals/okr
- [x] Add framework-agnostic Goals wrapper
- [x] Extract common types to common/
- [x] Integrate Goals wrapper into PRD roadmap

### v0.4.0 (2026-01-29)

- [x] Add OKR structure replacing legacy goals
- [x] Add tag-based filtering
- [x] Add roadmap swimlane tables
- [x] Add extended PRD sections

### v0.3.0 (2026-01-26)

- [x] Add structured-evaluation integration
- [x] Add merge command

### v0.2.0 (2026-01-25)

- [x] Add schema package with //go:embed
- [x] Add schema generator from Go types

### v0.1.0 (2026-01-25)

- [x] Initial release with PRD, MRD, TRD
- [x] Marp slide renderer
- [x] Document completeness scoring

## Related Projects

| Project | Description | Status |
|---------|-------------|--------|
| [structured-tasks](https://github.com/grokify/structured-tasks) | AI agent task tracking | Planned rename from structured-roadmap |
| [structured-changelog](https://github.com/grokify/structured-changelog) | Release management | Active |
| [structured-evaluation](https://github.com/grokify/structured-evaluation) | Quality assessment | Active |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to the roadmap.
