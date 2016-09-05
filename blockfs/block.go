package blockfs

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type Checksum string
type Block struct {
	Checksum Checksum
	Data     []byte
}

func NewChecksum(data []byte) Checksum {
	b := sha256.Sum256(data)
	return Checksum(hex.EncodeToString(b[:]))
}

func NewBlock(data []byte) *Block {
	return &Block{
		Checksum: NewChecksum(data),
		Data:     data,
	}
}

func (b *Block) WriteTo(w io.Writer) (int64, error) {
	wrote, err := w.Write(b.Data)
	return int64(wrote), err
}
