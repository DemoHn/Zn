package syntax

import "github.com/DemoHn/Zn/lex"
import zerr "github.com/DemoHn/Zn/pkg/error"

// Lexer is used to tokenizing the code and yield tokens, after lexing process, Parser will analysis
// tokens and transform to AST.
type Lexer struct {
	// Source - code source
	Source []rune
	// IndentType - intent type TAB or SPACE
	IndentType uint8
	// Lines - store lines
	Lines []LineInfo
	// current cursor
	cursor int
	// current line
	currentLine int
}

// LineInfo - stores the (absolute) start & end index of this line
//// this should be added to the scanner after all parsing is done
type LineInfo struct {
	// the indent number (at the beginning) of this line.
	// all lines should have indents to differentiate scopes.
	Indents int
	// startIdx - start index of lineBuffer
	StartIdx int
}

// define constants
const (
	IndentUnknown uint8 = 0
	IndentTab     uint8 = 9
	IndentSpace   uint8 = 32

	// define scanStates
	// example:
	//
	// [ INDENTS ] [ TEXT TEXT TEXT ] [ CR LF ]
	// ^         ^                  ^
	// 0         1                  2
	//
	// 0: ScanInit
	// 1: ScanIndent
	// 2: ScanEnd
	ScanInit   uint8 = 0
	ScanIndent uint8 = 1
	ScanEnd    uint8 = 2

	RuneEOF rune = 0
	RuneSP  rune = 0x0020 // <SP>
	RuneTAB rune = 0x0009 // <TAB>
	RuneCR  rune = 0x000D // \r
	RuneLF  rune = 0x000A // \n
)

func NewLexer(source []rune) *Lexer {
	return &Lexer{
		Source: source,
		IndentType: IndentUnknown,
		Lines: []LineInfo{},
		cursor: -1,
		currentLine: 0,
	}
}

// Next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) Next() rune {
	l.cursor++

	// still no data, return EOF directly
	return l.getChar(l.cursor)
}

// Peek - get the character of the cursor
func (l *Lexer) Peek() rune {
	return l.getChar(l.cursor + 1)
}

// Peek2 - get the next next character without moving the cursor
func (l *Lexer) Peek2() rune {
	return l.getChar(l.cursor + 2)
}

// GetCursor - get current cursor
func (l *Lexer) GetCursor() int {
	return l.cursor
}

// SetCursor - set cursor
func (l *Lexer) SetCursor(cursor int) {
	l.cursor = cursor
}

// getChar - get value from lineBuffer
func (l *Lexer) getChar(idx int) rune {
	if idx >= len(l.Source) {
		return RuneEOF
	}
	return l.Source[idx]
}

func (l *Lexer) getIndentNum(count int, t lex.IndentType) (int, error) {
	switch t {
	case IndentUnknown:
		if count > 0 && l.IndentType != t {
			return 0, zerr.InvalidIndentType(l.IndentType, t)
		}
	case IndentSpace, IndentTab:
		// init ls.IndentType
		if l.IndentType == IndentUnknown {
			l.IndentType = t
		}
		// when t = space, the character count must be 4 * N
		if t == IndentSpace && count%4 != 0 {
			return 0, zerr.InvalidIndentSpaceCount(count)
		}
		// when t does not match indentType, throw error
		if l.IndentType != t {
			return 0, zerr.InvalidIndentType(l.IndentType, t)
		}
	}
	// when indentType = TAB, count = indents
	// otherwise, count = indents * 4
	indentNum := count
	if l.IndentType == IndentSpace {
		indentNum = count / 4
	}

	return indentNum, nil
}

// ParseLine - when cursor parsed to the end line (ch = CR, LF or EOF)
// then record line info, parse CRLFs and indents, move cursor to the start of
// next line
// NOTE: ch would be CR or LF only.
func (l *Lexer) ParseLine(ch rune) error {
	p := l.Peek()
	switch p {
	case RuneCR, RuneLF:
		// read line break chars
		if (ch == RuneCR && p == RuneLF) || (ch == RuneLF && p == RuneCR) {
			l.Next()
			p = l.Peek()
		}
		// or indent chars
	}

	if p == RuneEOF {
		return nil
	}

	// append next line info
	l.Lines = append(l.Lines, LineInfo{
		Indents:  0,
		StartIdx: l.GetCursor() + 1,
	})
	l.currentLine += 1

	// parse next line's indents
	count := 0
	if containsRune(p, []rune{RuneSP, RuneTAB}) {
		for l.Peek() == p {
			count += 1
			l.Next()
		}
	}

	indentType := IndentUnknown
	switch p {
	case RuneTAB:
		indentType = IndentTab
	case RuneSP:
		indentType = IndentSpace
	}

	// get indent
	indentNum, err := l.getIndentNum(count, indentType)
	if err != nil {
		return err
	}

	// set indent num of current line
	lastIdx := l.currentLine - 1
	l.Lines[lastIdx].Indents = indentNum
	return nil
}