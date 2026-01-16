package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/lexar"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

func compareStatement(outputStatement ast.Statement, expectedStatement ast.Statement) error {
	if reflect.TypeOf(outputStatement) != reflect.TypeOf(expectedStatement) {
		return fmt.Errorf("types are diffrent, expecting statement type: %v, output statement type: %v", reflect.TypeOf(expectedStatement), reflect.TypeOf(outputStatement))
	}

	switch v := outputStatement.(type) {
	case ast.BlockStatement:
		return compareBlockStatement(v, expectedStatement.(ast.BlockStatement))
	case ast.ExpressionStatement:
		return compareExpression(v.Expression, expectedStatement.(ast.ExpressionStatement).Expression)
	case ast.DeclarationStatement:
		return compareDeclarationStatement(v, expectedStatement.(ast.DeclarationStatement))
	case ast.IfStatement:
		return compareIfStatement(v, expectedStatement.(ast.IfStatement))
	case ast.WhileStatement:
		return compareWhileStatement(v, expectedStatement.(ast.WhileStatement))
	case ast.ForStatement:
		return compareForStatement(v, expectedStatement.(ast.ForStatement))
	default:
		panic(fmt.Sprintf("unknown statement type: %v", reflect.TypeOf(outputStatement)))
	}
}

func compareForStatement(outputStatement ast.ForStatement, expectedStatement ast.ForStatement) error {
	err := compareExpression(outputStatement.Condition, expectedStatement.Condition)
	if err != nil {
		return fmt.Errorf("for statement condition err \n%v", err)
	}
	err = compareBlockStatement(outputStatement.Block, expectedStatement.Block)
	if err != nil {
		return fmt.Errorf("for statement block err \n%v", err)
	}
	if outputStatement.Declaration != nil {
		err = compareStatement(*outputStatement.Declaration, *expectedStatement.Declaration)
		if err != nil {
			return fmt.Errorf("for declaration block err \n%v", err)
		}
	}

	if outputStatement.Evaluator != nil {
		err = compareExpression(*outputStatement.Evaluator, *expectedStatement.Evaluator)
		if err != nil {
			return fmt.Errorf("for evaluator block err \n%v", err)
		}
	}
	return nil
}

func compareWhileStatement(outputStatement ast.WhileStatement, expectedStatement ast.WhileStatement) error {
	err := compareExpression(outputStatement.Condition, expectedStatement.Condition)
	if err != nil {
		return fmt.Errorf("While statement condition err \n%v", err)
	}
	err = compareBlockStatement(outputStatement.Block, expectedStatement.Block)
	if err != nil {
		return fmt.Errorf("While statement block err \n%v", err)
	}
	return nil
}

func compareIfStatement(outputStatement ast.IfStatement, expectedStatement ast.IfStatement) error {
	err := compareExpression(outputStatement.Condition, expectedStatement.Condition)
	if err != nil {
		return fmt.Errorf("If statement condition err \n%v", err)
	}
	err = compareBlockStatement(outputStatement.Then, expectedStatement.Then)
	if err != nil {
		return fmt.Errorf("If statement then err \n%v", err)
	}
	if outputStatement.Else != nil {
		err = compareBlockStatement(*outputStatement.Else, *expectedStatement.Else)
		if err != nil {
			return fmt.Errorf("If statement else err \n%v", err)
		}
	}

	return nil
}

func compareDeclarationStatement(outputStatement ast.DeclarationStatement, expectedStatement ast.DeclarationStatement) error {
	if !reflect.DeepEqual(expectedStatement.Names, outputStatement.Names) {
		return fmt.Errorf("Declaration stement expected names: %v, output names: %v", expectedStatement.Names, outputStatement.Names)
	}

	err := compareExpression(outputStatement.Expression, expectedStatement.Expression)

	if err != nil {
		return fmt.Errorf("declaration statement err: %v", err)
	}

	return nil
}

