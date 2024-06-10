package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blumid/gosub/runner"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT) // ctrl + c

	// Handle SIGINT

	go func() {
		for {
			<-sigs

			// Save the original stdout
			originalStdout := os.Stdout

			// Create a new file for the temporary stdout
			// tempFile, err := os.CreateTemp("", "temp_stdout")
			// if err != nil {
			// 	fmt.Println("Error creating temporary file:", err)
			// 	return
			// }

			// os.Stdout = tempFile

			if runner.MenuShown {
				runner.MenuShown = false
				fmt.Println("let's display sth")
				runner.DisplayMenu()
			}
			time.Sleep(5 * time.Second)
			// Restore the original stdout
			os.Stdout = originalStdout

			time.Sleep(5 * time.Second)
			// os.Exit(1)

		}
	}()

	// Keep the main function running
	// select {}

	options := runner.ParseOptions()
	runner.Run(options)

}
