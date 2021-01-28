package ast_test

import (
	"testing"

	. "github.com/jameslahm/glox/ast"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/visitor"
	"gopkg.in/go-playground/assert.v1"
)

func TestAstVisitor(t *testing.T) {
	var node = BinaryExpr{
		Left: &UnaryExpr{
			lexer.Token{lexer.MINUS, "-", nil, 1},
			&LiteralExpr{123},
		},
		Operator: lexer.Token{lexer.STAR, "*", nil, 1},
		Right: &GroupExpr{
			&LiteralExpr{45.67},
		},
	}
	output := node.Accept(&visitor.AstPrinter{})
	assert.Equal(t, output, "(* (- 123) (group 45.67))")
}
