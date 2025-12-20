package parser

import "github.com/codecrafters-io/interpreter-starter-go/app/pkg/lexar"

//
// let <identifier> = <expression>
// <expression> - produces value, statement don't
// let x = 5 doesnt product value
// wheras 5 does
// return 5 does not product value
// but add(5,5) does

// lets start with expression 2 + 3
// 2
// 2 +3 / 4
// 2/4*5-1+5

// BNF
// expression -> term ((+ | -) term)*
// term -> literals ((* | /) literals)*

// phase 1 only literals
// literals -> "true" | "false" | "nil"

type ASTExpression interface {
	expression()
	String() string
}

type ASTBoolean struct {
	value bool
	token lexar.Token
}

func (astBool ASTBoolean) expression() {}
func (astBool ASTBoolean) String() string {
	return astBool.token.Lexeme
}

type ASTNil struct{}

func (astNil ASTNil) expression() {}
func (astNil ASTNil) String() string {
	return "nil"
}

type Parser struct {
	lexar     lexar.Lexar
	currToken lexar.Token
}

func NewParser(lexar lexar.Lexar) *Parser {
	p := Parser{
		lexar: lexar,
	}
	p.next()

	return &p
}

func (p *Parser) next() {
	p.currToken = p.lexar.NextToken()
}

func (p *Parser) Parse() ASTExpression {
	for p.currToken.TokenType != lexar.Eof {
		return p.parseExpression()
		p.next()
	}
	panic("")
}

func (p *Parser) parseExpression() ASTExpression {
	switch p.currToken.TokenType {
	case lexar.TrueToken, lexar.FalseToken:
		return p.parseBool()
	case lexar.NilToken:
		return p.parseNil()
	default:
		panic("err")
	}
}

func (p *Parser) parseBool() ASTBoolean {
	return ASTBoolean{value: p.currToken.IsToken(lexar.TrueToken), token: p.currToken}
}

func (p *Parser) parseNil() ASTNil {
	return ASTNil{}
}
