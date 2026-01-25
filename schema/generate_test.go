package schema

import (
	"encoding/json"
	"testing"
)

func TestGeneratePRDSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GeneratePRDSchema()
	if err != nil {
		t.Fatalf("GeneratePRDSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Structured PRD" {
		t.Errorf("expected title 'Structured PRD', got %q", schema.Title)
	}

	if string(schema.ID) != PRDSchemaID {
		t.Errorf("expected ID %q, got %q", PRDSchemaID, schema.ID)
	}
}

func TestGeneratePRDSchemaJSON(t *testing.T) {
	gen := NewGenerator()

	data, err := gen.GeneratePRDSchemaJSON()
	if err != nil {
		t.Fatalf("GeneratePRDSchemaJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("generated JSON is empty")
	}

	// Verify it's valid JSON
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("generated JSON is invalid: %v", err)
	}

	// Check for expected top-level keys (invopop/jsonschema uses $ref pattern)
	expectedKeys := []string{"$schema", "$id", "$ref", "$defs", "title"}
	for _, key := range expectedKeys {
		if _, ok := schema[key]; !ok {
			t.Errorf("generated schema missing key: %s", key)
		}
	}

	// Verify $ref points to Document definition
	ref, ok := schema["$ref"].(string)
	if !ok || ref != "#/$defs/Document" {
		t.Errorf("expected $ref to be '#/$defs/Document', got %v", schema["$ref"])
	}
}

func TestGeneratorReflectsExtendedSections(t *testing.T) {
	gen := NewGenerator()

	data, err := gen.GeneratePRDSchemaJSON()
	if err != nil {
		t.Fatalf("GeneratePRDSchemaJSON failed: %v", err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("generated JSON is invalid: %v", err)
	}

	// Navigate to $defs/Document/properties (invopop/jsonschema structure)
	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("$defs is not an object")
	}

	doc, ok := defs["Document"].(map[string]any)
	if !ok {
		t.Fatal("$defs/Document is not an object")
	}

	props, ok := doc["properties"].(map[string]any)
	if !ok {
		t.Fatal("$defs/Document/properties is not an object")
	}

	// Check that extended sections are reflected from Go types
	extendedSections := []string{"problem", "market", "solution", "decisions", "reviews", "revision_history", "goals"}
	for _, section := range extendedSections {
		if _, ok := props[section]; !ok {
			t.Errorf("generated schema missing extended section: %s", section)
		}
	}
}
