package exec

import (
	"errors"
	"fmt"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
)

type ExtraInfo struct {
	ModuleName string
	File       string
	LineNum    int
	ColNum     int
	Text       string
	ErrorClass string
}

// SyntaxErrorWrapper - wrap IO errors with file info (current lexer etc.)
type SyntaxErrorWrapper struct {
	parser     *syntax.Parser
	moduleName string
	err        error
}

func WrapSyntaxError(parser *syntax.Parser, moduleName string, err error) error {
	switch serr := err.(type) {
	// DO NOT wrap a wrapped error
	case *SyntaxErrorWrapper, *RuntimeErrorWrapper:
		return serr
	default:
		return &SyntaxErrorWrapper{
			parser:     parser,
			moduleName: moduleName,
			err:        err,
		}
	}
}

func (sw *SyntaxErrorWrapper) Error() string {
	errClass := "语法错误"
	var errLines []string

	code := 0
	if sw.parser != nil {
		if serr, ok := sw.err.(*zerr.SyntaxError); ok {
			code = serr.Code
			lineIdx := sw.parser.FindLineIdx(serr.Cursor, 0)
			// add line 1
			errLines = append(errLines, fmtErrorLocationHeadLine(false, sw.moduleName, lineIdx+1))
			// add line 2
			errLines = append(errLines, fmtErrorSourceLineWithParser(sw.parser, serr.Cursor, true))
		}
	}

	if sw.err != nil {
		errLines = append(errLines, fmtErrorMessageLine(code, errClass, sw.err.Error()))
		return strings.Join(errLines, "\n")
	} else {
		return fmtErrorMessageLine(-99, errClass, "未知语法错误 :-(")
	}
}

type RuntimeErrorWrapper struct {
	vm  *r.VM
	err error
}

func WrapRuntimeError(vm *r.VM, err error) error {
	switch serr := err.(type) {
	// DO NOT wrap a wrapped error
	case *RuntimeErrorWrapper, *SyntaxErrorWrapper:
		return serr
	default:
		exception, realErr := extractSignalValue(err, zerr.SigTypeException)
		if realErr == nil {

			// read exception content
			errContent := ""
			// build an error from exception value
			// for default *value.Exception class
			if objE, ok := exception.(*value.Exception); ok {
				errContent = objE.Error()
			} else {
				// other custom exception class
				s, _ := exception.GetProperty(EVConstExceptionContentProperty)
				if s != nil {
					if ss, ok := s.(*value.String); ok {
						errContent = ss.GetValue()
					}
				}
			}

			return &RuntimeErrorWrapper{
				vm:  vm,
				err: errors.New(errContent),
			}
		}

		return &RuntimeErrorWrapper{
			vm:  vm,
			err: realErr,
		}
	}
}

func (rw *RuntimeErrorWrapper) Error() string {
	errClass := "运行异常"
	var errLines []string
	code := 0

	if werr, ok := rw.err.(*zerr.RuntimeError); ok {
		code = werr.Code
	}

	callStack := rw.vm.GetCallStack()
	if len(callStack) > 0 {
		// append head lines
		headTrace := callStack[0]
		module := headTrace.GetModule()
		isNativeModule := module.GetID() == r.NATIVE_CODE_MODULE_ID || module.GetProgram() == nil
		errLines = append(errLines, fmtErrorLocationHeadLine(isNativeModule, module.GetName(), headTrace.GetCurrentLine()+1))

		// get line text
		if !isNativeModule {
			if lineText := headTrace.GetSourceTextLine(headTrace.GetCurrentLine()); lineText != "" {
				errLines = append(errLines, fmtErrorSourceTextLine(lineText))
			}
		}

		// append body
		for _, tr := range callStack[1:] {
			trModule := tr.GetModule()
			isNativeModule = trModule.GetID() == r.NATIVE_CODE_MODULE_ID || module.GetProgram() == nil
			if trModule != nil {
				errLines = append(errLines, fmtErrorLocationBodyLine(isNativeModule, trModule.GetName(), tr.GetCurrentLine()+1))
				// get line text
				if !isNativeModule {
					if lineText := tr.GetSourceTextLine(tr.GetCurrentLine()); lineText != "" {
						errLines = append(errLines, fmtErrorSourceTextLine(lineText))
					}
				}
			}
		}
	}

	if rw.err != nil {
		errLines = append(errLines, fmtErrorMessageLine(code, errClass, rw.err.Error()))
	}

	return strings.Join(errLines, "\n")
}

