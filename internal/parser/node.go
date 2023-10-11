package parser

type Node struct {
	Kind     string
	Value    string
	Children []Node
}

func (n Node) IsParentheses() bool {
	return n.Kind == KindParentheses
}

func (n Node) IsVar() bool {
	return n.Kind == KindVar
}

func (n Node) IsNumber() bool {
	return n.Kind == KindInteger || n.Kind == KindFloat
}

func (n Node) IsOperation() bool {
	switch n.Kind {
	case KindOpEqual:
		fallthrough
	case KindOpPlus:
		fallthrough
	case KindOpMinus:
		fallthrough
	case KindOpDivide:
		fallthrough
	case KindOpMultiply:
		return true
	default:
		return false
	}
}

func (n Node) needSecondOperand() bool {
	return n.IsParentheses() && len(n.Children) == 2 && (n.Children[1].Kind == KindOpMultiply || n.Children[1].Kind == KindOpDivide)
}

type Tree []Node

func (t Tree) Pop() Tree {
	if len(t) == 0 {
		return t
	}
	return t[:len(t)-1]
}

func (t Tree) Last() (*Node, bool) {
	if len(t) == 0 {
		return nil, false
	}
	return &t[len(t)-1], true
}

func (t Tree) expectsNextNode() bool {
	node, ok := t.Last()
	if ok && !node.IsOperation() && !node.needSecondOperand() {
		return false
	}
	return true
}
