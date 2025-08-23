package cmd

import (
	"github.com/bart-lute/addigy-tools/internal/pkg"
	"github.com/spf13/cobra"
)

var (
	werkplekProCmd = &cobra.Command{
		Use:     "werkplek-pro",
		Aliases: []string{"wpp"},
		Short:   "Commands for Werkplek Pro",
	}

	// A Command to List Werkplek Pro Clients and some Metadata
	werkplekProClientsCmd = &cobra.Command{
		Use:   "clients",
		Short: "List of all Werkplek Pro clients",
		Run:   pkg.WerkplekProClients,
	}
)
