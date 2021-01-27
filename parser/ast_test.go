package parser

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/jameslahm/glox"
	"gopkg.in/go-playground/assert.v1"
)

type AstPrinter struct{}

func (v *AstPrinter) VisitBinaryExpr(node *BinaryExpr) interface{} {
	return v.Parenthesize(node.operator.Lexeme, node.left, node.right)
}

func (v *AstPrinter) VisitUnaryExpr(node *UnaryExpr) interface{} {
	return v.Parenthesize(node.operator.Lexeme, node.right)
}

func (v *AstPrinter) VisitLiteralExpr(node *LiteralExpr) interface{} {
	if node.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", node.value)
}

func (v *AstPrinter) VisitGroupExpr(node *GroupExpr) interface{} {
	return v.Parenthesize("group", node.expr)
}

func (v *AstPrinter) Parenthesize(lexeme string, exprs ...Node) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(lexeme)
	for _, node := range exprs {
		sb.WriteString(" ")
		s, ok := node.Accept(v).(string)
		if !ok {
			log.Fatalf("Error get output string")
		}
		sb.WriteString(s)
	}
	sb.WriteString(")")
	return sb.String()
}

func TestAstVisitor(t *testing.T) {
	var node = BinaryExpr{
		left: &UnaryExpr{
			glox.Token{glox.MINUS, "-", nil, 1},
			&LiteralExpr{123},
		},
		operator: glox.Token{glox.STAR, "*", nil, 1},
		right: &GroupExpr{
			&LiteralExpr{45.67},
		},
	}
	output := node.Accept(&AstPrinter{})
	assert.Equal(t, output, "(* (- 123) (group 45.67))")
}
