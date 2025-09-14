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

    // A Command to List Werkplek Pro Clients and some Metadata
    policiesListCmd = &cobra.Command{
        Use:   "list",
        Short: "List of all Policies",
        Run:   pkg.PoliciesList,
    }
)
