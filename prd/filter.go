package prd

// FilterByTags returns a new Document containing only entities that have at least
// one of the specified tags (OR logic). If no tags are provided, returns a copy
// of the original document. Metadata tags are always preserved.
func (d Document) FilterByTags(tags ...string) Document {
	if len(tags) == 0 {
		return d
	}

	filtered := Document{
		Metadata:         d.Metadata,
		ExecutiveSummary: d.ExecutiveSummary,
		Objectives:       Objectives{}, // Will be filtered below
	}

	// Filter Personas
	filtered.Personas = filterSliceByTags(d.Personas, tags, func(p Persona) []string { return p.Tags })

	// Filter UserStories
	filtered.UserStories = filterSliceByTags(d.UserStories, tags, func(us UserStory) []string { return us.Tags })

	// Filter Requirements
	filtered.Requirements = Requirements{
		Functional:    filterSliceByTags(d.Requirements.Functional, tags, func(fr FunctionalRequirement) []string { return fr.Tags }),
		NonFunctional: filterSliceByTags(d.Requirements.NonFunctional, tags, func(nfr NonFunctionalRequirement) []string { return nfr.Tags }),
	}

	// Filter Roadmap
	filteredPhases := make([]Phase, 0, len(d.Roadmap.Phases))
	for _, phase := range d.Roadmap.Phases {
		if hasAnyTag(phase.Tags, tags) {
			filteredPhases = append(filteredPhases, phase)
		} else {
			// Check if any deliverables match
			filteredDeliverables := filterSliceByTags(phase.Deliverables, tags, func(del Deliverable) []string { return del.Tags })
			if len(filteredDeliverables) > 0 {
				phaseCopy := phase
				phaseCopy.Deliverables = filteredDeliverables
				filteredPhases = append(filteredPhases, phaseCopy)
			}
		}
	}
	filtered.Roadmap = Roadmap{Phases: filteredPhases}

	// Filter Objectives.OKRs
	filteredOKRs := filterOKRsByTags(d.Objectives.OKRs, tags)
	if len(filteredOKRs) > 0 {
		filtered.Objectives.OKRs = filteredOKRs
	}

	// Filter optional sections

	// Assumptions (no tags on Assumption struct, keep all if any other content matches)
	filtered.Assumptions = d.Assumptions

	// OutOfScope (no tags, keep all)
	filtered.OutOfScope = d.OutOfScope

	// TechArchitecture (no tags on subsections, keep all)
	filtered.TechArchitecture = d.TechArchitecture

	// UXRequirements (no tags, keep all)
	filtered.UXRequirements = d.UXRequirements

	// Filter Risks
	filtered.Risks = filterSliceByTags(d.Risks, tags, func(r Risk) []string { return r.Tags })

	// Glossary (no tags, keep all)
	filtered.Glossary = d.Glossary

	// CustomSections (no tags, keep all)
	filtered.CustomSections = d.CustomSections

	// Extended sections
	filtered.Problem = d.Problem
	filtered.Market = d.Market
	filtered.Solution = d.Solution
	filtered.Decisions = d.Decisions
	filtered.Reviews = d.Reviews
	filtered.RevisionHistory = d.RevisionHistory
	filtered.Goals = d.Goals

	return filtered
}

// hasAnyTag returns true if entityTags contains at least one of filterTags.
func hasAnyTag(entityTags, filterTags []string) bool {
	tagSet := make(map[string]struct{}, len(filterTags))
	for _, t := range filterTags {
		tagSet[t] = struct{}{}
	}
	for _, et := range entityTags {
		if _, ok := tagSet[et]; ok {
			return true
		}
	}
	return false
}

