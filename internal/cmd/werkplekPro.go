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

	// A Command to List Devices With Slack Installed
	werkplekProDevicesWithSlackCmd = &cobra.Command{
		Use:   "slack",
		Short: "List of Devices with Slack installed",
		Run:   pkg.WerkplekProDevicesWithSlack,
	}

	// WerkplekProDevicesLocalAdminCmd Device has a Local Admin user configured
	werkplekProDevicesLocalAdminCmd = &cobra.Command{
		Use:   "local-admin",
		Short: "Local Admin configured for Werkplek Pro Devices",
		Run:   pkg.WerkplekProDevicesLocalAdmin,
	}

	// Kind of a dashboard for Werkplek Pro Devices
	werkplekProDevicesSecureBootlevelCmd = &cobra.Command{
		Use: "secure-boot-level",
		Aliases: []string{
			"bootlevel",
			"sbl",
		},
		Short: "Secure Boot Level information for Werkplek Pro Devices",
		Run:   pkg.WerkplekProDevicesSecureBootLevel,
	}

	// A Command to List Devices With Slack Installed
	werkplekProDevicesWithDropboxCmd = &cobra.Command{
		Use:   "dropbox",
		Short: "List of Devices with Dropbox installed",
		Run:   pkg.WerkplekProDevicesWithDropbox,
	}

	// List Online Devices from a Slice of Device IDS
	werkplekProDevicesOnlineCmd = &cobra.Command{
		Use:   "online",
		Short: "A list of Online Werkplek Pro Devices",
		Run:   pkg.WerkplekProDevicesOnline,
	}
)
