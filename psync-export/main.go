package main

import (
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/user"
	"path/filepath"
)

func PsyncBlocksDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".psync/blocks")
}

func InitHome() {
	os.MkdirAll(PsyncBlocksDir(), 0755)
}

func Checksum(data []byte) string {
	b := sha256.Sum256(data)
	return hex.EncodeToString(b[:])
}

func check(err error) {
	if err != nil {
		print(err)
		print("\n")
		os.Exit(1)
	}
}

func main() {
	InitHome()
	args := os.Args[1:]
	if len(args) != 1 {
		print("usage: psync-export <filename>\n")
		os.Exit(1)
	}

	filename := args[0]

	blocks_dir := PsyncBlocksDir()
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	for {
		buffer := make([]byte, 4096)
		actual, _ := f.Read(buffer)
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
