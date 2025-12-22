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
		{input: "123", expectedVal: ast.Inteager{Value: 123, Token: token.Token{TokenType: token.NumberInt, Lexeme: "123", Literal: int64(123)}}},
		{input: "\"123\"", expectedVal: ast.StringLiteral{Value: "123", Token: token.Token{TokenType: token.StringToken, Lexeme: "\"123\"", Literal: "123"}}},
		{input: "let", expectedVal: ast.Identifier{Value: "let", Token: token.Token{TokenType: token.Identifier, Lexeme: "let", Literal: nil}}},
		{input: "(let)", expectedVal: ast.GroupingExpression{Exp: ast.Identifier{Value: "let", Token: token.Token{TokenType: token.Identifier, Lexeme: "let", Literal: nil}}}},
		{input: "!true", expectedVal: ast.PrefixExpression{
			Right:    ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
			Operator: "!",
			Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
		}},
		{input: "!!true", expectedVal: ast.PrefixExpression{
			Right: ast.PrefixExpression{
				Right:    ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
				Operator: "!",
				Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
			},
			Operator: "!",
			Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
		}},
		{input: "!!(true)", expectedVal: ast.PrefixExpression{
			Right: ast.PrefixExpression{
				Right: ast.GroupingExpression{
					Exp: ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
				},
				Operator: "!",
				Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
			},
			Operator: "!",
			Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
		}},
		{input: "(!!true)", expectedVal: ast.GroupingExpression{
			Exp: ast.PrefixExpression{
				Right: ast.PrefixExpression{
					Right:    ast.Boolean{Value: true, Token: token.Token{TokenType: token.TrueToken, Lexeme: "true", Literal: nil}},
					Operator: "!",
					Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
				},
				Operator: "!",
				Token:    token.Token{TokenType: token.Bang, Lexeme: "!"},
			}},
		},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.Parse()

			if !reflect.DeepEqual(resp, testCase.expectedVal) {
				t.Errorf("Expected to get:\n %#v, got: \n%#v", testCase.expectedVal, resp)
			}
		})
	}
}

func TestParserExpression2(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal ast.Expression
	}{
		// {input: "123", expectedVal: ast.Inteager{Value: 123, Token: token.Token{TokenType: token.NumberInt, Lexeme: "123", Literal: int64(123)}}},
		// {input: "-5", expectedVal: ast.PrefixExpression{
		// 	Operator: "-",
		// 	Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
		// 	Right: ast.Inteager{
		// 		Value: 5,
		// 		Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)},
		// 	},
		// }},
		// {input: "--5", expectedVal: ast.PrefixExpression{
		// 	Operator: "-",
		// 	Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
		// 	Right: ast.PrefixExpression{
		// 		Operator: "-",
		// 		Token:    token.Token{TokenType: token.Minus, Lexeme: "-"},
		// 		Right: ast.Inteager{
		// 			Value: 5,
		// 			Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)},
		// 		},
		// 	},
		// }},
		// {input: "5+6", expectedVal: ast.TreeExpression{
		// 	Operator: "+",
		// 	Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
		// 	Left:     ast.Inteager{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
		// 	Right:    ast.Inteager{Value: 6, Token: token.Token{TokenType: token.NumberInt, Lexeme: "6", Literal: int64(6)}},
		// }},

		// {input: "5/2+3", expectedVal: ast.TreeExpression{
		// 	Operator: "+",
		// 	Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
		// 	Left: ast.TreeExpression{
		// 		Operator: "/",
		// 		Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
		// 		Right:    ast.Inteager{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
		// 		Left:     ast.Inteager{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
		// 	},
		// 	Right: ast.Inteager{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
		// }},

		// {input: "2+3/5", expectedVal: ast.TreeExpression{
		// 	Operator: "+",
		// 	Token:    token.Token{TokenType: token.Plus, Lexeme: "+"},
		// 	Left:     ast.Inteager{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
		// 	Right: ast.TreeExpression{
		// 		Operator: "/",
		// 		Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
		// 		Left:     ast.Inteager{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
		// 		Right:    ast.Inteager{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
		// 	},
		// }},
		// {input: "5*3/2", expectedVal: ast.TreeExpression{
		// 	Operator: "/",
		// 	Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
		// 	Right:    ast.Inteager{Value: 2, Token: token.Token{TokenType: token.NumberInt, Lexeme: "2", Literal: int64(2)}},
		// 	Left: ast.TreeExpression{
		// 		Operator: "*",
		// 		Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
		// 		Left:     ast.Inteager{Value: 5, Token: token.Token{TokenType: token.NumberInt, Lexeme: "5", Literal: int64(5)}},
		// 		Right:    ast.Inteager{Value: 3, Token: token.Token{TokenType: token.NumberInt, Lexeme: "3", Literal: int64(3)}},
		// 	},
		// }},

		// {input: "51*36/71", expectedVal: ast.TreeExpression{
		// 	Operator: "/",
		// 	Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
		// 	Right:    ast.Inteager{Value: 71, Token: token.Token{TokenType: token.NumberInt, Lexeme: "71", Literal: int64(71)}},
		// 	Left: ast.TreeExpression{
		// 		Operator: "*",
		// 		Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
		// 		Left:     ast.Inteager{Value: 51, Token: token.Token{TokenType: token.NumberInt, Lexeme: "51", Literal: int64(51)}},
		// 		Right:    ast.Inteager{Value: 36, Token: token.Token{TokenType: token.NumberInt, Lexeme: "36", Literal: int64(36)}},
		// 	},
		// }},

		{input: "(37 * -21 / (45 * 16))", expectedVal: ast.TreeExpression{
			Operator: "/",
			Token:    token.Token{TokenType: token.Slash, Lexeme: "/"},
			Right:    ast.Inteager{Value: 71, Token: token.Token{TokenType: token.NumberInt, Lexeme: "71", Literal: int64(71)}},
			Left: ast.TreeExpression{
				Operator: "*",
				Token:    token.Token{TokenType: token.Star, Lexeme: "*"},
				Left:     ast.Inteager{Value: 51, Token: token.Token{TokenType: token.NumberInt, Lexeme: "51", Literal: int64(51)}},
				Right:    ast.Inteager{Value: 36, Token: token.Token{TokenType: token.NumberInt, Lexeme: "36", Literal: int64(36)}},
			},
		}},

		// 51 * 36 / 72
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.parseExpressionTwo(Lowest)

			if !reflect.DeepEqual(resp, testCase.expectedVal) {
				t.Errorf("Expected to \nget: %#v, \ngot: %#v", testCase.expectedVal, resp)
			}
		})
	}
}
