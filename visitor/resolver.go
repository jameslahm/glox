package visitor

import (
	"errors"

	"github.com/jameslahm/glox/ast"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/utils"
)

const (
	FunctionInit = iota
	FunctionNormal
	None
	Class
)

type Resolver struct {
	Scopes                   []map[string]bool
	Errors                   []error
	VariableBindingDistances map[ast.Node]int
	InFunctionType           int
	InClassType              int
}

func NewResolver() *Resolver {
	var scopes = []map[string]bool{{}}

	return &Resolver{
		Scopes:                   scopes,
		VariableBindingDistances: make(map[ast.Node]int),
		InFunctionType:           None,
		InClassType:              None,
	}
}

func (v *Resolver) VisitBlockStatement(node *ast.BlockStatement) interface{} {
	v.EnterScope()
	for _, statement := range node.Statements {
		statement.Accept(v)
	}
	v.ExitScope()
	return nil
}

func (v *Resolver) VisitVarDeclaration(node *ast.VarDeclaration) interface{} {
	v.Declare(node.Name)
	if node.Expr != nil {
		node.Expr.Accept(v)
	}
	v.Define(node.Name)
	return nil
}

func (v *Resolver) Declare(token lexer.Token) {
	scope := v.GetCurrentScope()
	if _, ok := scope[token.Lexeme]; ok {
		err := errors.New(utils.ALREADY_DECLARE_VARIABLE)
		v.Errors = append(v.Errors, err)
		utils.Error(token.Line, err.Error())
	}
	scope[token.Lexeme] = false
}

func (v *Resolver) Define(token lexer.Token) {
	scope := v.GetCurrentScope()
	scope[token.Lexeme] = true
}

func (v *Resolver) EnterScope() {
	var scope = make(map[string]bool)
	v.Scopes = append(v.Scopes, scope)
}

func (v *Resolver) ExitScope() {
	v.Scopes = v.Scopes[:len(v.Scopes)-1]
}

func (v *Resolver) GetCurrentScope() map[string]bool {
	return v.Scopes[len(v.Scopes)-1]
}

func (v *Resolver) VisitVariable(node *ast.Variable) interface{} {
	scope := v.GetCurrentScope()
	if value, ok := scope[node.Name.Lexeme]; ok && !value {
		err := errors.New(utils.WARN_READ_VARIABLE_BEFORE_DEFINE)
		v.Errors = append(v.Errors, err)
		utils.Error(node.Name.Line, err.Error())
	}
	v.Resolve(node, node.Name.Lexeme)

	return nil
}

func (v *Resolver) VisitAssignment(node *ast.Assignment) interface{} {
	node.Expr.Accept(v)
	v.Resolve(node, node.Name.Lexeme)
	return nil
}

func (v *Resolver) Resolve(node ast.Node, name string) {
	scopesLen := len(v.Scopes)
	for i := scopesLen - 1; i >= 0; i-- {
		if _, ok := v.Scopes[i][name]; ok {
			v.VariableBindingDistances[node] = scopesLen - 1 - i
			return
		}
	}
}

func (v *Resolver) VisitFuncDeclaration(node *ast.FuncDeclaration) interface{} {
	v.Declare(node.Name)
	v.Define(node.Name)

	v.EnterScope()
	for _, param := range node.Params {
		v.Declare(param)
		v.Define(param)
	}
	inFunctionTypeBackup := v.InFunctionType
	if v.InClassType != None && node.Name.Lexeme == "init" {
		v.InFunctionType = FunctionInit
	} else {
		v.InFunctionType = FunctionNormal
	}
	node.Body.Accept(v)
	v.InFunctionType = inFunctionTypeBackup
	v.ExitScope()
	return nil
}

func (v *Resolver) VisitThisExpr(node *ast.ThisExpr) interface{} {
	if v.InClassType != Class {
		err := errors.New(utils.WARN_USE_THIS_OUT_CLASS)
		v.Errors = append(v.Errors, err)
		utils.Error(node.Keyword.Line, err.Error())
	}
	v.Resolve(node, node.Keyword.Lexeme)
	return nil
}

func (v *Resolver) VisitExprStatement(node *ast.ExprStatement) interface{} {
	node.Expr.Accept(v)
	return nil
}

func (v *Resolver) VisitIfStatement(node *ast.IfStatement) interface{} {
	node.Expr.Accept(v)
	node.Then.Accept(v)
	if node.Else != nil {
		node.Else.Accept(v)
	}
	return nil
}

func (v *Resolver) VisitPrintStatement(node *ast.PrintStatement) interface{} {
	node.Node.Accept(v)
	return nil
}

func (v *Resolver) VisitReturnStatement(node *ast.ReturnStatement) interface{} {
	if v.InFunctionType == None {
		err := errors.New(utils.WARN_RETURN_FROM_NOFUNCTION)
		v.Errors = append(v.Errors, err)
		utils.Error(node.Keyword.Line, err.Error())
	}

	if node.Expr != nil {
		if v.InFunctionType == FunctionInit {
			err := errors.New(utils.WARN_RETURN_VALUE_FROM_INIT)
			v.Errors = append(v.Errors, err)
			utils.Error(node.Keyword.Line, err.Error())
		}
		node.Expr.Accept(v)
	}
	return nil
}

func (v *Resolver) VisitWhileStatement(node *ast.WhileStatement) interface{} {
	node.Expr.Accept(v)
	node.Then.Accept(v)
	return nil
}

func (v *Resolver) VisitBinaryExpr(node *ast.BinaryExpr) interface{} {
	node.Left.Accept(v)
	node.Right.Accept(v)
	return nil
}

func (v *Resolver) VisitCallExpr(node *ast.CallExpr) interface{} {
	node.Callee.Accept(v)

	for _, arg := range node.Arguments {
		arg.Accept(v)
	}
	return nil
}

func (v *Resolver) VisitGroupExpr(node *ast.GroupExpr) interface{} {
	node.Expr.Accept(v)
	return nil
}

func (v *Resolver) VisitLiteralExpr(node *ast.LiteralExpr) interface{} {
	return nil
}

func (v *Resolver) VisitLogicalExpr(node *ast.LogicalExpr) interface{} {
	node.Left.Accept(v)
	node.Right.Accept(v)
	return nil
}

func (v *Resolver) VisitUnaryExpr(node *ast.UnaryExpr) interface{} {
	node.Right.Accept(v)
	return nil
}

func (v *Resolver) VisitProgram(node *ast.Program) interface{} {
	for _, statement := range node.Statements {
		statement.Accept(v)
	}
	return nil
}

func (v *Resolver) VisitClassDeclaration(node *ast.ClassDeclaration) interface{} {

	inClassTypeBackUp := v.InClassType

	v.InClassType = Class
	v.Declare(node.Name)
	v.Define(node.Name)

	v.EnterScope()

	scope := v.GetCurrentScope()
	scope["this"] = true
	for _, method := range node.Methods {
		method.Accept(v)
	}

	v.ExitScope()
	v.InClassType = inClassTypeBackUp
	return nil
}

func (v *Resolver) VisitGetExpr(node *ast.GetExpr) interface{} {
	node.Expr.Accept(v)
	return nil
}

func (v *Resolver) VisitSetExpr(node *ast.SetExpr) interface{} {
	node.Expr.Accept(v)
	node.Value.Accept(v)
	return nil
}
