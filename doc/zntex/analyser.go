package zntex

// Node - 定义节点
type Node interface {
	getType() int
}

// declare some types
const (
	TypeTextNode    = 20
	TypeEnvironNode = 21
	TypeCommandNode = 22
)

// TextNode - 文本节点
type TextNode struct {
	text string
}

func (t *TextNode) getType() int {
	return TypeTextNode
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
	return TypeEnvironNode
}

// CommandNode - 命令节点
type CommandNode struct {
	command string
	options []string
	args    [][]Node
}

func (t *CommandNode) getType() int {
	return TypeCommandNode
}

// Analyse - 分析得到其AST
func Analyse(tokens []Token) []Node {
	var ast = make([]Node, 0)

	/**for _, token := range tokens {
		switch v := token.(type) {
		//case *CommandNode
		}
	}*/
	return ast
}
