package blockfs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const BlocksDir string = "blocks"
const BlockSize int = 1024 * 1024 * 2

type FS struct {
	Path string
}

func NewFS(path string) *FS {
	os.MkdirAll(filepath.Join(path, BlocksDir), 0755)
	return &FS{
		Path: path,
	}
}

func (fs *FS) WriteBlock(b *Block) error {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	_, err = b.WriteTo(tmp)
	if err != nil {
		return err
	}
	path := filepath.Join(fs.Path, BlocksDir, string(b.Checksum))
	return os.Link(
		tmp.Name(),
		path,
	)
}

func (fs *FS) Export(r io.Reader) (HashList, error) {
	hashes := HashList{}
	buffer := make([]byte, BlockSize)
	for {
		length, err := r.Read(buffer)
		if length == 0 {
			break
		}
		if err != nil && err != io.EOF {
			return hashes, err
		}
		b := NewBlock(buffer[:length])
		fs.WriteBlock(b)
		hashes = append(hashes, b.Checksum)
	}
	return hashes, nil
}

func (fs *FS) GetBlock(c Checksum) (*Block, error) {
	f, err := os.Open(filepath.Join(fs.Path, BlocksDir, string(c)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := make([]byte, BlockSize)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	return NewBlock(b), nil
}
