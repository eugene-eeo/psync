package main

import (
	"github.com/docopt/docopt-go"
	"github.com/oxtoacart/bpool"
	"log"
	"net/http"
)

var bufpool *bpool.BytePool

func main() {
	InitHome()
	usage := `usage:
	psync-server export <filename>
	psync-server up <addr>`
	arguments, _ := docopt.Parse(usage, nil, true, "0.1", false)
	if arguments["export"].(bool) {
		Export(arguments["<filename>"].(string))
		return
	}
	if arguments["up"].(bool) {
		bufpool = bpool.NewBytePool(4096*4, 4096)
		http.HandleFunc("/", HttpService)
		err := http.ListenAndServe(arguments["<addr>"].(string), nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
