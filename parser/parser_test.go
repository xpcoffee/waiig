package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

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

	idnt, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("Statement is not an Identifier. Got %T", program.Statements[0])
	}

	if idnt.Value != "foobar" {
		t.Errorf("Unexpected value. Expected %q, Got %q", "foobar", idnt.Value)
	}

	if idnt.TokenLiteral() != "foobar" {
		t.Errorf("Unexpected TokenLiteral. Expected %q, Got %q", "foobar", idnt.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not an expression. Got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Statement is not an IntegerLiteral. Got %T", program.Statements[0])
	}

	if literal.Value != 5 {
		t.Errorf("Unexpected value. Expected %d, Got %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("Unexpected TokenLiteral. Expected %s, Got %s", "5", literal.TokenLiteral())
	}

}

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
