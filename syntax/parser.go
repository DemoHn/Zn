package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// Parser - parse all nodes
type Parser struct {
	*lex.Lexer
	tokens   [3]*lex.Token
	lineMask uint16
}

const (
	modeInline uint16 = 0x01
)

// NewParser -
func NewParser(l *lex.Lexer) *Parser {
	return &Parser{
		Lexer: l,
	}
}

// Parse - parse all tokens into an AST (stored as ProgramNode)
func (p *Parser) Parse() (pg *Program, err *error.Error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(*error.Error)
		}
		handleDeferError(p, err)
	}()

	// advance tokens TWICE
	p.next()
	p.next()

	pg = new(Program)

	peekIndent := p.getPeekIndent()
	// parse global block
	pg.Content, err = ParseBlockStmt(p, peekIndent)
	if err != nil {
		return
	}

	// ensure there's no remaining token after parsing global block
	if p.peek().Type != lex.TypeEOF {
		err = error.UnexpectedEOF()
	}

	return
}

func (p *Parser) next() *lex.Token {
	var tk *lex.Token
	var err *error.Error

	tk, err = p.NextToken()
	if err != nil {
		panic(err)
	}

	// after retrieving next token successfully, check if current token has
	// violate lineMasks
	// check the comment of validateLineMask() for details
	if p.tokens[0] != nil && p.tokens[1] != nil {
		if err = p.validateLineMask(p.tokens[0], p.tokens[1]); err != nil {
			panic(err)
		}
	}

	// move advanced token buffer
	p.tokens[0] = p.tokens[1]
	p.tokens[1] = p.tokens[2]
	p.tokens[2] = tk

	return p.tokens[0]
}

func (p *Parser) current() *lex.Token {
	return p.tokens[0]
}

func (p *Parser) peek() *lex.Token {
	return p.tokens[1]
}

func (p *Parser) peek2() *lex.Token {
	return p.tokens[2]
}

func (p *Parser) setLineMask(loc uint16) {
	p.lineMask = p.lineMask | loc
}

func (p *Parser) unsetLineMask(loc uint16) {
	p.lineMask = p.lineMask & (^loc)
}

func (p *Parser) validateLineMask(lastToken *lex.Token, newToken *lex.Token) *error.Error {

	var line1 = lastToken.Range.EndLine
	var line2 = newToken.Range.StartLine
	// if new token has entered a new line
	if line2 > line1 {
		// for modeInline, all tokens should have no explicit newline (CRLF)
		if (p.lineMask & modeInline) > 0 {
			return error.IncompleteStmt()
		}
	}
	return nil
}

// consume one token with denoted validTypes
// if not, return syntaxError
func (p *Parser) consume(validTypes ...lex.TokenType) {
	tkType := p.peek().Type
	for _, item := range validTypes {
		if item == tkType {
			p.next()
			return
		}
	}
	err := error.InvalidSyntax()
	panic(err)
}

// trying to consume one token. if the token is valid in the given range of tokenTypes,
// will return its tokenType; if not, then nothing will happen.
//
// returns (matched, tokenType)
func (p *Parser) tryConsume(validTypes []lex.TokenType) (bool, *lex.Token) {
	tk := p.peek()
	for _, vt := range validTypes {
		if vt == tk.Type {
			p.next()
			return true, tk
		}
	}

	return false, nil
}

// expectBlockIndent - detect if the Indent(peek) == Indent(current) + 1
// returns (validBlockIndent, newIndent)
func (p *Parser) expectBlockIndent() (bool, int) {
	var peekLine = p.peek().Range.StartLine
	var currLine = p.current().Range.StartLine

	var peekIndent = p.GetLineIndent(peekLine)
	var currIndent = p.GetLineIndent(currLine)

	if peekIndent == currIndent+1 {
		return true, peekIndent
	}
	return false, 0
}

// getPeekIndent -
func (p *Parser) getPeekIndent() int {
	var peekLine = p.peek().Range.StartLine

	return p.GetLineIndent(peekLine)
}

// similar to lexer's version, but with given line & col
func moveAndSetCursor(p *Parser, line int, col int, err *error.Error) {
	buf := p.GetLineBuffer()
	cursor := error.Cursor{
		File:    p.Lexer.InputStream.Scope,
		ColNum:  col,
		LineNum: line,
		Text:    string(buf),
	}

	defer func() {
		// recover but not handle it
		recover()
		err.SetCursor(cursor)
	}()

	endCursor := p.SlideToLineEnd()
	cursor.Text = string(buf[:endCursor+1])
	err.SetCursor(cursor)
}

func handleDeferError(p *Parser, err *error.Error) {
	var tk *lex.Token

	if err != nil && err.GetErrorClass() == error.SyntaxErrorClass {
		if cursorType, ok := err.GetInfo()["cursor"]; ok {
			if cursorType == "peek" {
				tk = p.peek()
			} else if cursorType == "current" {
				tk = p.current()
			}
			if tk != nil {
				line := tk.Range.StartLine
				col := tk.Range.StartCol
				moveAndSetCursor(p, line, col, err)
			}
		}
	}
}
