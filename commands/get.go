package commands

import (
	"os"
	"log"
	"path/filepath"
	"github.com/parnurzeal/gorequest"
	"github.com/eugene-eeo/psync/lib"
)

type blob struct {
	expected lib.Checksum
	data     []byte
}

func fetchBlock(addr string, hashes <-chan lib.Checksum, blobs chan<- *blob, done chan<- bool) {
	for checksum := range hashes {
		resp, body, errors := gorequest.New().
			Get(addr + "/" + string(checksum)).
			EndBytes()
		if len(errors) != 0 || resp.StatusCode != 200 {
			log.Fatal("cannot fetch block: ", checksum)
		}
		blobs <- &blob{
			expected: checksum,
			data: body,
		}
	}
	done <- true
}

func writeBlobs(blobs <-chan *blob, done chan<- bool) {
	root := lib.PsyncBlocksDir()
	for b := range blobs {
		block := lib.NewBlock(b.data)
		if block.Checksum != b.expected {
			log.Fatal("invalid block: ", b.expected)
		}
		f, err := os.Create(filepath.Join(root, string(b.expected)))
		CheckError(err)
		block.WriteTo(f)
	}
	done <- true
}

func Get(addr string, hashlist_path string) {
	f, err := os.Open(hashlist_path)
	CheckError(err)

	num_workers := 10

	hashes := make(chan lib.Checksum, num_workers)
	blobs := make(chan *blob, num_workers)
	worker_done := make(chan bool)
	writer_done := make(chan bool)

	for i := 0; i < num_workers; i++ {
		go fetchBlock(addr, hashes, blobs, worker_done)
	}
	go writeBlobs(blobs, writer_done)

	go func() {
		lib.ParseHashList(f, func(hash lib.Checksum) {
			//_, err := os.Stat(filepath.Join(root, string(hash)));
			//if err != nil {
			hashes <- hash
			//}
		})
		close(hashes)
	}()

	for i := 0; i < num_workers; i++ {
		<-worker_done
	}
	close(blobs)
	<-writer_done
}
