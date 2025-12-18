package lexar

import (
	"fmt"
)

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
	slash      TokenType = "SLASH"

	// one or two character tokens
	bang         TokenType = "BANG"
	bangEqual    TokenType = "BANG_EQUAL"
	equal        TokenType = "EQUAL"
	equalEqual   TokenType = "EQUAL_EQUAL"
	greater      TokenType = "GREATER"
	greaterEqual TokenType = "GREATER_EQUAL"
	less         TokenType = "LESS"
	lessEqual    TokenType = "LESS_EQUAL"

	// literals
	identifier TokenType = "IDENTIFIER"
	stringTok  TokenType = "STRING"
	number     TokenType = "NUMBER"

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

func (l *Lexar) next() byte {
	char := l.peek()
	l.index++
	return char
}

func (l *Lexar) peek() byte {
	if l.eof() {
		return '$'
	}
	return l.input[l.index]
}

func NewLexar(input []byte) *Lexar {
	return &Lexar{
		input: input,
		index: 0,
		line:  1,
	}
}

func (l *Lexar) match(expected byte) bool {
	isMatch := l.peek() == expected
	if isMatch {
		l.index++
	}

	return isMatch
}

func (l *Lexar) Scan() ([]Token, []error) {
	tokens := []Token{}
	errorArr := []error{}
	for !l.eof() {
		currentIndex := l.index
		inputChar := l.next()
		token := Token{}
		switch inputChar {
		case '(':
			token.tokenType = leftParen
		case ')':
			token.tokenType = rightParen
		case '{':
			token.tokenType = leftBrace
		case '}':
			token.tokenType = rightBrace
		case '*':
			token.tokenType = star
		case '+':
			token.tokenType = plus
		case ',':
			token.tokenType = comma
		case '.':
			token.tokenType = dot
		case '-':
			token.tokenType = minus
		case ';':
			token.tokenType = semicolon
		case '\n':
			l.line++
			continue
		case '=':
			token.tokenType = equal
			if l.match('=') {
				token.tokenType = equalEqual
			}
		case '!':
			token.tokenType = bang
			if l.match('=') {
				token.tokenType = bangEqual
			}
		case '>':
			token.tokenType = greater
			if l.match('=') {
				token.tokenType = greaterEqual
			}
		case '<':
			token.tokenType = less
			if l.match('=') {
				token.tokenType = lessEqual
			}
		case '\t', ' ', '\r':
			continue
		case '"':
			for l.peek() != '"' && !l.eof() {
				l.next()
			}
			token.tokenType = stringTok
			token.literal = string(l.input[currentIndex+1 : l.index])

			if l.peek() != '"' {
				errorArr = append(errorArr, fmt.Errorf("[line %d] Error: Unterminated string.", l.line))
				continue
			}
			l.next()

		case '/':
			if l.match('/') {
				for l.peek() != '\n' && !l.eof() {
					// skip comments
					l.next()
				}
				continue
			}
			token.tokenType = slash
		default:
			errorArr = append(errorArr, fmt.Errorf("[line %d] Error: Unexpected character: %v", l.line, string(inputChar)))
			continue
		}
		token.lexeme = string(l.input[currentIndex:l.index])
		tokens = append(tokens, token)
	}
	tokens = append(tokens, Token{tokenType: eof, lexeme: ""})

	return tokens, errorArr
}
