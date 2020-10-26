package cmd

import (
	"fmt"
	"time"
)

// Spinner displays indicator that command is processing
func Spinner() {
	for {
		for _, r := range ` _\|/` {
			fmt.Printf("\rProcessing... %c", r)
			time.Sleep(100 * time.Millisecond)
		}
	}
}