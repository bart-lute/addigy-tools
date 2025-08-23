package pkg

import (
    "fmt"
    "github.com/bart-lute/addigy-tools/internal/api"
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/spf13/viper"
    "log"
    "log/slog"
    "os"
    "strings"
    "time"
)

var apiDateTimeFormat = "2006-01-02T15:04:05.999999999Z"

func renderTable(header *table.Row, rows *[]table.Row) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetStyle(table.StyleLight)
    t.AppendHeader(*header)
    t.AppendRows(*rows)
    t.Render()
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

func getPoliciesSecurityAndPrivacyMap(configurationsProfileResponse *api.ConfigurationsProfilesResponse) *map[string][]string {
    payloadTypePrefix := "com.addigy.securityAndPrivacy"
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
