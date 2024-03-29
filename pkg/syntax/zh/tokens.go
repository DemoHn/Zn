package zh

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

//// 1. markers
// declare marks
const (
	Comma             rune = 0xFF0C // ，
	PauseComma        rune = 0x3001 // 、
	Colon             rune = 0xFF1A // ：
	Semicolon         rune = 0xFF1B // ；
	QuestionMark      rune = 0xFF1F // ？
	RefMark           rune = 0x0026 // &
	BangMark          rune = 0xFF01 // ！
	AnnotationMark    rune = 0x0040 // @
	HashMark          rune = 0x0023 // #
	LeftBracket       rune = 0x3010 // 【
	RightBracket      rune = 0x3011 // 】
	LeftParen         rune = 0xFF08 // （
	RightParen        rune = 0xFF09 // ）
	Equal             rune = 0x003D // =
	LeftCurlyBracket  rune = 0x007B // {
	RightCurlyBracket rune = 0x007D // }
	LessThanMark      rune = 0x003C // <
	GreaterThanMark   rune = 0x003E // >
	PlusMark          rune = '+'    // +
	MinusMark         rune = '-'    // -
	MultiplyMark      rune = '*'    // *
	Slash             rune = 0x002F // /
	CharZHU           rune = 0x6CE8 // 注
)

//// 2. quotes
// declare quotes
const (
	LeftLibQuoteI      rune = 0x300A //《
	RightLibQuoteI     rune = 0x300B // 》
	LeftDoubleQuoteI   rune = 0x300C // 「
	RightDoubleQuoteI  rune = 0x300D // 」
	LeftDoubleQuoteII  rune = 0x201C // “
	RightDoubleQuoteII rune = 0x201D // “
	LeftSingleQuoteI   rune = 0x300E // 『
	RightSingleQuoteI  rune = 0x300F // 』
	LeftSingleQuoteII  rune = 0x2018 // ‘
	RightSingleQuoteII rune = 0x2019 // ‘
)

//// 3. commentType
const (
	commentTypeSingle  = 1 // single line
	commentTypeSlash   = 2 // multiple line, starts with '/*'
	commentTypeQuoteI  = 3 // multiple line, starts with '注：「'
	commentTypeQuoteII = 4 // multiple line, starts with '注：“'
)

//// 4. var quote
const (
	MiddleDot rune = 0x00B7 // ·
	BackTick  rune = 0x0060 // `
)

//// token constants and constructors (without keyword token)
// token types -
// for special type Tokens, its range varies from 0 - 9
// for keyword types, check lex/keyword.go for details
const (
	TypeEOF           uint8 = 0
	TypeString        uint8 = 2  // string (only double quotes)
	TypeNumber        uint8 = 4  // numbers
	TypeIdentifier    uint8 = 5  //
	TypeEnumString    uint8 = 6  // string (with single quotes)
	TypeLibString     uint8 = 7  // string (with guillemots)
	TypeComment       uint8 = 10 // 注：
	TypeCommaSep      uint8 = 11 // ，
	TypeStmtSep       uint8 = 12 // ；
	TypeFuncCall      uint8 = 13 // ：
	TypeFuncDeclare   uint8 = 14 // ？
	TypeObjRef        uint8 = 15 // &
	TypeExceptionT    uint8 = 16 // ！
	TypeAnnotationT   uint8 = 17 // @
	TypeMapHash       uint8 = 18 // #
	TypeArrayQuoteL   uint8 = 20 // 【
	TypeArrayQuoteR   uint8 = 21 // 】
	TypeFuncQuoteL    uint8 = 22 // （
	TypeFuncQuoteR    uint8 = 23 // ）
	TypeStmtQuoteL    uint8 = 25 // {
	TypeStmtQuoteR    uint8 = 26 // }
	TypePauseCommaSep uint8 = 28 // 、
	TypeAssignMark    uint8 = 29 // =
	TypeGTMark        uint8 = 30 // >
	TypeLTMark        uint8 = 31 // <
	TypeGTEMark       uint8 = 32 // >=
	TypeLTEMark       uint8 = 33 // <=
	TypeNEMark        uint8 = 34 // /=
	TypeEqualMark     uint8 = 35 // ==
	TypePlus          uint8 = 36 // +
	TypeMinus         uint8 = 37 // -
	TypeMultiply      uint8 = 38 // *
	TypeDivision      uint8 = 39 // /
	//// from 40 - 78, reserved for keywords
)

