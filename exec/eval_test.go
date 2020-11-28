package exec

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

type programOKSuite struct {
	name           string
	program        string
	symbols        map[string]Value
	expReturnValue Value
	expProbe       map[string][][]string
}

func Test_DuplicateValue(t *testing.T) {
	suites := []struct {
		name      string
		input     Value
		outputStr string
	}{
		{
			name:      "copy decimal",
			input:     NewDecimalFromInt(1217543, -9),
			outputStr: "0.001217543",
		},
		{
			name: "copy decimal #2",
			input: &Decimal{
				co:  big.NewInt(12345),
				exp: 4,
			},
			outputStr: "123450000",
		},
		{
			name:      "copy string",
			input:     &String{value: "这是「一个」测试"},
			outputStr: "「这是「一个」测试」",
		},
		{
			name:      "copy bool",
			input:     &Bool{value: false},
			outputStr: "假",
		},
		{
			name: "copy array",
			input: &Array{
				value: []Value{&Bool{value: true}, &String{value: "哈哈哈哈"}, NewDecimalFromInt(1234, -3)},
			},
			outputStr: "【真，「哈哈哈哈」，1.234】",
		},
		{
			name: "copy array (nested)",
			input: NewArray([]Value{
				NewArray([]Value{
					NewDecimalFromInt(123, 0),
					NewDecimalFromInt(1234, 0),
					NewDecimalFromInt(12345, 0),
				}),
				NewString("ASDF"),
			}),
			outputStr: "【【123，1234，12345】，「ASDF」】",
		},
		{
			name: "copy hashmap (nested)",
			input: NewHashMap([]KVPair{
				{
					Key:   "猪",
					Value: NewDecimalFromInt(100, 2),
				},
				{
					Key: "锅",
					Value: NewHashMap([]KVPair{
						{
							Key:   "SH",
							Value: NewBool(true),
						},
					}),
				},
			}),
			outputStr: "【猪 == 10000，锅 == 【SH == 真】】",
		},
	}

	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			out := duplicateValue(suite.input)
			expectStr := StringifyValue(out)
			if expectStr != suite.outputStr {
				t.Errorf("duplicateValue() result expect -> %s, got -> %s", suite.outputStr, expectStr)
			}
		})
	}
}

