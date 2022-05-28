package exec

import (
	"fmt"
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
	"strings"
)

type ExtraInfo struct {
	ModuleName string
	File    string
	LineNum int
	ColNum  int
	Text    string
	ErrorClass string
}

// SyntaxErrorWrapper - wrap IO errors with file info (current lexer etc.)
type SyntaxErrorWrapper struct {
	lexer *syntax.Lexer
	moduleName string
	err error
}

func WrapSyntaxError(lexer *syntax.Lexer, moduleName string, err error) *SyntaxErrorWrapper {
	return &SyntaxErrorWrapper{
		lexer:      lexer,
		moduleName: moduleName,
		err:        err,
	}
}

func (sw *SyntaxErrorWrapper) Error() string {
	errClass := "语法错误"
	var errLines []string

	code := 0
	if sw.lexer != nil {
		if serr, ok := sw.err.(*zerr.SyntaxError); ok {
			code = serr.Code
			lineIdx := sw.lexer.FindLineIdx(serr.Cursor, 0)
			// add line 1
			errLines = append(errLines, fmtErrorLocationHeadLine(sw.moduleName, lineIdx + 1))
			// add line 2
			errLines = append(errLines, fmtErrorSourceTextLine(sw.lexer, serr.Cursor, true))
		}
	}

	if sw.err != nil {
		errLines = append(errLines, fmtErrorMessageLine(code, errClass, sw.err.Error()))
	}

	return strings.Join(errLines, "\n")
}


// print error lines - display detailed error info to user
// general format:
//
// 在 [FILE] 模块中，位于第 [LINE] 行：
//     [ LINE TEXT WITHOUT INDENTS AND CRLF ]
// [[ERRCODE]] [ERRCLASS]：[ERRTEXT]
//
// example error:
//
// 在 draft/example.zn 中，位于第 12 行：
//     如果代码不为空：
//    ^
// [2021] 语法错误：此行现行缩进类型为「TAB」，与前设缩进类型「空格」不符！

// fmtErrorLocationHeadLine -
// e.g. 在 draft/example.zn 中，位于第 12 行发生异常：
func fmtErrorLocationHeadLine(moduleName string, lineNum int) string {
	return fmt.Sprintf("在「%s」模块中，位于第 %d 行发生异常：", moduleName, lineNum)
}

// fmtErrorSourceTextLine -
// cursorIdx: global index inside the source text from denoted lexer
// if withCursorMark == false, hide the "^" mark that indicates the specific location where error occurs.
// e.g.:
//     如果代码不为空：
//        ^
func fmtErrorSourceTextLine(l *syntax.Lexer, cursorIdx int, withCursorMark bool) string {
	startIdx := cursorIdx
	endIdx := startIdx
	// find prev until meeting first CR/LF
	for startIdx > 0 {
		if l.Source[startIdx] == syntax.RuneCR || l.Source[startIdx] == syntax.RuneLF {
			startIdx += 1
			// skip indent chars
			for l.Source[startIdx] == syntax.RuneSP || l.Source[startIdx] == syntax.RuneTAB {
				startIdx += 1
			}
			break
		}
		startIdx -= 1
	}
	// find next until meeting first CR/LF
	for endIdx < len(l.Source) {
		if l.Source[endIdx] == syntax.RuneCR || l.Source[endIdx] == syntax.RuneLF {
			break
		}
		endIdx += 1
	}

	// get relative cursor offset (notice one Chinese char counts for 2 unit offsets)
	lineText := string(l.Source[startIdx:endIdx])
	fmtLine := fmt.Sprintf("    %s", lineText)
	if withCursorMark {
		cursorText := fmt.Sprintf("\n    %s^", strings.Repeat(" ", calcCursorOffset(lineText, cursorIdx - startIdx)))
		fmtLine += cursorText
	}

	return fmtLine
}

// fmtErrorMessageLine - format error message line
// NOTE: if code == 0, "[code]" is not shown
// [<code>] <errName>：<errMessage>
// e.g.: [2021] 语法错误：此行现行缩进类型为「TAB」，与前设缩进类型「空格」不符！
func fmtErrorMessageLine(code int, errName string, errMessage string) string {
	fmtCode := ""
	if code != 0 {
		fmtCode = fmt.Sprintf("[%d]", code)
	}
	return fmt.Sprintf("%s%s：%s", errName, fmtCode, errMessage)
}

func calcCursorOffset(text string, col int) int {
	if col < 0 {
		return col
	}
	widthBorders := []int32{
		126, 159, 687, 710, 711, 727, 733, 879, 1154, 1161,
		4347, 4447, 7467, 7521, 8369, 8426, 9000, 9002, 11021, 12350,
		12351, 12438, 12442, 19893, 19967, 55203, 63743, 64106, 65039, 65059,
		65131, 65279, 65376, 65500, 65510, 120831, 262141, 1114109,
	}

	widths := []int{
		1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
		1, 2, 1, 0, 1, 0, 1, 2, 1, 2,
		1, 2, 0, 2, 1, 2, 1, 2, 1, 0,
		2, 1, 2, 1, 2, 1, 2, 1,
	}

	offsets := 0

	getOffset := func(t rune) int {
		if t == 0xE || t == 0xF {
			return 0
		}
		for idx, b := range widthBorders {
			if t <= b {
				return widths[idx]
			}
		}
		return 1
	}
	for _, t := range []rune(text)[:col] {
		offsets = offsets + getOffset(t)
	}

	return offsets
}

func onMask(target uint16, mask uint16) bool {
	return (target & mask) > 0
}