// MarkLeads -
var MarkLeads = []rune{
	Comma, PauseComma, Colon, Semicolon, QuestionMark, RefMark, BangMark,
	AnnotationMark, HashMark, LeftBracket,
	RightBracket, LeftParen, RightParen, Equal,
	LeftCurlyBracket, RightCurlyBracket, LessThanMark, GreaterThanMark,
}

// NextToken -
func NextToken(l *syntax.Lexer) (syntax.Token, error) {
	// parse non-keyword tokens e.g.: Spaces, LineBreaks
	if err := l.PreNextToken(); err != nil {
		return syntax.Token{}, err
	}

	ch := l.GetCurrentChar()
	switch ch {
	case syntax.RuneEOF:
		return syntax.Token{Type: TypeEOF, StartIdx: l.GetCursor(), EndIdx: l.GetCursor()}, nil
	// handle 'A + B' case
	// for numbers like '+1234', this will be handled by parseNumber()
	case PlusMark, MinusMark, MultiplyMark:
		startIdx := l.GetCursor()
		chn := l.Peek()

		t := TypePlus
		if ch == MinusMark {
			t = TypeMinus
		} else if ch == MultiplyMark {
			t = TypeMultiply
		}
		// NOTE: the next char must be space to ensure it's not a part of
		// identifier
		if syntax.IsWhiteSpace(chn) {
			l.Next()
			return syntax.Token{Type: t, StartIdx: startIdx, EndIdx: l.GetCursor()}, nil
		}
	case Slash:
		startIdx := l.GetCursor()
		chn := l.Peek()

		// parse /=, example usage: '如果 X /= 10'
		if chn == Equal {
			l.Next()
			l.Next()
			return syntax.Token{
				Type:     TypeNEMark,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
			}, nil
		}

		// parse / (as div) only, example usage: '25 / 8'
		if syntax.IsWhiteSpace(chn) {
			l.Next()
			return syntax.Token{
				Type:     TypeDivision,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
			}, nil
		}

		// maybe / would be the start of comment block (e.g. // This is a comment)
		isComment, tk, err := parseComment(l)
		if err != nil {
			return syntax.Token{}, err
		}
		if isComment {
			return tk, nil
		}
	case CharZHU:
		isComment, tk, err := parseComment(l)
		if err != nil {
			return syntax.Token{}, err
		}
		if isComment {
			return tk, nil
		}
	case LeftLibQuoteI, LeftDoubleQuoteI, LeftDoubleQuoteII, LeftSingleQuoteI, LeftSingleQuoteII:
		return parseString(l)
	case MiddleDot, BackTick:
		return parseVarQuote(l)
	}

	// other token types
	if isNumber(ch) {
		return parseNumber(l)
	}
	if syntax.ContainsRune(ch, MarkLeads) {
		return parseMarkers(l)
	}

	// suppose it's a keyword
	isKeyword, tk, err := parseKeyword(l, true)
	if err != nil {
		return syntax.Token{}, err
	}
	if isKeyword {
		return tk, nil
	}

	return parseIdentifier(l)
}

// regex: ^[-+]?[0-9]*\.?[0-9]+((([eE][-+])|(\*(10)?\^[-+]?))[0-9]+)?$
// ref: https://github.com/DemoHn/Zn/issues/4
func parseNumber(l *syntax.Lexer) (syntax.Token, error) {
	ch := l.GetCurrentChar()
	startIdx := l.GetCursor()

	var literal []rune
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
		case syntax.RuneEOF:
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
		case ',':
			ch = l.Next()
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
		// add item to literal
		literal = append(literal, ch)
		ch = l.Next()
	}

