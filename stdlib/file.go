package stdlib

import (
	"io/ioutil"
	"os"

	"github.com/DemoHn/Zn/pkg/value"

	r "github.com/DemoHn/Zn/pkg/runtime"
)

var fileModuleName = "文件"
var fileModule = r.NewInternalModule(fileModuleName)

func readTextFromFileFunc(c *r.Context, values []r.Element) (r.Element, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	v := values[0].(*value.String)
	// open file
	file, err := os.Open(v.String())
	if err != nil {
		return nil, value.NewException("打开文件失败：" + err.Error())
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, value.NewException("读取文件失败：" + err.Error())
	}
	return value.NewString(string(data)), nil
}

func writeTextFromFileFunc(c *r.Context, values []r.Element) (r.Element, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	fileName := values[0].(*value.String)
	content := values[1].(*value.String)

	err := ioutil.WriteFile(fileName.String(), []byte(content.String()), 0644)
	if err != nil {
		return nil, value.NewException("写入文件失败：" + err.Error())
	}
	return nil, nil
}

func readDirFunc(c *r.Context, values []r.Element) (r.Element, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	dirName := values[0].(*value.String)
	dirs, err := os.ReadDir(dirName.String())
	if err != nil {
		return nil, value.NewException("读取目录失败：" + err.Error())
	}

	info := value.NewArray([]r.Element{})
	for _, dir := range dirs {
		info.AppendValue(value.NewString(dir.Name()))
	}
	return info, nil
}

func init() {
	// register functions
	RegisterFunctionForModule(fileModule, "读取文件", readTextFromFileFunc)
	RegisterFunctionForModule(fileModule, "写入文件", writeTextFromFileFunc)
	RegisterFunctionForModule(fileModule, "读取目录", readDirFunc)

	// 2023/6/11 - NOT add this module to the standard library for now
	//	RegisterModule(fileModuleName, fileModule)
}
