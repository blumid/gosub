package main

import (
	"fmt"

	"github.com/blumid/gosub/runner"
	"github.com/eiannone/keyboard"
)

func main() {

	// Handle SIGINT
	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	go func() {
		for {
			_, key, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}

			if key == keyboard.KeyEsc {
				if runner.MenuShown {
					runner.MenuShown = false
					runner.DisplayMenu()
					continue
				}
			}

		}
	}()
	options := runner.ParseOptions()
	runner.Run(options)

}
