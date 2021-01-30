package visitor

import (
	"fmt"

	"github.com/jameslahm/glox/ast"
	"github.com/jameslahm/glox/environment"
)

type DefaultVisitor struct {
	*environment.Env
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

func (v *DefaultVisitor) VisitIfStatement(node *ast.IfStatement) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitLogicalExpr(node *ast.LogicalExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitWhileStatement(node *ast.WhileStatement) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitCallExpr(node *ast.CallExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitProgram(node *ast.Program) interface{} {
	for _, statement := range node.Statements {
		fmt.Printf("%T", v)
		statement.Accept(v)
	}
	return nil
}

func (v *DefaultVisitor) VisitFuncDeclaration(node *ast.FuncDeclaration) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitReturnStatement(node *ast.ReturnStatement) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitClassDeclaration(node *ast.ClassDeclaration) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitGetExpr(node *ast.GetExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitSetExpr(node *ast.SetExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitThisExpr(node *ast.ThisExpr) interface{} {
	return nil
}

func (v *DefaultVisitor) VisitSuperExpr(node *ast.SuperExpr) interface{} {
	return nil
}
