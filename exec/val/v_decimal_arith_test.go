package val

import "testing"

type arithCase struct {
	name        string
	inputs      []string // a1 / a2 / a3 / ...
	expectError bool
	result      string
}

func TestDecimal_Arith_Div(t *testing.T) {
	cases := []arithCase{
		{
			name:   "normal integer div",
			inputs: []string{"84", "3", "2"},
			result: "14",
		},
		{
			name:   "only one number",
			inputs: []string{"123456.78"},
			result: "123456.78",
		},
		{
			name:   "divisor is 0",
			inputs: []string{"0.00", "12"},
			result: "0",
		},
		{
			name:   "divisor less than dividents",
			inputs: []string{"20", "250"},
			result: "0.08",
		},
		{
			name:   "divisor equals to divients",
			inputs: []string{"1", "1"},
			result: "1",
		},
		{
			name:   "divide with rouding (to-nearest)",
			inputs: []string{"27", "7"},
			result: "3.857142857142857",
		},
		{
			name:   "divide with rouding #2",
			inputs: []string{"78e+1238", "56e+98"},
			result: "1.392857142857*10^1140",
		},
		{
			name:   "result > 100",
			inputs: []string{"2e+9", "4e+4"},
			result: "50000",
		},
		{
			name:   "divide towards zero",
			inputs: []string{"3e-8", "98831927397324"},
			result: "3.035456333801*10^-22",
		},
		{
			name:   "divide negative number",
			inputs: []string{"-25", "8"},
			result: "-3.125",
		},
		{
			name:   "divident is negative",
			inputs: []string{"80", "-2", "-4", "-5"},
			result: "-2",
		},
		{
			name:   "divisor and divident both negative",
			inputs: []string{"-25", "-48"},
			result: "0.5208333333333333",
		},
		{
			name:        "error: div 0",
			inputs:      []string{"2", "0"},
			expectError: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			decimals := []*Decimal{}
			for _, in := range tt.inputs {
				zd, _ := NewDecimal(in)
				decimals = append(decimals, zd)
			}

			result, err := decimals[0].Div(decimals[1:]...)
			if err != nil {
				if tt.expectError == false {
					t.Errorf("expect no error, got error: %s", err.Error())
				}
			} else {
				output := result.String()
				if tt.expectError == true {
					t.Errorf("expect error, got no error")
				} else if output != tt.result {
					t.Errorf("expect value: %s, got: %s", tt.result, output)
				}
			}
		})
	}
}

func TestDecimal_Arith_Mul(t *testing.T) {
	cases := []arithCase{
		{
			name:   "normal multiplication",
			inputs: []string{"2", "3"},
			result: "6",
		},
		{
			name:   "multiply negative number",
			inputs: []string{"2", "-3"},
			result: "-6",
		},
		{
			name:   "multiply multiple itmes",
			inputs: []string{"0.2", "2.5", "7"},
			result: "3.5",
		},
		{
			name:   "multiply 0",
			inputs: []string{"25", "0.00"},
			result: "0",
		},
		{
			name:   "negative numbers",
			inputs: []string{"-25*10^-8", "-6e-244"},
			result: "1.5*10^-250",
		},
		{
			name:   "only one number",
			inputs: []string{"123456.78"},
			result: "123456.78",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			decimals := []*Decimal{}
			for _, in := range tt.inputs {
				zd, _ := NewDecimal(in)
				decimals = append(decimals, zd)
			}

			result := decimals[0].Mul(decimals[1:]...)
			output := result.String()
			if tt.expectError == true {
				t.Errorf("expect error, got no error")
			} else if output != tt.result {
				t.Errorf("expect value: %s, got: %s", tt.result, output)
			}
		})
	}
}
