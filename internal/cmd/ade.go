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
		Short: "Commands for automatic enrollment policies",
	}

	adeListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all Automatic Device Enrollments",
		Run:   pkg.ADEList,
	}
)
