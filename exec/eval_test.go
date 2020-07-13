package exec

import (
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/lex"
)

type programOKSuite struct {
	name           string
	program        string
	symbols        map[string]ZnValue
	expReturnValue ZnValue
	expProbe       map[string][][]string
}

func Test_ExecPrimeExpr(t *testing.T) {
	suites := []programOKSuite{
		{
			name:           "simple string",
			program:        "「香港记者跑得快」",
			symbols:        map[string]ZnValue{},
			expReturnValue: NewZnString("香港记者跑得快"),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "simple decimal",
			program:        "314159*10^-8",
			symbols:        map[string]ZnValue{},
			expReturnValue: NewZnDecimalFromInt(314159, -8),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple variable",
			program: "X-AE-A11",
			symbols: map[string]ZnValue{
				"X-AE-A11": NewZnBool(true),
				"X-AE":     NewZnString("HelloWorld"),
			},
			expReturnValue: NewZnBool(true),
			expProbe:       map[string][][]string{},
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
			expProbe: map[string][][]string{},
		},
		{
			name:    "simple empty hashmap",
			program: "【==】",
			symbols: map[string]ZnValue{
				"X-AE-A11": NewZnBool(true),
				"X-AE":     NewZnString("HelloWorld"),
			},
			expReturnValue: NewZnHashMap([]KVPair{}),
			expProbe:       map[string][][]string{},
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
			expProbe: map[string][][]string{},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func Test_IterateStmt(t *testing.T) {
	suites := []programOKSuite{
		{
			name: "with no lead variables (array)",
			program: `
遍历诸变量：
	令X为100
	Y为（X+Y：Y，5）
	（__probe：「$KEY」，此之索引）
	（__probe：「$VAL」，此之值）
	（__probe：「$X」，X）
	（__probe：「$Y」，Y）
			`,
			symbols: map[string]ZnValue{
				"Y": NewZnDecimalFromInt(255, -1), // 25.5
				"诸变量": NewZnArray([]ZnValue{
					NewZnString("一"),
					NewZnString("地"),
					NewZnString("在"),
					NewZnString("要"),
					NewZnString("工"),
				}),
			},
			expReturnValue: NewZnNull(),
			expProbe: map[string][][]string{
				"$KEY": {
					{"0", "*exec.ZnDecimal"},
					{"1", "*exec.ZnDecimal"},
					{"2", "*exec.ZnDecimal"},
					{"3", "*exec.ZnDecimal"},
					{"4", "*exec.ZnDecimal"},
				},
				"$VAL": {
					{"「一」", "*exec.ZnString"},
					{"「地」", "*exec.ZnString"},
					{"「在」", "*exec.ZnString"},
					{"「要」", "*exec.ZnString"},
					{"「工」", "*exec.ZnString"},
				},
				"$X": {
					{"100", "*exec.ZnDecimal"},
					{"100", "*exec.ZnDecimal"},
					{"100", "*exec.ZnDecimal"},
					{"100", "*exec.ZnDecimal"},
					{"100", "*exec.ZnDecimal"},
				},
				"$Y": {
					{"30.5", "*exec.ZnDecimal"},
					{"35.5", "*exec.ZnDecimal"},
					{"40.5", "*exec.ZnDecimal"},
					{"45.5", "*exec.ZnDecimal"},
					{"50.5", "*exec.ZnDecimal"},
				},
			},
		},
		{
			name: "with no lead variables (hashmap)",
			program: `
遍历示例列表：
	（__probe：「$KEY」，此之索引）
	（__probe：「$VAL」，此之值）
			`,
			symbols: map[string]ZnValue{
				"示例列表": NewZnHashMap([]KVPair{
					{
						Key:   "积分",
						Value: NewZnDecimalFromInt(1000, 0),
					},
					{
						Key:   "年龄",
						Value: NewZnDecimalFromInt(24, 0),
					},
					{
						Key:   "穿着",
						Value: NewZnString("蕾丝边裙子"),
					},
				}),
			},
			expReturnValue: NewZnNull(),
			expProbe: map[string][][]string{
				"$KEY": {
					{"「积分」", "*exec.ZnString"},
					{"「年龄」", "*exec.ZnString"},
					{"「穿着」", "*exec.ZnString"},
				},
				"$VAL": {
					{"1000", "*exec.ZnDecimal"},
					{"24", "*exec.ZnDecimal"},
					{"「蕾丝边裙子」", "*exec.ZnString"},
				},
			},
		},
		{
			name: "with one var lead (array, hashmap)",
			program: `
以V遍历【30， 40， 50】：
    （__probe：「$L1V」，V）
    以V遍历【「甲」 == 20，「乙」 == 30】：
        （__probe：「$L2V」，V）`,
			symbols:        map[string]ZnValue{},
			expReturnValue: NewZnNull(),
			expProbe: map[string][][]string{
				"$L1V": {
					{"30", "*exec.ZnDecimal"},
					{"40", "*exec.ZnDecimal"},
					{"50", "*exec.ZnDecimal"},
				},
				"$L2V": {
					{"20", "*exec.ZnDecimal"},
					{"30", "*exec.ZnDecimal"},
					{"20", "*exec.ZnDecimal"},
					{"30", "*exec.ZnDecimal"},
					{"20", "*exec.ZnDecimal"},
					{"30", "*exec.ZnDecimal"},
				},
			},
		},
		{
			name: "with two vars lead (array)",
			program: `
以K，V遍历【「土」，「地」】：
    （__probe：「K1」，K）
    （__probe：「V1」，V）`,
			symbols:        map[string]ZnValue{},
			expReturnValue: NewZnNull(),
			expProbe: map[string][][]string{
				"K1": {
					{"0", "*exec.ZnDecimal"},
					{"1", "*exec.ZnDecimal"},
				},
				"V1": {
					{"「土」", "*exec.ZnString"},
					{"「地」", "*exec.ZnString"},
				},
			},
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
				t.Errorf("probe log [%s] length not match, expect -> %d, got -> %d", tag, len(pLog), len(gotLog))
				return
			}
			// then check item one by one
			for idx, pLogItem := range pLog {
				if !reflect.DeepEqual(pLogItem[0], gotLog[idx].ValueStr) {
					t.Errorf("probe log [%s] `valueStr` not match at #%d, expect -> %s, got -> %s", tag, idx, pLogItem[0], gotLog[idx].ValueStr)
					return
				}
				if !reflect.DeepEqual(pLogItem[1], gotLog[idx].ValueType) {
					t.Errorf("probe log [%s] `valueType` not match at #%d, expect -> %s, got -> %s", tag, idx, pLogItem[1], gotLog[idx].ValueType)
					return
				}
			}
		}
	})
}
