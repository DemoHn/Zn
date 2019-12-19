package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/util"
)

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	lines      *LineScanner
	currentPos int
	peekPos    int
	quoteStack *util.RuneStack
	chBuffer   []rune
	code       []rune // source code
}

// NewLexer - new lexer
func NewLexer(code []rune) *Lexer {
	return &Lexer{
		lines:      NewLineScanner(),
		currentPos: 0,
		peekPos:    0,
		quoteStack: util.NewRuneStack(32),
		chBuffer:   []rune{},
		code:       append(code, EOF),
	}
}

// next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) next() rune {
	if l.peekPos >= len(l.code) {
		return EOF
	}

	data := l.code[l.peekPos]

	l.currentPos = l.peekPos
	l.peekPos++
	return data
}

// peek - get the character of the cursor
func (l *Lexer) peek() rune {
	if l.peekPos >= len(l.code) {
		return EOF
	}
	data := l.code[l.peekPos]

	return data
}

// peek2 - get the next next character without moving the cursor
func (l *Lexer) peek2() rune {
	if l.peekPos >= len(l.code)-1 {
		return EOF
	}
	data := l.code[l.peekPos+1]

	return data
}

// peek3 - get the next next next character without moving the cursor
func (l *Lexer) peek3() rune {
	if l.peekPos >= len(l.code)-2 {
		return EOF
	}
	data := l.code[l.peekPos+2]

	return data
}

// current - get current char value
func (l *Lexer) current() rune {
	if l.peekPos >= len(l.code) {
		return EOF
	}

	return l.code[l.currentPos]
}

// rebase - rebase cursor
func (l *Lexer) rebase(cursor int) {
	l.currentPos = cursor
	l.peekPos = cursor + 1
}

func (l *Lexer) clearBuffer() {
	l.chBuffer = []rune{}
}

func (l *Lexer) pushBuffer(ch rune) {
	l.chBuffer = append(l.chBuffer, ch)
}

// pushBufferRange - push buffer for a range
// @returns bool - if push success
func (l *Lexer) pushBufferRange(start int, end int) bool {
	if end >= len(l.code) || start > end {
		return false
	}

	l.chBuffer = append(l.chBuffer, l.code[start:end+1]...)
	return true
}

