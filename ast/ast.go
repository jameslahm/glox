package ast

import "github.com/jameslahm/glox/lexer"

type Visitor interface {
	VisitBinaryExpr(node *BinaryExpr) interface{}
	VisitUnaryExpr(node *UnaryExpr) interface{}
	VisitGroupExpr(node *GroupExpr) interface{}
	VisitLiteralExpr(node *LiteralExpr) interface{}
	VisitExprStatement(node *ExprStatement) interface{}
	VisitPrintStatement(node *PrintStatement) interface{}
	VisitProgram(node *Program) interface{}
	VisitVarDeclaration(node *VarDeclaration) interface{}
	VisitVariable(node *Variable) interface{}
	VisitAssignment(node *Assignment) interface{}
}

type Node interface {
	Accept(v Visitor) interface{}
}

type BinaryExpr struct {
	Left     Node
	Right    Node
	Operator lexer.Token
}

func (node *BinaryExpr) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(node)
}

type UnaryExpr struct {
	Operator lexer.Token
	Right    Node
}

func (node *UnaryExpr) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(node)
}

type GroupExpr struct {
	Expr Node
}

func (node *GroupExpr) Accept(v Visitor) interface{} {
	return v.VisitGroupExpr(node)
}

type LiteralExpr struct {
	Value interface{}
}

func (node *LiteralExpr) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(node)
}

type Variable struct {
	Name lexer.Token
}

func (node *Variable) Accept(v Visitor) interface{} {
	return v.VisitVariable(node)
}

type ExprStatement struct {
	Expr Node
}

func (node *ExprStatement) Accept(v Visitor) interface{} {
	return v.VisitExprStatement(node)
}

type PrintStatement struct {
	Node Node
}

func (node *PrintStatement) Accept(v Visitor) interface{} {
	return v.VisitPrintStatement(node)
}

type VarDeclaration struct {
	Name lexer.Token
	Expr Node
}

func (node *VarDeclaration) Accept(v Visitor) interface{} {
	return v.VisitVarDeclaration(node)
}

type Program struct {
	Statements []Node
}

func (node *Program) Accept(v Visitor) interface{} {
	return v.VisitProgram(node)
}

type Assignment struct {
	Name lexer.Token
	Expr Node
}

func (node *Assignment) Accept(v Visitor) interface{} {
	return v.VisitAssignment(node)
}
