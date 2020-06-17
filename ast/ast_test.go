package ast

import (
	"dara/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&IfStatement{
				Token: token.Token{Type: token.DECLARE, Literal: "if"},
				Condition: &Boolean{
					Token: token.Token{Type: token.TRUE, Literal: "true"},
					Value: true,
				},
				Consequence: &BlockStatement{},
			},
		},
	}

	if s := program.String(); s != "if true {  }" {
		t.Errorf("program.String() wrong. Got: %q", s)
	}
}
