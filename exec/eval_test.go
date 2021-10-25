package exec

import (
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
	"github.com/DemoHn/Zn/lex"
)

type programOKSuite struct {
	name           string
	program        string
	symbols        map[string]ctx.Value
	expReturnValue ctx.Value
	expProbe       map[string][][]string
}

func Test_DuplicateValue(t *testing.T) {
	suites := []struct {
		name      string
		input     ctx.Value
		outputStr string
	}{
		{
			name:      "copy decimal",
			input:     val.NewDecimalFromInt(1217543, -9),
			outputStr: "0.001217543",
		},
		{
			name:      "copy decimal #2",
			input:     val.NewDecimalFromInt(12345, 4),
			outputStr: "123450000",
		},
		{
			name:      "copy string",
			input:     val.NewString("这是「一个」测试"),
			outputStr: "「这是「一个」测试」",
		},
		{
			name:      "copy bool",
			input:     val.NewBool(false),
			outputStr: "假",
		},
		{
			name: "copy array",
			input: val.NewArray([]ctx.Value{
				val.NewBool(true),
				val.NewString("哈哈哈哈"),
				val.NewDecimalFromInt(1234, -3),
			}),
			outputStr: "【真，「哈哈哈哈」，1.234】",
		},
		{
			name: "copy array (nested)",
			input: val.NewArray([]ctx.Value{
				val.NewArray([]ctx.Value{
					val.NewDecimalFromInt(123, 0),
					val.NewDecimalFromInt(1234, 0),
					val.NewDecimalFromInt(12345, 0),
				}),
				val.NewString("ASDF"),
			}),
			outputStr: "【【123，1234，12345】，「ASDF」】",
		},
		{
			name: "copy hashmap (nested)",
			input: val.NewHashMap([]val.KVPair{
				{
					Key:   "猪",
					Value: val.NewDecimalFromInt(100, 2),
				},
				{
					Key: "锅",
					Value: val.NewHashMap([]val.KVPair{
						{
							Key:   "SH",
							Value: val.NewBool(true),
						},
					}),
				},
			}),
			outputStr: "【猪 == 10000，锅 == 【SH == 真】】",
		},
	}

	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			out := val.DuplicateValue(suite.input)
			expectStr := val.StringifyValue(out)
			if expectStr != suite.outputStr {
				t.Errorf("val.DuplicateValue() result expect -> %s, got -> %s", suite.outputStr, expectStr)
			}
		})
	}
}

