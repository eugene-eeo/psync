package commands

import (
	"os"
	"log"
	"net/http"
	"path/filepath"
	"github.com/gorilla/handlers"
)

func Serve(addr string) {
	path := filepath.Join(fs.Path, "blocks")

	r := http.NewServeMux()
	r.Handle("/", handlers.LoggingHandler(
		os.Stdout,
		http.FileServer(http.Dir(path)),
	))

	log.Fatal(http.ListenAndServe(addr, handlers.CompressHandler(r)))
}
