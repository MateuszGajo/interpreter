package evaluator

import (
	"fmt"
	"reflect"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/object"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

func Eval(node ast.Node) object.Object {
	switch v := node.(type) {
	case ast.InfixExpression:
		left := Eval(v.Left)
		right := Eval(v.Right)
		return evalInfixExpression(v.Operator, left, right)
	case ast.Boolean:
		return object.Boolean{Value: v.Value}
	case ast.Nil:
		return object.Nill{}
	case ast.Integer:
		return object.Integer{Value: v.Value}
	case ast.Float:
		return object.Float{Value: v.Value}
	case ast.StringLiteral:
		return object.String{Value: v.Value}
	default:
		panic(fmt.Sprintf("unsported node type: %v", reflect.TypeOf(node)))
	}
}

func evalInfixExpression(operator token.TokenType, left, right object.Object) object.Object {
	if left.Type() == object.IntegerType && left.Type() == object.IntegerType {
		return evalInfixExpression(operator, left, right)
	} else {
		panic(fmt.Sprintf("not supported infix, left type: %v, right type: %v, operator: %v", left.Type(), right.Type(), operator))
	}
}

func evalIntegerInfixExpression(operator token.TokenType, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case token.Plus:
		return object.Integer{Value: leftVal + rightVal}

	default:
		panic(fmt.Sprintf("not supported operator %v", operator))
	}
}