func compareExpression(outputExpression ast.Expression, expectedExpression ast.Expression) error {
	if reflect.TypeOf(outputExpression) != reflect.TypeOf(expectedExpression) {
		return fmt.Errorf("Types are diffrent, expecting expression type: %v, ouput expression type: %v", reflect.TypeOf(expectedExpression), reflect.TypeOf(outputExpression))
	}

	switch v := outputExpression.(type) {
	case ast.CallExpression:
		expected := expectedExpression.(ast.CallExpression)
		for i, arg := range v.Arguments {
			err := compareExpression(arg, expected.Arguments[i])
			if err != nil {
				return fmt.Errorf("CallExpression argument index: %v failed \n %v", i, err)
			}
		}
		err := compareExpression(v.Function, expected.Function)
		if err != nil {
			return fmt.Errorf("callExpression function err: %v", err)

		}
	case ast.InfixExpression:
		expected := expectedExpression.(ast.InfixExpression)
		err := compareExpression(expected.Left, v.Left)
		if err != nil {
			return fmt.Errorf("infixExpression left err: %v", err)

		}
		err = compareExpression(expected.Right, v.Right)
		if err != nil {
			return fmt.Errorf("infixExpression right err: %v", err)

		}
		// err = compareToken(expected.Token, v.Token)
		// if err != nil {
		// 	return fmt.Errorf("Infixexpression token err: %v", err)
		// }
		if expected.Operator != v.Operator {
			return fmt.Errorf("Expected operator to be: %v output operator: %v", expected.Operator, v.Operator)
		}
	case ast.PrefixExpression:
		expected := expectedExpression.(ast.PrefixExpression)
		err := compareExpression(expected.Right, v.Right)
		if err != nil {
			return fmt.Errorf("infixExpression right err: %v", err)

		}
		// err = compareToken(expected.Token, v.Token)
		// if err != nil {
		// 	return fmt.Errorf("Infixexpression token err: %v", err)

		// }
		if expected.Operator != v.Operator {
			return fmt.Errorf("Expected operator to be: %v output operator: %v", expected.Operator, v.Operator)
		}
	case ast.Integer:
		expected := expectedExpression.(ast.Integer)
		if !reflect.DeepEqual(expected.Value, v.Value) {
			return fmt.Errorf("Integer value are diffrent, expected type: %v, number: %v, output type: %v, number: %v", reflect.TypeOf(expected.Value), expected.Value, reflect.TypeOf(v.Value), v.Value)
		}
		// err := compareToken(v.Token, expected.Token)

		// if err != nil {
		// 	return fmt.Errorf("integer err: %v", err)
		// }
	case ast.Identifier:
		expected := expectedExpression.(ast.Identifier)
		if expected.Value != v.Value {
			return fmt.Errorf("Identifier value are diffrent, expected value: %v, output value: %v", expected.Value, v.Value)
		}

		// err := compareToken(v.Token, expected.Token)
		// if err != nil {
		// 	return fmt.Errorf("identifier err: %v", err)
		// }
	case ast.StringLiteral:
		expected := expectedExpression.(ast.StringLiteral)
		if expected.Value != v.Value {
			return fmt.Errorf("string literal value are diffrent, expected value: %v, output value: %v", expected.Value, v.Value)
		}

		// err := compareToken(v.Token, expected.Token)
		// if err != nil {
		// 	return fmt.Errorf("String literal err: %v", err)
		// }
	case ast.AssignExpression:
		expected := expectedExpression.(ast.AssignExpression)
		if v.IdentifierName != expected.IdentifierName {
			return fmt.Errorf("assign expression, expected identifier name: %v, output identifier name: %v", expected.IdentifierName, v.IdentifierName)
		}

		err := compareExpression(v.Value, expected.Value)
		if err != nil {
			return fmt.Errorf("AssignExpression err \n%v", err)
		}
	case ast.Boolean:
		expected := expectedExpression.(ast.Boolean)
		if v.Value != expected.Value {
			return fmt.Errorf("boolean expression, expected val: %v, output val: %v", expected.Value, v.Value)
		}

		// err := compareToken(v.Token, expected.Token)
		// if err != nil {
		// 	return fmt.Errorf("Boolean expression err: %v", err)
		// }
	case ast.Float:
		expected := expectedExpression.(ast.Float)

		if v.Value != expected.Value {
			return fmt.Errorf("Float expression, expected val: %v, output val: %v", expected.Value, v.Value)
		}

		// err := compareToken(v.Token, expected.Token)
		// if err != nil {
		// 	return fmt.Errorf("Float expression err: %v", err)
		// }
	case ast.GroupingExpression:
		expected := expectedExpression.(ast.GroupingExpression)
		return compareExpression(v.Exp, expected.Exp)
	case ast.Nil:
		return nil
	default:
		panic(fmt.Sprintf("unknown expression type: %v", reflect.TypeOf(outputExpression)))

	}
	return nil
}

