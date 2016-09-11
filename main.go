package main

import (
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

func glue(fs *blockfs.FS, filename string) {
	f, err := os.Open(filename)
	checkErr(err)
	defer f.Close()
	hashlist, err := blockfs.NewHashList(f)
	checkErr(err)
	for _, checksum := range hashlist {
		block, err := fs.GetBlock(checksum)
		checkErr(err)
		block.WriteTo(os.Stdout)
	}
}

func main() {
	usage := `Block and hashlist tool.

Usage:
  psync export <filename>
  psync glue <hashlist>
  psync -h | --help
  psync --version

Options:
  -h --help   Show this screen.
  --version   Show version.`
	args, _ := docopt.Parse(usage, nil, true, "psync 0.1.0", false)
	user, err := user.Current()
	checkErr(err)
	fs := blockfs.NewFS(filepath.Join(user.HomeDir, ".psync"))
	if args["export"].(bool) {
		export(fs, args["<filename>"].(string))
	}
	if args["glue"].(bool) {
		glue(fs, args["<hashlist>"].(string))
	}
}
