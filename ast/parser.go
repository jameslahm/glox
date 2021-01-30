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
	if parser.Match(lexer.IF) {
		return parser.IfStatement()
	}
	if parser.Match(lexer.WHILE) {
		return parser.WhileStatement()
	}
	if parser.Match(lexer.FOR) {
		return parser.ForStatement()
	}
	if parser.Match(lexer.RETURN) {
		return parser.ReturnStatement()
	}
	return parser.ExprStatement()
}

func (parser *Parser) ReturnStatement() Node {
	keyword := parser.Previous()
	var expr Node
	if !parser.Check(lexer.SEMICOLON) {
		expr = parser.Expression()
	}
	parser.MustConsume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_RETURN)
	return &ReturnStatement{
		Keyword: keyword,
		Expr:    expr,
	}
}

func (parser *Parser) ForStatement() Node {
	parser.MustConsume(lexer.LEFT_PAREN, utils.EXPECT_LEFT_PAREN_AFTER_FOR)

	var initializer Node
	if parser.Match(lexer.VAR) {
		initializer = parser.VarDeclaration()

	} else if parser.Match(lexer.SEMICOLON) {
		initializer = nil
	} else {
		initializer = parser.ExprStatement()
	}

	var condition Node
	if parser.Match(lexer.SEMICOLON) {
		condition = &LiteralExpr{
			Value: true,
		}
	} else {
		condition = parser.Expression()
		parser.MustConsume(lexer.SEMICOLON, utils.EXPECT_SEMICOLON_AFTER_LOOP_CONDITION)
	}

	var increment Node
	if parser.Match(lexer.RIGHT_PAREN) {
		increment = nil
	} else {
		increment = parser.Expression()
		parser.MustConsume(lexer.RIGHT_PAREN, utils.EXPECT_RIGHT_PAREN_AFTER_CLAUSES)
	}

	body := parser.Statement()
	if increment != nil {
		body = &BlockStatement{
			Statements: []Node{body, increment},
		}
	}

	body = &WhileStatement{
		Expr: condition,
		Then: body,
	}

	if initializer != nil {
		body = &BlockStatement{
			Statements: []Node{initializer, body},
		}
	}

	return body
}

func (parser *Parser) WhileStatement() Node {
	parser.MustConsume(lexer.LEFT_PAREN, utils.EXPECT_LEFT_PAREN_AFTER_WHILE)
	expr := parser.Expression()
	parser.MustConsume(lexer.RIGHT_PAREN, utils.EXPECT_RIGHT_PAREN_AFTER_CONDITION)

	statement := parser.Statement()

	return &WhileStatement{
		Expr: expr,
		Then: statement,
	}
}

func (parser *Parser) IfStatement() Node {
	parser.MustConsume(lexer.LEFT_PAREN, utils.EXPECT_LEFT_PAREN_AFTER_IF)
	expr := parser.Expression()
	parser.MustConsume(lexer.RIGHT_PAREN, utils.EXPECT_RIGHT_PAREN_AFTER_IF_CONDITION)
	thenStatement := parser.Statement()

	var elseStatement Node
	if parser.Match(lexer.ELSE) {
		elseStatement = parser.Statement()
	}
	return &IfStatement{
		Expr: expr,
		Then: thenStatement,
		Else: elseStatement,
	}
}

func (parser *Parser) LogicOr() Node {
	node := parser.LogicAnd()
	for !parser.isAtEnd() && parser.Match(lexer.OR) {
		right := parser.LogicAnd()
		node = &LogicalExpr{
			Left:     node,
			Right:    right,
			Operator: parser.Previous(),
		}
	}
	return node
}

func (parser *Parser) LogicAnd() Node {
	node := parser.Equality()
	for !parser.isAtEnd() && parser.Match(lexer.AND) {
		right := parser.Equality()
		node = &LogicalExpr{
			Left:     node,
			Right:    right,
			Operator: parser.Previous(),
		}
	}
	return node
}

