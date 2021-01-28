package environment

import (
	"fmt"

	"github.com/jameslahm/glox/glox_error"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/utils"
)

type Environment struct {
	Values map[string]interface{}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Assign(token lexer.Token, value interface{}) {
	if _, ok := e.Values[token.Lexeme]; ok {
		e.Values[token.Lexeme] = value
	} else {
		panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_VARIABLE, &token.Lexeme), token))
	}
}

func (e *Environment) Get(token lexer.Token) interface{} {
	if v, ok := e.Values[token.Lexeme]; ok {
		return v
	} else {
		panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_VARIABLE, &token.Lexeme), token))
	}
}
