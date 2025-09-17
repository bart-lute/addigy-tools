package pkg

import (
    "fmt"
    "github.com/bart-lute/addigy-tools/internal/api"
    "github.com/bart-lute/addigy-tools/internal/models"
    "github.com/bart-lute/addigy-tools/internal/models/device_entities"
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/spf13/cobra"
    "log"
    "log/slog"
    "slices"
    "sort"
    "strings"
)

func getWerkplekProDevices(policies *[]models.Policy, desiredFactIdentifiers *[]string, queryFilter *device_entities.QueryFilter, sortDirection string, sortField string) *[]device_entities.DeviceAudit {

    // Check if desiredFactIdentifiers contains "policy_id". If not add it. We depend on it
    if !slices.Contains(*desiredFactIdentifiers, "policy_id") {
        *desiredFactIdentifiers = append(*desiredFactIdentifiers, "policy_id")
    }

    // A slice of strings containing all Client sub Policy IDs
    policyIds := getChildPolicies(policies, getClientPolicyId())

    var items []device_entities.DeviceAudit
    var werkplekProItems []device_entities.DeviceAudit
    api.Devices(desiredFactIdentifiers, queryFilter, sortDirection, sortField, &items)
    for _, item := range items {
        if slices.Contains(*policyIds, fmt.Sprintf("%s", item.Facts["policy_id"].Value)) {
            werkplekProItems = append(werkplekProItems, item)
        }
    }

    return &werkplekProItems

}

