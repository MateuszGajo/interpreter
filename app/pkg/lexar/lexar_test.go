package lexar

import (
	"errors"
	"fmt"
	"testing"
)

type TestCase struct {
	input  string
	output []Token
	errors []error
}

func TestLexar(t *testing.T) {
	cases := []TestCase{
		{
			input: "(()",
			output: []Token{
				{TokenType: LeftParen, Lexeme: "(", Literal: nil},
				{TokenType: LeftParen, Lexeme: "(", Literal: nil},
				{TokenType: RightParen, Lexeme: ")", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "{{}}",
			output: []Token{
				{TokenType: LeftBrace, Lexeme: "{", Literal: nil},
				{TokenType: LeftBrace, Lexeme: "{", Literal: nil},
				{TokenType: RightBrace, Lexeme: "}", Literal: nil},
				{TokenType: RightBrace, Lexeme: "}", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "({*.,+*-;})",
			output: []Token{
				{TokenType: LeftParen, Lexeme: "(", Literal: nil},
				{TokenType: LeftBrace, Lexeme: "{", Literal: nil},
				{TokenType: Star, Lexeme: "*", Literal: nil},
				{TokenType: Dot, Lexeme: ".", Literal: nil},
				{TokenType: Comma, Lexeme: ",", Literal: nil},
				{TokenType: Plus, Lexeme: "+", Literal: nil},
				{TokenType: Star, Lexeme: "*", Literal: nil},
				{TokenType: Minus, Lexeme: "-", Literal: nil},
				{TokenType: Semicolon, Lexeme: ";", Literal: nil},
				{TokenType: RightBrace, Lexeme: "}", Literal: nil},
				{TokenType: RightParen, Lexeme: ")", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "{{$}}#",
			output: []Token{
				{TokenType: LeftBrace, Lexeme: "{", Literal: nil},
				{TokenType: LeftBrace, Lexeme: "{", Literal: nil},
				{TokenType: RightBrace, Lexeme: "}", Literal: nil},
				{TokenType: RightBrace, Lexeme: "}", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
			errors: []error{errors.New("[line 1] Error: Unexpected character: $"), errors.New("[line 1] Error: Unexpected character: #")},
		},
		{
			input: "==,=,!!==,>>=<<=",
			output: []Token{
				{TokenType: EqualEqual, Lexeme: "==", Literal: nil},
				{TokenType: Comma, Lexeme: ",", Literal: nil},
				{TokenType: Equal, Lexeme: "=", Literal: nil},
				{TokenType: Comma, Lexeme: ",", Literal: nil},
				{TokenType: Bang, Lexeme: "!", Literal: nil},
				{TokenType: BangEqual, Lexeme: "!=", Literal: nil},
				{TokenType: Equal, Lexeme: "=", Literal: nil},
				{TokenType: Comma, Lexeme: ",", Literal: nil},
				{TokenType: Greater, Lexeme: ">", Literal: nil},
				{TokenType: GreaterEqual, Lexeme: ">=", Literal: nil},
				{TokenType: Less, Lexeme: "<", Literal: nil},
				{TokenType: LessEqual, Lexeme: "<=", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "/,\n//coMment\n,\t \r",
			output: []Token{
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Comma, Lexeme: ",", Literal: nil},
				{TokenType: Comma, Lexeme: ",", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "/\",,123\"/",
			output: []Token{
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: StringToken, Lexeme: "\",,123\"", Literal: ",,123"},
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "\"baz\" . \"unterminated",
			output: []Token{
				{TokenType: StringToken, Lexeme: "\"baz\"", Literal: "baz"},
				{TokenType: Dot, Lexeme: ".", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
			errors: []error{errors.New("[line 1] Error: Unterminated string.")},
		},
		{
			input: "/123.123/134/123.00000",
			output: []Token{
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Number, Lexeme: "123.123", Literal: 123.123},
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Number, Lexeme: "134", Literal: 134.0},
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Number, Lexeme: "123.00000", Literal: 123.0},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "/foo bar foo_123/",
			output: []Token{
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Identifier, Lexeme: "foo", Literal: nil},
				{TokenType: Identifier, Lexeme: "bar", Literal: nil},
				{TokenType: Identifier, Lexeme: "foo_123", Literal: nil},
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
		{
			input: "/foo whiLe foo_123/ nil",
			output: []Token{
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: Identifier, Lexeme: "foo", Literal: nil},
				{TokenType: WhileToken, Lexeme: "while", Literal: nil},
				{TokenType: Identifier, Lexeme: "foo_123", Literal: nil},
				{TokenType: Slash, Lexeme: "/", Literal: nil},
				{TokenType: NilToken, Lexeme: "nil", Literal: nil},
				{TokenType: Eof, Lexeme: "", Literal: nil},
			},
		},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("running input: %v", testCase.input), func(t *testing.T) {
			lexar := NewLexar([]byte(testCase.input))
			tokens, errors := lexar.Scan()

			if len(testCase.output) != len(tokens) {
				t.Errorf("length mismatch: expected %d tokens, got %d",
					len(testCase.output), len(tokens))
			}

			if len(testCase.errors) != len(errors) {
				t.Errorf("length mismatch: expected %d errors, got %d",
					len(testCase.errors), len(errors))
			}

			for i := range testCase.errors {
				expected := testCase.errors[i]
				got := errors[i]

				if expected.Error() != got.Error() {
					t.Errorf("\nerror[%d] mismatch:\n expected: %+q\n got:      %+q", i, expected, got)
				}
			}

			for i := range testCase.output {
				expected := testCase.output[i]
				got := tokens[i]

				if expected != got {
					t.Errorf("\ntoken[%d] mismatch:\n  expected: %+v\n got:      %+v", i, expected, got)
				}
			}
		})
	}
}
