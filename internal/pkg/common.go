package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"
)

var apiDateTimeFormat = "2006-01-02T15:04:05.999999999Z"

func renderTable(header *table.Row, rows *[]table.Row, csv bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(*header)
	t.AppendRows(*rows)

	if csv {
		t.RenderCSV()
	} else {
		t.Render()
	}
}

func getLocation() *time.Location {
	location, err := time.LoadLocation(viper.GetString("location"))
	if err != nil {
		log.Fatal(err)
	}
	return location
}

func localDateTimeString(apiDate string) string {
	t, err := time.Parse(apiDateTimeFormat, apiDate)
	if err != nil {
		return ""
	}
	return t.In(getLocation()).Format("2006-01-02 15:04:05")
}

func getPoliciesProfileMap(configurationsProfileResponse *api.ConfigurationsProfilesResponse, payloadTypePrefix string) *map[string][]string {
	//payloadTypePrefix := "com.addigy.securityAndPrivacy"
	payloads := &configurationsProfileResponse.Payloads
	policiesMdmPayloads := &configurationsProfileResponse.PoliciesMdmPayloads
	groups := make(map[string]string)

	// Populate a Map with relevant Payload Group Names
	for _, payload := range *payloads {
		if strings.HasPrefix(payload.AddigyPayloadType, payloadTypePrefix) {
			if _, ok := groups[payload.PayloadGroupID]; !ok {
				slog.Debug(fmt.Sprintf("Found Profile with Name %s, Group Id: %s", payload.PayloadDisplayName, payload.PayloadGroupID))
				groups[payload.PayloadGroupID] = payload.PayloadDisplayName
			}
		}
	}

	// Next we are building a Map with Key = PolicyId and Value = Slice of Strings with Profile Names
	// There should only be 1 per Policy by convention, but you never know...
	m := make(map[string][]string)
	for _, policiesMdmPayload := range *policiesMdmPayloads {
		if _, ok := groups[policiesMdmPayload.ConfigurationID]; ok {
			slog.Debug(fmt.Sprintf("Found Profile with Name: %s, for Policy Id: %s", groups[policiesMdmPayload.ConfigurationID], policiesMdmPayload.PolicyID))
			m[policiesMdmPayload.PolicyID] = append(m[policiesMdmPayload.PolicyID], groups[policiesMdmPayload.ConfigurationID])
		}
	}
	return &m
}

func getPoliciesSecurityAndPrivacyMap(configurationsProfileResponse *api.ConfigurationsProfilesResponse) *map[string][]string {
	return getPoliciesProfileMap(configurationsProfileResponse, "com.addigy.securityAndPrivacy")
}

func getPoliciesSoftwareUpdateMap(configurationsProfileResponse *api.ConfigurationsProfilesResponse) *map[string][]string {
	return getPoliciesProfileMap(configurationsProfileResponse, "com.addigy.softwareupdate.com.apple.softwareupdate")
}

// Fetch all policies and create a Map, to easily retrieve data
func getPoliciesMap() *map[string]models.Policy {
	var policyQueryRequest models.PolicyQueryRequest
	var policies []models.Policy
	api.PoliciesQuery(&policyQueryRequest, &policies)

	return getPoliciesMapFromPolicies(&policies)
}

// Fetch all policies and create a Map, to easily retrieve data
func getPoliciesMapFromPolicies(policies *[]models.Policy) *map[string]models.Policy {
	policiesMap := make(map[string]models.Policy)
	for _, policy := range *policies {
		policiesMap[policy.PolicyID] = policy
	}
	return &policiesMap
}

// getKeyString returns a viper.GetString or a FATAL if not found
func getKeyString(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Fatalf("Key %s not found in config", key)
	}
	return value
}

// getChildPolicies Get a list of all Child Policies, recursively
func getChildPolicies(policies *[]models.Policy, policyId *string) *[]string {

	var policyIds []string

	for _, policy := range *policies {
		if policy.Parent == *policyId {
			policyIds = append(policyIds, policy.PolicyID)
			policyIds = append(policyIds, *getChildPolicies(policies, &policy.PolicyID)...)
		}
	}

	return &policyIds
}
