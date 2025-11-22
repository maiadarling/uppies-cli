package utils

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"uppies/cli/config"
)

func RequireLogin(cmd *cobra.Command, args []string) {
	if config.Token == "" {
		fmt.Println("You must be logged in to use this command. Run 'uppies login' to authenticate.")
		os.Exit(1)
	}
}