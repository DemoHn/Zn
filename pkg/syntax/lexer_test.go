package syntax

import (
	"testing"
)

func TestLexer_ParseLine(t *testing.T) {
	cases := []struct {
		name        string
		text        []rune
		expectLines []LineInfo
		err         error
	}{
		{
			name: "one-line string",
			text: []rune("        This is one line233"),
			expectLines: []LineInfo{
				{
					Indents:  0,
					StartIdx: 0,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.text)
			_ = l.parseBeginLex()
		H:
			for {
				ch := l.Next()
				switch ch {
				case RuneEOF:
					break H
				case RuneCR, RuneLF:
					e := l.parseLine(ch, true)
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

func TestLexer_FindLineIdx(t *testing.T) {
	cases := []struct {
		name         string
		cursor       int
		startLoopIdx int
		lines        []LineInfo
		expectedIdx  int
	}{
		{
			name:         "normal lines (final line)",
			cursor:       10,
			startLoopIdx: 0,
			lines: []LineInfo{
				{
					StartIdx: 0,
				},
				{
					StartIdx: 4,
				},
				{
					StartIdx: 7,
				},
			},
			expectedIdx: 2,
		},
		{
			name:         "normal lines (first line)",
			cursor:       3,
			startLoopIdx: 0,
			lines: []LineInfo{
				{
					StartIdx: 0,
				},
				{
					StartIdx: 4,
				},
				{
					StartIdx: 7,
				},
			},
			expectedIdx: 0,
		},
		{
			name:         "normal lines with same startIdx (rare case)",
			cursor:       9,
			startLoopIdx: 0,
			lines: []LineInfo{
				{
					StartIdx: 0,
				},
				{
					StartIdx: 4,
				},
				{
					StartIdx: 4,
				},
				{
					StartIdx: 7,
				},
				{
					StartIdx: 11,
				},
			},
			expectedIdx: 3,
		},
		{
			name:         "empty lines",
			cursor:       0,
			startLoopIdx: 0,
			lines:        []LineInfo{},
			expectedIdx:  0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			lex := &Lexer{
				Lines: tt.lines,
			}

			got := lex.FindLineIdx(tt.cursor, tt.startLoopIdx)
			if got != tt.expectedIdx {
				t.Errorf("expect %d, got %d", tt.expectedIdx, got)
			}
		})
	}
}