func DisplayError(err error) string {
	switch e := err.(type) {
	case *SyntaxErrorWrapper, *RuntimeErrorWrapper:
		return e.Error()
	case *zerr.IOError:
		cls := "IO错误"
		return fmtErrorMessageLine(e.Code, cls, e.Error())
	case *zerr.SyntaxError:
		cls := "语法错误"
		return fmtErrorMessageLine(e.Code, cls, e.Error())
	default:
		return err.Error()
	}
}

// print error lines - display detailed error info to user
// general format:
//
// 在模块[FILE]中，位于第 [LINE] 行：
//     [ LINE TEXT WITHOUT INDENTS AND CRLF ]
// [ERRCLASS][[ERRCODE]]：[ERRTEXT]
//
// example error:
//
// 在模块draft/example.zn 中，位于第 12 行：
//     如果代码不为空：
//    ^
// 语法错误[20]：此行现行缩进类型为「TAB」，与前设缩进类型「空格」不符！

// fmtErrorLocationHeadLine -
// e.g. 在「example」模块中，位于第 12 行发生异常：
func fmtErrorLocationHeadLine(isNativeModule bool, moduleName string, lineNum int) string {
	if isNativeModule {
		return "在 <内置模块> 中发生异常：\n    [ --内部程序-- ]"
	}
	if moduleName == MODULE_NAME_MAIN {
		return fmt.Sprintf("在主模块中，位于第 %d 行发生异常：", lineNum)
	}
	return fmt.Sprintf("在模块“%s”中，位于第 %d 行发生异常：", moduleName, lineNum)
}

// fmtErrorLocationBodyLine -
// e.g. 来自「example2」模块，第 12 行：
func fmtErrorLocationBodyLine(isNativeModule bool, moduleName string, lineNum int) string {
	if isNativeModule {
		return "来自 <内置模块>：\n    [ --内部程序-- ]"
	}
	if moduleName == MODULE_NAME_MAIN {
		return fmt.Sprintf("来自主模块，第 %d 行：", lineNum)
	}

	return fmt.Sprintf("来自“%s”模块，第 %d 行：", moduleName, lineNum)
}

// fmtErrorSourceTextLine -
// cursorIdx: global index inside the source text from denoted lexer
// if withCursorMark == false, hide the "^" mark that indicates the specific location where error occurs.
// e.g.:
//
//	如果代码不为空：
//	   ^
func fmtErrorSourceLineWithParser(p *syntax.Parser, cursorIdx int, withCursorMark bool) string {
	startIdx := cursorIdx
	endIdx := startIdx
	// append EOF to source to avoid index exceed exception
	sourceT := append(p.GetSource(), 0)
	for sourceT[startIdx] == syntax.RuneCR || sourceT[startIdx] == syntax.RuneLF {
		startIdx -= 1
	}
	// find prev until meeting first CR/LF
	for startIdx > 0 {
		if sourceT[startIdx] == syntax.RuneCR || sourceT[startIdx] == syntax.RuneLF {
			startIdx += 1
			// skip indent chars
			for sourceT[startIdx] == syntax.RuneSP || sourceT[startIdx] == syntax.RuneTAB {
				startIdx += 1
			}
			break
		}
		startIdx -= 1
	}
	// find next until meeting first CR/LF
	for endIdx < len(sourceT) {
		if sourceT[endIdx] == syntax.RuneCR || sourceT[endIdx] == syntax.RuneLF {
			break
		}
		endIdx += 1
	}

	// get relative cursor offset (notice one Chinese char counts for 2 unit offsets)
	lineText := string(sourceT[startIdx:endIdx])
	fmtLine := fmt.Sprintf("    %s", lineText)
	if withCursorMark {
		cursorText := fmt.Sprintf("\n    %s^", strings.Repeat(" ", calcCursorOffset(lineText, cursorIdx-startIdx)))
		fmtLine += cursorText
	}

	return fmtLine
}

// fmtErrorSourceTextLine -
// cursorIdx: global index inside the source text from denoted lexer
// if withCursorMark == false, hide the "^" mark that indicates the specific location where error occurs.
// e.g.:
//
//	如果代码不为空：
func fmtErrorSourceTextLine(lineText string) string {
	return fmt.Sprintf("    %s", lineText)
}

// fmtErrorMessageLine - format error message line
// NOTE: if code == 0, "[code]" is not shown
// <errName>[<code>]：<errMessage>
// e.g.: 语法错误[20]：此行现行缩进类型为「TAB」，与前设缩进类型「空格」不符！
func fmtErrorMessageLine(code int, errName string, errMessage string) string {
	fmtCode := ""
	if code != 0 {
		fmtCode = fmt.Sprintf("[%d]", code)
	}
	return fmt.Sprintf("\n%s%s：%s\n", errName, fmtCode, errMessage)
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
