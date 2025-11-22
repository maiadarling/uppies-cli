package spinners

import (
	"fmt"
	"time"
)

type Spinner struct {
	Chars []string
	Delay time.Duration
}

var BigDots = Spinner{
	Chars: []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
	Delay: 100 * time.Millisecond,
}

var Dots = Spinner{
	Chars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	Delay: 80 * time.Millisecond,
}

var Lines = Spinner{
	Chars: []string{"|", "/", "-", "\\"},
	Delay: 200 * time.Millisecond,
}

func StartSpinner(s Spinner) chan bool {
	stop := make(chan bool)
	go func() {
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				fmt.Printf("\r%s", s.Chars[i%len(s.Chars)])
				time.Sleep(s.Delay)
				i++
			}
		}
	}()
	return stop
}


