package error

// Error model
type Error struct {
	code   uint16
	text   string
	cursor Cursor
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

// SetCursor - set error occurance location
// for better readability
func (e *Error) SetCursor(cursor Cursor) {
	e.cursor = cursor
}

// Display - display detailed error info to user
// general format:
//
// 在 [FILE] 中，位于第 [LINE] 行：
//     [ LINE TEXT WITHOUT INDENTS AND CRLF ]
// ‹[ERRCODE]› [ERRCLASS]：[ERRTEXT]
//
// example error:
//
// 在 draft/example.zn 中，位于第 12 行：
//     如果代码不为空：
//    ^
// ‹2021› 语法错误：此行现行缩进类型为「TAB」，与前设缩进类型「空格」不符！
func (e *Error) Display() string {
	return ""
}

// Cursor denotes the indicator where the error occurs
type Cursor struct {
	file    string
	lineNum int
	colNum  int
	text    string
}

// ErrorClass defines the prefix of error code
type errorClass struct {
	prefix uint8
}

// NewError - new error with subcode
func (ec *errorClass) NewError(subcode uint8, model Error) *Error {
	var code uint16
	// code = prefix << 8 + subcode
	code = uint16(ec.prefix)
	code = code*256 + uint16(subcode)

	model.code = code
	return &model
}

// definitions of all error classes inside the Zn Programming language
var (
	// 0x20 - lexError
	// this error class displays all errors occur during lexing stage. (including input data)
	lexError errorClass
	// 0x21 - ioError
	// I/O related error (e.g. FileNotFound, OpenFileError, ReadFileError)
	ioError errorClass
)

// NewErrorSLOT - a tmp placeholder for adding errors quickly while the
// details has not been thought carefully.
// 简单来说就是临时错误加个坑位，等到正式写代码的时候再用
func NewErrorSLOT(text string) *Error {
	return &Error{
		code: 0xFFFE,
		text: text,
		info: nil,
	}
}

func init() {
	lexError = errorClass{0x20}
	ioError = errorClass{0x21}
}
