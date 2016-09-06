package blockfs_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/eugene-eeo/psync/blockfs"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestChecksum(t *testing.T) {
	data := []byte("abc")
	c := blockfs.NewChecksum(data)
	b := sha256.Sum256(data)
	h := hex.EncodeToString(b[:])
	if h != string(c) {
		t.Error("expected NewChecksum('abc') ==", h, "got", c)
	}
}

type Params struct {
	Data  []byte
	Block *blockfs.Block
}

func (p *Params) Generate(r *rand.Rand, size int) reflect.Value {
	buff := make([]byte, r.Intn(1000))
	r.Read(buff)
	block := blockfs.NewBlock(buff)
	return reflect.ValueOf(&Params{
		Block: block,
		Data:  buff,
	})
}

func TestBlockChecksum(t *testing.T) {
	assertion := func(p *Params) bool {
		return len(p.Block.Checksum) == 64 &&
			bytes.Equal(p.Data, p.Block.Data)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}

func TestBlockWriteTo(t *testing.T) {
	assertion := func(p *Params) bool {
		b := bytes.NewBuffer([]byte{})
		p.Block.WriteTo(b)
		return bytes.Equal(b.Bytes(), p.Data)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}
