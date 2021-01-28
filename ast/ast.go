package ast

import "github.com/jameslahm/glox"

type Visitor interface {
	VisitBinaryExpr(node *BinaryExpr) interface{}
	VisitUnaryExpr(node *UnaryExpr) interface{}
	VisitGroupExpr(node *GroupExpr) interface{}
	VisitLiteralExpr(node *LiteralExpr) interface{}
}

type Node interface {
	Accept(v Visitor) interface{}
}

type BinaryExpr struct {
	Left     Node
	Right    Node
	Operator glox.Token
}

func (node *BinaryExpr) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(node)
}

type UnaryExpr struct {
	Operator glox.Token
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
