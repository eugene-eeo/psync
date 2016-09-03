package lib

import (
	"io"
	"bufio"
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

func ParseHashList(r io.Reader, fn func(Checksum)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		fn(Checksum(line))
	}
}
