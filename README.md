# Structured Requirements Documents

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

A structured format for requirements documents with Go data types, JSON serialization, and markdown generation. Supports three document types that form a natural workflow:

**MRD** (Market Requirements) → **PRD** (Product Requirements) → **TRD** (Technical Requirements)

## Overview

This library provides comprehensive, machine-readable formats for requirements documents:

- **MRD** - Market Requirements Document: Market analysis, competitive landscape, buyer personas, positioning
- **PRD** - Product Requirements Document: Personas, user stories, functional/non-functional requirements, roadmap
- **TRD** - Technical Requirements Document: Architecture, technology stack, APIs, security design, deployment

Each document type supports:

- Mandatory and optional sections for flexibility
- JSON serialization with Go types
- Markdown generation with Pandoc-compatible YAML frontmatter
- Validation of required fields

## Installation

```bash
go install github.com/grokify/structured-plan/cmd/srequirements@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/structured-plan.git
cd structured-prd
go build -o srequirements ./cmd/srequirements
```

## CLI Usage

The `srequirements` CLI provides commands for each document type:

```bash
# PRD commands
srequirements prd generate <file.json>      # Generate markdown from PRD
srequirements prd validate <file.json>      # Validate PRD structure
srequirements prd check <file.json>         # Check PRD completeness

# MRD commands
srequirements mrd generate <file.json>      # Generate markdown from MRD
srequirements mrd validate <file.json>      # Validate MRD structure

# TRD commands
srequirements trd generate <file.json>      # Generate markdown from TRD
srequirements trd validate <file.json>      # Validate TRD structure
```

### Generate Options

```bash
srequirements prd generate input.json -o output.md    # Custom output path
srequirements prd generate input.json --no-frontmatter # Without YAML frontmatter
srequirements prd generate input.json --margin 1in    # Custom page margin
srequirements prd generate input.json --mainfont Arial # Custom font
```

### Check Options (PRD only)

```bash
srequirements prd check input.json          # Human-readable completeness report
srequirements prd check input.json --json   # JSON output for programmatic use
```

### Examples

```bash
# Validate and generate markdown
srequirements mrd validate examples/agent-platform.mrd.json
srequirements mrd generate examples/agent-platform.mrd.json

srequirements prd validate examples/agent-control-plane.prd.json
srequirements prd generate examples/agent-control-plane.prd.json
srequirements prd check examples/agent-control-plane.prd.json

srequirements trd validate examples/agent-control-plane.trd.json
srequirements trd generate examples/agent-control-plane.trd.json
```

## Library Usage

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/grokify/structured-plan/prd"
    "github.com/grokify/structured-plan/mrd"
    "github.com/grokify/structured-plan/trd"
)

func main() {
    // Create a PRD programmatically
    doc := prd.Document{
        Metadata: prd.Metadata{
            ID:      "prd-001",
            Title:   "User Authentication System",
            Version: "1.0.0",
            Status:  prd.StatusDraft,
            Authors: []prd.Person{{Name: "Jane Doe"}},
        },
        ExecutiveSummary: prd.ExecutiveSummary{
            ProblemStatement: "Users need secure authentication",
            ProposedSolution: "Implement OAuth 2.0 with MFA",
        },
        // ... additional fields
    }

    // Generate markdown
    opts := prd.MarkdownOptions{
        IncludeFrontmatter: true,
        Margin:             "2cm",
    }
    markdown := doc.ToMarkdown(opts)

    // Or marshal to JSON
    data, _ := json.MarshalIndent(doc, "", "  ")
    os.WriteFile("output.prd.json", data, 0600)
}
```

## Evaluation Integration

The library integrates with `structured-evaluation` for standardized quality reports:

```go
import "github.com/grokify/structured-plan/prd"

// Load and score a PRD
doc, _ := prd.Load("my-product.prd.json")

// Convert deterministic scoring to EvaluationReport format
report := prd.ScoreToEvaluationReport(doc, "my-product.prd.json")

