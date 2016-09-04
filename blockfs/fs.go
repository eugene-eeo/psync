package blockfs

import (
	"io"
	"os"
	"path/filepath"
)

const TAGS_DIR   string = "tags"
const BLOCKS_DIR string = "blocks"
const BLOCK_SIZE int = 1024 * 1024

type FS struct {
	Path string
}

func NewFS(path string) *FS {
	os.MkdirAll(filepath.Join(path, BLOCKS_DIR), 0755)
	os.MkdirAll(filepath.Join(path, TAGS_DIR), 0755)
	return &FS{
		Path: path,
	}
}

func (fs *FS) WriteBlock(b *Block) error {
	f, err := os.Create(filepath.Join(fs.Path, BLOCKS_DIR, string(b.Checksum)))
	if err != nil {
		return err
	}
	_, err = b.WriteTo(f)
	return err
}

func (fs *FS) Export(r io.Reader) (*HashList, error) {
	hashes := HashList{}
	buffer := make([]byte, BLOCK_SIZE)
	for {
		bits, err := r.Read(buffer)
		if bits == 0 {
			break
		}
		if err != nil && err != io.EOF {
			return nil, err
		}
		b := NewBlock(buffer[:bits])
		fs.WriteBlock(b)
		hashes = append(hashes, b.Checksum)
	}
	return &hashes, nil
}

func (fs *FS) ExportNamed(r io.Reader, name string) (*HashList, error) {
	f, err := os.Create(filepath.Join(fs.Path, TAGS_DIR, name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	hl, err := fs.Export(r)	
	hl.WriteTo(f)
	return hl, err
}

func (fs *FS) GetBlock(c Checksum) (*Block, error) {
	f, err := os.Open(filepath.Join(fs.Path, BLOCKS_DIR, string(c)))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := make([]byte, BLOCK_SIZE)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	return NewBlock(b), nil
}
