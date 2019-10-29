package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex/tokens"
)

// Tokenize is the main function of lex who transforms code to a sequence
// of pre-defined tokens.
func Tokenize(code []rune) ([]Token, *error.Error) {
	lexer := NewLexer(code)

	for !lexer.End() {

	}
	lexer.AppendToken(tokens.EOFToken{})
	return lexer.Tokens, nil
}

//// private parsers

//// TODO: consider different levels of quoting
func parseOperators(lexer *Lexer) {
	var tk Token
	var idx = lexer.GetIndex()
	switch lexer.Next() {
	// map Quote I - V
	case tokens.LeftQuoteI:
		tk = tokens.QuoterToken{1, true, []rune{tokens.LeftQuoteI}, idx, idx}
	case tokens.RightQuoteI:
		tk = tokens.QuoterToken{1, false, []rune{tokens.RightQuoteI}, idx, idx}
	case tokens.LeftQuoteII:
		tk = tokens.QuoterToken{2, true, []rune{tokens.LeftQuoteII}, idx, idx}
	case tokens.RightQuoteII:
		tk = tokens.QuoterToken{2, false, []rune{tokens.RightQuoteII}, idx, idx}
	case tokens.LeftQuoteIII:
		tk = tokens.QuoterToken{3, true, []rune{tokens.LeftQuoteIII}, idx, idx}
	case tokens.RightQuoteIII:
		tk = tokens.QuoterToken{3, false, []rune{tokens.RightQuoteIII}, idx, idx}
	case tokens.LeftQuoteIV:
		tk = tokens.QuoterToken{4, true, []rune{tokens.LeftQuoteIV}, idx, idx}
	case tokens.RightQuoteIV:
		tk = tokens.QuoterToken{4, false, []rune{tokens.RightQuoteIV}, idx, idx}
	case tokens.LeftQuoteV:
		tk = tokens.QuoterToken{5, true, []rune{tokens.LeftQuoteV}, idx, idx}
	case tokens.RightQuoteV:
		tk = tokens.QuoterToken{5, false, []rune{tokens.RightQuoteV}, idx, idx}

	}
}
