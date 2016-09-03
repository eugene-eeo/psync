package lib

import (
	"io"
)

func Chunked(r io.Reader, fn func([]byte, error)) {
	for {
		buffer := make([]byte, BLOCK_SIZE)
		n, err := r.Read(buffer)
		if n == 0 {
			break
		}
		fn(buffer, err)
	}
}
