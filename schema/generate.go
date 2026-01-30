package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	"github.com/grokify/structured-requirements/requirements/prd"
)

// Generator creates JSON Schema files from Go types.
type Generator struct {
	// Reflector is the jsonschema reflector used for generation.
	Reflector *jsonschema.Reflector
}

// NewGenerator creates a new schema generator with default settings.
func NewGenerator() *Generator {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		ExpandedStruct:             false,
		RequiredFromJSONSchemaTags: true,
	}
	return &Generator{Reflector: r}
}

// GeneratePRDSchema generates JSON Schema for the PRD Document type.
func (g *Generator) GeneratePRDSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&prd.Document{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for prd.Document")
	}

	// Set schema metadata
	schema.ID = jsonschema.ID(PRDSchemaID)
	schema.Title = "Structured PRD"
	schema.Description = "Schema for structured Product Requirements Documents"

	return schema, nil
}

// GeneratePRDSchemaJSON generates JSON Schema for PRD and returns it as JSON bytes.
func (g *Generator) GeneratePRDSchemaJSON() ([]byte, error) {
	schema, err := g.GeneratePRDSchema()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(schema, "", "  ")
}

// WritePRDSchema generates and writes the PRD schema to a file.
func (g *Generator) WritePRDSchema(path string) error {
	data, err := g.GeneratePRDSchemaJSON()
	if err != nil {
		return fmt.Errorf("generating schema: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// GenerateAll generates all schema files to the specified directory.
func (g *Generator) GenerateAll(dir string) error {
	// Generate PRD schema
	prdPath := filepath.Join(dir, "prd.schema.json")
	if err := g.WritePRDSchema(prdPath); err != nil {
		return fmt.Errorf("generating PRD schema: %w", err)
	}

	// TODO: Add MRD and TRD schema generation when types are ready
	// mrdPath := filepath.Join(dir, "mrd.schema.json")
	// trdPath := filepath.Join(dir, "trd.schema.json")

	return nil
}
