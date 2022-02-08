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


const defaultReadBlock = 1024

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
		hasRead: false,
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

func (f *FileStream) read(n int) ([]rune, error) {
	data, remains, err := readRune(f.reader, f.encBuffer, n)
	if err != nil {
		return []rune{}, err
	}
	f.encBuffer = remains
	// detect BOM
	if !f.hasRead {
		f.hasRead = true
		if data[0] == 0xFEFF {
			data = data[1:]
		}
	}
	return data, nil
}