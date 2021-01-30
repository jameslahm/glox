package visitor

import (
	"github.com/jameslahm/glox/ast"
	"github.com/jameslahm/glox/environment"
)

type LoxFunction struct {
	Node *ast.FuncDeclaration
	Env  *environment.Env
}

type GloxCallable interface {
	Call(v *AstInterpreter, arguments []interface{}) interface{}
	Arity() int
}

func (f *LoxFunction) Call(v *AstInterpreter, arguments []interface{}) (ret interface{}) {
	v.NewExecuteScope(f.Env)
	for i, param := range f.Node.Params {
		v.Env.Define(param.Lexeme, arguments[i])
	}
	defer func() {
		r := recover()
		if r != nil {
			v.RestoreExecuteScope()
			ret = r
		}
	}()
	f.Node.Body.Accept(v)
	v.RestoreExecuteScope()
	return nil
}

func (f *LoxFunction) Arity() int {
	return len(f.Node.Params)
}
