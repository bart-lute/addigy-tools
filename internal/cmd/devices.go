package cmd

import (
	"github.com/bart-lute/addigy-tools/internal/pkg"
	"github.com/spf13/cobra"
)

var (
	devicesCmd = &cobra.Command{
		Use:     "devices",
		Aliases: []string{"dev"},
		Short:   "Commands for Devices",
	}

	// A Command to List Custom Facts
	devicesWithSlackCmd = &cobra.Command{
		Use:   "with-slack",
		Short: "List of Devices with Slack installed",
		Run:   pkg.DevicesWithSlack,
	}
)
