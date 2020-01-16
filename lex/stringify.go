package lex

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// StringifyToken is the process that transforms an abstract token type to a readable string sequence.
// usually used for debugging.
func StringifyToken(tk *Token) string {
	return stringifyToken(tk)
}

// StringifyAllTokens - similar with StringifyToken, but this time it's for all tokens
func StringifyAllTokens(tks []*Token) string {
	var tokenStrs = []string{}
	for _, tk := range tks {
		tokenStrs = append(tokenStrs, stringifyToken(tk))
	}
	return strings.Join(tokenStrs, " ")
}

// ParseTokenStr - from token str to tokens
func ParseTokenStr(str string) []Token {
	tks := make([]Token, 0)
	r := regexp.MustCompile(`\$(\d+)\[(.+?)\]`)
	matches := r.FindAllStringSubmatch(str, -1)

	for _, match := range matches {
		n, _ := strconv.Atoi(match[1])
		tks = append(tks, Token{
			Type:    TokenType(n),
			Literal: []rune(match[2]),
		})
	}

	return tks
}

// StringifyLines - stringify current parsed lines into readable string info
// format::=
//   {lineInfo1} {lineInfo2} {lineInfo3} ...
//
// lineInfo ::=
//   SP<2>[text1] or
//   T<4>[text2] or
//   E<0>
func StringifyLines(ls *LineStack) string {
	ss := []string{}
	var indentChar string
	// get indent type
	switch ls.IndentType {
	case IdetUnknown:
		indentChar = "U"
	case IdetSpace:
		indentChar = "SP"
	case IdetTab:
		indentChar = "T"
	}

	for _, line := range ls.lines {
		if len(line.Source) == 0 {
			ss = append(ss, fmt.Sprintf("E<%d>", line.Indents))
		} else {
			ss = append(ss,
				fmt.Sprintf("%s<%d>[%s]", indentChar, line.Indents, string(line.Source)))
		}
	}
	return strings.Join(ss, " ")
}

func stringifyToken(tk *Token) string {
	return fmt.Sprintf("$%d[%s]", tk.Type, string(tk.Literal))
}
