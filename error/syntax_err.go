package error

// InvalidSyntax -
func InvalidSyntax() *Error {
	return syntaxError.NewError(0x50, Error{
		text: "不合规范之语法",
	})
}

// UnexpectedIndent -
func UnexpectedIndent() *Error {
	return syntaxError.NewError(0x51, Error{
		text: "意外出现之缩进",
	})
}

// IncompleteStmt -
func IncompleteStmt() *Error {
	return syntaxError.NewError(0x52, Error{
		text: "语句仍未结束",
	})
}

// ExprMustTypeID -
func ExprMustTypeID() *Error {
	return syntaxError.NewError(0x53, Error{
		text: "表达式须为「泛标识符」〈如‘变量’、‘对象之名’之类〉",
	})
}
