package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func Export(filename string) {
	root := PsyncBlocksDir()
	f, err := os.Open(filename)
	CheckError(err)
	defer f.Close()
	for {
		buffer := make([]byte, 8192)
		actual, _ := f.Read(buffer)
		if actual == 0 {
			break
		}
		digest := Checksum(buffer)
		fmt.Println(digest)
		w, err := os.Create(filepath.Join(root, digest))
		CheckError(err)
		w.Write(buffer)
		w.Close()
	}
}
