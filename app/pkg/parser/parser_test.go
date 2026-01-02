package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/lexar"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

func TestParserExpression(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal ast.Expression
	}{
		{input: "true", expectedVal: ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}}},
		{input: "false", expectedVal: ast.Boolean{Value: false, Token: token.Token{TokenType: token.FalseToken, Lexeme: "false", Literal: nil}}},
		{input: "nil", expectedVal: ast.Nil{}},
		{input: "123.0", expectedVal: ast.Float{Value: 123, Token: token.Token{TokenType: token.NumberFloat, Lexeme: "123.0", Literal: 123.0}}},
		{input: "123.12", expectedVal: ast.Float{Value: 123.12, Token: token.Token{TokenType: token.NumberFloat, Lexeme: "123.12", Literal: 123.12}}},
		{input: "123", expectedVal: ast.Integer{Value: 123, Token: token.Token{TokenType: token.NumberInt, Lexeme: "123", Literal: int64(123)}}},
		{input: "\"123\"", expectedVal: ast.StringLiteral{Value: "123", Token: token.Token{TokenType: token.StringToken, Lexeme: "\"123\"", Literal: "123"}}},
		{input: "let", expectedVal: ast.Identifier{Value: "let", Token: token.Token{TokenType: token.Identifier, Lexeme: "let", Literal: nil}}},
		{input: "(let)", expectedVal: ast.GroupingExpression{Exp: ast.Identifier{Value: "let", Token: token.Token{TokenType: token.Identifier, Lexeme: "let", Literal: nil}}}},
		{input: "!true", expectedVal: ast.PrefixExpression{
			Right:    ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
			Operator: token.Bang,
			Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
		}},
		{input: "!!true", expectedVal: ast.PrefixExpression{
			Right: ast.PrefixExpression{
				Right:    ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
				Operator: token.Bang,
				Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
			},
			Operator: token.Bang,
			Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
		}},
		{input: "!!(true)", expectedVal: ast.PrefixExpression{
			Right: ast.PrefixExpression{
				Right: ast.GroupingExpression{
					Exp: ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
				},
				Operator: token.Bang,
				Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
			},
			Operator: token.Bang,
			Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
		}},
		{input: "(!!true)", expectedVal: ast.GroupingExpression{
			Exp: ast.PrefixExpression{
				Right: ast.PrefixExpression{
					Right:    ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
					Operator: token.Bang,
					Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
				},
				Operator: token.Bang,
				Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
			}},
		},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.Parse()

			respExpression := resp.Statements[0].(ast.ExpressionStatement).Expression

			if !reflect.DeepEqual(respExpression, testCase.expectedVal) {
				t.Errorf("Expected to get:\n %#v, got: \n%#v", testCase.expectedVal, resp)
			}
		})
	}
}

