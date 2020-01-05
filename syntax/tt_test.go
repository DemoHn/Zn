package syntax

import (
	"fmt"
	"testing"

	"github.com/DemoHn/Zn/lex"
)

func TestRandomly(t *testing.T) {
	input := "令甲，乙为（【12，34，【“测试到底”，10】】）\n  令丙为“23”；a 自 设为12；【12】，得到利益"
	l := lex.NewLexer([]rune(input))

	parser := NewParser(l)
	ast, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}

	data := StringifyAST(ast)
	fmt.Println(data)
}
