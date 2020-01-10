package lex

import (
	"fmt"
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

// StringifyLines - stringify current parsed lines into readable string info
// format::=
//   {lineInfo1} {lineInfo2} {lineInfo3} ...
//
// lineInfo ::=
//   Space<2>[23,45] or
//   Tab<4>[0,1] or
//   Empty<0>
func StringifyLines(ls *LineStack) string {
	ss := []string{}
	var indentChar string
	// get indent type
	switch ls.IndentType {
	case IdetUnknown:
		indentChar = "Unknown"
	case IdetSpace:
		indentChar = "Space"
	case IdetTab:
		indentChar = "Tab"
	}

	for _, line := range ls.lines {
		if line.Indents == 0 && len(line.Source) == 0 {
			ss = append(ss, "Empty<0>")
		} else {
			ss = append(ss, fmt.Sprintf(
				"%s<%d>",
				indentChar, line.Indents,
			))
		}

	}
	return strings.Join(ss, " ")
}

func stringifyToken(tk *Token) string {
	return fmt.Sprintf("$%d[%s]", tk.Type, string(tk.Literal))
}
