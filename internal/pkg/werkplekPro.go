package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models"
	"github.com/bart-lute/addigy-tools/internal/models/device_entities"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"slices"
	"sort"
	"strings"
)

func WerkplekProClients(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Listing WerkPlek Pro Clients"))

	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatalln(err)
	}

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
	policiesSoftwareUpdateMap := *getPoliciesSoftwareUpdateMap(&configurationsProfilesResponse)

	tHeader := table.Row{
		"Policy",
		"Security and Privacy Profiles",
		"Software Update Policies",
	}
	var tRows []table.Row
	for _, policy := range clientPolicies {
		tRows = append(tRows, table.Row{
			policy.Name,
			fmt.Sprintf("%s", strings.Join(policiesSecurityAndPrivacyMap[policy.PolicyID], ", ")),
			fmt.Sprintf("%s", strings.Join(policiesSoftwareUpdateMap[policy.PolicyID], ", ")),
		})
	}
	renderTable(&tHeader, &tRows, csv)

}

func WerkplekProDevicesLocalAdmin(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Listing Local Admin status for all devices"))

	// Sanity Check
	clientPoliciesId := viper.GetString("policies.werkplekPro.clients.id")
	if clientPoliciesId == "" {
		log.Fatal("Config does not contain key: policies.werkplekPro.clients.id")
	}

	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatalln(err)
	}

	// Custom Facts use the ID as Identifier
	hasLocalAdminId := getKeyString("customFacts.ids.hasLocalAdmin")
	secureTokenEnabledForLocalAdminId := getKeyString("customFacts.ids.secureTokenEnabledForLocalAdmin")

	desiredFactIdentifiers := []string{
		"device_name",
		"serial_number",
		"online",
		"policy_id",
		"last_online",
		"device_model_name",
		hasLocalAdminId,
		secureTokenEnabledForLocalAdminId,
	}

	var queryFilter device_entities.QueryFilter

	var items []device_entities.DeviceAudit
	api.Devices(&desiredFactIdentifiers, &queryFilter, "asc", "device_name", &items)

	var policyQueryRequest models.PolicyQueryRequest
	var policies []models.Policy
	api.PoliciesQuery(&policyQueryRequest, &policies)
	policyIds := getChildPolicies(&policies, &clientPoliciesId)

	policiesMap := *getPoliciesMapFromPolicies(&policies)

	tHead := table.Row{
		"DEVICE NAME",
		"SERIAL NUMBER",
		"ONLINE",
		"POLICY",
		"LAST ONLINE",
		"DEVICE_MODEL_NAME",
		"HAS LOCAL_ADMIN USER",
		"LOCAL_ADMIN HAS SECURE_TOKEN",
	}

	var tRows []table.Row
	for _, item := range items {
		if slices.Contains(*policyIds, fmt.Sprintf("%s", item.Facts["policy_id"].Value)) {
			tRow := table.Row{}
			for _, i := range desiredFactIdentifiers {
				value := item.Facts[i].Value
				if value == nil {
					value = "n/a"
				}
				if i == "last_online" {
					//value = item.Facts[i].Value
					value = localDateTimeString(fmt.Sprintf("%s", item.Facts[i].Value))
				} else if i == "policy_id" {
					value = policiesMap[fmt.Sprint(item.Facts[i].Value)].Name
				}
				tRow = append(tRow, value)
			}
			tRows = append(tRows, tRow)

		}
	}
	renderTable(&tHead, &tRows, csv)

}

// WerkplekProDevicesSecureBootLevel A list of devices with data related to Secure Boot Level
func WerkplekProDevicesSecureBootLevel(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Listing Secure Bootlevel status for all devices "))

	// Sanity Check
	clientPoliciesId := viper.GetString("policies.werkplekPro.clients.id")
	if clientPoliciesId == "" {
		log.Fatal("Config does not contain key: policies.werkplekPro.clients.id")
	}

	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatalln(err)
	}

	desiredFactIdentifiers := []string{
		"device_name",
		"online",
		"policy_id",
		"last_online",
		"secure_boot_level",
		"mac_os_x_version",
		"device_model_name",
		"enrolled_via_dep",
	}

	var queryFilter device_entities.QueryFilter

	var items []device_entities.DeviceAudit
	api.Devices(&desiredFactIdentifiers, &queryFilter, "asc", "device_name", &items)

	var policyQueryRequest models.PolicyQueryRequest
	var policies []models.Policy
	api.PoliciesQuery(&policyQueryRequest, &policies)
	policyIds := getChildPolicies(&policies, &clientPoliciesId)

	policiesMap := *getPoliciesMapFromPolicies(&policies)

	tHead := table.Row{
		"DEVICE NAME",
		"ONLINE",
		"POLICY",
		"LAST ONLINE",
		"SECURE_BOOT_LEVEL",
		"MACOS_VERSION",
		"DEVICE_MODEL_NAME",
		"ENROLLED_VIA_ADE",
	}

	var tRows []table.Row
	for _, item := range items {
		if slices.Contains(*policyIds, fmt.Sprintf("%s", item.Facts["policy_id"].Value)) {
			tRow := table.Row{}
			for _, i := range desiredFactIdentifiers {
				value := item.Facts[i].Value
				if value == nil {
					value = "n/a"
				}
				if i == "last_online" {
					//value = item.Facts[i].Value
					value = localDateTimeString(fmt.Sprintf("%s", item.Facts[i].Value))
				} else if i == "policy_id" {
					value = policiesMap[fmt.Sprint(item.Facts[i].Value)].Name
				}
				tRow = append(tRow, value)
			}
			tRows = append(tRows, tRow)
		}
	}
	renderTable(&tHead, &tRows, csv)

}
