package prd

// SecurityModel documents security architecture and threat model.
// This section is required for all PRDs.
type SecurityModel struct {
	// Overview provides a high-level summary of the security approach.
	Overview string `json:"overview"`

	// ThreatModel identifies assets, threat actors, and threats.
	ThreatModel ThreatModel `json:"threatModel"`

	// AccessControl defines access control strategy.
	AccessControl AccessControl `json:"accessControl"`

	// Encryption specifies encryption requirements.
	Encryption EncryptionRequirements `json:"encryption"`

	// AuditLogging defines audit logging requirements.
	AuditLogging AuditLogging `json:"auditLogging"`

	// ComplianceControls maps to compliance frameworks.
	// Key is framework name (e.g., "SOC2", "GDPR"), value is list of controls.
	ComplianceControls map[string][]string `json:"complianceControls,omitempty"`

	// DataClassification defines data sensitivity levels.
	DataClassification []DataClassification `json:"dataClassification,omitempty"`

	// AppendixRefs references appendices with additional security details.
	AppendixRefs []string `json:"appendixRefs,omitempty"`
}

// ThreatModel identifies security threats and mitigations.
type ThreatModel struct {
	// Assets are the valuable resources to protect.
	Assets []string `json:"assets"`

	// ThreatActors are potential attackers.
	ThreatActors []string `json:"threatActors"`

	// KeyThreats lists major threats with mitigations.
	KeyThreats []SecurityThreat `json:"keyThreats"`

	// TrustBoundaries identifies trust boundaries in the system.
	TrustBoundaries []string `json:"trustBoundaries,omitempty"`
}

// SecurityThreat represents a security threat.
type SecurityThreat struct {
	// ID is the unique identifier for this threat.
	ID string `json:"id,omitempty"`

	// Threat description.
	Threat string `json:"threat"`

	// Category is the threat category (e.g., "STRIDE" categories).
	Category string `json:"category,omitempty"`

	// Mitigation strategy.
	Mitigation string `json:"mitigation"`

	// Severity level (critical, high, medium, low).
	Severity string `json:"severity,omitempty"`

	// Status of mitigation (planned, implemented, verified).
	Status string `json:"status,omitempty"`

	// RelatedIDs links to related requirements or risks.
	RelatedIDs []string `json:"relatedIds,omitempty"`
}

// AccessControl defines access control strategy.
type AccessControl struct {
	// Model is the access control model (RBAC, ABAC, ReBAC, etc.).
	Model string `json:"model"`

	// Description provides details on the access control approach.
	Description string `json:"description,omitempty"`

	// Layers describes access control at different layers.
	Layers []AccessControlLayer `json:"layers,omitempty"`

	// Roles defines available roles and permissions.
	Roles []SecurityRole `json:"roles,omitempty"`

	// Policies describes policy enforcement (e.g., Cedar, OPA).
	Policies string `json:"policies,omitempty"`
}

// AccessControlLayer describes access control at a specific layer.
type AccessControlLayer struct {
	// Layer name (e.g., "API Gateway", "Application", "Data").
	Layer string `json:"layer"`

	// Controls implemented at this layer.
	Controls []string `json:"controls"`

	// Description provides additional context.
	Description string `json:"description,omitempty"`
}

// SecurityRole defines a role with permissions.
type SecurityRole struct {
	// ID is the unique identifier for this role.
	ID string `json:"id,omitempty"`

	// Role name.
	Role string `json:"role"`

	// Description of the role.
	Description string `json:"description,omitempty"`

	// Permissions granted to this role.
	Permissions []string `json:"permissions"`

	// Scope defines where this role applies.
	Scope string `json:"scope,omitempty"`
}

// EncryptionRequirements specifies encryption requirements.
type EncryptionRequirements struct {
	// AtRest describes encryption at rest.
	AtRest EncryptionSpec `json:"atRest"`

	// InTransit describes encryption in transit.
	InTransit EncryptionSpec `json:"inTransit"`

	// FieldLevel describes field-level encryption if applicable.
	FieldLevel *EncryptionSpec `json:"fieldLevel,omitempty"`
}

// EncryptionSpec describes encryption configuration.
type EncryptionSpec struct {
	// Method is the encryption method (e.g., "AES-256-GCM").
	Method string `json:"method"`

	// KeyManagement describes key management approach.
	KeyManagement string `json:"keyManagement"`

	// Rotation describes key rotation policy.
	Rotation string `json:"rotation,omitempty"`

	// Provider is the encryption provider (e.g., "AWS KMS", "HashiCorp Vault").
	Provider string `json:"provider,omitempty"`
}

// AuditLogging defines audit logging requirements.
type AuditLogging struct {
	// Scope describes what is logged.
	Scope string `json:"scope"`

	// Events lists specific events that are logged.
	Events []string `json:"events,omitempty"`

	// Format is the log format (e.g., "OCSF", "JSON", "CEF").
	Format string `json:"format,omitempty"`

	// Retention is how long logs are retained.
	Retention string `json:"retention"`

	// Immutability describes tamper-proofing approach.
	Immutability string `json:"immutability,omitempty"`

	// Destination is where logs are stored.
	Destination string `json:"destination,omitempty"`
}

// DataClassification defines data sensitivity classification.
type DataClassification struct {
	// Level is the classification level (e.g., "public", "internal", "confidential", "restricted").
	Level string `json:"level"`

	// Description explains this classification level.
	Description string `json:"description"`

	// Examples of data at this level.
	Examples []string `json:"examples,omitempty"`

	// Handling requirements for this level.
	Handling string `json:"handling,omitempty"`
}
