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

// Cursor denotes the indicator where the error occurs
type Cursor struct {
	file   string
	line   int32
	column int32
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
)

func init() {
	inputError = errorClass{0x10}
}
