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
	*environment.Env
}

func NewAstInterpreter() *AstInterpreter {
	interpreter := &AstInterpreter{
		Env: environment.NewEnvironment(nil),
	}

	interpreter.Define("clock", &Clock{})

	return interpreter
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

func (v *AstInterpreter) VisitFuncDeclaration(node *ast.FuncDeclaration) interface{} {
	v.Env.Define(node.Name.Lexeme, node)
	return nil
}

func (v *AstInterpreter) VisitReturnStatement(node *ast.ReturnStatement) interface{} {
	var value interface{}
	if node.Expr != nil {
		value = node.Expr.Accept(v)
		fmt.Printf("%#v\n", value)
	}
	panic(value)
}

func (v *AstInterpreter) VisitAssignment(node *ast.Assignment) interface{} {
	value := node.Expr.Accept(v)
	v.Env.Assign(node.Name, value)
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
	return value
}

func (v *AstInterpreter) VisitVarDeclaration(node *ast.VarDeclaration) interface{} {
	if node.Expr != nil {
		value := node.Expr.Accept(v)
		v.Env.Define(node.Name.Lexeme, value)
	} else {
		v.Env.Define(node.Name.Lexeme, nil)
	}
	return nil
}

func (v *AstInterpreter) VisitVariable(node *ast.Variable) interface{} {
	return v.Env.Get(node.Name)
}

func (v *AstInterpreter) VisitProgram(node *ast.Program) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	for _, statement := range node.Statements {
		statement.Accept(v)
	}
	return nil
}

func (v *AstInterpreter) VisitBlockStatement(node *ast.BlockStatement) interface{} {
	v.EnterScope()
	for _, statement := range node.Statements {
		statement.Accept(v)
	}
	v.ExitScope()
	return nil
}

func (v *AstInterpreter) VisitCallExpr(node *ast.CallExpr) interface{} {
	callee := node.Callee.Accept(v)

	var arguments []interface{}
	for _, param := range node.Arguments {
		value := param.Accept(v)
		arguments = append(arguments, value)
	}
	if f, ok := callee.(ast.GloxCallable); ok {
		if f.Arity() == len(node.Arguments) {
			return f.Call(v, arguments)
		} else {
			panic(glox_error.NewRuntimeError(utils.MISMATCH_CALL_PARAMS_LENGTH, node.Paren))
		}
	} else {
		panic(glox_error.NewRuntimeError(utils.ONLY_CALL_FUNCTION_AND_CLASS, node.Paren))
	}
}

func (v *AstInterpreter) VisitIfStatement(node *ast.IfStatement) interface{} {
	value := cast.ToBool(node.Expr.Accept(v))
	if value {
		node.Then.Accept(v)
	} else {
		if node.Else != nil {
			node.Else.Accept(v)
		}
	}
	return nil
}

func (v *AstInterpreter) VisitLogicalExpr(node *ast.LogicalExpr) interface{} {
	if node.Operator.Type == lexer.AND {
		leftValue := node.Left.Accept(v)
		if !cast.ToBool(leftValue) {
			return leftValue
		} else {
			rightValue := node.Right.Accept(v)
			return rightValue
		}
	}
	if node.Operator.Type == lexer.OR {
		leftValue := node.Left.Accept(v)
		if !cast.ToBool(leftValue) {
			rightValue := cast.ToBool(node.Right.Accept(v))
			return rightValue
		} else {
			return leftValue
		}

	}
	return nil
}

func (v *AstInterpreter) VisitWhileStatement(node *ast.WhileStatement) interface{} {
	value := cast.ToBool(node.Expr.Accept(v))
	for value {
		node.Then.Accept(v)
		value = cast.ToBool(node.Expr.Accept(v))
	}
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

func (v *AstInterpreter) EnterScope() {
	newEnv := environment.NewEnvironment(v.Env)
	v.Env = newEnv
}

func (v *AstInterpreter) ExitScope() {
	// TODO parent = nil
	parent := v.Env.Parent
	v.Env = parent
}
