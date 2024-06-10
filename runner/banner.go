package runner

import (
	"fmt"

	"moul.io/banner"
)

// Name
const toolName = `gosub`

// version is the current version of dnsx
const version = `2.1.1`

// showBanner is used to show the banner to the user
func ShowBanner() {
	fmt.Print(banner.Inline(toolName))
	fmt.Printf("\tblumid - %s v%s\n\n", toolName, version)
}
