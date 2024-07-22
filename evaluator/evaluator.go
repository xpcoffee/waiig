package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		return evalInfixExpression(left, node.Operator, right)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.ReturnStatement:
		return evalReturnStatement(node)
	}

	return nil
}

// returns the evalutation of the LAST statement
func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(blockStatement *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range blockStatement.Statements {
		result = Eval(stmt)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
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
	}

	return NULL
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixOperator(
			left.(*object.Integer),
			operator,
			right.(*object.Integer),
		)

	case operator == "==":
		// the == and != operators do pointer comparison for boolean and NULL
		// other evaluations (string, objects etc) need to happen before this point
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	}

	return NULL
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
		return NULL
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
	}

	return NULL
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
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

func evalReturnStatement(rs *ast.ReturnStatement) object.Object {
	value := Eval(rs.ReturnValue)
	return &object.ReturnValue{Value: value}
}
