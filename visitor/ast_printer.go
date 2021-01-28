package visitor

import (
	"fmt"
	"log"
	"strings"

	"github.com/jameslahm/glox/ast"
)

type AstPrinter struct {
	DefaultVisitor
}

func (v *AstPrinter) VisitBinaryExpr(node *ast.BinaryExpr) interface{} {
	return v.Parenthesize(node.Operator.Lexeme, node.Left, node.Right)
}

func (v *AstPrinter) VisitUnaryExpr(node *ast.UnaryExpr) interface{} {
	return v.Parenthesize(node.Operator.Lexeme, node.Right)
}

func (v *AstPrinter) VisitLiteralExpr(node *ast.LiteralExpr) interface{} {
	if node.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", node.Value)
}

func (v *AstPrinter) VisitGroupExpr(node *ast.GroupExpr) interface{} {
	return v.Parenthesize("group", node.Expr)
}

func (v *AstPrinter) Parenthesize(lexeme string, exprs ...ast.Node) string {
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
