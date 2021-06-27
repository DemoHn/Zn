package stdlib

import (
	"io/ioutil"
	"os"
	"runtime"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

var libfileValueMap = map[string]ctx.Value{
	"换行符":     val.NewString("\n"),
	"打开文件":    val.NewFunction("打开文件", openFileFunc),
	"自文件读取文本": val.NewFunction("自文件读取文本", readTextFromFileFunc),
	"自文件读取数据": val.NewFunction("自文件读取数据", readDataFromFileFunc),
	"自文件写入文本": val.NewFunction("自文件写入文本", writeTextFromFileFunc),
	"自文件写入数据": val.NewFunction("自文件写入数据", writeDataFromFileFunc),
}

var fileObjectRef val.ClassRef

func openFileFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	// validate one param: string ONLY
	if err := val.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	fileName := values[0]
	obj := val.NewObject(fileObjectRef)
	obj.SetPropList(map[string]ctx.Value{
		"文件名称": fileName,
	})

	return obj, nil
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

func readDataFromFileFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
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

	arr := val.NewArray([]ctx.Value{})
	for _, b := range data {
		arr.AppendValue(val.NewDecimalFromInt(int(b), 0))
	}
	return arr, nil
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

func writeDataFromFileFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	// validate one param: string ONLY
	if err := val.ValidateExactParams(values, "string", "array"); err != nil {
		return nil, err
	}
	fileName := values[0].(*val.String)
	content := values[1].(*val.Array)

	// convert content to byte array
	byteArr := []byte{}
	for _, item := range content.GetValue() {
		item, ok := item.(*val.Decimal)
		if !ok {
			return nil, error.NewErrorSLOT("输入数据不符合要求")
		}
		intVal, err := item.AsInteger()
		if err != nil {
			return nil, error.NewErrorSLOT("输入数据不符合要求")
		}
		if intVal > 255 || intVal < 0 {
			return nil, error.NewErrorSLOT("输入数据不符合要求")
		}

		byteArr = append(byteArr, byte(intVal))
	}

	err := ioutil.WriteFile(fileName.String(), byteArr, 0644)
	if err != nil {
		return nil, error.NewErrorSLOT("写入文件失败" + err.Error())
	}
	return nil, nil
}

func initFileObjectRef() val.ClassRef {
	ref := val.NewClassRef("文件对象")
	ref.MethodList = map[string]val.ClosureRef{
		"读取文本": val.NewClosure(nil, func(c *ctx.Context, v []ctx.Value) (ctx.Value, *error.Error) {
			thisValue := c.GetScope().GetThisValue()
			thisObj := thisValue.(*val.Object)

			objProps := thisObj.GetPropList()
			fileName := objProps["文件名称"]

			return readTextFromFileFunc(c, []ctx.Value{fileName})
		}),
	}
	return ref
}

func init() {
	// 确定换行符
	lineSep := "\n"
	if runtime.GOOS == "windows" {
		lineSep = "\r\n"
	}
	libfileValueMap["换行符"] = val.NewString(lineSep)

	fileObjectRef = initFileObjectRef()
}
