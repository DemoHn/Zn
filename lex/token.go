package lex

import "github.com/DemoHn/Zn/util"

// TokenType - general token type
type TokenType int

// Token - general token type
type Token struct {
	Type    TokenType
	Literal []rune
	Range   TokenRange
}

// TokenRange locates the absolute position of a token
type TokenRange struct {
	// startLine - line num (start from 1) of first char
	StartLine int
	StartIdx  int
	// endLine - line num (start from 1) of last char
	EndLine int
	EndIdx  int
}

// newTokenRange creates new TokenRange struct
// with startLine & startIdx initialized.
func newTokenRange(l *Lexer) TokenRange {
	return TokenRange{
		StartLine: l.getCurrentLine(),
		StartIdx:  l.cursor,
	}
}

func (r *TokenRange) setRangeEnd(l *Lexer) {
	r.EndLine = l.CurrentLine
	r.EndIdx = l.cursor + 1
}

// GetStartLine -
func (r *TokenRange) GetStartLine() int {
	return r.StartLine
}

// GetEndLine -
func (r *TokenRange) GetEndLine() int {
	return r.EndLine
}

//// 0. EOF

// EOF - mark as end of file, should only exists at the end of sequence
const EOF rune = 0

//// 1. keywords
// TokenTypePrefix: 0x10
// keywords are all ideoglyphs that its length varies from its definitions.
// so here we define all possible chars that may be an element of one keyword.
const (
	// GlyphXXX - 关键词中文名 - 可能出现的关键词位置
	// GlyphLING - 令 - 令
	GlyphLING rune = 0x4EE4
	// GlyphWEI - 为 - 为，不为，成为，作为，是为，何为
	GlyphWEI rune = 0x4E3A
	// GlyphSHI - 是 - 是为
	GlyphSHI rune = 0x662F
	// GlyphRU - 如 - 如何，如果
	GlyphRU rune = 0x5982
	// GlyphHE - 何 - 如何，何为
	GlyphHE rune = 0x4F55
	// GlyphYI - 已 - 已知
	GlyphYI rune = 0x5DF2
	// GlyphZHIy - 知 - 已知
	GlyphZHIy rune = 0x77E5
	// GlyphFAN - 返 - 返回
	GlyphFAN rune = 0x8FD4
	// GlyphHUI - 回 - 返回
	GlyphHUI rune = 0x56DE
	// GlyphBU - 不 - 不为，不是，不等于，不大于，不小于
	GlyphBU rune = 0x4E0D
	// GlyphDENG - 等 - 等于，不等于
	GlyphDENG rune = 0x7B49
	// GlyphYU - 于 - 不等于，大于，不大于，小于，不小于
	GlyphYU rune = 0x4E8E
	// GlyphDA - 大 - 大于，不大于
	GlyphDA rune = 0x5927
	// GlyphXIAO - 小 - 小于，不小于
	GlyphXIAO rune = 0x5C0F
	// GlyphYIi - 以 - 以
	GlyphYIi rune = 0x4EE5
	// GlyphDE - 得 - 得到
	GlyphDE rune = 0x5F97
	// GlyphDAO - 到 - 得到
	GlyphDAO rune = 0x5230
	// GlyphGUO - 果 - 如果
	GlyphGUO rune = 0x679C
	// GlyphZE - 则 - 否则
	GlyphZE rune = 0x5219
	// GlyphFOU - 否 - 否则
	GlyphFOU rune = 0x5426
	// GlyphMEI - 每 - 每当
	GlyphMEI rune = 0x6BCF
	// GlyphDANG - 当 - 每当
	GlyphDANG rune = 0x5F53
	// GlyphCHENG - 成 - 成为
	GlyphCHENG rune = 0x6210
	// GlyphZUO - 作 - 作为
	GlyphZUO rune = 0x4F5C
	// GlyphDING - 定 - 定义
	GlyphDING rune = 0x5B9A
	// GlyphYIy - 义 - 定义
	GlyphYIy rune = 0x4E49
	// GlyphLEI - 类 - 类比
	GlyphLEI rune = 0x7C7B
	// GlyphBI - 比 - 类比
	GlyphBI rune = 0x6BD4
	// GlyphQI - 其 - 其
	GlyphQI rune = 0x5176
	// GlyphCI - 此 - 此，此之
	GlyphCI rune = 0x6B64
	// GlyphZHU - 注 - 注
	GlyphZHU rune = 0x6CE8
	// GlyphDUI - 对 - 对于
	GlyphDUI rune = 0x5BF9
	// GlyphHUO - 或 - 或
	GlyphHUO rune = 0x6216
	// GlyphQIE - 且 - 且
	GlyphQIE rune = 0x4E14
	// GlyphZHI - 之 - 之
	GlyphZHI rune = 0x4E4B
	// GlyphZAI - 再 - 再如
	GlyphZAI rune = 0x518D
)

