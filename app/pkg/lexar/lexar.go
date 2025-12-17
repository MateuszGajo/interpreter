package lexar

import "fmt"

type TokenType string

const (
	// single character tokens
	leftParen  TokenType = "LEFT_PAREN"
	rightParen TokenType = "RIGHT_PAREN"
	leftBrace  TokenType = "LEFT_BRACE"
	rightBrace TokenType = "RIGHT_BRACE"
	star       TokenType = "STAR"
	dot        TokenType = "DOT"
	comma      TokenType = "COMMA"
	plus       TokenType = "PLUS"
	minus      TokenType = "MINUS"
	semicolon  TokenType = "SEMICOLON"

	eof TokenType = "EOF"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   string
}

func (t *Token) ToString() string {
	literal := t.literal
	if t.literal == "" {
		literal = "null"
	}
	return string(t.tokenType) + " " + t.lexeme + " " + literal
}

type Lexar struct {
	input []byte
	index int
	line  int
}

func (l *Lexar) eof() bool {
	return l.index >= len(l.input)
}

func (l *Lexar) peekNext() byte {
	if l.eof() {
		return '$'
	}
	index := l.index
	index++

	return l.input[index]
}

func (l *Lexar) advance() byte {
	char := l.input[l.index]
	l.index++
	return char
}

func (l *Lexar) peek() byte {
	return l.input[l.index]
}

func NewLexar(input []byte) *Lexar {
	return &Lexar{
		input: input,
		index: 0,
		line:  1,
	}
}

func (l *Lexar) Scan() ([]Token, []error) {
	tokens := []Token{}
	errors := []error{}
	for !l.eof() {
		token := l.advance()
		switch token {
		case '(':
			tokens = append(tokens, Token{tokenType: leftParen, lexeme: string(token)})
		case ')':
			tokens = append(tokens, Token{tokenType: rightParen, lexeme: string(token)})
		case '{':
			tokens = append(tokens, Token{tokenType: leftBrace, lexeme: string(token)})
		case '}':
			tokens = append(tokens, Token{tokenType: rightBrace, lexeme: string(token)})
		case '*':
			tokens = append(tokens, Token{tokenType: star, lexeme: string(token)})
		case '+':
			tokens = append(tokens, Token{tokenType: plus, lexeme: string(token)})
		case ',':
			tokens = append(tokens, Token{tokenType: comma, lexeme: string(token)})
		case '.':
			tokens = append(tokens, Token{tokenType: dot, lexeme: string(token)})
		case '-':
			tokens = append(tokens, Token{tokenType: minus, lexeme: string(token)})
		case ';':
			tokens = append(tokens, Token{tokenType: semicolon, lexeme: string(token)})
		case '\n':
			l.line++
		default:
			errors = append(errors, fmt.Errorf("[line %d] Error: Unexpected character: %v", l.line, string(token)))
		}
	}
	tokens = append(tokens, Token{tokenType: eof, lexeme: ""})

	return tokens, errors
}
