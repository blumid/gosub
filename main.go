package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/blumid/gosub/runner"
)

var stopMenuShown bool

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT) // ctrl + c

	options := runner.ParseOptions()
	runner.Run(options)

	// Handle SIGINT
	go func() {
		for {
			<-sigs
			if !stopMenuShown {
				stopMenuShown = true
				runner.DisplayMenu(stopMenuShown)
			}
		}
	}()

	// Keep the main function running
	select {}
}
