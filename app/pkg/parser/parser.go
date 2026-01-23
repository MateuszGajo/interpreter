package parser

import (
	"fmt"
	"reflect"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/lexar"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

type Parser struct {
	lexar     lexar.Lexar
	currToken token.Token
	peekToken token.Token
	errors    []string
}

func NewParser(lexar lexar.Lexar) *Parser {
	p := Parser{
		lexar:  lexar,
		errors: []string{},
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

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("[line %v] %v", p.lexar.Line, msg))
}

func (p *Parser) noPrefixParseError(t token.Token) {
	msg := fmt.Sprintf("Error at '%v': Expect expression.", t.Lexeme)
	p.addError(msg)
}

func (p *Parser) peekError(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tokenType, p.peekToken.TokenType)
	p.addError(msg)
}

func (p *Parser) expectError(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected token to be %s, got %s instead", tokenType, p.currToken.TokenType)
	p.addError(msg)
}

func (p *Parser) expect(tokenType token.TokenType) bool {
	if p.currToken.TokenType == tokenType {
		return true
	}

	p.expectError(tokenType)

	return false
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekToken.TokenType == tokenType {
		p.next()
		return true
	}

	p.peekError(tokenType)

	return false
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currToken.TokenType != token.Eof {
		stmt := p.parseStatement()
		if stmt == nil {
			return program
		}
		program.Statements = append(program.Statements, stmt)
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	var stmt ast.Statement
	switch p.currToken.TokenType {
	case token.LeftCurlyBrace:
		return p.parseBlockStatement()
	case token.IfToken:
		return p.parseIfStatement()
	case token.WhileToken:
		return p.parseWhileStatement()
	case token.ForToken:
		return p.parseForStatement()
	case token.FunToken:
		return p.parseFunctionStatement()
	case token.VarToken:
		stmt = p.parseDeclarationStatement()
	case token.ReturnToken:
		stmt = p.parseReturnStatement()
	default:
		stmt = p.parseExpressionStatement()
	}

	if !p.expect(token.Semicolon) {
		return nil
	}
	p.next()

	return stmt
}

func (p *Parser) parseReturnStatement() ast.ReturnStatement {
	returnStatement := ast.ReturnStatement{}
	p.next()

	returnStatement.Expression = p.parseExpression(Lowest)
	if p.peekToken.TokenType == token.Semicolon {
		p.next()
	}

	return returnStatement
}

func (p *Parser) parseFunctionStatement() ast.FunctionStatement {
	function := ast.FunctionStatement{}
	p.next()

	callExp := p.parseExpression(Lowest)

	c, ok := callExp.(ast.CallExpression)

	if !ok {
		panic("expected call expression")
	}

	function.Name = c.Function.(ast.Identifier).Value

	for _, arg := range c.Arguments {
		identifier, ok := arg.(ast.Identifier)
		if !ok {
			panic(fmt.Sprintf("expected arg to be identifier, got: %v", reflect.TypeOf(arg)))
		}
		function.Args = append(function.Args, &identifier)
	}

	if !p.expectPeek(token.LeftCurlyBrace) {
		return function
	}
	function.Block = p.parseBlockStatement().(ast.BlockStatement)

	return function
}

func (p *Parser) parseForStatement() ast.ForStatement {
	forStatement := ast.ForStatement{}

	if !p.expectPeek(token.LeftParen) {
		return forStatement
	}

	p.next()

	if p.currToken.TokenType == token.Semicolon {
		forStatement.Declaration = nil
		if !p.expect(token.Semicolon) {
			return forStatement
		}
		p.next()
	} else {
		declaration := p.parseStatement()
		forStatement.Declaration = &declaration

	}

	forStatement.Condition = p.parseExpression(Lowest)

	if !p.expectPeek(token.Semicolon) {
		return forStatement
	}

	p.next()

	if p.currToken.TokenType != token.RightParen {
		exp := p.parseExpression(Lowest)
		forStatement.Evaluator = &exp
		p.next()
	}

	if !p.expect(token.RightParen) {
		return forStatement
	}

	p.next()

	forStatement.Block = p.parseBranchBlock()

	return forStatement
}

func (p *Parser) parseWhileStatement() ast.WhileStatement {
	while := ast.WhileStatement{}
	if !p.expectPeek(token.LeftParen) {
		return while
	}

	while.Condition = p.parseExpression(Lowest)

	if !p.expect(token.RightParen) {
		return while
	}
	p.next()

	while.Block = p.parseBranchBlock()

	return while
}

func (p *Parser) parseIfStatement() ast.IfStatement {

	ifStatement := ast.IfStatement{}

	if !p.expectPeek(token.LeftParen) {
		return ifStatement
	}

	ifStatement.Condition = p.parseExpression(Lowest)

	if !p.expect(token.RightParen) {
		return ifStatement
	}

	p.next()

	ifStatement.Then = p.parseBranchBlock()
	if p.currToken.TokenType == token.ElseToken {
		p.next()
		elseStmt := p.parseBranchBlock()
		ifStatement.Else = &elseStmt
	}

	return ifStatement
}

func (p *Parser) parseDeclarationStatement() ast.DeclarationStatement {
	p.expectPeek(token.Identifier)
	names := []string{p.currToken.Lexeme}

	if p.peekToken.TokenType == token.Semicolon {
		p.next()
		return ast.DeclarationStatement{Expression: ast.Nil{}, Names: names}
	}

	p.expectPeek(token.Equal)
	// var aa = 65
	// var aa = bb = cc = 65
	// var bb = 77 + 100;
	// var bb = aa + 10;
	p.next()

	for p.currToken.TokenType == token.Identifier && p.peekToken.TokenType == token.Equal {
		name := p.currToken.Lexeme
		names = append(names, name)
		p.next()
		p.next()
	}

	exp := p.parseExpression(Lowest)
	if p.peekToken.TokenType == token.Semicolon {
		p.next()
	}
	return ast.DeclarationStatement{Expression: exp, Names: names}
}

func (p *Parser) parseExpressionStatement() ast.ExpressionStatement {
	exp := p.parseExpression(Lowest)

	if p.peekToken.TokenType == token.Semicolon {
		p.next()
	}

	return ast.ExpressionStatement{Expression: exp}
}

var precedence = map[token.TokenType]int{
	token.Less:         LessOrGreater,
	token.LessEqual:    LessOrGreater,
	token.Greater:      LessOrGreater,
	token.GreaterEqual: LessOrGreater,
	token.Equal:        Assign,
	token.EqualEqual:   Equals,
	token.BangEqual:    Equals,
	token.Plus:         Sum,
	token.Minus:        Sum,
	token.Slash:        Product,
	token.Star:         Product,
	token.LeftParen:    Call,
	token.OrToken:      Logicals,
	token.AndToken:     Logicals,
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
	Assign        // =
	Logicals      // and or
	Equals        // ==
	LessOrGreater // < >
	Sum           // +, -
	Product       // *, /
	Call          // myFunction(x)
	Prefix        // -X !X
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
	prefix[token.Identifier] = p.parseIdentifier
	prefix[token.PrintToken] = p.parsePrintExpression

	infix[token.Plus] = p.parseInfix
	infix[token.Slash] = p.parseInfix
	infix[token.Star] = p.parseInfix
	infix[token.Minus] = p.parseInfix
	infix[token.Greater] = p.parseInfix
	infix[token.GreaterEqual] = p.parseInfix
	infix[token.LessEqual] = p.parseInfix
	infix[token.Less] = p.parseInfix
	infix[token.Equal] = p.parseAssign
	infix[token.EqualEqual] = p.parseInfix
	infix[token.BangEqual] = p.parseInfix
	infix[token.LeftParen] = p.parseCallExpression
	infix[token.OrToken] = p.parseInfix
	infix[token.AndToken] = p.parseInfix
}

// {
// var hello = "baz";
// print hello;
//  }

func (p *Parser) parseBlockStatement() ast.Statement {
	statements := []ast.Statement{}
	p.next()
	for p.currToken.TokenType != token.RightCurlyBrace && p.currToken.TokenType != token.Eof {
		stmt := p.parseStatement()
		if stmt == nil {
			return nil
		}
		statements = append(statements, stmt)
	}

	if !p.expect(token.RightCurlyBrace) {
		return nil
	}
	p.next()

	return ast.BlockStatement{Statements: statements}
}

func (p *Parser) parseBranchBlock() ast.BlockStatement {
	statement := p.parseStatement()
	if v, ok := statement.(ast.BlockStatement); ok {
		return v
	}
	return ast.BlockStatement{Statements: []ast.Statement{statement}}

}

func (p *Parser) parseAssign(left ast.Expression) ast.Expression {
	identifier := left.(ast.Identifier)

	expr := ast.AssignExpression{
		IdentifierName: identifier.Value,
	}
	p.next()

	expr.Value = p.parseExpression(Assign - 1)

	return expr

}

func (p *Parser) parsePrintExpression() ast.Expression {
	p.next()
	return ast.CallExpression{Function: ast.Identifier{Value: "print", Token: token.Token{TokenType: token.Identifier, Lexeme: "print"}}, Arguments: p.parseExpressionList()}
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExp := ast.CallExpression{Function: function}
	callExp.Arguments = p.parseExpressionList()
	return callExp
}

// func (p *Parser) parseGrouping() ast.Expression {
// 	if p.peekToken.TokenType == token.RightParen {
// 		return nil
// 	}
// 	p.next()

// 	groupingExpression := ast.GroupingExpression{
// 		Exp: p.parseExpression(Lowest),
// 	}

// 	if !p.expectPeek(token.RightParen) {
// 		return nil
// 	}

// 	return groupingExpression
// }

func (p *Parser) parseExpressionList() []ast.Expression {
	list := []ast.Expression{}

	if p.currToken.TokenType == token.LeftParen {
		p.next()
	}
	if p.currToken.TokenType == token.RightParen {
		return list
	}
	exp := p.parseExpression(Lowest)
	// if v, ok := exp.(ast.GroupingExpression); ok {
	// 	return v.Exp
	// }
	if exp == nil {
		if p.peekToken.TokenType == token.RightParen {
			p.next()
		}
		return list
	}
	list = append(list, exp)
	p.next()

	for p.currToken.TokenType == token.Comma {
		p.next()
		exp := p.parseExpression(Lowest)
		list = append(list, exp)
		p.next()
	}

	return list
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	ast := ast.InfixExpression{
		Left:     left,
		Operator: p.currToken.TokenType,
		Token:    p.currToken,
	}

	p.next()

	right := p.parseExpression(getPrecedence(ast.Token.TokenType))

	ast.Right = right

	return ast
}

// Prat precedence
// Start at leaves: Parse numbers/unary ops first as base expressions
// Climb precedence ladder: While next operator matches or exceeds current level:
// Higher precedence: Builds tight subtrees (like * inside +)
// Same precedence: Chains left-to-right (left spine: (1*2)*3)
// Drop stops building: Lower precedence operator halts current level, passes complete subtree up

// 1*2*3 → (((1*2)*3)
//      *
//     / \
//    *   3
//   / \
//  1   2

// 1+2*3 → (1+(2*3))
//      +
//     / \
//    1  *
//      /\
//     2  3

// 1*2/3*4 → (((1*2)/3)*4)
//        *
//       / \
//     "/"   4
//   	/ \
//     *   3
// 	  / \
//    1  2

// 52 +80 -94
func (p *Parser) parseExpression(predecence int) ast.Expression {
	prefixFunc, ok := prefix[p.currToken.TokenType]
	if !ok {
		p.noPrefixParseError(p.currToken)
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
		Operator: p.currToken.TokenType,
		Token:    p.currToken,
	}
	p.next()
	right := p.parseExpression(Prefix)

	ast.Right = right

	return ast
}

// func (p *Parser) parseExpressionList() []ast.Expression {
// 	list := []ast.Expression{}

// 	exp := p.parseExpression(Lowest)
// 	if v, ok := exp.(ast.GroupingExpression); ok {
// 		return v.Exp
// 	}
// 	if exp == nil {
// 		if p.peekToken.TokenType == token.RightParen {
// 			p.next()
// 		}
// 		return list
// 	}
// 	list = append(list, exp)
// 	p.next()

// 	for p.currToken.TokenType == token.Comma {
// 		exp := p.parseExpression(Lowest)
// 		list = append(list, exp)
// 	}

func (p *Parser) parseGrouping() ast.Expression {
	grouping := ast.GroupingExpression{Exp: []ast.Expression{}}
	if p.peekToken.TokenType == token.RightParen {
		return nil
	}
	p.next()

	grouping.Exp = append(grouping.Exp, p.parseExpression(Lowest))

	p.next()
	for p.currToken.TokenType == token.Comma {
		p.next()
		exp := p.parseExpression(Lowest)
		grouping.Exp = append(grouping.Exp, exp)
		p.next()
	}

	if !p.expect(token.RightParen) {
		return nil
	}

	return grouping
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.Identifier{Value: p.currToken.Lexeme, Token: p.currToken}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return ast.StringLiteral{Value: p.currToken.Literal.(string), Token: p.currToken}
}

func (p *Parser) parseInt() ast.Expression {
	return ast.Integer{Value: p.currToken.Literal.(int64), Token: p.currToken}
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