// filterSliceByTags filters a slice of items, keeping only those with at least one matching tag.
func filterSliceByTags[T any](items []T, tags []string, getTags func(T) []string) []T {
	if len(items) == 0 {
		return nil
	}
	result := make([]T, 0, len(items))
	for _, item := range items {
		if hasAnyTag(getTags(item), tags) {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// FilterByTagsAll returns a new Document containing only entities that have ALL
// of the specified tags (AND logic). If no tags are provided, returns a copy
// of the original document.
func (d Document) FilterByTagsAll(tags ...string) Document {
	if len(tags) == 0 {
		return d
	}

	filtered := Document{
		Metadata:         d.Metadata,
		ExecutiveSummary: d.ExecutiveSummary,
		Objectives:       Objectives{},
	}

	// Filter Personas
	filtered.Personas = filterSliceByTagsAll(d.Personas, tags, func(p Persona) []string { return p.Tags })

	// Filter UserStories
	filtered.UserStories = filterSliceByTagsAll(d.UserStories, tags, func(us UserStory) []string { return us.Tags })

	// Filter Requirements
	filtered.Requirements = Requirements{
		Functional:    filterSliceByTagsAll(d.Requirements.Functional, tags, func(fr FunctionalRequirement) []string { return fr.Tags }),
		NonFunctional: filterSliceByTagsAll(d.Requirements.NonFunctional, tags, func(nfr NonFunctionalRequirement) []string { return nfr.Tags }),
	}

	// Filter Roadmap
	filteredPhases := make([]Phase, 0, len(d.Roadmap.Phases))
	for _, phase := range d.Roadmap.Phases {
		if hasAllTags(phase.Tags, tags) {
			filteredPhases = append(filteredPhases, phase)
		} else {
			// Check if any deliverables match
			filteredDeliverables := filterSliceByTagsAll(phase.Deliverables, tags, func(del Deliverable) []string { return del.Tags })
			if len(filteredDeliverables) > 0 {
				phaseCopy := phase
				phaseCopy.Deliverables = filteredDeliverables
				filteredPhases = append(filteredPhases, phaseCopy)
			}
		}
	}
	filtered.Roadmap = Roadmap{Phases: filteredPhases}

	// Filter Objectives.OKRs
	filteredOKRs := filterOKRsByTagsAll(d.Objectives.OKRs, tags)
	if len(filteredOKRs) > 0 {
		filtered.Objectives.OKRs = filteredOKRs
	}

	// Filter optional sections (same as FilterByTags - no tags on these)
	filtered.Assumptions = d.Assumptions
	filtered.OutOfScope = d.OutOfScope
	filtered.TechArchitecture = d.TechArchitecture
	filtered.UXRequirements = d.UXRequirements
	filtered.Risks = filterSliceByTagsAll(d.Risks, tags, func(r Risk) []string { return r.Tags })
	filtered.Glossary = d.Glossary
	filtered.CustomSections = d.CustomSections

	// Extended sections
	filtered.Problem = d.Problem
	filtered.Market = d.Market
	filtered.Solution = d.Solution
	filtered.Decisions = d.Decisions
	filtered.Reviews = d.Reviews
	filtered.RevisionHistory = d.RevisionHistory
	filtered.Goals = d.Goals

	return filtered
}

// hasAllTags returns true if entityTags contains ALL of filterTags.
func hasAllTags(entityTags, filterTags []string) bool {
	tagSet := make(map[string]struct{}, len(entityTags))
	for _, t := range entityTags {
		tagSet[t] = struct{}{}
	}
	for _, ft := range filterTags {
		if _, ok := tagSet[ft]; !ok {
			return false
		}
	}
	return true
}

// filterSliceByTagsAll filters a slice of items, keeping only those with ALL matching tags.
func filterSliceByTagsAll[T any](items []T, tags []string, getTags func(T) []string) []T {
	if len(items) == 0 {
		return nil
	}
	result := make([]T, 0, len(items))
	for _, item := range items {
		if hasAllTags(getTags(item), tags) {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// filterOKRsByTagsAll filters OKRs keeping those where the Objective or any KeyResult has ALL matching tags.
func filterOKRsByTagsAll(okrs []OKR, tags []string) []OKR {
	if len(okrs) == 0 {
		return nil
	}
	result := make([]OKR, 0, len(okrs))
	for _, okr := range okrs {
		// Check if Objective has all matching tags
		if hasAllTags(okr.Objective.Tags, tags) {
			result = append(result, okr)
			continue
		}
		// Check if any KeyResult has all matching tags
		filteredKRs := filterSliceByTagsAll(okr.KeyResults, tags, func(kr KeyResult) []string { return kr.Tags })
		if len(filteredKRs) > 0 {
			okrCopy := OKR{
				Objective:  okr.Objective,
				KeyResults: filteredKRs,
			}
			result = append(result, okrCopy)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// filterOKRsByTags filters OKRs keeping those where the Objective or any KeyResult has matching tags.
func filterOKRsByTags(okrs []OKR, tags []string) []OKR {
	if len(okrs) == 0 {
		return nil
	}
	result := make([]OKR, 0, len(okrs))
	for _, okr := range okrs {
		// Check if Objective has matching tags
		if hasAnyTag(okr.Objective.Tags, tags) {
			result = append(result, okr)
			continue
		}
		// Check if any KeyResult has matching tags
		filteredKRs := filterSliceByTags(okr.KeyResults, tags, func(kr KeyResult) []string { return kr.Tags })
		if len(filteredKRs) > 0 {
			// Include OKR with only matching KeyResults
			okrCopy := OKR{
				Objective:  okr.Objective,
				KeyResults: filteredKRs,
			}
			result = append(result, okrCopy)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