// func compareToken(outputToken token.Token, expectedToken token.Token) error {
// 	if outputToken.Lexeme != expectedToken.Lexeme {
// 		return fmt.Errorf("Token lexeme are diffrent expecting lexeme: %v, output lexeme: %v", expectedToken.Lexeme, outputToken.Lexeme)
// 	}

// 	if outputToken.TokenType != expectedToken.TokenType {
// 		return fmt.Errorf("Token type are diffrent expecting type: %v, output type: %v", expectedToken.TokenType, outputToken.TokenType)
// 	}

// 	if outputToken.Literal != expectedToken.Literal {
// 		return fmt.Errorf("Token literal are diffrent expecting literal: %v, output literal: %v", expectedToken.Literal, outputToken.Literal)
// 	}

// 	return nil
// }

func compareBlockStatement(outputStatement ast.BlockStatement, expectedStatement ast.BlockStatement) error {
	for i, statement := range outputStatement.Statements {
		err := compareStatement(statement, expectedStatement.Statements[i])

		if err != nil {
			return fmt.Errorf("Block statement index: %v failed \n %v", i, err)
		}
	}

	return nil
}

func TestParserExpression(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal ast.Expression
	}{
		{input: "true;", expectedVal: ast.Boolean{Value: true}},
		{input: "false;", expectedVal: ast.Boolean{Value: false}},
		{input: "nil;", expectedVal: ast.Nil{}},
		{input: "123.0;", expectedVal: ast.Float{Value: 123}},
		{input: "123.12;", expectedVal: ast.Float{Value: 123.12}},
		{input: "123;", expectedVal: ast.Integer{Value: 123}},
		{input: "\"123\";", expectedVal: ast.StringLiteral{Value: "123"}},
		{input: "let;", expectedVal: ast.Identifier{Value: "let"}},
		{input: "(let);", expectedVal: ast.GroupingExpression{Exp: ast.Identifier{Value: "let"}}},
		{input: "!true;", expectedVal: ast.PrefixExpression{
			Right:    ast.Boolean{Value: true},
			Operator: token.Bang,
		}},
		{input: "!!true;", expectedVal: ast.PrefixExpression{
			Right: ast.PrefixExpression{
				Right:    ast.Boolean{Value: true},
				Operator: token.Bang,
			},
			Operator: token.Bang,
		}},
		{input: "!!(true);", expectedVal: ast.PrefixExpression{
			Right: ast.PrefixExpression{
				Right: ast.GroupingExpression{
					Exp: ast.Boolean{Value: true},
				},
				Operator: token.Bang,
			},
			Operator: token.Bang,
		}},
		{input: "(!!true);", expectedVal: ast.GroupingExpression{
			Exp: ast.PrefixExpression{
				Right: ast.PrefixExpression{
					Right:    ast.Boolean{Value: true},
					Operator: token.Bang,
				},
				Operator: token.Bang,
			}},
		},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.Parse()

			if len(parser.Errors()) > 0 {
				for _, err := range parser.Errors() {
					t.Fatal(err)
				}
			}

			respExpression := resp.Statements[0].(ast.ExpressionStatement).Expression
			err := compareExpression(respExpression, testCase.expectedVal)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestParserExpressionArithmetic(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal ast.Expression
	}{
		{input: "123", expectedVal: ast.Integer{Value: 123}},
		{input: "-5", expectedVal: ast.PrefixExpression{
			Operator: token.Minus,

			Right: ast.Integer{
				Value: 5,
			},
		}},
		{input: "--5", expectedVal: ast.PrefixExpression{
			Operator: token.Minus,

			Right: ast.PrefixExpression{
				Operator: token.Minus,

				Right: ast.Integer{
					Value: 5,
				},
			},
		}},
		{input: "5+6", expectedVal: ast.InfixExpression{
			Operator: token.Plus,

			Left:  ast.Integer{Value: 5},
			Right: ast.Integer{Value: 6},
		}},

		{input: "5/2+3", expectedVal: ast.InfixExpression{
			Operator: token.Plus,

			Left: ast.InfixExpression{
				Operator: token.Slash,

				Right: ast.Integer{Value: 2},
				Left:  ast.Integer{Value: 5},
			},
			Right: ast.Integer{Value: 3},
		}},

		{input: "2+3/5", expectedVal: ast.InfixExpression{
			Operator: token.Plus,

			Left: ast.Integer{Value: 2},
			Right: ast.InfixExpression{
				Operator: token.Slash,

				Left:  ast.Integer{Value: 3},
				Right: ast.Integer{Value: 5},
			},
		}},
		{input: "5*3/2", expectedVal: ast.InfixExpression{
			Operator: token.Slash,

			Right: ast.Integer{Value: 2},
			Left: ast.InfixExpression{
				Operator: token.Star,

				Left:  ast.Integer{Value: 5},
				Right: ast.Integer{Value: 3},
			},
		}},

		{input: "51*36/71", expectedVal: ast.InfixExpression{
			Operator: token.Slash,

			Right: ast.Integer{Value: 71},
			Left: ast.InfixExpression{
				Operator: token.Star,

				Left:  ast.Integer{Value: 51},
				Right: ast.Integer{Value: 36},
			},
		}},
		{input: "(37 * -21 / (45 * 16))", expectedVal: ast.GroupingExpression{
			Exp: ast.InfixExpression{
				Operator: token.Slash,

				Right: ast.GroupingExpression{
					Exp: ast.InfixExpression{
						Operator: token.Star,

						Left:  ast.Integer{Value: 45},
						Right: ast.Integer{Value: 16},
					},
				},
				Left: ast.InfixExpression{
					Operator: token.Star,

					Left: ast.Integer{Value: 37},
					Right: ast.PrefixExpression{
						Operator: token.Minus,

						Right: ast.Integer{Value: 21},
					},
				},
			},
		}},
		{input: "1+2*3+4", expectedVal: ast.InfixExpression{
			Operator: token.Plus,

			Right: ast.Integer{Value: 4},
			Left: ast.InfixExpression{
				Operator: token.Plus,

				Left: ast.Integer{Value: 1},
				Right: ast.InfixExpression{
					Operator: token.Star,

					Right: ast.Integer{Value: 3},
					Left:  ast.Integer{Value: 2},
				},
			},
		},
		},
		{input: "52+80-94", expectedVal: ast.InfixExpression{
			Operator: token.Minus,

			Right: ast.Integer{Value: 94},
			Left: ast.InfixExpression{
				Operator: token.Plus,

				Left:  ast.Integer{Value: 52},
				Right: ast.Integer{Value: 80},
			},
		}},
		{input: "(1+2)*(3+4)", expectedVal: ast.InfixExpression{
			Operator: token.Star,

			Left: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,

					Left:  ast.Integer{Value: 1},
					Right: ast.Integer{Value: 2},
				},
			},
			Right: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,

					Left:  ast.Integer{Value: 3},
					Right: ast.Integer{Value: 4},
				},
			},
		}},
		{input: "(1+2)>(3+4)", expectedVal: ast.InfixExpression{
			Operator: token.Greater,

			Left: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,

					Left:  ast.Integer{Value: 1},
					Right: ast.Integer{Value: 2},
				},
			},
			Right: ast.GroupingExpression{
				Exp: ast.InfixExpression{
					Operator: token.Plus,

					Left:  ast.Integer{Value: 3},
					Right: ast.Integer{Value: 4},
				},
			},
		}},
		{input: "22 == 20", expectedVal: ast.InfixExpression{
			Operator: token.EqualEqual,

			Right: ast.Integer{Value: 20},
			Left:  ast.Integer{Value: 22},
		}},
		{input: "22 != 20", expectedVal: ast.InfixExpression{
			Operator: token.BangEqual,

			Right: ast.Integer{Value: 20},
			Left:  ast.Integer{Value: 22},
		}},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.parseExpression(Lowest)

			err := compareExpression(resp, testCase.expectedVal)
			if err != nil {
				t.Fatal(err)
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
						Expression: ast.StringLiteral{Value: "hello world"},
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
						Expression: ast.Integer{Value: 68},
					},
					ast.DeclarationStatement{
						Names:      []string{"world"},
						Expression: ast.Identifier{Value: "bar"},
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
						Expression: ast.Integer{Value: 68},
					},
					ast.DeclarationStatement{
						Names:      []string{"world"},
						Expression: ast.Integer{Value: 11},
					},
					ast.ExpressionStatement{
						Expression: ast.AssignExpression{IdentifierName: "bar", Value: ast.Identifier{Value: "world"}},
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
			for i, statement := range resp.Statements {
				expectedStatement := testCase.expectedVal.Statements[i]
				err := compareStatement(statement, expectedStatement)

				if err != nil {
					t.Fatalf("\n Response statement index: %v failed \n %v", i, err)
				}
			}
		})
	}
}

