package zh

import (
	"fmt"
	"github.com/DemoHn/Zn/pkg/syntax"
	"testing"
)

func TestNum(t *testing.T) {
	b := TokenBuilderZH{}
	l := syntax.NewLexer([]rune("\n    \n\r        2333"))
	tk, err := b.NextToken(l)
	fmt.Println(tk, err)
	fmt.Println(l.Lines)
	fmt.Println(l.GetCurrentChar())
}
