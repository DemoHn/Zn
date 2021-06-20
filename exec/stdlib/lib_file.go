package stdlib

import (
	"runtime"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

var libfileValueMap = map[string]ctx.Value{
	"换行符":  val.NewString("\n"),
	"打开文件": val.NewFunction("打开文件", openFileFunc),
}

func openFileFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	return nil, nil
}

func init() {
	// 确定换行符
	lineSep := "\n"
	if runtime.GOOS == "windows" {
		lineSep = "\r\n"
	}
	libfileValueMap["换行符"] = val.NewString(lineSep)
}
