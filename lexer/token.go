package lexer

import (
	"fmt"
)

const (
	// Single Character tokens
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// Relation And Logic tokens
	AND
	OR

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals tokens
	IDENTIFIER
	STRING
	NUMBER
	NIL
	FALSE
	TRUE

	// Keywords tokens
	CLASS

	// Control Flow tokens
	IF
	ELSE
	FOR
	WHILE
	RETURN

	// Internal support tokens
	PRINT
	SUPER
	THIS

	// Declaration tokens
	FUN
	VAR

	EOF
)

type Token struct {
	Type   int
	Lexeme string
	Value  interface{}
	Line   int
}

func (token *Token) String() string {
	return fmt.Sprintf("%d %s %d", token.Type, token.Lexeme, token.Line)
}

func NewToken(tokenType int, lexeme string, value interface{}, line int) *Token {
	return &Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   line,
		Value:  value,
	}
}
