package lex

import (
	"reflect"
	"testing"
)

type nextTokenCase struct {
	name        string
	input       string
	expectError bool
	token       Token
	lineInfo    string
}

// mainly for testing parseCommentHead()
func TestNextToken_CommentsONLY(t *testing.T) {
	cases := []nextTokenCase{
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
				Literal: []rune("注：注：nach nach"),
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
				Literal: []rune("假设这是一个注"),
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
				Literal: []rune("假设这又是一个注"),
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
				Literal: []rune("假设这又是一个注"),
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
				Literal: []rune(""),
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
				Literal: []rune("一一\r\n    二二\n三三\n四四"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "Unknown<0>[0,4] Unknown<0>[7,12] Unknown<0>[14,15]",
		},
		{
			name:        "mutlLine comment with quote stack",
			input:       "注：“一一「2233」《某本书》注：“”二二\n     ”",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("一一「2233」《某本书》注：“”二二\n     "),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "Unknown<0>[0,21]",
		},
		{
			name:        "mutlLine comment with straight quote",
			input:       "注：「PK」",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("PK"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "",
		},
		{
			name:        "mutlLine comment unfinished quote",
			input:       "注：「PKG“”",
			expectError: false,
			token: Token{
				Type:    TokenComment,
				Literal: []rune("PKG“”"),
				Info: map[string]bool{
					"isMultiLine": true,
				},
			},
			lineInfo: "Unknown<0>[0,7]",
		},
	}
	assertNextToken(cases, t)
}

func TestNextToken_StringONLY(t *testing.T) {
	cases := []nextTokenCase{
		{
			name:        "normal quote string",
			input:       "“LSK” 多出来的",
			expectError: false,
			token: Token{
				Type:    TokenString,
				Literal: []rune("LSK"),
				Info:    '“',
			},
			lineInfo: "",
		},
		{
			name:        "normal quote string (with whitespaces)",
			input:       "“这 是 一 个 字 符 串”",
			expectError: false,
			token: Token{
				Type:    TokenString,
				Literal: []rune("这 是 一 个 字 符 串"),
				Info:    '“',
			},
			lineInfo: "",
		},
		{
			name:        "normal quote string (with multiple quotes)",
			input:       "“「233」 ‘456’ 《〈who〉》『『is』』”",
			expectError: false,
			token: Token{
				Type:    TokenString,
				Literal: []rune("「233」 ‘456’ 《〈who〉》『『is』』"),
				Info:    '“',
			},
			lineInfo: "",
		},
		{
			name:        "multiple-line string",
			input:       "『233\n456\r\n7  』",
			expectError: false,
			token: Token{
				Type:    TokenString,
				Literal: []rune("233\n456\r\n7  "),
				Info:    '『',
			},
			lineInfo: "",
		},
	}

	assertNextToken(cases, t)
}

func assertNextToken(cases []nextTokenCase, t *testing.T) {
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
