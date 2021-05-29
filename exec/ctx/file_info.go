package ctx

import "github.com/DemoHn/Zn/lex"

// FileInfo records current file info, usually for displaying error
type FileInfo struct {
	//// lexical scope
	// file - current execution file directory
	File string
	// currentLine - current exeuction line
	CurrentLine int
	// lineStack - lexical info of (parsed) current file
	LineStack *lex.LineStack
}

func InitFileInfo(l *lex.Lexer) *FileInfo {
	return &FileInfo{
		File:        l.InputStream.GetFile(),
		CurrentLine: 0,
		LineStack:   l.LineStack,
	}
}

// SetCurrentLine -
func (f *FileInfo) SetCurrentLine(line int) {
	f.CurrentLine = line
}
