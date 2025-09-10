package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models"
	"github.com/bart-lute/addigy-tools/internal/models/ade"
	"github.com/bart-lute/addigy-tools/internal/models/ade_service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"sort"
	"time"
)

func getADETokens() *[]ade_service.AdeToken {
	// Empty Filter (get ALL)
	var automaticEnrollmentRequest ade.AutomaticEnrollmentRequest
	var adeTokens []ade_service.AdeToken
	api.ADETokensPoliciesQuery(&automaticEnrollmentRequest, &adeTokens)
	sort.Slice(adeTokens, func(i, j int) bool {
		return adeTokens[i].Account.OrgName < adeTokens[j].Account.OrgName
	})
	return &adeTokens
}

func getPolicyMap(adeTokens *[]ade_service.AdeToken) map[string]*models.Policy {
	// We need a bit of Metadata from the Policy
	// Instead of bluntly retrieving all we are going to build a filter
	var policyQueryRequest models.PolicyQueryRequest
	var policyIds []string
	for _, adeToken := range *adeTokens {
		policyIds = append(policyIds, adeToken.PolicyID)
	}
	policyQueryRequest.Policies = policyIds
	var policies []models.Policy
	api.PoliciesQuery(&policyQueryRequest, &policies)

	// Now that we have the Token and the Policy, we need a way to simply retrieve the Metadata
	// I'm sure this can be done much simpler, but I'm not smart enough
	policyMap := make(map[string]*models.Policy)
	for _, policy := range policies {
		policyMap[policy.PolicyID] = &policy
	}
	return policyMap
}

func hasError(adeToken *ade_service.AdeToken) bool {
	lastScanTime, err := time.Parse(apiDateTimeFormat, adeToken.LastScanTime)
	if err != nil {
		return true
	}
	if !adeToken.DevicesSyncCompleted {
		return true
	}
	if adeToken.SyncingError != "" {
		return true
	}
	if int(time.Since(lastScanTime).Hours()) > 1 {
		return true
	}
	return false
}

func ADEList(cmd *cobra.Command, args []string) {

	brokenOnly, err := cmd.Flags().GetBool("broken-only")
	if err != nil {
		log.Fatal(err)
	}
	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatalln(err)
	}
	slog.Debug(fmt.Sprintf("broken-only: %v", brokenOnly))

	adeTokens := getADETokens()
	policyMap := getPolicyMap(adeTokens) // Used to retrieve metadata from the Policy

	tableHeader := table.Row{
		"policy name",
		"last scan",
	}
	var rows []table.Row
	for _, adeToken := range *adeTokens {
		if !brokenOnly || hasError(&adeToken) {
			rows = append(rows, table.Row{
				policyMap[adeToken.PolicyID].Name,
				localDateTimeString(adeToken.LastScanTime),
			})
		}
	}
	renderTable(&tableHeader, &rows, csv)
}
