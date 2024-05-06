package exec

import (
	"fmt"
	"testing"

	"github.com/DemoHn/Zn/pkg/syntax"
)

type tryMatchNumberCase struct {
	literal      string
	expectError  bool
	expectResult bool
}

func TestTryMatchNumber(t *testing.T) {
	cases := []tryMatchNumberCase{
		// #1. [OK] numbers
		{
			literal:      "045789",
			expectError:  false,
			expectResult: true,
		},
		{
			literal:      "-50",
			expectError:  false,
			expectResult: true,
		},
		{
			literal:      "-3.4e-8",
			expectError:  false,
			expectResult: true,
		},
		{
			literal:      "-3.4*10^8",
			expectError:  false,
			expectResult: true,
		},
		{
			literal:      "2.5*10^7",
			expectError:  false,
			expectResult: true,
		},
		// #2. [OK] identifier, but not number
		{
			literal:      "非常6+1",
			expectError:  false,
			expectResult: false,
		},
		{
			literal:      "标识符X",
			expectError:  false,
			expectResult: false,
		},
		{
			literal:      "Y3",
			expectError:  false,
			expectResult: false,
		},
		{
			literal:      "$15.87",
			expectError:  false,
			expectResult: false,
		},
		{
			literal:      "-Q3-",
			expectError:  false,
			expectResult: false,
		},
		{
			literal:      "-..",
			expectError:  false,
			expectResult: false,
		},
		// #3. [FAIL] more chars after number
		{
			literal:      "+2X",
			expectError:  true,
			expectResult: false,
		},
		// (we don't support unit type yet)
		{
			literal:      "25.8km/h",
			expectError:  true,
			expectResult: false,
		},
		{
			literal:      "9..",
			expectError:  true,
			expectResult: false,
		},
	}

	for idx, tt := range cases {
		t.Run(fmt.Sprintf("test TryMatchNumber#%d", idx+1), func(t *testing.T) {
			id := &syntax.ID{}
			id.SetLiteral([]rune(tt.literal))

			res, err := tryParseNumber(id)
			if (tt.expectError == false && err != nil) || (tt.expectError == true && err == nil) {
				t.Errorf("expect no error, got error: %s", err)
				return
			}

			if res != tt.expectResult {
				t.Errorf("expect result = %v, got %v", tt.expectResult, res)
			}
		})
	}
}
