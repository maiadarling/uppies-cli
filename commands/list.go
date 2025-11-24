package commands

import (
	"fmt"
	"os"

	"uppies/cli/api"
	"uppies/cli/internal/utils"

	"github.com/spf13/cobra"
	"github.com/olekukonko/tablewriter"
)

func listRun(cmd *cobra.Command, args []string) {
	client := api.NewAPIClient()
	sites, err := client.ListSites()
	if err != nil {
		fmt.Println("Error listing sites:", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "Domains", "Status"})


	for _, site := range sites.Data {
		// Join domains by new line
		domains := ""
		for i, domain := range site.Domains {
			if i > 0 {
				domains += "\n"
			}
			domains += domain
		}
		table.Append([]string{site.Name, domains, site.Status})
	}
	table.Render()
}

func ListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all sites",
		Args:    cobra.NoArgs,
		PreRun:  utils.RequireLogin,
		Run:     listRun,
	}
}