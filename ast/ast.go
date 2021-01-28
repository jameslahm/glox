package ast

import (
	"github.com/jameslahm/glox/lexer"
)

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
	VisitBlockStatement(node *BlockStatement) interface{}
	VisitIfStatement(node *IfStatement) interface{}
	VisitLogicalExpr(node *LogicalExpr) interface{}
	VisitWhileStatement(node *WhileStatement) interface{}
	VisitCallExpr(node *CallExpr) interface{}
	VisitFuncDeclaration(node *FuncDeclaration) interface{}
	VisitReturnStatement(node *ReturnStatement) interface{}

	EnterScope()
	ExitScope()
	Define(name string, value interface{})
	Assign(token lexer.Token, value interface{})
	Get(token lexer.Token) interface{}
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

type BlockStatement struct {
	Statements []Node
}

func (node *BlockStatement) Accept(v Visitor) interface{} {
	return v.VisitBlockStatement(node)
}

type IfStatement struct {
	Expr Node
	Then Node
	Else Node
}

func (node *IfStatement) Accept(v Visitor) interface{} {
	return v.VisitIfStatement(node)
}

type LogicalExpr struct {
	Left     Node
	Operator lexer.Token
	Right    Node
}

func (node *LogicalExpr) Accept(v Visitor) interface{} {
	return v.VisitLogicalExpr(node)
}

type WhileStatement struct {
	Expr Node
	Then Node
}

func (node *WhileStatement) Accept(v Visitor) interface{} {
	return v.VisitWhileStatement(node)
}

type CallExpr struct {
	Callee    Node
	Arguments []Node

	// ? For Location
	Paren lexer.Token
}

func (node *CallExpr) Accept(v Visitor) interface{} {
	return v.VisitCallExpr(node)
}

type GloxCallable interface {
	Call(v Visitor, arguments []interface{}) interface{}
	Arity() int
}

type FuncDeclaration struct {
	Name   lexer.Token
	Params []lexer.Token
	Body   Node
}

func (f *FuncDeclaration) Call(v Visitor, arguments []interface{}) (ret interface{}) {
	v.EnterScope()
	for i, param := range f.Params {
		v.Define(param.Lexeme, arguments[i])
	}
	defer func() {
		r := recover()
		ret = r
	}()
	f.Body.Accept(v)
	v.ExitScope()
	return nil
}

func (f *FuncDeclaration) Arity() int {
	return len(f.Params)
}

func (f *FuncDeclaration) Accept(v Visitor) interface{} {
	return v.VisitFuncDeclaration(f)
}

type ReturnStatement struct {
	Keyword lexer.Token
	Expr    Node
}

func (node *ReturnStatement) Accept(v Visitor) interface{} {
	return v.VisitReturnStatement(node)
}
