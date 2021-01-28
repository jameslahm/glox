package environment

import (
	"fmt"

	"github.com/jameslahm/glox/glox_error"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/utils"
)

type Env struct {
	Values map[string]interface{}
	Parent *Env
}

func NewEnvironment(parent *Env) *Env {
	return &Env{
		Parent: parent,
		Values: make(map[string]interface{}),
	}
}

func (e *Env) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Env) Assign(token lexer.Token, value interface{}) {
	if _, ok := e.Values[token.Lexeme]; ok {
		e.Values[token.Lexeme] = value
	} else {
		if e.Parent != nil {
			e.Parent.Assign(token, value)
			return
		}
		panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_VARIABLE, &token.Lexeme), token))
	}
}

func (e *Env) Get(token lexer.Token) interface{} {
	if v, ok := e.Values[token.Lexeme]; ok {
		return v
	} else {
		if e.Parent != nil {
			return e.Parent.Get(token)
		}
		panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_VARIABLE, &token.Lexeme), token))
	}
}
