package commands

import (
	"os"
	"path/filepath"
	"github.com/parnurzeal/gorequest"
	"github.com/eugene-eeo/psync/lib"
)

type blob struct {
	checksum lib.Checksum
	data     []byte
}

func fetchBlock(addr string, src <-chan lib.Checksum, dst chan<- *blob) {
	for checksum := range src {
		_, body, _ := gorequest.New().
			Get(addr + "/" + string(checksum)).
			EndBytes()
		dst <- &blob{
			checksum: checksum,
			data: body,
		}
	}
}

func Get(addr string, hashlist_path string) {
	f, err := os.Open(hashlist_path)
	CheckError(err)
	hashes := make(chan lib.Checksum, 10)
	blobs := make(chan *blob, 10)
	num := uint64(0)
	for i := 0; i < 10; i++ {
		go fetchBlock(addr, hashes, blobs)
	}
	root := lib.PsyncBlocksDir()
	go func() {
		lib.ParseHashList(f, func(hash lib.Checksum) {
			_, err := os.Stat(filepath.Join(root, string(hash)));
			if err != nil {
				hashes <- hash
				num++
			}
		})
		close(hashes)
	}()
	for i := uint64(0); i < num; i++ {
		b := <-blobs
		if lib.Checksum(b.data) == b.checksum {
			f, err := os.Open(filepath.Join(root, string(b.checksum)))
			CheckError(err)
			f.Write(b.data)
			f.Close()
		}
	}
}
