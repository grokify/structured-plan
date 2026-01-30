package prd

// Requirements contains both functional and non-functional requirements.
type Requirements struct {
	Functional    []FunctionalRequirement    `json:"functional"`
	NonFunctional []NonFunctionalRequirement `json:"nonFunctional"`
}

// FunctionalRequirement represents a functional requirement.
type FunctionalRequirement struct {
	ID                 string                `json:"id"` // e.g., FR-001
	Title              string                `json:"title"`
	Description        string                `json:"description"`
	Category           string                `json:"category"` // Feature category
	Priority           MoSCoW                `json:"priority"`
	UserStoryIDs       []string              `json:"userStoryIds"` // Related user stories
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptanceCriteria"`
	PhaseID            string                `json:"phaseId"` // Target roadmap phase
	Dependencies       []string              `json:"dependencies,omitempty"`
	Assumptions        []string              `json:"assumptions,omitempty"`
	Tags               []string              `json:"tags,omitempty"` // For filtering by topic/domain
	Notes              string                `json:"notes,omitempty"`

	// AppendixRefs references appendices with additional details for this requirement.
	AppendixRefs []string `json:"appendixRefs,omitempty"`
}

// NFRCategory represents categories of non-functional requirements.
type NFRCategory string

const (
	NFRPerformance      NFRCategory = "performance"
	NFRScalability      NFRCategory = "scalability"
	NFRReliability      NFRCategory = "reliability"
	NFRAvailability     NFRCategory = "availability"
	NFRSecurity         NFRCategory = "security"
	NFRMultiTenancy     NFRCategory = "multi_tenancy"
	NFRObservability    NFRCategory = "observability"
	NFRMaintainability  NFRCategory = "maintainability"
	NFRUsability        NFRCategory = "usability"
	NFRCompatibility    NFRCategory = "compatibility"
	NFRCompliance       NFRCategory = "compliance"
	NFRDisasterRecovery NFRCategory = "disaster_recovery"
	NFRCostEfficiency   NFRCategory = "cost_efficiency"
	NFRPortability      NFRCategory = "portability"
	NFRTestability      NFRCategory = "testability"
	NFRExtensibility    NFRCategory = "extensibility"
	NFRInteroperability NFRCategory = "interoperability"
	NFRLocalization     NFRCategory = "localization"
)

