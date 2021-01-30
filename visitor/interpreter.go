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
	Env                  *environment.Env
	_originEnvStack      []*environment.Env
	VariableBindDistance map[ast.Node]int
}

func NewAstInterpreter(variableBindDistances map[ast.Node]int) *AstInterpreter {
	interpreter := &AstInterpreter{
		Env:                  environment.NewEnvironment(nil),
		VariableBindDistance: variableBindDistances,
	}

	interpreter.Env.Define("clock", &Clock{})

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
	loxFunction := NewLoxFunction(node, v.Env)
	v.Env.Define(node.Name.Lexeme, loxFunction)
	return nil
}

func (v *AstInterpreter) VisitReturnStatement(node *ast.ReturnStatement) interface{} {
	var value interface{}
	if node.Expr != nil {
		value = node.Expr.Accept(v)
	}
	panic(value)
}

func (v *AstInterpreter) VisitAssignment(node *ast.Assignment) interface{} {
	value := node.Expr.Accept(v)
	distance := v.VariableBindDistance[node]
	v.Env.Assign(node.Name, value, distance)
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
		v.Env.Define(node.Name.Lexeme, value)
	} else {
		v.Env.Define(node.Name.Lexeme, nil)
	}
	return nil
}

func (v *AstInterpreter) VisitVariable(node *ast.Variable) interface{} {
	distance, ok := v.VariableBindDistance[node]
	if !ok {
		panic(glox_error.NewRuntimeError("No Distance", node.Name))
	}
	return v.Env.Get(node.Name, distance)
}

func (v *AstInterpreter) VisitProgram(node *ast.Program) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%#v\n", r)
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
	if f, ok := callee.(LoxCallable); ok {
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

func (v *AstInterpreter) VisitGetExpr(node *ast.GetExpr) interface{} {
	expr := node.Expr.Accept(v)
	if instance, ok := expr.(*LoxInstance); ok {
		return instance.Get(node.Name)
	} else {
		panic(glox_error.NewRuntimeError(utils.ONLY_INSTANCES_HAVE_PROPERTIES, node.Name))
	}
}

func (v *AstInterpreter) VisitSetExpr(node *ast.SetExpr) interface{} {
	expr := node.Expr.Accept(v)
	if instance, ok := expr.(*LoxInstance); ok {
		value := node.Value.Accept(v)
		instance.Set(node.Name, value)
		return value
	} else {
		panic(glox_error.NewRuntimeError(utils.ONLY_INSTANCES_HAVE_PROPERTIES, node.Name))
	}
}

func (v *AstInterpreter) VisitClassDeclaration(node *ast.ClassDeclaration) interface{} {
	v.Env.Define(node.Name.Lexeme, nil)

	var superClass *LoxClass

	if node.SuperClass != nil {
		var ok = false
		superClass, ok = node.SuperClass.Accept(v).(*LoxClass)
		if !ok {
			panic(glox_error.NewRuntimeError(utils.SUPER_CLASS_MUST_BE_CLASS, node.Name))
		}

		v.EnterScope()
		v.Env.Define("super", superClass)
	}

	class := NewLoxClass(node.Name.Lexeme, superClass)
	for _, method := range node.Methods {
		var isInitializer = false
		if method.Name.Lexeme == "init" {
			isInitializer = true
		}
		class.Methods[method.Name.Lexeme] = NewLoxFunction(method, v.Env, isInitializer)
	}
	if node.SuperClass != nil {
		v.ExitScope()
	}
	v.Env.Assign(node.Name, class, 0)
	return nil
}

func (v *AstInterpreter) VisitSuperExpr(node *ast.SuperExpr) interface{} {
	distance := v.VariableBindDistance[node]
	superClass := v.Env.Get(node.Keyword, distance)

	token := node.Keyword
	token.Lexeme = "this"
	token.Type = lexer.THIS
	instance := v.Env.Get(token, distance-1)

	if class, ok := superClass.(*LoxClass); ok {
		if ins, ok := instance.(*LoxInstance); ok {
			return class.GetMethod(node.Method).Bind(ins)
		}
	}
	panic(glox_error.NewRuntimeError(utils.SUPER_CLASS_MUST_BE_CLASS, node.Keyword))

}

func (v *AstInterpreter) VisitWhileStatement(node *ast.WhileStatement) interface{} {
	value := cast.ToBool(node.Expr.Accept(v))
	for value {
		node.Then.Accept(v)
		value = cast.ToBool(node.Expr.Accept(v))
	}
	return nil
}

func (v *AstInterpreter) VisitThisExpr(node *ast.ThisExpr) interface{} {
	distance := v.VariableBindDistance[node]
	return v.Env.Get(node.Keyword, distance)

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

func (v *AstInterpreter) NewExecuteScope(e *environment.Env) {
	v._originEnvStack = append(v._originEnvStack, v.Env)
	v.Env = environment.NewEnvironment(e)
}

func (v *AstInterpreter) RestoreExecuteScope() {
	originEnvLens := len(v._originEnvStack)
	v.Env = v._originEnvStack[originEnvLens-1]
	v._originEnvStack = v._originEnvStack[:originEnvLens-1]
}
