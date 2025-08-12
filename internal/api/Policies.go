package api

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/models"
	"log/slog"
	"net/http"
)

func PoliciesQuery(policyQueryRequest *models.PolicyQueryRequest, policies *[]models.Policy) {
	doRequest(http.MethodPost, "oa/policies/query", &policyQueryRequest, &policies)
	slog.Debug(fmt.Sprintf("Found %d Policies", len(*policies)))
}
