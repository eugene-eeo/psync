package commands

import (
	"os"
)

func Export(filename string) {
	f, err := os.Open(filename)
	CheckError(err)
	defer f.Close()
	hashlist, err := fs.Export(f)
	CheckError(err)
	hashlist.WriteTo(os.Stdout)
}
