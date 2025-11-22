package commands

import (
	"fmt"

	"uppies/cli/config"
	"uppies/cli/internal/utils"
	"uppies/cli/api"
	"uppies/cli/internal/terminal"

	"github.com/spf13/cobra"
)

func loginRun(cmd *cobra.Command, args []string) {
	accessKey := utils.GetUserInput("Access Key: ", true)

	if accessKey == "" {
		fmt.Println("No access key provided.")
		return
	}

	config.Token = accessKey
	config.SaveConfig()

	terminal.RunStage("Authenticating", func() {
		client := api.NewAPIClient()
		_, err := client.ListSites()
		if err != nil {
			fmt.Println("Error during authentication:", err)
			return
		}
	})

	fmt.Println("Logged in successfully!")
}

func LoginCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "login",
		Short:   "Authenticate with Uppies",
		Args:    cobra.NoArgs,
		PreRun:  nil,
		Run:     loginRun,
	}
}