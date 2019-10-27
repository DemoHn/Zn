package lex

import (
	"unicode/utf8"

	"github.com/DemoHn/Zn/error"
)

// Source stores all source code inputs (whatever from REPL, file, or CLI etc.) as an array.
type Source struct {
	Inputs []SourceInput
}

// SourceInput stores code text with utf-8 encoding
type SourceInput struct {
	Scope string // TODO - define the scope
	Text  []rune
}

// the only instance of source object
var source Source

//// methods

// AddSourceInput transforms and adds one input (e.g. code file) from raw byte format to
// utf-8 encoding source. Throws error if transforming failed
func (s *Source) AddSourceInput(rawData []byte) *error.Error {
	var input = SourceInput{
		// TODO: define scope
		Scope: "",
		Text:  make([]rune, 0),
	}

	b := rawData
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			return error.DecodeUTF8Fail()
		}
		input.Text = append(input.Text, r)
		b = b[size:]
	}

	s.Inputs = append(s.Inputs, input)
	return nil
}

// init global source object
func init() {
	source = Source{
		Inputs: []SourceInput{},
	}
}
