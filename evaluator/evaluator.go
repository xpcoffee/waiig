package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalInfixExpression(left, node.Operator, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	case *ast.LetStatement:
		return evalLetStatement(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteralExpression:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.FunctionCallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			fmt.Printf("problem inital Eval: %s\n", function.Inspect())
			return function
		}

		args := evalExpressions(node.Parameters, env)
		if len(args) == 1 && isError(args[0]) {
			fmt.Printf("problem with parameters: %s\n", args[0].Inspect())
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		pairs := make(map[object.HashKey]object.HashPair)
		for k, v := range node.Pairs {
			value := Eval(v, env)
			keyObj := Eval(k, env)
			if hashableObj, ok := keyObj.(object.Hashable); !ok {
				return newError("Cannot use as key %s", keyObj.Type())
			} else {
				pairs[hashableObj.HashKey()] = object.HashPair{Key: keyObj, Value: value}
			}
		}
		return &object.Hash{Pairs: pairs}

	case *ast.IndexingExpression:
		target := Eval(node.Target, env)
		switch target := target.(type) {
		case *object.Array:
			evaluatedIndex := Eval(node.Index, env)
			if evaluatedIndex.Type() != object.INTEGER_OBJ {
				return newError("Cannot use as index %s", evaluatedIndex.Type())
			}
			index := evaluatedIndex.(*object.Integer)

			if index.Value < 0 {
				return newError("Cannot index with a negative number %d", index.Value)
			}

			if index.Value >= int64(len(target.Elements)) {
				return newError("Index is larger than the max. index=%d, max=%d", index.Value, len(target.Elements)-1)
			}

			return target.Elements[index.Value]
		case *object.Hash:
			evaluatedIndex := Eval(node.Index, env)

			if hashableObj, ok := evaluatedIndex.(object.Hashable); !ok {
				return newError("Cannot use as index %s", evaluatedIndex.Type())
			} else {
				return target.Pairs[hashableObj.HashKey()].Value
			}
		default:
			return newError("Cannot index type %s", target.Type())
		}
	}

	return nil
}

func isHashIndexType(obj object.Object) bool {
	switch obj.Type() {
	case object.INTEGER_OBJ:
		return true
	case object.STRING_OBJ:
		return true
	case object.BOOLEAN_OBJ:
		return true
	default:
		return false
	}
}

// returns the evalutation of the LAST statement
func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(blockStatement *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range blockStatement.Statements {
		result = Eval(stmt, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unkown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixOperator(
			left.(*object.Integer),
			operator,
			right.(*object.Integer),
		)

	case right.Type() == object.STRING_OBJ && left.Type() == object.STRING_OBJ:
		if operator == "+" {
			return &object.String{Value: left.(*object.String).Value + right.(*object.String).Value}
		}
		return newError("unkown operator: %s %s %s", left.Type(), operator, right.Type())

	case operator == "==":
		// the == and != operators do pointer comparison for boolean and NULL
		// other evaluations (string, objects etc) need to happen before this point
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unkown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalBangOperatorExpression(exp object.Object) object.Object {
	switch exp {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		// this default catches all non-null values and flips them to "false"
		// that means that they are truthy by default
		return FALSE
	}
}

func evalMinusOperatorExpression(exp object.Object) object.Object {
	if exp.Type() != object.INTEGER_OBJ {
		return newError("unkown operator: -%s", exp.Type())
	}

	value := exp.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixOperator(left *object.Integer, operator string, right *object.Integer) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: left.Value + right.Value}
	case "-":
		return &object.Integer{Value: left.Value - right.Value}
	case "*":
		return &object.Integer{Value: left.Value * right.Value}
	case "/":
		return &object.Integer{Value: left.Value / right.Value}
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	case ">":
		return nativeBoolToBooleanObject(left.Value > right.Value)
	case "<":
		return nativeBoolToBooleanObject(left.Value < right.Value)
	default:
		return newError("unkown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) object.Object {
	value := Eval(rs.ReturnValue, env)
	if isError(value) {
		return value
	}
	return &object.ReturnValue{Value: value}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}

	return obj.Type() == object.ERROR_OBJ
}

func evalLetStatement(ls *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(ls.Value, env)
	if isError(val) {
		return val
	}
	env.Set(ls.Name.Value, val)

	return val
}

func evalIdentifier(ie *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(ie.Value); ok {
		return val
	}

	if builtin, ok := builtins[ie.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + ie.Value)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	results := []object.Object{}

	for _, exp := range expressions {
		result := Eval(exp, env)
		if isError(result) {
			return []object.Object{result}
		}
		results = append(results, result)
	}

	return results
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		closure := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, closure)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %T", fn)
	}

}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIndex, param := range fn.Parameters {
		env.Set(param.Value, args[paramIndex])
	}

	return env
}

// Prevents a value returned from a function from short-circuiting
// parent blocks
func unwrapReturnValue(obj object.Object) object.Object {
	rtnVal, ok := obj.(*object.ReturnValue)
	if !ok {
		return obj
	}
	return rtnVal.Value
}
