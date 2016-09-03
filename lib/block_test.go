package lib_test

import (
	"testing"
	"testing/quick"
	"math/rand"
	"reflect"
	"github.com/eugene-eeo/psync/lib"
	"bytes"
)

type BlockTest struct {
	Block *lib.Block
	Data  []byte
}

func (b *BlockTest) Generate(r *rand.Rand, size int) reflect.Value {
	buff := make([]byte, r.Intn(1000))
	r.Read(buff)
	block := lib.NewBlock(buff)
	return reflect.ValueOf(&BlockTest{
		Block: block,
		Data:  buff,
	})
}

func TestBlockChecksum(t *testing.T) {
	assertion := func (b *BlockTest) bool {
		return reflect.DeepEqual(b.Data, b.Block.Data) && len(b.Block.Checksum) == 64
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}

func TestBlockWriteTo(t *testing.T) {
	assertion := func (b *BlockTest) bool {
		buff := bytes.NewBuffer([]byte{})
		written, err := b.Block.WriteTo(buff)
		if err != nil {
			t.Log(err)
			return false
		}
		if written != int64(len(b.Block.Data)) {
			return false
		}
		return reflect.DeepEqual(
			buff.Bytes(),
			b.Block.Data,
		)
	}
	if err := quick.Check(assertion, nil); err != nil {
		t.Error(err)
	}
}
