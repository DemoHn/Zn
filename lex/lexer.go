package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/util"
)

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	*LineStack
	quoteStack   *util.RuneStack
	*InputStream        // input stream
	chBuffer     []rune // the buffer for parsing & generating tokens
	lineBuffer   []rune
	cursor       int
	blockSize    int
}

// NewLexer - new lexer
func NewLexer(in *InputStream) *Lexer {
	return &Lexer{
		LineStack:   NewLineStack(),
		quoteStack:  util.NewRuneStack(32),
		InputStream: in,
		chBuffer:    []rune{},
		lineBuffer:  []rune{},
		cursor:      0,
		blockSize:   256,
	}
}

// next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) next() rune {
	if l.cursor+2 >= len(l.lineBuffer) {
		if l.End() {
			return EOF
		}
		if b, err := l.Read(l.blockSize); err == nil {
			l.lineBuffer = append(l.lineBuffer, b...)
		} else {
			// TODO: handle error
		}
	}

	// still no data, return EOF directly
	if l.cursor >= len(l.lineBuffer) {
		return EOF
	}
	data := l.lineBuffer[l.cursor]
	l.cursor = l.cursor + 1
	return data
}

// peek - get the character of the cursor
func (l *Lexer) peek() rune {
	if l.cursor+1 >= len(l.lineBuffer) {
		return EOF
	}
	return l.lineBuffer[l.cursor+1]
}

// peek2 - get the next next character without moving the cursor
func (l *Lexer) peek2() rune {
	if l.cursor+2 >= len(l.lineBuffer) {
		return EOF
	}
	return l.lineBuffer[l.cursor+2]
}

// current - get current char value
func (l *Lexer) current() rune {
	if l.cursor >= len(l.lineBuffer) {
		return EOF
	}
	return l.lineBuffer[l.cursor]
}

// rebase - rebase cursor within the same line
func (l *Lexer) rebase(cursor int) {
	l.cursor = cursor
}

func (l *Lexer) clearBuffer() {
	l.chBuffer = []rune{}
}

func (l *Lexer) pushBuffer(ch ...rune) {
	l.chBuffer = append(l.chBuffer, ch...)
}

// NextToken - parse and generate the next token (including comments)
func (l *Lexer) NextToken() (*Token, *error.Error) {
head:
	var ch = l.next()
	switch ch {
	case EOF:
		l.PushLine(l.lineBuffer, l.cursor)
		return NewTokenEOF(), nil
	case SP, TAB:
		// if indent has been scanned, it should be regarded as whitespaces
		// (it's totally ignored)
		if !l.HasScanIndent() {
			l.parseIndents(ch)
		} else {
			l.consumeWhiteSpace(ch)
		}
		goto head
	case CR, LF:
		l.parseCRLF(ch)
		goto head
	// meet with 注, it may be possibly a lead character of a comment block
	// notice: it would also be a normal identifer (if 注[number]：) does not satisfy.
	case GlyphZHU:
		cursor := l.cursor
		isComment, isMultiLine := l.validateComment(ch)
		if isComment {
			return l.parseComment(l.current(), isMultiLine)
		}

		l.rebase(cursor)
		// goto normal identifier
	// left quotes
	case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
		return l.parseString(ch)
	case MiddleDot:
		return l.parseVarQuote(ch)
	default:
		// skip whitespaces
		if isWhiteSpace(ch) {
			l.consumeWhiteSpace(ch)
			goto head
		}
		// parse number
		if isNumber(ch) || ch == '+' || ch == '-' {
			return l.parseNumber(ch)
		}
		if util.Contains(ch, MarkLeads) {
			return l.parseMarkers(ch)
		}
		// suppose it's a keyword
		if isKeyword, tk := l.parseKeyword(ch, true); isKeyword {
			return tk, nil
		}
	}
	return l.parseIdentifier(ch)
}

