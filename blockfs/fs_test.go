package blockfs_test

import (
	"bytes"
	"github.com/eugene-eeo/psync/blockfs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func allocTempDir(t *testing.T) string {
	dirname, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return dirname
}

func allocTempFs(t *testing.T) (*blockfs.FS, func()) {
	dirname := allocTempDir(t)
	fs, err := blockfs.NewFS(dirname)
	cleanup := func() {
		os.RemoveAll(dirname)
	}
	if err != nil {
		t.Error("unexpected error:", err)
		cleanup()
		t.Fail()
	}
	return fs, cleanup
}

func TestNewFS(t *testing.T) {
	fs, cleanup := allocTempFs(t)
	defer cleanup()
	_, err := os.Stat(filepath.Join(fs.Path, "blocks"))
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestWriteBlock(t *testing.T) {
	fs, cleanup := allocTempFs(t)
	defer cleanup()
	data := []byte("test-data")
	block := blockfs.NewBlock(data)
	err := fs.WriteBlock(block)
	if err != nil {
		t.Error("unexpected error:", err)
		t.Fail()
	}
	// test that re-writing the same block produces no errors.
	err = fs.WriteBlock(block)
	if err != nil {
		t.Error("unexpected error:", err)
		t.Fail()
	}
	b, _ := fs.GetBlock(block.Checksum)
	if !bytes.Equal(b.Data, data) {
		t.Error("expected data to equal", data, "got", b.Data)
	}
	if b.Checksum != block.Checksum {
		t.Error("expected checksum to equal", block.Checksum, "got", b.Checksum)
	}
}

func TestExport(t *testing.T) {
	fs, cleanup := allocTempFs(t)
	defer cleanup()

	buff := make([]byte, blockfs.BlockSize+blockfs.BlockSize>>1)
	rand.Read(buff)
	hashlist, err := fs.Export(bytes.NewBuffer(buff))

	if err != nil {
		t.Error("unexpected error:", err)
	}

	if len(hashlist) != 2 {
		t.Error("expected 2 blocks to be written, got", len(hashlist))
	}

	dst := bytes.NewBuffer([]byte{})
	for _, checksum := range hashlist {
		block, err := fs.GetBlock(checksum)
		if err != nil {
			t.Error("unexpected error during GetBlock:", err)
			t.Fail()
		}
		dst.Write(block.Data)
	}

	if !bytes.Equal(dst.Bytes(), buff) {
		t.Error("expected resolved chunks to equal")
	}
}

func TestExists(t *testing.T) {
	fs, cleanup := allocTempFs(t)
	defer cleanup()
	block := blockfs.NewBlock([]byte("abc"))
	err := fs.WriteBlock(block)
	if err != nil {
		t.Error("cannot write block:", err)
		t.Fail()
	}
	if !fs.Exists(block.Checksum) {
		t.Error("expected written block to exist")
		t.Fail()
	}
	if fs.Exists(blockfs.Checksum("random-string")) {
		t.Error("expected non-existent block to not exist")
		t.Fail()
	}
}

func TestMissingBlocks(t *testing.T) {
	fs, cleanup := allocTempFs(t)
	defer cleanup()
	blocks := []*blockfs.Block{
		blockfs.NewBlock([]byte("abc")),
		blockfs.NewBlock([]byte("def")),
		blockfs.NewBlock([]byte("ghi")),
	}
	hashlist := blockfs.HashList([]blockfs.Checksum{})
	for _, block := range blocks {
		hashlist = append(hashlist, block.Checksum)
	}
	if !reflect.DeepEqual(fs.MissingBlocks(hashlist), hashlist) {
		t.Error("expected unwritten blocks to be missing")
		t.Fail()
	}
	fs.WriteBlock(blocks[0])
	fs.WriteBlock(blocks[1])
	if !reflect.DeepEqual(fs.MissingBlocks(hashlist), hashlist[2:]) {
		t.Error("expected written blocks to not be included")
		t.Fail()
	}
}
