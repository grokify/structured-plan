package prd

import (
	"testing"
)

func TestValidateTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		// Valid tags
		{"simple lowercase", "mvp", false},
		{"with hyphen", "phase-1", false},
		{"multiple hyphens", "backend-api-v2", false},
		{"leading digit", "2024-q1", false},
		{"all digits", "123", false},
		{"digit in middle", "v2-api", false},
		{"single char", "a", false},
		{"single digit", "1", false},

		// Invalid tags
		{"empty", "", true},
		{"uppercase", "MVP", true},
		{"mixed case", "myTag", true},
		{"leading hyphen", "-mvp", true},
		{"trailing hyphen", "mvp-", true},
		{"double hyphen", "my--tag", true},
		{"space", "my tag", true},
		{"underscore", "my_tag", true},
		{"special char", "my@tag", true},
		{"dot", "v1.0", true},
		{"colon", "phase:1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTag(%q) error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTags(t *testing.T) {
	tests := []struct {
		name      string
		tags      []string
		wantCount int
	}{
		{"all valid", []string{"mvp", "phase-1", "2024-q1"}, 0},
		{"one invalid", []string{"mvp", "My Tag", "phase-1"}, 1},
		{"multiple invalid", []string{"MVP", "my--tag", "phase-1"}, 2},
		{"empty slice", []string{}, 0},
		{"nil slice", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateTags(tt.tags)
			if len(errs) != tt.wantCount {
				t.Errorf("ValidateTags(%v) returned %d errors, want %d", tt.tags, len(errs), tt.wantCount)
			}
		})
	}
}
