package roadmap

import (
	"strings"
	"testing"
)

func TestToSwimlaneTable(t *testing.T) {
	r := createTestRoadmap()
	opts := DefaultTableOptions()

	table := r.ToSwimlaneTable(opts)

	// Check header contains phase names
	if !strings.Contains(table, "**Phase 1**") {
		t.Error("Missing Phase 1 header")
	}
	if !strings.Contains(table, "**Phase 2**") {
		t.Error("Missing Phase 2 header")
	}

	// Check swimlane labels
	if !strings.Contains(table, "**Features**") {
		t.Error("Missing Features swimlane")
	}
	if !strings.Contains(table, "**Infrastructure**") {
		t.Error("Missing Infrastructure swimlane")
	}

	// Check deliverable titles
	if !strings.Contains(table, "Auth") {
		t.Error("Missing Auth deliverable")
	}
	if !strings.Contains(table, "Dashboard") {
		t.Error("Missing Dashboard deliverable")
	}
}

func TestToSwimlaneTableEmpty(t *testing.T) {
	r := &Roadmap{}
	opts := DefaultTableOptions()

	table := r.ToSwimlaneTable(opts)

	if table != "" {
		t.Errorf("Expected empty table for empty roadmap, got: %s", table)
	}
}

func TestToSwimlaneTableWithStatus(t *testing.T) {
	r := createTestRoadmap()
	opts := DefaultTableOptions()
	opts.IncludeStatus = true

	table := r.ToSwimlaneTable(opts)

	// Check status icons
	if !strings.Contains(table, "✅") {
		t.Error("Missing completed status icon")
	}
	if !strings.Contains(table, "🔄") {
		t.Error("Missing in-progress status icon")
	}
}

func TestToPhaseTable(t *testing.T) {
	r := createTestRoadmap()
	opts := DefaultTableOptions()

	table := r.ToPhaseTable(opts)

	// Check header
	if !strings.Contains(table, "| Phase | Status | Deliverables |") {
		t.Error("Missing table header")
	}

	// Check phase names
	if !strings.Contains(table, "Foundation") {
		t.Error("Missing Foundation phase")
	}
	if !strings.Contains(table, "Core Features") {
		t.Error("Missing Core Features phase")
	}

	// Check status
	if !strings.Contains(table, "in_progress") {
		t.Error("Missing in_progress status")
	}
}

func TestToPhaseTableEmpty(t *testing.T) {
	r := &Roadmap{}
	opts := DefaultTableOptions()

	table := r.ToPhaseTable(opts)

	if table != "" {
		t.Errorf("Expected empty table for empty roadmap, got: %s", table)
	}
}

func TestSwimlaneLabel(t *testing.T) {
	tests := []struct {
		input    DeliverableType
		expected string
	}{
		{DeliverableFeature, "Features"},
		{DeliverableDocumentation, "Documentation"},
		{DeliverableInfrastructure, "Infrastructure"},
		{DeliverableIntegration, "Integrations"},
		{DeliverableMilestone, "Milestones"},
		{DeliverableRollout, "Rollout"},
		{DeliverableType("custom"), "Custom"},
		{DeliverableType(""), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			result := SwimlaneLabel(tt.input)
			if result != tt.expected {
				t.Errorf("SwimlaneLabel(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		input    DeliverableStatus
		expected string
	}{
		{DeliverableCompleted, "✅"},
		{DeliverableInProgress, "🔄"},
		{DeliverableBlocked, "🚫"},
		{DeliverableNotStarted, "⏳"},
		{DeliverableStatus(""), ""},
		{DeliverableStatus("unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			result := StatusIcon(tt.input)
			if result != tt.expected {
				t.Errorf("StatusIcon(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPhaseTargetStatusIcon(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"achieved", "✅"},
		{"in_progress", "🔄"},
		{"missed", "❌"},
		{"not_started", "⏳"},
		{"", ""},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := PhaseTargetStatusIcon(tt.input)
			if result != tt.expected {
				t.Errorf("PhaseTargetStatusIcon(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStatusLegend(t *testing.T) {
	legend := StatusLegend()

	// Check it contains expected icons
	if !strings.Contains(legend, "✅") {
		t.Error("Legend missing completed icon")
	}
	if !strings.Contains(legend, "🔄") {
		t.Error("Legend missing in-progress icon")
	}
	if !strings.Contains(legend, "⏳") {
		t.Error("Legend missing not-started icon")
	}
	if !strings.Contains(legend, "🚫") {
		t.Error("Legend missing blocked icon")
	}
	if !strings.Contains(legend, "❌") {
		t.Error("Legend missing missed icon")
	}
}

func TestMaxTitleLen(t *testing.T) {
	r := &Roadmap{
		Phases: []Phase{
			{
				ID:   "p1",
				Name: "Test Phase",
				Deliverables: []Deliverable{
					{
						ID:    "d1",
						Title: "This is a very long deliverable title that should be truncated",
						Type:  DeliverableFeature,
					},
				},
			},
		},
	}

	opts := TableOptions{
		IncludeStatus: false,
		MaxTitleLen:   20,
	}

	table := r.ToSwimlaneTable(opts)

	// Should be truncated with "..."
	if !strings.Contains(table, "...") {
		t.Error("Expected truncated title with ellipsis")
	}
	// Should not contain the full title
	if strings.Contains(table, "that should be truncated") {
		t.Error("Title was not truncated")
	}
}

func createTestRoadmap() *Roadmap {
	return &Roadmap{
		Phases: []Phase{
			{
				ID:     "phase-1",
				Name:   "Foundation",
				Type:   PhaseTypeGeneric,
				Status: PhaseStatusInProgress,
				Goals:  []string{"Establish base infrastructure"},
				Deliverables: []Deliverable{
					{
						ID:     "d1",
						Title:  "Auth",
						Type:   DeliverableFeature,
						Status: DeliverableCompleted,
					},
					{
						ID:     "d2",
						Title:  "CI/CD",
						Type:   DeliverableInfrastructure,
						Status: DeliverableInProgress,
					},
				},
				SuccessCriteria: []string{"All tests passing"},
			},
			{
				ID:     "phase-2",
				Name:   "Core Features",
				Type:   PhaseTypeGeneric,
				Status: PhaseStatusPlanned,
				Goals:  []string{"Build core functionality"},
				Deliverables: []Deliverable{
					{
						ID:     "d3",
						Title:  "Dashboard",
						Type:   DeliverableFeature,
						Status: DeliverableNotStarted,
					},
					{
						ID:     "d4",
						Title:  "Monitoring",
						Type:   DeliverableInfrastructure,
						Status: DeliverableNotStarted,
					},
				},
				SuccessCriteria: []string{"Feature complete"},
			},
		},
	}
}
