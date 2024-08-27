package exec

import (
	"fmt"
	"testing"

	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type fmtCase struct {
	formatter string
	params    []runtime.Element
}

func TestFormatStr(t *testing.T) {
	cases := []fmtCase{
		{
			"HK{}-{xyz}",
			[]runtime.Element{
				value.NewString("香港记者"),
				value.NewNumber(12.208945),
			},
		},
	}

	for _, c := range cases {
		paramArr := value.NewArray(c.params)
		res, _ := formatString(value.NewString(c.formatter), paramArr)
		fmt.Println(res)
	}
}
