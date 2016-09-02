package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
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

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
