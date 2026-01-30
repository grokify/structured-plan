package trd

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
		TechnologyStack:  d.TechnologyStack,
		Deployment:       d.Deployment,
	}

	// Filter Architecture
	filtered.Architecture = Architecture{
		Overview:      d.Architecture.Overview,
		Principles:    d.Architecture.Principles,
		Patterns:      d.Architecture.Patterns,
		Components:    filterSliceByTags(d.Architecture.Components, tags, func(c Component) []string { return c.Tags }),
		Diagrams:      d.Architecture.Diagrams,
		DataFlows:     d.Architecture.DataFlows,
		ArchDecisions: filterSliceByTags(d.Architecture.ArchDecisions, tags, func(ad ArchDecision) []string { return ad.Tags }),
	}

	// Filter APISpecifications
	filtered.APISpecifications = filterSliceByTags(d.APISpecifications, tags, func(api APISpec) []string { return api.Tags })

	// Filter DataModel
	if d.DataModel != nil {
		dmCopy := *d.DataModel
		dmCopy.Entities = filterSliceByTags(d.DataModel.Entities, tags, func(e Entity) []string { return e.Tags })
		dmCopy.DataStores = filterSliceByTags(d.DataModel.DataStores, tags, func(ds DataStore) []string { return ds.Tags })
		filtered.DataModel = &dmCopy
	}

	// Filter SecurityDesign
	filtered.SecurityDesign = SecurityDesign{
		Overview:         d.SecurityDesign.Overview,
		AuthN:            d.SecurityDesign.AuthN,
		AuthZ:            d.SecurityDesign.AuthZ,
		Encryption:       d.SecurityDesign.Encryption,
		NetworkSecurity:  d.SecurityDesign.NetworkSecurity,
		Compliance:       d.SecurityDesign.Compliance,
		ThreatModel:      filterSliceByTags(d.SecurityDesign.ThreatModel, tags, func(t Threat) []string { return t.Tags }),
		SecurityControls: filterSliceByTags(d.SecurityDesign.SecurityControls, tags, func(sc SecurityControl) []string { return sc.Tags }),
	}

	// Filter Performance
	filtered.Performance = Performance{
		Overview:      d.Performance.Overview,
		Requirements:  filterSliceByTags(d.Performance.Requirements, tags, func(pr PerfRequirement) []string { return pr.Tags }),
		Benchmarks:    d.Performance.Benchmarks,
		Optimizations: d.Performance.Optimizations,
	}

	// Keep Scalability as-is (no tagged entities)
	filtered.Scalability = d.Scalability

	// Filter Integration
	filtered.Integration = filterSliceByTags(d.Integration, tags, func(i Integration) []string { return i.Tags })

	// Keep Development and Testing as-is (no tagged entities)
	filtered.Development = d.Development
	filtered.Testing = d.Testing

	// Filter Risks
	filtered.Risks = filterSliceByTags(d.Risks, tags, func(r Risk) []string { return r.Tags })

	// Filter Constraints
	filtered.Constraints = filterSliceByTags(d.Constraints, tags, func(c Constraint) []string { return c.Tags })

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