//// parsing logics
func (l *Lexer) parseIndents(ch rune) *error.Error {
	count := 0
	for l.peek() == ch {
		count++
		l.next()
	}
	// determine indentType
	indentType := IdetUnknown
	switch ch {
	case TAB:
		indentType = IdetTab
	case SP:
		indentType = IdetSpace
	}
	return l.SetIndent(count, indentType)
}

// parseCRLF and return the newline chars by the way
func (l *Lexer) parseCRLF(ch rune) []rune {
	p := l.peek()
	// for CRLF <windows type> or LFCR
	if (ch == CR && p == LF) ||
		(ch == LF && p == CR) {

		// skip one char since we have judge two chars
		l.next()
		l.PushLine(l.lineBuffer, l.cursor)
		return []rune{ch, p}
	}
	// for LF or CR only
	// LF: <linux>, CR:<old mac>
	l.PushLine(l.lineBuffer, l.cursor)
	return []rune{ch}
}

// validate if the coming block is a comment block
// valid comment block are listed below:
// (single-line)
// 1. 注：
// 2. 注123456：
//
// (multi-line)
//
// 3. 注：“
// 4. 注123456：“
//
// @returns (isValid, isMultiLine)
func (l *Lexer) validateComment(ch rune) (bool, bool) {
	// “ or 「
	lquotes := []rune{LeftQuoteIV, LeftQuoteII}
	for {
		ch = l.next()
		// match pattern 1, 2
		if ch == Colon {
			// match pattern 3, 4
			if util.Contains(l.peek(), lquotes) {
				l.next()
				return true, true
			}
			return true, false
		}
		if isNumber(ch) || isWhiteSpace(ch) {
			continue
		}
		return false, false
	}
}

// parseComment until its end
func (l *Lexer) parseComment(ch rune, isMultiLine bool) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	if isMultiLine {
		if !l.quoteStack.Push(ch) {
			return nil, error.NewErrorSLOT("push stack is full")
		}
	}
	// iterate
	for {
		ch = l.next()
		switch ch {
		case EOF:
			l.PushLine(l.lineBuffer, l.cursor)
			return NewCommentToken(l.chBuffer, isMultiLine), nil
		case CR, LF:
			// parse CR,LF first
			nl := l.parseCRLF(ch)
			if !isMultiLine {
				return NewCommentToken(l.chBuffer, isMultiLine), nil
			}
			// for multi-line comment blocks, CRLF is also included
			l.pushBuffer(nl...)

			// manually set no indents
			l.SetIndent(0, IdetUnknown)
		default:
			// for mutli-line comment, calculate quotes is necessary.
			if isMultiLine {
				// push left quotes
				if util.Contains(ch, LeftQuotes) {
					if !l.quoteStack.Push(ch) {
						return nil, error.NewErrorSLOT("quote stack if full")
					}
				}
				// pop right quotes if possible
				if util.Contains(ch, RightQuotes) {
					currentL, _ := l.quoteStack.Current()
					if QuoteMatchMap[currentL] == ch {
						l.quoteStack.Pop()
					}
					// stop quoting
					if l.quoteStack.IsEmpty() {
						return NewCommentToken(l.chBuffer, isMultiLine), nil
					}
				}
			}
			l.pushBuffer(ch)
		}
	}
}

// parseString -
func (l *Lexer) parseString(ch rune) (*Token, *error.Error) {
	// start up
	l.clearBuffer()
	l.quoteStack.Push(ch)
	firstChar := ch
	// iterate
	for {
		ch := l.next()
		switch ch {
		case EOF:
			// after meeting with EOF
			l.PushLine(l.lineBuffer, l.cursor)
			return NewStringToken(l.chBuffer, firstChar), nil
		// push quotes
		case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
			l.pushBuffer(ch)
			if !l.quoteStack.Push(ch) {
				return nil, error.NewErrorSLOT("quote stack is full")
			}
		// pop quotes if match
		case RightQuoteI, RightQuoteII, RightQuoteIII, RightQuoteIV, RightQuoteV:
			currentL, _ := l.quoteStack.Current()
			if QuoteMatchMap[currentL] == ch {
				l.quoteStack.Pop()
			}
			// stop quoting
			if l.quoteStack.IsEmpty() {
				return NewStringToken(l.chBuffer, firstChar), nil
			}
			l.pushBuffer(ch)
		case CR, LF:

			nl := l.parseCRLF(ch)
			// push buffer & mark new line
			l.pushBuffer(nl...)
			l.SetIndent(0, IdetUnknown)
		default:
			l.pushBuffer(ch)
		}
	}
}

