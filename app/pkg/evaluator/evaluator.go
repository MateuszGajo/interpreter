package evaluator

import (
	"fmt"
	"reflect"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/object"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/token"
)

var builtins = map[string]func(param []object.Object) object.Object{
	"print": func(args []object.Object) object.Object {
		if len(args) != 1 {
			return &object.Error{Message: "print should only have 1 arg"}
		}
		fmt.Println(args[0].Inspect())
		return &object.Nill{}
	},
}

// what can arguments be??
// literal string, interger, function, so basicaly object.Object array of them
//

func Eval(node ast.Node) object.Object {
	switch v := node.(type) {
	case *ast.Program:
		return evalProgram(v)
	case ast.ExpressionStatement:
		return Eval(v.Expression)
	case ast.InfixExpression:
		left := Eval(v.Left)
		right := Eval(v.Right)
		return evalInfixExpression(v.Operator, left, right)
	case ast.PrefixExpression:
		right := Eval(v.Right)
		return evalPrefixExpression(v.Operator, right)
	case ast.CallExpression:
		function := v.Function.(ast.Identifier).Value
		args := []object.Object{}

		for _, item := range v.Arguments {
			data := Eval(item)
			args = append(args, data)
		}

		builtIn, ok := builtins[function]

		if !ok {
			panic("not exists")
		}

		return builtIn(args)
	case ast.Boolean:
		return &object.Boolean{Value: v.Value}
	case ast.Nil:
		return &object.Nill{}
	case ast.Integer:
		return &object.Integer{Value: v.Value}
	case ast.Float:
		return &object.Float{Value: v.Value}
	case ast.StringLiteral:
		return &object.String{Value: v.Value}
	case ast.GroupingExpression:
		return Eval(v.Exp)
	default:
		panic(fmt.Sprintf("unsported node type: %v", reflect.TypeOf(node)))
		// return object.Error{Message: "fdsf"}
	}
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, item := range program.Statements {
		result = Eval(item)
	}

	return result

}

func evalPrefixExpression(operator token.TokenType, right object.Object) object.Object {

	switch operator {
	case token.Minus:
		if right.Type() != object.IntegerType && right.Type() != object.FloatType {
			return &object.Error{Message: "Operand must be a number"}
		}

		return &object.Integer{Value: -right.(*object.Integer).Value}
	case token.Bang:
		if right.Type() == object.IntegerType {
			return &object.Boolean{Value: right.(*object.Integer).Value == 0}
		} else if right.Type() == object.FloatType {
			return &object.Boolean{Value: right.(*object.Float).Value == 0}

		} else if right.Type() == object.NillType {
			return &object.Boolean{Value: true}

		} else if right.Type() == object.BooleanType {
			return &object.Boolean{Value: !right.(*object.Boolean).Value}
		}
		panic("operand must be a nil, boolean or int")
	default:
		panic(fmt.Sprintf("not supported infix, right type: %v, operator: %v", right.Type(), operator))

	}
}

var specialInfixOperatorErrors = map[token.TokenType]string{
	token.Star:         "Operands must be a number",
	token.Minus:        "Operands must be a number",
	token.Plus:         "Operands must be a number",
	token.Slash:        "Operands must be a number",
	token.Greater:      "Operands must be a number",
	token.GreaterEqual: "Operands must be a number",
	token.Less:         "Operands must be a number",
	token.LessEqual:    "Operands must be a number",
}

func evalInfixExpression(operator token.TokenType, left, right object.Object) object.Object {
	if left.Type() == object.FloatType || right.Type() == object.FloatType {
		leftFloat := left
		if left.Type() == object.IntegerType {
			leftFloat = &object.Float{Value: float64(left.(*object.Integer).Value)}
		}

		rightFloat := right
		if right.Type() == object.IntegerType {
			rightFloat = &object.Float{Value: float64(right.(*object.Integer).Value)}
		}

		return evalFLoatInfixExpression(operator, leftFloat, rightFloat)
	} else if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		return evalIntegerInfixExpression(operator, left, right)
	} else if left.Type() == object.StringType && right.Type() == object.StringType {
		resp := evalStringInfixExpression(operator, left, right)
		if resp != nil {
			return resp
		}
	} else if left.Type() == object.BooleanType && right.Type() == object.BooleanType {
		resp := evalBoolInfixExpression(operator, left, right)
		if resp != nil {
			return resp
		}
	} else if operator == token.EqualEqual {
		return object.Boolean{Value: false}
	}

	if _, ok := specialInfixOperatorErrors[operator]; ok {
		return &object.Error{Message: specialInfixOperatorErrors[operator]}
	}
	panic(fmt.Sprintf("not supported infix, left type: %v, right type: %v, operator: %v", left.Type(), right.Type(), operator))
}

func evalBoolInfixExpression(operator token.TokenType, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value
	switch operator {
	case token.EqualEqual:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.BangEqual:
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return nil

	}
}

func evalStringInfixExpression(operator token.TokenType, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case token.Plus:
		return &object.String{Value: leftVal + rightVal}
	case token.EqualEqual:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.BangEqual:
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return nil

	}
}

func evalIntegerInfixExpression(operator token.TokenType, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case token.Plus:
		return &object.Integer{Value: leftVal + rightVal}
	case token.Star:
		return &object.Integer{Value: leftVal * rightVal}
	case token.Minus:
		return &object.Integer{Value: leftVal - rightVal}
	case token.Greater:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GreaterEqual:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.Less:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LessEqual:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.EqualEqual:
		return &object.Boolean{Value: leftVal == rightVal}
	case token.BangEqual:
		return &object.Boolean{Value: leftVal != rightVal}
	case token.Slash:
		result := float64(leftVal) / float64(rightVal)
		if result == float64(int64(result)) {
			return &object.Integer{Value: int64(result)}
		} else {
			return &object.Float{Value: result}
		}
	default:
		panic(fmt.Sprintf("eval integer infix expression not supported operator %v", operator))
	}
}

func evalFLoatInfixExpression(operator token.TokenType, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case token.Plus:
		return &object.Float{Value: leftVal + rightVal}
	case token.Star:
		return &object.Float{Value: leftVal * rightVal}
	case token.Minus:
		return &object.Float{Value: leftVal - rightVal}
	case token.Greater:
		return &object.Boolean{Value: leftVal > rightVal}
	case token.GreaterEqual:
		return &object.Boolean{Value: leftVal >= rightVal}
	case token.Less:
		return &object.Boolean{Value: leftVal < rightVal}
	case token.LessEqual:
		return &object.Boolean{Value: leftVal <= rightVal}
	case token.Slash:
		return &object.Float{Value: leftVal / rightVal}
	default:
		panic(fmt.Sprintf("not supported operator %v", operator))
	}
}
