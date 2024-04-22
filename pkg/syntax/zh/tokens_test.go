package zh

import (
	"testing"

	"github.com/DemoHn/Zn/pkg/syntax"
)

type nextTokenCase struct {
	name        string
	input       string
	expectError bool
	// [(type, startIdx, endIdx), (type, startIdx, endIdx), ...]
	tokens [][]int
}

func TestNextToken_NumberONLY(t *testing.T) {
	// NOTE 1:
	// nums such as 2..3 will be regarded as `2.`(2.0) and `.3`(0.3) combination
	cases := []nextTokenCase{
		{
			name:        "normal number (all digits)",
			input:       "123456七",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 7},
			},
		},
		{
			name:        "normal number (start to end)",
			input:       "1234567",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 7},
			},
		},
		{
			name:        "normal number (with dot and minus)",
			input:       "/* comment */ -123.456km",
			expectError: false,
			tokens: [][]int{
				{int(TypeComment), 0, 13},
				{int(TypeIdentifier), 14, 24},
			},
		},
		{
			name:        "normal number (with plus at beginning)",
			input:       "+00000.456km",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 12},
			},
		},
		{
			name:        "normal number (with plus)",
			input:       "+000003 Rs",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 7},
			},
		},
		{
			name:        "normal number (with E+)",
			input:       "+000003E+05 Rs",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 11},
			},
		},
		{
			name:        "normal number (with e-)",
			input:       "+000003e-25 Rs",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 11},
			},
		},
		{
			name:        "normal number (decimal with e+)",
			input:       "-003.0452e+25 Rs",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 13},
			},
		},
		{
			name:        "arithmetic expression",
			input:       "25 / +3",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 2},
				{int(TypeDivision), 3, 4},
				{int(TypeIdentifier), 5, 7},
			},
		},
		{
			name:        "*10^ as E",
			input:       "23.5*10^8",
			expectError: false,
			tokens: [][]int{
				{int(TypeIdentifier), 0, 9},
			},
		},
		/**
		// test fail cases
		{
			name:        "operater only",
			input:       "---",
			expectError: true,
			errCursor:   1,
		},
		{
			name:        "operater only #2",
			input:       "-++",
			expectError: true,
			errCursor:   1,
		},
		{
			name:        "E first",
			input:       "-E+3",
			expectError: true,
			errCursor:   1,
		},
		{
			name:        "E without following PM mark",
			input:       "2395.234E34",
			expectError: true,
			errCursor:   9,
		},
		{
			name:        "number with other weird char",
			input:       "23.r",
			expectError: true,
			errCursor:   3,
		},
		{
			name:        "numbers *9^",
			input:       "1111*9^23",
			expectError: true,
			errCursor:   5,
		},
		{
			name:        "incomplete *10^",
			input:       "1234*10^",
			expectError: true,
			errCursor:   8,
		},
		*/
	}

	assertParseTokens(cases, t)
}

func assertParseTokens(cases []nextTokenCase, t *testing.T) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tks, _, err := parseTokens([]rune(tt.input))

			// validate error
			if tt.expectError == false {
				if len(tks) == 0 {
					t.Errorf("expect token, but no token")
				}
				if err != nil {
					t.Errorf("NextToken() failed! expected no error, but got error")
					t.Error(err)
				} else {
					for i, tk := range tt.tokens {
						tm := tks[i]
						// tk[0] <--> type
						if tk[0] != int(tm.Type) {
							t.Errorf("idx[%d] token: type not match, expect:%d, got:%d", i, tk[0], int(tm.Type))
						}
						if tk[1] != tm.StartIdx {
							t.Errorf("idx[%d] token: startIdx not match, expect:%d, got:%d", i, tk[1], tm.StartIdx)
						}
						if tk[2] != tm.EndIdx {
							t.Errorf("idx[%d] token: endIdx not match, expect:%d, got:%d", i, tk[2], tm.EndIdx)
						}
					}
				}
			} else {
				if err == nil {
					t.Errorf("NextToken() failed! expected error, but got no error")
				}
			}
		})
	}
}

func parseTokens(source []rune) ([]syntax.Token, []syntax.LineInfo, error) {
	l := syntax.NewLexer(source)
	var tks []syntax.Token
	for {
		tk, err := NextToken(l)
		if err != nil {
			return tks, l.Lines, err
		}

		if tk.Type == TypeEOF {
			break
		}
		tks = append(tks, tk)
	}

	return tks, l.Lines, nil
}
