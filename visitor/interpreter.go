package visitor

import (
	"fmt"

	"github.com/jameslahm/glox/ast"
	"github.com/jameslahm/glox/environment"
	"github.com/jameslahm/glox/glox_error"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/utils"
	"github.com/spf13/cast"
)

type AstInterpreter struct {
	// DefaultVisitor
	env *environment.Environment
}

func NewAstInterpreter() *AstInterpreter {
	return &AstInterpreter{
		env: environment.NewEnvironment(nil),
	}
}

func (v *AstInterpreter) VisitBinaryExpr(node *ast.BinaryExpr) interface{} {
	leftValue := node.Left.Accept(v)
	rightValue := node.Right.Accept(v)

	switch node.Operator.Type {
	case lexer.MINUS:
		v.CheckNumberOperands(node.Operator, leftValue, rightValue)
		return cast.ToFloat64(leftValue) - cast.ToFloat64(rightValue)
	case lexer.SLASH:
		v.CheckNumberOperands(node.Operator, leftValue, rightValue)
		return cast.ToFloat64(leftValue) / cast.ToFloat64(rightValue)
	case lexer.STAR:
		v.CheckNumberOperands(node.Operator, leftValue, rightValue)
		return cast.ToFloat64(leftValue) * cast.ToFloat64(rightValue)
	case lexer.PLUS:
		if _, ok := leftValue.(float64); ok {
			return cast.ToFloat64(leftValue) + cast.ToFloat64(rightValue)
		} else {
			return cast.ToString(leftValue) + cast.ToString(rightValue)
		}
	case lexer.GREATER:
		v.CheckNumberOperands(node.Operator, leftValue, rightValue)
		return cast.ToFloat64(leftValue) > cast.ToFloat64(rightValue)
	case lexer.GREATER_EQUAL:
		v.CheckNumberOperands(node.Operator, leftValue, rightValue)
		return cast.ToFloat64(leftValue) >= cast.ToFloat64(rightValue)
	case lexer.LESS:
		v.CheckNumberOperands(node.Operator, leftValue, rightValue)
		return cast.ToFloat64(leftValue) < cast.ToFloat64(rightValue)
	case lexer.LESS_EQUAL:
		v.CheckNumberOperands(node.Operator, leftValue, rightValue)
		return cast.ToFloat64(leftValue) <= cast.ToFloat64(rightValue)
	case lexer.BANG_EQUAL:
		return leftValue != rightValue
	case lexer.EQUAL_EQUAL:
		return leftValue == rightValue
	default:
		return nil
	}
}

func (v *AstInterpreter) VisitAssignment(node *ast.Assignment) interface{} {
	value := node.Expr.Accept(v)
	v.env.Assign(node.Name, value)
	return value
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
	case lexer.MINUS:
		v.CheckNumberOperand(node.Operator, value)
		return -cast.ToFloat64(value)
	case lexer.BANG:
		return !cast.ToBool(value)
	default:
		return nil
	}
}

func (v *AstInterpreter) VisitExprStatement(node *ast.ExprStatement) interface{} {
	return node.Expr.Accept(v)
}

func (v *AstInterpreter) VisitPrintStatement(node *ast.PrintStatement) interface{} {
	value := node.Node.Accept(v)
	fmt.Println(value)
	return value
}

func (v *AstInterpreter) VisitVarDeclaration(node *ast.VarDeclaration) interface{} {
	if node.Expr != nil {
		value := node.Expr.Accept(v)
		v.env.Define(node.Name.Lexeme, value)
	} else {
		v.env.Define(node.Name.Lexeme, nil)
	}
	return nil
}

func (v *AstInterpreter) VisitVariable(node *ast.Variable) interface{} {
	return v.env.Get(node.Name)
}

func (v *AstInterpreter) VisitProgram(node *ast.Program) interface{} {
	for _, statement := range node.Statements {
		statement.Accept(v)
	}
	return nil
}

func (v *AstInterpreter) VisitBlockStatement(node *ast.BlockStatement) interface{} {
	parentEnv := v.env
	newEnv := environment.NewEnvironment(v.env)
	v.env = newEnv
	for _, statement := range node.Statements {
		statement.Accept(v)
	}
	v.env = parentEnv
	return nil
}

func (v *AstInterpreter) CheckNumberOperand(token lexer.Token, value interface{}) {
	if _, ok := value.(float64); !ok {
		panic(glox_error.NewRuntimeError(utils.INVALID_OPERAND_NUMBER, token))
	}
}

func (v *AstInterpreter) CheckNumberOperands(token lexer.Token, lefValue interface{}, rightValue interface{}) {
	_, leftOk := lefValue.(float64)
	_, rightOk := rightValue.(float64)
	if !leftOk || !rightOk {
		panic(glox_error.NewRuntimeError(utils.INVALID_OPERAND_NUMBERS, token))
	}
}
