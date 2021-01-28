package visitor

import (
	"github.com/jameslahm/glox"
	"github.com/jameslahm/glox/ast"
	"github.com/spf13/cast"
)

type AstInterpreter struct {
}

func (v *AstInterpreter) VisitBinaryExpr(node *ast.BinaryExpr) interface{} {
	leftValue := node.Left.Accept(v)
	rightValue := node.Right.Accept(v)

	switch node.Operator.Type {
	case glox.MINUS:
		return cast.ToFloat64(leftValue) - cast.ToFloat64(rightValue)
	case glox.SLASH:
		return cast.ToFloat64(leftValue) / cast.ToFloat64(rightValue)
	case glox.STAR:
		return cast.ToFloat64(leftValue) * cast.ToFloat64(rightValue)
	case glox.PLUS:
		if _, ok := leftValue.(float64); ok {
			return cast.ToFloat64(leftValue) + cast.ToFloat64(rightValue)
		} else {
			return cast.ToString(leftValue) + cast.ToString(rightValue)
		}
	case glox.GREATER:
		return cast.ToFloat64(leftValue) > cast.ToFloat64(rightValue)
	case glox.GREATER_EQUAL:
		return cast.ToFloat64(leftValue) >= cast.ToFloat64(rightValue)
	case glox.LESS:
		return cast.ToFloat64(leftValue) < cast.ToFloat64(rightValue)
	case glox.LESS_EQUAL:
		return cast.ToFloat64(leftValue) <= cast.ToFloat64(rightValue)
	case glox.BANG_EQUAL:
		return leftValue != rightValue
	case glox.EQUAL_EQUAL:
		return leftValue == rightValue
	default:
		return nil
	}
}

func (v *AstInterpreter) VisitLiteralExpr(node *ast.LiteralExpr) interface{} {
	return node.Value
}

func (v *AstInterpreter) VisitGroupExpr(node *ast.GroupExpr) interface{} {
	return node.Expr.Accept(v)
}

func (v *AstInterpreter) VisitUnaryExpr(node *ast.UnaryExpr) interface{} {
	value := node.Right.Accept(v)
	switch node.Operator.Type {
	case glox.MINUS:
		return -cast.ToFloat64(value)
	case glox.BANG:
		return !cast.ToBool(value)
	default:
		return nil
	}
}