end:
	if syntax.ContainsInt(state, endStates) {
		return syntax.Token{
			Type:     TypeNumber,
			Literal:  literal,
			StartIdx: startIdx,
			EndIdx:   l.GetCursor(),
		}, nil
	}
	return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
}

func parseMarkers(l *syntax.Lexer) (syntax.Token, error) {
	startIdx := l.GetCursor()
	ch := l.GetCurrentChar()
	var tokenType uint8

	switch ch {
	case Comma:
		tokenType = TypeCommaSep
	case PauseComma:
		tokenType = TypePauseCommaSep
	case Colon:
		tokenType = TypeFuncCall
	case Semicolon:
		tokenType = TypeStmtSep
	case QuestionMark:
		tokenType = TypeFuncDeclare
	case RefMark:
		tokenType = TypeObjRef
	case BangMark:
		tokenType = TypeExceptionT
	case AnnotationMark:
		tokenType = TypeAnnotationT
	case HashMark:
		tokenType = TypeMapHash
	case LeftBracket:
		tokenType = TypeArrayQuoteL
	case RightBracket:
		tokenType = TypeArrayQuoteR
	case LeftParen:
		tokenType = TypeFuncQuoteL
	case RightParen:
		tokenType = TypeFuncQuoteR
	case LeftCurlyBracket:
		tokenType = TypeStmtQuoteL
	case RightCurlyBracket:
		tokenType = TypeStmtQuoteR
	case Equal:
		if l.Peek() == Equal {
			l.Next()
			tokenType = TypeEqualMark
		} else {
			tokenType = TypeAssignMark
		}
	case LessThanMark:
		if l.Peek() == Equal {
			l.Next()
			tokenType = TypeLTEMark
		} else {
			tokenType = TypeLTMark
		}
	case GreaterThanMark:
		if l.Peek() == Equal {
			l.Next()
			tokenType = TypeGTEMark
		} else {
			tokenType = TypeGTMark
		}
	default:
		return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
	}
	// include all necessary chars
	l.Next()

	return syntax.Token{
		Type:     tokenType,
		StartIdx: startIdx,
		EndIdx:   l.GetCursor(),
	}, nil
}

// parse · <identifier> ·
func parseVarQuote(l *syntax.Lexer) (syntax.Token, error) {
	startIdx := l.GetCursor()
	var literal []rune
	for {
		ch := l.Next()
		switch ch {
		case syntax.RuneEOF:
			return syntax.Token{
				Type:     TypeIdentifier,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
				Literal:  literal,
			}, nil
		case MiddleDot, BackTick:
			l.Next()
			return syntax.Token{
				Type:     TypeIdentifier,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
				Literal:  literal,
			}, nil
		default:
			literal = append(literal, ch)
		}
	}
}

