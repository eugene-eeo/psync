package commands

import (
	"os"
	"log"
	"strings"
	"path/filepath"
	"github.com/eugene-eeo/psync/blockfs"
	"github.com/parnurzeal/gorequest"
)

type request struct {
	url      string
	checksum blockfs.Checksum
}

type response struct {
	checksum blockfs.Checksum
	block    *blockfs.Block
}

func fetchBlock(requests <-chan *request, responses chan<- *response, done chan<- bool) {
	all_ok := true
	for req := range requests {
		resp, body, errors := gorequest.New().Get(req.url).EndBytes()
		if len(errors) != 0 || resp.StatusCode != 200 {
			log.Println("cannot fetch block: ", req.checksum)
			all_ok = false
			continue
		}
		responses <- &response{
			checksum: req.checksum,
			block:    blockfs.NewBlock(body),
		}
	}
	done <- all_ok
}

func writeBlocks(responses <-chan *response, done chan<- bool) {
	root := filepath.Join(fs.Path, "blocks")
	all_ok := true
	for b := range responses {
		if b.checksum != b.block.Checksum {
			log.Println("invalid block: ", b.checksum)
			all_ok = false
			continue
		}
		f, err := os.Create(filepath.Join(root, string(b.checksum)))
		CheckError(err)
		b.block.WriteTo(f)
		f.Close()
	}
	done <- all_ok
}

func Get(addr string, hashlist_path string, force bool) {
	f, err := os.Open(hashlist_path)
	CheckError(err)
	hashlist, err := blockfs.NewHashList(f)
	CheckError(err)

	requests := make(chan *request, 10)
	responses := make(chan *response, 10)
	fetch_done := make(chan bool)
	write_done := make(chan bool)

	missing := hashlist
	if !force {
		missing = fs.MissingBlocks(hashlist)
	}

	for i := 0; i < 10; i++ {
		go fetchBlock(requests, responses, fetch_done)
	}

	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	go func() {
		for _, checksum := range missing {
			requests <- &request{
				url: addr + "/" + string(checksum),
				checksum: checksum,
			}
		}
		close(requests)
	}()

	all_ok := true
	go writeBlocks(responses, write_done)
	for i := 0; i < 10; i++ {
		fetch_ok := <-fetch_done
		all_ok = all_ok && fetch_ok
	}
	close(responses)
	write_ok := <-write_done
	all_ok = all_ok && write_ok
	if !all_ok {
		os.Exit(1)
	}
}
