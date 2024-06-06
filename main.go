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
		for {
			<-sigs
			if runner.MenuShown {
				runner.MenuShown = false
				fmt.Println("let's display sth")
				runner.DisplayMenu()
			}
		}
	}()

	// Keep the main function running
	select {}
}
