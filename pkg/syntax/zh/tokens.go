package zh

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/util"
)

// TokenBuilderZH
type TokenBuilderZH struct{}

//// 2. markers
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
)

//// token constants and constructors (without keyword token)
// token types -
// for special type Tokens, its range varies from 0 - 9
// for keyword types, check lex/keyword.go for details
const (
	TypeEOF           uint8 = 0
	TypeSpace         uint8 = 1  //
	TypeString        uint8 = 2  // string (only double quotes)
	TypeNumber        uint8 = 4  // numbers
	TypeIdentifier    uint8 = 5  //
	TypeEnumString    uint8 = 6  // string (with single quotes)
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
	TypeAssignMark    uint8 = 29 // uint8 =
	TypeGTMark        uint8 = 30 // >
	TypeLTMark        uint8 = 31 // <
	TypeGTEMark       uint8 = 32 // >uint8 =
	TypeLTEMark       uint8 = 33 // <uint8 =
	TypeNEMark        uint8 = 34 // /uint8 =
	TypeEqualMark     uint8 = 35 // uint8 =uint8 =
	TypePlus          uint8 = 36 // +
	TypeMinus         uint8 = 37 // -
	TypeMultiply      uint8 = 38 // *
	TypeDivision      uint8 = 39 // /
)

// MarkLeads -
var MarkLeads = []rune{
	Comma, PauseComma, Colon, Semicolon, QuestionMark, RefMark, BangMark,
	AnnotationMark, HashMark, LeftBracket,
	RightBracket, LeftParen, RightParen, Equal,
	LeftCurlyBracket, RightCurlyBracket, LessThanMark, GreaterThanMark,
}

// WhiteSpaces - all kinds of valid spaces
var WhiteSpaces = []rune{
	// where 0x0020 <--> SP
	0x0009, 0x000B, 0x000C, 0x0020, 0x00A0,
	0x2000, 0x2001, 0x2002, 0x2003, 0x2004,
	0x2005, 0x2006, 0x2007, 0x2008, 0x2009,
	0x200A, 0x200B, 0x202F, 0x205F, 0x3000,
}

// NextToken -
func (tb *TokenBuilderZH) NextToken(l *syntax.Lexer) (syntax.Token, error) {
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
		if isWhiteSpace(chn) {
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
		if isWhiteSpace(chn) {
			l.Next()
			return syntax.Token{
				Type:     TypeDivision,
				StartIdx: startIdx,
				EndIdx:   l.GetCursor(),
			}, nil
		}
	}

	// other token types
	if isNumberLead(ch) {
		return parseNumber(l)
	}
	if util.Contains(ch, MarkLeads) {
		return parseMarkers(l)
	}

	return syntax.Token{}, nil
}

// regex: ^[-+]?[0-9]*\.?[0-9]+((([eE][-+])|(\*(10)?\^[-+]?))[0-9]+)?$
// ref: https://github.com/DemoHn/Zn/issues/4
func parseNumber(l *syntax.Lexer) (syntax.Token, error) {
	ch := l.GetCurrentChar()
	startIdx := l.GetCursor()
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
		case '_':
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
		ch = l.Next()
	}

end:
	if util.ContainsInt(state, endStates) {
		return syntax.Token{
			Type:     TypeNumber,
			StartIdx: startIdx,
			EndIdx:   l.GetCursor(),
		}, nil
	}
	return syntax.Token{}, zerr.InvalidChar(ch)
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
		return syntax.Token{}, zerr.InvalidChar(ch)
	}
	// include all necessary chars
	l.Next()

	return syntax.Token{
		Type:     tokenType,
		StartIdx: startIdx,
		EndIdx:   l.GetCursor(),
	}, nil
}

//// utils
func isNumberLead(ch rune) bool {
	return (ch >= '0' && ch <= '9') || util.Contains(ch, []rune{'.', '-', '+'})
}

// helpers
func isWhiteSpace(ch rune) bool {
	for _, whiteSpace := range WhiteSpaces {
		if ch == whiteSpace {
			return true
		}
	}

	return false
}
