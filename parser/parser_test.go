package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors:", len(errors))
	for _, error := range errors {
		t.Errorf("parser error: %q", error)
	}
	t.FailNow()
}

func TestLetStatements(t *testing.T) {
	input := `
        let x = 5;
        let y = 10;
        let fooBar = 838383;
    `

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"fooBar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !(testLetStatement(t, stmt, tt.expectedIdentifier)) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("token literal is not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not an ast.LetStatement. got=%s", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not %s. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral is not %s. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
        return 5;
        return 10;
        return 9933322;
    `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedExpression string
	}{
		{"5"},
		{"10"},
		{"9933322"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !(testReturnStatement(t, stmt, tt.expectedExpression)) {
			return
		}
	}
}

func testReturnStatement(t *testing.T, s ast.Statement, name string) bool {
	rtrnStmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("statement is not an ast.ReturnStatement. got=%s", s)
		return false
	}

	if rtrnStmt.TokenLiteral() != "return" {
		t.Errorf("token literal is not 'return'. got=%s", s.TokenLiteral())
		return false
	}

	// if rtrnStmt.ReturnValue.String() != name {
	// 	t.Errorf("Unexpected return expression. expected=%q got=%q", name, rtrnStmt.ReturnValue.String())
	// 	return false
	// }

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected a single statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not an expression. Got %T", program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	idnt, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Statement is not an Identifier. Got %T", exp)
		return false
	}

	if idnt.Value != value {
		t.Errorf("Unexpected value. Expected %q, Got %q", value, idnt.Value)
		return false
	}

	if idnt.TokenLiteral() != value {
		t.Errorf("Unexpected TokenLiteral. Expected %q, Got %q", value, idnt.TokenLiteral())
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected a single statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not an expression. Got %T", program.Statements[0])
	}

	testBoolean(t, stmt.Expression, true)
}

func testBoolean(t *testing.T, exp ast.Expression, b bool) bool {
	idnt, ok := exp.(*ast.BooleanExpression)
	if !ok {
		t.Errorf("Statement is not an Identifier. Got %T", exp)
		return false
	}

	if idnt.Value != b {
		t.Errorf("Unexpected value. Expected %t, Got %t", b, idnt.Value)
		return false
	}

	expected := fmt.Sprintf("%t", b)
	if idnt.TokenLiteral() != expected {
		t.Errorf("Unexpected TokenLiteral. Expected %t, Got %q", b, expected)
		return false
	}

	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected a single statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not an expression. Got %T", program.Statements[0])
	}

	if !testLiteralExpression(t, stmt.Expression, 5) {
		return
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}
	t.Errorf("Type of literal expression is unkown got=%T", expected)
	return false
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("Expected a single statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Statement is not an expression. Got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("Statement is not a PrefixExpression. Got %T", program.Statements[0])
		}

		if exp.Operator != tt.operator {
			t.Errorf("Unexpected operator. wanted=%q got=%q", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Not an IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("Value incorrect. wanted=%d got=%d", value, integ.Value)
		return false
	}

	expectedLiteral := fmt.Sprintf("%d", value)
	if integ.TokenLiteral() != expectedLiteral {
		t.Errorf("TokenLiteral incorrect. wanted=%q got=%q", expectedLiteral, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected a single statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not an expression. Got %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.left, tt.operator, tt.right) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	ie, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Statement is not a InfixExpression. Got %T", exp)
		return false
	}

	if !testLiteralExpression(t, ie.Left, left) {
		return false
	}

	if ie.Operator != operator {
		t.Errorf("Operator incorrect. wanted=%q got=%q", operator, ie.Operator)
		return false
	}

	if !testLiteralExpression(t, ie.Right, right) {
		return false
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b / c + d * e - f",
			"(((a + (b / c)) + (d * e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		// boolean expressions
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		// grouped expressions
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 1",
			"((5 + 5) * 1)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Parsing result is unexpected. wanted=%q got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expected := "if (x < y) x"
	actual := program.String()
	if actual != expected {
		t.Errorf("Parsing result is unexpected. wanted=%q got=%q", expected, actual)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected a single statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression. Got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Statement is not an IfExpression. Got %T", program.Statements[0])
	}

	cnd, ok := exp.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("Condition is not a InfixExpression. Got %T", exp.Condition)
	}

	if !testInfixExpression(t, cnd, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("Expected a single consequence statement, got %d", len(program.Statements))
	}

	cqc, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Condition is not a BlockStatement. Got %T", exp.Consequence)
	}

	if !testIdentifier(t, cqc.Expression, "x") {
		return
	}

	if exp.Alternative != nil && exp.Alternative.Statements != nil {
		t.Fatalf("Expected a nil alternative, got=%q", exp.Alternative.String())
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expected := "if (x < y) x else y"
	actual := program.String()
	if actual != expected {
		t.Errorf("Parsing result is unexpected. wanted=%q got=%q", expected, actual)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected a single statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression. Got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Statement is not an IfExpression. Got %T", program.Statements[0])
	}

	cnd, ok := exp.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("Condition is not a InfixExpression. Got %T", exp.Condition)
	}

	if !testInfixExpression(t, cnd, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("Expected a single consequence statement, got %d", len(program.Statements))
	}

	cqc, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence is not a BlockStatement. Got %T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, cqc.Expression, "x") {
		return
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative is not a BlockStatement. Got %T", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `fn (x, y) { x + y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expected := "fn(x,y)(x + y)"
	actual := program.String()
	if actual != expected {
		t.Errorf("Parsing result is unexpected. wanted=%q got=%q", expected, actual)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected a single statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression. Got %T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteralExpression)
	if !ok {
		t.Fatalf("Statement is not a FunctionLiteralExpression. Got %T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("Expected two parameters, got=%d", len(function.Parameters))
	}
	testLiteralExpression(t, &function.Parameters[0], "x")
	testLiteralExpression(t, &function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Expected single statement in body, got=%d", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Body statement is not a expression. Got %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input              string
		expectedParameters []string
	}{
		{input: "fn(){}", expectedParameters: []string{}},
		{input: "fn(x){}", expectedParameters: []string{"x"}},
		{input: "fn(x,y,z){}", expectedParameters: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected a single statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not an expression. Got %T", program.Statements[0])
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteralExpression)
		if !ok {
			t.Fatalf("Statement is not FunctionLiteralExpression. Got %T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.expectedParameters) {
			t.Fatalf("Unexpected number of parameters, expected=%d got=%d", len(tt.expectedParameters), len(function.Parameters))
		}

		for i, p := range function.Parameters {
			testLiteralExpression(t, &p, tt.expectedParameters[i])
		}
	}
}
