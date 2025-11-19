package file

import (
	"io"
	"os"

	"github.com/DemoHn/Zn/pkg/value"

	r "github.com/DemoHn/Zn/pkg/runtime"
)

const FILE_LIB_NAME = "@文件"

var fileLIB *r.Library

func FN_readTextFromFile(receiver r.Element, values []r.Element) (r.Element, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	v := values[0].(*value.String)
	// open file
	file, err := os.Open(v.String())
	if err != nil {
		return nil, value.ThrowException("打开文件失败：" + err.Error())
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, value.ThrowException("读取文件失败：" + err.Error())
	}
	return value.NewString(string(data)), nil
}

func FN_writeTextFromFile(receiver r.Element, values []r.Element) (r.Element, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	fileName := values[0].(*value.String)
	content := values[1].(*value.String)

	err := os.WriteFile(fileName.String(), []byte(content.String()), 0644)
	if err != nil {
		return nil, value.ThrowException("写入文件失败：" + err.Error())
	}
	return nil, nil
}

func FN_readDir(receiver r.Element, values []r.Element) (r.Element, error) {
	// validate one param: string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	dirName := values[0].(*value.String)
	dirs, err := os.ReadDir(dirName.String())
	if err != nil {
		return nil, value.ThrowException("读取目录失败：" + err.Error())
	}

	info := value.NewArray([]r.Element{})
	for _, dir := range dirs {
		info.AppendValue(value.NewString(dir.Name()))
	}
	return info, nil
}

func Export() *r.Library {
	return fileLIB
}

func init() {
	fileLIB = r.NewLibrary(FILE_LIB_NAME)

	fileLIB.RegisterFunction("读取文件", value.NewFunction(FN_readTextFromFile)).
		RegisterFunction("写入文件", value.NewFunction(FN_writeTextFromFile)).
		RegisterFunction("读取目录", value.NewFunction(FN_readDir))
}
