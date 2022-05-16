package io

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"io"
	"os"
)

type FileStream struct {
	reader io.Reader
	encBuffer []byte
	path string
	hasRead bool
}


const (
	defaultReadBlock = 4096
	BOM = 0xFEFF
)

// NewFileStream - create file stream
func NewFileStream(path string) (*FileStream, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, zerr.FileNotFound(path)
	}

	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &FileStream{
		reader:    reader,
		encBuffer: []byte{},
		path:      path,
		hasRead:   false,
	}, nil
}

func (f *FileStream) ReadAll() ([]rune, error) {
	var result []rune
	for {
		res, err := f.read(defaultReadBlock)
		if err != nil {
			return []rune{}, err
		}

		if len(res) == 0 {
			break
		}
		result = append(result, res...)
 	}

 	return result, nil
}

// Read - read some chars
func (f *FileStream) Read(n int) ([]rune, error) {
	return f.read(n)
}

// GetPath -
func (f *FileStream) GetPath() string {
	return f.path
}

func (f *FileStream) read(n int) ([]rune, error) {
	data, remains, err := readRune(f.reader, f.encBuffer, n)
	if err != nil {
		return []rune{}, err
	}
	f.encBuffer = remains

	if !f.hasRead {
		f.hasRead = true
		// detect BOM, if BOM on the first char, then remove it directly.
		if data[0] == BOM {
			data = data[1:]
		}
	}
	return data, nil
}