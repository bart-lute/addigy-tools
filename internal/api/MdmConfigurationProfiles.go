package api

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/models/policy_service"
	"log/slog"
	"net/http"
)

type ConfigurationsProfilesResponse struct {
	Payloads            []policy_service.MdmPayloadJSONResult   `json:"payloads"`
	PoliciesMdmPayloads []policy_service.PolicyMdmConfiguration `json:"policies_mdm_payloads"`
	StagedPayloads      []policy_service.MdmPayloadJSONResult   `json:"staged_payloads"`
}

func ConfigurationsProfiles(configurationProfilesResponse *ConfigurationsProfilesResponse) {
	doRequest(http.MethodGet, "mdm/configurations/profiles", nil, &configurationProfilesResponse)
	slog.Debug(fmt.Sprintf("Found %d Payloads", len(configurationProfilesResponse.Payloads)))
}
