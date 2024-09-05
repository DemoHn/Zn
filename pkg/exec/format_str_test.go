package exec

import (
	"testing"

	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type fmtCase struct {
	formatter string
	params    []runtime.Element
	expected  string
}

func TestFormatStr(t *testing.T) {
	cases := []fmtCase{
		{
			"HK{#}-{}",
			[]runtime.Element{
				value.NewNumber(13.208945),
				value.NewString("香港记者"),
			},
			"HK13.2089-香港记者",
		},
		{
			"HK{#.2}-{}",
			[]runtime.Element{
				value.NewNumber(-13.208945),
				value.NewString("香港记者"),
			},
			"HK-13.21-香港记者",
		},
	}

	for _, c := range cases {
		paramArr := value.NewArray(c.params)
		res, _ := formatString(value.NewString(c.formatter), paramArr)

		if res.String() != c.expected {
			t.Errorf("formatString('%s'): expect '%s', result: '%s'", c.formatter, c.expected, res.String())
		}
	}
}
