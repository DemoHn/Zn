package zh

import (
	"strconv"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

// // 1. punctuations
const (
	Comma             rune = 0xFF0C // ，
	Comma_EN          rune = 0x002C // en ,
	PauseComma        rune = 0x3001 // 、
	Colon             rune = 0xFF1A // ：
	Colon_EN          rune = 0x003A // en :
	Semicolon         rune = 0xFF1B // ；
	Semicolon_EN      rune = 0x003B // en ;
	QuestionMark      rune = 0xFF1F // ？
	QuestionMark_EN   rune = 0x003F // en ?
	BangMark          rune = 0xFF01 // ！
	BangMark_EN       rune = 0x0021 // en !
	LeftBracket       rune = 0x3010 // 【
	LeftBracket_EN    rune = 0x005B // en [
	RightBracket      rune = 0x3011 // 】
	RightBracket_EN   rune = 0x005D // en ]
	LeftParen         rune = 0xFF08 //（
	LeftParen_EN      rune = 0x0028 //  en (
	RightParen        rune = 0xFF09 // ）
	RightParen_EN     rune = 0x0029 // en )
	LeftCurlyBracket  rune = 0x007B // {
	RightCurlyBracket rune = 0x007D // }
)

// // 2. quotes
// declare quotes

const (
	LeftLibQuoteI      rune = 0x300A //《
	RightLibQuoteI     rune = 0x300B // 》
	LeftDoubleQuoteI   rune = 0x300C // 「
	RightDoubleQuoteI  rune = 0x300D // 」
	LeftDoubleQuoteII  rune = 0x201C // “
	RightDoubleQuoteII rune = 0x201D // ”
	LeftSingleQuoteI   rune = 0x300E // 『
	RightSingleQuoteI  rune = 0x300F // 』
	LeftSingleQuoteII  rune = 0x2018 // ‘
	RightSingleQuoteII rune = 0x2019 // ’
)

// // 3. operators
const (
	RefOp         rune = 0x0026 // &
	AnnotationOp  rune = 0x0040 // @
	HashOp        rune = 0x0023 // #
	EqualOp       rune = 0x003D // =
	LessThanOp    rune = 0x003C // <
	GreaterThanOp rune = 0x003E // >
	PlusOp        rune = 0x002B // +
	MinusOp       rune = 0x002D // -
	MultiplyOp    rune = 0x002A // *
	SlashOp       rune = 0x002F // /
	IntDivOp      rune = 0x007C // |
	RemainderOp   rune = 0x0025 // %
)

// // 4. var quote
const (
	BackTick rune = 0x0060 // `
)

// // 5. comment keyword
const (
	CharZHU rune = 0x6CE8 // 注
)

// // token constants and constructors (without keyword token)
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
	TypeStmtQuoteL    uint8 = 24 // {
	TypeStmtQuoteR    uint8 = 25 // }
	TypePauseCommaSep uint8 = 26 // 、
	TypeIntDivMark    uint8 = 27 // |
	TypeModuloMark    uint8 = 28 // %
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

// // Comment Types -
const (
	commentTypeSingle  = 1 // single line
	commentTypeSlash   = 2 // multiple line, starts with '/*'
	commentTypeQuoteI  = 3 // multiple line, starts with '注：「'
	commentTypeQuoteII = 4 // multiple line, starts with '注：“'
)

var markPunctuations = []rune{
	Comma,
	Comma_EN,
	PauseComma,
	Colon,
	Colon_EN,
	Semicolon,
	Semicolon_EN,
	QuestionMark,
	QuestionMark_EN,
	BangMark,
	BangMark_EN,
	LeftBracket,
	LeftBracket_EN,
	RightBracket,
	RightBracket_EN,
	LeftParen,
	LeftParen_EN,
	RightParen,
	RightParen_EN,
	LeftCurlyBracket,
	RightCurlyBracket,
}

var markOperators = []rune{
	RefOp,
	AnnotationOp,
	HashOp,
	EqualOp,
	LessThanOp,
	GreaterThanOp,
	PlusOp,
	MinusOp,
	MultiplyOp,
	SlashOp,
	IntDivOp,
	RemainderOp,
}

var markQuotes = []rune{
	LeftLibQuoteI,
	RightLibQuoteI,
	LeftDoubleQuoteI,
	RightDoubleQuoteI,
	LeftDoubleQuoteII,
	RightDoubleQuoteII,
	LeftSingleQuoteI,
	RightSingleQuoteI,
	LeftSingleQuoteII,
	RightSingleQuoteII,
}

// quote match map
var quoteMatchMap = map[rune]rune{
	LeftDoubleQuoteI:  RightDoubleQuoteI,
	LeftDoubleQuoteII: RightDoubleQuoteII,
	LeftSingleQuoteI:  RightSingleQuoteI,
	LeftSingleQuoteII: RightSingleQuoteII,
	LeftLibQuoteI:     RightLibQuoteI,
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
		return parseEOF(l)
	case CharZHU, SlashOp:
		// save current cursor location (as) savepoint - when parsing 注-like
		// token as comment failed, it's time to turn back (to the savepoint) and try to treat it
		// as an identifier
		savePointLoc := l.GetCursor()
		// try to parse 注 or / as comment, if not, try to parse as other types (e.g. identifier)
		isComment, tk, err := parseComment(l)
		if err != nil {
			return syntax.Token{}, err
		}
		if isComment {
			return tk, nil
		} else {
			l.SetCursor(savePointLoc)
			// then fallthrough to the next logic - parse as an operator or identifier
			// DO NOT WRITE return-statement HERE!!!
		}
	case LeftLibQuoteI, LeftDoubleQuoteI, LeftDoubleQuoteII, LeftSingleQuoteI, LeftSingleQuoteII:
		return parseString(l)
	case BackTick:
		return parseVarQuote(l)
	}

	// other token types
	if syntax.ContainsRune(ch, markPunctuations) {
		return parsePunctuations(l)
	}
	// NOTICE: if err == nil && isOperator == false, The logic will PASS THROUGH
	// instead of return TOKEN or return error!
	if syntax.ContainsRune(ch, markOperators) {
		isOperator, tk, err := parseOperators(l)
		if err != nil {
			return syntax.Token{}, err
		}
		if isOperator {
			return tk, nil
		}
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

func parsePunctuations(l *syntax.Lexer) (syntax.Token, error) {
	startIdx := l.GetCursor()
	ch := l.GetCurrentChar()

	punctuationTypeMap := map[rune]uint8{
		Comma:             TypeCommaSep,
		Comma_EN:          TypeCommaSep,
		PauseComma:        TypePauseCommaSep,
		Colon:             TypeFuncCall,
		Colon_EN:          TypeFuncCall,
		Semicolon:         TypeStmtSep,
		Semicolon_EN:      TypeStmtSep,
		QuestionMark:      TypeFuncDeclare,
		QuestionMark_EN:   TypeFuncDeclare,
		BangMark:          TypeExceptionT,
		BangMark_EN:       TypeExceptionT,
		LeftBracket:       TypeArrayQuoteL,
		LeftBracket_EN:    TypeArrayQuoteL,
		RightBracket:      TypeArrayQuoteR,
		RightBracket_EN:   TypeArrayQuoteR,
		LeftParen:         TypeFuncQuoteL,
		LeftParen_EN:      TypeFuncQuoteL,
		RightParen:        TypeFuncQuoteR,
		RightParen_EN:     TypeFuncQuoteR,
		LeftCurlyBracket:  TypeStmtQuoteL,
		RightCurlyBracket: TypeStmtQuoteR,
	}

	if pType, ok := punctuationTypeMap[ch]; ok {
		l.Next()
		return syntax.Token{
			Type:     pType,
			StartIdx: startIdx,
			EndIdx:   l.GetCursor(),
		}, nil
	}

	return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
}

func parseOperators(l *syntax.Lexer) (bool, syntax.Token, error) {
	startIdx := l.GetCursor()
	ch := l.GetCurrentChar()
	var tokenType uint8

	switch ch {
	case RefOp:
		tokenType = TypeObjRef
	case AnnotationOp:
		tokenType = TypeAnnotationT
	case HashOp:
		tokenType = TypeMapHash
	case EqualOp:
		if l.Peek() == EqualOp { // op: ==
			l.Next()
			tokenType = TypeEqualMark
		} else { // op: =
			tokenType = TypeAssignMark
		}
	case LessThanOp:
		if l.Peek() == EqualOp { // op: <=
			l.Next()
			tokenType = TypeLTEMark
		} else { // op: <
			tokenType = TypeLTMark
		}
	case GreaterThanOp:
		if l.Peek() == EqualOp {
			l.Next()
			tokenType = TypeGTEMark
		} else {
			tokenType = TypeGTMark
		}
	case IntDivOp:
		tokenType = TypeIntDivMark
	case RemainderOp:
		tokenType = TypeModuloMark
	case PlusOp, MinusOp, MultiplyOp, SlashOp: // op: + - * /
		chn := l.Peek()

		t := TypePlus
		switch ch {
		case PlusOp:
			t = TypePlus
		case MinusOp:
			t = TypeMinus
		case MultiplyOp:
			t = TypeMultiply
		case SlashOp:
			// parse /=, example usage: '如果 X /= 10'
			if chn == EqualOp {
				l.Next()
				l.Next()
				return true, syntax.Token{
					Type:     TypeNEMark,
					StartIdx: startIdx,
					EndIdx:   l.GetCursor(),
				}, nil
			} else {
				t = TypeDivision
			}
		}
		// NOTE: the next char must be space/punctuations/quotes to ensure it's not a part of
		// identifier
		if syntax.IsWhiteSpace(chn) || syntax.ContainsRune(chn, markPunctuations) || syntax.ContainsRune(chn, markQuotes) {
			l.Next()
			return true, syntax.Token{Type: t, StartIdx: startIdx, EndIdx: l.GetCursor()}, nil
		} else {
			return false, syntax.Token{}, nil
		}

	default:
		return false, syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
	}

	l.Next()
	return true, syntax.Token{
		Type:     tokenType,
		StartIdx: startIdx,
		EndIdx:   l.GetCursor(),
	}, nil
}

// parse ` <identifier> `
func parseVarQuote(l *syntax.Lexer) (syntax.Token, error) {
	startIdx := l.GetCursor()
	var literal []rune
	for {
		ch := l.Next()

		if isIdentifierChar(ch) || syntax.ContainsRune(ch, syntax.IDContinue) {
			literal = append(literal, ch)
		} else if ch == BackTick {
			l.Next()
			return syntax.Token{
				Type:     TypeIdentifier,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
				Literal:  literal,
			}, nil
		} else {
			return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
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
			return syntax.Token{}, zerr.IncompleteString(l.GetCursor())
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
		case BackTick:
			literal = unescapeBackTickSpecialStr(l, literal)
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

	terminateMarkers := append([]rune{
		// EOF or newline chars
		syntax.RuneEOF, syntax.RuneCR, syntax.RuneLF,
		// operators except for PlusOp, MinusOp, MultiplyOp, SlashOp
		RefOp, AnnotationOp, HashOp, EqualOp, LessThanOp, GreaterThanOp, IntDivOp,
	}, markPunctuations...)

	// 0. first char must be an identifier
	if !isIdentifierChar(ch) {
		return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
	}
	for {
		ch = l.Next()
		// 1. when next char is space, stop here
		if syntax.IsWhiteSpace(ch) {
			goto ID_end
		}

		// 2. when next char is a part of keyword, stop here
		isKeyword, _, err := parseKeyword(l, false)
		if err != nil {
			return syntax.Token{}, err
		}
		if isKeyword {
			goto ID_end
		}
		// 3. when next char is a start of comment, stop here
		// only 「//」 and 「/*」 type is available
		// NOTE: we will regard comment type「注」 as a regular identifier
		if ch == SlashOp && syntax.ContainsRune(l.Peek(), []rune{SlashOp, MultiplyOp, EqualOp}) {
			goto ID_end
		}
		// 4. when next char is a mark, stop here
		if syntax.ContainsRune(ch, terminateMarkers) {
			goto ID_end
		}
		// 5. otherwise, if it's an identifier with *, /, .
		// add char to literal
		if isIdentifierChar(ch) || syntax.ContainsRune(ch, syntax.IDContinue) {
			literal = append(literal, ch)
			continue
		}
		return syntax.Token{}, zerr.InvalidChar(ch, l.GetCursor())
	}
ID_end:
	// SlashOp ('/') COULD NOT be the last char of an identifer token
	// to avoid confusion
	if literal[len(literal)-1] == SlashOp {
		return syntax.Token{}, zerr.InvalidChar(SlashOp, l.GetCursor()-1)
	}
	// else return a complete identifier token
	return syntax.Token{
		Type:     TypeIdentifier,
		StartIdx: startIdx,
		EndIdx:   l.GetCursor(),
		Literal:  literal,
	}, nil
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
			return false, syntax.Token{}, nil
		}
	case SlashOp:
		p := l.Peek()
		if p == SlashOp {
			l.Next()
			// single line comment
			isComment = true
			multiCommentType = commentTypeSingle
		} else if p == MultiplyOp {
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
			case MultiplyOp:
				if multiCommentType == commentTypeSlash && l.Peek() == SlashOp {
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

func parseEOF(l *syntax.Lexer) (syntax.Token, error) {
	// set last lineText
	if len(l.Lines) > 0 {
		lastLine := &(l.Lines[len(l.Lines)-1])
		startIdx := lastLine.StartIdx

		// skip INDENTS when inserting source text
		if l.IndentType == syntax.IndentSpace {
			startIdx += 4 * lastLine.Indents
		} else if l.IndentType == syntax.IndentTab {
			startIdx += lastLine.Indents
		}

		lastLine.LineText = l.Source[startIdx:l.GetCursor()]
	}
	return syntax.Token{Type: TypeEOF, StartIdx: l.GetCursor(), EndIdx: l.GetCursor()}, nil
}

//// parseKeyword logic in keyword.go

// unescape backtick-marked special string (e.g. `CR`, `U+1F005`) to its original character
func unescapeBackTickSpecialStr(l *syntax.Lexer, srcLiteral []rune) []rune {
	// store all characters after the backtick (including the 'backtick' marker itself)
	literalBuffer := []rune{l.GetCurrentChar()}

	// hand-written CR state matchine
	const (
		sBegin  = 1
		sC      = 2
		sR      = 3
		sL      = 4
		sF      = 5
		sT      = 6
		sA      = 7
		sB      = 8
		sK      = 9
		sS      = 10
		sP      = 11
		sU      = 12
		smP     = 13 // U+xxxx
		sHexNum = 14
	)
	var state = sBegin
	var hexCount = 0

	for {
		// peek the next char - if the next char is a quote mark, there are 2 cases:
		//   a) the next next char is a backtick (`“`) - the mark is "considered" as a standalone quote (no need to pair)
		//      and we add the char to srcLiteral directly
		//   b) the next next char is other string (`”balhbalh)- NO WAY, stop before parsing the quote mark
		switch l.Peek() {
		case LeftDoubleQuoteI, LeftDoubleQuoteII, LeftSingleQuoteI, LeftSingleQuoteII, LeftLibQuoteI,
			RightDoubleQuoteI, RightDoubleQuoteII, RightSingleQuoteI, RightSingleQuoteII, RightLibQuoteI:
			qch := l.Peek()
			if l.GetCurrentChar() == BackTick && l.Peek2() == BackTick {
				l.Next()
				l.Next()
				return append(srcLiteral, qch)
			} else {
				goto UNDONE_end
			}
		}

		cch := l.Next()
		literalBuffer = append(literalBuffer, cch)
		// to match U+xxxx, the char range is [0-9A-Fa-f]
		// because A,B,C,F could be either part of unicode Number or inside the word "TAB", "CR", "LF"
		if (cch >= '0' && cch <= '9') || (cch >= 'A' && cch <= 'F') {
			switch state {
			case smP:
				state = sHexNum
				hexCount = 1
				continue
			case sHexNum:
				hexCount += 1
				continue
			}
			// for other cases, fallthrough to CR/LF/TAB etc... keyword handling
		}

		switch cch {
		case 'C': // CR / CRLF
			switch state {
			case sBegin:
				state = sC
			default:
				goto UNDONE_end
			}
		case 'L': // LF / CRLF
			switch state {
			case sBegin, sR:
				state = sL
			default:
				goto UNDONE_end
			}
		case 'T': // TAB
			switch state {
			case sBegin:
				state = sT
			default:
				goto UNDONE_end
			}
		case 'S': // SP
			switch state {
			case sBegin:
				state = sS
			default:
				goto UNDONE_end
			}
		case 'B': // BK / TAB
			switch state {
			case sBegin, sA:
				state = sB
			default:
				goto UNDONE_end
			}
		case 'U': // U+xxxx
			switch state {
			case sBegin:
				state = sU
			default:
				goto UNDONE_end
			}
		case 'R':
			switch state {
			case sC:
				state = sR
			default:
				goto UNDONE_end
			}
		case 'F':
			switch state {
			case sL:
				state = sF
			default:
				goto UNDONE_end
			}
		case 'A': // TAB
			switch state {
			case sT:
				state = sA
			default:
				goto UNDONE_end
			}
		case 'P':
			switch state {
			case sS:
				state = sP
			default:
				goto UNDONE_end
			}
		case 'K':
			switch state {
			case sB:
				state = sK
			default:
				goto UNDONE_end
			}
		case '+':
			switch state {
			case sU:
				state = smP
			default:
				goto UNDONE_end
			}
		case '`':
			// for normal cases, unescape the character and append
			switch string(literalBuffer) {
			case "`TAB`":
				return append(srcLiteral, '\t')
			case "`BK`":
				return append(srcLiteral, '`')
			case "`SP`":
				return append(srcLiteral, ' ')
			case "`CR`":
				return append(srcLiteral, '\r')
			case "`LF`":
				return append(srcLiteral, '\n')
			case "`CRLF`":
				return append(srcLiteral, []rune{'\r', '\n'}...)
			default: // U+xxxx
				if state == sHexNum {
					if hexCount >= 1 && hexCount <= 8 {
						hexStr := string(literalBuffer[3 : len(literalBuffer)-1])
						hexNum, _ := strconv.ParseInt(hexStr, 16, 32)
						return append(srcLiteral, rune(hexNum))
					}
				}
				goto UNDONE_end
			}

		default:
			goto UNDONE_end
		}
	}
UNDONE_end:
	// unescape the string fails, KEEP the original string to the final literal
	return append(srcLiteral, literalBuffer...)
}

// // utils
func isPureNumber(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// @params: ch - input char
func isIdentifierChar(ch rune) bool {
	return syntax.IdInRange(ch)
}
