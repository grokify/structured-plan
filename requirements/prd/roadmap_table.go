package prd

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/grokify/structured-requirements/roadmap"
)

// RoadmapTableOptions configures roadmap table generation.
// This is an alias for backward compatibility.
type RoadmapTableOptions = roadmap.TableOptions

// DefaultRoadmapTableOptions returns sensible defaults for roadmap table generation.
func DefaultRoadmapTableOptions() RoadmapTableOptions {
	return roadmap.DefaultTableOptions()
}

// StatusLegend returns a markdown table explaining the status icons.
func StatusLegend() string {
	return roadmap.StatusLegend()
}

// ToSwimlaneTableWithOKRs generates a markdown table with phases as columns,
// deliverable types as swimlane rows, and OKR swimlanes auto-derived from PhaseTargets.
// This is a PRD-specific method that requires access to Document.Objectives.OKRs.
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
		sb.WriteString(fmt.Sprintf("| **%s** |", roadmap.SwimlaneLabel(swimlane)))

		for _, phase := range d.Roadmap.Phases {
			var items []string
			for _, del := range phase.Deliverables {
				if del.Type == swimlane {
					item := del.Title
					if opts.MaxTitleLen > 0 && len(item) > opts.MaxTitleLen {
						item = item[:opts.MaxTitleLen-3] + "..."
					}
					if opts.IncludeStatus && del.Status != "" {
						item = fmt.Sprintf("%s %s", roadmap.StatusIcon(del.Status), item)
					}
					items = append(items, "• "+item)
				}
			}
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")
	}

	// Add OKR swimlanes if we have OKRs with PhaseTargets
	if hasOKRsWithPhaseTargets {
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
								krLabel = fmt.Sprintf("%s %s", roadmap.PhaseTargetStatusIcon(pt.Status), krLabel)
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
