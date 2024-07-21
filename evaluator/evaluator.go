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
		return evalStatements(node.Statements)

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
	}

	return nil
}

// returns the evalutation of the LAST statement
func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
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
	switch operator {
	case "+":
		return evalAddition(left, right)
	case "-":
		return evalSubtraction(left, right)
	case "*":
		return evalProduct(left, right)
	case "/":
		return evalDivision(left, right)
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

func evalAddition(left object.Object, right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ || left.Type() != object.INTEGER_OBJ {
		return NULL
	}

	a := left.(*object.Integer).Value
	b := right.(*object.Integer).Value
	return &object.Integer{Value: a + b}
}

func evalSubtraction(left object.Object, right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ || left.Type() != object.INTEGER_OBJ {
		return NULL
	}

	a := left.(*object.Integer).Value
	b := right.(*object.Integer).Value
	return &object.Integer{Value: a - b}
}

func evalProduct(left object.Object, right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ || left.Type() != object.INTEGER_OBJ {
		return NULL
	}

	a := left.(*object.Integer).Value
	b := right.(*object.Integer).Value
	return &object.Integer{Value: a * b}
}

func evalDivision(left object.Object, right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ || left.Type() != object.INTEGER_OBJ {
		return NULL
	}

	a := left.(*object.Integer).Value
	b := right.(*object.Integer).Value
	return &object.Integer{Value: a / b}
}
