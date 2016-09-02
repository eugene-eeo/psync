package main

import (
	"github.com/oxtoacart/bpool"
	"log"
	"os"
	"os/user"
	"strings"
	"path/filepath"
	"net/http"
)

var bufpool *bpool.BytePool

func PsyncBlocksDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".psync/blocks")
}

func InitHome() {
	os.MkdirAll(PsyncBlocksDir(), 0755)
}

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

func main() {
	InitHome()
	args := os.Args[1:]
	if len(args) != 1 {
		print("usage: psync-server <addr>\n")
		os.Exit(1)
	}
	addr := args[0]
	bufpool = bpool.NewBytePool(4096*4, 4096)
	http.HandleFunc("/", HttpService)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
