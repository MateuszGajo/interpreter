package token

import "fmt"

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
	NumberInt   TokenType = "NUMBER_INT"
	NumberFloat TokenType = "NUMBER_FLOAT"

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

var ReservedKeywords = map[string]TokenType{
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
	if t.TokenType == NumberFloat {
		literal := ""
		num := t.Literal.(float64)
		if num == float64(int64(num)) {
			literal = fmt.Sprintf("%.1f", num)
		} else {
			literal = fmt.Sprintf("%g", num)
		}

		return literal

	} else if t.TokenType == NumberInt {
		return fmt.Sprintf("%.1f", float64(t.Literal.(int64)))
	} else if t.Literal == nil {
		return "null"
	}
	return t.Literal.(string)
}