// Or generate a template for LLM judge evaluation
template := prd.GenerateEvaluationTemplate(doc, "my-product.prd.json")
```

**Standard Evaluation Categories:**

| Category | Weight | Description |
|----------|--------|-------------|
| problem_definition | 20% | Problem statement clarity and evidence |
| solution_fit | 15% | Solution alignment with problem |
| user_understanding | 10% | Persona depth and user insights |
| market_awareness | 10% | Competitive analysis |
| scope_discipline | 10% | Clear objectives and boundaries |
| requirements_quality | 10% | Functional and non-functional specs |
| metrics_quality | 10% | Success metrics with targets |
| ux_coverage | 5% | Design and accessibility |
| technical_feasibility | 5% | Architecture and integrations |
| risk_management | 5% | Risk identification and mitigation |
```

## Document Types

### MRD - Market Requirements Document

Defines the market opportunity and business justification.

| Section | Required | Description |
|---------|----------|-------------|
| `metadata` | Yes | Document ID, title, version, authors |
| `executive_summary` | Yes | Market opportunity, proposed offering, key findings |
| `market_overview` | Yes | TAM/SAM/SOM, growth rate, trends |
| `target_market` | Yes | Primary/secondary segments, buyer personas |
| `competitive_landscape` | Yes | Competitors, strengths/weaknesses, differentiators |
| `market_requirements` | Yes | Market-level requirements with priorities |
| `positioning` | Yes | Positioning statement, key benefits |
| `go_to_market` | No | Launch strategy, pricing, distribution |
| `success_metrics` | Yes | Revenue targets, market share goals |
| `risks` | No | Market and competitive risks |
| `glossary` | No | Term definitions |

### PRD - Product Requirements Document

Defines what the product should do and for whom.

| Section | Required | Description |
|---------|----------|-------------|
| `metadata` | Yes | Document ID, title, version, authors |
| `executive_summary` | Yes | Problem statement, proposed solution, outcomes |
| `objectives` | Yes | Business objectives, product goals, success metrics |
| `personas` | Yes | User personas with goals and pain points |
| `user_stories` | Yes | User stories with acceptance criteria |
| `requirements.functional` | Yes | Functional requirements (MoSCoW priority) |
| `requirements.non_functional` | Yes | NFRs (performance, security, etc.) |
| `roadmap` | Yes | Phases with deliverables and success criteria |
| `assumptions` | No | Assumptions, constraints, dependencies |
| `out_of_scope` | No | Explicitly excluded items |
| `technical_architecture` | No | System overview, integrations |
| `risks` | No | Product and technical risks |
| `glossary` | No | Term definitions |

### TRD - Technical Requirements Document

Defines how the product will be built.

| Section | Required | Description |
|---------|----------|-------------|
| `metadata` | Yes | Document ID, title, version, authors |
| `executive_summary` | Yes | Purpose, scope, technical approach |
| `architecture` | Yes | Overview, principles, components, data flows |
| `technology_stack` | Yes | Languages, frameworks, databases, infrastructure |
| `api_specifications` | No | API definitions with endpoints |
| `data_model` | No | Entities, attributes, data stores |
| `security_design` | Yes | AuthN, AuthZ, encryption, compliance |
| `performance` | Yes | Performance requirements and benchmarks |
| `scalability` | No | Horizontal/vertical scaling, limits |
| `deployment` | Yes | Environments, strategy, regions |
| `integrations` | No | External system integrations |
| `development` | No | Coding standards, branch strategy |
| `testing` | No | Testing strategy and coverage |
| `risks` | No | Technical risks |
| `glossary` | No | Term definitions |

## File Naming Convention

Use these extensions for automatic type detection:

- `*.prd.json` - Product Requirements Document
- `*.mrd.json` - Market Requirements Document
- `*.trd.json` - Technical Requirements Document

## PRD Details

### Personas

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique persona identifier |
| `name` | Yes | Persona name (e.g., "Developer Dan") |
| `role` | Yes | Job title or role |
| `description` | Yes | Background and context |
| `goals` | Yes | What they want to achieve |
| `pain_points` | Yes | Current frustrations |
| `behaviors` | No | Typical behaviors and patterns |
| `technical_proficiency` | No | Low, Medium, High, Expert |

