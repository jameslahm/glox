package glox_error

import (
	"github.com/jameslahm/glox/lexer"
)

type RuntimeError struct {
	message string
	token   lexer.Token
}

func NewRuntimeError(message string, token lexer.Token) *RuntimeError {
	return &RuntimeError{
		message: message,
		token:   token,
	}
}
