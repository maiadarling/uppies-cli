package commands

import (
	"time"

	"uppies/cli/internal/terminal"

	"github.com/spf13/cobra"
)

func testRun(cmd *cobra.Command, args []string) {
	terminal.RunStage("Logging in", func() {
		// Delay for 5 seconds to simulate work
		time.Sleep(5 * time.Second)
	})
}

func TestCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "test",
		Short:   "test command",
		Args:    cobra.NoArgs,
		PreRun:  nil,
		Run:     testRun,
	}
}