// 4 types of string:
// 1. 「 ... 」 or “ ... ”
// 2. 『 ... 』 or ‘ ... ‘
// 3. 《 ... 》
func parseString(l *syntax.Lexer) (syntax.Token, error) {
	sch := l.GetCurrentChar()
	startIdx := l.GetCursor()
	literal := []rune{}

	quoteNum := 1
	tkType := TypeString
	quoteMatchMap := map[rune]rune{
		LeftDoubleQuoteI:  RightDoubleQuoteI,
		LeftDoubleQuoteII: RightDoubleQuoteII,
		LeftSingleQuoteI:  RightSingleQuoteI,
		LeftSingleQuoteII: RightSingleQuoteII,
		LeftLibQuoteI:     RightLibQuoteI,
	}

	// get token type
	if sch == LeftSingleQuoteI || sch == LeftSingleQuoteII {
		tkType = TypeEnumString
	} else if sch == LeftLibQuoteI {
		tkType = TypeLibString
	}

	for {
		ch := l.Next()
		switch ch {
		case syntax.RuneEOF:
			return syntax.Token{
				Type:     tkType,
				Literal:  literal,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
			}, nil
		case syntax.RuneCR, syntax.RuneLF:
			p := l.Peek()
			if (ch == syntax.RuneCR && p == syntax.RuneLF) || (ch == syntax.RuneLF && p == syntax.RuneCR) {
				literal = append(literal, ch)
				l.Next()
			}
			l.Lines = append(l.Lines, syntax.LineInfo{
				Indents:  0,
				StartIdx: l.GetCursor() + 1,
			})
			// add literal (for CR/LF only, append oneChar; for CR+LF, append LF)
			literal = append(literal, l.GetCurrentChar())
		case LeftDoubleQuoteI, LeftDoubleQuoteII, LeftSingleQuoteI, LeftSingleQuoteII, LeftLibQuoteI:
			if sch == ch {
				quoteNum += 1
			}
			literal = append(literal, ch)
		case RightDoubleQuoteI, RightDoubleQuoteII, RightSingleQuoteI, RightSingleQuoteII, RightLibQuoteI:
			if quoteMatchMap[sch] == ch {
				quoteNum -= 1
				if quoteNum == 0 {
					// return strings
					l.Next()
					return syntax.Token{
						Type:     tkType,
						Literal:  literal,
						StartIdx: startIdx,
						EndIdx:   l.GetCursor(),
					}, nil
				}
			}
			literal = append(literal, ch)
		default:
			literal = append(literal, ch)
		}
	}
}

// parseIdentifier -
func parseIdentifier(l *syntax.Lexer) (syntax.Token, error) {
	ch := l.GetCurrentChar()
	startIdx := l.GetCursor()

	literal := []rune{ch}

	terminators := append([]rune{
		syntax.RuneEOF, syntax.RuneCR, syntax.RuneLF,
		LeftLibQuoteI, LeftSingleQuoteI, LeftSingleQuoteII,
		LeftDoubleQuoteI, LeftDoubleQuoteII,
	}, MarkLeads...)

	// 0. first char must be an identifier
	if !isIdentifierChar(ch) {
		return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
	}
	for {
		ch = l.Next()
		// 1. when next char is space, stop here
		if syntax.IsWhiteSpace(ch) {
			return syntax.Token{
				Type:     TypeIdentifier,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
				Literal:  literal,
			}, nil
		}

		// 2. when next char is a part of keyword, stop here
		isKeyword, _, err := parseKeyword(l, false)
		if err != nil {
			return syntax.Token{}, err
		}
		if isKeyword {
			return syntax.Token{
				Type:     TypeIdentifier,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
				Literal:  literal,
			}, nil
		}
		// 3. when next char is a start of comment, stop here
		// only 「//」 and 「/*」 type is available
		// NOTE: we will regard comment type「注」 as a regular identifier
		if ch == Slash && syntax.ContainsRune(l.Peek(), []rune{Slash, MultiplyMark, Equal}) {
			return syntax.Token{
				Type:     TypeIdentifier,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
				Literal:  literal,
			}, nil
		}
		// 4. when next char is a mark, stop here
		if syntax.ContainsRune(ch, terminators) {
			return syntax.Token{
				Type:     TypeIdentifier,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
				Literal:  literal,
			}, nil
		}
		// 5. otherwise, if it's an identifier with +, -, *, /, .
		// add char to literal
		if isIdentifierChar(ch) || syntax.ContainsRune(ch, []rune{PlusMark, MinusMark, MultiplyMark, Slash, '.'}) {
			literal = append(literal, ch)
			continue
		}
		return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
	}
}

