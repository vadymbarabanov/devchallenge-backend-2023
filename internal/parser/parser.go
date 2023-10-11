package parser

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidParentheses = errors.New("invalid parentheses")
	ErrInvalidOperation   = errors.New("invalid operation")
)

// Parse parses given input string into an abstract syntax tree.
// If the input is not a valid Excel formula an error will be returned.
func Parse(input string) (Tree, error) {
	input = strings.Trim(input, " ")
	nodes := make(Tree, 0)
	buffer := make([]rune, 0)
	parenStack := make([]rune, 0)
	parenBuffer := make([]rune, 0)

	for i, char := range input {
		if char == Space {
			continue
		}

		// catches all chars between '(' and ')'
		// including others '(' and ')'
		if len(parenStack) > 0 {
			if char == OpenParen {
				parenBuffer = append(parenBuffer, char)
				parenStack = append(parenStack, char)
				continue
			}

			if char == CloseParen && len(parenStack) > 1 {
				parenBuffer = append(parenBuffer, char)
				parenStack = parenStack[:len(parenStack)-1]
				continue
			}

			if char != CloseParen {
				parenBuffer = append(parenBuffer, char)
				continue
			}
		}

		isLastChar := len(input)-1 == i

		// continue fill variable name or number if already started
		if len(buffer) > 0 {
			switch {
			case isLetter(char) || unicode.IsNumber(char) && !isLastChar:
				buffer = append(buffer, char)
			case char == Dot && unicode.IsNumber(buffer[0]) && !isLastChar:
				buffer = append(buffer, char)
			case (isLetter(char) || unicode.IsNumber(char)) && isLastChar:
				buffer = append(buffer, char)
				if !nodes.expectsNextNode() {
					return nil, ErrInvalidOperation
				}

				node := createVarOrNumberNode(buffer)
				nodes = satisfyOperators(nodes, node)
				return nodes, nil
			default:
				if !nodes.expectsNextNode() {
					return nil, ErrInvalidOperation
				}
				node := createVarOrNumberNode(buffer)
				nodes = satisfyOperators(nodes, node)
				buffer = make([]rune, 0)
			}
		}

		node := Node{}

		switch char {
		case OpEqual:
			node.Kind = KindOpEqual
			nodes = append(nodes, node)

		case OpPlus:
			if isLastChar {
				return nil, ErrInvalidOperation
			}

			node.Kind = KindOpPlus
			nodes = append(nodes, node)
		case OpMinus:
			if isLastChar {
				return nil, ErrInvalidOperation
			}

			lastNode, ok := nodes.Last()
			if !ok || (!lastNode.IsNumber() && !lastNode.IsVar()) {
				node.Kind = KindParentheses
				node.Children = []Node{
					{
						Kind:  KindInteger,
						Value: "-1",
					},
					{
						Kind: KindOpMultiply,
					},
				}
				nodes = append(nodes, node)
			} else {
				node.Kind = KindOpMinus
				nodes = append(nodes, node)
			}

		case OpDivide:
			if isLastChar {
				return nil, ErrInvalidOperation
			}

			node.Kind = KindOpDivide
			tree, err := wrapLastNode(nodes, node)
			if err != nil {
				return nil, err
			}
			nodes = tree

		case OpMultiply:
			if isLastChar {
				return nil, ErrInvalidOperation
			}

			node.Kind = KindOpMultiply
			tree, err := wrapLastNode(nodes, node)
			if err != nil {
				return nil, err
			}
			nodes = tree

		// catches first opened parenthesis
		case OpenParen:
			parenStack = append(parenStack, char)
		// catches last closed parenthesis
		case CloseParen:
			if !nodes.expectsNextNode() {
				return nil, ErrInvalidOperation
			}

			if len(parenStack) != 1 || parenStack[0] != OpenParen {
				return nil, ErrInvalidParentheses
			}

			node.Kind = KindParentheses
			parenStack = make([]rune, 0)

			parsedChildren, err := Parse(string(parenBuffer))
			if err != nil {
				return nil, err
			}
			node.Children = parsedChildren
			parenBuffer = make([]rune, 0)

			nodes = satisfyOperators(nodes, node)
		}

		// start parsing variable name or number
		if len(buffer) == 0 && (isLetter(char) || unicode.IsNumber(char)) {
			buffer = append(buffer, char)

			if isLastChar {
				if !nodes.expectsNextNode() {
					return nil, ErrInvalidOperation
				}
				node := createVarOrNumberNode(buffer)
				nodes = satisfyOperators(nodes, node)
			}
		}
	}

	if len(parenStack) > 0 {
		return nil, ErrInvalidParentheses
	}

	return nodes, nil
}

func wrapLastNode(nodes Tree, operation Node) (Tree, error) {
	lastNode, ok := nodes.Last()
	if !ok {
		return nil, ErrInvalidOperation
	}

	node := Node{
		Kind:     KindParentheses,
		Children: []Node{*lastNode, operation},
	}
	nodes = nodes.Pop()
	return append(nodes, node), nil

}

func satisfyOperators(nodes Tree, node Node) Tree {
	lastNode, ok := nodes.Last()
	if !ok {
		return append(nodes, node)
	}

	if lastNode.needSecondOperand() {
		lastNode.Children = append(lastNode.Children, node)
		return nodes
	}

	return append(nodes, node)
}

func createVarOrNumberNode(buffer []rune) Node {
	node := Node{}

	if unicode.IsLetter(buffer[0]) {
		node.Kind = KindVar
	} else if containsDot(buffer) {
		node.Kind = KindFloat
	} else {
		node.Kind = KindInteger
	}

	node.Value = string(buffer)
	return node
}
