package exec

import (
	"fmt"
	"testing"
)

type decimalCase struct {
	name        string
	input       string
	expectError bool
	value       string
}

func TestDecimal_ParseValue(t *testing.T) {
	cases := []decimalCase{
		{
			name:  "normal number only",
			input: "000123456789",
			value: "(false, 123456789, 0)",
		},
		{
			name:  "negative number",
			input: "-12345678",
			value: "(true, 12345678, 0)",
		},
		{
			name:  "number with underscore",
			input: "2345_000_000_000",
			value: "(false, 2345000000000, 0)",
		},
		{
			name:  "0 only",
			input: "000",
			value: "(false, 0, 0)",
		},
		{
			name:  "positive-sign number",
			input: "+123456",
			value: "(false, 123456, 0)",
		},
		{
			name:  "with dots",
			input: "-123.4567",
			value: "(true, 1234567, -4)",
		},
		{
			name:  "dot first",
			input: ".23",
			value: "(false, 23, -2)",
		},
		{
			name:  "dot first - negative number",
			input: "-.234",
			value: "(true, 234, -3)",
		},
		{
			name:  "dot first - positive-sign number",
			input: "+.2345",
			value: "(false, 2345, -4)",
		},
		{
			name:  "0.000x",
			input: "000.0000000",
			value: "(false, 0, -7)",
		},
		{
			name:  "*10^",
			input: "-2345.678*10^38",
			value: "(true, 2345678, 35)",
		},
		{
			name:  "Exp negative",
			input: "-0.23458E-5",
			value: "(true, 23458, -10)",
		},
		{
			name:  "Exp positive",
			input: ".56e+9",
			value: "(false, 56, 7)",
		},
	}

	for _, tt := range cases {
		zd := new(ZnDecimal)
		t.Run(tt.name, func(t *testing.T) {
			err := zd.SetValue(tt.input)
			if err != nil {
				if tt.expectError == false {
					t.Errorf("expect no error, got error: %s", err.Error())
				}
			} else {
				output := stringify(zd)
				if tt.expectError == true {
					t.Errorf("expect error, got no error")
				} else if output != tt.value {
					t.Errorf("expect value: %s, got: %s", tt.value, output)
				}
			}
		})
	}
}

func stringify(zd *ZnDecimal) string {
	return fmt.Sprintf("(%v, %s, %d)", zd.sign, zd.co.String(), zd.exp)
}
