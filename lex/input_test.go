package lex

import (
	"testing"

	"github.com/DemoHn/Zn/error"
)

func TestSource_AddSourceInput(t *testing.T) {
	type args struct {
		rawData []byte
	}

	source := Source{
		Inputs: []SourceInput{},
	}
	tests := []struct {
		name string
		args args
		want *error.Error
	}{
		{
			name: "normal text",
			args: args{
				rawData: []byte("Hello, 世界！"),
			},
			want: nil,
		},
		{
			name: "normal - only 0",
			args: args{
				rawData: []byte{0},
			},
			want: nil,
		},
		{
			name: "empty text",
			args: args{
				rawData: []byte{},
			},
			want: nil,
		},
		{
			name: "wrong utf-8 sequence",
			args: args{
				rawData: []byte{0xD0, 0x81, 0x81},
			},
			want: error.NewError(0x1001),
		},
		{
			name: "wrong utf-8 sequence #2",
			args: args{
				rawData: []byte{0x81, 0x81},
			},
			want: error.NewError(0x1001),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := source.AddSourceInput(tt.args.rawData)
			if got == nil && tt.want != nil {
				t.Errorf("Source.AddSourceInput() = %v, want error: code(%v)", got, tt.want.GetCode())
			}
			if got != nil && tt.want == nil {
				t.Errorf("Source.AddSourceInput() = %v, want nil", got)
			}
			if got != nil && tt.want != nil && got.GetCode() != tt.want.GetCode() {
				t.Errorf("Source.AddSourceInput() = code(%v), want code(%v)", got.GetCode(), tt.want.GetCode())
			}
		})
	}
}
