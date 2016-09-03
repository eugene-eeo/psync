package lib

import (
	"os"
	"os/user"
	"path/filepath"
)

func BlocksDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".psync/blocks")
}

func InitHome() {
	os.MkdirAll(BlocksDir(), 0755)
}
