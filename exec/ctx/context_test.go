package ctx

import (
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/exec/val"
	"github.com/DemoHn/Zn/lex"
)

func TestExecuteCode_OK(t *testing.T) {
	cases := []struct {
		name        string
		text        string
		resultValue Value
	}{
		{
			"normal oneline expression",
			"令A为10；A为10241024",
			newDecimal("10241024"),
		},
		{
			"function call",
			"如何测试？\n    （X+Y：2，3）\n\n（测试）",
			newDecimal("5"),
		},
		{
			"with return",
			`如何测试？
	已知阈值
	如果阈值大于10：
		返回「大于」
	返回「小于」
	「等于」  注：这是一个干扰项
	
（测试：6）`,
			NewString("小于"),
		},
	}
	ctx := NewContext()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			in := lex.NewTextStream(tt.text)
			res, err := ctx.ExecuteCode(in)
			if err != nil {
				t.Errorf("expect no error, has got error: %v", err)
				return
			}

			if !reflect.DeepEqual(tt.resultValue, res) {
				t.Errorf("expect value: %v, got: %v", tt.resultValue, res)
				return
			}

			ctx.resetScopeValue()
		})
	}
}

// display full error info
func TestExecuteCode_FAIL(t *testing.T) {
	text := `令变量名-甲为10
令变量名-乙为20
（X+Y：变量名-未定，变量名-甲）`

	in := lex.NewTextStream(text)
	ctx := NewContext()
	_, err := ctx.ExecuteCode(in)

	if err == nil {
		t.Errorf("should got error, return no error")
		return
	}

	displayText := `在「$repl」中，位于第 3 行发现错误：
    （X+Y：变量名-未定，变量名-甲）
    
‹2501› 标识错误：标识「变量名-未定」未有定义`

	if err.Display() != displayText {
		t.Errorf("should return \n%s\n, got \n%s\n", displayText, err.Display())
		return
	}
}

// create decimal (and ignore errors)
func newDecimal(value string) *Decimal {
	dat, _ := val.NewDecimal(value)
	return dat
}
