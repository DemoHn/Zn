package zntex

// Node - 定义节点
type Node interface {
	getType() int
}

// declare some types
const (
	typeTextNode    = 20
	typeEnvironNode = 21
	typeCommandNode = 22
)

// TextNode - 文本节点
type TextNode struct {
	text string
}

func (t *TextNode) getType() int {
	return typeTextNode
}

func (t *TextNode) getText() string {
	return t.text
}

// EnvironNode - env节点
type EnvironNode struct {
	tag      string
	options  []string
	args     [][]Node
	children []Node
}

func (t *EnvironNode) getType() int {
	return typeEnvironNode
}

// CommandNode - 命令节点
type CommandNode struct {
	command string
	options []string
	args    [][]Node
}

func (t *CommandNode) getType() int {
	return typeCommandNode
}

// Analyse - 分析得到其AST
func Analyse(tokens []Token) []Node {
	var ast = make([]Node, 0)

	for _, token := range tokens {
		switch v := token.(type) {
			case *CommandNode
		}
	}
	return ast
}
