package cmd

import (
	"github.com/bart-lute/addigy-tools/internal/pkg"
	"github.com/spf13/cobra"
)

var (
	factsCmd = &cobra.Command{
		Use:   "facts",
		Short: "Commands for Facts",
	}

	// A Command to List Custom Facts
	factsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List of all Facts",
		Run:   pkg.FactsList,
	}
)
