package error

// Error model
type Error struct {
	code   uint16
	text   string
	cursor *Cursor
	info   interface{}
}

// Error - display error text
func (e *Error) Error() string {
	return e.text
}

// GetCode - get error code
func (e *Error) GetCode() uint16 {
	return e.code
}

// NewError model - only set code - usually for tests
func NewError(code uint16) *Error {
	return &Error{
		code:   code,
		text:   "",
		cursor: nil,
		info:   nil,
	}
}

// Cursor denotes the indicator where the error occurs
type Cursor struct {
	file   string
	line   int
	column int
}

// ErrorClass defines the prefix of error code
type errorClass struct {
	prefix uint8
}

// NewError - new error with subcode
func (ec *errorClass) NewError(subcode uint8, model Error) Error {
	var code uint16
	// code = prefix << 8 + subcode
	code = uint16(ec.prefix)
	code = code*256 + uint16(subcode)

	model.code = code
	return model
}

// definitions of all error classes inside the Zn Programming language
var (
	// 0x10 - inputError
	// this error class handles all errors before transforming file inputs
	// to utf-8 encoding text
	inputError errorClass
	// 0x11 - lexError
	// this error class displays all errors occur during lexing stage.
	lexError errorClass
)

// NewErrorSLOT - a tmp placeholder for adding errors quickly while the
// details has not been thought carefully.
// 简单来说就是临时错误的坑
func NewErrorSLOT(text string) *Error {
	return &Error{
		code:   0xFFFE,
		text:   text,
		cursor: nil,
		info:   nil,
	}
}
func init() {
	inputError = errorClass{0x10}
	lexError = errorClass{0x11}
}
