package main

import (
	"github.com/docopt/docopt-go"
	"github.com/eugene-eeo/psync/commands"
)

func main() {
	usage := `usage:
	psync up <addr>
	psync export <filename>
	psync get [--force] <addr> <hashlist>
	psync glue <hashlist>
	`
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
			arguments["--force"].(bool),
		)
		return
	}
	if arguments["glue"].(bool) {
		commands.Glue(
			arguments["<hashlist>"].(string),
		)
	}
}
