package zh

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

// ParserZH - parse all nodes
type ParserZH struct {
	*syntax.Lexer
	// current token
	TokenP1 *syntax.Token
	// which line is startIdx located (P1)
	StartLineIdxP1 int
	// which line is endIdx located
	EndLineIdxP1 int
	// peek token
	TokenP2 *syntax.Token
	// which line is startIdx located (P2)
	StartLineIdxP2 int
	// which line is endIdx located
	EndLineIdxP2 int
	// statement completion flag - if set true, next statement MUST
	// 1) start from another line OR
	// 2) seperate former statement with '；'
	stmtCompleteFlag bool
}

// NewParserZH -
func NewParserZH() *ParserZH {
	return &ParserZH{
		stmtCompleteFlag: false,
	}
}

// Parse - parse all tokens into an AST (stored as ProgramNode)
func (p *ParserZH) ParseAST(l *syntax.Lexer) (pg *syntax.Program, err error) {
	// set lexer
	p.Lexer = l
	// advance tokens ONCE
	p.next()

	peekIndent := p.getPeekIndent()
	// parse global block
	block := ParseBlockStmt(p, peekIndent)
	pg = &syntax.Program{
		Lexer:   l,
		Content: block,
	}

	// ensure there's no remaining token after parsing global block
	if p.peek().Type != TypeEOF {
		err = p.getInvalidSyntaxCurr()
	}
	return
}

func (p *ParserZH) next() *syntax.Token {
	var tk syntax.Token // default tk.Type = 0 (TypeEOF)
	var err error
	// init next tk first
	tk, err = NextToken(p.Lexer)
	if err != nil {
		panic(err)
	}

	// skip comment token until meet non-comment one
	for tk.Type == TypeComment {
		tk, err = NextToken(p.Lexer)
		if err != nil {
			panic(err)
		}
	}

	// move advanced token buffer
	p.TokenP1 = p.TokenP2
	p.TokenP2 = &tk

	// get peek token's startLine and endLine
	p.StartLineIdxP1 = p.StartLineIdxP2
	p.EndLineIdxP1 = p.EndLineIdxP2
	p.StartLineIdxP2 = p.FindLineIdx(tk.StartIdx, p.StartLineIdxP2)
	p.EndLineIdxP2 = p.FindLineIdx(tk.EndIdx, p.EndLineIdxP2)

	if p.meetStmtLineBreak() {
		p.setStmtCompleteFlag()
	}
	return p.TokenP1
}

func (p *ParserZH) current() *syntax.Token {
	return p.TokenP1
}

func (p *ParserZH) peek() *syntax.Token {
	return p.TokenP2
}

// meetStmtLineBreak - if there's a statement line-break at the end of token.
//
// StmtLineBreak is not an existing token, i.e. it wouldn't insert into token stream.
// It's a virtual mark that separates parsing statements, thus you can image it like a
// "virtually" inserted semicolon, like the following code:
//
// 如果此IDE之 ';' <-- here is the statement line-break
//     名为「VSCODE」
//
// There's stmt line-break at the end of first line, thus the process of parsing IF-statement
// should terminate due to lacking matched tokens, similar to the behaviour that an semicolon is inserted.
//
// Theoretically, any type of token that located at the end of line
//
//   i.e.  this token is the last one of current line,
//   or    $token.Range.EndLine < ($token+1).Range.StartLine,
//
// should meet StmtLineBreak, which means meetStmtLineBreak() = true. Still, there are some
// exceptions listed below:
//
//   1.    the current token type is one of the following 6 punctuations:  ， 、  {  【  ：  ？
//   or
//   2.    the next token type if one of the following 3 marks:  】  }  EOF
//
// Example 1#, 2#, 3# illustrates those exceptions that even if there are two or more lines, it's still
// *ONE* valid and complete statement:
//
// Example 1#
//
// （显示并连缀：
//     「结果为」、
//      人口-中位数）
//
// Example 2#
//
// 令天干地支表为【
//     「子」 == 「甲」，
//     「丑」 == 「乙」，
//     「寅」 == 「丙」，
//     「卯」 == 「丁」
// 】
//
// Example 3#
//
// `时`等于12 且{
//     `分`等于0 或 `分`等于30
// }等于真
func (p *ParserZH) meetStmtLineBreak() bool {
	current := p.current()
	peek := p.peek()

	exceptCurrentTokenTypes := []uint8{
		TypeCommaSep,
		TypePauseCommaSep,
		TypeStmtQuoteL,
		TypeArrayQuoteL,
		TypeFuncCall,
		TypeFuncDeclare,
	}

	exceptFollowingTokenTypes := []uint8{
		TypeArrayQuoteR,
		TypeStmtQuoteR,
	}

	if peek == nil || current == nil {
		return false
	}

	// while parsing last token of the file, there must be a line break
	// to terminate the statement -- no reason that the parsing progress
	// be continued.
	if current.Type == TypeEOF || peek.Type == TypeEOF {
		return true
	}

	// current token is at line end
	if p.StartLineIdxP2 > p.EndLineIdxP1 {
		// exception rule 1
		for _, currTk := range exceptCurrentTokenTypes {
			if currTk == current.Type {
				return false
			}
		}
		// exception rule 2
		for _, followingTk := range exceptFollowingTokenTypes {
			if followingTk == peek.Type {
				return false
			}
		}
		return true
	}
	return false
}

