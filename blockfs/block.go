package blockfs

import (
	"io"
	"crypto/sha256"
	"encoding/hex"
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
		Data: data,
	}
}

func (b *Block) WriteTo(w io.Writer) (int64, error) {
	wrote, err := w.Write(b.Data)
	return int64(wrote), err
}

func BlockStream(r io.Reader) func() (*Block, error) {
	buffer := make([]byte, BLOCK_SIZE)
	return func() (*Block, error) {
		length, err := r.Read(buffer)
		if length == 0 {
			return nil, nil
		}
		b := NewBlock(buffer[:length])
		return b, err
	}
}
