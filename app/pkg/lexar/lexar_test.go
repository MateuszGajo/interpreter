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
				{tokenType: leftParen, lexeme: "(", literal: ""},
				{tokenType: leftParen, lexeme: "(", literal: ""},
				{tokenType: rightParen, lexeme: ")", literal: ""},
				{tokenType: eof, lexeme: "", literal: ""},
			},
		},
		{
			input: "{{}}",
			output: []Token{
				{tokenType: leftBrace, lexeme: "{", literal: ""},
				{tokenType: leftBrace, lexeme: "{", literal: ""},
				{tokenType: rightBrace, lexeme: "}", literal: ""},
				{tokenType: rightBrace, lexeme: "}", literal: ""},
				{tokenType: eof, lexeme: "", literal: ""},
			},
		},
		{
			input: "({*.,+*-;})",
			output: []Token{
				{tokenType: leftParen, lexeme: "(", literal: ""},
				{tokenType: leftBrace, lexeme: "{", literal: ""},
				{tokenType: star, lexeme: "*", literal: ""},
				{tokenType: dot, lexeme: ".", literal: ""},
				{tokenType: comma, lexeme: ",", literal: ""},
				{tokenType: plus, lexeme: "+", literal: ""},
				{tokenType: star, lexeme: "*", literal: ""},
				{tokenType: minus, lexeme: "-", literal: ""},
				{tokenType: semicolon, lexeme: ";", literal: ""},
				{tokenType: rightBrace, lexeme: "}", literal: ""},
				{tokenType: rightParen, lexeme: ")", literal: ""},
				{tokenType: eof, lexeme: "", literal: ""},
			},
		},
		{
			input: "{{$}}#",
			output: []Token{
				{tokenType: leftBrace, lexeme: "{", literal: ""},
				{tokenType: leftBrace, lexeme: "{", literal: ""},
				{tokenType: rightBrace, lexeme: "}", literal: ""},
				{tokenType: rightBrace, lexeme: "}", literal: ""},
				{tokenType: eof, lexeme: "", literal: ""},
			},
			errors: []error{errors.New("[line 1] Error: Unexpected character: $"), errors.New("[line 1] Error: Unexpected character: #")},
		},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("running input: %v", testCase.input), func(t *testing.T) {
			lexar := NewLexar([]byte(testCase.input))
			tokens, errors := lexar.Scan()

			if len(testCase.output) != len(tokens) {
				t.Fatalf("length mismatch: expected %d tokens, got %d",
					len(testCase.output), len(tokens))
			}

			if len(testCase.errors) != len(errors) {
				t.Fatalf("length mismatch: expected %d errors, got %d",
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
