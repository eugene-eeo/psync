package main

import (
	"github.com/docopt/docopt-go"
	"github.com/eugene-eeo/psync/commands"
	"github.com/eugene-eeo/psync/lib"
)

func main() {
	usage := `usage:
	psync up <addr>
	psync export <filename>
	psync get <addr> <hashlist>`
	lib.InitHome()
	arguments, _ := docopt.Parse(usage, nil, true, "psync 0.1", false)
	if arguments["up"].(bool) {
		commands.Serve(arguments["<addr>"].(string))
		return
	}
	if arguments["export"].(bool) {
		commands.Export(arguments["<filename>"].(string))
		return
	}
	if arguments["get"].(bool) {
		commands.Get(
			arguments["<addr>"].(string),
			arguments["<hashlist>"].(string),
		)
		return
	}
}
