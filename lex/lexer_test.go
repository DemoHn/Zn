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
			name:        "singleLine empty comment",
			input:       "注：",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune(""),
				Info: map[string]bool{
					"isMultiLine": false,
				},
			},
			lineInfo: "Unknown<0>[0,1]",
		},
		{
			name:        "singleLine empty comment (single quote)",
			input:       "注： “",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune(" “"),
				Info: map[string]bool{
					"isMultiLine": false,
				},
			},
			lineInfo: "Unknown<0>[0,3]",
		},
		{
			name:        "singleLine empty comment (with number)",
			input:       "注 1024 2048 ：",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune(""),
				Info: map[string]bool{
					"isMultiLine": false,
				},
			},
			lineInfo: "Unknown<0>[0,12]",
		},
		{
			name:        "singleLine comment with newline",
			input:       "注：注：注：nach nach\r\n注：又是一个注",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("注：注：nach nach\r\n"),
				Info: map[string]bool{
					"isMultiLine": false,
				},
			},
			lineInfo: "Unknown<0>[0,14]",
		},
		//// multi-line comment
		{
			name:        "mutlLine comment with no new line",
			input:       "注：“假设这是一个注” 后面假设又是一些数",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("“假设这是一个注”"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "",
		},
		{
			name:        "mutlLine comment with no other string",
			input:       "注：“假设这又是一个注”",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("“假设这又是一个注”"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "",
		},
		{
			name:        "mutlLine comment (with number)",
			input:       "注 1234 5678 ：“假设这又是一个注”",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("“假设这又是一个注”"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "",
		},
		{
			name:        "mutlLine comment with empty string",
			input:       "注：“”",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("“”"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "",
		},
		{
			name:        "mutlLine comment with multiple lines",
			input:       "注：“一一\r\n    二二\n三三\n四四”",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("“一一\r\n    二二\n三三\n四四”"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "Unknown<0>[0,4] Unknown<0>[6,12] Unknown<0>[13,15]",
		},
		{
			name:        "mutlLine comment with quote stack",
			input:       "注：“一一「2233」《某本书》注：“”二二\n     ”",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("“一一「2233」《某本书》注：“”二二\n     ”"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "Unknown<0>[0,21]",
		},
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