func Test_ExecPrimeExpr(t *testing.T) {
	suites := []programOKSuite{
		{
			name:           "simple string",
			program:        "「香港记者跑得快」",
			symbols:        map[string]Value{},
			expReturnValue: NewString("香港记者跑得快"),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "simple decimal",
			program:        "314159*10^-8",
			symbols:        map[string]Value{},
			expReturnValue: NewDecimalFromInt(314159, -8),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple variable",
			program: "X-AE-A11",
			symbols: map[string]Value{
				"X-AE-A11": NewBool(true),
				"X-AE":     NewString("HelloWorld"),
			},
			expReturnValue: NewBool(true),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple array",
			program: "【10，20，300】",
			symbols: map[string]Value{
				"X-AE-A11": NewBool(true),
				"X-AE":     NewString("HelloWorld"),
			},
			expReturnValue: NewArray([]Value{
				NewDecimalFromInt(10, 0),
				NewDecimalFromInt(20, 0),
				NewDecimalFromInt(300, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name:    "simple empty hashmap",
			program: "【==】",
			symbols: map[string]Value{
				"X-AE-A11": NewBool(true),
				"X-AE":     NewString("HelloWorld"),
			},
			expReturnValue: NewHashMap([]KVPair{}),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple hashmap",
			program: "【「1」 == 2】",
			symbols: map[string]Value{
				"X-AE-A11": NewBool(true),
				"X-AE":     NewString("HelloWorld"),
			},
			expReturnValue: NewHashMap([]KVPair{
				{
					Key:   "1",
					Value: NewDecimalFromInt(2, 0),
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
			symbols: map[string]Value{
				"Y": NewDecimalFromInt(255, -1), // 25.5
				"诸变量": NewArray([]Value{
					NewString("一"),
					NewString("地"),
					NewString("在"),
					NewString("要"),
					NewString("工"),
				}),
			},
			expReturnValue: NewNull(),
			expProbe: map[string][][]string{
				"$KEY": {
					{"0", "*exec.Decimal"},
					{"1", "*exec.Decimal"},
					{"2", "*exec.Decimal"},
					{"3", "*exec.Decimal"},
					{"4", "*exec.Decimal"},
				},
				"$VAL": {
					{"一", "*exec.String"},
					{"地", "*exec.String"},
					{"在", "*exec.String"},
					{"要", "*exec.String"},
					{"工", "*exec.String"},
				},
				"$X": {
					{"100", "*exec.Decimal"},
					{"100", "*exec.Decimal"},
					{"100", "*exec.Decimal"},
					{"100", "*exec.Decimal"},
					{"100", "*exec.Decimal"},
				},
				"$Y": {
					{"30.5", "*exec.Decimal"},
					{"35.5", "*exec.Decimal"},
					{"40.5", "*exec.Decimal"},
					{"45.5", "*exec.Decimal"},
					{"50.5", "*exec.Decimal"},
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
			symbols: map[string]Value{
				"示例列表": NewHashMap([]KVPair{
					{
						Key:   "积分",
						Value: NewDecimalFromInt(1000, 0),
					},
					{
						Key:   "年龄",
						Value: NewDecimalFromInt(24, 0),
					},
					{
						Key:   "穿着",
						Value: NewString("蕾丝边裙子"),
					},
				}),
			},
			expReturnValue: NewNull(),
			expProbe: map[string][][]string{
				"$KEY": {
					{"积分", "*exec.String"},
					{"年龄", "*exec.String"},
					{"穿着", "*exec.String"},
				},
				"$VAL": {
					{"1000", "*exec.Decimal"},
					{"24", "*exec.Decimal"},
					{"蕾丝边裙子", "*exec.String"},
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
			symbols:        map[string]Value{},
			expReturnValue: NewNull(),
			expProbe: map[string][][]string{
				"$L1V": {
					{"30", "*exec.Decimal"},
					{"40", "*exec.Decimal"},
					{"50", "*exec.Decimal"},
				},
				"$L2V": {
					{"20", "*exec.Decimal"},
					{"30", "*exec.Decimal"},
					{"20", "*exec.Decimal"},
					{"30", "*exec.Decimal"},
					{"20", "*exec.Decimal"},
					{"30", "*exec.Decimal"},
				},
			},
		},
		{
			name: "with two vars lead (array)",
			program: `
以K，V遍历【「土」，「地」】：
    （__probe：「K1」，K）
    （__probe：「V1」，V）`,
			symbols:        map[string]Value{},
			expReturnValue: NewNull(),
			expProbe: map[string][][]string{
				"K1": {
					{"0", "*exec.Decimal"},
					{"1", "*exec.Decimal"},
				},
				"V1": {
					{"土", "*exec.String"},
					{"地", "*exec.String"},
				},
			},
		},
		{
			name: "with two vars lead (hashmap)",
			program: `
以K，V遍历【「上」==「下」，「左」==「右」】：
    （__probe：「K1」，K）
    （__probe：「V1」，V）`,
			symbols:        map[string]Value{},
			expReturnValue: NewNull(),
			expProbe: map[string][][]string{
				"K1": {
					{"上", "*exec.String"},
					{"左", "*exec.String"},
				},
				"V1": {
					{"下", "*exec.String"},
					{"右", "*exec.String"},
				},
			},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func Test_VarDeclareStmt(t *testing.T) {
	suites := []programOKSuite{
		{
			name:           "normal one var declaration",
			program:        `令金克木为「森林」；（__probe：「$K1」，金克木）`,
			symbols:        map[string]Value{},
			expReturnValue: NewString("森林"),
			expProbe: map[string][][]string{
				"$K1": {
					{"森林", "*exec.String"},
				},
			},
		},
		{
			name:           "normal one var with compound expression",
			program:        `令_B52为（X+Y：2008，1963）；_B52`,
			symbols:        map[string]Value{},
			expReturnValue: NewDecimalFromInt(3971, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "normal multiple vars",
			program:        "令A为5；令B为2；令C为3；（X*Y：A，B，C）",
			symbols:        map[string]Value{},
			expReturnValue: NewDecimalFromInt(30, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "normal multiple vars (with reference)",
			program:        "令A为10；令B为A；令C为B；（X*Y：A，B，C）",
			symbols:        map[string]Value{},
			expReturnValue: NewDecimalFromInt(1000, 0),
			expProbe:       map[string][][]string{},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func Test_VarAssignExpr(t *testing.T) {
	suites := []programOKSuite{
		{
			name:    "normal var assign",
			program: `A为200`,
			symbols: map[string]Value{
				"A": NewBool(true),
			},
			expReturnValue: NewDecimalFromInt(200, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "normal var assign with computed value",
			program: `A为（X+Y：100，200）`,
			symbols: map[string]Value{
				"A": NewBool(true),
			},
			expReturnValue: NewDecimalFromInt(300, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "normal var assign with array value",
			program: `A为【2，4，6，8】`,
			symbols: map[string]Value{
				"A": NewBool(true),
			},
			expReturnValue: NewArray([]Value{
				NewDecimalFromInt(2, 0),
				NewDecimalFromInt(4, 0),
				NewDecimalFromInt(6, 0),
				NewDecimalFromInt(8, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name:    "normal var assign with reference",
			program: `令B为【2，4，6，8】；A为&B；B#2为60；A`,
			symbols: map[string]Value{
				"A": NewBool(true),
			},
			// value of A should be same as value of B, since A is B's reference
			expReturnValue: NewArray([]Value{
				NewDecimalFromInt(2, 0),
				NewDecimalFromInt(4, 0),
				NewDecimalFromInt(60, 0),
				NewDecimalFromInt(8, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name:    "normal var assign without reference",
			program: `令B为【2，4，6，8】；A为B；B#2为60；A`,
			symbols: map[string]Value{
				"A": NewBool(true),
			},
			// value of A has been copied from value of B, so there's no changing effect
			// when B's value has been changed
			expReturnValue: NewArray([]Value{
				NewDecimalFromInt(2, 0),
				NewDecimalFromInt(4, 0),
				NewDecimalFromInt(6, 0),
				NewDecimalFromInt(8, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name: "var assign object with/without reference",
			program: `
定义城市：
	其名为「正定」
	是为名

令A成为城市：「正定」
B为A
C为&A

A之名为「保定」

注： 显示结果，「B之名」 和 「C之名」 应都为 「保定」
（__probe：「B」，B之名）
（__probe：「C」，C之名）
A之名
`,
			symbols: map[string]Value{
				"B": NewBool(true),
				"C": NewBool(true),
			},
			// for objects, there's no difference between "assign by value" and "assign by reference"
			// which means all objects are transferred by reference. Thus when A's property changes,
			// B and C's properties also change.
			expReturnValue: NewString("保定"),
			expProbe: map[string][][]string{
				"B": {
					{"保定", "*exec.String"},
				},
				"C": {
					{"保定", "*exec.String"},
				},
			},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func Test_WhileLoopStmt(t *testing.T) {
	suites := []programOKSuite{
		{
			name: "simple while loop",
			program: `
每当X大于0：
	（__probe：「$X」，X）
	X为（X-Y：X，1）`,
			symbols: map[string]Value{
				"X": NewDecimalFromInt(3, 0),
			},
			expReturnValue: NewNull(),
			expProbe: map[string][][]string{
				"$X": {
					{"3", "*exec.Decimal"},
					{"2", "*exec.Decimal"},
					{"1", "*exec.Decimal"},
				},
			},
		},
		{
			name: "test break",
			program: `
每当X大于0：
	Y为1
	每当Y大于0：
		Y为（X+Y：Y，1）
		如果Y为4：
			此之（结束）
		（__probe：「VY」，Y）
		
	X为（X+Y：X，-1）
	（__probe：「VX」，X）
			`,
			symbols: map[string]Value{
				"X": NewDecimalFromInt(2, 0),
				"Y": NewDecimalFromInt(0, 0),
			},
			expReturnValue: NewNull(),
			expProbe: map[string][][]string{
				"VY": {
					{"2", "*exec.Decimal"},
					{"3", "*exec.Decimal"},
					{"2", "*exec.Decimal"},
					{"3", "*exec.Decimal"},
				},
				"VX": {
					{"1", "*exec.Decimal"},
					{"0", "*exec.Decimal"},
				},
			},
		},
	}
	for _, tt := range suites {
		assertSuite(t, tt)
	}
}

func assertSuite(t *testing.T, suite programOKSuite) {
	t.Run(suite.name, func(t *testing.T) {
		var e2 *error.Error
		ctx := NewContext()
		in := lex.NewTextStream(suite.program)
		// parseCode
		program, err := ctx.parseCode(in)
		if err != nil {
			e2 = err
		}
		// init scope
		ctx.initScope(program.Lexer)

		// impose symbols
		for k, v := range suite.symbols {
			ctx.scope.symbolMap[k] = SymbolInfo{
				value:   v,
				isConst: false,
			}
		}

		result, err := ctx.execProgram(program)
		if err != nil {
			e2 = err
		}

		// assert result
		if e2 != nil {
			t.Errorf("program should have no error, got error: %s", e2.Display())
			return
		}
		if !reflect.DeepEqual(ctx.scope.returnValue, suite.expReturnValue) {
			t.Errorf("return value expect -> %s, got -> %s", suite.expReturnValue, result)
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
