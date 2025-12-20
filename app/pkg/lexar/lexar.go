package lexar

import (
	"fmt"
	"strconv"
)

type TokenType string

const (
	// single character tokens
	LeftParen  TokenType = "LEFT_PAREN"
	RightParen TokenType = "RIGHT_PAREN"
	LeftBrace  TokenType = "LEFT_BRACE"
	RightBrace TokenType = "RIGHT_BRACE"
	Star       TokenType = "STAR"
	Dot        TokenType = "DOT"
	Comma      TokenType = "COMMA"
	Plus       TokenType = "PLUS"
	Minus      TokenType = "MINUS"
	Semicolon  TokenType = "SEMICOLON"
	Slash      TokenType = "SLASH"

	// one or two character tokens
	Bang         TokenType = "BANG"
	BangEqual    TokenType = "BANG_EQUAL"
	Equal        TokenType = "EQUAL"
	EqualEqual   TokenType = "EQUAL_EQUAL"
	Greater      TokenType = "GREATER"
	GreaterEqual TokenType = "GREATER_EQUAL"
	Less         TokenType = "LESS"
	LessEqual    TokenType = "LESS_EQUAL"

	// literals
	Identifier  TokenType = "IDENTIFIER"
	StringToken TokenType = "STRING"
	Number      TokenType = "NUMBER"

	// keywords
	AndToken    TokenType = "AND"
	ClassToken  TokenType = "CLASS"
	ElseToken   TokenType = "ELSE"
	FalseToken  TokenType = "FALSE"
	FunToken    TokenType = "FUN"
	ForToken    TokenType = "FOR"
	IfToken     TokenType = "IF"
	NilToken    TokenType = "NIL"
	OrToken     TokenType = "OR"
	PrintToken  TokenType = "PRINT"
	ReturnToken TokenType = "RETURN"
	SuperToken  TokenType = "SUPER"
	ThisToken   TokenType = "THIS"
	TrueToken   TokenType = "TRUE"
	VarToken    TokenType = "VAR"
	WhileToken  TokenType = "WHILE"

	Eof        TokenType = "EOF"
	ErrorToken TokenType = "ERROR"
)

var reservedKeywords = map[string]TokenType{
	"and":    AndToken,
	"class":  ClassToken,
	"else":   ElseToken,
	"false":  FalseToken,
	"fun":    FunToken,
	"for":    ForToken,
	"if":     IfToken,
	"nil":    NilToken,
	"or":     OrToken,
	"print":  PrintToken,
	"return": ReturnToken,
	"super":  SuperToken,
	"this":   ThisToken,
	"true":   TrueToken,
	"var":    VarToken,
	"while":  WhileToken,
}

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
}

func (t *Token) IsToken(tokenType TokenType) bool {
	return t.TokenType == tokenType
}

func (t *Token) ToString() string {
	literal := ""
	if t.TokenType == Number {
		num := t.Literal.(float64)
		if num == float64(int64(num)) {
			literal = fmt.Sprintf("%.1f", num)
		} else {
			literal = fmt.Sprintf("%g", num)
		}
	} else if t.Literal == nil {
		literal = "null"
	} else if t.TokenType == StringToken {
		literal = t.Literal.(string)
	}
	return string(t.TokenType) + " " + t.Lexeme + " " + literal
}

type Lexar struct {
	input []byte
	index int
	start int
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
		return 0
	}
	return l.input[l.index]
}

func (l *Lexar) peekNext() byte {
	if l.index+1 >= len(l.input) {
		return 0
	}
	return l.input[l.index+1]
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

func isNumber(char byte) bool {
	return char >= '0' && char <= '9'
}

func (l *Lexar) NextToken() Token {
	l.start = l.index
	inputChar := l.next()
	token := Token{}
	switch inputChar {
	case '(':
		token.TokenType = LeftParen
	case ')':
		token.TokenType = RightParen
	case '{':
		token.TokenType = LeftBrace
	case '}':
		token.TokenType = RightBrace
	case '*':
		token.TokenType = Star
	case '+':
		token.TokenType = Plus
	case ',':
		token.TokenType = Comma
	case '.':
		token.TokenType = Dot
	case '-':
		token.TokenType = Minus
	case ';':
		token.TokenType = Semicolon
	case '\n':
		l.line++
		return l.NextToken()
	case '=':
		token.TokenType = Equal
		if l.match('=') {
			token.TokenType = EqualEqual
		}
	case '!':
		token.TokenType = Bang
		if l.match('=') {
			token.TokenType = BangEqual
		}
	case '>':
		token.TokenType = Greater
		if l.match('=') {
			token.TokenType = GreaterEqual
		}
	case '<':
		token.TokenType = Less
		if l.match('=') {
			token.TokenType = LessEqual
		}
	case '\t', ' ', '\r':
		return l.NextToken()
	case '"':
		for l.peek() != '"' && !l.eof() {
			l.next()
		}
		token.TokenType = StringToken
		token.Literal = string(l.input[l.start+1 : l.index])

		if l.peek() != '"' {
			token.TokenType = ErrorToken
			token.Literal = fmt.Errorf("[line %d] Error: Unterminated string.", l.line)
		}
		l.next()

	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.eof() {
				// skip comments
				l.next()
			}
			return l.NextToken()
		}
		token.TokenType = Slash
	case 0:
		token.TokenType = Eof
	default:
		if isNumber(inputChar) {
			token.TokenType = Number
			token.Literal = l.parseNumber()

		} else if isAlpha(inputChar) {
			for isAlphaNumeric(l.peek()) {
				l.next()
			}
			if val, ok := reservedKeywords[string(l.input[l.start:l.index])]; ok {
				token.TokenType = val
			} else {
				token.TokenType = Identifier
			}
		} else {
			token.TokenType = ErrorToken
			token.Literal = fmt.Errorf("[line %d] Error: Unexpected character: %v", l.line, string(inputChar))
		}
	}

	token.Lexeme = string(l.input[l.start:l.index])

	return token

}
func (l *Lexar) Scan() ([]Token, []error) {
	tokens := []Token{}
	errorArr := []error{}
	token := l.NextToken()
	for token.TokenType != Eof {
		if token.TokenType == ErrorToken {
			errorArr = append(errorArr, token.Literal.(error))
		} else {
			tokens = append(tokens, token)
		}
		token = l.NextToken()
	}

	tokens = append(tokens, Token{TokenType: Eof, Lexeme: ""})

	return tokens, errorArr
}

func isAlphaNumeric(char byte) bool {
	return isAlpha(char) || isNumber(char)
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func (l *Lexar) parseNumber() float64 {
	for isNumber(l.peek()) {
		l.next()
	}

	if l.peek() == '.' && isNumber(l.peekNext()) {
		l.next()

		for isNumber(l.peek()) {
			l.next()
		}
	}

	// dont need to handle error we check validate it above
	num, _ := strconv.ParseFloat(string(l.input[l.start:l.index]), 64)

	return num
}
