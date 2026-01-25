package schema

import (
	"encoding/json"
	"testing"
)

func TestPRDSchemaEmbedded(t *testing.T) {
	if len(PRDSchemaJSON) == 0 {
		t.Fatal("PRDSchemaJSON is empty")
	}

	var schema map[string]any
	if err := json.Unmarshal(PRDSchemaJSON, &schema); err != nil {
		t.Fatalf("PRDSchemaJSON is not valid JSON: %v", err)
	}

	expectedKeys := []string{"$schema", "$id", "title", "type", "properties", "$defs"}
	for _, key := range expectedKeys {
		if _, ok := schema[key]; !ok {
			t.Errorf("PRDSchemaJSON missing expected key: %s", key)
		}
	}
}

func TestPRDSchemaExtendedSections(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(PRDSchemaJSON, &schema); err != nil {
		t.Fatalf("PRDSchemaJSON is not valid JSON: %v", err)
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("properties is not an object")
	}

	extendedSections := []string{"problem", "market", "solution", "decisions", "reviews", "revision_history", "goals"}
	for _, section := range extendedSections {
		if _, ok := props[section]; !ok {
			t.Errorf("PRDSchemaJSON properties missing extended section: %s", section)
		}
	}
}

func TestPRDSchemaFunction(t *testing.T) {
	schema := PRDSchema()
	if schema == "" {
		t.Fatal("PRDSchema() returned empty string")
	}
	if schema != string(PRDSchemaJSON) {
		t.Error("PRDSchema() does not match PRDSchemaJSON")
	}
}

func TestPRDSchemaBytesFunction(t *testing.T) {
	schemaBytes := PRDSchemaBytes()
	if len(schemaBytes) == 0 {
		t.Fatal("PRDSchemaBytes() returned empty slice")
	}
	if string(schemaBytes) != string(PRDSchemaJSON) {
		t.Error("PRDSchemaBytes() does not match PRDSchemaJSON")
	}
}
