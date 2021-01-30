package visitor

import (
	"github.com/jameslahm/glox/ast"
	"github.com/jameslahm/glox/environment"
	"github.com/jameslahm/glox/lexer"
)

type LoxFunction struct {
	Node          *ast.FuncDeclaration
	Env           *environment.Env
	IsInitializer bool
}

type LoxCallable interface {
	Call(v *AstInterpreter, arguments []interface{}) interface{}
	Arity() int
}

func (f *LoxFunction) Call(v *AstInterpreter, arguments []interface{}) (ret interface{}) {
	ok := false

	v.NewExecuteScope(f.Env)
	for i, param := range f.Node.Params {
		v.Env.Define(param.Lexeme, arguments[i])
	}
	defer func() {
		r := recover()
		if !ok {
			v.RestoreExecuteScope()
			ret = r
			if f.IsInitializer {
				ret = f.GetThis()
			}
		}
	}()
	f.Node.Body.Accept(v)
	v.RestoreExecuteScope()
	ok = true
	if f.IsInitializer {
		return f.GetThis()
	}
	return nil
}

func (f *LoxFunction) GetThis() interface{} {
	if !f.IsInitializer {
		panic("No this in normal function")
	}
	token := f.Node.Name
	token.Lexeme = "this"
	token.Type = lexer.THIS
	return f.Env.Get(token, 0)
}

func (f *LoxFunction) Arity() int {
	return len(f.Node.Params)
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	env := f.Env
	newEnv := environment.NewEnvironment(env)
	newEnv.Define("this", instance)
	return &LoxFunction{
		Node: f.Node,
		Env:  newEnv,
	}
}

func NewLoxFunction(node *ast.FuncDeclaration, env *environment.Env, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		Node:          node,
		Env:           env,
		IsInitializer: isInitializer,
	}
}
