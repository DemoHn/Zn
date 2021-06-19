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
