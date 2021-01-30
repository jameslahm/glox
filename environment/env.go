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

func (e *Env) Assign(token lexer.Token, value interface{}, distance int) {
	env := e
	for i := 0; i < distance; i++ {
		env = env.Parent
	}
	if _, ok := env.Values[token.Lexeme]; ok {
		env.Values[token.Lexeme] = value
	} else {
		panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_VARIABLE, token.Lexeme), token))
	}
}

func (e *Env) Get(token lexer.Token, distance int) interface{} {
	env := e
	for i := 0; i < distance; i++ {
		env = env.Parent
	}
	if v, ok := env.Values[token.Lexeme]; ok {
		return v
	} else {
		panic(glox_error.NewRuntimeError(fmt.Sprintf(utils.UNDEFINED_VARIABLE, token.Lexeme), token))
	}
}
