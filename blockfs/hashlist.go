package blockfs

import (
	"io"
)

type HashList []Checksum

func (h *HashList) WriteTo(w io.Writer) (int64, error) {
	t := []byte{}
	for _, checksum := range *h {
		b := []byte(checksum)
		t = append(t, b...)
		t = append(t, 10) // newline
	}
	total, err := w.Write(t)
	return int64(total), err
}
