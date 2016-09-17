package main

import (
	"net/http"
	"regexp"
	"log"
	"github.com/eugene-eeo/psync/blockfs"
)

func Serve(fs *blockfs.FS, addr string) {
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
