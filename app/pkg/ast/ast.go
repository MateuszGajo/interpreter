package ast

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

type Expression interface {
	Expression()
	String() string
}

type Boolean struct {
	Value bool
	Token token.Token
}

func (astBool Boolean) Expression() {}
func (astBool Boolean) String() string {
	return astBool.Token.Lexeme
}

type Nil struct{}

func (astNil Nil) Expression() {}
func (astNil Nil) String() string {
	return "nil"
}

type Inteager struct {
	Token token.Token
	Value int64
}

func (astInteager Inteager) Expression() {}
func (astInteager Inteager) String() string {
	return astInteager.Token.ToString()
}

type Float struct {
	Token token.Token
	Value float64
}

func (astFloat Float) Expression() {}
func (astFloat Float) String() string {
	return astFloat.Token.ToString()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (astStringLiteral StringLiteral) Expression() {}
func (astStringLiteral StringLiteral) String() string {
	return astStringLiteral.Token.ToString()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (astIdentifier Identifier) Expression() {}
func (astIdentifier Identifier) String() string {
	return astIdentifier.Token.Lexeme
}

type Grouping struct {
	Exp Expression
}

func (astGrouping Grouping) Expression() {}
func (astGrouping Grouping) String() string {
	return "(" + astGrouping.Exp.String() + ")"
}
