package ast

import (
	"dara/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "a"},
					Value: "a",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "b"},
					Value: "b",
				},
			},
		},
	}

	if s := program.String(); s != "let a = b;" {
		t.Errorf("program.String() wrong. Got: %q", s)
	}
}
