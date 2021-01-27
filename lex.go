package glox

import (
	"fmt"
	"strconv"

	"github.com/jameslahm/glox/utils"
)

type Lexer struct {
	Source string
	Tokens []Token

	start   int
	current int
	line    int

	hasError bool
}

func NewLexer(source string) *Lexer {
	var lexer = Lexer{
		Source:  source,
		current: 0,
		line:    1,
	}
	return &lexer
}

func (lexer *Lexer) Lex() {
	for !lexer.IsAtEnd() {
		lexer.start = lexer.current
		lexer.Scan()
	}
}

func (lexer *Lexer) Scan() {
	c := lexer.Advance()
	switch c {
	case '(':
		lexer.AddToken(LEFT_PAREN, nil)
	case ')':
		lexer.AddToken(RIGHT_PAREN, nil)
	case '{':
		lexer.AddToken(LEFT_BRACE, nil)
	case '}':
		lexer.AddToken(RIGHT_BRACE, nil)
	case ',':
		lexer.AddToken(COMMA, nil)
	case '.':
		lexer.AddToken(DOT, nil)
	case '-':
		lexer.AddToken(MINUS, nil)
	case '+':
		lexer.AddToken(PLUS, nil)
	case ';':
		lexer.AddToken(SEMICOLON, nil)
	case '*':
		lexer.AddToken(STAR, nil)
	case '/':
		if lexer.Match('/') {
			for !lexer.Match('\n') && !lexer.IsAtEnd() {
				lexer.Advance()
			}
		} else {
			lexer.AddToken(SLASH, nil)
		}
	case '!':
		if lexer.Match('=') {
			lexer.AddToken(BANG_EQUAL, nil)
		} else {
			lexer.AddToken(BANG, nil)
		}
	case '=':
		if lexer.Match('=') {
			lexer.AddToken(EQUAL_EQUAL, nil)
		} else {
			lexer.AddToken(EQUAL, nil)
		}
	case '<':
		if lexer.Match('=') {
			lexer.AddToken(LESS_EQUAL, nil)
		} else {
			lexer.AddToken(LESS, nil)
		}
	case '>':
		if lexer.Match('=') {
			lexer.AddToken(GREATER_EQUAL, nil)
		} else {
			lexer.AddToken(GREATER, nil)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		lexer.line++
		break
	case '"':
		lexer.start++
		lexer.AddStringToken()
	case 'o':
		lexer.AddToken(OR)
	case 
	default:
		if utils.IsDigit(c) {
			lexer.AddNumberToken()
			break
		}
		lexer.hasError = true
		utils.Error(lexer.line, fmt.Sprintf("%s %c", utils.UNEXPECTED_CHARACTER_MESSAGE, c))
	}
}

func (lexer *Lexer) Advance() byte {
	lexer.current++
	return lexer.Source[lexer.current-1]
}

func (lexer *Lexer) Peek() byte {
	return lexer.Source[lexer.current]
}

func (lexer *Lexer) PeekNext() byte {
	return lexer.Source[lexer.current+1]
}

func (lexer *Lexer) AddToken(tokenType int, value interface{}) {
	lexeme := lexer.Source[lexer.start:lexer.current]
	token := NewToken(tokenType, lexeme, value, lexer.line)
	lexer.Tokens = append(lexer.Tokens, *token)
}

func (lexer *Lexer) AddTokenWithLexeme(tokenType int, value interface{}, lexeme string) {
	token := NewToken(tokenType, lexeme, value, lexer.line)
	lexer.Tokens = append(lexer.Tokens, *token)
}

func (lexer *Lexer) Match(c byte) bool {
	if lexer.IsAtEnd() {
		return false
	}
	currentByte := lexer.Source[lexer.current]
	return currentByte == c
}

func (lexer *Lexer) IsAtEnd() bool {
	return lexer.current >= len(lexer.Source)
}

func (lexer *Lexer) AddStringToken() {
	for !lexer.IsAtEnd() && !lexer.Match('"') {
		c := lexer.Advance()
		if c == '\n' {
			lexer.line++
		}
	}

	if lexer.IsAtEnd() {
		utils.Error(lexer.line, utils.UNTERMINATED_STRING)
	}

	lexeme := lexer.Source[lexer.start : lexer.current-1]
	value := lexeme
	lexer.AddTokenWithLexeme(STRING, value, lexeme)
}

func (lexer *Lexer) AddNumberToken() {
	for utils.IsDigit(lexer.Peek()) {
		lexer.Advance()
	}

	if lexer.Peek() == '.' && utils.IsDigit(lexer.PeekNext()) {
		lexer.Advance()
		lexer.Advance()

		for utils.IsDigit(lexer.Peek()) {
			lexer.Advance()
		}
	}

	lexeme := lexer.Source[lexer.start:lexer.current]
	value, err := strconv.ParseDouble(lexeme)
	if err != nil {
		utils.Error(lexer.line, fmt.Sprintf("%s %s", utils.INVALID_NUMBER, lexeme))
	}

	lexer.AddToken(NUMBER, value)
}
