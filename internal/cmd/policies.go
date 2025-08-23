package cmd

import (
	"github.com/bart-lute/addigy-tools/internal/pkg"
	"github.com/spf13/cobra"
)

var (
	policiesCmd = &cobra.Command{
		Use:   "policies",
		Short: "Commands for Policies",
	}

	//policiesListCmd = &cobra.Command{
	//	Use:   "list",
	//	Short: "List Policies",
	//	Run:   pkg.PoliciesList,
	//}

	policiesWerkplekProCmd = &cobra.Command{
		Use: "werkplek-pro",
		Aliases: []string{
			"wpp",
		},
		Short: "Commands for WerkPlek Pro",
	}

	policiesWerkplekProClientsCmd = &cobra.Command{
		Use:   "clients",
		Short: "List of all Werkplek Pro clients",
		Run:   pkg.PoliciesWerkplekProClients,
	}
)
