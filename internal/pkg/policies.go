package pkg

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/api"
	"github.com/bart-lute/addigy-tools/internal/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"sort"
	"strings"
)

type treeItem struct {
	policy   *models.Policy
	children *[]treeItem
}

func buildTree(policies *[]models.Policy, parentId string) *[]treeItem {
	// Find Sub Policies
	var tree []treeItem
	for _, policy := range *policies {
		if policy.Parent == parentId {
			tree = append(tree, treeItem{
				&policy,
				buildTree(policies, policy.PolicyID), // Traverse
			})
		}
	}

	// Sort the Slice by Policy Name
	sort.Slice(tree, func(i, j int) bool {
		return strings.ToLower(tree[i].policy.Name) < strings.ToLower(tree[j].policy.Name)
	})

	// Return the Tree (recursively)
	return &tree
}

func printTree(tree *[]treeItem, indent string) {
	for _, x := range *tree {
		fmt.Printf("%s%s (%s)\n", indent, x.policy.Name, x.policy.PolicyID)
		if len(*x.children) > 0 {
			printTree(x.children, indent+"  ")
		}
	}
}

func fetchRows(tree *[]treeItem, indent string, tRows *[]table.Row) {
	for _, item := range *tree {
		tRow := table.Row{
			fmt.Sprintf("%s%s", indent, item.policy.Name),
		}
		*tRows = append(*tRows, tRow)
		if len(*item.children) > 0 {
			fetchRows(item.children, indent+"  ", tRows)
		}
	}
}

func PoliciesList(cmd *cobra.Command, args []string) {
	slog.Debug(fmt.Sprintf("Listing Policies"))

	csv, err := cmd.Flags().GetBool("csv")
	if err != nil {
		log.Fatal(err)
	}

	var policyQueryRequest models.PolicyQueryRequest
	var policies []models.Policy
	api.PoliciesQuery(&policyQueryRequest, &policies)

	tree := buildTree(&policies, "")
	//printTree(tree, "")

	tHeader := table.Row{
		"POLICY",
	}
	var tRows []table.Row
	fetchRows(tree, "", &tRows)

	renderTable(&tHeader, &tRows, csv)
}
