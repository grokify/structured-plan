#!/bin/bash
# migrate_json_to_camelcase.sh
# Migrates user JSON files from snake_case to camelCase field names
#
# Usage: ./scripts/migrate_json_to_camelcase.sh input.json [output.json]
#
# If output.json is not specified, the input file is modified in place.
# Requires: jq

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 input.json [output.json]"
    echo ""
    echo "Converts snake_case JSON field names to camelCase."
    echo "If output.json is not specified, modifies input.json in place."
    exit 1
fi

INPUT="$1"
OUTPUT="${2:-$1}"

if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed."
    echo "Install with: brew install jq"
    exit 1
fi

if [ ! -f "$INPUT" ]; then
    echo "Error: File not found: $INPUT"
    exit 1
fi

# Convert snake_case keys to camelCase using jq
jq '
def snake_to_camel:
  gsub("_(?<a>[a-z])"; .a | ascii_upcase);

walk(
  if type == "object" then
    with_entries(.key |= snake_to_camel)
  else .
  end
)
' "$INPUT" > "${OUTPUT}.tmp" && mv "${OUTPUT}.tmp" "$OUTPUT"

echo "Converted: $INPUT -> $OUTPUT"
