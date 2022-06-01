package io

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func setup() (string, error) {
	tmpdir, err := ioutil.TempDir("", "tmp_file_stream")
	if err != nil {
		return "", err
	}
	fileName := path.Join(tmpdir, "tmp.txt")

	f, err := os.Create(fileName)
	if f != nil {
		defer f.Close()
	}

	return fileName, err
}

func teardown(file string) {
	_ = os.Remove(file)
}

func TestFileStream_ReadAll(t *testing.T) {
	cases := []struct{
		name string
		data []byte
		assertD []rune
		assertE error
	}{
		{
			name: "normal chars",
			data: []byte("ABCD123456"),
			assertD: []rune("ABCD123456"),
			assertE: nil,
		},
		{
			name: "large blocks",
			data: bytes.Repeat([]byte{0xE7, 0x8C, 0xAA, 0xE5, 0xA4, 0xB4}, 10284),
			assertD: []rune(strings.Repeat("猪头", 10284)),
			assertE: nil,
		},
		{
			name: "file starts with BOM",
			data: []byte{0xEF, 0xBB, 0xBF, 0xE7, 0x8C, 0xAA, 0xE5, 0xA4, 0xB4},
			assertD: []rune("猪头"),
			assertE: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// setup and teardown
			file, _ := setup()
			defer teardown(file)
			// write file
			_ = ioutil.WriteFile(file, tt.data, 0644)

			s, _ := NewFileStream(file)

			expect, err := s.ReadAll()
			if tt.assertE == nil &&  err != nil {
				t.Fatalf("expect no error, but error occured: %s", err)
			} else if tt.assertE != nil && err == nil {
				t.Fatalf("expect error, but NO error occured")
			} else {
				if strings.Compare(string(tt.assertD), string(expect)) != 0 {
					t.Fatalf("data output not match!")
				}
			}
		})
	}
}
