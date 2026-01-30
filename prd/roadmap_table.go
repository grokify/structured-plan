package prd

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// RoadmapTableOptions configures roadmap table generation.
type RoadmapTableOptions struct {
	// IncludeStatus adds status indicators to deliverables
	IncludeStatus bool
	// IncludeEmptySwimlanes shows rows even if no deliverables of that type exist
	IncludeEmptySwimlanes bool
	// SwimlaneOrder specifies the order of swimlanes (nil = alphabetical)
	SwimlaneOrder []DeliverableType
	// MaxTitleLen truncates deliverable titles (0 = no limit)
	MaxTitleLen int
	// IncludeOKRs adds Objectives and Key Results swimlanes derived from PhaseTargets
	IncludeOKRs bool
}

// DefaultRoadmapTableOptions returns sensible defaults for roadmap table generation.
func DefaultRoadmapTableOptions() RoadmapTableOptions {
	return RoadmapTableOptions{
		IncludeStatus:         true,
		IncludeEmptySwimlanes: false,
		SwimlaneOrder: []DeliverableType{
			DeliverableFeature,
			DeliverableIntegration,
			DeliverableInfrastructure,
			DeliverableDocumentation,
			DeliverableMilestone,
			DeliverableRollout,
		},
		MaxTitleLen: 0, // No truncation by default
	}
}

