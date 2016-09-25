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

func NewFS(path string) (*FS, error) {
	os.Mkdir(path, 0755)
	err := os.Mkdir(filepath.Join(path, BlocksDir), 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	fs := FS{Path: path}
	return &fs, nil
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
	err = os.Link(tmp.Name(), path)
	if os.IsExist(err) {
		return nil
	}
	return err
}

func (fs *FS) Export(r io.Reader) (HashList, error) {
	hashes := HashList{}
	buffer := make([]byte, BlockSize)
	for {
		length, err := io.ReadFull(r, buffer)
		if length == 0 {
			break
		}
		b := NewBlock(buffer[:length])
		fs.WriteBlock(b)
		hashes = append(hashes, b.Checksum)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return hashes, err
		}
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
	length, err := f.Read(b)
	if err != nil {
		return nil, err
	}
	return NewBlock(b[:length]), nil
}

func (fs *FS) Exists(c Checksum) bool {
	_, err := os.Stat(filepath.Join(fs.Path, BlocksDir, string(c)))
	if err != nil {
		return false
	}
	return true
}

func (fs *FS) MissingBlocks(h HashList) HashList {
	missing := HashList{}
	for _, checksum := range h {
		if !fs.Exists(checksum) {
			missing = append(missing, checksum)
		}
	}
	return missing
}
