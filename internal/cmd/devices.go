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

    // Kind of a dashboard for Werkplek Pro Devices
    devicesSecureBootlevelCmd = &cobra.Command{
        Use: "secure-boot-level",
        Aliases: []string{
            "bootlevel",
            "sbl",
        },
        Short: "Secure Boot Level information for all Devices",
        Run:   pkg.DevicesSecureBootLevel,
    }
)
