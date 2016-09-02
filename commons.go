package main

import (
	"os/user"
	"path/filepath"
	"crypto/sha256"
	"encoding/hex"
	"os"
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
