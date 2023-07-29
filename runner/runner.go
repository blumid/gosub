package runner

import "github.com/projectdiscovery/gologger"

const banner = `
					   _
  __ _  ___  ___ _   _| |__
 / _| |/ _ \/ __| | | | |_ \
| (_| | (_) \__ \ |_| | |_) |
 \__, |\___/|___/\__,_|_.__/
 |___/
 `

// Name
const ToolName = `gosub`

// version is the current version of dnsx
const version = `1.1`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\blumid\n\n")
}
