package api

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/models/auditor_service"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
)

func Facts(facts *[]auditor_service.Fact) {
	slog.Debug(fmt.Sprintf("Listing Facts"))

	organisationId := viper.GetString("api.organisationId")
	doRequest(http.MethodGet, fmt.Sprintf("o/%s/facts", organisationId), nil, &facts)

	slog.Debug(fmt.Sprintf("Found %d Facts", len(*facts)))
}