### User Stories

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique story identifier |
| `persona_id` | Yes | Reference to persona |
| `title` | Yes | Short descriptive title |
| `story` | Yes | "As a [persona], I want [goal] so that [reason]" |
| `acceptance_criteria` | Yes | Testable conditions (Given/When/Then) |
| `priority` | Yes | Critical, High, Medium, Low |
| `phase_id` | Yes | Reference to roadmap phase |

### Roadmap and Swimlane Table

The PRD roadmap is rendered as a swimlane table with phases as columns and deliverable types as rows.

#### Roadmap Structure

```json
{
  "roadmap": {
    "phases": [
      {
        "id": "phase-1",
        "name": "MVP",
        "deliverables": [
          {
            "id": "d1",
            "title": "User Authentication",
            "type": "feature",
            "status": "completed"
          }
        ]
      }
    ]
  }
}
```

#### Ensuring Items Appear in the Roadmap Table

For a deliverable to appear in the swimlane table:

1. **Add to a Phase**: The deliverable must be in a phase's `deliverables` array
2. **Set the Type**: The `type` field determines which swimlane row the item appears in
3. **Set Status (optional)**: The `status` field adds a status icon

#### Deliverable Types (Swimlanes)

| Type Value | Swimlane Row | Description |
|------------|--------------|-------------|
| `feature` | Features | Product features and capabilities |
| `integration` | Integrations | Third-party integrations |
| `infrastructure` | Infrastructure | Platform, CI/CD, monitoring |
| `documentation` | Documentation | User guides, API docs |
| `milestone` | Milestones | Release milestones, checkpoints |
| `rollout` | Rollout | Customer/segment deployment phases |

#### Deliverable Status Icons

| Status Value | Icon | Description |
|--------------|------|-------------|
| `completed` | ✅ | Work is done |
| `in_progress` | 🔄 | Currently being worked on |
| `not_started` | ⏳ | Planned but not started |
| `blocked` | 🚫 | Blocked by dependency |

#### Example: Complete Deliverable

```json
{
  "id": "auth-feature",
  "title": "OAuth 2.0 Authentication",
  "description": "Implement OAuth 2.0 with support for Google and GitHub providers",
  "type": "feature",
  "status": "in_progress"
}
```

This appears in the **Features** row under the phase it belongs to, with a 🔄 icon.

#### Common Issues

| Problem | Cause | Solution |
|---------|-------|----------|
| Item not appearing | Missing or invalid `type` | Set `type` to a valid value |
| Item in wrong row | Wrong `type` value | Check spelling (e.g., `feature` not `Feature`) |
| Item in wrong column | Wrong phase | Move deliverable to correct phase's array |
| No status icon | Missing `status` field | Add `status` field with valid value |

#### Operational Rollout Swimlane

The `rollout` type enables tracking customer/segment deployments across phases. This is useful for phased go-to-market strategies where features are deployed to different customer segments over time.

**Recommended approach for calendar-tied phases:**

When phases represent calendar periods (quarters, months), place rollouts in the phase when deployment actually occurs, not when development completes:

| Swimlane | **Phase 1**<br>Q1 2026 | **Phase 2**<br>Q2 2026 | **Phase 3**<br>Q3 2026 |
|----------|------------------------|------------------------|------------------------|
| **Features** | • Auth<br>• Dashboard | • Reporting<br>• API v2 | • Analytics |
| **Rollout** | | • ✅ Auth → Enterprise<br>• 🔄 Dashboard → Pilot | • Reporting → All |

**Rationale:**

- **Reflects reality** - Development and rollout rarely happen in the same calendar window
- **Shows dependencies** - Clearly communicates "build first, then deploy"
- **Planning accuracy** - Resource allocation aligns with actual work timing

**Naming convention for rollout deliverables:**

Use the `→` notation to distinguish rollout targets:

