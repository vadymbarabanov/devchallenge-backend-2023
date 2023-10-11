package parser_test

import (
	"dev-challenge/internal/parser"
	"errors"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		err   error
		want  []parser.Node
	}{
		{
			name:  "invalid parentheses",
			input: "2+(4/2(",
			err:   parser.ErrInvalidParentheses,
			want:  nil,
		},
		{
			name:  "vaild formula",
			input: "=A1*(-A2+cell_3)/0.5",
			err:   nil,
			want: []parser.Node{
				{
					Kind: parser.KindOpEqual,
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
											Value: "cell_3",
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
			},
		},
	}

	invalidOperations := []string{"5+", "5-", "*5", "5*", "/5", "5/", "5(2+2)", "(2+2)5"}

	t.Run("invalid operations", func(t *testing.T) {
		for _, invalidOp := range invalidOperations {
			_, err := parser.Parse(invalidOp)
			if !errors.Is(parser.ErrInvalidOperation, err) {
				t.Fatalf("expected %v to bo invalid operation: want (%v) got (%v)", invalidOp, parser.ErrInvalidOperation, err)
			}
		}
	})

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, err := parser.Parse(test.input)
			if !errors.Is(test.err, err) {
				t.Fatalf("want (%v) get (%v)", test.err, err)
			}
			compareNodes(t, test.want, got)
		})
	}
}

func compareNodes(t *testing.T, want, got []parser.Node) {
	if got == nil && want != nil {
		t.Fatalf("want (%v) got (%v)", want, got)
	}

	if len(want) != len(got) {
		t.Fatalf("nodes length does not match, want (%d) got (%d)", len(want), len(got))
	}

	for i, node := range got {
		switch {
		case node.Kind != want[i].Kind:
			t.Fatalf("node kind mismatch: want (%v) got (%v)", want[i].Kind, node.Kind)
		case node.Value != want[i].Value:
			t.Fatalf("node value mismatch: want (%v) got (%v)", want[i].Value, node.Value)
		case len(node.Children) > 0:
			compareNodes(t, want[i].Children, node.Children)
		}
	}
}
