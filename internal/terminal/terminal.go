package terminal

import (
	"fmt"
	"time"
)

const SpinnerDelay = 100 * time.Millisecond
const StatusRetryInterval = 2 * time.Second

type AnsiCode string

const (
	// Colors
	Reset   AnsiCode = "\033[0m"
	Red     AnsiCode = "\033[31m"
	Green   AnsiCode = "\033[32m"
	Yellow  AnsiCode = "\033[33m"
	Blue    AnsiCode = "\033[34m"
	Magenta AnsiCode = "\033[35m"
	Cyan    AnsiCode = "\033[36m"
	Gray    AnsiCode = "\033[37m"
	White   AnsiCode = "\033[97m"

	// Cursor control
	HideCursor AnsiCode = "\033[?25l"
	ShowCursor AnsiCode = "\033[?25h"

	// Screen control
	ClearScreen AnsiCode = "\033[2J"
	ClearLine   AnsiCode = "\033[2K"
)

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
				moveBOL()
				fmt.Printf("%s", chars[i%8])
				time.Sleep(SpinnerDelay)
				i++
			}
		}
	}()
	return stop
}

func moveBOL() {
	fmt.Print("\r")
}

func clearLine() {
	fmt.Print(string(ClearLine))
}

func RunStage(label string, fn func()) {
	// Hide the cursor
	fmt.Print("  " + string(Cyan) + label + "... " + string(Reset))
	stop := StartSpinner()
	fn()
	stop <- true

	moveBOL()

	fmt.Print("✓ " + string(Green) + label + "      " + string(Reset))



	fmt.Println("Done")
}