func TestParserExpressionArithmetic(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal ast.Expression
	}{
		{input: "123", expectedVal: ast.Integer{Value: 123, Token: token.Token{TokenType: token.NumberInt, Lexeme: "123", Literal: int64(123)}}},
		{input: "-5", expectedVal: ast.PrefixExpression{
			Operator: token.Minus,
			Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
			Right: ast.Integer{
				Value: 5,
				Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)},
			},
		}},
		{input: "--5", expectedVal: ast.PrefixExpression{
			Operator: token.Minus,
			Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
			Right: ast.PrefixExpression{
				Operator: token.Minus,
				Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
				Right: ast.Integer{
					Value: 5,
					Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)},
				},
			},
		}},
		{input: "5+6", expectedVal: ast.InfixExpression{
			Operator: token.Plus,
			Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
			Left:     ast.Integer{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
			Right:    ast.Integer{Value: 6, Token: token.Token{TokenType: token.NumberInt, Lexeme: "6", Literal: int64(6)}},
		}},

		{input: "5/2+3", expectedVal: ast.InfixExpression{
			Operator: token.Plus,
			Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
			Left: ast.InfixExpression{
				Operator: token.Slash,
				Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
				Right:    ast.Integer{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
				Left:     ast.Integer{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
			},
			Right: ast.Integer{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
		}},

		{input: "2+3/5", expectedVal: ast.InfixExpression{
			Operator: token.Plus,
			Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
			Left:     ast.Integer{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
			Right: ast.InfixExpression{
				Operator: token.Slash,
				Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
				Left:     ast.Integer{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
				Right:    ast.Integer{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
			},
		}},
		{input: "5*3/2", expectedVal: ast.InfixExpression{
			Operator: token.Slash,
			Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
			Right:    ast.Integer{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
			Left: ast.InfixExpression{
				Operator: token.Star,
				Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
				Left:     ast.Integer{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
				Right:    ast.Integer{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
			},
		}},

		{input: "51*36/71", expectedVal: ast.InfixExpression{
			Operator: token.Slash,
			Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
			Right:    ast.Integer{Value: 71, Token: token.Token{TokenType: token.NumberInt, Lexeme: "71", Literal: int64(71)}},
			Left: ast.InfixExpression{
				Operator: token.Star,
				Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
				Left:     ast.Integer{Value: 51, Token: token.Token{TokenType: token.NumberInt, Lexeme: "51", Literal: int64(51)}},
				Right:    ast.Integer{Value: 36, Token: token.Token{TokenType: token.NumberInt, Lexeme: "36", Literal: int64(36)}},
			},
		}},
		{input: "(37 * -21 / (45 * 16))", expectedVal: ast.GroupingExpression{
			Exp: ast.InfixExpression{
				Operator: token.Slash,
				Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
				Right: ast.GroupingExpression{
					Exp: ast.InfixExpression{
						Operator: token.Star,
						Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
						Left:     ast.Integer{Value: 45, Token: token.Token{TokenType: token.NumberInt, Lexeme: "45", Literal: int64(45)}},
						Right:    ast.Integer{Value: 16, Token: token.Token{TokenType: token.NumberInt, Lexeme: "16", Literal: int64(16)}},
					},
				},
				Left: ast.InfixExpression{
					Operator: token.Star,
					Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
					Left:     ast.Integer{Value: 37, Token: token.Token{TokenType: token.NumberInt, Lexeme: "37", Literal: int64(37)}},
					Right: ast.PrefixExpression{
						Operator: token.Minus,
						Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
						Right:    ast.Integer{Value: 21, Token: token.Token{TokenType: token.NumberInt, Lexeme: "21", Literal: int64(21)}},
					},
				},
			},
		}},
		{input: "1+2*3+4", expectedVal: ast.InfixExpression{
			Operator: token.Plus,
			Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
			Right:    ast.Integer{Value: 4, Token: token.Token{TokenType: token.NumberInt, Lexeme: "4", Literal: int64(4)}},
			Left: ast.InfixExpression{
				Operator: token.Plus,
				Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
				Left:     ast.Integer{Value: 1, Token: token.Token{TokenType: token.NumberInt, Lexeme: "1", Literal: int64(1)}},
				Right: ast.InfixExpression{
					Operator: token.Star,
					Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
					Right:    ast.Integer{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
					Left:     ast.Integer{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
				},
			},
		},
		},
		{input: "52+80-94", expectedVal: ast.InfixExpression{
			Operator: token.Minus,
			Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
			Right:    ast.Integer{Value: 94, Token: token.Token{TokenType: token.NumberInt, Lexeme: "94", Literal: int64(94)}},
			Left: ast.InfixExpression{
				Operator: token.Plus,
				Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
				Left:     ast.Integer{Value: 52, Token: token.Token{TokenType: token.NumberInt, Lexeme: "52", Literal: int64(52)}},
				Right:    ast.Integer{Value: 80, Token: token.Token{TokenType: token.NumberInt, Lexeme: "80", Literal: int64(80)}},
			},
		}},
		{input: "(1+2)*(3+4)", expectedVal: ast.InfixExpression{
			Operator: token.Star,
			Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
			Left: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,
					Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
					Left:     ast.Integer{Value: 1, Token: token.Token{TokenType: token.NumberInt, Lexeme: "1", Literal: int64(1)}},
					Right:    ast.Integer{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
				},
			},
			Right: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,
					Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
					Left:     ast.Integer{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
					Right:    ast.Integer{Value: 4, Token: token.Token{TokenType: token.NumberInt, Lexeme: "4", Literal: int64(4)}},
				},
			},
		}},
		{input: "(1+2)>(3+4)", expectedVal: ast.InfixExpression{
			Operator: token.Greater,
			Token:    token.Token{TokenType: token.Greater, Lexeme: ">"},
			Left: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,
					Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
					Left:     ast.Integer{Value: 1, Token: token.Token{TokenType: token.NumberInt, Lexeme: "1", Literal: int64(1)}},
					Right:    ast.Integer{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
				},
			},
			Right: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,
					Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
					Left:     ast.Integer{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
					Right:    ast.Integer{Value: 4, Token: token.Token{TokenType: token.NumberInt, Lexeme: "4", Literal: int64(4)}},
				},
			},
		}},
		{input: "22 == 20", expectedVal: ast.InfixExpression{
			Operator: token.EqualEqual,
			Token:    token.Token{TokenType: token.EqualEqual, Lexeme: "=="},
			Right:    ast.Integer{Value: 20, Token: token.Token{TokenType: token.NumberInt, Lexeme: "20", Literal: int64(20)}},
			Left:     ast.Integer{Value: 22, Token: token.Token{TokenType: token.NumberInt, Lexeme: "22", Literal: int64(22)}},
		}},
		{input: "22 != 20", expectedVal: ast.InfixExpression{
			Operator: token.BangEqual,
			Token:    token.Token{TokenType: token.BangEqual, Lexeme: "!="},
			Right:    ast.Integer{Value: 20, Token: token.Token{TokenType: token.NumberInt, Lexeme: "20", Literal: int64(20)}},
			Left:     ast.Integer{Value: 22, Token: token.Token{TokenType: token.NumberInt, Lexeme: "22", Literal: int64(22)}},
		}},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.parseExpression(Lowest)

			if !reflect.DeepEqual(resp, testCase.expectedVal) {
				t.Errorf("Expected to \nget: %#v, \ngot: %#v", testCase.expectedVal, resp)
			}
		})
	}
}
func TestParserDeclarationStatment(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal *ast.Program
	}{
		{
			input: "var aa = 53;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.DeclarationStatement{
						Names: []string{"aa"},
						Expression: ast.Integer{
							Value: 53,
							Token: token.Token{TokenType: token.NumberInt, Lexeme: "53", Literal: int64(53)},
						},
					},
				},
			},
		},
		{
			input: "var aa;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.DeclarationStatement{
						Names:      []string{"aa"},
						Expression: ast.Nil{},
					},
				},
			},
		},
		{
			input: "var aa =bb =\"hello world\";",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.DeclarationStatement{
						Names:      []string{"aa", "bb"},
						Expression: ast.StringLiteral{Value: "hello world", Token: token.Token{Lexeme: "\"hello world\"", Literal: "hello world", TokenType: token.StringToken}},
					},
				},
			},
		},
		{
			input: "var bar =68;var world=bar;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.DeclarationStatement{
						Names:      []string{"bar"},
						Expression: ast.Integer{Value: 68, Token: token.Token{Lexeme: "68", Literal: int64(68), TokenType: token.NumberInt}},
					},
					ast.DeclarationStatement{
						Names:      []string{"world"},
						Expression: ast.Identifier{Value: "bar", Token: token.Token{Lexeme: "bar", TokenType: token.Identifier}},
					},
				},
			},
		},
		{
			input: "var bar =68;var world=11;bar=world",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.DeclarationStatement{
						Names:      []string{"bar"},
						Expression: ast.Integer{Value: 68, Token: token.Token{Lexeme: "68", Literal: int64(68), TokenType: token.NumberInt}},
					},
					ast.DeclarationStatement{
						Names:      []string{"world"},
						Expression: ast.Integer{Value: 11, Token: token.Token{Lexeme: "11", Literal: int64(11), TokenType: token.NumberInt}},
					},
					ast.ExpressionStatement{
						Expression: ast.AssignExpression{IdentifierName: "bar", Value: ast.Identifier{Value: "world", Token: token.Token{Lexeme: "world", TokenType: token.Identifier}}},
					},
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.Parse()

			if !reflect.DeepEqual(resp, testCase.expectedVal) {
				t.Errorf("Expected to \nget: %#v, \ngot: %#v", testCase.expectedVal, resp)
			}
		})
	}
}