func (parser *Parser) BlockStatement() Node {
	var statements []Node
	for !parser.Check(lexer.RIGHT_BRACE) && !parser.isAtEnd() {
		statement := parser.Declaration()
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

	if parser.Match(lexer.FUN) {
		return parser.FuncDeclaration()
	}

	if parser.Match(lexer.CLASS) {
		return parser.ClassDeclaration()
	}

	return parser.Statement()
}

func (parser *Parser) ClassDeclaration() Node {
	name := parser.MustConsume(lexer.IDENTIFIER, utils.EXPECT_CLASS_NAME)
	var superClass *Variable

	if parser.Match(lexer.LESS) {
		token := parser.MustConsume(lexer.IDENTIFIER, utils.EXPECT_SUPER_CLASS_NAME)
		superClass = &Variable{
			Name: token,
		}
	}

	parser.MustConsume(lexer.LEFT_BRACE, utils.EXPECT_LEFT_BRACE_BEFORE_CLASS_BODY)

	var methods []*FuncDeclaration
	for !parser.isAtEnd() && !parser.Check(lexer.RIGHT_BRACE) {
		method := parser.FuncDeclaration()
		methods = append(methods, method)
	}
	parser.MustConsume(lexer.RIGHT_BRACE, utils.EXPECT_RIGHT_BRACE_AFTER_CLASS_BODY)
	return &ClassDeclaration{
		Name:       name,
		Methods:    methods,
		SuperClass: superClass,
	}
}

func (parser *Parser) FuncDeclaration() *FuncDeclaration {
	name := parser.MustConsume(lexer.IDENTIFIER, utils.EXPECT_FUNCTION_NAME)
	parser.MustConsume(lexer.LEFT_PAREN, utils.EXPECT_LEFT_PAREN_AFTER_FUNCTION_NAME)

	var parameters []lexer.Token

	for !parser.isAtEnd() && !parser.Check(lexer.RIGHT_PAREN) {
		for {
			param := parser.MustConsume(lexer.IDENTIFIER, utils.EXPECT_PARAM_NAME)
			parameters = append(parameters, param)
			if !parser.Match(lexer.COMMA) {
				break
			}
		}
	}
	parser.MustConsume(lexer.RIGHT_PAREN, utils.EXPECT_RIGHT_PAREN_AFTER_PARAMETERS)
	body := parser.Statement()

	return &FuncDeclaration{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
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
	for parser.Match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
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
	expr := parser.LogicOr()
	if parser.Match(lexer.EQUAL) {
		value := parser.Assignment()
		if v, ok := expr.(*Variable); ok {
			return &Assignment{
				Name: v.Name,
				Expr: value,
			}
		} else if v, ok := expr.(*GetExpr); ok {
			return &SetExpr{
				Expr:  v.Expr,
				Name:  v.Name,
				Value: value,
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
	return parser.Call()
}

func (parser *Parser) Call() Node {
	expr := parser.Primary()

	for parser.Match(lexer.LEFT_PAREN) || parser.Match(lexer.DOT) {
		if parser.Previous().Type == lexer.LEFT_PAREN {
			var arguments []Node
			for !parser.Check(lexer.RIGHT_PAREN) && !parser.isAtEnd() {
				arg := parser.Expression()
				arguments = append(arguments, arg)
				if !parser.Match(lexer.COMMA) {
					break
				}
			}
			parser.MustConsume(lexer.RIGHT_PAREN, utils.EXPECT_RIGHT_PAREN_AFTER_ARGUMENTS)
			if len(arguments) > 255 {
				utils.Error(parser.Previous().Line, utils.WARN_NO_MORE_THAN_MAXIMUM_ARGUMENTS)
			}
			expr = &CallExpr{
				Callee:    expr,
				Arguments: arguments,
				Paren:     parser.Previous(),
			}
		} else {
			name := parser.MustConsume(lexer.IDENTIFIER, utils.EXPECT_PROPERTY_NAME_AFTER_DOT)
			expr = &GetExpr{
				Expr: expr,
				Name: name,
			}
		}

	}

	return expr
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
	if parser.Match(lexer.THIS) {
		return &ThisExpr{
			Keyword: parser.Previous(),
		}
	}
	if parser.Match(lexer.IDENTIFIER) {
		return &Variable{
			Name: parser.Previous(),
		}
	}
	if parser.Match(lexer.SUPER) {
		keyword := parser.Previous()
		parser.MustConsume(lexer.COMMA, utils.EXPECT_DOT_AFTER_SUPER)
		token := parser.MustConsume(lexer.IDENTIFIER, utils.EXPECT_SUPER_CLASS_NAME)

		return &SuperExpr{
			Keyword: keyword,
			Method:  token,
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
