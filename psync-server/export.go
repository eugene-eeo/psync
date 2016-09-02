package main

import (
	"os"
	"os/user"
	"path/filepath"
	"log"
	"fmt"
	"io"
	"crypto/sha256"
	"encoding/hex"
)

func Checksum(data []byte) string {
	b := sha256.Sum256(data)
	return hex.EncodeToString(b[:])
}

func PsyncBlocksDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".psync/blocks")
}

func InitHome() {
	os.MkdirAll(PsyncBlocksDir(), 0755)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Export(filename string) {
	blocks_dir := PsyncBlocksDir()
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	for {
		buffer := make([]byte, 4096)
		actual, err := f.Read(buffer)
		if err != io.EOF {
			check(err)
		}
		if actual == 0 {
			break
		}
		digest := Checksum(buffer)
		fmt.Println(digest)
		w, err := os.Create(filepath.Join(blocks_dir, digest))
		check(err)
		defer w.Close()
		w.Write(buffer)
	}
}
