package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"github.com/eugene-eeo/psync/blockfs"
	"github.com/docopt/docopt-go"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "psync: ", err)
		os.Exit(1)
	}
}

func export(fs *blockfs.FS, filename string) {
	f, err := os.Open(filename)
	checkErr(err)
	defer f.Close()
	hashlist, err := fs.Export(f)
	checkErr(err)
	hashlist.WriteTo(os.Stdout)
}

func glue(fs *blockfs.FS, filename string, verify bool) {
	f, err := os.Open(filename)
	checkErr(err)
	defer f.Close()
	hashlist, err := blockfs.NewHashList(f)
	checkErr(err)
	for _, checksum := range hashlist {
		cat(fs, checksum, verify)
	}
}

func cat(fs *blockfs.FS, checksum blockfs.Checksum, verify bool) {
	block, err := fs.GetBlock(blockfs.Checksum(checksum))
	checkErr(err)
	if verify {
		if block.Checksum != blockfs.NewChecksum(block.Data) {
			checkErr(errors.New("invalid block: " + string(block.Checksum)))
		}
	}
	block.WriteTo(os.Stdout)
}

func main() {
	usage := `Block and hashlist tool.

Usage:
  psync export <filename>
  psync glue [--verify] <hashlist>
  psync cat [--verify] <checksum>

Options:
  --verify    Verify block contents.
  -h --help   Show this screen.`
	args, _ := docopt.Parse(usage, nil, true, "psync 0.1.0", false)
	user, err := user.Current()
	checkErr(err)
	fs := blockfs.NewFS(filepath.Join(user.HomeDir, ".psync"))
	if args["export"].(bool) {
		export(fs, args["<filename>"].(string))
	}
	if args["glue"].(bool) {
		glue(fs, args["<hashlist>"].(string), args["--verify"].(bool))
	}
	if args["cat"].(bool) {
		cat(
			fs,
			blockfs.Checksum(args["<checksum>"].(string)),
			args["--verify"].(bool),
		)
	}
}
