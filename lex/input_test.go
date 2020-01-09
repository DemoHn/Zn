package lex

import (
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/error"
)

func TestInputStream_BasicUTF8Parsing(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want *error.Error
	}{
		{
			name: "normal text",
			args: []byte("Hello, 世界！"),
			want: nil,
		},
		{
			name: "normal - the \\0 char",
			args: []byte{0},
			want: nil,
		},
		{
			name: "empty text",
			args: []byte{},
			want: nil,
		},
		{
			name: "wrong utf-8 sequence",
			args: []byte{0xD0, 0x81, 0x81},
			want: error.NewError(0x1001),
		},
		{
			name: "wrong utf-8 sequence #2",
			args: []byte{0x81, 0x81},
			want: error.NewError(0x1001),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := NewBufferStream(tt.args)
			var got *error.Error
			for !stream.End() {
				_, got = stream.Read(256)
				if got != nil {
					break
				}
			}

			if got == nil && tt.want != nil {
				t.Errorf("stream.Read() = %v, want error: code(%v)", got, tt.want.GetCode())
			}
			if got != nil && tt.want == nil {
				t.Errorf("stream.Read() = %v, want nil", got)
			}
			if got != nil && tt.want != nil && got.GetCode() != tt.want.GetCode() {
				t.Errorf("stream.Read() = code(%v), want code(%v)", got.GetCode(), tt.want.GetCode())
			}
		})
	}
}

func TestInputStream_MultiTimeUTF8Parsing(t *testing.T) {
	// read char block size
	tests := []struct {
		name      string
		args      []byte
		blockSize int
		err       *error.Error
		runeList  []string
	}{
		{
			name:      "normal char sequence",
			args:      []byte("千里之行233"),
			blockSize: 4,
			err:       nil,
			runeList: []string{
				"千", "里", "之行", "233", "",
			},
		},
		{
			name:      "normal char sequence /read in 1 time",
			args:      []byte("千里之行"),
			blockSize: 256,
			runeList: []string{
				"千里之行", "",
			},
		},
		{
			name:      "normal char sequence /read by 1 byte",
			args:      []byte("A测试"),
			blockSize: 1,
			runeList: []string{
				"A", "", "", "测", "", "", "试", "",
			},
		},
		{
			name:      "normal char sequence /perfectly by 3 bytes",
			args:      []byte("始于足下"),
			blockSize: 3,
			runeList: []string{
				"始", "于", "足", "下", "",
			},
		},
		// fail cases
		{
			name:      "non UTF-8 sequences /from begin to end",
			args:      []byte{0x9B, 0x03, 0x20, 0x83},
			blockSize: 6,
			runeList:  []string{""},
			err:       error.NewError(0x1001),
		},
		{
			name:      "non UTF-8 sequences in the middle",
			args:      append([]byte("千里"), []byte{0xFC, 0x81, 0x81, 0x20}...),
			blockSize: 2,
			runeList:  []string{"", "千", "里", "", ""},
			err:       error.NewError(0x1001),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := NewBufferStream(tt.args)
			var e *error.Error
			var data = []rune{}
			var dataList = []string{}

			for !stream.End() {
				data, e = stream.Read(tt.blockSize)
				if e != nil {
					break
				}
				dataList = append(dataList, string(data))
			}

			// about error
			if e == nil && tt.err != nil {
				t.Errorf("stream.Read() = %v, want error: code(%v)", e, tt.err.GetCode())
			}
			if e != nil && tt.err == nil {
				t.Errorf("stream.Read() = %v, want nil", e)
			}
			if e != nil && tt.err != nil && e.GetCode() != tt.err.GetCode() {
				t.Errorf("stream.Read() = code(%v), want code(%v)", e.GetCode(), tt.err.GetCode())
			}
			// about dataList
			if !reflect.DeepEqual(dataList, tt.runeList) {
				t.Errorf("dataList = %v, want = %v", dataList, tt.runeList)
			}
		})
	}
}