func Test_ExecPrimeExpr(t *testing.T) {
	suites := []programOKSuite{
		{
			name:           "simple string",
			program:        "「香港记者跑得快」",
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewString("香港记者跑得快"),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "simple decimal",
			program:        "314159*10^-8",
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(314159, -8),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple variable",
			program: "X-AE-A11",
			symbols: map[string]ctx.Value{
				"X-AE-A11": val.NewBool(true),
				"X-AE":     val.NewString("HelloWorld"),
			},
			expReturnValue: val.NewBool(true),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple array",
			program: "【10、20、300】",
			symbols: map[string]ctx.Value{
				"X-AE-A11": val.NewBool(true),
				"X-AE":     val.NewString("HelloWorld"),
			},
			expReturnValue: val.NewArray([]ctx.Value{
				val.NewDecimalFromInt(10, 0),
				val.NewDecimalFromInt(20, 0),
				val.NewDecimalFromInt(300, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name:           "empty array",
			program:        "【】",
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewArray([]ctx.Value{}),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple empty hashmap",
			program: "【=】",
			symbols: map[string]ctx.Value{
				"X-AE-A11": val.NewBool(true),
				"X-AE":     val.NewString("HelloWorld"),
			},
			expReturnValue: val.NewHashMap([]val.KVPair{}),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "simple hashmap",
			program: "【「1」 == 2】",
			symbols: map[string]ctx.Value{
				"X-AE-A11": val.NewBool(true),
				"X-AE":     val.NewString("HelloWorld"),
			},
			expReturnValue: val.NewHashMap([]val.KVPair{
				{
					Key:   "1",
					Value: val.NewDecimalFromInt(2, 0),
				},
			}),
			expProbe: map[string][][]string{},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func Test_MemberExpr(t *testing.T) {
	var exampleClassDef = `
定义示例：
	其名 为 “示例”
	其总和 为 0
	
	如何累加？
		已知 累加数
		其总和 为 【其总和、累加数】之和
		返回空

	如何获取总和？
		返回 其总和
`
	suites := []programOKSuite{
		{
			name:           "array index expr (normal)",
			program:        `【3、5、7、9】#2`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(7, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "array index expr (variable)",
			program: `A#1`,
			symbols: map[string]ctx.Value{
				"A": val.NewArray([]ctx.Value{
					val.NewString("XX"),
					val.NewString("YY"),
				}),
			},
			expReturnValue: val.NewString("YY"),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "consecutive array index expr",
			program:        `【【30、40】、100】#0 #0`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(30, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "hashmap index expr (normal)",
			program:        `【“L” == 7，“M” == 8】# “L”`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(7, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "hashmap index expr (variable)",
			program: `V # “L”`,
			symbols: map[string]ctx.Value{
				"V": val.NewHashMap([]val.KVPair{
					{Key: "L", Value: val.NewDecimalFromInt(7, 0)},
					{Key: "M", Value: val.NewDecimalFromInt(8, 0)},
				}),
			},
			expReturnValue: val.NewDecimalFromInt(7, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "consecutive hashmap index expr",
			program:        `【“X” ==【“Y” == 20】】# “X” # “Y”`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(20, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "expr as index",
			program:        `【3、5、7、9】#{（X-Y：8、7）}`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(5, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name: "object get property",
			program: exampleClassDef + `
令X 成为示例

X之名
			`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewString("示例"),
			expProbe:       map[string][][]string{},
		},
		/**
				{
					name: "object run methods",
					program: exampleClassDef + `
		令X 成为示例
		X之总和为20
		X之（累加：25）
		X之（获取总和）
					`,
					symbols:        map[string]ctx.Value{},
					expReturnValue: val.NewDecimalFromInt(45, 0),
					expProbe:       map[string][][]string{},
				},
		*/
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
	Y为（X+Y：Y、5）
	（__probe：「$X」、X）
	（__probe：「$Y」、Y）
			`,
			symbols: map[string]ctx.Value{
				"Y": val.NewDecimalFromInt(255, -1), // 25.5
				"诸变量": val.NewArray([]ctx.Value{
					val.NewString("一"),
					val.NewString("地"),
					val.NewString("在"),
					val.NewString("要"),
					val.NewString("工"),
				}),
			},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"$X": {
					{"100", "*val.Decimal"},
					{"100", "*val.Decimal"},
					{"100", "*val.Decimal"},
					{"100", "*val.Decimal"},
					{"100", "*val.Decimal"},
				},
				"$Y": {
					{"30.5", "*val.Decimal"},
					{"35.5", "*val.Decimal"},
					{"40.5", "*val.Decimal"},
					{"45.5", "*val.Decimal"},
					{"50.5", "*val.Decimal"},
				},
			},
		},
		{
			name: "with one var lead (array, hashmap)",
			program: `
以V遍历【30、 40、 50】：
    （__probe：「$L1V」、V）
    以V遍历【「甲」 == 20，「乙」 == 30】：
        （__probe：「$L2V」、V）`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"$L1V": {
					{"30", "*val.Decimal"},
					{"40", "*val.Decimal"},
					{"50", "*val.Decimal"},
				},
				"$L2V": {
					{"20", "*val.Decimal"},
					{"30", "*val.Decimal"},
					{"20", "*val.Decimal"},
					{"30", "*val.Decimal"},
					{"20", "*val.Decimal"},
					{"30", "*val.Decimal"},
				},
			},
		},
		{
			name: "with two vars lead (array)",
			program: `
以K、V遍历【「土」、「地」】：
    （__probe：「K1」、K）
    （__probe：「V1」、V）`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"K1": {
					{"0", "*val.Decimal"},
					{"1", "*val.Decimal"},
				},
				"V1": {
					{"土", "*val.String"},
					{"地", "*val.String"},
				},
			},
		},
		{
			name: "with two vars lead (hashmap)",
			program: `
以K、V遍历【「上」==「下」，「左」==「右」】：
    （__probe：「K1」、K）
    （__probe：「V1」、V）`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"K1": {
					{"上", "*val.String"},
					{"左", "*val.String"},
				},
				"V1": {
					{"下", "*val.String"},
					{"右", "*val.String"},
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
			program:        `令金克木为「森林」；（__probe：「$K1」、金克木）`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewString("森林"),
			expProbe: map[string][][]string{
				"$K1": {
					{"森林", "*val.String"},
				},
			},
		},
		{
			name:           "normal one var with compound expression",
			program:        `令_B52为（X+Y：2008、1963）；_B52`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(3971, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "normal multiple vars",
			program:        "令A为5；令B为2；令C为3；（X*Y：A、B、C）",
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(3, 1),
			expProbe:       map[string][][]string{},
		},
		{
			name:           "normal multiple vars (with reference)",
			program:        "令A为10；令B为A；令C为B；（X*Y：A、B、C）",
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(1, 3),
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
			symbols: map[string]ctx.Value{
				"A": val.NewBool(true),
			},
			expReturnValue: val.NewDecimalFromInt(200, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "normal var assign with computed value",
			program: `A为（X+Y：100、200）`,
			symbols: map[string]ctx.Value{
				"A": val.NewBool(true),
			},
			expReturnValue: val.NewDecimalFromInt(300, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name:    "normal var assign with array value",
			program: `A为【2、4、6、8】`,
			symbols: map[string]ctx.Value{
				"A": val.NewBool(true),
			},
			expReturnValue: val.NewArray([]ctx.Value{
				val.NewDecimalFromInt(2, 0),
				val.NewDecimalFromInt(4, 0),
				val.NewDecimalFromInt(6, 0),
				val.NewDecimalFromInt(8, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name:    "normal var assign with reference",
			program: `令B为【2、4、6、8】；A为&B；B#2为60；A`,
			symbols: map[string]ctx.Value{
				"A": val.NewBool(true),
			},
			// value of A should be same as value of B, since A is B's reference
			expReturnValue: val.NewArray([]ctx.Value{
				val.NewDecimalFromInt(2, 0),
				val.NewDecimalFromInt(4, 0),
				val.NewDecimalFromInt(60, 0),
				val.NewDecimalFromInt(8, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name:    "normal var assign without reference",
			program: `令B为【2、4、6、8】；A为B；B#2为60；A`,
			symbols: map[string]ctx.Value{
				"A": val.NewBool(true),
			},
			// value of A has been copied from value of B, so there's no changing effect
			// when B's value has been changed
			expReturnValue: val.NewArray([]ctx.Value{
				val.NewDecimalFromInt(2, 0),
				val.NewDecimalFromInt(4, 0),
				val.NewDecimalFromInt(6, 0),
				val.NewDecimalFromInt(8, 0),
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
（__probe：「B」、B之名）
（__probe：「C」、C之名）
A之名
`,
			symbols: map[string]ctx.Value{
				"B": val.NewBool(true),
				"C": val.NewBool(true),
			},
			// for objects, there's no difference between "assign by value" and "assign by reference"
			// which means all objects are transferred by reference. Thus when A's property changes,
			// B and C's properties also change.
			expReturnValue: val.NewString("保定"),
			expProbe: map[string][][]string{
				"B": {
					{"保定", "*val.String"},
				},
				"C": {
					{"保定", "*val.String"},
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
	（__probe：「$X」、X）
	X为（X-Y：X、1）`,
			symbols: map[string]ctx.Value{
				"X": val.NewDecimalFromInt(3, 0),
			},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"$X": {
					{"3", "*val.Decimal"},
					{"2", "*val.Decimal"},
					{"1", "*val.Decimal"},
				},
			},
		},
		{
			name: "test break",
			program: `
每当X大于0：
	Y为1
	每当Y大于0：
		Y为（X+Y：Y、1）
		如果Y = 4：
			（结束循环）
		（__probe：「VY」、Y）
		
	X为（X+Y：X、-1）
	（__probe：「VX」、X）
			`,
			symbols: map[string]ctx.Value{
				"X": val.NewDecimalFromInt(2, 0),
				"Y": val.NewDecimalFromInt(0, 0),
			},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"VY": {
					{"2", "*val.Decimal"},
					{"3", "*val.Decimal"},
					{"2", "*val.Decimal"},
					{"3", "*val.Decimal"},
				},
				"VX": {
					{"1", "*val.Decimal"},
					{"0", "*val.Decimal"},
				},
			},
		},
	}
	for _, tt := range suites {
		assertSuite(t, tt)
	}
}

func Test_BranchStmt(t *testing.T) {
	suites := []programOKSuite{
		{
			name: "exec true expr",
			program: `
如果 变量A = “真实”：
	（__probe：“TAG”、变量A）	
			`,
			symbols: map[string]ctx.Value{
				"变量A": val.NewString("真实"),
			},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"TAG": {
					{"真实", "*val.String"},
				},
			},
		},
		{
			name: "exec false expr",
			program: `
如果 变量A = “真实”：
	（__probe：“TAG”、 “走过真逻辑”）
（__probe：“TAG”、 “走过公共逻辑”）
			`,
			symbols: map[string]ctx.Value{
				"变量A": val.NewString("不真实"),
			},
			expReturnValue: val.NewString("走过公共逻辑"),
			expProbe: map[string][][]string{
				"TAG": {
					{"走过公共逻辑", "*val.String"},
				},
			},
		},
		{
			name: "if-else expr",
			program: `
如果 变量A 大于 100：
	（__probe：“TAG_A”、真）
否则：
	（__probe：“TAG_A”、假）

如果 变量B 大于 100：
	（__probe：“TAG_B”、真）
否则：
	（__probe：“TAG_B”、假）
			`,
			symbols: map[string]ctx.Value{
				"变量A": val.NewDecimalFromInt(120, 0), // true expression
				"变量B": val.NewDecimalFromInt(80, 0),  // false expression
			},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"TAG_A": {
					{"真", "*val.Bool"},
				},
				"TAG_B": {
					{"假", "*val.Bool"},
				},
			},
		},
		{
			name: "if-elseif expr",
			program: `
以成绩遍历【40、95、70、82】：
	如果 成绩 大于 90：
		评级 为 “优秀”
	再如 成绩 大于 80：
		评级 为 “良好”
	再如 成绩 大于 60：
		评级 为 “及格”
	否则：
		评级 为 “不及格”

	（__probe：“TAG”、 评级）
			`,
			symbols: map[string]ctx.Value{
				"评级": val.NewString("一般"),
			},
			expReturnValue: val.NewNull(),
			expProbe: map[string][][]string{
				"TAG": {
					{"不及格", "*val.String"},
					{"优秀", "*val.String"},
					{"及格", "*val.String"},
					{"良好", "*val.String"},
				},
			},
		},
		{
			name: "if-stmt: new scope",
			program: `
（__probe：“TAG”、评级）  注1：初始变量设置
如果成绩大于70：
	令评级为“优秀”
	（__probe：“TAG”、评级） 注2：在新作用域内定义变量并赋值

	成绩为85
	（__probe：“TAG”、成绩）	

（__probe：“TAG”、成绩） 注3：成绩 变量已经在全局作用域被修改，其值应为85
（__probe：“TAG”、评级）
			`,
			symbols: map[string]ctx.Value{
				"评级": val.NewString("一般"),
				"成绩": val.NewDecimalFromInt(80, 0),
			},
			expReturnValue: val.NewString("一般"),
			expProbe: map[string][][]string{
				"TAG": {
					{"一般", "*val.String"},
					{"优秀", "*val.String"},
					{"85", "*val.Decimal"},
					{"85", "*val.Decimal"},
					{"一般", "*val.String"},
				},
			},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func Test_FunctionCall(t *testing.T) {
	suites := []programOKSuite{
		{
			name: "normal function call (with return value)",
			program: `
如何执行方法？
	令A 为【20、30】
	A#1 为 40
	返回 A

（执行方法）
`,
			symbols: map[string]ctx.Value{},
			expReturnValue: val.NewArray([]ctx.Value{
				val.NewDecimalFromInt(20, 0),
				val.NewDecimalFromInt(40, 0),
			}),
			expProbe: map[string][][]string{},
		},
		{
			name: "normal function call (without return value)",
			program: `
如何执行方法？
	令A 为【20、30】
	如果 A#1 = 40：
		（X+Y：1、2）

（执行方法）
`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewNull(),
			expProbe:       map[string][][]string{},
		},
		{
			name: "function call with args",
			program: `
如何执行方法？
	已知A、B
	返回（X+Y：A、B）

（执行方法：10、30）
`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(40, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name: "function call (last expression as return value)",
			program: `
如何执行方法？
	已知A、B
	（X+Y：A、B）

（执行方法：10、30）
`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(40, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name: "function call hoisting (call expr exceeds definition)",
			program: `
令A为（执行方法：10、30）

如何执行方法？
	已知A、B
	（X+Y：A、B）

A`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(40, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name: "internal function declaration",
			program: `
如何执行方法？
	已知A、B

	如何加数据？
		已知A、B
		（X+Y：A、B）

	（加数据：A、B）

（执行方法：10、40）
`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(50, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name: "function call another function",
			program: `
如何执行方法？
	已知C、D
	（乘数据：（乘数据：C、D）、C）

如何乘数据？
	已知A、B
	（X*Y：A、B）

（执行方法：5、3）
			`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(75, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name: "function call new scope",
			program: `
如何执行方法？
	已知A、B
	令C为10
	D为20

	（X+Y：A、B、C）

（__probe：“RETURN”、（执行方法：3、4））
（__probe：“TAG_A”、A）
（__probe：“TAG_B”、B）
（__probe：“TAG_C”、C）
（__probe：“TAG_D”、D）
			`,
			symbols: map[string]ctx.Value{
				"A": val.NewDecimalFromInt(10, 0),
				"B": val.NewDecimalFromInt(20, 0),
				"C": val.NewDecimalFromInt(30, 0),
				"D": val.NewDecimalFromInt(40, 0),
			},
			expReturnValue: val.NewDecimalFromInt(20, 0),
			expProbe: map[string][][]string{
				"RETURN": {
					{"17", "*val.Decimal"},
				},
				"TAG_A": {
					{"10", "*val.Decimal"},
				},
				"TAG_B": {
					{"20", "*val.Decimal"},
				},
				"TAG_C": {
					{"30", "*val.Decimal"},
				},
				"TAG_D": {
					{"20", "*val.Decimal"},
				},
			},
		},
		{
			name: "function recursion call",
			program: `
如何调用FIB？
	已知X
	如果X不大于1：
		返回1
	否则：
		返回【（调用FIB：【X、-1】之和）、（调用FIB：【X、-2】之和）】之和


（调用FIB：10）`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(89, 0),
			expProbe:       map[string][][]string{},
		},
		{
			name: "function name alias",
			program: `
如何·乘以2·？
	已知X
	返回（X+Y：X、X）
	
令累加成双 为 ·乘以2· 
（累加成双：24）`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewDecimalFromInt(48, 0),
			expProbe:       map[string][][]string{},
		},
	}

	for _, suite := range suites {
		assertSuite(t, suite)
	}
}

func Test_CreateObject(t *testing.T) {
	suites := []programOKSuite{
		{
			name: "create object with empty constructor",
			program: `
定义模型：
	其名为 “乐高”

令X、Y 成为 模型

（__probe：“TAG”、X之名）
（__probe：“TAG”、Y之名）
			`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewString("乐高"),
			expProbe: map[string][][]string{
				"TAG": {
					{"乐高", "*val.String"},
					{"乐高", "*val.String"},
				},
			},
		},
		{
			name: "create object with one param",
			program: `
定义模型：
	其名为 “乐高”

	是为 名

令X、Y 成为 模型：“香港记者”

（__probe：“TAG”、X之名）
（__probe：“TAG”、Y之名）
			`,
			symbols:        map[string]ctx.Value{},
			expReturnValue: val.NewString("香港记者"),
			expProbe: map[string][][]string{
				"TAG": {
					{"香港记者", "*val.String"},
					{"香港记者", "*val.String"},
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
		c := ctx.NewContext(GlobalValues)
		in := lex.NewTextStream(suite.program)
		// impose symbols
		for k, v := range suite.symbols {
			c.BindSymbol(k, v)
		}

		result, err := ExecuteCode(c, in)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(c.GetScope().GetReturnValue(), suite.expReturnValue) {
			t.Errorf("return value expect -> %s, got -> %s", suite.expReturnValue, result)
			return
		}
		// assert probe value
		for tag, pLog := range suite.expProbe {
			gotLog := c.GetProbe().GetProbeLog(tag)
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
