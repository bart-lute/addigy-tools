package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models/device_entities"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
)

func DevicesWithSlack(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Devices With Slack Installed"))

	// Custom Facts use the ID as Identifier
	slackVersionId := getKeyString("customFacts.ids.slackVersion")
	hasSlackInstalledId := getKeyString("customFacts.ids.hasSlackInstalled")

	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatalln(err)
	}

	desiredFactIdentifiers := []string{
		"device_name",
		"online",
		"policy_id",
		slackVersionId,
	}

	var auditFilters []device_entities.AuditFilter

	auditFilters = append(auditFilters, device_entities.AuditFilter{
		AuditField: hasSlackInstalledId,
		Operation:  "=",
		Type:       "boolean",
		Value:      true,
	})

	queryFilter := device_entities.QueryFilter{
		Filters: &auditFilters,
	}

	var items []device_entities.DeviceAudit
	api.Devices(&desiredFactIdentifiers, &queryFilter, "asc", "device_name", &items)
	policiesMap := *getPoliciesMap()

	tHead := table.Row{
		"DEVICE NAME",
		"ONLINE",
		"POLICY",
		//"LAST ONLINE",
		"SLACK VERSION",
	}
	var tRows []table.Row
	for _, item := range items {
		tRow := table.Row{}
		for _, i := range desiredFactIdentifiers {
			if i == "policy_id" {
				tRow = append(tRow, policiesMap[fmt.Sprint(item.Facts[i].Value)].Name)
			} else {
				tRow = append(tRow, item.Facts[i].Value)
			}
		}
		tRows = append(tRows, tRow)
	}
	renderTable(&tHead, &tRows, csv)
}
