package ast

import (
	"testing"

	"github.com/oliversabler/apa/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "låt",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "minVariabel",
					},
					Value: "minVariabel",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "enAnnanVariabel",
					},
					Value: "enAnnanVariabel",
				},
			},
		},
	}

	if program.String() != "låt minVariabel = enAnnanVariabel;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
