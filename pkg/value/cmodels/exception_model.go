package cmodels

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

func NewExceptionModel() *value.ClassModel {
	constructorFunc := value.NewFunction(nil, func(c *r.Context, values []r.Element) (r.Element, error) {
		if err := value.ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}

		message := values[0].(*value.String)
		return value.NewException(message.String()), nil
	})

	return value.NewClassModel("异常", nil).
		SetConstructorFunc(constructorFunc).
		DefineProperty("内容", value.NewString(""))
}