func TestParserExpressionFunctions(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal ast.Expression
	}{
		{input: "print(\"aa\")", expectedVal: ast.CallExpression{
			Function:  ast.Identifier{Value: "print", Token: token.Token{TokenType: token.Identifier, Lexeme: "print"}},
			Arguments: []ast.Expression{ast.GroupingExpression{Exp: ast.StringLiteral{Value: "aa", Token: token.Token{TokenType: token.StringToken, Lexeme: "\"aa\"", Literal: "aa"}}}},
		}},
		{input: "print \"aa\"", expectedVal: ast.CallExpression{
			Function:  ast.Identifier{Value: "print", Token: token.Token{TokenType: token.Identifier, Lexeme: "print"}},
			Arguments: []ast.Expression{ast.StringLiteral{Value: "aa", Token: token.Token{TokenType: token.StringToken, Lexeme: "\"aa\"", Literal: "aa"}}},
		}},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.parseExpression(Lowest)

			if !reflect.DeepEqual(resp, testCase.expectedVal) {
				t.Errorf("Expected to \nget: %#v, \ngot: %#v", testCase.expectedVal, resp)
			}
		})
	}
}

func TestExpressionErrors(t *testing.T) {

	tests := []struct {
		input  string
		errors []string
	}{
		{input: "(72 +)", errors: []string{"[line 1] Error at ')': Expect expression.", "[line 1] expected next token to be RIGHT_PAREN, got EOF instead"}},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			parser.parseExpression(Lowest)

			if !reflect.DeepEqual(parser.Errors(), testCase.errors) {
				t.Errorf("Expected to \nget: %#v, \ngot: %#v", testCase.errors, parser.Errors())
			}
		})
	}
}
