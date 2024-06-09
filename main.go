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
	signal.Notify(sigs, syscall.SIGINT) // ctrl + c

	options := runner.ParseOptions()
	runner.Run(options)

	// Handle SIGINT

	go func() {
		// for {
		<-sigs

		// Create a new os.Stdout and write the message
		tempStdout := os.Stdout
		os.Stdout = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
		// defer func() { os.Stdout = tempStdout }()

		if runner.MenuShown {
			runner.MenuShown = false
			fmt.Println("let's display sth")
			runner.DisplayMenu()
		}
		os.Stdout = tempStdout
		// }
	}()

	// Keep the main function running
	select {}
}