// NonFunctionalRequirement represents a non-functional requirement.
type NonFunctionalRequirement struct {
	ID                string      `json:"id"` // e.g., NFR-001
	Category          NFRCategory `json:"category"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Metric            string      `json:"metric"` // What is measured
	Target            string      `json:"target"` // Target value (e.g., "P95 < 200ms")
	MeasurementMethod string      `json:"measurementMethod,omitempty"`
	Priority          MoSCoW      `json:"priority"`
	PhaseID           string      `json:"phaseId"`
	CurrentBaseline   string      `json:"currentBaseline,omitempty"`
	Notes             string      `json:"notes,omitempty"`

	// SLO-specific fields (for observability/reliability)
	SLO *SLOSpec `json:"slo,omitempty"`

	// Multi-tenancy specific fields
	MultiTenancy *MultiTenancySpec `json:"multiTenancy,omitempty"`

	// Security specific fields
	Security *SecuritySpec `json:"security,omitempty"`

	Tags []string `json:"tags,omitempty"` // For filtering by topic/domain

	// AppendixRefs references appendices with additional details for this requirement.
	AppendixRefs []string `json:"appendixRefs,omitempty"`
}

// SLOSpec defines Service Level Objective specifications.
type SLOSpec struct {
	SLI            string `json:"sli"`       // Service Level Indicator
	SLOTarget      string `json:"sloTarget"` // e.g., "99.9%"
	Window         string `json:"window"`    // e.g., "30 days rolling"
	ErrorBudget    string `json:"errorBudget,omitempty"`
	Consequences   string `json:"consequences,omitempty"` // What happens on breach
	AlertThreshold string `json:"alertThreshold,omitempty"`
}

// MultiTenancySpec defines multi-tenancy requirements.
type MultiTenancySpec struct {
	IsolationModel          IsolationModel   `json:"isolationModel"`
	DataSegregation         DataSegregation  `json:"dataSegregation"`
	EncryptionModel         EncryptionModel  `json:"encryptionModel,omitempty"`
	NetworkIsolation        NetworkIsolation `json:"networkIsolation,omitempty"`
	NoisyNeighborProtection string           `json:"noisyNeighborProtection,omitempty"`
}

// IsolationModel represents tenant isolation strategies.
type IsolationModel string

const (
	IsolationPool   IsolationModel = "pool"   // Shared resources
	IsolationSilo   IsolationModel = "silo"   // Dedicated resources
	IsolationBridge IsolationModel = "bridge" // Hybrid approach
)

// DataSegregation represents database isolation levels.
type DataSegregation string

const (
	DataSharedSchema      DataSegregation = "shared_schema"       // Single schema with tenant ID
	DataSchemaPerTenant   DataSegregation = "schema_per_tenant"   // Separate schema per tenant
	DataDatabasePerTenant DataSegregation = "database_per_tenant" // Separate database per tenant
)

// EncryptionModel represents cryptographic isolation levels.
type EncryptionModel string

const (
	EncryptionSharedKeys     EncryptionModel = "shared_keys"
	EncryptionTenantSpecific EncryptionModel = "tenant_specific_keys"
	EncryptionBYOK           EncryptionModel = "byok" // Bring Your Own Key
)

// NetworkIsolation represents network-level isolation.
type NetworkIsolation string

const (
	NetworkShared             NetworkIsolation = "shared"
	NetworkVPCPerTenant       NetworkIsolation = "vpc_per_tenant"
	NetworkNamespaceIsolation NetworkIsolation = "namespace_isolation"
)

// SecuritySpec defines security-specific requirements.
type SecuritySpec struct {
	AuthenticationMethods  []string `json:"authenticationMethods,omitempty"` // OAuth2, SAML, MFA
	AuthorizationModel     string   `json:"authorizationModel,omitempty"`    // RBAC, ABAC
	EncryptionAtRest       bool     `json:"encryptionAtRest,omitempty"`
	EncryptionInTransit    bool     `json:"encryptionInTransit,omitempty"`
	ComplianceStandards    []string `json:"complianceStandards,omitempty"` // SOC2, GDPR, HIPAA
	VulnerabilityScanning  bool     `json:"vulnerabilityScanning,omitempty"`
	PenetrationTesting     bool     `json:"penetrationTesting,omitempty"`
	SecurityAuditFrequency string   `json:"securityAuditFrequency,omitempty"`
}

// ObservabilitySpec defines observability requirements.
type ObservabilitySpec struct {
	Logging    LoggingSpec  `json:"logging"`
	Metrics    MetricsSpec  `json:"metrics"`
	Tracing    TracingSpec  `json:"tracing"`
	Alerting   AlertingSpec `json:"alerting"`
	Dashboards []string     `json:"dashboards,omitempty"`
}

// LoggingSpec defines logging requirements.
type LoggingSpec struct {
	Format                string   `json:"format"`          // JSON, structured
	RetentionPeriod       string   `json:"retentionPeriod"` // e.g., "90 days"
	LogLevels             []string `json:"logLevels"`
	CorrelationID         bool     `json:"correlationId"` // Include trace/correlation IDs
	SensitiveDataHandling string   `json:"sensitiveDataHandling,omitempty"`
}

// MetricsSpec defines metrics requirements.
type MetricsSpec struct {
	Format             string   `json:"format"`             // Prometheus, OpenTelemetry
	CollectionInterval string   `json:"collectionInterval"` // e.g., "15s"
	RetentionPeriod    string   `json:"retentionPeriod"`
	CustomMetrics      []string `json:"customMetrics,omitempty"`
}

// TracingSpec defines distributed tracing requirements.
type TracingSpec struct {
	Enabled           bool   `json:"enabled"`
	SamplingRate      string `json:"samplingRate"`      // e.g., "100%", "10%"
	PropagationFormat string `json:"propagationFormat"` // W3C, B3
	ExportFormat      string `json:"exportFormat"`      // OTLP, Jaeger, Zipkin
}

// AlertingSpec defines alerting requirements.
type AlertingSpec struct {
	Channels          []string `json:"channels"` // PagerDuty, Slack, Email
	EscalationPolicy  string   `json:"escalationPolicy,omitempty"`
	OnCallIntegration bool     `json:"onCallIntegration,omitempty"`
}

// ReliabilitySpec defines reliability requirements.
type ReliabilitySpec struct {
	TargetUptime         string `json:"targetUptime"`   // e.g., "99.9%"
	MTBF                 string `json:"mtbf,omitempty"` // Mean Time Between Failures
	MTTR                 string `json:"mttr,omitempty"` // Mean Time To Recovery
	RTO                  string `json:"rto,omitempty"`  // Recovery Time Objective
	RPO                  string `json:"rpo,omitempty"`  // Recovery Point Objective
	FailoverStrategy     string `json:"failoverStrategy,omitempty"`
	BackupFrequency      string `json:"backupFrequency,omitempty"`
	DisasterRecoveryPlan bool   `json:"disasterRecoveryPlan,omitempty"`
}
