package lex

import "github.com/DemoHn/Zn/util"

// TokenType - general token type
type TokenType int

// Token - general token type
type Token struct {
	Type    TokenType
	Literal []rune
	Info    interface{}
}

// token types -
// for special type Tokens, its range varies from 0 - 9
const (
	TypeEOF TokenType = 0
)

//// 0. EOF

// EOF - mark as end of file, should only exists at the end of sequence
const EOF rune = 0

// TokenEOF - new EOF token
func TokenEOF() *Token {
	return &Token{
		Type:    TypeEOF,
		Literal: []rune{},
		Info:    nil,
	}
}

//// 1. keywords
// TokenTypePrefix: 0x10
// keywords are all ideoglyphs that its length varies from its definitions.
// so here we define all possible chars that may be an element of one keyword.
const (
	// GlyphXXX - 关键词中文名 - 可能出现的关键词位置
	// GlyphLING - 令 - 令
	GlyphLING rune = 0x4EE4
	// GlyphWEI - 为 - 为，设为，不为，成为，作为，是为，何为
	GlyphWEI rune = 0x4E3A
	// GlyphSHI - 是 - 是，不是，是为
	GlyphSHI rune = 0x662F
	// GlyphSHE - 设 - 设为
	GlyphSHE rune = 0x8BBE
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
	// GlyphDENG - 等 - 不等于
	GlyphDENG rune = 0x7B49
	// GlyphYU - 于 - 不等于，大于，不大于，小于，不小于
	GlyphYU rune = 0x4E8E
	// GlyphDA - 大 - 大于，不大于
	GlyphDA rune = 0x5927
	// GlyphXIAO - 小 - 小于，不小于
	GlyphXIAO rune = 0x5C0F
	// GlyphYIi - 以 - 以
	GlyphYIi rune = 0x4EE5
	// GlyphER - 而 - 而
	GlyphER rune = 0x800C
	// GlyphDE - 得 - 得
	GlyphDE rune = 0x5F97
	// GlyphGUO - 果 - 如果
	GlyphGUO rune = 0x679C
	// GlyphZE - 则 - 则，否则
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
	// GlyphCI - 此 - 此
	GlyphCI rune = 0x6B64
	// GlyphZHU - 注 - 注
	GlyphZHU rune = 0x6CE8
	// GlyphZAI - 在 - 在
	GlyphZAI rune = 0x5728
	// GlyphZHONG - 中 - 中
	GlyphZHONG rune = 0x4E2D
	// GlyphHUO - 或 - 或
	GlyphHUO rune = 0x6216
	// GlyphQIE - 且 - 且
	GlyphQIE rune = 0x4E14
	// GlyphZHI - 之 - 之
	GlyphZHI rune = 0x4E4B
	// GlyphDEo - 的 - 的
	GlyphDEo rune = 0x7684
)

// KeywordLeads - all glyphs that would be possible of the first character of one keyword.
var KeywordLeads = []rune{
	GlyphLING, GlyphWEI, GlyphSHI, GlyphSHE,
	GlyphRU, GlyphYI, GlyphFAN, GlyphBU,
	GlyphDA, GlyphXIAO, GlyphYIi, GlyphER,
	GlyphDE, GlyphZE, GlyphFOU, GlyphMEI,
	GlyphCHENG, GlyphZUO, GlyphDING, GlyphLEI,
	GlyphQI, GlyphCI, GlyphZHU, GlyphHE,
	GlyphZAI, GlyphZHONG, GlyphHUO, GlyphQIE,
	GlyphZHI, GlyphDEo,
}

// helpers
func isKeywordLead(ch rune) bool {
	for _, keyword := range KeywordLeads {
		if ch == keyword {
			return true
		}
	}

	return false
}

// for keyword token types, its range varies from 10 - 50
const (
	TokenComment TokenType = 10
)

// NewCommentToken -
func NewCommentToken(buf []rune, isMultiLine bool) *Token {
	cpBuf := util.Copy(buf)
	return &Token{
		Type:    TokenComment,
		Literal: cpBuf,
		Info: map[string]bool{
			"isMultiLine": isMultiLine,
		},
	}
}

//// 2. markers
// declare marks
const (
	Comma          rune = 0xFF0C //，
	Colon          rune = 0xFF1A //：
	Semicolon      rune = 0xFF1B //；
	QuestionMark   rune = 0xFF1F //？
	RefMark        rune = 0x0026 // &
	BangMark       rune = 0xFF01 // ！
	AnnotationMark rune = 0x0040 // @
	HashMark       rune = 0x0023 // #
	EmMark         rune = 0x2014 // —
	EllipsisMark   rune = 0x2026 // …
)

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
	LeftBracket   rune = 0x3010 // 【
	RightBracket  rune = 0x3011 // 】
	LeftParen     rune = 0xFF08 // （
	RightParen    rune = 0xFF09 // ）
	VarRemark     rune = 0x00B7 // ·
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

//// 5. numbers
func isNumber(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}
