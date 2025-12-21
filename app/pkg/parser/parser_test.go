package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/lexar"
)

func TestParser(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal interface{}
	}{
		{input: "true", expectedVal: true},
		{input: "false", expectedVal: false},
		{input: "nil"},
		{input: "123.0", expectedVal: 123.0},
		{input: "123.12", expectedVal: 123.12},
		{input: "123", expectedVal: int64(123)},
		{input: "\"123\"", expectedVal: "123"},
		{input: "let", expectedVal: "let"},
		{input: "(let)", expectedVal: "(let)"},
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.Parse()
			var err error = testaa(resp, testCase.input, testCase.expectedVal)

			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func testaa(astNode ast.Expression, input string, expectedVal interface{}) error {
	var err error
	switch v := expectedVal.(type) {
	case bool:
		err = testBoolExpr(astNode, v)
	case nil:
		err = testNilExpr(astNode)
	case int64:
		err = testIntExpr(astNode, v)
	case float64:
		err = testFloatExpr(astNode, v)
	case string:
		if strings.Contains(input, "\"") {
			err = testStringLiteral(astNode, v)
		} else if strings.Contains(input, "(") {
			grouping, ok := astNode.(ast.Grouping)
			if !ok {
				return fmt.Errorf("Expected gropuing got: %v", reflect.TypeOf(astNode))
			}
			if grouping.String() != v {
				return fmt.Errorf("expected val: %v got:%v", v, grouping.String())
			}

		} else {
			err = testIdentifierLiteral(astNode, v)
		}
	default:
		err = fmt.Errorf("unsuported type: %v", reflect.TypeOf(v))
	}

	return err
}

func testIdentifierLiteral(exp ast.Expression, expectedVal string) error {
	astLiteral, ok := exp.(ast.Identifier)

	if !ok {
		return fmt.Errorf("Expected to get: %v got: %v", "identifier literal", reflect.TypeOf(exp))
	}

	if astLiteral.Value != expectedVal {
		return fmt.Errorf("Expected val:%v, got: %v", expectedVal, astLiteral.Value)
	}

	return nil
}

func testStringLiteral(exp ast.Expression, expectedVal string) error {
	astLiteral, ok := exp.(ast.StringLiteral)

	if !ok {
		return fmt.Errorf("Expected to get: %v got: %v", "string literal", reflect.TypeOf(exp))
	}

	if astLiteral.Value != expectedVal {
		return fmt.Errorf("Expected val:%v, got: %v", expectedVal, astLiteral.Value)
	}

	return nil
}

func testFloatExpr(exp ast.Expression, expectedVal float64) error {
	astBool, ok := exp.(ast.Float)

	if !ok {
		return fmt.Errorf("Expected to get: %v got: %v", ast.Inteager{}, reflect.TypeOf(exp))
	}

	if astBool.Value != expectedVal {
		return fmt.Errorf("Expected val:%v, got: %v", expectedVal, astBool.Value)
	}

	return nil
}

func testIntExpr(exp ast.Expression, expectedVal int64) error {
	astBool, ok := exp.(ast.Inteager)

	if !ok {
		return fmt.Errorf("Expected to get: %v got: %v", ast.Inteager{}, reflect.TypeOf(exp))
	}

	if astBool.Value != expectedVal {
		return fmt.Errorf("Expected val:%v, got: %v", expectedVal, astBool.Value)
	}

	return nil
}

func testBoolExpr(exp ast.Expression, expectedVal bool) error {
	astBool, ok := exp.(ast.Boolean)

	if !ok {
		return fmt.Errorf("Expected to get astboolean got: %v", reflect.TypeOf(exp))
	}

	if astBool.Value != expectedVal {
		return fmt.Errorf("Expected boolean expression val:%v, got: %v", expectedVal, astBool.Value)
	}

	return nil
}

func testNilExpr(exp ast.Expression) error {
	_, ok := exp.(ast.Nil)

	if !ok {
		return fmt.Errorf("Expected to get astboolean got: %v", reflect.TypeOf(exp))
	}

	return nil
}
