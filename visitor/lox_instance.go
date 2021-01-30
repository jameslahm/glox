package visitor

import (
	"github.com/jameslahm/glox/lexer"
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
		return instance.Class.GetMethod(token).Bind(instance)
	}
}

func (instance *LoxInstance) Set(token lexer.Token, value interface{}) {
	instance.Fields[token.Lexeme] = value
}
