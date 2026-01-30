package mrd

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
		MarketOverview:   d.MarketOverview,
		Positioning:      d.Positioning,
	}

	// Filter TargetMarket
	filtered.TargetMarket = TargetMarket{
		PrimarySegments:   filterSliceByTags(d.TargetMarket.PrimarySegments, tags, func(s MarketSegment) []string { return s.Tags }),
		SecondarySegments: filterSliceByTags(d.TargetMarket.SecondarySegments, tags, func(s MarketSegment) []string { return s.Tags }),
		BuyerPersonas:     filterSliceByTags(d.TargetMarket.BuyerPersonas, tags, func(bp BuyerPersona) []string { return bp.Tags }),
		Verticals:         d.TargetMarket.Verticals,
		GeographicFocus:   d.TargetMarket.GeographicFocus,
		CompanySize:       d.TargetMarket.CompanySize,
	}

	// Filter CompetitiveLandscape
	filtered.CompetitiveLandscape = CompetitiveLandscape{
		Overview:        d.CompetitiveLandscape.Overview,
		Competitors:     filterSliceByTags(d.CompetitiveLandscape.Competitors, tags, func(c Competitor) []string { return c.Tags }),
		MarketPosition:  d.CompetitiveLandscape.MarketPosition,
		Differentiators: d.CompetitiveLandscape.Differentiators,
		CompetitiveGaps: d.CompetitiveLandscape.CompetitiveGaps,
	}

	// Filter MarketRequirements
	filtered.MarketRequirements = filterSliceByTags(d.MarketRequirements, tags, func(mr MarketRequirement) []string { return mr.Tags })

	// Filter GoToMarket milestones
	if d.GoToMarket != nil {
		gtmCopy := *d.GoToMarket
		gtmCopy.Milestones = filterSliceByTags(d.GoToMarket.Milestones, tags, func(m Milestone) []string { return m.Tags })
		filtered.GoToMarket = &gtmCopy
	}

	// Filter SuccessMetrics
	filtered.SuccessMetrics = filterSliceByTags(d.SuccessMetrics, tags, func(sm SuccessMetric) []string { return sm.Tags })

	// Filter Risks
	filtered.Risks = filterSliceByTags(d.Risks, tags, func(r Risk) []string { return r.Tags })

	// Keep untagged sections as-is
	filtered.Assumptions = d.Assumptions
	filtered.Glossary = d.Glossary
	filtered.CustomSections = d.CustomSections

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
