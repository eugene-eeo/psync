package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func error500(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Print(err)
}

func HttpService(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.URL.Path)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	if strings.Count(path, "/") > 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	blocks_dir := PsyncBlocksDir()
	f, err := os.Open(filepath.Join(blocks_dir, path[1:]))
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		error500(err, w)
		return
	}
	defer f.Close()
	buff := bufpool.Get()
	f.Read(buff)
	w.Write(buff)
	bufpool.Put(buff)
}
