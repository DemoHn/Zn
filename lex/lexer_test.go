package lex

import (
	"fmt"
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeComment,
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
				Type:    TypeString,
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
				Type:    TypeString,
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
				Type:    TypeString,
				Literal: []rune("「233」 ‘456’ 《〈who〉》『『is』』"),
				Info:    '“',
			},
			lineInfo: "",
		},
		{
			name:        "multiple-line string",
			input:       "『233\n    456\r\n7  』",
			expectError: false,
			token: Token{
				Type:    TypeString,
				Literal: []rune("233\n    456\r\n7  "),
				Info:    '『',
			},
			lineInfo: "Unknown<0>[0,3] Unknown<0>[5,11]",
		},
	}

	assertNextToken(cases, t)
}

func TestNextToken_VarQuoteONLY(t *testing.T) {
	cases := []nextTokenCase{
		{
			name:        "normal variable quote",
			input:       "·正常之变量·",
			expectError: false,
			token: Token{
				Type:    TypeVarQuote,
				Literal: []rune("正常之变量"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal variable quote (with spaces)",
			input:       "· 正常 之 变量  ·",
			expectError: false,
			token: Token{
				Type:    TypeVarQuote,
				Literal: []rune("正常之变量"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal variable quote (with slashs)",
			input:       "· 知/其/不- 可/而*为+ _abcd_之1235 AJ·",
			expectError: false,
			token: Token{
				Type:    TypeVarQuote,
				Literal: []rune("知/其/不-可/而*为+_abcd_之1235AJ"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal variable quote - english variable",
			input:       "·_korea_char102·",
			expectError: false,
			token: Token{
				Type:    TypeVarQuote,
				Literal: []rune("_korea_char102"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "invalid quote - number at first",
			input:       "·123ABC·",
			expectError: true,
			token:       Token{},
			lineInfo:    "",
		},
		{
			name:        "invalid quote - invalid punctuation",
			input:       "·正（大）光明·",
			expectError: true,
			token:       Token{},
			lineInfo:    "",
		},
		{
			name:        "invalid quote - char buffer overflow",
			input:       "·这是一个很长变量这是一个很长变量这是一个很长变量这是一个很长变量这是一个很长变量·",
			expectError: true,
			token:       Token{},
			lineInfo:    "",
		},
		{
			name:        "invalid quote - CR, LFs are not allowed inside quotes",
			input:       "·变量\r\n又是变量名·",
			expectError: true,
			token:       Token{},
			lineInfo:    "",
		},
	}
	assertNextToken(cases, t)
}

func TestNextToken_NumberONLY(t *testing.T) {
	// NOTE 1:
	// nums such as 2..3 will be regarded as `2.`(2.0) and `.3`(0.3) combination
	cases := []nextTokenCase{
		{
			name:        "normal number (all digits)",
			input:       "123456七",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("123456"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (start to end)",
			input:       "123456",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("123456"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (with dot and minus)",
			input:       "-123.456km",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("-123.456"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (with plus at beginning)",
			input:       "+00000.456km",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("+00000.456"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (with plus)",
			input:       "+000003 Rs",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("+000003"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (with E+)",
			input:       "+000003E+05 Rs",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("+000003E+05"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (with e-)",
			input:       "+000003e-25 Rs",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("+000003e-25"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (decimal with e+)",
			input:       "-003.0452e+25 Rs",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("-003.0452e+25"),
				Info:    nil,
			},
			lineInfo: "",
		},
		{
			name:        "normal number (ignore underscore)",
			input:       "-12_300_500_800_900 RSU",
			expectError: false,
			token: Token{
				Type:    TypeNumber,
				Literal: []rune("-12300500800900"),
				Info:    nil,
			},
			lineInfo: "",
		},
		// test fail cases
		{
			name:        "operater only",
			input:       "---",
			expectError: true,
			token:       Token{},
			lineInfo:    "",
		},
		{
			name:        "operater only #2",
			input:       "-++",
			expectError: true,
			token:       Token{},
			lineInfo:    "",
		},
		{
			name:        "E first",
			input:       "-E+3",
			expectError: true,
			token:       Token{},
			lineInfo:    "",
		},
	}

	assertNextToken(cases, t)
}

func TestNextToken_MarkerONLY(t *testing.T) {
	// 01. generate TRUE cases
	var markerMap = map[string]TokenType{
		"，":  TypeCommaSep,
		"：":  TypeFuncCall,
		"；":  TypeStmtSep,
		"？":  TypeFuncDeclare,
		"&":  TypeObjRef,
		"！":  TypeMustT,
		"@":  TypeAnnoT,
		"#":  TypeMapHash,
		"……": TypeMoreParam,
		"【":  TypeArrayQuoteL,
		"】":  TypeArrayQuoteR,
		"（":  TypeStmtQuoteL,
		"）":  TypeStmtQuoteR,
		"==": TypeMapData,
		"⟺":  TypeMapData,
	}

	var cases = make([]nextTokenCase, 0)
	for k, v := range markerMap {
		cases = append(cases, nextTokenCase{
			name:        fmt.Sprintf("generate token %s", k),
			input:       fmt.Sprintf("%s EE", k),
			expectError: false,
			token: Token{
				Type:    v,
				Literal: []rune(k),
				Info:    nil,
			},
		})
	}

	assertNextToken(cases, t)
}

func TestNextToken_IdentifierONLY_SUCCESS(t *testing.T) {
	cases := []nextTokenCase{
		{
			name:        "normal identifier",
			input:       "反",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("反"),
			},
		},
		{
			name:        "normal identifier #2",
			input:       "正定县",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县"),
			},
		},
		{
			name:        "normal identifier #3 with spaces",
			input:       "正  定  县",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县"),
			},
		},
		{
			name:        "normal identifier with number followed",
			input:       "正定县2345",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县2345"),
			},
		},
		{
			name:        "normal identifier with + - * /",
			input:       "正定/+_县/2345",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定/+_县/2345"),
			},
		},
		{
			name:        "normal identifier (quote as terminator)",
			input:       "正定县「」",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县"),
			},
		},
		{
			name:        "normal identifier (var quote as terminator)",
			input:       "正定县·如果·",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县"),
			},
		},
		{
			name:        "normal identifier (marker) as terminator)",
			input:       "正定县（河北）",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县"),
			},
		},
		{
			name:        "normal identifier (keyword as terminator)",
			input:       "正定县作为大县",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县"),
			},
		},
		{
			name:        "normal identifier (following keyword lead but not keyword formed)",
			input:       "正定县如大县",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("正定县如大县"),
			},
		},
		{
			name:        "normal identifier (like keyword)",
			input:       "如不果返回",
			expectError: false,
			token: Token{
				Type:    TypeIdentifier,
				Literal: []rune("如不果"),
			},
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
				} else if !reflect.DeepEqual(*tk, tt.token) {
					t.Errorf("NextToken() return Token failed! expect: %v, got: %v", tt.token, *tk)
				}
			} else {
				if err == nil {
					t.Errorf("NextToken() failed! expected error, but got no error")
				}
			}
		})
	}
}
