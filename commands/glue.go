package commands

import (
	"os"
	"github.com/eugene-eeo/psync/blockfs"
)

func Glue(hashlist_path string) {
	f, err := os.Open(hashlist_path)
	CheckError(err)
	hl, err := blockfs.NewHashList(f)
	CheckError(err)
	for _, checksum := range hl {
		block, err := fs.GetBlock(checksum)
		CheckError(err)
		block.WriteTo(os.Stdout)
	}
}