// ToSwimlaneTable generates a markdown table with phases as columns and
// deliverable types as swimlane rows.
//
// Example output:
//
//	| Swimlane       | **Phase 1**<br>Foundation | **Phase 2**<br>Core Features |
//	|----------------|---------------------------|------------------------------|
//	| Features       | • Auth<br>• Search        | • Dashboard                  |
//	| Infrastructure | • CI/CD                   | • Monitoring                 |
func (r *Roadmap) ToSwimlaneTable(opts RoadmapTableOptions) string {
	if len(r.Phases) == 0 {
		return ""
	}

	// Collect all unique swimlanes (deliverable types) across all phases
	swimlaneSet := make(map[DeliverableType]bool)
	for _, phase := range r.Phases {
		for _, del := range phase.Deliverables {
			swimlaneSet[del.Type] = true
		}
	}

	// Determine swimlane order
	var swimlanes []DeliverableType
	if len(opts.SwimlaneOrder) > 0 {
		for _, st := range opts.SwimlaneOrder {
			if swimlaneSet[st] || opts.IncludeEmptySwimlanes {
				swimlanes = append(swimlanes, st)
			}
		}
		// Add any swimlanes not in the specified order
		for st := range swimlaneSet {
			if !slices.Contains(opts.SwimlaneOrder, st) {
				swimlanes = append(swimlanes, st)
			}
		}
	} else {
		for st := range swimlaneSet {
			swimlanes = append(swimlanes, st)
		}
		sort.Slice(swimlanes, func(i, j int) bool {
			return string(swimlanes[i]) < string(swimlanes[j])
		})
	}

	if len(swimlanes) == 0 {
		return ""
	}

	var sb strings.Builder

	// Header row: | Swimlane | **Phase 1**<br>Description | **Phase 2**<br>Description | ...
	sb.WriteString("| Swimlane |")
	for i, phase := range r.Phases {
		// Format: **Phase N**<br>Phase Name/Description
		header := fmt.Sprintf(" **Phase %d**<br>%s |", i+1, phase.Name)
		sb.WriteString(header)
	}
	sb.WriteString("\n")

	// Separator row
	sb.WriteString("|----------|")
	for range r.Phases {
		sb.WriteString("----------|")
	}
	sb.WriteString("\n")

	// Data rows: one per swimlane
	for _, swimlane := range swimlanes {
		sb.WriteString(fmt.Sprintf("| **%s** |", swimlaneLabel(swimlane)))

		for _, phase := range r.Phases {
			// Collect deliverables of this type in this phase
			var items []string
			for _, del := range phase.Deliverables {
				if del.Type == swimlane {
					item := del.Title
					if opts.MaxTitleLen > 0 && len(item) > opts.MaxTitleLen {
						item = item[:opts.MaxTitleLen-3] + "..."
					}
					if opts.IncludeStatus && del.Status != "" {
						item = fmt.Sprintf("%s %s", statusIcon(del.Status), item)
					}
					// Add bullet point prefix
					items = append(items, "• "+item)
				}
			}
			// Join with <br> for line breaks within cell
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ToPhaseTable generates a traditional phase-based table showing
// each phase with its deliverables listed.
//
// Example output:
//
//	| Phase   | Status      | Deliverables                          |
//	|---------|-------------|---------------------------------------|
//	| Phase 1 | In Progress | • ✅ Auth<br>• 🔄 Search<br>• ⏳ Docs |
//	| Phase 2 | Planned     | • Dashboard<br>• Monitoring           |
func (r *Roadmap) ToPhaseTable(opts RoadmapTableOptions) string {
	if len(r.Phases) == 0 {
		return ""
	}

	var sb strings.Builder

	// Header
	sb.WriteString("| Phase | Status | Deliverables |\n")
	sb.WriteString("|-------|--------|---------------|\n")

	for i, phase := range r.Phases {
		var items []string
		for _, del := range phase.Deliverables {
			item := del.Title
			if opts.MaxTitleLen > 0 && len(item) > opts.MaxTitleLen {
				item = item[:opts.MaxTitleLen-3] + "..."
			}
			if opts.IncludeStatus && del.Status != "" {
				item = fmt.Sprintf("%s %s", statusIcon(del.Status), item)
			}
			// Add bullet point prefix
			items = append(items, "• "+item)
		}

		status := string(phase.Status)
		if status == "" {
			status = "planned"
		}

		// Join with <br> for line breaks within cell
		sb.WriteString(fmt.Sprintf("| **Phase %d**<br>%s | %s | %s |\n",
			i+1, phase.Name, status, strings.Join(items, "<br>")))
	}

	return sb.String()
}

// swimlaneLabel converts a DeliverableType to a human-readable label.
func swimlaneLabel(dt DeliverableType) string {
	switch dt {
	case DeliverableFeature:
		return "Features"
	case DeliverableDocumentation:
		return "Documentation"
	case DeliverableInfrastructure:
		return "Infrastructure"
	case DeliverableIntegration:
		return "Integrations"
	case DeliverableMilestone:
		return "Milestones"
	case DeliverableRollout:
		return "Rollout"
	default:
		// Capitalize the first letter
		s := string(dt)
		if len(s) > 0 {
			return strings.ToUpper(s[:1]) + s[1:]
		}
		return s
	}
}

// statusIcon returns an emoji/icon for the deliverable status.
func statusIcon(status DeliverableStatus) string {
	switch status {
	case DeliverableCompleted:
		return "✅"
	case DeliverableInProgress:
		return "🔄"
	case DeliverableBlocked:
		return "🚫"
	case DeliverableNotStarted:
		return "⏳"
	default:
		return ""
	}
}

// phaseTargetStatusIcon returns an emoji/icon for the phase target status.
func phaseTargetStatusIcon(status string) string {
	switch status {
	case "achieved":
		return "✅"
	case "in_progress":
		return "🔄"
	case "missed":
		return "❌"
	case "not_started":
		return "⏳"
	default:
		return ""
	}
}

// StatusLegend returns a markdown table explaining the status icons.
func StatusLegend() string {
	return `| Icon | Status |
|------|--------|
| ✅ | Completed / Achieved |
| 🔄 | In Progress |
| ⏳ | Not Started |
| 🚫 | Blocked |
| ❌ | Missed |
`
}

// ToSwimlaneTableWithOKRs generates a markdown table with phases as columns,
// deliverable types as swimlane rows, and OKR swimlanes auto-derived from PhaseTargets.
//
// Example output:
//
//	| Swimlane       | **Phase 1**<br>Q1 2026 | **Phase 2**<br>Q2 2026 |
//	|----------------|------------------------|------------------------|
//	| Features       | • Auth<br>• Dashboard  | • Reporting            |
//	| Objectives     | • O1: Market leader    |                        |
//	| Key Results    | • KR1.1: Share → 15%   | • KR1.1: Share → 20%   |
func (d *Document) ToSwimlaneTableWithOKRs(opts RoadmapTableOptions) string {
	if len(d.Roadmap.Phases) == 0 {
		return ""
	}

	// Build phase ID to index map
	phaseIDToIndex := make(map[string]int)
	for i, phase := range d.Roadmap.Phases {
		phaseIDToIndex[phase.ID] = i
	}

	// Collect all unique swimlanes (deliverable types) across all phases
	swimlaneSet := make(map[DeliverableType]bool)
	for _, phase := range d.Roadmap.Phases {
		for _, del := range phase.Deliverables {
			swimlaneSet[del.Type] = true
		}
	}

	// Check if we have OKRs with PhaseTargets
	hasOKRsWithPhaseTargets := false
	if opts.IncludeOKRs {
		for _, okr := range d.Objectives.OKRs {
			for _, kr := range okr.KeyResults {
				if len(kr.PhaseTargets) > 0 {
					hasOKRsWithPhaseTargets = true
					break
				}
			}
			if hasOKRsWithPhaseTargets {
				break
			}
		}
	}

	// Determine swimlane order
	var swimlanes []DeliverableType
	if len(opts.SwimlaneOrder) > 0 {
		for _, st := range opts.SwimlaneOrder {
			if swimlaneSet[st] || opts.IncludeEmptySwimlanes {
				swimlanes = append(swimlanes, st)
			}
		}
		// Add any swimlanes not in the specified order
		for st := range swimlaneSet {
			if !slices.Contains(opts.SwimlaneOrder, st) {
				swimlanes = append(swimlanes, st)
			}
		}
	} else {
		for st := range swimlaneSet {
			swimlanes = append(swimlanes, st)
		}
		sort.Slice(swimlanes, func(i, j int) bool {
			return string(swimlanes[i]) < string(swimlanes[j])
		})
	}

	var sb strings.Builder

	// Header row: | Swimlane | **Phase 1**<br>Description | **Phase 2**<br>Description | ...
	sb.WriteString("| Swimlane |")
	for i, phase := range d.Roadmap.Phases {
		header := fmt.Sprintf(" **Phase %d**<br>%s |", i+1, phase.Name)
		sb.WriteString(header)
	}
	sb.WriteString("\n")

	// Separator row
	sb.WriteString("|----------|")
	for range d.Roadmap.Phases {
		sb.WriteString("----------|")
	}
	sb.WriteString("\n")

	// Data rows: one per deliverable swimlane
	for _, swimlane := range swimlanes {
		sb.WriteString(fmt.Sprintf("| **%s** |", swimlaneLabel(swimlane)))

		for _, phase := range d.Roadmap.Phases {
			var items []string
			for _, del := range phase.Deliverables {
				if del.Type == swimlane {
					item := del.Title
					if opts.MaxTitleLen > 0 && len(item) > opts.MaxTitleLen {
						item = item[:opts.MaxTitleLen-3] + "..."
					}
					if opts.IncludeStatus && del.Status != "" {
						item = fmt.Sprintf("%s %s", statusIcon(del.Status), item)
					}
					items = append(items, "• "+item)
				}
			}
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")
	}

	// Add OKR swimlanes if enabled and we have OKRs with PhaseTargets
	if opts.IncludeOKRs && hasOKRsWithPhaseTargets {
		// Objectives swimlane: show objective in every phase where any of its KRs has a target
		sb.WriteString("| **Objectives** |")
		objectivesByPhase := make(map[int][]string) // phase index -> objective descriptions
		for i, okr := range d.Objectives.OKRs {
			// Find all phases where this objective's KRs have targets
			phasesWithKRs := make(map[int]bool)
			for _, kr := range okr.KeyResults {
				for _, pt := range kr.PhaseTargets {
					if idx, ok := phaseIDToIndex[pt.PhaseID]; ok {
						phasesWithKRs[idx] = true
					}
				}
			}
			// Add objective label to each phase where its KRs appear
			objLabel := fmt.Sprintf("O%d: %s", i+1, okr.Objective.Description)
			for phaseIdx := range phasesWithKRs {
				objectivesByPhase[phaseIdx] = append(objectivesByPhase[phaseIdx], "• "+objLabel)
			}
		}
		for i := range d.Roadmap.Phases {
			items := objectivesByPhase[i]
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")

		// Key Results swimlane: show KR targets for each phase
		sb.WriteString("| **Key Results** |")
		for phaseIdx, phase := range d.Roadmap.Phases {
			var items []string
			for i, okr := range d.Objectives.OKRs {
				for j, kr := range okr.KeyResults {
					for _, pt := range kr.PhaseTargets {
						if pt.PhaseID == phase.ID {
							// Format: KR1.1: Description → Target
							krLabel := fmt.Sprintf("KR%d.%d: %s → %s",
								i+1, j+1,
								kr.Description,
								pt.Target)
							if opts.IncludeStatus && pt.Status != "" {
								krLabel = fmt.Sprintf("%s %s", phaseTargetStatusIcon(pt.Status), krLabel)
							}
							items = append(items, "• "+krLabel)
						}
					}
				}
			}
			_ = phaseIdx // unused but kept for clarity
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

