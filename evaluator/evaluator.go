package evaluator

import (
	"dara/ast"
	"fmt"
	"math"
)

var (
	NIL   = &Nil{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

func Eval(node ast.Node) Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfStatement:
		return evalIfStatement(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &ReturnValue{Value: val}

	// Expressions
	case *ast.NumberLiteral:
		return &Number{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Nil:
		return &Nil{}
	case *ast.StringLiteral:
		return &String{Value: node.Value}
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func evalProgram(program *ast.Program) Object {
	var result Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch r := result.(type) {
		case *ReturnValue:
			return r.Value
		case *Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) Object {
	var result Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			if rt := result.Type(); rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfStatement(is *ast.IfStatement) Object {
	condition := Eval(is.Condition)

	if condition.Type() != BOOLEAN_OBJ {
		return newError("type mismatch: non-boolean condition %s (%s) in if statement",
			condition.Inspect(), condition.Type())
	}
	if condition == TRUE {
		return Eval(is.Consequence)
	}
	if is.Alternative != nil {
		return Eval(is.Alternative)
	}
	return NIL
}

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("invalid operation: operator %s is not defined for %s (%s)",
			operator, right.Inspect(), right.Type())
	}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == NUMBER_OBJ && right.Type() == NUMBER_OBJ:
		return evalArithmeticInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("invalid operation: operator %s is not defined for %s (%s)",
			operator, left.Inspect(), left.Type())
	}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return newError("invalid operation: operator %s is not defined for %s (%s)",
			"!", right.Inspect(), right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != NUMBER_OBJ {
		return newError("invalid operation: operator %s is not defined for %s (%s)",
			"-", right.Inspect(), right.Type())
	}

	value := right.(*Number).Value
	return &Number{Value: -value}
}

func evalArithmeticInfixExpression(operator string, left, right Object) Object {
	var (
		leftVal  = left.(*Number).Value
		rightVal = right.(*Number).Value
	)

	switch operator {
	case "+":
		return &Number{Value: leftVal + rightVal}
	case "-":
		return &Number{Value: leftVal - rightVal}
	case "*":
		return &Number{Value: leftVal * rightVal}
	case "/":
		return &Number{Value: leftVal / rightVal}
	case "%":
		return &Number{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("invalid operation: operator %s is not defined for %s (%s)",
			operator, right.Inspect(), right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	var (
		leftVal  = left.(*String).Value
		rightVal = right.(*String).Value
	)

	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("invalid operation: operator %s is not defined for %s (%s)",
			operator, left.Inspect(), left.Type())
	}
}
