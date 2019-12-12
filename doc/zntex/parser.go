package zntex

import "fmt"

// Parser -
type Parser struct {
	currentPos int
	peekPos    int
	quoteStack *RuneStack
	chBuffer   []rune
	input      []rune
}

// common char
const (
	EOF = 0
	LF  = 10
)

const (
	typeComment = 1
	typeText    = 2
	typeEnviron = 3
	typeCommand = 4
)

// Token defines an abstract type of token
type Token interface {
	getType() int
}

// TextToken is the normal text
type TextToken struct {
	Text []rune
}

func (t *TextToken) getType() int {
	return typeText
}

// CommentToken parse a one-line token
type CommentToken struct {
	Text []rune
}

func (t *CommentToken) getType() int {
	return typeComment
}

// CommandToken defines a general type of command (including its options and arguments)
type CommandToken struct {
	Literal []rune
	Command string
	Options []string
	Args    []string
}

func (t *CommandToken) getType() int {
	return typeCommand
}

// EnvironToken defines an envrionment
type EnvironToken struct {
	Literal []rune
	IsBegin bool
	Tag     string
	Options []string
	Args    []string
}

func (t *EnvironToken) getType() int {
	return typeEnviron
}

//// public methods

// NextToken - get next token
func (p *Parser) NextToken() (Token, error) {
	var ch = p.next()
	pch := p.peek()
	switch ch {
	case '\r', '\n', '\t', '\v', '\f', ' ':
		p.skipBlanks(ch)
		if p.peek() != EOF {
			return &TextToken{Text: []rune{' '}}, nil
		}
		return nil, nil
	case '\\':
		switch pch {
		case '\\':
			p.next()
			return p.parseText(LF), nil
		case '{', '}', '[', ']', '%':
			p.next()
			return p.parseText(pch), nil
		default:
			return p.parseCommand(ch)
		}
	case EOF:
		return nil, nil
	case '%':
		return p.parseComment(ch), nil
	default:
		return p.parseText(ch), nil
	}
}

//// private helpers
func (p *Parser) next() rune {
	if p.peekPos >= len(p.input) {
		return EOF
	}

	data := p.input[p.peekPos]
	p.currentPos = p.peekPos
	p.peekPos++
	return data
}

func (p *Parser) peek() rune {
	if p.peekPos >= len(p.input) {
		return EOF
	}

	data := p.input[p.peekPos]
	return data
}

func (p *Parser) peek2() rune {
	if p.peekPos+1 >= len(p.input) {
		return EOF
	}

	data := p.input[p.peekPos+1]
	return data
}

func (p *Parser) clearBuffer() {
	p.chBuffer = []rune{}
}

func (p *Parser) pushBuffer(ch rune) {
	p.chBuffer = append(p.chBuffer, ch)
}

