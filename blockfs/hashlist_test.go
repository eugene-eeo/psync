package blockfs_test

import (
	"github.com/eugene-eeo/psync/blockfs"
	"testing"
	"bytes"
	"strings"
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

func TestNewHashList(t *testing.T) {
	hashes := []string{
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d",
		"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6",
	}
	data := strings.Join(hashes, "\n")
	buff := bytes.NewBuffer([]byte(data))
	hl, err := blockfs.NewHashList(buff)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}
	for i, checksum := range hl {
		if string(checksum) != hashes[i] {
			t.Error(
				"expected hashes to equal:",
				checksum,
				hashes[i],
			)
		}
	}
}
