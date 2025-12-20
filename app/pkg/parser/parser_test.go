package parser

import (
	"fmt"
	"reflect"
	"testing"

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
	}

	for _, testCase := range tests {
		t.Run(fmt.Sprintf("Running input: %v", testCase.input), func(t *testing.T) {
			lexar := lexar.NewLexar([]byte(testCase.input))
			parser := NewParser(*lexar)

			resp := parser.Parse()
			var err error
			switch v := testCase.expectedVal.(type) {
			case bool:
				err = testBoolExpr(resp, v)
			case nil:
				err = testNilExpr(resp)
			}

			if err != nil {
				t.Error(err.Error())
			}
		})
	}

}

func testBoolExpr(exp ASTExpression, expectedVal bool) error {
	astBool, ok := exp.(ASTBoolean)

	if !ok {
		return fmt.Errorf("Expected to get astboolean got: %v", reflect.TypeOf(exp))
	}

	if astBool.value != expectedVal {
		return fmt.Errorf("Expected boolean expression val:%v, got: %v", expectedVal, astBool.value)
	}

	return nil
}

func testNilExpr(exp ASTExpression) error {
	_, ok := exp.(ASTNil)

	if !ok {
		return fmt.Errorf("Expected to get astboolean got: %v", reflect.TypeOf(exp))
	}

	return nil
}