// WerkplekProClients
// This is a somewhat simple filter. It assumes that Werkplek Pro Clients
// (and ONLY Clients) are located directly underneath the Werkplek Pro / Clients policy
func WerkplekProClients(cmd *cobra.Command, args []string) {
    slog.Debug(fmt.Sprintf("Listing WerkPlek Pro Clients"))

    csv, err := cmd.Flags().GetBool("csv")
    if err != nil {
        log.Fatalln(err)
    }

    var policyQueryRequest *models.PolicyQueryRequest
    var policies []models.Policy
    var clientPolicies []models.Policy

    api.PoliciesQuery(policyQueryRequest, &policies)

    // Populate the Clients Slice
    for _, policy := range policies {
        if policy.Parent == *getClientPolicyId() {
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
        "filevault_key_escrowed",
    }

    // Get all Policies
    var policyQueryRequest models.PolicyQueryRequest
    var policies []models.Policy
    api.PoliciesQuery(&policyQueryRequest, &policies)

    // Create a map, which we can use to retrieve Metadata
    policiesMap := *getPoliciesMapFromPolicies(&policies)

    // We are going to filter the devices to show only Werkplek Pro (Client) devices
    var queryFilter device_entities.QueryFilter
    items := getWerkplekProDevices(&policies, &desiredFactIdentifiers, &queryFilter, "asc", "device_name")

    tHead := table.Row{
        "DEVICE NAME",
        "SERIAL NUMBER",
        "ONLINE",
        "POLICY",
        "LAST ONLINE",
        "DEVICE_MODEL_NAME",
        "HAS LOCAL_ADMIN USER",
        "LOCAL_ADMIN HAS SECURE_TOKEN",
        "FILEVAULT_KEY_ESCROWED",
    }

    var tRows []table.Row
    for _, item := range *items {
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
    renderTable(&tHead, &tRows, csv)

}

// WerkplekProDevicesSecureBootLevel A list of devices with data related to Secure Boot Level
func WerkplekProDevicesSecureBootLevel(cmd *cobra.Command, args []string) {
    slog.Debug(fmt.Sprintf("Listing Secure Bootlevel status for all devices "))

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

    // Get all Policies
    var policyQueryRequest models.PolicyQueryRequest
    var policies []models.Policy
    api.PoliciesQuery(&policyQueryRequest, &policies)

    // Create a map, useful for later
    policiesMap := *getPoliciesMapFromPolicies(&policies)

    // We are going to filter the devices to show only Werkplek Pro (Client) devices
    var queryFilter device_entities.QueryFilter
    items := getWerkplekProDevices(&policies, &desiredFactIdentifiers, &queryFilter, "asc", "device_name")

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
    for _, item := range *items {
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
    renderTable(&tHead, &tRows, csv)
}

// werkplekProSimpleList A wrapper command for simple lists
func werkplekProSimpleList(cmd *cobra.Command, desiredFactIdentifiers *[]string, auditFilters *[]device_entities.AuditFilter, tableHeader *table.Row, sortDirection string, sortField string) {

    slog.Debug(fmt.Sprintf("Creating a Simple List"))
    csv, err := cmd.Flags().GetBool("csv")
    if err != nil {
        log.Fatalln(err)
    }

    // Get all Policies
    var policyQueryRequest models.PolicyQueryRequest
    var policies []models.Policy
    api.PoliciesQuery(&policyQueryRequest, &policies)

    // Create a map, useful for later
    policiesMap := *getPoliciesMapFromPolicies(&policies)

    queryFilter := device_entities.QueryFilter{
        Filters: auditFilters,
    }

    items := getWerkplekProDevices(&policies, desiredFactIdentifiers, &queryFilter, sortDirection, sortField)

    var tRows []table.Row
    for _, item := range *items {
        tRow := table.Row{}
        for _, i := range *desiredFactIdentifiers {
            value := item.Facts[i].Value
            if value == nil {
                value = "n/a"
            }
            if i == "policy_id" {
                tRow = append(tRow, policiesMap[fmt.Sprint(value)].Name)
            } else {
                tRow = append(tRow, value)
            }
        }
        tRows = append(tRows, tRow)
    }
    renderTable(tableHeader, &tRows, csv)

}

func WerkplekProDevicesWithSlack(cmd *cobra.Command, args []string) {
    slog.Debug(fmt.Sprintf("Devices With Slack Installed"))

    // Custom Facts use the ID as Identifier
    slackVersionId := getKeyString("customFacts.ids.slackVersion")
    hasSlackInstalledId := getKeyString("customFacts.ids.hasSlackInstalled")

    desiredFactIdentifiers := []string{
        "device_name",
        "online",
        "policy_id",
        slackVersionId,
    }

    auditFilters := []device_entities.AuditFilter{
        {
            AuditField: hasSlackInstalledId,
            Operation:  "=",
            Type:       "boolean",
            Value:      true,
        },
    }

    tableHeader := table.Row{
        "DEVICE NAME",
        "ONLINE",
        "POLICY",
        "SLACK VERSION",
    }

    werkplekProSimpleList(cmd, &desiredFactIdentifiers, &auditFilters, &tableHeader, "asc", "device_name")
}

func WerkplekProDevicesWithDropbox(cmd *cobra.Command, args []string) {
    slog.Debug(fmt.Sprintf("Devices With Dropbox Installed"))

    // Custom Facts use the ID as Identifier
    dropboxVersionId := getKeyString("customFacts.ids.dropboxVersion")
    hasDropboxInstalledId := getKeyString("customFacts.ids.hasDropboxInstalled")

    desiredFactIdentifiers := []string{
        "device_name",
        "online",
        "policy_id",
        dropboxVersionId,
    }

    auditFilters := []device_entities.AuditFilter{
        {
            AuditField: hasDropboxInstalledId,
            Operation:  "=",
            Type:       "boolean",
            Value:      true,
        },
    }

    tableHeader := table.Row{
        "DEVICE NAME",
        "ONLINE",
        "POLICY",
        "DROPBOX VERSION",
    }

    werkplekProSimpleList(cmd, &desiredFactIdentifiers, &auditFilters, &tableHeader, "asc", "device_name")
}

func WerkplekProDevicesOnline(cmd *cobra.Command, args []string) {
    slog.Debug(fmt.Sprintf("Werkplek Pro Online Devices"))

    serials, err := cmd.Flags().GetStringSlice("serials")
    if err != nil {
        log.Fatalln(err)
    }

    desiredFactIdentifiers := []string{
        "device_name",
        "policy_id",
        "agentid",
    }

    auditFilters := []device_entities.AuditFilter{
        {
            AuditField: "online",
            Operation:  "=",
            Type:       "boolean",
            Value:      true,
        },
    }

    if len(serials) > 0 {
        if len(serials) > 0 {
            auditFilters = append(auditFilters, device_entities.AuditFilter{
                AuditField: "serial_number",
                Operation:  "contains",
                Type:       "string",
                Value:      serials,
            })
        }
    }

    tableHeader := table.Row{
        "DEVICE NAME",
        "POLICY",
        "LINK",
    }

    // Get all Policies
    var policyQueryRequest models.PolicyQueryRequest
    var policies []models.Policy
    api.PoliciesQuery(&policyQueryRequest, &policies)

    // Create a map, useful for later
    policiesMap := *getPoliciesMapFromPolicies(&policies)

    queryFilter := device_entities.QueryFilter{
        Filters: &auditFilters,
    }

    items := getWerkplekProDevices(&policies, &desiredFactIdentifiers, &queryFilter, "asc", "device_name")

    portalUrl := getKeyString("portal.url")

    var tRows []table.Row
    for _, item := range *items {
        tRow := table.Row{
            item.Facts["device_name"].Value,
            policiesMap[fmt.Sprint(item.Facts["policy_id"].Value)].Name,
            fmt.Sprintf("%s/devices/%s", portalUrl, item.Facts["agentid"].Value),
        }
        tRows = append(tRows, tRow)
    }
    renderTable(&tableHeader, &tRows, false)

}
