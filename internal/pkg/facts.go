package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models/auditor_service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"sort"
	"strings"
)

func FactsList(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Listing Facts"))

	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatal(err)
	}

	var facts []auditor_service.Fact
	api.Facts(&facts)

	// Some sorting
	sort.Slice(facts, func(i, j int) bool {
		return strings.ToLower(facts[i].Name) < strings.ToLower(facts[j].Name)
	})

	tHead := table.Row{
		"NAME",
		"ID",
		"RETURN TYPE",
	}

	var tRow []table.Row
	for _, fact := range facts {
		tRow = append(tRow, table.Row{
			fact.Name,
			fact.Identifier,
			fact.ReturnType,
		})
	}

	renderTable(&tHead, &tRow, csv)
}
