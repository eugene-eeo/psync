package blockfs_test

import (
	"bytes"
	"github.com/eugene-eeo/psync/blockfs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func allocTempDir(t *testing.T) string {
	dirname, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return dirname
}

func TestNewFS(t *testing.T) {
	dirname := allocTempDir(t)
	defer os.RemoveAll(dirname)
	_, err := blockfs.NewFS(dirname)
	if err != nil {
		t.Error("unexpected error:", err)
		t.Fail()
	}
	_, err = os.Stat(filepath.Join(dirname, "blocks"))
	if err != nil {
		t.Error("unexpected error:", err)
	}
}

func TestWriteBlock(t *testing.T) {
	dirname := allocTempDir(t)
	defer os.RemoveAll(dirname)
	fs, err := blockfs.NewFS(dirname)
	if err != nil {
		t.Error("unexpected error:", err)
		t.Fail()
	}
	data := []byte("test-data")
	block := blockfs.NewBlock(data)
	err = fs.WriteBlock(block)
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
	dirname := allocTempDir(t)
	defer os.RemoveAll(dirname)
	fs, err := blockfs.NewFS(dirname)
	if err != nil {
		t.Error("unexpected error:", err)
		t.Fail()
	}

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