// validate if the coming block is a comment block then parse comment block
// valid comment block are listed below:
// (single-line)
// 1. 注：...
// 2. 注123456：...
// 3. // ...
//
// (multi-line)
//
// 4. 注：“ or 「 ...
// 5. 注123456：“ or 「 ...
// 6. /* ...
//
// @returns (isValid, Token, error)
func parseComment(l *syntax.Lexer) (bool, syntax.Token, error) {
	startIdx := l.GetCursor()
	ch := l.GetCurrentChar()
	// for multi-line comment, we have to ensure all left quotes and right quotes
	// shall be matched
	quoteCount := 0
	isComment := false
	var multiCommentType int
	switch ch {
	case CharZHU:
		// parse number marks; e.g. 注123456：
		for {
			if !isPureNumber(l.Next()) {
				break
			}
		}

		// parse ：after 「注」
		if l.GetCurrentChar() == Colon {
			isComment = true
			switch l.Next() {
			case LeftDoubleQuoteI:
				multiCommentType = commentTypeQuoteI
				quoteCount = 1
			case LeftDoubleQuoteII:
				multiCommentType = commentTypeQuoteII
				quoteCount = 1
			default:
				multiCommentType = commentTypeSingle
			}
		} else {
			return false, syntax.Token{}, zerr.InvalidChar(l.GetCurrentChar(), l.GetCursor())
		}
	case Slash:
		p := l.Peek()
		if p == Slash {
			l.Next()
			// single line comment
			isComment = true
			multiCommentType = commentTypeSingle
		} else if p == MultiplyMark {
			l.Next()
			// multiple line comment
			isComment = true
			multiCommentType = commentTypeSlash
		}
	}

	// parse comment content
	if isComment {
		for {
			ch = l.Next()
			switch ch {
			case syntax.RuneEOF:
				return true, syntax.Token{
					Type:     TypeComment,
					StartIdx: startIdx,
					EndIdx:   l.GetCursor(),
				}, nil
			case syntax.RuneCR, syntax.RuneLF:
				// single line
				if multiCommentType == commentTypeSingle {
					return true, syntax.Token{
						Type:     TypeComment,
						StartIdx: startIdx,
						EndIdx:   l.GetCursor(),
					}, nil
				}
				p := l.Peek()
				if (ch == syntax.RuneCR && p == syntax.RuneLF) || (ch == syntax.RuneLF && p == syntax.RuneCR) {
					l.Next()
				}
				l.Lines = append(l.Lines, syntax.LineInfo{
					Indents:  0,
					StartIdx: l.GetCursor() + 1,
				})
			case LeftDoubleQuoteI:
				if multiCommentType == commentTypeQuoteI {
					quoteCount += 1
				}
			case LeftDoubleQuoteII:
				if multiCommentType == commentTypeQuoteII {
					quoteCount += 1
				}
			case RightDoubleQuoteI:
				if multiCommentType == commentTypeQuoteI {
					quoteCount -= 1
					if quoteCount == 0 {
						l.Next()
						return true, syntax.Token{
							Type:     TypeComment,
							StartIdx: startIdx,
							EndIdx:   l.GetCursor(),
						}, nil
					}
				}
			case RightDoubleQuoteII:
				if multiCommentType == commentTypeQuoteII {
					quoteCount -= 1
					if quoteCount == 0 {
						l.Next()
						return true, syntax.Token{
							Type:     TypeComment,
							StartIdx: startIdx,
							EndIdx:   l.GetCursor(),
						}, nil
					}
				}
			case MultiplyMark:
				if multiCommentType == commentTypeSlash && l.Peek() == Slash {
					l.Next()
					l.Next()
					return true, syntax.Token{
						Type:     TypeComment,
						StartIdx: startIdx,
						EndIdx:   l.GetCursor(),
					}, nil
				}
			}
		}
	}

	return false, syntax.Token{}, nil
}

//// parseKeyword logic in keyword.go

//// utils
func isPureNumber(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isNumber(ch rune) bool {
	return (ch >= '0' && ch <= '9') || syntax.ContainsRune(ch, []rune{'.', '-', '+'})
}

// @params: ch - input char
func isIdentifierChar(ch rune) bool {
	return syntax.IdInRange(ch)
}