```json
{
  "id": "rollout-auth-enterprise",
  "title": "Auth → Enterprise customers",
  "description": "Roll out Phase 1 Auth feature to enterprise segment",
  "type": "rollout",
  "status": "completed"
}
```

**Example: Multi-phase customer rollout**

```json
{
  "phases": [
    {
      "id": "phase-1",
      "name": "Q1 2026 - Build",
      "deliverables": [
        { "id": "f1", "title": "User Authentication", "type": "feature", "status": "completed" },
        { "id": "f2", "title": "Dashboard", "type": "feature", "status": "completed" }
      ]
    },
    {
      "id": "phase-2",
      "name": "Q2 2026 - Pilot",
      "deliverables": [
        { "id": "f3", "title": "Reporting", "type": "feature", "status": "in_progress" },
        { "id": "r1", "title": "Auth → Enterprise (Acme, TechCo)", "type": "rollout", "status": "completed" },
        { "id": "r2", "title": "Dashboard → Pilot customers", "type": "rollout", "status": "in_progress" }
      ]
    },
    {
      "id": "phase-3",
      "name": "Q3 2026 - GA",
      "deliverables": [
        { "id": "r3", "title": "Auth → All customers", "type": "rollout", "status": "not_started" },
        { "id": "r4", "title": "Dashboard → All customers", "type": "rollout", "status": "not_started" },
        { "id": "r5", "title": "Reporting → Enterprise", "type": "rollout", "status": "not_started" }
      ]
    }
  ]
}
```

### Non-Functional Requirements

| Category | Description | Example Metrics |
|----------|-------------|-----------------|
| `performance` | Response time, throughput | P95 < 200ms |
| `scalability` | Scaling capability | 10K concurrent users |
| `reliability` | Uptime, MTBF, MTTR | 99.9% uptime |
| `security` | AuthN, AuthZ, encryption | SOC 2 compliance |
| `multi_tenancy` | Tenant isolation | Schema-per-tenant |
| `observability` | Logging, metrics, tracing | 100% trace coverage |
| `compliance` | Regulatory requirements | GDPR, HIPAA |

## MRD Details

### Market Size (TAM/SAM/SOM)

| Field | Required | Description |
|-------|----------|-------------|
| `value` | Yes | Market size (e.g., "$10B") |
| `year` | No | Reference year |
| `source` | No | Data source citation |
| `notes` | No | Additional context |

### Buyer Personas

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier |
| `name` | Yes | Persona name |
| `title` | Yes | Job title |
| `buying_role` | Yes | Decision Maker, Influencer, User, Gatekeeper |
| `budget_authority` | Yes | Has budget authority (boolean) |
| `pain_points` | Yes | Business pain points |
| `goals` | Yes | Business goals |
| `buying_criteria` | No | Purchase decision criteria |

### Competitors

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier |
| `name` | Yes | Competitor name |
| `category` | No | Direct, Indirect, Substitute |
| `strengths` | Yes | Competitive strengths |
| `weaknesses` | Yes | Competitive weaknesses |
| `market_share` | No | Market share percentage |
| `threat_level` | No | High, Medium, Low |

## TRD Details

### Architecture Components

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Component identifier |
| `name` | Yes | Component name |
| `description` | Yes | What it does |
| `type` | No | Service, Library, Database, Queue, etc. |
| `responsibilities` | No | List of responsibilities |
| `dependencies` | No | IDs of dependent components |
| `technology` | No | Implementation technology |

### API Specifications

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | API identifier |
| `name` | Yes | API name |
| `type` | Yes | REST, gRPC, GraphQL, WebSocket |
| `version` | No | API version |
| `base_url` | No | Base URL |
| `auth` | No | Authentication method |
| `endpoints` | No | List of endpoints |

### Security Design

| Field | Required | Description |
|-------|----------|-------------|
| `overview` | Yes | Security approach summary |
| `authentication` | No | AuthN method, provider, MFA |
| `authorization` | No | AuthZ model (RBAC, ABAC) |
| `encryption` | No | At-rest and in-transit encryption |
| `compliance` | No | Compliance standards (SOC2, GDPR) |

