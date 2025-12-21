package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/lexar"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

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

type Parser struct {
	lexar     lexar.Lexar
	currToken token.Token
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

func (p *Parser) Parse() ast.Expression {
	for p.currToken.TokenType != token.Eof {
		return p.parseExpression()
		p.next()
	}
	panic("")
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.currToken.TokenType {
	case token.TrueToken, token.FalseToken:
		return p.parseBool()
	case token.NilToken:
		return p.parseNil()
	case token.NumberInt:
		return p.parseInt()
	case token.NumberFloat:
		return p.parseFloat()
	case token.StringToken:
		return p.parseStringLiteral()
	case token.Identifier:
		return p.parseIdentifier()
	case token.LeftParen:
		p.next()
		aa := ast.Grouping{
			Exp: p.parseExpression(),
		}
		p.next()
		if p.currToken.TokenType != token.RightParen {
			panic("expecting right paren")
		}

		return aa
	default:
		panic(fmt.Sprintf("unknown token type: %v", p.currToken.TokenType))
	}
}

func (p *Parser) parseIdentifier() ast.Identifier {
	return ast.Identifier{Value: p.currToken.Lexeme, Token: p.currToken}
}

func (p *Parser) parseStringLiteral() ast.StringLiteral {
	return ast.StringLiteral{Value: p.currToken.Literal.(string), Token: p.currToken}
}

func (p *Parser) parseInt() ast.Inteager {
	return ast.Inteager{Value: p.currToken.Literal.(int64), Token: p.currToken}
}

func (p *Parser) parseFloat() ast.Float {
	return ast.Float{Value: p.currToken.Literal.(float64), Token: p.currToken}
}

func (p *Parser) parseBool() ast.Boolean {
	return ast.Boolean{Value: p.currToken.IsToken(token.TrueToken), Token: p.currToken}
}

func (p *Parser) parseNil() ast.Nil {
	return ast.Nil{}
}
