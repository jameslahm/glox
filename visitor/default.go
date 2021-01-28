package visitor

import (
	"fmt"

	"github.com/jameslahm/glox/ast"
)

type DefaultVisitor struct {
}

func (v *DefaultVisitor) VisitBinaryExpr(node *ast.BinaryExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitUnaryExpr(node *ast.UnaryExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitLiteralExpr(node *ast.LiteralExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitGroupExpr(node *ast.GroupExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitExprStatement(node *ast.ExprStatement) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitPrintStatement(node *ast.PrintStatement) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitVarDeclaration(node *ast.VarDeclaration) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitVariable(node *ast.Variable) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitAssignment(node *ast.Assignment) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitBlockStatement(node *ast.BlockStatement) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitProgram(node *ast.Program) interface{} {
	for _, statement := range node.Statements {
		fmt.Printf("%T", v)
		statement.Accept(v)
	}
	return nil
}
