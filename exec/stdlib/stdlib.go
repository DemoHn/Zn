package stdlib

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

type valueMap map[string]ctx.Value

type funcExecutor func(*ctx.Context, []ctx.Value) (ctx.Value, *error.Error)

// PackageList -
var PackageList = map[string]valueMap{}

func init() {

}
