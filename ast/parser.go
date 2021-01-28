package ast

import (
	"errors"

	"github.com/jameslahm/glox"
	"github.com/jameslahm/glox/utils"
)

type Parser struct {
	Tokens  []glox.Token
	current int
}

func (parser *Parser) Expression() Node {
	return parser.Equality()
}

func (parser *Parser) Equality() Node {
	node := parser.Comparison()
	for parser.Match(glox.BANG_EQUAL, glox.EQUAL) {
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
	for parser.Match(glox.GREATER, glox.GREATER_EQUAL, glox.LESS, glox.LESS_EQUAL) {
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
	for parser.Match(glox.MINUS, glox.PLUS) {
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
	for parser.Match(glox.SLASH, glox.STAR) {
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
	if parser.Match(glox.BANG, glox.MINUS) {
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
	if parser.Match(glox.TRUE) {
		return &LiteralExpr{
			Value: true,
		}
	}
	if parser.Match(glox.FALSE) {
		return &LiteralExpr{
			Value: false,
		}
	}
	if parser.Match(glox.NIL) {
		return &LiteralExpr{
			Value: nil,
		}
	}
	if parser.Match(glox.LEFT_PAREN) {
		node := parser.Expression()
		// TODO: error handle
		err := parser.Consume(glox.RIGHT_PAREN, utils.UNMATCHED_PAREN)
		if err != nil {
			return &GroupExpr{
				Expr: node,
			}
		}
		return &GroupExpr{
			Expr: node,
		}

	}
	if parser.Match(glox.NUMBER) {
		return &LiteralExpr{
			Value: parser.Previous().Value,
		}
	}
	if parser.Match(glox.STRING) {
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

func (parser *Parser) Peek() glox.Token {
	return parser.Tokens[parser.current]
}

func (parser *Parser) Advance() glox.Token {
	parser.current++
	return parser.Previous()
}

func (parser *Parser) Previous() glox.Token {
	return parser.Tokens[parser.current-1]
}

func (parser *Parser) Synchronize() {
	parser.Advance()
	for !parser.isAtEnd() {
		if parser.Previous().Type == glox.SEMICOLON {
			return
		}

		switch parser.Peek().Type {
		case glox.CLASS:
		case glox.FUN:
		case glox.VAR:
		case glox.FOR:
		case glox.IF:
		case glox.WHILE:
		case glox.PRINT:
		case glox.RETURN:
			return
		}

		parser.Advance()
	}
}
