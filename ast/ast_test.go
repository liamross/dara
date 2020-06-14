package ast

import (
	"dara/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&DeclareStatement{
				Token: token.Token{Type: token.DECLARE, Literal: ":="},
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

	if s := program.String(); s != "a := b;" {
		t.Errorf("program.String() wrong. Got: %q", s)
	}
}
