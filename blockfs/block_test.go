package blockfs_test

import (
	"math/rand"
	"github.com/eugene-eeo/psync/blockfs"
	"testing"
	"testing/quick"
	"reflect"
	"bytes"
	"fmt"
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

func TestBlockStream(t *testing.T) {
	buff := make(
		[]byte,
		blockfs.BLOCK_SIZE * 2 + blockfs.BLOCK_SIZE >> 2,
	)
	rand.Read(buff)
	stream := blockfs.BlockStream(bytes.NewBuffer(buff))
	blocks := 0
	for {
		block, err := stream()
		if block == nil {
			break
		}
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		blocks++
		start := (blocks - 1) * blockfs.BLOCK_SIZE
		end := start + blockfs.BLOCK_SIZE
		if blocks == 3 {
			end = start + blockfs.BLOCK_SIZE >> 2
		}
		if !bytes.Equal(block.Data, buff[start:end]) {
			t.Error("expected blocks to equal")
			t.Fail()
		}
	}
	if blocks != 3 {
		t.Error("expected 3 blocks, got", blocks)
	}
}

type ErrorStream struct {
	Error error
}

func (e *ErrorStream) Read(data []byte) (int, error) {
	data[0] = 1
	return 1, e.Error
}

func TestBlockStreamError(t *testing.T) {
	errstream := ErrorStream{
		Error: fmt.Errorf("given"),
	}
	stream := blockfs.BlockStream(&errstream)
	block, err := stream()
	if block.Data[0] != 1 {
		t.Error("expected first byte to be written")
	}
	if err == nil {
		t.Error("expected error to equal", errstream.Error)
	}
}
