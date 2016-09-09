package commands

import (
	"github.com/eugene-eeo/psync/blockfs"
	"os/user"
	"path/filepath"
	"log"
)

var fs *blockfs.FS

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	usr, err := user.Current()
	CheckError(err)
	path := filepath.Join(usr.HomeDir, ".psync")
	fs = blockfs.NewFS(path)
}
