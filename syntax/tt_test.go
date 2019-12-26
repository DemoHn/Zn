package syntax

/**
func TestRandomly(t *testing.T) {
	input := "令甲，乙为234"
	l := lex.NewLexer([]rune(input))

	parser := NewParser(l)

	n, e := parser.ParseVarAssign()

	if e != nil {
		t.Error(e)
		return
	}
	// print all data
	for _, item := range n.Variables {
		fmt.Println(*item)
	}
	fmt.Println(e)
}
*/
