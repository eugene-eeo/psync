package main

import (
	"net/http"
	"log"
	"os/user"
	"path/filepath"
	"regexp"
	"github.com/eugene-eeo/psync/blockfs"
	"flag"
)

func main() {
	addrPtr := flag.String("addr", ":8000", "address")
	flag.Parse()

	addr := *addrPtr

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(usr.HomeDir, ".psync")
	fs, err := blockfs.NewFS(path)
	if err != nil {
		log.Fatal(err)
	}

	pat := regexp.MustCompile("^[a-f0-9]{64}$")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path[1:]
		if !pat.MatchString(p) {
			w.WriteHeader(404)
			return
		}
		block, err := fs.GetBlock(blockfs.Checksum(p))
		if err != nil {
			w.WriteHeader(404)
			return
		}
		w.Write(block.Data)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
