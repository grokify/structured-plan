// Package schema provides embedded JSON Schema files for structured requirements documents.
// These schemas can be used for validation and documentation purposes.
package schema

import (
	_ "embed"
)

// PRD Schema

//go:embed prd.schema.json
var PRDSchemaJSON []byte

// PRDSchema returns the PRD JSON Schema as a string.
func PRDSchema() string {
	return string(PRDSchemaJSON)
}

// PRDSchemaBytes returns the PRD JSON Schema as a byte slice.
func PRDSchemaBytes() []byte {
	return PRDSchemaJSON
}

// SchemaID constants for referencing schemas.
const (
	// PRDSchemaID is the canonical ID for the PRD schema.
	PRDSchemaID = "https://github.com/grokify/structured-requirements/schema/prd.schema.json"

	// MRDSchemaID is the canonical ID for the MRD schema (placeholder).
	MRDSchemaID = "https://github.com/grokify/structured-requirements/schema/mrd.schema.json"

	// TRDSchemaID is the canonical ID for the TRD schema (placeholder).
	TRDSchemaID = "https://github.com/grokify/structured-requirements/schema/trd.schema.json"
)

// TODO: Add MRD and TRD schemas when created.
// When mrd.schema.json is added:
//
//	//go:embed mrd.schema.json
//	var MRDSchemaJSON []byte
//
//	func MRDSchema() string { return string(MRDSchemaJSON) }
//
// When trd.schema.json is added:
//
//	//go:embed trd.schema.json
//	var TRDSchemaJSON []byte
//
//	func TRDSchema() string { return string(TRDSchemaJSON) }
