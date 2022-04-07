package zh

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/util"
)

// TokenBuilderZH
type TokenBuilderZH struct {
	noBeginLex bool
}

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
	Slash             rune = 0x002F // /
	LessThanMark      rune = 0x003C // <
	GreaterThanMark   rune = 0x003E // >
)

//// token constants and constructors (without keyword token)
// token types -
// for special type Tokens, its range varies from 0 - 9
// for keyword types, check lex/keyword.go for details
const (
	TypeEOF           = 0
	TypeSpace         = 1  //
	TypeString        = 2  // string (only double quotes)
	TypeNumber        = 4  // numbers
	TypeIdentifier    = 5  //
	TypeEnumString    = 6  // string (with single quotes)
	TypeComment       = 10 // 注：
	TypeCommaSep      = 11 // ，
	TypeStmtSep       = 12 // ；
	TypeFuncCall      = 13 // ：
	TypeFuncDeclare   = 14 // ？
	TypeObjRef        = 15 // &
	TypeExceptionT    = 16 // ！
	TypeAnnotationT   = 17 // @
	TypeMapHash       = 18 // #
	TypeArrayQuoteL   = 20 // 【
	TypeArrayQuoteR   = 21 // 】
	TypeFuncQuoteL    = 22 // （
	TypeFuncQuoteR    = 23 // ）
	TypeStmtQuoteL    = 25 // {
	TypeStmtQuoteR    = 26 // }
	TypePauseCommaSep = 28 // 、
	TypeAssignMark    = 29 // =
	TypeGTMark        = 30 // >
	TypeLTMark        = 31 // <
	TypeGTEMark       = 32 // >=
	TypeLTEMark       = 33 // <=
	TypeNEMark        = 34 // /=
	TypeEqualMark     = 35 // ==
)

// NextToken -
func (tb *TokenBuilderZH) NextToken(l *syntax.Lexer) (syntax.Token, error) {
	// ParseBeginLex
	if !tb.noBeginLex {
		tb.noBeginLex = true
		if err := l.ParseBeginLex(); err != nil {
			return syntax.Token{}, err
		}
	}

	ch := l.GetCurrentChar()
	switch ch {
	default:
		if isNumber(ch) {
			return parseNumber(l)
		}
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
func isNumber(ch rune) bool {
	return (ch >= '0' && ch <= '9') || util.Contains(ch, []rune{'.', '-', '+'})
}
