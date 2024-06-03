package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
				displayMenu()
			}
		}
	}()

	// Keep the main function running
	select {}
}

func displayMenu() {
	fmt.Println("Options:")
	fmt.Println("1. resume")
	fmt.Println("2. targets")
	fmt.Println("3. quit")
	fmt.Println("Enter your choice:")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')

	if choice == "2" {
		displayTargets()
		return
	}

	if choice == "\n" {
		fmt.Println("Continuing...")
		stopMenuShown = false
		return
	}
	index, err := strconv.Atoi(choice[:len(choice)-1])
	if err != nil || index < 1 || index > 10 {
		fmt.Println("Invalid choice. Continuing...")
		stopMenuShown = false
		return
	}

	stopMenuShown = false
}

func displayTargets() {
	fmt.Println("displaying targets...")
}

// I should to check it out this part later:
// for _, cmdArgs := range commandList {
//     ctx, cancel := context.WithCancel(context.Background())
//     cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
//     cmd.Stdout = os.Stdout
//     cmd.Stderr = os.Stderr
//     commands = append(commands, cmd)
//     cancelFuncs = append(cancelFuncs, cancel)
//     go runCommand(cmd)
// }
