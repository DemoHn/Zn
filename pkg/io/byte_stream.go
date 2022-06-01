package io

import (
	"bytes"
	"io"
)

// ByteStream - import a string as code source
type ByteStream struct {
	reader io.Reader
	length int
	encBuffer []byte
}

// NewByteStream - new text stream
func NewByteStream(b []byte) *ByteStream {
	return &ByteStream{
		reader: bytes.NewReader(b),
		length: len(b),
		encBuffer: []byte{},
	}
}

func (b *ByteStream) ReadAll() ([]rune, error) {
	data, _, err := readRune(b.reader, b.encBuffer, b.length)
	if err != nil {
		return []rune{}, err
	}
	return data, nil
}

func (b *ByteStream) Read(n int) ([]rune, error) {
	data, remains, err := readRune(b.reader, b.encBuffer, n)
	if err != nil {
		return []rune{}, err
	}
	b.encBuffer = remains
	return data, nil
}