## PRD Completeness Check

The `srequirements prd check` command analyzes a PRD for completeness and quality, providing:

- **Overall score** (0-100%) and letter grade (A-F)
- **Section-by-section breakdown** for both required and optional sections
- **Specific recommendations** prioritized by severity

### Scoring

The completeness check evaluates:

| Section | Weight | What's Checked |
|---------|--------|----------------|
| Metadata | 10% | ID, title, version, status, authors |
| Executive Summary | 10% | Problem statement depth, proposed solution, outcomes |
| Objectives | 10% | Business objectives, product goals, success metrics with targets |
| Personas | 10% | Number of personas, completeness of goals/pain points |
| User Stories | 10% | Acceptance criteria coverage, persona/phase linkage |
| Requirements | 10% | Functional/non-functional count, essential NFR categories |
| Roadmap | 10% | Phases with deliverables, success criteria, goals |
| Optional sections | 30% | Assumptions, out of scope, tech architecture, UX, risks, glossary |

### Example Output

```
=============================================================
PRD COMPLETENESS REPORT
=============================================================

Overall Score: 90.8% (Grade: A)
Required Sections: 7/7 complete
Optional Sections: 4/6 complete

-------------------------------------------------------------
SECTION BREAKDOWN
-------------------------------------------------------------

Required Sections:
  [+] Metadata                  100.0% (complete)
  [+] Executive Summary         100.0% (complete)
  [+] Objectives                100.0% (complete)
  [+] Personas                  100.0% (complete)
  [+] User Stories              100.0% (complete)
  [+] Requirements               83.3% (complete)
  [+] Roadmap                   100.0% (complete)

Optional Sections:
  [+] Assumptions & Constraints 100.0% (complete)
  [+] Out of Scope              100.0% (complete)
  [~] Technical Architecture     50.0% (partial)
  [ ] UX Requirements             0.0% (missing)
  [+] Risks                     100.0% (complete)
  [+] Glossary                  100.0% (complete)

-------------------------------------------------------------
RECOMMENDATIONS
-------------------------------------------------------------

HIGH (should fix):
  [*] Requirements: Missing NFR categories: reliability

=============================================================
```

## PDF Generation

The generated markdown includes YAML frontmatter compatible with Pandoc:

```bash
# Generate markdown
srequirements prd generate myproduct.prd.json -o myproduct.md

# Convert to PDF with Pandoc
pandoc myproduct.md -o myproduct.pdf --pdf-engine=xelatex
```

## Examples

See the `examples/` directory for complete examples:

- `examples/agent-platform.mrd.json` - Market requirements for an AI governance platform
- `examples/agent-control-plane.prd.json` - Product requirements for the control plane
- `examples/agent-control-plane.trd.json` - Technical requirements for implementation

## References

### Requirements Documents

- [Modern Analyst - 9 Types of Requirements Documents](https://modernanalyst.com/Resources/Articles/tabid/115/ID/5464/9-Types-Of-Requirements-Documents-What-They-Mean-And-Who-Writes-Them.aspx)
- [Product School - PRD Template](https://productschool.com/blog/product-strategy/product-template-requirements-document-prd)
- [Atlassian - Product Requirements](https://www.atlassian.com/agile/product-management/requirements)

### Technical Documentation

- [AWS - Architecture Documentation](https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html)
- [C4 Model - Software Architecture](https://c4model.com/)
- [ADR - Architecture Decision Records](https://adr.github.io/)

## License

MIT License

 [build-status-svg]: https://github.com/grokify/structured-plan/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/structured-plan/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/structured-plan/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/structured-plan/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/structured-plan
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/structured-plan
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/structured-plan
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/structured-plan
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fstructured-requirements
 [loc-svg]: https://tokei.rs/b1/github/grokify/structured-requirements
 [repo-url]: https://github.com/grokify/structured-plan
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/structured-plan/blob/master/LICENSE
