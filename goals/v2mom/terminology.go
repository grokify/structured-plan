package v2mom

// Terminology defines display labels for V2MOM components.
// This allows rendering in V2MOM, OKR, or hybrid terminology.
type Terminology struct {
	Methods          string // "Methods", "Objectives", or "Methods (Objectives)"
	MethodSingular   string // "Method", "Objective", or "Method (Objective)"
	Measures         string // "Measures", "Key Results", or "Measures (Key Results)"
	MeasureSingular  string // "Measure", "Key Result", or "Measure (Key Result)"
	Obstacles        string // "Obstacles" or "Risks"
	ObstacleSingular string // "Obstacle" or "Risk"
}

// GetTerminologyLabels returns terminology labels based on mode.
func GetTerminologyLabels(mode string) Terminology {
	switch mode {
	case TerminologyOKR:
		return Terminology{
			Methods:          "Objectives",
			MethodSingular:   "Objective",
			Measures:         "Key Results",
			MeasureSingular:  "Key Result",
			Obstacles:        "Risks",
			ObstacleSingular: "Risk",
		}
	case TerminologyHybrid:
		return Terminology{
			Methods:          "Methods (Objectives)",
			MethodSingular:   "Method (Objective)",
			Measures:         "Measures (Key Results)",
			MeasureSingular:  "Measure (Key Result)",
			Obstacles:        "Obstacles",
			ObstacleSingular: "Obstacle",
		}
	default: // TerminologyV2MOM
		return Terminology{
			Methods:          "Methods",
			MethodSingular:   "Method",
			Measures:         "Measures",
			MeasureSingular:  "Measure",
			Obstacles:        "Obstacles",
			ObstacleSingular: "Obstacle",
		}
	}
}

// V2MOMTerminology returns V2MOM terminology labels.
func V2MOMTerminology() Terminology {
	return GetTerminologyLabels(TerminologyV2MOM)
}

// OKRTerminology returns OKR terminology labels.
func OKRTerminology() Terminology {
	return GetTerminologyLabels(TerminologyOKR)
}

// HybridTerminology returns hybrid terminology labels.
func HybridTerminology() Terminology {
	return GetTerminologyLabels(TerminologyHybrid)
}
