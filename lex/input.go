package lex

import (
	"io/ioutil"
	"os"
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

// ReadTextFromFile - read code text from file
func (s *Source) ReadTextFromFile(absPath string) ([]byte, *error.Error) {
	// stat if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return []byte{}, error.FileNotFound(absPath)
	}

	// open file
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return []byte{}, error.FileOpenError(absPath, err)
	}

	return data, nil
}

// init global source object
func init() {
	source = Source{
		Inputs: []SourceInput{},
	}
}
