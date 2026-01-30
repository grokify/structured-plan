package trd

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestDocumentParsing tests JSON unmarshaling of TRD documents.
func TestDocumentParsing(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(t *testing.T, doc Document)
	}{
		{
			name: "minimal valid document",
			json: `{
				"metadata": {
					"id": "trd-001",
					"title": "Test TRD",
					"version": "1.0.0",
					"status": "draft",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Test Author"}]
				},
				"executiveSummary": {
					"purpose": "Define technical architecture",
					"scope": "Backend services",
					"technicalApproach": "Microservices"
				},
				"architecture": {
					"overview": "Microservices architecture",
					"components": [{"id": "c1", "name": "API Gateway", "description": "Entry point"}]
				},
				"technologyStack": {},
				"securityDesign": {
					"overview": "Defense in depth"
				},
				"performance": {
					"requirements": [{"id": "p1", "name": "Latency", "metric": "p99 latency", "target": "<100ms"}]
				},
				"deployment": {
					"overview": "Kubernetes deployment",
					"environments": [{"name": "Production", "purpose": "Live traffic"}]
				}
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if doc.Metadata.ID != "trd-001" {
					t.Errorf("expected ID 'trd-001', got %q", doc.Metadata.ID)
				}
				if doc.Metadata.Title != "Test TRD" {
					t.Errorf("expected Title 'Test TRD', got %q", doc.Metadata.Title)
				}
				if doc.ExecutiveSummary.Purpose != "Define technical architecture" {
					t.Errorf("expected Purpose, got %q", doc.ExecutiveSummary.Purpose)
				}
				if len(doc.Architecture.Components) != 1 {
					t.Errorf("expected 1 component, got %d", len(doc.Architecture.Components))
				}
				if len(doc.Performance.Requirements) != 1 {
					t.Errorf("expected 1 performance requirement, got %d", len(doc.Performance.Requirements))
				}
			},
		},
		{
			name: "document with API specifications",
			json: `{
				"metadata": {
					"id": "trd-002",
					"title": "Test",
					"version": "1.0.0",
					"status": "draft",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author"}]
				},
				"executiveSummary": {"purpose": "P", "scope": "S", "technicalApproach": "T"},
				"architecture": {
					"overview": "Overview",
					"components": [{"id": "c1", "name": "Service", "description": "D"}]
				},
				"technologyStack": {},
				"apiSpecifications": [
					{
						"id": "api1",
						"name": "User API",
						"type": "REST",
						"version": "v1",
						"endpoints": [
							{"method": "GET", "path": "/users", "description": "List users"},
							{"method": "POST", "path": "/users", "description": "Create user"}
						]
					}
				],
				"securityDesign": {"overview": "O"},
				"performance": {"requirements": [{"id": "p1", "name": "N", "metric": "M", "target": "T"}]},
				"deployment": {"overview": "O", "environments": [{"name": "prod"}]}
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if len(doc.APISpecifications) != 1 {
					t.Errorf("expected 1 API spec, got %d", len(doc.APISpecifications))
				}
				api := doc.APISpecifications[0]
				if api.Type != "REST" {
					t.Errorf("expected API type 'REST', got %q", api.Type)
				}
				if len(api.Endpoints) != 2 {
					t.Errorf("expected 2 endpoints, got %d", len(api.Endpoints))
				}
			},
		},
		{
			name: "document with security details",
			json: `{
				"metadata": {
					"id": "trd-003",
					"title": "Test",
					"version": "1.0.0",
					"status": "draft",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author"}]
				},
				"executiveSummary": {"purpose": "P", "scope": "S", "technicalApproach": "T"},
				"architecture": {"overview": "O", "components": [{"id": "c1", "name": "N", "description": "D"}]},
				"technologyStack": {},
				"securityDesign": {
					"overview": "Zero trust architecture",
					"authentication": {
						"method": "OAuth2",
						"provider": "Auth0",
						"mfa": true
					},
					"authorization": {
						"model": "RBAC",
						"roles": ["admin", "user", "viewer"]
					},
					"encryption": {
						"atRest": "AES-256",
						"inTransit": "TLS 1.3"
					},
					"compliance": ["SOC2", "GDPR"]
				},
				"performance": {"requirements": [{"id": "p1", "name": "N", "metric": "M", "target": "T"}]},
				"deployment": {"overview": "O", "environments": [{"name": "prod"}]}
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if doc.SecurityDesign.AuthN == nil {
					t.Fatal("expected authentication to be present")
				}
				if doc.SecurityDesign.AuthN.Method != "OAuth2" {
					t.Errorf("expected auth method 'OAuth2', got %q", doc.SecurityDesign.AuthN.Method)
				}
				if !doc.SecurityDesign.AuthN.MFA {
					t.Error("expected MFA to be true")
				}
				if doc.SecurityDesign.AuthZ == nil {
					t.Fatal("expected authorization to be present")
				}
				if doc.SecurityDesign.AuthZ.Model != "RBAC" {
					t.Errorf("expected auth model 'RBAC', got %q", doc.SecurityDesign.AuthZ.Model)
				}
				if len(doc.SecurityDesign.Compliance) != 2 {
					t.Errorf("expected 2 compliance standards, got %d", len(doc.SecurityDesign.Compliance))
				}
			},
		},
		{
			name: "document with integrations",
			json: `{
				"metadata": {
					"id": "trd-004",
					"title": "Test",
					"version": "1.0.0",
					"status": "draft",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"authors": [{"name": "Author"}]
				},
				"executiveSummary": {"purpose": "P", "scope": "S", "technicalApproach": "T"},
				"architecture": {"overview": "O", "components": [{"id": "c1", "name": "N", "description": "D"}]},
				"technologyStack": {},
				"securityDesign": {"overview": "O"},
				"performance": {"requirements": [{"id": "p1", "name": "N", "metric": "M", "target": "T"}]},
				"deployment": {"overview": "O", "environments": [{"name": "prod"}]},
				"integrations": [
					{
						"id": "int1",
						"name": "Payment Gateway",
						"type": "API",
						"direction": "Outbound",
						"protocol": "HTTPS",
						"authMethod": "API Key"
					}
				]
			}`,
			wantErr: false,
			check: func(t *testing.T, doc Document) {
				if len(doc.Integration) != 1 {
					t.Errorf("expected 1 integration, got %d", len(doc.Integration))
				}
				if doc.Integration[0].Direction != "Outbound" {
					t.Errorf("expected direction 'Outbound', got %q", doc.Integration[0].Direction)
				}
			},
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc Document
			err := json.Unmarshal([]byte(tt.json), &doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, doc)
			}
		})
	}
}

// TestDocumentMarshaling tests JSON marshaling of TRD documents.
func TestDocumentMarshaling(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:        "trd-test",
			Title:     "Test Document",
			Version:   "1.0.0",
			Status:    StatusDraft,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Authors:   []Person{{Name: "Test Author"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			Purpose:           "Define architecture",
			Scope:             "Backend services",
			TechnicalApproach: "Microservices",
		},
		Architecture: Architecture{
			Overview:   "Microservices architecture",
			Components: []Component{{ID: "c1", Name: "API", Description: "Gateway"}},
		},
		SecurityDesign: SecurityDesign{
			Overview: "Defense in depth",
		},
		Performance: Performance{
			Requirements: []PerfRequirement{{ID: "p1", Name: "Latency", Metric: "p99", Target: "<100ms"}},
		},
		Deployment: Deployment{
			Overview:     "Kubernetes",
			Environments: []Environment{{Name: "Production"}},
		},
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var doc2 Document
	if err := json.Unmarshal(data, &doc2); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if doc2.Metadata.ID != doc.Metadata.ID {
		t.Errorf("round-trip ID mismatch: got %q, want %q", doc2.Metadata.ID, doc.Metadata.ID)
	}
	if doc2.Architecture.Overview != doc.Architecture.Overview {
		t.Errorf("round-trip Architecture.Overview mismatch: got %q, want %q", doc2.Architecture.Overview, doc.Architecture.Overview)
	}
}

// TestStatusConstants verifies status constant values.
func TestStatusConstants(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{StatusDraft, "draft"},
		{StatusInReview, "in_review"},
		{StatusApproved, "approved"},
		{StatusDeprecated, "deprecated"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("Status = %q, want %q", tt.status, tt.want)
			}
		})
	}
}

// TestMarkdownGeneration tests the ToMarkdown method.
func TestMarkdownGeneration(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:        "trd-test",
			Title:     "Test Technical Requirements",
			Version:   "1.0.0",
			Status:    StatusDraft,
			CreatedAt: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			Authors:   []Person{{Name: "John Engineer", Role: "Tech Lead"}},
			Tags:      []string{"architecture", "microservices"},
		},
		ExecutiveSummary: ExecutiveSummary{
			Purpose:           "Define the technical architecture for the platform.",
			Scope:             "Backend services and infrastructure.",
			TechnicalApproach: "Microservices deployed on Kubernetes.",
			KeyDecisions:      []string{"Use Go for services", "PostgreSQL for data"},
		},
		Architecture: Architecture{
			Overview:   "Microservices architecture with API gateway.",
			Principles: []string{"Defense in depth", "Least privilege"},
			Components: []Component{
				{ID: "c1", Name: "API Gateway", Description: "Entry point", Type: "Service"},
				{ID: "c2", Name: "User Service", Description: "User management", Type: "Service"},
			},
		},
		TechnologyStack: TechnologyStack{
			Languages: []Technology{
				{Name: "Go", Version: "1.22", Purpose: "Backend services"},
			},
			Databases: []Technology{
				{Name: "PostgreSQL", Version: "16", Purpose: "Primary data store"},
			},
		},
		APISpecifications: []APISpec{
			{ID: "api1", Name: "User API", Type: "REST", Version: "v1"},
		},
		SecurityDesign: SecurityDesign{
			Overview:   "Zero trust architecture.",
			Compliance: []string{"SOC2", "GDPR"},
		},
		Performance: Performance{
			Requirements: []PerfRequirement{
				{ID: "p1", Name: "API Latency", Metric: "p99 latency", Target: "<100ms"},
			},
		},
		Deployment: Deployment{
			Overview: "Kubernetes deployment.",
			Environments: []Environment{
				{Name: "Development", Purpose: "Developer testing"},
				{Name: "Production", Purpose: "Live traffic"},
			},
			Strategy: "Blue-green",
		},
		Integration: []Integration{
			{ID: "int1", Name: "Payment Gateway", Type: "API", Direction: "Outbound"},
		},
		Glossary: []GlossaryTerm{
			{Term: "mTLS", Definition: "Mutual TLS"},
		},
	}

	t.Run("with frontmatter", func(t *testing.T) {
		opts := MarkdownOptions{
			IncludeFrontmatter: true,
			Margin:             "2cm",
			MainFont:           "Helvetica",
		}
		md := doc.ToMarkdown(opts)

		// Check frontmatter
		if !strings.HasPrefix(md, "---\n") {
			t.Error("expected markdown to start with YAML frontmatter")
		}
		if !strings.Contains(md, `title: "Test Technical Requirements"`) {
			t.Error("expected frontmatter to contain title")
		}

		// Check content sections
		if !strings.Contains(md, "# Test Technical Requirements") {
			t.Error("expected markdown to contain document title as H1")
		}
		if !strings.Contains(md, "Executive Summary") {
			t.Error("expected markdown to contain Executive Summary section")
		}
		if !strings.Contains(md, "Architecture") {
			t.Error("expected markdown to contain Architecture section")
		}
		if !strings.Contains(md, "API Gateway") {
			t.Error("expected markdown to contain component name")
		}
		if !strings.Contains(md, "Technology Stack") {
			t.Error("expected markdown to contain Technology Stack section")
		}
		if !strings.Contains(md, "API Specifications") {
			t.Error("expected markdown to contain API Specifications section")
		}
		if !strings.Contains(md, "Security") {
			t.Error("expected markdown to contain Security section")
		}
		if !strings.Contains(md, "Performance") {
			t.Error("expected markdown to contain Performance section")
		}
		if !strings.Contains(md, "Deployment") {
			t.Error("expected markdown to contain Deployment section")
		}
		if !strings.Contains(md, "Integration") {
			t.Error("expected markdown to contain Integration section")
		}
		if !strings.Contains(md, "Glossary") {
			t.Error("expected markdown to contain Glossary section")
		}
	})

	t.Run("without frontmatter", func(t *testing.T) {
		opts := MarkdownOptions{
			IncludeFrontmatter: false,
		}
		md := doc.ToMarkdown(opts)

		if strings.HasPrefix(md, "---\n") {
			t.Error("expected markdown to NOT start with YAML frontmatter")
		}
		if !strings.HasPrefix(md, "# Test Technical Requirements") {
			t.Error("expected markdown to start with document title")
		}
	})
}

// TestValidation tests document validation logic.
func TestValidation(t *testing.T) {
	tests := []struct {
		name       string
		doc        Document
		wantErrors []string
	}{
		{
			name: "valid document",
			doc: Document{
				Metadata: Metadata{
					ID:      "trd-001",
					Title:   "Valid TRD",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{
					Purpose: "Purpose",
					Scope:   "Scope",
				},
				Architecture: Architecture{
					Overview:   "Overview",
					Components: []Component{{ID: "c1", Name: "C", Description: "D"}},
				},
				SecurityDesign: SecurityDesign{
					Overview: "Security overview",
				},
				Performance: Performance{
					Requirements: []PerfRequirement{{ID: "p1", Name: "N", Metric: "M", Target: "T"}},
				},
				Deployment: Deployment{
					Environments: []Environment{{Name: "Production"}},
				},
			},
			wantErrors: nil,
		},
		{
			name: "missing metadata.id",
			doc: Document{
				Metadata: Metadata{
					Title:   "Test",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{Purpose: "P", Scope: "S"},
				Architecture:     Architecture{Overview: "O", Components: []Component{{ID: "c1"}}},
				SecurityDesign:   SecurityDesign{Overview: "O"},
				Performance:      Performance{Requirements: []PerfRequirement{{ID: "p1"}}},
				Deployment:       Deployment{Environments: []Environment{{Name: "prod"}}},
			},
			wantErrors: []string{"metadata.id is required"},
		},
		{
			name: "missing architecture overview",
			doc: Document{
				Metadata: Metadata{
					ID:      "trd-001",
					Title:   "Test",
					Version: "1.0.0",
					Authors: []Person{{Name: "Author"}},
				},
				ExecutiveSummary: ExecutiveSummary{Purpose: "P", Scope: "S"},
				Architecture:     Architecture{Components: []Component{{ID: "c1"}}},
				SecurityDesign:   SecurityDesign{Overview: "O"},
				Performance:      Performance{Requirements: []PerfRequirement{{ID: "p1"}}},
				Deployment:       Deployment{Environments: []Environment{{Name: "prod"}}},
			},
			wantErrors: []string{"architecture.overview is required"},
		},
		{
			name: "multiple errors",
			doc: Document{
				Metadata:     Metadata{Version: "1.0.0"},
				Architecture: Architecture{},
				Performance:  Performance{},
				Deployment:   Deployment{},
			},
			wantErrors: []string{
				"metadata.id is required",
				"metadata.title is required",
				"metadata.authors is required",
				"executive_summary.purpose is required",
				"executive_summary.scope is required",
				"architecture.overview is required",
				"architecture.components is required",
				"security_design.overview is required",
				"performance.requirements is required",
				"deployment.environments is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validateDocument(tt.doc)

			if len(tt.wantErrors) == 0 {
				if len(errors) > 0 {
					t.Errorf("expected no errors, got %v", errors)
				}
				return
			}

			if len(errors) != len(tt.wantErrors) {
				t.Errorf("expected %d errors, got %d: %v", len(tt.wantErrors), len(errors), errors)
				return
			}

			for i, wantErr := range tt.wantErrors {
				if errors[i] != wantErr {
					t.Errorf("error[%d] = %q, want %q", i, errors[i], wantErr)
				}
			}
		})
	}
}

// validateDocument mirrors CLI validation logic.
func validateDocument(doc Document) []string {
	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required")
	}
	if doc.ExecutiveSummary.Purpose == "" {
		errors = append(errors, "executive_summary.purpose is required")
	}
	if doc.ExecutiveSummary.Scope == "" {
		errors = append(errors, "executive_summary.scope is required")
	}
	if doc.Architecture.Overview == "" {
		errors = append(errors, "architecture.overview is required")
	}
	if len(doc.Architecture.Components) == 0 {
		errors = append(errors, "architecture.components is required")
	}
	if doc.SecurityDesign.Overview == "" {
		errors = append(errors, "security_design.overview is required")
	}
	if len(doc.Performance.Requirements) == 0 {
		errors = append(errors, "performance.requirements is required")
	}
	if len(doc.Deployment.Environments) == 0 {
		errors = append(errors, "deployment.environments is required")
	}

	return errors
}
