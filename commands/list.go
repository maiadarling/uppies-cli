package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"uppies/cli/api"
	"uppies/cli/internal/utils"
)

func listRun(cmd *cobra.Command, args []string) {
	client := api.NewAPIClient()
	sites, err := client.ListSites()
	if err != nil {
		fmt.Println("Error listing sites:", err)
		os.Exit(1)
	}

	for _, site := range sites.Data {
		fmt.Printf("Name: %s, URL: %s, Status: %s\n", site.Name, site.URL, site.Status)
	}
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