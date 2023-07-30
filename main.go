package main

import (
	"github.com/blumid/tools/gosub/runner"
)

func main() {

	options := runner.ParseOptions()
	runner.Run(options)

}