// parseCommand (and Environ token)
func (p *Parser) parseCommand(ch rune) (Token, error) {
	// startup
	p.clearBuffer()
	var startPos = p.currentPos
	var s1Text, s4Text string
	var s2Text = make([]string, 0)
	var s3Text = make([]string, 0)
	// state defs:
	// A. COMMANDS
	//
	// \commandName[opt1,opt2]{arg1}{arg2}
	// <--- 1 ----><--- 2 ---><---- 3 --->
	//
	// B. ENVIRONS
	// \begin{tagName}[opt1,opt2]{arg1}{arg2}
	// <--1-><-- 4 --><--- 2 ---><---- 3 --->
	var state = 1
	// iterate
	for {
		ch = p.next()
		switch ch {
		case '\r', '\n':
			continue
		case '{':
			switch state {
			case 1:
				// fetch content of chBuffer
				buf := string(p.chBuffer)
				if buf == "begin" || buf == "end" {
					state = 4
				} else {
					state = 3
				}
				s1Text = buf
				p.clearBuffer()
				p.quoteStack.Push(ch)
				if p.quoteStack.Num() > 1 {
					return nil, fmt.Errorf("nest not supported")
				}
			case 2:
				state = 3
				p.clearBuffer()
				fallthrough
			case 3:
				p.quoteStack.Push(ch)
				if p.quoteStack.Num() > 1 {
					return nil, fmt.Errorf("nest not supported")
				}
			case 4:
				state = 3
				p.quoteStack.Push(ch)
				if p.quoteStack.Num() > 1 {
					return nil, fmt.Errorf("nest not supported")
				}
			default:
				return nil, fmt.Errorf("invalid char")
			}
		case '[':
			switch state {
			case 1:
				s1Text = string(p.chBuffer)
				fallthrough
			case 4:
				state = 2
				p.clearBuffer()
				p.quoteStack.Push(ch)
				if p.quoteStack.Num() > 1 {
					return nil, fmt.Errorf("nest not supported")
				}
			default:
				return nil, fmt.Errorf("invalid char")
			}
		case '}':
			switch state {
			case 3:
				// match quotes
				if data, _ := p.quoteStack.Current(); data == '{' {
					p.quoteStack.Pop()
				} else {
					return nil, fmt.Errorf("invalid } match")
				}
				s3Text = append(s3Text, string(p.chBuffer))
				p.clearBuffer()
			case 4:
				// match quotes
				if data, _ := p.quoteStack.Current(); data == '{' {
					p.quoteStack.Pop()
				} else {
					return nil, fmt.Errorf("invalid } match")
				}
				s4Text = string(p.chBuffer)
				p.clearBuffer()
			default:
				return nil, fmt.Errorf("invalid char")
			}
			if peekEnd(state, p.peek()) {
				goto end
			}
		case ']':
			switch state {
			case 2:
				if data, _ := p.quoteStack.Pop(); data == '[' {
					p.quoteStack.Pop()
				} else {
					return nil, fmt.Errorf("invalid ] quote")
				}
				s2Text = append(s2Text, string(p.chBuffer))
				p.clearBuffer()
			default:
				return nil, fmt.Errorf("invalid state")
			}
		case ',':
			switch state {
			case 2:
				s2Text = append(s2Text, string(p.chBuffer))
				p.clearBuffer()
			default:
				p.pushBuffer(ch)
			}
		default:
			if ch == EOF {
				return nil, fmt.Errorf("invalid end")
			}
			// options不接受空格
			if state == 2 && Contains(ch, []rune{' ', '\t', '\r', '\n'}) {
				continue
			}
			p.pushBuffer(ch)
		}
	}
end:
	var endPos = p.currentPos
	if s1Text == "begin" || s1Text == "end" {
		return &EnvironToken{
			Literal: p.input[startPos : endPos+1],
			IsBegin: s1Text == "begin",
			Tag:     s4Text,
			Options: s2Text,
			Args:    s3Text,
		}, nil
	}

	return &CommandToken{
		Literal: p.input[startPos : endPos+1],
		Command: s1Text,
		Options: s2Text,
		Args:    s3Text,
	}, nil
}

func peekEnd(state int, peekCh rune) bool {
	switch state {
	case 3:
		return !(peekCh == '{')
	case 4:
		return !Contains(peekCh, []rune{'[', '{'})
	default:
		return false
	}
}
func (p *Parser) parseComment(ch rune) Token {
	p.clearBuffer()
	for {
		ch = p.next()
		switch ch {
		case EOF, '\r', '\n':
			if ch == '\r' && p.peek() == '\n' {
				p.next()
			}
			return &CommentToken{
				Text: Copy(p.chBuffer),
			}
		default:
			p.pushBuffer(ch)
		}
	}
}

func (p *Parser) parseText(ch rune) Token {
	// setup
	p.clearBuffer()
	p.pushBuffer(ch)
	// iterate
	for {
		switch p.peek() {
		case '\\':
			switch p.peek2() {
			case '\\':
				p.next()
				p.next()
				p.pushBuffer(LF)
			case '{', '}', '[', ']', '%':
				pch := p.peek2()
				p.next()
				p.next()
				p.pushBuffer(pch)
			default:
				return &TextToken{
					Text: Copy(p.chBuffer),
				}
			}
		case '\r', '\n', '\v', '\f', ' ', '\t', '%', EOF:
			return &TextToken{
				Text: Copy(p.chBuffer),
			}
		default:
			ch = p.next()
			p.pushBuffer(ch)
		}
	}
}

func (p *Parser) skipBlanks(ch rune) {
	var blanks = []rune{
		'\r', '\n', ' ', '\v', '\f', '\t',
	}
	for {
		if !Contains(p.peek(), blanks) {
			break
		}
		ch = p.next()
	}
}
