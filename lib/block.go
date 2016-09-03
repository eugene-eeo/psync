package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

const BLOCK_SIZE = 8192

type Checksum string

func NewChecksum(data []byte) Checksum {
	b := sha256.Sum256(data)
	return Checksum(hex.EncodeToString(b[:]))
}

type Block struct {
	Checksum Checksum
	Data []byte
}

func NewBlock(data []byte) *Block {
	return &Block{
		Checksum: NewChecksum(data),
		Data: data,
	}
}

func (block *Block) WriteTo(w io.Writer) (int64, error) {
	b, err := w.Write(block.Data)
	return int64(b), err
}
