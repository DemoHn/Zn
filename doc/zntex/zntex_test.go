package zntex

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tt := &ZnTex{}
	data := []byte(
		`\begin{tag}{xx}{yy}{zz} Hello, Hello  t   
\command{DemoHn}{from}{mingchuanlu}
\command2[okr, okd, _oks]{diu,diu}{der,,
, der}
		

    YYYY
\end{tag}{DD}    
`)

	tt.ReadInput(data)
	err := tt.Parse()
	if err != nil {
		t.Error(err)
	}

	for _, tk := range tt.tokens {
		switch v := tk.(type) {
		case *CommandToken:
			fmt.Printf("CommandToken: literal(%s), options(%v) args(%v)\n", string(v.Literal), v.Options, v.Args)
		case *EnvironToken:
			beginFlag := "begin"
			if v.IsBegin == false {
				beginFlag = "end"
			}
			fmt.Printf("EnvToken: literal(%s) tag(%v) beginFlag(%v) options(%v) args(%v)\n", string(v.Literal), v.Tag, beginFlag, v.Options, v.Args)
		case *TextToken:
			fmt.Printf("TextToken: text(%s)\n", string(v.Text))
		case *CommentToken:
			fmt.Printf("CommentToken: text(%s)\n", string(v.Text))
		}

	}
}