// NextToken - parse and generate the next token (including comments)
func (l *Lexer) NextToken() (*Token, *error.Error) {
head:
	var ch = l.next()
	switch ch {
	case EOF:
		l.lines.PushLine(l.currentPos)
		return NewTokenEOF(), nil
	case SP, TAB:
		// if indent has been scanned, it should be regarded as whitespaces
		// (it's totally ignored)
		if !l.lines.HasScanIndent() {
			l.parseIndents(ch)
		} else {
			l.consumeWhiteSpace(ch)
		}
		goto head
	case CR, LF:
		l.parseCRLF(ch)
		l.parseIndents(l.peek())
		goto head
	// meet with 注, it may be possibly a lead character of a comment block
	// notice: it would also be a normal identifer (if 注[number]：) does not satisfy.
	case GlyphZHU:
		cursor := l.currentPos
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
	return l.lines.SetIndent(count, indentType, l.currentPos+1)
}

func (l *Lexer) parseCRLF(ch rune) {
	cursor := l.currentPos
	// for CRLF <windows type> or LFCR
	if (ch == CR && l.peek() == LF) ||
		(ch == LF && l.peek() == CR) {

		// skip one char since we have judge two chars
		l.next()
		l.lines.PushLine(cursor - 1)
		return
	}
	// for LF or CR only
	// LF: <linux>, CR:<old mac>
	l.lines.PushLine(cursor - 1)
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
			l.lines.PushLine(l.currentPos - 1)
			return NewCommentToken(l.chBuffer, isMultiLine), nil
		case CR, LF:
			c1 := l.currentPos
			// parse CR,LF first
			l.parseCRLF(ch)
			if !isMultiLine {
				return NewCommentToken(l.chBuffer, isMultiLine), nil
			}
			// for multi-line comment blocks, CRLF is also included
			l.pushBufferRange(c1, l.currentPos)
			// manually set no indents
			l.lines.SetIndent(0, IdetUnknown, l.currentPos+1)
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
			l.lines.PushLine(l.currentPos - 1)
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
			c1 := l.currentPos
			l.parseCRLF(ch)
			// push buffer & mark new line
			l.pushBufferRange(c1, l.currentPos)
			l.lines.SetIndent(0, IdetUnknown, l.currentPos+1)
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
			l.lines.PushLine(l.currentPos - 1)
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
		l.rebase(l.currentPos - 1)
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
		return NewMarkToken(l.chBuffer, typeCommaSep), nil
	case Colon:
		return NewMarkToken(l.chBuffer, typeFuncCall), nil
	case Semicolon:
		return NewMarkToken(l.chBuffer, typeStmtSep), nil
	case QuestionMark:
		return NewMarkToken(l.chBuffer, typeFuncDeclare), nil
	case RefMark:
		return NewMarkToken(l.chBuffer, typeObjRef), nil
	case BangMark:
		return NewMarkToken(l.chBuffer, typeMustT), nil
	case AnnotationMark:
		return NewMarkToken(l.chBuffer, typeAnnoT), nil
	case HashMark:
		return NewMarkToken(l.chBuffer, typeMapHash), nil
	case EllipsisMark:
		if l.peek() == EllipsisMark {
			l.pushBuffer(l.next())
			return NewMarkToken(l.chBuffer, typeMoreParam), nil
		}
		return nil, error.NewErrorSLOT("invalid ellipsis")
	case LeftBracket:
		return NewMarkToken(l.chBuffer, typeArrayQuoteL), nil
	case RightBracket:
		return NewMarkToken(l.chBuffer, typeArrayQuoteR), nil
	case LeftParen:
		return NewMarkToken(l.chBuffer, typeStmtQuoteL), nil
	case RightParen:
		return NewMarkToken(l.chBuffer, typeStmtQuoteR), nil
	case Equal:
		if l.peek() == Equal {
			l.pushBuffer(l.next())
			return NewMarkToken(l.chBuffer, typeMapData), nil
		}
		return nil, error.NewErrorSLOT("invalid single euqal")
	case DoubleArrow:
		return NewMarkToken(l.chBuffer, typeMapData), nil
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
		tk = NewKeywordToken(typeDeclareW)
	case GlyphWEI:
		tk = NewKeywordToken(typeLogicYesW)
	case GlyphSHI:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(typeObjConstructW)
		} else {
			tk = NewKeywordToken(typeLogicYesIIW)
		}
	case GlyphSHE:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(typeAssignW)
		} else {
			return false, nil
		}
	case GlyphRU:
		switch l.peek() {
		case GlyphGUO:
			wordLen = 2
			tk = NewKeywordToken(typeCondW)
		case GlyphHE:
			wordLen = 2
			tk = NewKeywordToken(typeFuncW)
		default:
			return false, nil
		}
	case GlyphHE:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(typeStaticFuncW)
		} else {
			return false, nil
		}
	case GlyphYI:
		if l.peek() == GlyphZHIy {
			wordLen = 2
			tk = NewKeywordToken(typeParamAssignW)
		} else {
			return false, nil
		}
	case GlyphFAN:
		if l.peek() == GlyphHUI {
			wordLen = 2
			tk = NewKeywordToken(typeReturnW)
		} else {
			return false, nil
		}
	case GlyphBU:
		switch l.peek() {
		case GlyphWEI:
			wordLen = 2
			tk = NewKeywordToken(typeLogicNotW)
		case GlyphSHI:
			wordLen = 2
			tk = NewKeywordToken(typeLogicNotIIW)
		case GlyphDENG:
			if l.peek2() == GlyphYU {
				wordLen = 3
				tk = NewKeywordToken(typeLogicNotEqW)
			} else {
				return false, nil
			}
		case GlyphDA:
			if l.peek2() == GlyphYU {
				wordLen = 3
				tk = NewKeywordToken(typeLogicLteW)
			} else {
				return false, nil
			}
		case GlyphXIAO:
			if l.peek2() == GlyphYU {
				wordLen = 3
				tk = NewKeywordToken(typeLogicGteW)
			} else {
				return false, nil
			}
		}
	case GlyphDENG:
		if l.peek() == GlyphYU {
			wordLen = 2
			tk = NewKeywordToken(typeLogicEqualW)
		} else {
			return false, nil
		}
	case GlyphDA:
		if l.peek() == GlyphYU {
			wordLen = 2
			tk = NewKeywordToken(typeLogicGtW)
		} else {
			return false, nil
		}
	case GlyphXIAO:
		if l.peek() == GlyphYU {
			wordLen = 2
			tk = NewKeywordToken(typeLogicLtW)
		} else {
			return false, nil
		}
	case GlyphYIi:
		tk = NewKeywordToken(typeVarOneW)
	case GlyphER:
		tk = NewKeywordToken(typeVarTwoW)
	case GlyphDE:
		tk = NewKeywordToken(typeFuncYieldW)
	case GlyphFOU:
		if l.peek() == GlyphZE {
			wordLen = 2
			tk = NewKeywordToken(typeCondElseW)
		} else {
			return false, nil
		}
	case GlyphMEI:
		if l.peek() == GlyphDANG {
			wordLen = 2
			tk = NewKeywordToken(typeWhileLoopW)
		} else {
			return false, nil
		}
	case GlyphCHENG:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(typeObjNewW)
		} else {
			return false, nil
		}
	case GlyphZUO:
		if l.peek() == GlyphWEI {
			wordLen = 2
			tk = NewKeywordToken(typeVarAliasW)
		} else {
			return false, nil
		}
	case GlyphDING:
		if l.peek() == GlyphYIy {
			wordLen = 2
			tk = NewKeywordToken(typeObjDefineW)
		} else {
			return false, nil
		}
	case GlyphLEI:
		if l.peek() == GlyphBI {
			wordLen = 2
			tk = NewKeywordToken(typeObjTraitW)
		} else {
			return false, nil
		}
	case GlyphQI:
		tk = NewKeywordToken(typeObjThisW)
	case GlyphCI:
		tk = NewKeywordToken(typeObjSelfW)
	case GlyphZAI:
		tk = NewKeywordToken(typeFuncCallOneW)
	case GlyphZHONG:
		tk = NewKeywordToken(typeFuncCallTwoW)
	case GlyphHUO:
		tk = NewKeywordToken(typeLogicOrW)
	case GlyphQIE:
		tk = NewKeywordToken(typeLogicAndW)
	case GlyphZHI:
		tk = NewKeywordToken(typeObjDotW)
	case GlyphDEo:
		tk = NewKeywordToken(typeObjDotIIW)
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
		prev := l.currentPos
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
