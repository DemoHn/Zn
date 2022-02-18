package syntax

import (
	"testing"
)

func TestLexer_ParseLine(t *testing.T) {
	cases := []struct {
		name string
		text []rune
		expectLines []LineInfo
		err error
	}{
		{
			name: "one-line string",
			text: []rune("This is one line\r\n233"),
			expectLines: []LineInfo{
				{
					Indents: 0,
					StartIdx: 0,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func (t *testing.T) {
			l := NewLexer(tt.text)
			H:
			for {
				ch := l.Next()
				switch ch {
				case RuneEOF:
					break H
				case RuneCR, RuneLF:
					e := l.ParseLine(ch)
					if e == nil && tt.err != nil {
						t.Errorf("expect error, finally no error")
					} else if e != nil && tt.err == nil {
						t.Errorf("expect no error, finally meet error: %s", e)
					}
				}
			}

			// after parsing all lines, it's time to assert lines
			if len(l.Lines) != len(tt.expectLines) {
				t.Errorf("expect %d lines, actual %d lines", len(tt.expectLines), len(l.Lines))
			}

			// check line info one-by-one
			t.Logf("%v\n", l.Lines)
		})
	}
}