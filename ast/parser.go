package ast

import (
	"errors"
	"fmt"

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
	var statements []Node
	for !parser.isAtEnd() {
		statement := parser.Declaration()
		if statement != nil {
			statements = append(statements, statement)
		}
	}
	return &Program{
		Statements: statements,
	}
}

func (parser *Parser) Statement() Node {
	if parser.Match(lexer.PRINT) {
		return parser.PrintStatement()
	}
	return parser.ExprStatement()
}

func (parser *Parser) PrintStatement() Node {
	node := parser.Expression()
	parser.Consume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_VALUE)
	return &PrintStatement{
		Node: node,
	}
}

func (parser *Parser) ExprStatement() Node {
	node := parser.Expression()
	parser.Consume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_VALUE)
	return &ExprStatement{
		Expr: node,
	}
}

func (parser *Parser) Declaration() Node {
	if parser.Match(lexer.VAR) {
		return parser.VarDeclaration()
	}

	defer func() {
		if r := recover(); r != nil {
			parser.Synchronize()
		}

	}()

	return parser.Statement()
}

func (parser *Parser) VarDeclaration() Node {
	name, err := parser.Consume(lexer.IDENTIFIER, utils.EXPECT_VARIABLE_NAME)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var initializer Node
	if parser.Match(lexer.EQUAL) {
		initializer = parser.Expression()
	}
	_, err = parser.Consume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_VARIABLE_DECLARATION)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return &VarDeclaration{
		Name: name,
		Expr: initializer,
	}

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
		_, err := parser.Consume(lexer.RIGHT_PAREN, utils.UNMATCHED_PAREN)
		if err != nil {
			fmt.Println(err)
			panic(err)
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
	if parser.Match(lexer.IDENTIFIER) {
		return &Variable{
			Name: parser.Previous(),
		}
	}
	return nil
}

func (parser *Parser) Consume(tokenType int, message string) (lexer.Token, error) {
	if parser.Check(tokenType) {
		return parser.Peek(), nil
	}
	return lexer.Token{}, errors.New(message)
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
