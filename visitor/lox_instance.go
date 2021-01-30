package visitor

import (
	"fmt"

	"github.com/jameslahm/glox/glox_error"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/utils"
)

type LoxInstance struct {
	Class  *LoxClass
	Fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		Class:  class,
		Fields: make(map[string]interface{}),
	}
}

func (instance *LoxInstance) Get(token lexer.Token) interface{} {
	if v, ok := instance.Fields[token.Lexeme]; ok {
		return v
	} else {
		if v, ok := instance.Class.Methods[token.Lexeme]; ok {
			return v.Bind(instance)
		}
		panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_PROPERTY, token.Lexeme), token))
	}
}

func (instance *LoxInstance) Set(token lexer.Token, value interface{}) {
	instance.Fields[token.Lexeme] = value
}