// KeywordLeads - all glyphs that would be possible of the first character of one keyword.
var KeywordLeads = []rune{
	GlyphLING, GlyphWEI, GlyphSHI,
	GlyphRU, GlyphYI, GlyphFAN, GlyphBU, GlyphDENG,
	GlyphDA, GlyphXIAO, GlyphYIi,
	GlyphDE, GlyphFOU, GlyphMEI,
	GlyphCHENG, GlyphZUO, GlyphDING, GlyphLEI,
	GlyphQI, GlyphCI, GlyphHE, GlyphHUO, GlyphQIE,
	GlyphDUI, GlyphZHI, GlyphZAI,
}

//// 2. markers
// declare marks
const (
	Comma             rune = 0xFF0C //，
	Colon             rune = 0xFF1A //：
	Semicolon         rune = 0xFF1B //；
	QuestionMark      rune = 0xFF1F //？
	RefMark           rune = 0x0026 // &
	BangMark          rune = 0xFF01 // ！
	AnnotationMark    rune = 0x0040 // @
	HashMark          rune = 0x0023 // #
	EllipsisMark      rune = 0x2026 // …
	LeftBracket       rune = 0x3010 // 【
	RightBracket      rune = 0x3011 // 】
	LeftParen         rune = 0xFF08 // （
	RightParen        rune = 0xFF09 // ）
	Equal             rune = 0x003D // =
	DoubleArrow       rune = 0x27FA // ⟺
	LeftCurlyBracket  rune = 0x007B // {
	RightCurlyBracket rune = 0x007D // }
)

// MarkLeads -
var MarkLeads = []rune{
	Comma, Colon, Semicolon, QuestionMark, RefMark, BangMark,
	AnnotationMark, HashMark, EllipsisMark, LeftBracket,
	RightBracket, LeftParen, RightParen, Equal, DoubleArrow,
	LeftCurlyBracket, RightCurlyBracket,
}

//// 3. spaces
const (
	SP  rune = 0x0020 // <SP>
	TAB rune = 0x0009 // <TAB>
	CR  rune = 0x000D // \r
	LF  rune = 0x000A // \n
)

