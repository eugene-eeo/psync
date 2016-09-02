package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Export(filename string) {
	root := PsyncBlocksDir()
	f, err := os.Open(filename)
	check(err)
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
		check(err)
		w.Write(buffer)
		w.Close()
	}
}
