package evaluator_test

import (
	"dev-challenge/internal/evaluator"
	"dev-challenge/internal/parser"
	"errors"
	"testing"
)

func TestEvaluator_Evaluate(t *testing.T) {
	// =A1*(-A2+A3)/0.5

	var want float64 = -4

	input := []parser.Node{
		{
			Kind:  parser.KindOpEqual,
			Value: string(parser.OpEqual),
		},
		{
			Kind: parser.KindParentheses,
			Children: []parser.Node{
				{
					Kind: parser.KindParentheses,
					Children: []parser.Node{
						{
							Kind:  parser.KindVar,
							Value: "A1",
						},
						{
							Kind: parser.KindOpMultiply,
						},
						{
							Kind: parser.KindParentheses,
							Children: []parser.Node{
								{
									Kind: parser.KindParentheses,
									Children: []parser.Node{
										{
											Kind:  parser.KindInteger,
											Value: "-1",
										},
										{
											Kind: parser.KindOpMultiply,
										},
										{
											Kind:  parser.KindVar,
											Value: "A2",
										},
									},
								},
								{
									Kind: parser.KindOpPlus,
								},
								{
									Kind:  parser.KindVar,
									Value: "A3",
								},
							},
						},
					},
				},
				{
					Kind: parser.KindOpDivide,
				},
				{
					Kind:  parser.KindFloat,
					Value: "0.5",
				},
			},
		},
	}

	result, err := evaluator.Evaluate(input, getFormulaByID)
	if err != nil {
		t.Fatalf("want (<nil>) got (%v)", err)
	}

	if result != want {
		t.Fatalf("want (%v) got (%v)", want, result)
	}
}

func getFormulaByID(id string) (string, error) {
	switch id {
	case "A1":
		return "2", nil
	case "A2":
		return "=A1+A1", nil
	case "A3":
		return "=A2-1", nil

	default:
		return "", errors.New("cell not found")
	}
}
