package main

import "github.com/docopt/docopt-go"

func main() {
	usage := `usage:
	psync up <addr>
	psync export <filename>
	psync get <addr> <hashlist>`
	arguments, _ := docopt.Parse(usage, nil, true, "psync 0.1", false)
	if arguments["up"].(bool) {
		Serve(arguments["<addr>"].(string))
		return
	}
	if arguments["export"].(bool) {
		Export(arguments["<filename>"].(string))
	}
	if arguments["get"].(bool) {
		Get(
			arguments["<addr>"].(string),
			arguments["<hashlist>"].(string),
		)
	}
}
