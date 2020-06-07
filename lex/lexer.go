package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/util"
)

const (
	defaultBlockSize int = 512
)

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	*InputStream // input stream
	*LineStack
	quoteStack *util.RuneStack
	chBuffer   []rune // the buffer for parsing & generating tokens
	cursor     int
	blockSize  int
	beginLex   bool
}

// NewLexer - new lexer
func NewLexer(in *InputStream) *Lexer {
	return &Lexer{
		LineStack:   NewLineStack(in),
		quoteStack:  util.NewRuneStack(32),
		InputStream: in,
		chBuffer:    []rune{},
		cursor:      -1,
		blockSize:   defaultBlockSize,
		beginLex:    true,
	}
}

// next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) next() rune {
	l.cursor++

	if l.cursor+2 >= l.getLineBufferSize() {
		if !l.End() {
			if b, err := l.Read(l.blockSize); err == nil {
				l.AppendLineBuffer(b)
			} else {
				// throw the error globally
				// it will be handled (recovered) in NextToken(),
				// similiar with C++'s try-catch statement.
				panic(err)
			}
		}
	}

	// still no data, return EOF directly
	return l.getChar(l.cursor)
}

// peek - get the character of the cursor
func (l *Lexer) peek() rune {
	return l.getChar(l.cursor + 1)
}

// peek2 - get the next next character without moving the cursor
func (l *Lexer) peek2() rune {
	return l.getChar(l.cursor + 2)
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

// SetBlockSize -
func (l *Lexer) SetBlockSize(size int) {
	l.blockSize = size
}

// NextToken - parse and generate the next token (including comments)
func (l *Lexer) NextToken() (tok *Token, err *error.Error) {
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(*error.Error)
			// for other kinds of error (e.g. runtime error), panic it directly
			if !ok {
				panic(r)
			}
		}
		handleDeferError(l, err)
	}()

	// For the first line, we use some tricks to determine if this line
	// contains indents
	if l.beginLex {
		l.beginLex = false
		if !util.Contains(l.peek(), []rune{SP, TAB, EOF}) {
			if err := l.SetIndent(0, IdetUnknown); err != nil {
				return nil, err
			}
		}
	}
head:
	var ch = l.next()
	switch ch {
	case EOF:
		l.PushLine(l.cursor)
		tok = NewTokenEOF(l.CurrentLine, l.cursor)
		return
	case SP, TAB:
		// if indent has been scanned, it should be regarded as whitespaces
		// (it's totally ignored)
		if l.onIndentStage() {
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
		isComment, isMultiLine, note := l.validateComment(ch)
		if isComment {
			tok, err = l.parseComment(l.getChar(l.cursor), isMultiLine, note)
			return
		}

		l.rebase(cursor)
		// goto normal identifier
	// left quotes
	case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
		tok, err = l.parseString(ch)
		return
	case MiddleDot:
		tok, err = l.parseVarQuote(ch)
		return
	default:
		// skip whitespaces
		if isWhiteSpace(ch) {
			l.consumeWhiteSpace(ch)
			goto head
		}
		// parse number
		if isNumber(ch) || util.Contains(ch, []rune{'.', '+', '-'}) {
			tok, err = l.parseNumber(ch)
			return
		}
		if util.Contains(ch, MarkLeads) {
			tok, err = l.parseMarkers(ch)
			return
		}
		// suppose it's a keyword
		if isKeyword, tk := l.parseKeyword(ch, true); isKeyword {
			tok = tk
			return
		}
	}
	tok, err = l.parseIdentifier(ch)
	return
}

func handleDeferError(l *Lexer, err *error.Error) {
	if err != nil {
		if err.GetErrorClass() == error.IOErrorClass {
			// For I/O error, load current line buffer directly
			// instead of moving cursor to line end (since it's impossible to retrieve line end)
			err.SetCursor(error.Cursor{
				File:    l.InputStream.Scope,
				LineNum: l.CurrentLine,
				Text:    l.GetLineText(l.CurrentLine, false),
				ColNum:  0,
			})
		} else {
			l.moveAndSetCursor(err)
		}
	}
}

//// parsing logics
func (l *Lexer) parseIndents(ch rune) {
	count := 1
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
	if err := l.SetIndent(count, indentType); err != nil {
		panic(err)
	}
}

