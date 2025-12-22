package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/lexar"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

type Parser struct {
	lexar     lexar.Lexar
	currToken token.Token
	peekToken token.Token
}

func NewParser(lexar lexar.Lexar) *Parser {
	p := Parser{
		lexar: lexar,
	}
	p.next()
	p.next()
	p.registerfunc()

	return &p
}

func (p *Parser) next() {
	p.currToken = p.peekToken
	p.peekToken = p.lexar.NextToken()
}

func (p *Parser) Parse() ast.Expression {
	for p.currToken.TokenType != token.Eof {
		return p.parseExpressionTwo(Lowest)
	}
	panic("")
}

var precedence = map[token.TokenType]int{
	token.Plus:  Sum,
	token.Slash: Product,
	token.Star:  Product,
}

func getPrecedence(tokenType token.TokenType) int {
	if prec, ok := precedence[tokenType]; ok {
		return prec
	}

	return -1
}

type prefixParseFn func() ast.Expression
type infixParseFn func(left ast.Expression) ast.Expression

const (
	_ int = iota
	Lowest
	Sum
	Product
	Prefix
)

var prefix = map[token.TokenType]prefixParseFn{}
var infix = map[token.TokenType]infixParseFn{}

func (p *Parser) registerfunc() {
	prefix[token.NumberInt] = p.parseInt
	prefix[token.NumberFloat] = p.parseFloat
	prefix[token.Identifier] = p.parseIdentifier
	prefix[token.NilToken] = p.parseNil
	prefix[token.StringToken] = p.parseStringLiteral
	prefix[token.TrueToken] = p.parseBool
	prefix[token.FalseToken] = p.parseBool
	prefix[token.Minus] = p.parsePrefix
	prefix[token.Bang] = p.parsePrefix
	prefix[token.LeftParen] = p.parseGrouping

	infix[token.Plus] = p.parseInfix
	infix[token.Slash] = p.parseInfix
	infix[token.Star] = p.parseInfix
}

// func (p *Parser) parseExpressionFunc() ast.Expression {
// 	switch p.currToken.TokenType {
// 	case token.LeftParen:

// 	case token.Minus, token.Bang:
// 		return p.parsePrefixExpression()
// 	default:
// 		panic(fmt.Sprintf("unknown token type: %v", p.currToken.TokenType))
// 	}
// }

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	ast := ast.TreeExpression{
		Left:     left,
		Operator: p.currToken.Lexeme,
		Token:    p.currToken,
	}

	p.next()

	right := p.parseExpressionTwo(getPrecedence(ast.Token.TokenType))

	ast.Right = right

	return ast
}

func (p *Parser) parseExpressionTwo(predecence int) ast.Expression {
	prefixFunc, ok := prefix[p.currToken.TokenType]
	if !ok {
		return nil
	}

	left := prefixFunc()

	for predecence < getPrecedence(p.peekToken.TokenType) {
		infixFunc, ok := infix[p.peekToken.TokenType]
		if !ok {
			return nil
		}

		p.next()
		left = infixFunc(left)
	}

	return left

}

func (p *Parser) parsePrefix() ast.Expression {
	ast := ast.PrefixExpression{
		Operator: p.currToken.Lexeme,
		Token:    p.currToken,
	}
	p.next()
	right := p.parseExpressionTwo(Prefix)

	ast.Right = right

	return ast
}

func (p *Parser) parseGrouping() ast.Expression {
	p.next()
	groupingExpression := ast.GroupingExpression{
		Exp: p.parseExpressionTwo(Lowest),
	}

	if p.peekToken.TokenType != token.RightParen {
		panic(fmt.Sprintf("expecting right paren got:%v", p.currToken.TokenType))
	}

	return groupingExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.Identifier{Value: p.currToken.Lexeme, Token: p.currToken}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return ast.StringLiteral{Value: p.currToken.Literal.(string), Token: p.currToken}
}

func (p *Parser) parseInt() ast.Expression {
	return ast.Inteager{Value: p.currToken.Literal.(int64), Token: p.currToken}
}

func (p *Parser) parseFloat() ast.Expression {
	return ast.Float{Value: p.currToken.Literal.(float64), Token: p.currToken}
}

func (p *Parser) parseBool() ast.Expression {
	return ast.Boolean{Value: p.currToken.IsToken(token.TrueToken), Token: p.currToken}
}

func (p *Parser) parseNil() ast.Expression {
	return ast.Nil{}
}
