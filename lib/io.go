package lib

import (
	"io"
	"bufio"
	"strings"
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
	src := bufio.NewReader(r)
	for {
		line, _ := src.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if len(line) == 0 {
			break
		}
		fn(Checksum(line))
	}
}
