package trd

import (
	"testing"
)

func TestFilterByTags_EmptyTags(t *testing.T) {
	doc := Document{
		Architecture: Architecture{
			Components: []Component{
				{ID: "c1", Name: "Component 1", Tags: []string{"data-layer"}},
			},
		},
	}
	filtered := doc.FilterByTags()
	if len(filtered.Architecture.Components) != 1 {
		t.Errorf("expected 1 component, got %d", len(filtered.Architecture.Components))
	}
}

func TestFilterByTags_SingleTag(t *testing.T) {
	doc := Document{
		Architecture: Architecture{
			Components: []Component{
				{ID: "c1", Name: "API Gateway", Tags: []string{"api"}},
				{ID: "c2", Name: "Data Service", Tags: []string{"data-layer"}},
				{ID: "c3", Name: "Auth Service", Tags: []string{"security"}},
			},
			ArchDecisions: []ArchDecision{
				{ID: "adr1", Title: "ADR 1", Tags: []string{"data-layer"}},
				{ID: "adr2", Title: "ADR 2", Tags: []string{"api"}},
			},
		},
		APISpecifications: []APISpec{
			{ID: "api1", Name: "Public API", Tags: []string{"api", "external"}},
			{ID: "api2", Name: "Internal API", Tags: []string{"internal"}},
		},
	}

	filtered := doc.FilterByTags("data-layer")

	if len(filtered.Architecture.Components) != 1 {
		t.Errorf("expected 1 component with data-layer tag, got %d", len(filtered.Architecture.Components))
	}
	if len(filtered.Architecture.ArchDecisions) != 1 {
		t.Errorf("expected 1 ADR with data-layer tag, got %d", len(filtered.Architecture.ArchDecisions))
	}
	if len(filtered.APISpecifications) != 0 {
		t.Errorf("expected 0 API specs, got %d", len(filtered.APISpecifications))
	}
}

func TestFilterByTags_MultipleTags_ORLogic(t *testing.T) {
	doc := Document{
		SecurityDesign: SecurityDesign{
			ThreatModel: []Threat{
				{ID: "t1", Name: "Threat 1", Tags: []string{"injection"}},
				{ID: "t2", Name: "Threat 2", Tags: []string{"xss"}},
				{ID: "t3", Name: "Threat 3", Tags: []string{"csrf"}},
			},
		},
	}

	// OR logic: should return threats with injection OR xss
	filtered := doc.FilterByTags("injection", "xss")

	if len(filtered.SecurityDesign.ThreatModel) != 2 {
		t.Errorf("expected 2 threats (injection OR xss), got %d", len(filtered.SecurityDesign.ThreatModel))
	}
}

func TestFilterByTags_DataModel(t *testing.T) {
	doc := Document{
		DataModel: &DataModel{
			Entities: []Entity{
				{ID: "e1", Name: "User", Tags: []string{"core"}},
				{ID: "e2", Name: "AuditLog", Tags: []string{"audit"}},
			},
			DataStores: []DataStore{
				{ID: "ds1", Name: "Primary DB", Tags: []string{"core"}},
				{ID: "ds2", Name: "Archive DB", Tags: []string{"archiving"}},
			},
		},
	}

	filtered := doc.FilterByTags("core")

	if filtered.DataModel == nil {
		t.Fatal("expected DataModel to not be nil")
	}
	if len(filtered.DataModel.Entities) != 1 {
		t.Errorf("expected 1 entity with core tag, got %d", len(filtered.DataModel.Entities))
	}
	if len(filtered.DataModel.DataStores) != 1 {
		t.Errorf("expected 1 data store with core tag, got %d", len(filtered.DataModel.DataStores))
	}
}

func TestFilterByTags_Integration(t *testing.T) {
	doc := Document{
		Integration: []Integration{
			{ID: "i1", Name: "Payment Gateway", Tags: []string{"external", "payments"}},
			{ID: "i2", Name: "Analytics", Tags: []string{"internal"}},
			{ID: "i3", Name: "Email Service", Tags: []string{"external", "notifications"}},
		},
	}

	filtered := doc.FilterByTags("external")

	if len(filtered.Integration) != 2 {
		t.Errorf("expected 2 integrations with external tag, got %d", len(filtered.Integration))
	}
}

func TestFilterByTags_Risks(t *testing.T) {
	doc := Document{
		Risks: []Risk{
			{ID: "r1", Description: "Risk 1", Tags: []string{"scalability"}},
			{ID: "r2", Description: "Risk 2", Tags: []string{"security"}},
			{ID: "r3", Description: "Risk 3", Tags: []string{"scalability", "performance"}},
		},
	}

	filtered := doc.FilterByTags("scalability")

	if len(filtered.Risks) != 2 {
		t.Errorf("expected 2 risks with scalability tag, got %d", len(filtered.Risks))
	}
}
