package commands

import (
	"fmt"
	"time"

	"uppies/cli/internal/terminal"
	"uppies/cli/internal/terminal/spinners"

	"github.com/spf13/cobra"
)

func testRun(cmd *cobra.Command, args []string) {
	fmt.Println("Uppies \\o/")
	fmt.Println("")

	terminal.RunStage("Logging in", func() {
		// Delay for 5 seconds to simulate work
		time.Sleep(3 * time.Second)
	}, spinners.Lines)

	terminal.RunStage("Logging in", func() {
		// Delay for 5 seconds to simulate work
		time.Sleep(3 * time.Second)
	}, spinners.Dots)

	terminal.RunStage("Logging in", func() {
		// Delay for 5 seconds to simulate work
		time.Sleep(3 * time.Second)
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