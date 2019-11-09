package tokens

import "fmt"

// declare chars
const (
	SPACE rune = 0x0020 // <SP>
	TAB   rune = 0x0009 // <TAB>
	CR    rune = 0x000D // \r
	LF    rune = 0x000A // \n
)

// EOFToken is an artifical sign added manually to the end of tokens
// It indicates the end of file, and of the token sequence.
type EOFToken struct{}

func (eof EOFToken) String(detailed bool) string {
	return "<EOF>"
}

// NULLToken - indicates nothing
// normally it should not exist
type NULLToken struct{}

func (nul NULLToken) String(detailed bool) string {
	return "<NUL>"
}

// NewLineToken - the token that indicates a new line
type NewLineToken struct {
	start int
	end   int
}

func (nl NewLineToken) String(detailed bool) string {
	return "<CRLF>"
}

// IndentToken - marks the indent at the beginning of newline
type IndentToken struct {
	Space int
	start int
	end   int
}

func (idt IndentToken) String(detailed bool) string {
	return fmt.Sprintf("Indent<%d>", idt.Space)
}
