// Package common provides shared types used across PRD, MRD, and TRD documents.
package common

import (
	"fmt"
	"time"
)

// Person represents an individual contributor.
type Person struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
}

// Approver represents a person with approval authority.
type Approver struct {
	Person
	ApprovedAt *time.Time `json:"approvedAt,omitempty"`
	Approved   bool       `json:"approved"`
	Comments   string     `json:"comments,omitempty"`
}

// FormatPersonMarkdown formats a Person for markdown display.
// Output formats:
//   - "Name" (name only)
//   - "[Name](mailto:email)" (name + email)
//   - "Name (Role)" (name + role)
//   - "[Name](mailto:email) (Role)" (all fields)
func FormatPersonMarkdown(p Person) string {
	var result string

	if p.Email != "" {
		result = fmt.Sprintf("[%s](mailto:%s)", p.Name, p.Email)
	} else {
		result = p.Name
	}

	if p.Role != "" {
		result += fmt.Sprintf(" (%s)", p.Role)
	}

	return result
}

// FormatPeopleMarkdown formats a slice of Person for markdown display.
// Returns a comma-separated list of formatted persons.
func FormatPeopleMarkdown(people []Person) string {
	if len(people) == 0 {
		return ""
	}

	formatted := make([]string, len(people))
	for i, p := range people {
		formatted[i] = FormatPersonMarkdown(p)
	}

	result := formatted[0]
	for i := 1; i < len(formatted); i++ {
		result += ", " + formatted[i]
	}
	return result
}
