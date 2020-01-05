package lex

import (
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
		{
			name:        "comment 2 lines, one string",
			input:       "注：“可是都 \n  不为空”“是为”《淮南子》",
			expectError: false,
			tokens:      "$10[可是都 \n  不为空] $90[是为] $90[淮南子]",
		},
		{
			name:        "nest multiple strings",
			input:       "·显然在其中·“不为空”‘为\n\n空’「「「随意“嵌套”」233」456」",
			expectError: false,
			tokens:      "$91[显然在其中] $90[不为空] $90[为\n\n空] $90[「「随意“嵌套”」233」456]",
		},
		{
			name:        "incomplete var quote at end",
			input:       "如何·显然在其中",
			expectError: false,
			tokens:      "$45[如何] $91[显然在其中]",
		},
		{
			name:        "consecutive keywords",
			input:       "以其为",
			expectError: false,
			tokens:      "$56[以] $65[其] $41[为]",
		},
		{
			name:        "consecutive keywords #2",
			input:       "不以其为",
			expectError: false,
			tokens:      "$101[不] $56[以] $65[其] $41[为]",
		},
		{
			name:        "multi line string with var quote inside",
			input:       "“搞\n个\n    大新闻”《·焦点在哪里·》\n\t注：“又是一年\n    春来到”",
			expectError: false,
			tokens:      "$90[搞\n个\n    大新闻] $90[·焦点在哪里·] $10[又是一年\n    春来到]",
		},
		{
			name:        "markers with spaces",
			input:       "\n    （  ） ， A/B  #  25",
			expectError: false,
			tokens:      "$22[（] $23[）] $11[，] $101[A/B] $18[#] $100[25]",
		},
		{
			name:        "keyword after line",
			input:       "令甲，乙为（【12，34，【“测试到底”，10】】）\n令丙为“23”",
			expectError: false,
			tokens: "$40[令] $101[甲] $11[，] $101[乙] $41[为] $22[（] $20[【] $100[12]" +
				" $11[，] $100[34] $11[，] $20[【] $90[测试到底] $11[，] $100[10] $21[】]" +
				" $21[】] $23[）] $40[令] $101[丙] $41[为] $90[23]",
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
				if tk.Type == TypeEOF {
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
				var actualStr = StringifyAllTokens(tokens)
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
