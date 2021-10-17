package val

import "testing"

type strReplaceCase struct {
	name  string
	input string
	value string
}

func Test_String_ReplaceSepcialChars(t *testing.T) {
	cases := []strReplaceCase{
		{
			name:  "plain string",
			input: "Hello World",
			value: "Hello World",
		},
		{
			name:  "only with /CR",
			input: "{/CR}",
			value: "\r",
		},
		{
			name:  "only with /CR/LF",
			input: "{/CR}{/LF}",
			value: "\r\n",
		},
		{
			name:  "with unicode hexcode",
			input: "天下{/+4e4b}最{/+5927}",
			value: "天下之最大",
		},
		{
			name:  "invalid bracket #1",
			input: "boolean{ABC}",
			value: "boolean{ABC}",
		},
		{
			name:  "invalid bracket #2",
			input: "boolean{/CR",
			value: "boolean{/CR",
		},
		{
			name:  "nested bracket",
			input: "bool{/TAB}ean{/CR{/LF}}",
			value: "bool\tean{/CR\n}",
		},
		{
			name:  "express {/LF} exactly",
			input: "{{/s}LF}",
			value: "{/LF}",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceSpecialChars(tt.input)
			if result != tt.value {
				t.Errorf("expect: %s, got: %s", tt.value, result)
			}
		})
	}
}
