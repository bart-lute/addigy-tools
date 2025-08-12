package cmd

import (
	"github.com/bart-lute/addigy-tools/internal/pkg"
	"github.com/spf13/cobra"
)

var (
	adeCmd = &cobra.Command{
		Use: "automatic-device-enrollment",
		Aliases: []string{
			"ade",
		},
		Short: "Find broken ade synchronization",
	}

	adeListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all Automatic Device Enrollments",
		Run:   pkg.ADEList,
	}
)
