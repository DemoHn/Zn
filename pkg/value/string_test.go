package value

import (
	"testing"
)

func TestString_ReplaceSpecialChars(t *testing.T) {
	// cases
	cases := [][]string{
		{"香港记者跑得快", "香港记者跑得快"},
		{"too young`LF`, too simple", "too young\n, too simple"},
		{"too young`U+8DD1`, `U+5f97`too simple", "too young跑, 得too simple"},
		{"too young`CRLF``CRLF` hello", "too young\r\n\r\n hello"},
		{"too young`CRLF``BK` hello", "too young\r\n` hello"},
		{"too young`CBS`BK`BK` hello", "too young`CBS`BK` hello"},
	}

	for _, c := range cases {
		result := replaceSpecialChars(c[0])
		if result != c[1] {
			t.Errorf("replace '%s': expect '%s', result: '%s'", c[0], c[1], result)
		}
	}
}

func TestString_StrExecAtoi(t *testing.T) {
	cases := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{input: "123", expected: 123, hasError: false},
		{input: "123kg", expected: 0, hasError: true},
		{input: "abc123", expected: 0, hasError: true},
		{input: "+123.456", expected: 123.456, hasError: false},
		{input: "-123.456", expected: -123.456, hasError: false},
		{input: "0.456", expected: 0.456, hasError: false},
		{input: "-0.456", expected: -0.456, hasError: false},
		{input: "123.456e2", expected: 12345.6, hasError: false},
		{input: "123.486*10^3", expected: 123486, hasError: false},
		{input: "-123.456e-2", expected: -1.23456, hasError: false},
		{input: "1.23e4", expected: 12300, hasError: false},
		{input: "-1.23e-4", expected: -0.000123, hasError: false},
		{input: "1.23e", expected: 0, hasError: true},
		{input: "1.23e+", expected: 0, hasError: true},
		{input: "1.23e-+", expected: 0, hasError: true},
		{input: "1.23e-4.5", expected: 0, hasError: true},
	}

	for _, c := range cases {
		result, err := strExecAtoi(NewString(c.input), nil, nil)
		if c.hasError {
			if err == nil {
				t.Errorf("strExecAtoi('%s'): expect error, got result: '%f'", c.input, result.(*Number).GetValue())
			}
		} else {
			if err != nil {
				t.Errorf("strExecAtoi('%s'): expect '%f', got error: '%s'", c.input, c.expected, err)
			} else if numResult, ok := result.(*Number); !ok {
				t.Errorf("strExecAtoi('%s'): expect '*Number', got '%T'", c.input, result)
			} else if numResult.GetValue() != c.expected {
				t.Errorf("strExecAtoi('%s'): expect '%f', result: '%f'", c.input, c.expected, numResult.GetValue())
			}
		}
	}
}
