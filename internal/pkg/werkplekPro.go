package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"sort"
	"strings"
)

func WerkplekProClients(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Listing WerkPlek Pro Clients"))
	var policyQueryRequest *models.PolicyQueryRequest
	var policies []models.Policy
	var clientPolicies []models.Policy

	// Sanity Check
	clientPoliciesId := viper.GetString("policies.werkplekPro.clients.id")
	if clientPoliciesId == "" {
		log.Fatal("Config does not contain key: policies.werkplekPro.clients.id")
	}

	api.PoliciesQuery(policyQueryRequest, &policies)

	// Populate the Clients Slice
	for _, policy := range policies {
		if policy.Parent == clientPoliciesId {
			clientPolicies = append(clientPolicies, policy)
		}
	}
	slog.Debug(fmt.Sprintf("Found %d Werkplek Pro Client Policies", len(clientPolicies)))

	// Sorting the result (case-insensitive)
	sort.Slice(clientPolicies, func(i, j int) bool {
		return strings.ToLower(clientPolicies[i].Name) < strings.ToLower(clientPolicies[j].Name)
	})

	var configurationsProfilesResponse api.ConfigurationsProfilesResponse
	api.ConfigurationsProfiles(&configurationsProfilesResponse)

	policiesSecurityAndPrivacyMap := *getPoliciesSecurityAndPrivacyMap(&configurationsProfilesResponse)

	tHeader := table.Row{
		"Policy",
		"Security and Privacy Policies",
	}
	var tRows []table.Row
	for _, policy := range clientPolicies {
		tRows = append(tRows, table.Row{
			policy.Name,
			fmt.Sprintf("%s", strings.Join(policiesSecurityAndPrivacyMap[policy.PolicyID], ", ")),
		})
	}
	renderTable(&tHeader, &tRows)

}
