package main

import (
	"github.com/blumid/gosub/runner"
)

func main() {

	options := runner.ParseOptions()
	runner.Run(options)

}
