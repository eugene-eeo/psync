package commands

import (
	"os"
	"io"
	"github.com/eugene-eeo/psync/lib"
)

func Glue(hashlist_path string) {
	f, err := os.Open(hashlist_path)
	CheckError(err)
	hl := lib.NewHashList(f)
	hl.Resolve(func (r io.Reader, e error) error {
		CheckError(e)
		buff := make([]byte, lib.BLOCK_SIZE)
		r.Read(buff)
		_, err := os.Stdout.Write(buff)
		return err
	})
}
