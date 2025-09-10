package cmd

import (
	"github.com/bart-lute/addigy-tools/internal/pkg"
	"github.com/spf13/cobra"
)

var (
	customFactsCmd = &cobra.Command{
		Use:     "custom-facts",
		Aliases: []string{"cf"},
		Short:   "Commands for Custom Facts",
	}

	// A Command to List Custom Facts
	customFactsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List of all Custom Facts",
		Run:   pkg.CustomFactsList,
	}
)
