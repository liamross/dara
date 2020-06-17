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

func Eval(node ast.Node, env *Environment) Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfStatement:
		return evalIfStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	// Expressions
	case *ast.NumberLiteral:
		return &Number{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &String{Value: node.Value}

	case *ast.Nil:
		return &Nil{}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.DeclareExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.AssignExpression:
		return evalAssignExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}

	return nil
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func evalProgram(program *ast.Program, env *Environment) Object {
	var result Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch r := result.(type) {
		case *ReturnValue:
			return r.Value
		case *Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			if rt := result.Type(); rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfStatement(is *ast.IfStatement, env *Environment) Object {
	condition := Eval(is.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition.Type() != BOOLEAN_OBJ {
		return newError("type mismatch: non-boolean condition %s (%s) in if statement",
			condition.Inspect(), condition.Type())
	}
	if condition == TRUE {
		return Eval(is.Consequence, env)
	}
	if is.Alternative != nil {
		return Eval(is.Alternative, env)
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

func evalAssignExpression(node *ast.AssignExpression, env *Environment) Object {
	if _, ok := env.Get(node.Name.Value); !ok {
		return newError("undeclared name: %s", node.Name.Value)
	}
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	env.Set(node.Name.Value, val)
	return val
}

func evalIdentifier(node *ast.Identifier, env *Environment) Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("undeclared name: %s", node.Value)
	}
	return val
}

func evalExpressions(exps []ast.Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn Object, args []Object) Object {
	function, ok := fn.(*Function)
	if !ok {
		return newError("invalid operation: can not call non-function (%s)", fn.Type())
	}
	var (
		extendedEnv = extendedFunctionEnv(function, args)
		evaluated   = Eval(function.Body, extendedEnv)
	)
	return unwrapReturnValue(evaluated)
}

func extendedFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewScopedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
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
