package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/eugene-eeo/psync/lib"
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Export(filename string) {
	root := lib.BlocksDir()
	f, err := os.Open(filename)
	CheckError(err)
	defer f.Close()
	lib.Chunked(f, func(data []byte, err error) {
		block := lib.NewBlock(data)
		fmt.Println(block.Checksum)
		w, err := os.Create(filepath.Join(
			root,
			string(block.Checksum),
		))
		CheckError(err)
		block.WriteTo(w)
		w.Close()
	})
}
