package stdlib

import (
	"io/ioutil"
	"os"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

var libfileValueMap = map[string]ctx.Value{
	"读取文件": val.NewFunction("读取文件", readTextFromFileFunc),
	"写入文件": val.NewFunction("写入文件", writeTextFromFileFunc),
	"读取目录": val.NewFunction("读取目录", readDirFunc),
}

func readTextFromFileFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	// validate one param: string ONLY
	if err := val.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	v := values[0].(*val.String)
	// open file
	file, err := os.Open(v.String())
	if err != nil {
		return nil, error.NewErrorSLOT("打开文件失败：" + err.Error())
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, error.NewErrorSLOT("读取文件失败" + err.Error())
	}
	return val.NewString(string(data)), nil
}

func writeTextFromFileFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	// validate one param: string ONLY
	if err := val.ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	fileName := values[0].(*val.String)
	content := values[1].(*val.String)

	err := ioutil.WriteFile(fileName.String(), []byte(content.String()), 0644)
	if err != nil {
		return nil, error.NewErrorSLOT("写入文件失败" + err.Error())
	}
	return nil, nil
}

func readDirFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	// validate one param: string ONLY
	if err := val.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	dirName := values[0].(*val.String)
	dirs, err := os.ReadDir(dirName.String())
	if err != nil {
		return nil, error.NewErrorSLOT("读取目录失败" + err.Error())
	}

	info := val.NewArray([]ctx.Value{})
	for _, dir := range dirs {
		info.AppendValue(val.NewString(dir.Name()))
	}
	return info, nil
}
