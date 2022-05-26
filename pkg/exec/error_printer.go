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
	err *zerr.Error
}

func WrapSyntaxError(lexer *syntax.Lexer, moduleName string, err error) *SyntaxErrorWrapper {
	var e *zerr.Error
	if errZ, ok := err.(*zerr.Error); ok {
		e = errZ
	}

	return &SyntaxErrorWrapper{
		lexer:      lexer,
		moduleName: moduleName,
		err:        e,
	}
}

func (sw *SyntaxErrorWrapper) Error() string {
	info := ExtraInfo{
		ErrorClass: "语法错误",
		ModuleName: sw.moduleName,
	}

	if sw.lexer != nil {
		if zerr.IsSyntaxError(sw.err) {
			startIdx := sw.err.Extra.(int)
			lineIdx := sw.lexer.FindLineIdx(startIdx, 0)
			// read text from line
			info.Text, info.ColNum = extractLineTextAndCursor(sw.lexer, startIdx)
			info.LineNum = lineIdx + 1
		} else if zerr.IsLexError(sw.err) {
			startIdx := sw.lexer.GetCursor()
			lineIdx := sw.lexer.FindLineIdx(startIdx, 0)
			// read text from line
			info.Text, info.ColNum = extractLineTextAndCursor(sw.lexer, startIdx)
			info.LineNum = lineIdx + 1
		}
	}
	// show all error lines
	var mask uint16 = 0x0000
	return printError(sw.err, mask, info)
}

func extractLineTextAndCursor(l *syntax.Lexer, cursorIdx int) (string, int) {
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
	return lineText, cursorIdx - startIdx
}

// printError - display detailed error info to user
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
func printError(e *zerr.Error, mask uint16, info ExtraInfo) string {
	var line1, line2, line3, line4 string
	// line1
	if onMask(mask, dpHideFileName) {
		if onMask(mask, dpHideLineNum) {
			line1 = "发生异常："
		} else {
			line1 = fmt.Sprintf("在第 %d 行发生异常：", info.LineNum)
		}
	} else if onMask(mask, dpHideLineNum) {
		line1 = fmt.Sprintf("在「%s」模块中发生异常：", info.ModuleName)
	} else {
		line1 = fmt.Sprintf("在「%s」模块中，位于第 %d 行发生异常：", info.ModuleName, info.LineNum)
	}
	// line2
	if onMask(mask, dpHideLineText) {
		line2 = ""
	} else {
		line2 = fmt.Sprintf("    %s", info.Text)
	}
	// line3
	if onMask(mask, dpHideLineText) || onMask(mask, dpHideLineCursor) {
		line3 = ""
		if !onMask(mask, dpHideLineText) {
			line3 = "    "
		}
	} else {
		line3 = fmt.Sprintf("   %s^", strings.Repeat(" ", calcCursorOffset(info.Text, info.ColNum)+1))
	}
	// line4
	if onMask(mask, dpHideErrClass) {
		line4 = e.Message
	} else {
		errClassText := fmt.Sprintf("[%04X] %s", e.Code, info.ErrorClass)
		line4 = fmt.Sprintf("%s：%s", errClassText, e.Message)
	}

	lines := []string{line1, line2, line3, line4}
	var texts []string
	for _, line := range lines {
		if line != "" {
			texts = append(texts, line)
		}
	}
	return strings.Join(texts, "\n")
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

// declare display masks
//                    16 8 4 2 1
// X X X X X X X X X X O O O O O
const (
	dpHideFileName   uint16 = 0x0001
	dpHideLineCursor uint16 = 0x0002
	dpHideLineNum    uint16 = 0x0004
	dpHideLineText   uint16 = 0x0008
	dpHideErrClass   uint16 = 0x0010
)