// meetStmtBreak - similar to `meetStmtLineBreak`
func (p *ParserZH) meetStmtBreak() bool {
	peek := p.peek()
	if peek.Type == TypeStmtSep || peek.Type == TypeEOF {
		return true
	}
	return false
}

func (p *ParserZH) unsetStmtCompleteFlag() {
	p.stmtCompleteFlag = false
}

func (p *ParserZH) setStmtCompleteFlag() {
	p.stmtCompleteFlag = true
}

// trying to consume one token. if the token is valid in the given range of tokenTypes,
// will return its tokenType; if not, then nothing will happen.
//
// returns (matched, tokenType)
func (p *ParserZH) tryConsume(validTypes ...uint8) (bool, *syntax.Token) {
	tk := p.peek()
	// if next token is comma, then ignore comma (only once!) and
	// read next token
	if tk.Type == TypeCommaSep {
		p.next()
		tk = p.peek()
	}
	if p.stmtCompleteFlag {
		return false, nil
	}

	for _, vt := range validTypes {
		if vt == tk.Type {
			p.next()
			return true, tk
		}
	}

	return false, nil
}

// consume one token with denoted validTypes
// if not, return syntaxError
func (p *ParserZH) consume(validTypes ...uint8) {
	match, _ := p.tryConsume(validTypes...)
	if !match {
		panic(p.getInvalidSyntaxPeek())
	}
}

// expectBlockIndent - detect if the Indent(peek) == Indent(current) + 1
// returns (validBlockIndent, newIndent)
func (p *ParserZH) expectBlockIndent() (bool, int) {
	var peekLine = p.StartLineIdxP2
	var currLine = p.StartLineIdxP1

	var peekIndent = p.GetLineInfo(peekLine).Indents
	var currIndent = p.GetLineInfo(currLine).Indents

	if peekIndent == currIndent+1 {
		return true, peekIndent
	}
	return false, 0
}

// getPeekIndent -
func (p *ParserZH) getPeekIndent() int {
	var peekLine = p.StartLineIdxP2

	lineInfo := p.GetLineInfo(peekLine)
	if lineInfo == nil {
		return 0
	}
	return lineInfo.Indents
}

// equals to s.SetCurrentLine(<line of tk>)
func (p *ParserZH) setStmtCurrentLine(s syntax.Statement, tk *syntax.Token) {
	if tk != nil {
		idx := p.FindLineIdx(tk.StartIdx, 0)
		s.SetCurrentLine(idx)
	}
}

// wrap 0x2250 InvalidSyntaxCurr - with current token's startIdx
func (p *ParserZH) getInvalidSyntaxCurr() error {
	startIdx := p.TokenP1.StartIdx
	return zerr.InvalidSyntax(startIdx)
}

func (p *ParserZH) getInvalidSyntaxPeek() error {
	startIdx := p.TokenP1.StartIdx
	if p.TokenP2 != nil {
		startIdx = p.TokenP2.StartIdx
	}

	return zerr.InvalidSyntax(startIdx)
}

func (p *ParserZH) getUnexpectedIndentPeek() error {
	startIdx := p.TokenP1.StartIdx
	if p.TokenP2 != nil {
		startIdx = p.TokenP2.StartIdx
	}

	return zerr.UnexpectedIndent(startIdx)
}

func (p *ParserZH) getExprMustTypeIDPeek() error {
	startIdx := p.TokenP1.StartIdx
	if p.TokenP2 != nil {
		startIdx = p.TokenP2.StartIdx
	}

	return zerr.ExprMustTypeID(startIdx)
}
