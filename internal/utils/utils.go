package utils

import (
	"fmt"
	"time"
	"os"

	"github.com/spf13/cobra"
	"uppies/cli/config"
)

const SpinnerDelay = 100 * time.Millisecond
const StatusRetryInterval = 2 * time.Second

func RequireLogin(cmd *cobra.Command, args []string) {
	if config.Token == "" {
		fmt.Println("You must be logged in to use this command. Run 'uppies login' to authenticate.")
		os.Exit(1)
	}
}

func StartSpinner() chan bool {
	stop := make(chan bool)
	go func() {
		chars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				fmt.Printf("\b%s", chars[i%8])
				time.Sleep(SpinnerDelay)
				i++
			}
		}
	}()
	return stop
}

func RunStage(label string, fn func()) {
	fmt.Print("-> " + label + "... ")
	stop := StartSpinner()
	fn()
	stop <- true
	fmt.Println("\b Done")
}