package zh

import (
	"fmt"
	"github.com/DemoHn/Zn/pkg/syntax"
	"testing"
)

func TestNum(t *testing.T) {
	b := TokenBuilderZH{}
	l := &syntax.Lexer{
		Source:     []rune("    12345ABC"),
		IndentType: 0,
		Lines:      nil,
	}

	tk, err := b.NextToken(l)
	fmt.Println(tk, err)
	fmt.Println(l.GetCurrentChar())
}
