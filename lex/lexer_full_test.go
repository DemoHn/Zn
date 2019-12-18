package lex

import (
	"fmt"
	"strings"
	"testing"
)

type tokensCase struct {
	name        string
	input       string
	expectError bool
	tokens      string
}

// stringify token grammer:
// $type[literal] $type2[literal]
//
// example:
// $0[这是一个长长的单行注释] $1[引用一个文本]
func TestNextToken_MixedText(t *testing.T) {
	cases := []tokensCase{
		{
			name:        "1 number, 1 identifier",
			input:       `12.5rpm`,
			expectError: false,
			tokens:      `$100[12.5] $101[rpm]`,
		},
		{
			name:        "1 identifier with 1 inline comment",
			input:       `标识符名注：这是一个标识符啊 `,
			expectError: false,
			tokens:      `$101[标识符名] $10[这是一个标识符啊 ]`,
		},
		{
			name:        "1 identifier (mixed number) with 1 inline comment",
			input:       `标识符名12注：这是一个标识符啊 `,
			expectError: false,
			tokens:      `$101[标识符名12] $10[这是一个标识符啊 ]`,
		},
		{
			name:        "1 identifier, 注 is not comment",
			input:       `起居注23不为其和`,
			expectError: false,
			tokens:      `$101[起居注23] $49[不为] $65[其] $101[和]`,
		},
		{
			name:        "identifer in keyword",
			input:       `令变量不为空`,
			expectError: false,
			tokens:      `$40[令] $101[变量] $49[不为] $101[空]`,
		},
		{
			name:        "1 identifier sep 1 number",
			input:       `变量1为12.45E+3`,
			expectError: false,
			tokens:      `$101[变量1] $41[为] $100[12.45E+3]`,
		},
	}

	assertTokens(cases, t)
}

func assertTokens(cases []tokensCase, t *testing.T) {
	for _, tt := range cases {
		lex := NewLexer([]rune(tt.input))
		t.Run(tt.name, func(t *testing.T) {
			var tErr error
			var tokens = make([]*Token, 0)
			// iterate to get tokens
			for {
				tk, err := lex.NextToken()
				if err != nil {
					tErr = err
					break
				}
				if tk.Type == typeEOF {
					break
				}
				tokens = append(tokens, tk)
			}
			// assert data
			if tt.expectError == false {
				if tErr != nil {
					t.Errorf("parse Tokens failed! expected no error, but got error")
					t.Error(tErr)
					return
				}

				// conform all tokens to string
				var tokenStrs = []string{}
				for _, ptk := range tokens {
					tokenStrs = append(tokenStrs, fmt.Sprintf("$%d[%s]", ptk.Type, string(ptk.Literal)))
				}

				var actualStr = strings.Join(tokenStrs, " ")
				if actualStr != tt.tokens {
					t.Errorf("tokens not same! \nexpect->\n%s\ngot->\n%s", tt.tokens, actualStr)
				}

			} else {
				if tErr == nil {
					t.Errorf("NextToken() failed! expected error, but got no error")
				}
			}

		})
	}
}
