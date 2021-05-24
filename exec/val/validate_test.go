package val

import (
	"testing"

	"github.com/DemoHn/Zn/exec/ctx"
)

type validateCase struct {
	name        string
	values      []ctx.Value
	types       []string
	expectError bool
}

func TestValidateLeastParams(t *testing.T) {
	cases := []validateCase{
		{
			name: "no wildcard",
			values: []ctx.Value{
				NewBool(true),
				NewString("pig"),
			},
			types: []string{
				"bool",
				"string",
			},
			expectError: false,
		},
		{
			name: "with must wildcard - success",
			values: []ctx.Value{
				NewBool(true),
				NewString("pig"),
				NewString("pot"),
			},
			types: []string{
				"bool",
				"string+",
			},
			expectError: false,
		},
		{
			name: "with optional wildcard - success",
			values: []ctx.Value{
				NewBool(true),
			},
			types: []string{
				"bool",
				"string*",
			},
			expectError: false,
		},
		{
			name: "with optional wildcard - success#2",
			values: []ctx.Value{
				NewBool(true),
				NewString("pot"),
				NewString("pig"),
			},
			types: []string{
				"bool",
				"string*",
			},
			expectError: false,
		},
		{
			name: "with optional wildcard - success#2",
			values: []ctx.Value{
				NewBool(true),
				NewString("pot"),
				NewString("pig"),
			},
			types: []string{
				"bool",
				"string*",
			},
			expectError: false,
		},
		{
			name: "with optional wildcard - FAIL - one wildcard mismatch",
			values: []ctx.Value{
				NewBool(true),
				NewString("pot"),
				NewNull(),
			},
			types: []string{
				"bool",
				"string*",
			},
			expectError: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLeastParams(tt.values, tt.types...)
			if err != nil {
				if tt.expectError == false {
					t.Errorf("expect no error, got error: %s", err.Error())
				}
			} else {
				if tt.expectError == true {
					t.Errorf("expect error, got no error")
				}
			}
		})
	}
}
