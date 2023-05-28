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
