package io

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"io"
	"unicode/utf8"
)

// InputStream defines an abstract reader that read bytes from external sources
// (e.g. files, strings, etc.) and transform to unicode chars
type InputStream interface {
	Read(n int) ([]rune, error)
	ReadAll() ([]rune, error)
}

// readRune - read bytes and yield runes
func readRune(r io.Reader, remains []byte, b int) ([]rune, []byte, error) {
	p := make([]byte, b)
	rs := make([]rune, 0)

	t, err := r.Read(p)
	if err != nil && err != io.EOF {
		return rs, []byte{}, zerr.ReadFileError(err, " <buffer> ")
	}

	buf := append(remains, p[:t]...)
	for len(buf) > 0 {
		ru, size := utf8.DecodeRune(buf)
		if ru == utf8.RuneError {
			return rs, buf, nil
		}

		rs = append(rs, ru)
		buf = buf[size:]
	}
	return rs, buf, nil
}