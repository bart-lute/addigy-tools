package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models/fact_entities"
	"github.com/bart-lute/addigy-tools/internal/models/facts"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
)

func CustomFactsList(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Listing Custom Facts"))

	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatalln(err)
	}

	var filter facts.Filter
	sortDirection := "asc"
	sortField := "name"
	var items []fact_entities.Fact

	api.FactsCustomQuery(&filter, sortDirection, sortField, &items)

	tHead := table.Row{"NAME", "ID", "RETURN TYPE"}
	var tRows []table.Row
	for _, item := range items {
		tRows = append(tRows, table.Row{
			item.Name,
			item.Id,
			item.ReturnType,
		})
	}

	renderTable(&tHead, &tRows, csv)
}
