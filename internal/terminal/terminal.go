package terminal

import (
	"fmt"
	"time"

	"uppies/cli/internal/terminal/spinners"
)

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

func clearLine() {
	fmt.Print(string(ClearLine))
}

func hideCursor() {
	fmt.Print(string(HideCursor))
}

func showCursor() {
	fmt.Print(string(ShowCursor))
}

func RunStage(label string, fn func(), spinner ...spinners.Spinner) {
	hideCursor()

	// Hide the cursor
	fmt.Print("  " + string(Cyan) + label + " " + string(Reset))
	s := spinners.Dots
	if len(spinner) > 0 {
		s = spinner[0]
	}
	stop := spinners.StartSpinner(s)
	fn()
	stop <- true

	fmt.Print("\r")

	fmt.Print("✓ " + string(Green) + label + "      " + string(Reset))

	fmt.Println("Done")

	showCursor()
}

