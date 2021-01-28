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
	Errors  []error
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
	if parser.Match(lexer.LEFT_BRACE) {
		return parser.BlockStatement()
	}
	return parser.ExprStatement()
}

func (parser *Parser) BlockStatement() Node {
	var statements []Node
	for !parser.Check(lexer.RIGHT_BRACE) && !parser.isAtEnd() {
		statement := parser.Statement()
		if statement != nil {
			statements = append(statements, statement)
		}
	}
	parser.MustConsume(lexer.RIGHT_BRACE, utils.EXPECT_RIGHT_BRACE_AFTER_BLOCK)
	return &BlockStatement{
		Statements: statements,
	}
}

func (parser *Parser) PrintStatement() Node {
	node := parser.Expression()
	parser.MustConsume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_VALUE)
	return &PrintStatement{
		Node: node,
	}
}

func (parser *Parser) ExprStatement() Node {
	node := parser.Expression()
	parser.MustConsume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_VALUE)
	return &ExprStatement{
		Expr: node,
	}
}

func (parser *Parser) Declaration() Node {
	defer func() {
		if r := recover(); r != nil {
			parser.Synchronize()
		}

	}()

	if parser.Match(lexer.VAR) {
		return parser.VarDeclaration()
	}

	return parser.Statement()
}

func (parser *Parser) VarDeclaration() Node {
	name := parser.MustConsume(lexer.IDENTIFIER, utils.EXPECT_VARIABLE_NAME)
	var initializer Node
	if parser.Match(lexer.EQUAL) {
		initializer = parser.Expression()
	}
	parser.MustConsume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_VARIABLE_DECLARATION)
	return &VarDeclaration{
		Name: name,
		Expr: initializer,
	}
}

func (parser *Parser) Expression() Node {
	return parser.Assignment()
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

func (parser *Parser) Assignment() Node {
	expr := parser.Equality()
	if parser.Match(lexer.EQUAL) {
		value := parser.Assignment()
		if v, ok := expr.(*Variable); ok {
			return &Assignment{
				Name: v.Name,
				Expr: value,
			}
		}

		err := errors.New(utils.INVALID_ASSIGNMENT_TARGET)
		fmt.Println(err)
		panic(err)
	}
	return expr
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
		parser.MustConsume(lexer.RIGHT_PAREN, utils.UNMATCHED_PAREN)
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

func (parser *Parser) MustConsume(tokenType int, message string) lexer.Token {
	if parser.Check(tokenType) {
		parser.current++
		return parser.Previous()
	}

	message = fmt.Sprintf("[line %d] Error %s", parser.Previous().Line, message)
	err := errors.New(message)
	parser.Errors = append(parser.Errors, err)
	fmt.Println(err)
	panic(err)
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
	if parser.isAtEnd() {
		return
	}

	for !parser.isAtEnd() {
		switch parser.Peek().Type {
		case lexer.CLASS,
			lexer.FUN,
			lexer.VAR,
			lexer.FOR,
			lexer.IF,
			lexer.WHILE,
			lexer.PRINT,
			lexer.RETURN:
			return
		case lexer.SEMICOLON:
			parser.Advance()
			return
		}

		parser.Advance()
	}
}
