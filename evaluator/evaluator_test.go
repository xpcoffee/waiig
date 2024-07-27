package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 - 10", 50},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("evaluated object is not an object.Integer. got=%T", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Unexpected evaluated value. expected=%d got=%d", expected, result.Value)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"false != true", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("evaluated object is not an object.Boolean. got=%T", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Unexpected evaluated value. expected=%t got=%t", expected, result.Value)
		return false
	}
	return true
}

func TestEvalBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		// integers are truthy
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestIfExpresssions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if(true) { 10 }", 10},
		{"if(false) { 10 }", nil},
		{"if(1) { 10 }", 10},
		{"if(2 < 3) { 10 }", 10},
		{"if(2 > 3) { 10 }", nil},
		{"if(2 < 3) { 10 } else { 20 }", 10},
		{"if(2 > 3) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}

}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("evaluated object is not an object.Null. got=%T", obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 5; 9", 5},
		{"return 2 * 6; 9", 12},
		{"7; return 2 * 3; 9", 6},
		{
			`if(10 > 1) {
                if(10 > 1) {
                    return 10;
                }
                return 1;
            }`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input            string
		expected_message string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unkown operator: -BOOLEAN"},
		{"true + false", "unkown operator: BOOLEAN + BOOLEAN"},
		{"true + false; 5", "unkown operator: BOOLEAN + BOOLEAN"},
		{"if(10 > 1) { true + false }", "unkown operator: BOOLEAN + BOOLEAN"},
		{
			`if(10 > 1) {
                if(10 > 1) {
                    return true + false;
                }
                return 1;
            }`,
			"unkown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("%s", tt.input)
			t.Errorf("evaluated object is not an object.Error. got=%T", errObj)
			continue
		}
		if errObj.Message != tt.expected_message {
			t.Errorf("Unexpected evaluated message. expected=%s got=%s", tt.expected_message, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a", 5},
		{"let a = 5 * 5; a", 25},
		{"let a = 6; let b = a; b", 6},
		{"let a = 7; let b = a + 1; let c = 2 * a + b; c", 22},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2 ;};"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("evaluated object is not an object.Function. got=%T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("incorrect number of parameters. expected=1 got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("incorrect parameter. expected=x got=%s", fn.Parameters[0].String())
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("incorrect function body. expected=%s got=%s", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5)", 5},
		{"let identity = fn(x) { return x; }; identity(6)", 6},
		{"let double = fn(y) { 2 * y; }; double(7)", 14},
		{"let add = fn(x, y) { x + y; }; add(8, 8)", 16},
		{"let add = fn(x, y) { x + y; }; add(5 + 6, add(7, 8))", 26},
		{"fn(x){ x; }(9);", 9},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, world!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not a string, got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, world!" {
		t.Errorf("String has the wrong value. expected=%q got=%q", "Hello, world!", str.Value)
	}
}