func (l *Lexer) parseVarQuote(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	// iterate
	count := 0
	for {
		ch = l.next()
		// we should ensure the following chars to satisfy the condition
		// of an identifier
		switch ch {
		case EOF:
			l.PushLine(l.lineBuffer, l.cursor)
			return NewVarQuoteToken(l.chBuffer), nil
		case MiddleDot:
			return NewVarQuoteToken(l.chBuffer), nil
		default:
			// ignore white-spaces!
			if isWhiteSpace(ch) {
				continue
			}
			if isIdentifierChar(ch, count == 0) {
				l.pushBuffer(ch)
				count++
				if count > maxIdentifierLength {
					return nil, error.NewErrorSLOT("invalid syntax: 变量名太长")
				}
			} else {
				return nil, error.NewErrorSLOT("invalid syntax: invalid var")
			}
		}
	}
}

// regex: ^[-+]?[0-9]*\.?[0-9]+(E[-+]?[0-9]+)?$
func (l *Lexer) parseNumber(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	// hand-written regex parser
	// ref: https://cyberzhg.github.io/toolbox/min_dfa?regex=KG18cCk/TisuP04rKChlfEUpKG18cCk/TispPw==
	// or the documentation has declared that.
	var state = 1
	var endStates = []int{2, 4, 6}

	for {
		switch ch {
		case EOF:
			goto end
		case 'e', 'E':
			switch state {
			case 2, 4:
				state = 5
			default:
				goto end
			}
		case '.':
			switch state {
			case 2:
				state = 4
			default:
				goto end
			}
		case '-', '+':
			switch state {
			case 1:
				state = 3
			case 5:
				state = 7
			default:
				goto end
			}
		case '_':
			ch = l.next()
			continue
		default:
			if isNumber(ch) {
				switch state {
				case 1, 2, 3:
					state = 2
				case 5, 6, 7:
					state = 6
				}
			} else {
				goto end
			}
		}
		l.pushBuffer(ch)
		ch = l.next()
	}

end:
	if util.ContainsInt(state, endStates) {
		// back to last available char
		l.rebase(l.cursor - 1)
		return NewNumberToken(l.chBuffer), nil
	}
	return nil, error.NewErrorSLOT("invalid number")
}

// parseMarkers -
func (l *Lexer) parseMarkers(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	l.pushBuffer(ch)

	// switch
	switch ch {
	case Comma:
		return NewMarkToken(l.chBuffer, TypeCommaSep), nil
	case Colon:
		return NewMarkToken(l.chBuffer, TypeFuncCall), nil
	case Semicolon:
		return NewMarkToken(l.chBuffer, TypeStmtSep), nil
	case QuestionMark:
		return NewMarkToken(l.chBuffer, TypeFuncDeclare), nil
	case RefMark:
		return NewMarkToken(l.chBuffer, TypeObjRef), nil
	case BangMark:
		return NewMarkToken(l.chBuffer, TypeMustT), nil
	case AnnotationMark:
		return NewMarkToken(l.chBuffer, TypeAnnoT), nil
	case HashMark:
		return NewMarkToken(l.chBuffer, TypeMapHash), nil
	case EllipsisMark:
		if l.peek() == EllipsisMark {
			l.pushBuffer(l.next())
			return NewMarkToken(l.chBuffer, TypeMoreParam), nil
		}
		return nil, error.NewErrorSLOT("invalid ellipsis")
	case LeftBracket:
		return NewMarkToken(l.chBuffer, TypeArrayQuoteL), nil
	case RightBracket:
		return NewMarkToken(l.chBuffer, TypeArrayQuoteR), nil
	case LeftParen:
		return NewMarkToken(l.chBuffer, TypeStmtQuoteL), nil
	case RightParen:
		return NewMarkToken(l.chBuffer, TypeStmtQuoteR), nil
	case Equal:
		if l.peek() == Equal {
			l.pushBuffer(l.next())
			return NewMarkToken(l.chBuffer, TypeMapData), nil
		}
		return nil, error.NewErrorSLOT("invalid single euqal")
	case DoubleArrow:
		return NewMarkToken(l.chBuffer, TypeMapData), nil
	}

	return nil, error.NewErrorSLOT("invalid marker")
}

