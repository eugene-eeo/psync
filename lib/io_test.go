package lib_test

import (
	"github.com/eugene-eeo/psync/lib"
	"testing"
	"math/rand"
	"bytes"
	"reflect"
)

func TestChunked(t *testing.T) {
	b := make([]byte, lib.BLOCK_SIZE * 2.25)
	rand.Read(b)
	times := 0
	lib.Chunked(bytes.NewBuffer(b), func(s []byte, e error) {
		times++
		if e != nil {
			t.Error("unexpected error")
		}
		if len(s) != lib.BLOCK_SIZE {
			t.Error("expected chunk length to be BLOCK_SIZE")
		}
		start := lib.BLOCK_SIZE * (times-1)
		length := lib.BLOCK_SIZE
		if times == 3 {
			length = lib.BLOCK_SIZE * 0.25
		}
		end := start + length
		if !reflect.DeepEqual(b[start:end], s[:length]) {
			t.Error("expected chunks to equal")
		}
	})
	if times != 3 {
		t.Error("expected func to be called 3 times")
	}
}
