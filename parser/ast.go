package parser

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
	left     Node
	right    Node
	operator glox.Token
}

func (node *BinaryExpr) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(node)
}

type UnaryExpr struct {
	operator glox.Token
	right    Node
}

func (node *UnaryExpr) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(node)
}

type GroupExpr struct {
	expr Node
}

func (node *GroupExpr) Accept(v Visitor) interface{} {
	return v.VisitGroupExpr(node)
}

type LiteralExpr struct {
	value interface{}
}

func (node *LiteralExpr) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(node)
}
