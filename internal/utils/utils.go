package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"uppies/cli/config"
)

func RequireLogin(cmd *cobra.Command, args []string) {
	if config.Token == "" {
		fmt.Println("You must be logged in to use this command. Run 'uppies login' to authenticate.")
		os.Exit(1)
	}
}

func GetUserInput(prompt string, trim ...bool) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	if len(trim) > 0 && trim[0] {
		input = strings.TrimSpace(input)
	}
	return input
}