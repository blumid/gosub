package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/blumid/gosub/runner"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT) // ctrl + c

	// var (
	// 	stdout1 *os.File = os.Stdout
	// 	stdout2 *os.File
	// )
	// var buf bytes.Buffer
	// stdout2 = &buf

	// Handle SIGINT

	// Keep the main function running
	// select {}
	go func() {
		// for {
		<-sigs
		// signal.Ignore(syscall.SIGINT)
		fmt.Println("fuck you pressed ctrl+c")

		// Save the original stdout
		// originalStdout := os.Stdout

		// Create a new file for the temporary stdout
		// tempFile, err := os.CreateTemp("", "temp_stdout")
		// if err != nil {
		// 	fmt.Println("Error creating temporary file:", err)
		// 	return
		// }

		// os.Stdout = tempFile

		// if runner.MenuShown {
		// 	runner.MenuShown = false
		// 	fmt.Println("let's display sth")
		// 	runner.DisplayMenu()
		// }
		// time.Sleep(5 * time.Second)
		// Restore the original stdout
		// os.Stdout = originalStdout
		// os.Exit(1)

		// }
	}()

	options := runner.ParseOptions()
	runner.Run(options)

}