// parseKeyword -
// @return bool matchKeyword
// @return *Token token
//
// when matchKeyword = true, a keyword token will be generated
// matchKeyword = false, regard it as normal identifer
// and return directly.
func (l *Lexer) parseKeyword(ch rune, moveForward bool) (bool, *Token) {
	var tk *Token
	var wordLen = 1

	// manual matching one or consecutive keywords
	switch ch {
	case GlyphLING:
		tk = NewKeywordToken(TypeDeclareW)
	case GlyphWEI:
		tk = NewKeywordToken(TypeLogicYesW)
	case GlyphSHI:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(TypeObjConstructW)
		} else {
			tk = NewKeywordToken(TypeLogicYesIIW)
		}
	case GlyphRU:
		switch l.peek() {
		case GlyphGUO:
			wordLen = 2
			tk = NewKeywordToken(TypeCondW)
		case GlyphHE:
			wordLen = 2
			tk = NewKeywordToken(TypeFuncW)
		default:
			return false, nil
		}
	case GlyphHE:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(TypeStaticFuncW)
		} else {
			return false, nil
		}
	case GlyphYI:
		if l.peek() == GlyphZHIy {
			wordLen = 2
			tk = NewKeywordToken(TypeParamAssignW)
		} else {
			return false, nil
		}
	case GlyphFAN:
		if l.peek() == GlyphHUI {
			wordLen = 2
			tk = NewKeywordToken(TypeReturnW)
		} else {
			return false, nil
		}
	case GlyphBU:
		switch l.peek() {
		case GlyphWEI:
			wordLen = 2
			tk = NewKeywordToken(TypeLogicNotW)
		case GlyphSHI:
			wordLen = 2
			tk = NewKeywordToken(TypeLogicNotIIW)
		case GlyphDENG:
			if l.peek2() == GlyphYU {
				wordLen = 3
				tk = NewKeywordToken(TypeLogicNotEqW)
			} else {
				return false, nil
			}
		case GlyphDA:
			if l.peek2() == GlyphYU {
				wordLen = 3
				tk = NewKeywordToken(TypeLogicLteW)
			} else {
				return false, nil
			}
		case GlyphXIAO:
			if l.peek2() == GlyphYU {
				wordLen = 3
				tk = NewKeywordToken(TypeLogicGteW)
			} else {
				return false, nil
			}
		}
	case GlyphDENG:
		if l.peek() == GlyphYU {
			wordLen = 2
			tk = NewKeywordToken(TypeLogicEqualW)
		} else {
			return false, nil
		}
	case GlyphDA:
		if l.peek() == GlyphYU {
			wordLen = 2
			tk = NewKeywordToken(TypeLogicGtW)
		} else {
			return false, nil
		}
	case GlyphXIAO:
		if l.peek() == GlyphYU {
			wordLen = 2
			tk = NewKeywordToken(TypeLogicLtW)
		} else {
			return false, nil
		}
	case GlyphYIi:
		tk = NewKeywordToken(TypeVarOneW)
	case GlyphDE:
		if l.peek() == GlyphDAO {
			wordLen = 2
			tk = NewKeywordToken(TypeFuncYieldW)
		} else {
			return false, nil
		}
	case GlyphDUI:
		if l.peek() == GlyphYU {
			wordLen = 2
			tk = NewKeywordToken(TypeFuncCallOneW)
		} else {
			return false, nil
		}
	case GlyphFOU:
		if l.peek() == GlyphZE {
			wordLen = 2
			tk = NewKeywordToken(TypeCondElseW)
		} else {
			return false, nil
		}
	case GlyphMEI:
		if l.peek() == GlyphDANG {
			wordLen = 2
			tk = NewKeywordToken(TypeWhileLoopW)
		} else {
			return false, nil
		}
	case GlyphCHENG:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(TypeObjNewW)
		} else {
			return false, nil
		}
	case GlyphZUO:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(TypeVarAliasW)
		} else {
			return false, nil
		}
	case GlyphDING:
		if l.peek() == GlyphYIy {
			wordLen = 2
			tk = NewKeywordToken(TypeObjDefineW)
		} else {
			return false, nil
		}
	case GlyphLEI:
		if l.peek() == GlyphBI {
			wordLen = 2
			tk = NewKeywordToken(TypeObjTraitW)
		} else {
			return false, nil
		}
	case GlyphQI:
		tk = NewKeywordToken(TypeObjThisW)
	case GlyphCI:
		if l.peek() == GlyphZHI {
			wordLen = 2
			tk = NewKeywordToken(TypeStaticSelfW)
		} else {
			tk = NewKeywordToken(TypeObjSelfW)
		}
	case GlyphHUO:
		tk = NewKeywordToken(TypeLogicOrW)
	case GlyphQIE:
		tk = NewKeywordToken(TypeLogicAndW)
	case GlyphZHI:
		tk = NewKeywordToken(TypeObjDotW)
	case GlyphDEo:
		tk = NewKeywordToken(TypeObjDotIIW)
	}

	if tk != nil {
		if moveForward {
			switch wordLen {
			case 1:
				l.pushBuffer(ch)
			case 2:
				l.pushBuffer(ch)
				l.pushBuffer(l.next())
			case 3:
				l.pushBuffer(ch)
				l.pushBuffer(l.next())
				l.pushBuffer(l.next())
			}
		}

		return true, tk
	}
	return false, nil
}

