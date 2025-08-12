package api

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/models/ade"
	"github.com/bart-lute/addigy-tools/internal/models/ade_service"
	"log/slog"
	"net/http"
)

// ADETokensPoliciesQuery Returns All ADE Tokens in the Account
func ADETokensPoliciesQuery(automaticEnrollmentRequest *ade.AutomaticEnrollmentRequest, adeTokens *[]ade_service.AdeToken) {
	doRequest(http.MethodPost, "oa/ade/tokens/policies/query", &automaticEnrollmentRequest, &adeTokens)
	slog.Debug(fmt.Sprintf("Found %d Automatic Device Enrollment (ADE) Tokens", len(*adeTokens)))
}
