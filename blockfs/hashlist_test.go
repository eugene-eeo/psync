package blockfs_test

import (
	"github.com/eugene-eeo/psync/blockfs"
	"testing"
	"bytes"
)

func TestHashListWriteTo(t *testing.T) {
	hl := blockfs.HashList{
		blockfs.Checksum("abc"),
		blockfs.Checksum("def"),
		blockfs.Checksum("ghi"),
	}
	b := bytes.NewBuffer([]byte{})
	total, err := hl.WriteTo(b)
	if err != nil {
		t.Error("expected err to be nil, got", err)
	}
	if total != (3+1)*3 {
		t.Error("expected 12 bytes to be written, got", total)
	}
	given := b.Bytes()
	expected := []byte("abc\ndef\nghi\n")
	if !bytes.Equal(expected, given) {
		t.Error("expected written contents to equal", expected, "got", given)
	}
}
