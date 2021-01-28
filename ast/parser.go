package ast

import (
	"errors"

	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/utils"
)

type Parser struct {
	Tokens  []lexer.Token
	current int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		current: 0,
		Tokens:  tokens,
	}
}

func (parser *Parser) Parse() Node {
	return parser.Expression()
}

func (parser *Parser) Expression() Node {
	return parser.Equality()
}

func (parser *Parser) Equality() Node {
	node := parser.Comparison()
	for parser.Match(lexer.BANG_EQUAL, lexer.EQUAL) {
		token := parser.Previous()
		right := parser.Comparison()
		node = &BinaryExpr{
			Left:     node,
			Operator: token,
			Right:    right,
		}
	}
	return node
}

func (parser *Parser) Comparison() Node {
	node := parser.Term()
	for parser.Match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		token := parser.Previous()
		right := parser.Term()
		node = &BinaryExpr{
			Left:     node,
			Operator: token,
			Right:    right,
		}

	}
	return node
}

func (parser *Parser) Term() Node {
	node := parser.Factor()
	for parser.Match(lexer.MINUS, lexer.PLUS) {
		token := parser.Previous()
		right := parser.Factor()
		node = &BinaryExpr{
			Left:     node,
			Operator: token,
			Right:    right,
		}
	}
	return node
}

func (parser *Parser) Factor() Node {
	node := parser.Unary()
	for parser.Match(lexer.SLASH, lexer.STAR) {
		token := parser.Previous()
		right := parser.Unary()
		node = &BinaryExpr{
			Left:     node,
			Operator: token,
			Right:    right,
		}
	}
	return node
}

func (parser *Parser) Unary() Node {
	if parser.Match(lexer.BANG, lexer.MINUS) {
		token := parser.Previous()
		node := parser.Unary()
		return &UnaryExpr{
			Operator: token,
			Right:    node,
		}
	}
	return parser.Primary()
}

func (parser *Parser) Primary() Node {
	if parser.Match(lexer.TRUE) {
		return &LiteralExpr{
			Value: true,
		}
	}
	if parser.Match(lexer.FALSE) {
		return &LiteralExpr{
			Value: false,
		}
	}
	if parser.Match(lexer.NIL) {
		return &LiteralExpr{
			Value: nil,
		}
	}
	if parser.Match(lexer.LEFT_PAREN) {
		node := parser.Expression()
		// TODO: error handle
		err := parser.Consume(lexer.RIGHT_PAREN, utils.UNMATCHED_PAREN)
		if err != nil {
			return &GroupExpr{
				Expr: node,
			}
		}
		return &GroupExpr{
			Expr: node,
		}

	}
	if parser.Match(lexer.NUMBER) {
		return &LiteralExpr{
			Value: parser.Previous().Value,
		}
	}
	if parser.Match(lexer.STRING) {
		return &LiteralExpr{
			Value: parser.Previous().Value,
		}
	}
	return nil
}

func (parser *Parser) Consume(tokenType int, message string) error {
	if parser.Check(tokenType) {
		return nil
	}
	return errors.New(message)
}

func (parser *Parser) Match(tokenTypes ...int) bool {
	for _, tokenType := range tokenTypes {
		if parser.Check(tokenType) {
			parser.Advance()
			return true
		}
	}
	return false
}

func (parser *Parser) Check(tokenType int) bool {
	if parser.isAtEnd() {
		return false
	}
	return parser.Peek().Type == tokenType
}

func (parser *Parser) isAtEnd() bool {
	return parser.current >= len(parser.Tokens)
}

func (parser *Parser) Peek() lexer.Token {
	return parser.Tokens[parser.current]
}

func (parser *Parser) Advance() lexer.Token {
	parser.current++
	return parser.Previous()
}

func (parser *Parser) Previous() lexer.Token {
	return parser.Tokens[parser.current-1]
}

func (parser *Parser) Synchronize() {
	parser.Advance()
	for !parser.isAtEnd() {
		if parser.Previous().Type == lexer.SEMICOLON {
			return
		}

		switch parser.Peek().Type {
		case lexer.CLASS:
		case lexer.FUN:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return
		}

		parser.Advance()
	}
}