// parseCRLF and return the newline chars by the way
func (l *Lexer) parseCRLF(ch rune) []rune {
	var rtn = []rune{}
	p := l.peek()
	// for CRLF <windows type> or LFCR
	if (ch == CR && p == LF) || (ch == LF && p == CR) {
		// skip one char since we have judge two chars
		l.next()
		l.PushLine(l.cursor - 1)

		rtn = []rune{ch, p}
	} else {
		// for LF or CR only
		// LF: <linux>, CR:<old mac>
		l.PushLine(l.cursor)
		rtn = []rune{ch}
	}

	// new line and reset cursor
	l.NewLine(l.cursor + 1)

	// to see if next line contains (potential) indents
	if !util.Contains(l.peek(), []rune{SP, TAB, EOF}) {
		if err := l.SetIndent(0, IdetUnknown); err != nil {
			panic(err)
		}
	}
	return rtn
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
func (l *Lexer) validateComment(ch rune) (bool, bool, []rune) {
	note := []rune{}
	// “ or 「
	lquotes := []rune{LeftQuoteIV, LeftQuoteII}
	for {
		ch = l.next()
		// match pattern 1, 2
		if ch == Colon {
			// match pattern 3, 4
			if util.Contains(l.peek(), lquotes) {
				l.next()
				return true, true, note
			}
			return true, false, note
		}
		if isNumber(ch) || isWhiteSpace(ch) {
			note = append(note, ch)
			continue
		}
		return false, false, note
	}
}

// parseComment until its end
func (l *Lexer) parseComment(ch rune, isMultiLine bool, note []rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	if isMultiLine {
		if !l.quoteStack.Push(ch) {
			return nil, error.QuoteStackFull(l.quoteStack.GetMaxSize())
		}
		l.pushBuffer(ch)
	}
	rg := newTokenRange(l)
	// iterate
	for {
		ch = l.next()
		switch ch {
		case EOF:
			l.rebase(l.cursor - 1)
			rg.setRangeEnd(l)
			return NewCommentToken(l.chBuffer, note, rg), nil
		case CR, LF:
			// parse CR,LF first
			nl := l.parseCRLF(ch)
			if !isMultiLine {
				return NewCommentToken(l.chBuffer, note, rg), nil
			}
			// for multi-line comment blocks, CRLF is also included
			l.pushBuffer(nl...)

			// manually set no indents
			if err := l.SetIndent(0, IdetUnknown); err != nil {
				return nil, err
			}
		default:
			// for mutli-line comment, calculate quotes is necessary.
			if isMultiLine {
				// push left quotes
				if util.Contains(ch, LeftQuotes) {
					if !l.quoteStack.Push(ch) {
						return nil, error.QuoteStackFull(l.quoteStack.GetMaxSize())
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
						rg.setRangeEnd(l)
						l.pushBuffer(ch)
						return NewCommentToken(l.chBuffer, note, rg), nil
					}
				}
			}
			l.pushBuffer(ch)
			// cache rangeEnd location as the position of last char
			// Example: (for single-line comment token)
			// 注： ABCDEFG \r\n
			//
			// Here, the rangeEnd should be on char `G`, but it will stop iff the following `\r\n` is parsed!
			rg.setRangeEnd(l)
		}
	}
}

// parseString -
func (l *Lexer) parseString(ch rune) (*Token, *error.Error) {
	// start up
	l.clearBuffer()
	l.quoteStack.Push(ch)
	firstChar := ch
	rg := newTokenRange(l)

	// iterate
	for {
		ch := l.next()
		switch ch {
		case EOF:
			l.rebase(l.cursor - 1)
			// after meeting with EOF
			rg.setRangeEnd(l)
			return NewStringToken(l.chBuffer, firstChar, rg), nil
		// push quotes
		case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
			l.pushBuffer(ch)
			if !l.quoteStack.Push(ch) {
				return nil, error.QuoteStackFull(l.quoteStack.GetMaxSize())
			}
		// pop quotes if match
		case RightQuoteI, RightQuoteII, RightQuoteIII, RightQuoteIV, RightQuoteV:
			currentL, _ := l.quoteStack.Current()
			if QuoteMatchMap[currentL] == ch {
				l.quoteStack.Pop()
			}
			// stop quoting
			if l.quoteStack.IsEmpty() {
				rg.setRangeEnd(l)
				return NewStringToken(l.chBuffer, firstChar, rg), nil
			}
			l.pushBuffer(ch)
		case CR, LF:

			nl := l.parseCRLF(ch)
			// push buffer & mark new line
			l.pushBuffer(nl...)
			if err := l.SetIndent(0, IdetUnknown); err != nil {
				return nil, err
			}
		default:
			l.pushBuffer(ch)
		}
	}
}

func (l *Lexer) parseVarQuote(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	rg := newTokenRange(l)
	// iterate
	count := 0
	for {
		ch = l.next()
		// we should ensure the following chars to satisfy the condition
		// of an identifier
		switch ch {
		case EOF:
			l.rebase(l.cursor - 1)
			rg.setRangeEnd(l)
			return NewVarQuoteToken(l.chBuffer, rg), nil
		case MiddleDot:
			rg.setRangeEnd(l)
			return NewVarQuoteToken(l.chBuffer, rg), nil
		default:
			// ignore white-spaces!
			if isWhiteSpace(ch) {
				continue
			}
			if isIdentifierChar(ch, count == 0) {
				l.pushBuffer(ch)
				count++
				if count > maxIdentifierLength {
					return nil, error.IdentifierExceedLength(maxIdentifierLength)
				}
			} else {
				return nil, error.InvalidIdentifier()
			}
		}
	}
}

// regex: ^[Ff]?[-+]?[0-9]*\.?[0-9]+((([eE][-+])|(\*(10)?\^[-+]?))[0-9]+)?$
// ref: https://github.com/DemoHn/Zn/issues/4
func (l *Lexer) parseNumber(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	rg := newTokenRange(l)

	// hand-written regex parser
	// ref: https://cyberzhg.github.io/toolbox/min_dfa?regex=Rj9QP0QqLj9EKygoKEVQKXwocygxMCk/dVA/KSlEKyk/
	// hand-drawn min-DFA:
	// https://github.com/DemoHn/Zn/issues/6
	const (
		sBegin      = 1
		sDot        = 2
		sIntEnd     = 3
		sIntPMFlag  = 5
		sDotDecEnd  = 6
		sEFlag      = 7
		sSFlag      = 8
		sExpPMFlag  = 9
		sSciI       = 10
		sSciEndFlag = 11
		sExpEnd     = 12
		sSciII      = 13
	)
	var state = sBegin
	var endStates = []int{sIntEnd, sDotDecEnd, sExpEnd}

	for {
		switch ch {
		case EOF:
			goto end
		case 'e', 'E':
			switch state {
			case sDotDecEnd, sIntEnd:
				state = sEFlag
			default:
				goto end
			}
		case '.':
			switch state {
			case sBegin, sIntPMFlag, sIntEnd:
				state = sDot
			default:
				goto end
			}
		case '-', '+':
			switch state {
			case sBegin:
				state = sIntPMFlag
			case sEFlag, sSciEndFlag:
				state = sExpPMFlag
			default:
				goto end
			}
		case '_':
			ch = l.next()
			continue
		case '*':
			switch state {
			case sDotDecEnd, sIntEnd:
				state = sSFlag
			default:
				goto end
			}
		case '1':
			switch state {
			case sSFlag:
				state = sSciI
				// same with other numbers
			case sBegin, sIntEnd, sIntPMFlag:
				state = sIntEnd
			case sDot, sDotDecEnd:
				state = sDotDecEnd
			case sExpPMFlag, sSciEndFlag, sExpEnd:
				state = sExpEnd
			default:
				goto end
			}
		case '0':
			switch state {
			case sSciI:
				state = sSciII
			case sBegin, sIntEnd, sIntPMFlag:
				state = sIntEnd
			case sDot, sDotDecEnd:
				state = sDotDecEnd
			case sExpPMFlag, sSciEndFlag, sExpEnd:
				state = sExpEnd
			default:
				goto end
			}
		case '2', '3', '4', '5', '6', '7', '8', '9':
			switch state {
			case sBegin, sIntEnd, sIntPMFlag:
				state = sIntEnd
			case sDot, sDotDecEnd:
				state = sDotDecEnd
			case sExpPMFlag, sSciEndFlag, sExpEnd:
				state = sExpEnd
			default:
				goto end
			}
		case '^':
			switch state {
			case sSFlag, sSciII:
				state = sSciEndFlag
			default:
				goto end
			}
		default:
			goto end
		}
		l.pushBuffer(ch)
		ch = l.next()
	}

end:
	if util.ContainsInt(state, endStates) {
		// back to last available char
		l.rebase(l.cursor - 1)
		rg.setRangeEnd(l)
		return NewNumberToken(l.chBuffer, rg), nil
	}
	return nil, error.InvalidChar(ch)
}

// parseMarkers -
func (l *Lexer) parseMarkers(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	l.pushBuffer(ch)

	startR := newTokenRange(l)
	// switch
	switch ch {
	case Comma:
		return NewMarkToken(l.chBuffer, TypeCommaSep, startR, 1), nil
	case Colon:
		return NewMarkToken(l.chBuffer, TypeFuncCall, startR, 1), nil
	case Semicolon:
		return NewMarkToken(l.chBuffer, TypeStmtSep, startR, 1), nil
	case QuestionMark:
		return NewMarkToken(l.chBuffer, TypeFuncDeclare, startR, 1), nil
	case RefMark:
		return NewMarkToken(l.chBuffer, TypeObjRef, startR, 1), nil
	case BangMark:
		return NewMarkToken(l.chBuffer, TypeMustT, startR, 1), nil
	case AnnotationMark:
		return NewMarkToken(l.chBuffer, TypeAnnoT, startR, 1), nil
	case HashMark:
		if l.peek() == LeftCurlyBracket {
			l.pushBuffer(l.next())
			return NewMarkToken(l.chBuffer, TypeMapQHash, startR, 2), nil
		}
		return NewMarkToken(l.chBuffer, TypeMapHash, startR, 1), nil
	case EllipsisMark:
		if l.peek() == EllipsisMark {
			l.pushBuffer(l.next())
			return NewMarkToken(l.chBuffer, TypeMoreParam, startR, 2), nil
		}
		return nil, error.InvalidSingleEllipsis()
	case LeftBracket:
		return NewMarkToken(l.chBuffer, TypeArrayQuoteL, startR, 1), nil
	case RightBracket:
		return NewMarkToken(l.chBuffer, TypeArrayQuoteR, startR, 1), nil
	case LeftParen:
		return NewMarkToken(l.chBuffer, TypeFuncQuoteL, startR, 1), nil
	case RightParen:
		return NewMarkToken(l.chBuffer, TypeFuncQuoteR, startR, 1), nil
	case LeftCurlyBracket:
		return NewMarkToken(l.chBuffer, TypeStmtQuoteL, startR, 1), nil
	case RightCurlyBracket:
		return NewMarkToken(l.chBuffer, TypeStmtQuoteR, startR, 1), nil
	case Equal:
		if l.peek() == Equal {
			l.pushBuffer(l.next())
			return NewMarkToken(l.chBuffer, TypeMapData, startR, 2), nil
		}
		return nil, error.InvalidSingleEqual()
	case DoubleArrow:
		return NewMarkToken(l.chBuffer, TypeMapData, startR, 1), nil
	}
	return nil, error.InvalidChar(ch)
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

	rg := newTokenRange(l)
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
			return false, nil
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
	case GlyphZAI:
		if l.peek() == GlyphRU {
			wordLen = 2
			tk = NewKeywordToken(TypeCondOtherW)
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
	}

	if tk != nil {
		if moveForward {
			switch wordLen {
			case 1:
				l.pushBuffer(ch)
			case 2:
				l.pushBuffer(ch, l.next())
			case 3:
				l.pushBuffer(ch, l.next(), l.next())
			}
		}

		//rg.EndLine = rg.StartLine
		//rg.EndCol = rg.StartCol + wordLen - 1
		rg.EndLine = rg.StartLine
		rg.EndIdx = rg.StartIdx + wordLen
		tk.Range = rg
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
		return nil, error.InvalidIdentifier()
	}

	rg := newTokenRange(l)
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
			rg.setRangeEnd(l)
			return NewIdentifierToken(l.chBuffer, rg), nil
		}
		// parse 注
		if ch == GlyphZHU {
			if validComment, _, _ := l.validateComment(ch); validComment {
				l.rebase(prev)
				return NewIdentifierToken(l.chBuffer, rg), nil
			}
			l.rebase(prev + 1)
		}
		// other char as terminator
		if util.Contains(ch, terminators) {
			l.rebase(prev)
			return NewIdentifierToken(l.chBuffer, rg), nil
		}

		if isIdentifierChar(ch, false) {
			if count >= maxIdentifierLength {
				return nil, error.IdentifierExceedLength(maxIdentifierLength)
			}
			l.pushBuffer(ch)
			rg.setRangeEnd(l)
			count++
			continue
		}
		return nil, error.InvalidChar(ch)
	}
}

// moveAndSetCursor - retrieve full text of the line and set the current cursor
// to display errors
func (l *Lexer) moveAndSetCursor(err *error.Error) {
	cursor := error.Cursor{
		File:    l.InputStream.Scope,
		ColNum:  l.cursor - l.scanCursor.startIdx,
		LineNum: l.CurrentLine,
		Text:    l.GetLineText(l.CurrentLine, true),
	}
	err.SetCursor(cursor)
}
