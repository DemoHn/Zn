package lex

import (
	"reflect"
	"testing"
)

// mainly for testing parseCommentHead()
func TestNextToken_CommentsONLY(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expectError bool
		token       Token
		lineInfo    string
	}{
		{
			name:        "singleLine comment",
			input:       "注：这是一个长 长 的单行注释comment",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("这是一个长 长 的单行注释comment"),
				Info: map[string]bool{
					"isMultiLine": false,
				},
			},
			lineInfo: "Unknown<0>[0,21]",
		},
		{
			name: "singleLine comment with newline",
			input: "注：注：注：nach nach\r\n注：又是一个注",
			expectError: false,
			token: Token{
				Type: TokenComment,
				Literal: []rune("注：注：nach nach"),
				Info: map[string]bool{
					"isMultiLine": false,
				},
			},
			lineInfo: "Unknown<0>[0,14]",
		}
	}

	for _, tt := range cases {
		lex := NewLexer([]rune(tt.input))
		t.Run(tt.name, func(t *testing.T) {
			tk, err := lex.NextToken()
			// validate error
			if tt.expectError == false {
				if err != nil {
					t.Errorf("NextToken() failed! expected no error, but got error")
					t.Error(err)
				} else if tt.lineInfo != lex.lines.String() {
					t.Errorf("NextToken() lineInfo expect `%s`, actual `%s`", tt.lineInfo, lex.lines.String())
				}
			} else {
				if err == nil {
					t.Errorf("NextToken() failed! expected error, but got no error")
				}
			}

			if !reflect.DeepEqual(*tk, tt.token) {
				t.Errorf("NextToken() return Token failed! expect: %v, got: %v", tt.token, *tk)
			}
		})
	}
}

//// helpers
/**
// test indents
func TestLex_Tokenize(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
		lineInfo    string
	}{
		{
			name:        "empty input",
			source:      "",
			expectError: false,
			lineInfo:    "Empty<0>",
		},
		{
			name:        "all CRs & LFs",
			source:      "\r\n\r\n\n\n",
			expectError: false,
			lineInfo:    "Empty<0> Empty<0> Empty<0> Empty<0> Empty<0>",
		},
		{
			name:        "no indent",
			source:      "all is the char",
			expectError: false,
			lineInfo:    "Unknown<0>[0,14]",
		},
		{
			name:        "2 lines with no indent",
			source:      "line-1\r\nline-2",
			expectError: false,
			lineInfo:    "Unknown<0>[0,5] Unknown<0>[8,13]",
		},
		{
			name:        "5 lines with no indent",
			source:      "line-1\r\nline-2\nline-3n\rline-4r\n\rline-5",
			expectError: false,
			lineInfo:    "Unknown<0>[0,5] Unknown<0>[8,13] Unknown<0>[15,21] Unknown<0>[23,29] Unknown<0>[32,37]",
		},
		{
			name:        "with space indents",
			source:      "line1\r\n    line2",
			expectError: false,
			lineInfo:    "Space<0>[0,4] Space<1>[11,15]",
		},
		{
			name:        "with tab indents",
			source:      "line1\n\t\tline2",
			expectError: false,
			lineInfo:    "Tab<0>[0,4] Tab<2>[8,12]",
		},
		{
			name:        "multi lines with tab indents",
			source:      "line1\n\t\tline2\n\t\t\tline\t3",
			expectError: false,
			lineInfo:    "Tab<0>[0,4] Tab<2>[8,12] Tab<3>[17,22]",
		},
		{
			name:        "with non-null-indent empty line",
			source:      "line1\n        \nline3    ",
			expectError: false,
			lineInfo:    "Space<0>[0,4] Empty<0> Space<0>[15,23]",
		},
		{
			name:        "incorrect space nums: 3",
			source:      "line1\n   \nline3    ",
			expectError: true,
			lineInfo:    "",
		},
		{
			name:        "mixture of spaces & tabs",
			source:      "line1\n    \n\t\t\thello",
			expectError: true,
			lineInfo:    "",
		},
	}

	for _, tt := range tests {
		lex := NewLexer([]rune(tt.source))
		t.Run(tt.name, func(t *testing.T) {
			got := lex.Tokenize()

			if tt.expectError == false && got != nil {
				t.Errorf("Tokenize() failed! expected no error, but got error")
				t.Error(got)
				return
			}

			if tt.expectError == true && got == nil {
				t.Errorf("Tokenize() failed! expected error, but got no error")
			}

			if tt.expectError == false && tt.lineInfo != lex.lineScan.String() {
				t.Errorf("Tokenize() lineInfo expect `%s`, actual `%s`", tt.lineInfo, lex.lineScan.String())
				return
			}
		})
	}
}
*/
