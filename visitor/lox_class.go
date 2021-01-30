package visitor

import (
	"fmt"

	"github.com/jameslahm/glox/glox_error"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/utils"
)

type LoxClass struct {
	Name       string
	Methods    map[string]*LoxFunction
	SuperClass *LoxClass
}

func NewLoxClass(name string, superClass *LoxClass) *LoxClass {
	return &LoxClass{
		Name:       name,
		Methods:    make(map[string]*LoxFunction),
		SuperClass: superClass,
	}
}

func (c *LoxClass) Call(v *AstInterpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		Class: c,
	}
	if initializer, ok := c.Methods["init"]; ok {
		initializer.Bind(instance).Call(v, arguments)
	}
	return instance
}

func (c *LoxClass) Arity() int {
	if initialzer, ok := c.Methods["init"]; ok {
		return initialzer.Arity()
	}
	return 0
}

func (c *LoxClass) GetMethod(token lexer.Token) *LoxFunction {
	if v, ok := c.Methods[token.Lexeme]; ok {
		return v
	}
	if c.SuperClass != nil {
		return c.SuperClass.GetMethod(token)
	}
	panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_PROPERTY, token.Lexeme), token))
}
