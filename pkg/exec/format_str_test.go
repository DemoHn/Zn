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
		{
			"HK{#.2E}-{}",
			[]runtime.Element{
				value.NewNumber(13.208945),
				value.NewString("香港记者"),
			},
			"HK1.32E+01-香港记者",
		},
		{
			"HK{#+E}-{}",
			[]runtime.Element{
				value.NewNumber(973.208945),
				value.NewString("香港记者"),
			},
			"HK+9.732089E+02-香港记者",
		},
		{
			"HK{#%}-{}",
			[]runtime.Element{
				value.NewNumber(0.13208945),
				value.NewString("香港记者"),
			},
			"HK13.2089%-香港记者",
		},
		{
			"HK{#.0}-{}",
			[]runtime.Element{
				value.NewNumber(13.208945),
				value.NewString("香港记者"),
			},
			"HK13-香港记者",
		},
		{
			"HK{#.0%}-{}",
			[]runtime.Element{
				value.NewNumber(0.13208945),
				value.NewString("香港记者"),
			},
			"HK13%-香港记者",
		},
		{
			"HK{#}-{}",
			[]runtime.Element{
				value.NewNumber(-1234.56789),
				value.NewString("香港记者"),
			},
			"HK-1234.57-香港记者",
		},
		{
			"HK{#.2E}-{}",
			[]runtime.Element{
				value.NewNumber(-1234.56789),
				value.NewString("香港记者"),
			},
			"HK-1.23E+03-香港记者",
		},
		{
			"HK{#+E}-{}",
			[]runtime.Element{
				value.NewNumber(-1234.56789),
				value.NewString("香港记者"),
			},
			"HK-1.234568E+03-香港记者",
		},
		{
			"HK{#%}-{}",
			[]runtime.Element{
				value.NewNumber(-0.123456789),
				value.NewString("香港记者"),
			},
			"HK-12.3457%-香港记者",
		},
		{
			"HK{#.8}-{}",
			[]runtime.Element{
				value.NewNumber(-1234.56789),
				value.NewString("香港记者"),
			},
			"HK-1234.56789000-香港记者",
		},
		{
			"HK{#.0%}-{}",
			[]runtime.Element{
				value.NewNumber(-0.123456789),
				value.NewString("香港记者"),
			},
			"HK-12%-香港记者",
		},
	}

	for _, c := range cases {
		paramArr := value.NewArray(c.params)
		res, err := formatString(value.NewString(c.formatter), paramArr)

		if err != nil {
			t.Errorf("formatString('%s'): expect '%s', got error: %s", c.formatter, c.expected, err.Error())
		} else if res.String() != c.expected {
			t.Errorf("formatString('%s'): expect '%s', result: '%s'", c.formatter, c.expected, res.String())
		}
	}
}
