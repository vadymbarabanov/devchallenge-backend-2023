package evaluator

import (
	"dev-challenge/internal/parser"
	"strconv"
)

func Evaluate(tree parser.Tree, getFormulaByID func(string) (string, error)) (float64, error) {
	result := 0.0
	bufferedValue := 0.0
	operation := parser.Node{}

	for _, node := range tree {
		switch {
		case node.IsParentheses():
			res, err := Evaluate(node.Children, getFormulaByID)
			if err != nil {
				return 0, err
			}
			bufferedValue = res

		case node.IsOperation():
			if node.Kind != parser.KindOpEqual {
				operation = node
			}
			continue

		case node.IsVar():
			formula, err := getFormulaByID(node.Value)
			if err != nil {
				return 0, err
			}
			parsedFormula, err := parser.Parse(formula)
			if err != nil {
				return 0, err
			}
			res, err := Evaluate(parsedFormula, getFormulaByID)
			if err != nil {
				return 0, err
			}
			bufferedValue = res
		case node.IsNumber():
			val, err := strconv.ParseFloat(node.Value, 64)
			if err != nil {
				return 0, err
			}
			bufferedValue = val
		}

		if operation.Kind == "" {
			result = bufferedValue
		} else {
			switch operation.Kind {
			case parser.KindOpPlus:
				result += bufferedValue
			case parser.KindOpMinus:
				result -= bufferedValue
			case parser.KindOpMultiply:
				result *= bufferedValue
			case parser.KindOpDivide:
				result /= bufferedValue
			default:
				return 0, parser.ErrInvalidOperation
			}

			operation = parser.Node{}
		}
	}

	return result, nil
}