// WhiteSpaces - all kinds of valid spaces
var WhiteSpaces = []rune{
	// where 0x0020 <--> SP
	0x0009, 0x000B, 0x000C, 0x0020, 0x00A0,
	0x2000, 0x2001, 0x2002, 0x2003, 0x2004,
	0x2005, 0x2006, 0x2007, 0x2008, 0x2009,
	0x200A, 0x200B, 0x202F, 0x205F, 0x3000,
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

//// 4. quotes
// declare quotes
const (
	LeftQuoteI    rune = 0x300A //《
	RightQuoteI   rune = 0x300B // 》
	LeftQuoteII   rune = 0x300C // 「
	RightQuoteII  rune = 0x300D // 」
	LeftQuoteIII  rune = 0x300E // 『
	RightQuoteIII rune = 0x300F // 』
	LeftQuoteIV   rune = 0x201C // “
	RightQuoteIV  rune = 0x201D // ”
	LeftQuoteV    rune = 0x2018 // ‘
	RightQuoteV   rune = 0x2019 // ’
)

// LeftQuotes -
var LeftQuotes = []rune{
	LeftQuoteI,
	LeftQuoteII,
	LeftQuoteIII,
	LeftQuoteIV,
	LeftQuoteV,
}

// RightQuotes -
var RightQuotes = []rune{
	RightQuoteI,
	RightQuoteII,
	RightQuoteIII,
	RightQuoteIV,
	RightQuoteV,
}

// QuoteMatchMap -
var QuoteMatchMap = map[rune]rune{
	LeftQuoteI:   RightQuoteI,
	LeftQuoteII:  RightQuoteII,
	LeftQuoteIII: RightQuoteIII,
	LeftQuoteIV:  RightQuoteIV,
	LeftQuoteV:   RightQuoteV,
}

//// 5. var quote
const (
	MiddleDot rune = 0x00B7 // ·
)

//// 6. numbers
func isNumber(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

//// 7. identifiers
const maxIdentifierLength = 32

// @params: ch - input char
// @params: isFirst - is the first char of identifier
func isIdentifierChar(ch rune, isFirst bool) bool {
	// CJK unified ideograph
	if ch >= 0x4E00 && ch <= 0x9FFF {
		return true
	}
	// 〇, _
	if ch == 0x3007 || ch == '_' {
		return true
	}
	// A-Z
	if ch >= 'A' && ch <= 'Z' {
		return true
	}
	if ch >= 'a' && ch <= 'z' {
		return true
	}
	if !isFirst {
		if ch >= '0' && ch <= '9' {
			return true
		}
		if util.Contains(ch, []rune{'*', '+', '-', '/'}) {
			return true
		}
	}
	return false
}

//// token consts and constructors
// token types -
// for special type Tokens, its range varies from 0 - 9
const (
	TypeEOF          TokenType = 0
	TypeComment      TokenType = 10
	TypeCommaSep     TokenType = 11
	TypeStmtSep      TokenType = 12
	TypeFuncCall     TokenType = 13
	TypeFuncDeclare  TokenType = 14
	TypeObjRef       TokenType = 15
	TypeMustT        TokenType = 16
	TypeAnnoT        TokenType = 17
	TypeMapHash      TokenType = 18
	TypeMoreParam    TokenType = 19
	TypeArrayQuoteL  TokenType = 20
	TypeArrayQuoteR  TokenType = 21
	TypeFuncQuoteL   TokenType = 22 // （
	TypeFuncQuoteR   TokenType = 23 // ）
	TypeMapData      TokenType = 24 // ==
	TypeStmtQuoteL   TokenType = 25 // {
	TypeStmtQuoteR   TokenType = 26 // }
	TypeMapQHash     TokenType = 27 // #{
	TypeDeclareW     TokenType = 40 // 令
	TypeLogicYesW    TokenType = 41 // 为
	TypeCondOtherW   TokenType = 43 // 再如
	TypeCondW        TokenType = 44 // 如果
	TypeFuncW        TokenType = 45 // 如何
	TypeStaticFuncW  TokenType = 46 // 何为
	TypeParamAssignW TokenType = 47 // 已知
	TypeReturnW      TokenType = 48 // 返回
	TypeLogicNotW    TokenType = 49 // 不为
	TypeLogicNotEqW  TokenType = 51 // 不等于
	TypeLogicLteW    TokenType = 52 // 不大于
	TypeLogicGteW    TokenType = 53 // 不小于
	TypeLogicLtW     TokenType = 54 // 小于
	TypeLogicGtW     TokenType = 55 // 大于
	TypeVarOneW      TokenType = 56 // 以
	//
	TypeFuncYieldW    TokenType = 58 // 得到
	TypeCondElseW     TokenType = 59 // 否则
	TypeWhileLoopW    TokenType = 60 // 每当
	TypeObjNewW       TokenType = 61 // 成为
	TypeVarAliasW     TokenType = 62 // 作为
	TypeObjDefineW    TokenType = 63 // 定义
	TypeObjTraitW     TokenType = 64 // 类比
	TypeObjThisW      TokenType = 65 // 其
	TypeObjSelfW      TokenType = 66 // 此
	TypeFuncCallOneW  TokenType = 67 // 对于
	TypeLogicOrW      TokenType = 69 // 或
	TypeLogicAndW     TokenType = 70 // 且
	TypeObjDotW       TokenType = 71 // 之
	TypeObjConstructW TokenType = 73 // 是为
	TypeLogicEqualW   TokenType = 74 // 等于
	TypeStaticSelfW   TokenType = 75 // 此之
	TypeString        TokenType = 90
	TypeVarQuote      TokenType = 91
	TypeNumber        TokenType = 100
	TypeIdentifier    TokenType = 101
)

// KeywordTypeMap -
var KeywordTypeMap = map[TokenType][]rune{
	TypeDeclareW:      []rune{GlyphLING},
	TypeLogicYesW:     []rune{GlyphWEI},
	TypeCondW:         []rune{GlyphRU, GlyphGUO},
	TypeFuncW:         []rune{GlyphRU, GlyphHE},
	TypeStaticFuncW:   []rune{GlyphHE, GlyphWEI},
	TypeParamAssignW:  []rune{GlyphYI, GlyphZHIy},
	TypeReturnW:       []rune{GlyphFAN, GlyphHUI},
	TypeLogicNotW:     []rune{GlyphBU, GlyphWEI},
	TypeLogicNotEqW:   []rune{GlyphBU, GlyphDENG, GlyphYU},
	TypeLogicLteW:     []rune{GlyphBU, GlyphDA, GlyphYU},
	TypeLogicGteW:     []rune{GlyphBU, GlyphXIAO, GlyphYU},
	TypeLogicLtW:      []rune{GlyphXIAO, GlyphYU},
	TypeLogicGtW:      []rune{GlyphDA, GlyphYU},
	TypeVarOneW:       []rune{GlyphYIi},
	TypeFuncYieldW:    []rune{GlyphDE, GlyphDAO},
	TypeCondElseW:     []rune{GlyphFOU, GlyphZE},
	TypeWhileLoopW:    []rune{GlyphMEI, GlyphDANG},
	TypeObjNewW:       []rune{GlyphCHENG, GlyphWEI},
	TypeVarAliasW:     []rune{GlyphZUO, GlyphWEI},
	TypeObjDefineW:    []rune{GlyphDING, GlyphYIy},
	TypeObjTraitW:     []rune{GlyphLEI, GlyphBI},
	TypeObjThisW:      []rune{GlyphQI},
	TypeObjSelfW:      []rune{GlyphCI},
	TypeFuncCallOneW:  []rune{GlyphDUI, GlyphYU},
	TypeLogicOrW:      []rune{GlyphHUO},
	TypeLogicAndW:     []rune{GlyphQIE},
	TypeObjDotW:       []rune{GlyphZHI},
	TypeObjConstructW: []rune{GlyphSHI, GlyphWEI},
	TypeLogicEqualW:   []rune{GlyphDENG, GlyphYU},
	TypeStaticSelfW:   []rune{GlyphCI, GlyphZHI},
	TypeCondOtherW:    []rune{GlyphZAI, GlyphRU},
}

// NewTokenEOF - new EOF token
func NewTokenEOF(line int, col int) *Token {
	return &Token{
		Type:    TypeEOF,
		Literal: []rune{},
		Range: TokenRange{
			StartLine: line,
			StartIdx:  col,
			EndLine:   line,
			EndIdx:    col,
		},
	}
}

// NewStringToken -
func NewStringToken(buf []rune, quoteType rune, rg TokenRange) *Token {
	return &Token{
		Type:    TypeString,
		Literal: util.Copy(buf),
		Range:   rg,
	}
}

// NewVarQuoteToken -
func NewVarQuoteToken(buf []rune, rg TokenRange) *Token {
	return &Token{
		Type:    TypeVarQuote,
		Literal: util.Copy(buf),
		Range:   rg,
	}
}

// NewCommentToken -
func NewCommentToken(buf []rune, isMultiLine bool, rg TokenRange) *Token {
	return &Token{
		Type:    TypeComment,
		Literal: util.Copy(buf),
		Range:   rg,
	}
}

// NewNumberToken -
func NewNumberToken(buf []rune, rg TokenRange) *Token {
	return &Token{
		Type:    TypeNumber,
		Literal: util.Copy(buf),
		Range:   rg,
	}
}

// NewMarkToken -
func NewMarkToken(buf []rune, t TokenType, startR TokenRange, num int) *Token {
	rg := startR
	rg.EndLine = startR.StartLine
	rg.EndIdx = startR.StartIdx + num
	return &Token{
		Type:    t,
		Literal: util.Copy(buf),
		Range:   rg,
	}
}

// NewKeywordToken -
func NewKeywordToken(t TokenType) *Token {
	var l = []rune{}
	if item, ok := KeywordTypeMap[t]; ok {
		l = item
	}
	return &Token{
		Type:    t,
		Literal: l,
	}
}

// NewIdentifierToken -
func NewIdentifierToken(buf []rune, rg TokenRange) *Token {
	return &Token{
		Type:    TypeIdentifier,
		Literal: util.Copy(buf),
		Range:   rg,
	}
}
