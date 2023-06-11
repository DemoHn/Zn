package stdlib

import (
	"io/ioutil"
	"os"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/value"

	r "github.com/DemoHn/Zn/pkg/runtime"
)

var fileModuleName = "文件"
var fileModule = r.NewModule(fileModuleName, nil)

func readTextFromFileFunc(c *r.Context, values []r.Value) (r.Value, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	v := values[0].(*value.String)
	// open file
	file, err := os.Open(v.String())
	if err != nil {
		return nil, zerr.NewRuntimeException("打开文件失败：" + err.Error())
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, zerr.NewRuntimeException("读取文件失败：" + err.Error())
	}
	return value.NewString(string(data)), nil
}

func writeTextFromFileFunc(c *r.Context, values []r.Value) (r.Value, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	fileName := values[0].(*value.String)
	content := values[1].(*value.String)

	err := ioutil.WriteFile(fileName.String(), []byte(content.String()), 0644)
	if err != nil {
		return nil, zerr.NewRuntimeException("写入文件失败：" + err.Error())
	}
	return nil, nil
}

func readDirFunc(c *r.Context, values []r.Value) (r.Value, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	dirName := values[0].(*value.String)
	dirs, err := os.ReadDir(dirName.String())
	if err != nil {
		return nil, zerr.NewRuntimeException("读取目录失败：" + err.Error())
	}

	info := value.NewArray([]r.Value{})
	for _, dir := range dirs {
		info.AppendValue(value.NewString(dir.Name()))
	}
	return info, nil
}

/*
func init() {
	// register functions
	RegisterFunctionForModule(fileModule, "读取文件", readTextFromFileFunc)
	RegisterFunctionForModule(fileModule, "写入文件", writeTextFromFileFunc)
	RegisterFunctionForModule(fileModule, "读取目录", readDirFunc)

	RegisterModule(fileModuleName, fileModule)
}
*/
