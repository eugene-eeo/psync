package main

import (
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

func Checksum(data []byte) string {
	b := sha256.Sum256(data)
	return hex.EncodeToString(b[:])
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
