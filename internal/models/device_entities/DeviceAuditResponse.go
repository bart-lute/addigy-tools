package device_entities

import "github.com/bart-lute/addigy-tools/internal/models/response_entities"

type DeviceAuditResponse struct {
	Items    []DeviceAudit              `json:"items"`
	Metadata response_entities.Metadata `json:"metadata"`
}
