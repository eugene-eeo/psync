package main

import (
	"os"
	"strings"
	"bufio"
	"path/filepath"
	"github.com/parnurzeal/gorequest"
)

type blob struct {
	valid    bool
	expected string
	data     []byte
}

func fetchBlock(addr string, src <-chan string, dst chan<- *blob) {
	for checksum := range src {
		_, body, _ := gorequest.New().
			Get(addr + "/" + checksum).
			EndBytes()
		valid := checksum == Checksum(body)
		dst <- &blob{
			valid: valid,
			expected: checksum,
			data: body,
		}
	}
}

func Get(addr string, hashlist_path string) {
	f, err := os.Open(hashlist_path)
	CheckError(err)
	reader := bufio.NewReader(f)
	hashes := make(chan string, 10)
	blobs := make(chan *blob, 10)
	num := uint64(0)
	for i := 0; i < 10; i++ {
		go fetchBlock(addr, hashes, blobs)
	}
	root := PsyncBlocksDir()
	go func() {
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			if len(line) == 0 {
				break
			}
			if _, err := os.Stat(filepath.Join(root, line)); err != nil {
				hashes <- line
				num++
			}
		}
		close(hashes)
	}()
	for i := uint64(0); i < num; i++ {
		b := <-blobs
		if b.valid {
			f, err := os.Open(filepath.Join(root, b.expected))
			CheckError(err)
			f.Write(b.data)
			f.Close()
		}
	}
}
