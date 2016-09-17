package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/eugene-eeo/psync/blockfs"
	"os"
	"os/user"
	"path/filepath"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "psync: ", err)
		os.Exit(1)
	}
}

func export(fs *blockfs.FS) {
	hashlist, err := fs.Export(os.Stdin)
	checkErr(err)
	hashlist.WriteTo(os.Stdout)
}

func glue(fs *blockfs.FS, verify bool) {
	hashlist, err := blockfs.NewHashList(os.Stdin)
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
  psync export
  psync glue [--verify]
  psync cat [--verify] <checksum>
  psync up <addr>
  psync get <addr> [--force]

Options:
  --force     Redownload blocks even if they exist.
  --verify    Verify block contents.
  -h --help   Show this screen.`
	args, _ := docopt.Parse(usage, nil, true, "psync 0.1.0", false)
	user, err := user.Current()
	checkErr(err)
	fs, err := blockfs.NewFS(filepath.Join(user.HomeDir, ".psync"))
	checkErr(err)
	if args["export"].(bool) {
		export(fs)
	}
	if args["glue"].(bool) {
		glue(fs, args["--verify"].(bool))
	}
	if args["cat"].(bool) {
		cat(
			fs,
			blockfs.Checksum(args["<checksum>"].(string)),
			args["--verify"].(bool),
		)
	}
	if args["up"].(bool) {
		Serve(
			fs,
			args["<addr>"].(string),
		)
	}
	if args["get"].(bool) {
		Get(
			fs,
			args["<addr>"].(string),
			args["--force"].(bool),
		)
	}
}
