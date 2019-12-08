package zntex

// Parser -
type Parser struct {
	currentPos int
	peekPos    int
	quoteStack *RuneStack
	chBuffer   []rune
	input      []rune
}

// common consts
const (
	EOF = 0
	LF  = 10
)

// Token defines an abstract type of token
type Token interface {
	literal() []rune
}

// TextToken is the normal text
type TextToken struct {
	Text []rune
}

func (t *TextToken) literal() []rune {
	return t.Text
}

// CommentToken parse a one-line token
type CommentToken struct {
	Text []rune
}

func (t *CommentToken) literal() []rune {
	return t.Text
}

// CommandToken defines a general type of command (including its options and arguments)
type CommandToken struct {
	Literal []rune
	Command []rune
	Options [][]rune
	Args    [][]rune
}

func (t *CommandToken) literal() []rune {
	return t.Literal
}

// EnvironToken defines an envrionment
type EnvironToken struct {
	Literal []rune
	IsBegin bool
	Tag     []rune
	Options [][]rune
	Args    [][]rune
}

func (t *EnvironToken) literal() []rune {
	return t.Literal
}

//// public methods

// NextToken - get next token
func (p *Parser) NextToken() Token {
	ch := p.next()
	pch := p.peek()
	switch ch {
	case '\\':
		switch pch {
		case '\\':
			p.next()
			return p.parseText(LF)
		case '{', '}', '[', ']', '%':
			p.next()
			return p.parseText(pch)
		default:
			return p.parseCommand(ch)
		}
	case EOF:
		return nil
	case '%':
		return p.parseComment(ch)
	default:
		return p.parseText(ch)
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
func (p *Parser) parseCommand(ch rune) Token {
	return nil
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
		case '\r', '\n':
			if p.peek() == '\r' && p.peek2() == '\n' {
				p.next()
			}
			p.next()
		case '%', EOF:
			return &TextToken{
				Text: Copy(p.chBuffer),
			}
		default:
			ch = p.next()
			p.pushBuffer(ch)
		}
	}
}