// consume (and skip) whitespaces
func (l *Lexer) consumeWhiteSpace(ch rune) {
	for isWhiteSpace(l.peek()) {
		l.next()
	}
}

// parseIdentifier
func (l *Lexer) parseIdentifier(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	var count = 0
	var terminators = append([]rune{
		EOF, CR, LF, LeftQuoteI, LeftQuoteII, LeftQuoteIII,
		LeftQuoteIV, LeftQuoteV, MiddleDot,
	}, MarkLeads...)

	if !isIdentifierChar(ch, true) {
		return nil, error.NewErrorSLOT("invalid identifier")
	}
	// push first char
	l.pushBuffer(ch)
	count++

	// iterate
	for {
		prev := l.cursor
		ch = l.next()

		if isWhiteSpace(ch) {
			continue
		}
		// if the following chars are a keyword,
		// then terminate the identifier parsing process.
		if isKeyword, _ := l.parseKeyword(ch, false); isKeyword {
			l.rebase(prev)
			return NewIdentifierToken(l.chBuffer), nil
		}
		// parse 注
		if ch == GlyphZHU {
			if validComment, _ := l.validateComment(ch); validComment {
				l.rebase(prev)
				return NewIdentifierToken(l.chBuffer), nil
			}
			l.rebase(prev + 1)
		}
		// other char as terminator
		if util.Contains(ch, terminators) {
			l.rebase(prev)
			return NewIdentifierToken(l.chBuffer), nil
		}

		if isIdentifierChar(ch, false) {
			if count >= maxIdentifierLength {
				return nil, error.NewErrorSLOT("exceed length")
			}
			l.pushBuffer(ch)
			count++
			continue
		}
		return nil, error.NewErrorSLOT("invalid char")
	}
}
