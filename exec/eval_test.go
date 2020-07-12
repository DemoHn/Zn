package exec

import (
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/debug"
	"github.com/DemoHn/Zn/lex"
)

type programOKSuite struct {
	name           string
	program        string
	symbols        map[string]ZnValue
	expReturnValue ZnValue
	expProbe       map[string][]debug.ProbeLog
}

func Test_ExecPrimeExpr(t *testing.T) {
	suites := []programOKSuite{
		{
			name:           "simple string",
			program:        "「香港记者跑得快」",
			symbols:        map[string]ZnValue{},
			expReturnValue: NewZnString("香港记者跑得快"),
			expProbe:       map[string][]debug.ProbeLog{},
		},
		{
			name:           "simple decimal",
			program:        "314159*10^-8",
			symbols:        map[string]ZnValue{},
			expReturnValue: NewZnDecimalFromInt(314159, -8),
			expProbe:       map[string][]debug.ProbeLog{},
		},
		{
			name:    "simple variable",
			program: "X-AE-A11",
			symbols: map[string]ZnValue{
				"X-AE-A11": NewZnBool(true),
				"X-AE":     NewZnString("HelloWorld"),
			},
			expReturnValue: NewZnBool(true),
			expProbe:       map[string][]debug.ProbeLog{},
		},
		{
			name:    "simple array",
			program: "【10，20，300】",
			symbols: map[string]ZnValue{
				"X-AE-A11": NewZnBool(true),
				"X-AE":     NewZnString("HelloWorld"),
			},
			expReturnValue: NewZnArray([]ZnValue{
				NewZnDecimalFromInt(10, 0),
				NewZnDecimalFromInt(20, 0),
				NewZnDecimalFromInt(300, 0),
			}),
			expProbe: map[string][]debug.ProbeLog{},
		},
		{
			name:    "simple empty hashmap",
			program: "【==】",
			symbols: map[string]ZnValue{
				"X-AE-A11": NewZnBool(true),
				"X-AE":     NewZnString("HelloWorld"),
			},
			expReturnValue: NewZnHashMap([]KVPair{}),
			expProbe:       map[string][]debug.ProbeLog{},
		},
		{
			name:    "simple hashmap",
			program: "【「1」 == 2】",
			symbols: map[string]ZnValue{
				"X-AE-A11": NewZnBool(true),
				"X-AE":     NewZnString("HelloWorld"),
			},
			expReturnValue: NewZnHashMap([]KVPair{
				{
					Key:   "1",
					Value: NewZnDecimalFromInt(2, 0),
				},
			}),
			expProbe: map[string][]debug.ProbeLog{},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func assertSuite(t *testing.T, suite programOKSuite) {
	t.Run(suite.name, func(t *testing.T) {
		ctx := NewContext()
		scope := NewRootScope()
		// impose symbols
		for k, v := range suite.symbols {
			scope.SetSymbol(k, v, false)
		}

		in := lex.NewTextStream(suite.program)
		result := ctx.ExecuteCode(in, scope)

		// assert result
		if result.HasError {
			t.Errorf("program should have no error, got error: %s", result.Error)
			return
		}
		if !reflect.DeepEqual(result.Value, suite.expReturnValue) {
			t.Errorf("return value expect -> %s, got -> %s", suite.expReturnValue, result.Value)
			return
		}
		// assert probe value
		for tag, pLog := range suite.expProbe {
			gotLog := ctx._probe.GetProbeLog(tag)
			// ensure length is same
			if len(gotLog) != len(pLog) {
				t.Errorf("probe log length not match, expect -> %d, got -> %d", len(pLog), len(gotLog))
				return
			}
			// then check item one by one
			for idx, pLogItem := range pLog {
				if !reflect.DeepEqual(pLogItem.ValueStr, gotLog[idx].ValueStr) {
					t.Errorf("probe log `valueStr` not match at #%d, expect -> %s, got -> %s", idx, pLogItem.ValueStr, gotLog[idx].ValueStr)
					return
				}
				if !reflect.DeepEqual(pLogItem.ValueType, gotLog[idx].ValueType) {
					t.Errorf("probe log `valueType` not match at #%d, expect -> %s, got -> %s", idx, pLogItem.ValueType, gotLog[idx].ValueType)
					return
				}
			}
		}
	})
}
