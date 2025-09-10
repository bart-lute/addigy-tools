package device_entities

type DeviceAudit struct {
	AgentAuditDate string                `json:"agent_audit_date"`
	AgentId        string                `json:"agent_id"`
	AuditDate      string                `json:"audit_date"`
	Facts          map[string]DeviceFact `json:"facts"`
	OrgId          string                `json:"org_id"`
}
