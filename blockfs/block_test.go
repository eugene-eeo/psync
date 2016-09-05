package blockfs_test

import (
	"math/rand"
	"github.com/eugene-eeo/psync/blockfs"
	"testing"
	"testing/quick"
	"reflect"
	"bytes"
)

type Params struct {
	Data []byte
	Block *blockfs.Block
}

func (p *Params) Generate(r *rand.Rand, size int) reflect.Value {
	buff := make([]byte, r.Intn(1000))
	r.Read(buff)
	block := blockfs.NewBlock(buff)
	return reflect.ValueOf(&Params{
		Block: block,
		Data: buff,
	})
}

func TestBlockChecksum(t *testing.T) {
	assertion := func (p *Params) bool {
		return len(p.Block.Checksum) == 64 &&
			   bytes.Equal(p.Data, p.Block.Data)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}

func TestBlockWriteTo(t *testing.T) {
	assertion := func (p *Params) bool {
		b := bytes.NewBuffer([]byte{})
		p.Block.WriteTo(b)
		return bytes.Equal(b.Bytes(), p.Data)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}