// if(aa) print aa;
// if(aa) {print "fdsf"}
// if(aa) {print "fdsf"} else if bb {}
// if statement
// array of condition -> code
// rest (alternative else conditon) -> code

func TestParserIfStatment(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal *ast.Program
	}{

		{
			input: "if(48){print foo;}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.Integer{

								Value: 48,
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Function: ast.Identifier{Value: "print"},
										Arguments: []ast.Expression{
											ast.Identifier{Value: "foo"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "if(48>11){print \"aaa\";}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.InfixExpression{
								Operator: token.Greater,
								Left: ast.Integer{

									Value: 48,
								},
								Right: ast.Integer{

									Value: 11,
								},
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Function: ast.Identifier{Value: "print"},
										Arguments: []ast.Expression{
											ast.StringLiteral{Value: "aaa"},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		{
			input: "if(48)print foo;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.Integer{

								Value: 48,
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Function: ast.Identifier{Value: "print"},
										Arguments: []ast.Expression{
											ast.Identifier{Value: "foo"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "if(40){print foo;}else{print bar;}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.Integer{

								Value: 40,
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Function: ast.Identifier{Value: "print"},
										Arguments: []ast.Expression{
											ast.Identifier{Value: "foo"},
										},
									},
								},
							},
						},
						Else: &ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Function: ast.Identifier{Value: "print"},
										Arguments: []ast.Expression{
											ast.Identifier{Value: "bar"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "if(40){print foo;}else if(30){print bar;}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.Integer{

								Value: 40,
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Function: ast.Identifier{Value: "print"},
										Arguments: []ast.Expression{
											ast.Identifier{Value: "foo"},
										},
									},
								},
							},
						},
						Else: &ast.BlockStatement{
							Statements: []ast.Statement{
								ast.IfStatement{
									Condition: ast.GroupingExpression{
										Exp: ast.Integer{

											Value: 30,
										},
									},
									Then: ast.BlockStatement{
										Statements: []ast.Statement{
											ast.ExpressionStatement{
												Expression: ast.CallExpression{
													Function: ast.Identifier{Value: "print"},
													Arguments: []ast.Expression{
														ast.Identifier{Value: "bar"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "if(40)if(30)print bar;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.Integer{

								Value: 40,
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.IfStatement{
									Condition: ast.GroupingExpression{
										Exp: ast.Integer{

											Value: 30,
										},
									},
									Then: ast.BlockStatement{
										Statements: []ast.Statement{
											ast.ExpressionStatement{
												Expression: ast.CallExpression{
													Function: ast.Identifier{Value: "print"},
													Arguments: []ast.Expression{
														ast.Identifier{Value: "bar"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, {
			input: "if(40)if(30)print bar;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.Integer{

								Value: 40,
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.IfStatement{
									Condition: ast.GroupingExpression{
										Exp: ast.Integer{

											Value: 30,
										},
									},
									Then: ast.BlockStatement{
										Statements: []ast.Statement{
											ast.ExpressionStatement{
												Expression: ast.CallExpression{
													Function: ast.Identifier{Value: "print"},
													Arguments: []ast.Expression{
														ast.Identifier{Value: "bar"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "if (false or \"ok\") print foo;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.IfStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.InfixExpression{
								Operator: token.OrToken,

								Left: ast.Boolean{

									Value: false,
								},
								Right: ast.StringLiteral{

									Value: "ok",
								},
							},
						},
						Then: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Function: ast.Identifier{Value: "print"},
										Arguments: []ast.Expression{
											ast.Identifier{Value: "foo"},
										},
									},
								},
							},
						},
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

			if len(parser.Errors()) > 0 {
				for _, err := range parser.Errors() {
					t.Fatal(err)
				}
			}
			if len(resp.Statements) != len(testCase.expectedVal.Statements) {
				t.Fatalf("Expected number of statement to be: %v, got: %v", len(testCase.expectedVal.Statements), len(resp.Statements))
			}
			for i, statement := range resp.Statements {
				expectedStatement := testCase.expectedVal.Statements[i]
				err := compareStatement(statement, expectedStatement)

				if err != nil {
					t.Fatalf("\n Response statement index: %v failed \n %v", i, err)
				}
			}
		})
	}
}

func toStatementPtr(s ast.Statement) *ast.Statement {
	return &s
}

func toExpressionPtr(s ast.Expression) *ast.Expression {
	return &s
}

func TestParserLoopStatment(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal *ast.Program
	}{

		{
			input: "while(3<2) print \"bar\";",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.WhileStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.InfixExpression{
								Operator: token.Less,

								Left:  ast.Integer{Value: 3},
								Right: ast.Integer{Value: 2},
							},
						},
						Block: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Arguments: []ast.Expression{
											ast.StringLiteral{Value: "bar"},
										},
										Function: ast.Identifier{Value: "print"},
									},
								},
							},
						},
					},
				},
			},
		},

		{
			input: "while(3<2) {print \"bar\";}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.WhileStatement{
						Condition: ast.GroupingExpression{
							Exp: ast.InfixExpression{
								Operator: token.Less,

								Left:  ast.Integer{Value: 3},
								Right: ast.Integer{Value: 2},
							},
						},
						Block: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Arguments: []ast.Expression{
											ast.StringLiteral{Value: "bar"},
										},
										Function: ast.Identifier{Value: "print"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "for (var baz =0;3<2;) print baz = baz +1;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.ForStatement{
						Condition: ast.InfixExpression{
							Operator: token.Less,

							Left:  ast.Integer{Value: 3},
							Right: ast.Integer{Value: 2},
						},
						Declaration: toStatementPtr(ast.DeclarationStatement{
							Names:      []string{"baz"},
							Expression: ast.Integer{Value: 0},
						}),
						Block: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Arguments: []ast.Expression{
											ast.AssignExpression{
												IdentifierName: "baz",
												Value: ast.InfixExpression{
													Left: ast.Identifier{
														Value: "baz",
													},
													Operator: token.Plus,

													Right: ast.Integer{
														Value: 1,
													},
												},
											},
										},
										Function: ast.Identifier{Value: "print"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `
			var baz=0;
			for (baz =1;baz<3; baz= baz+1) print baz;`,
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.DeclarationStatement{
						Names:      []string{"baz"},
						Expression: ast.Integer{Value: int64(0)},
					},
					ast.ForStatement{
						Condition: ast.InfixExpression{
							Operator: token.Less,
							Left:     ast.Identifier{Value: "baz"},
							Right:    ast.Integer{Value: 3},
						},
						Declaration: toStatementPtr(ast.ExpressionStatement{
							Expression: ast.AssignExpression{
								IdentifierName: "baz",
								Value:          ast.Integer{Value: int64(1)},
							},
						}),
						Evaluator: toExpressionPtr(ast.AssignExpression{
							IdentifierName: "baz",
							Value: ast.InfixExpression{
								Operator: token.Plus,
								Left:     ast.Identifier{Value: "baz"},
								Right:    ast.Integer{Value: int64(1)},
							},
						}),
						Block: ast.BlockStatement{
							Statements: []ast.Statement{
								ast.ExpressionStatement{
									Expression: ast.CallExpression{
										Arguments: []ast.Expression{
											ast.Identifier{Value: "baz"},
										},
										Function: ast.Identifier{Value: "print"},
									},
								},
							},
						},
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

			if len(parser.Errors()) > 0 {
				for _, err := range parser.Errors() {
					t.Fatal(err)
				}
			}
			if len(resp.Statements) != len(testCase.expectedVal.Statements) {
				t.Fatalf("Expected number of statement to be: %v, got: %v", len(testCase.expectedVal.Statements), len(resp.Statements))
			}
			for i, statement := range resp.Statements {
				expectedStatement := testCase.expectedVal.Statements[i]
				err := compareStatement(statement, expectedStatement)

				if err != nil {
					t.Fatalf("\n Response statement index: %v failed \n %v", i, err)
				}
			}
		})
	}
}

func TestParserFunctionStatment(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal *ast.Program
	}{

		{
			input: "print clock();",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Expression: ast.CallExpression{
							Function: ast.Identifier{Value: "print"},
							Arguments: []ast.Expression{
								ast.CallExpression{
									Function:  ast.Identifier{Value: "clock"},
									Arguments: []ast.Expression{},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "print clock() + 53;",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Expression: ast.CallExpression{
							Function: ast.Identifier{Value: "print"},
							Arguments: []ast.Expression{
								ast.InfixExpression{
									Left: ast.CallExpression{
										Function:  ast.Identifier{Value: "clock"},
										Arguments: []ast.Expression{},
									},
									Right:    ast.Integer{Value: int64(53)},
									Operator: token.Plus,
								},
							},
						},
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

			if len(parser.Errors()) > 0 {
				for _, err := range parser.Errors() {
					t.Fatal(err)
				}
			}
			if len(resp.Statements) != len(testCase.expectedVal.Statements) {
				t.Fatalf("Expected number of statement to be: %v, got: %v", len(testCase.expectedVal.Statements), len(resp.Statements))
			}
			for i, statement := range resp.Statements {
				expectedStatement := testCase.expectedVal.Statements[i]
				err := compareStatement(statement, expectedStatement)

				if err != nil {
					t.Fatalf("\n Response statement index: %v failed \n %v", i, err)
				}
			}
		})
	}
}
func TestParserBlockStatements(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal *ast.Program
	}{
		{
			input: "{var bar = 11;}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.BlockStatement{
						Statements: []ast.Statement{
							ast.DeclarationStatement{
								Names:      []string{"bar"},
								Expression: ast.Integer{Value: 11},
							},
						},
					},
				},
			},
		},
		{
			input: "{var bar = 11; print baz;}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.BlockStatement{
						Statements: []ast.Statement{
							ast.DeclarationStatement{
								Names:      []string{"bar"},
								Expression: ast.Integer{Value: 11},
							},
							ast.ExpressionStatement{
								Expression: ast.CallExpression{
									Function: ast.Identifier{Value: "print"},
									Arguments: []ast.Expression{
										ast.Identifier{Value: "baz"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "{var bar = 11; var world = 12; print bar + world;}",
			expectedVal: &ast.Program{
				Statements: []ast.Statement{
					ast.BlockStatement{
						Statements: []ast.Statement{
							ast.DeclarationStatement{
								Names:      []string{"bar"},
								Expression: ast.Integer{Value: 11},
							},
							ast.DeclarationStatement{
								Names:      []string{"world"},
								Expression: ast.Integer{Value: 12},
							},
							ast.ExpressionStatement{
								Expression: ast.CallExpression{
									Function: ast.Identifier{Value: "print"},
									Arguments: []ast.Expression{
										ast.InfixExpression{
											Left:     ast.Identifier{Value: "bar"},
											Right:    ast.Identifier{Value: "world"},
											Operator: "PLUS",
										},
									},
								},
							},
						},
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

			if len(parser.Errors()) > 0 {
				for _, err := range parser.Errors() {
					t.Fatal(err)
				}
			}
			for i, statement := range resp.Statements {
				expectedStatement := testCase.expectedVal.Statements[i]
				err := compareStatement(statement, expectedStatement)

				if err != nil {
					t.Fatalf("\n Response statement index: %v failed \n %v", i, err)
				}
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
			Function:  ast.Identifier{Value: "print"},
			Arguments: []ast.Expression{ast.GroupingExpression{Exp: ast.StringLiteral{Value: "aa"}}},
		}},
		{input: "print()", expectedVal: ast.CallExpression{
			Function:  ast.Identifier{Value: "print"},
			Arguments: []ast.Expression{},
		}},
		{input: "print \"aa\"", expectedVal: ast.CallExpression{
			Function:  ast.Identifier{Value: "print"},
			Arguments: []ast.Expression{ast.StringLiteral{Value: "aa"}},
		}},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.parseExpression(Lowest)

			if len(parser.Errors()) > 0 {
				for _, err := range parser.Errors() {
					t.Fatal(err)
				}
			}

			err := compareExpression(resp, testCase.expectedVal)

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestExpressionErrors(t *testing.T) {

	tests := []struct {
		input  string
		errors []string
	}{
		{input: "(72 +)", errors: []string{"[line 1] Error at ')': Expect expression.", "[line 1] expected next token to be RIGHT_PAREN, got EOF instead", "[line 1] expected token to be SEMICOLON, got RIGHT_PAREN instead"}},
		{input: "{var baz=1;", errors: []string{"[line 1] expected token to be RIGHT_BRACE, got EOF instead"}},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			parser.Parse()

			if !reflect.DeepEqual(parser.Errors(), testCase.errors) {
				t.Errorf("Expected to \nget: %#v, \ngot: %#v", testCase.errors, parser.Errors())
			}
		})
	}
}
