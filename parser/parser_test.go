package parser

import (
	"fmt"
	"testing"

	"github.com/pwbrown/go-monkey/ast"
	"github.com/pwbrown/go-monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, 1)
		letStmt := testLetStatement(t, program.Statements[0], tt.expectedIdentifier)
		testLiteralExpression(t, letStmt.Value, tt.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, 1)
		retStmt := testReturnStatement(t, program.Statements[0])
		testLiteralExpression(t, retStmt.ReturnValue, tt.expectedValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := parseInput(t, input, 1)

	expStmt := testExpressionStatement(t, program.Statements[0])

	testIdentifier(t, expStmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program := parseInput(t, input, 1)

	expStmt := testExpressionStatement(t, program.Statements[0])

	testIntegerLiteral(t, expStmt.Expression, 5)
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		operator   string
		rightValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, 1)
		expStmt := testExpressionStatement(t, program.Statements[0])
		testPrefixExpression(t, expStmt.Expression, tt.operator, tt.rightValue)
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, 1)
		expStmt := testExpressionStatement(t, program.Statements[0])
		testInfixExpression(t, expStmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
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
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, -1)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, 1)
		expStmt := testExpressionStatement(t, program.Statements[0])
		testLiteralExpression(t, expStmt.Expression, tt.expectedBoolean)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := parseInput(t, input, 1)
	stmtExp := testExpressionStatement(t, program.Statements[0])
	ifExp := testIfExpression(t, stmtExp.Expression)
	testInfixExpression(t, ifExp.Condition, "x", "<", "y")
	cons := testBlockStatement(t, ifExp.Consequence, 1)
	innerStmtExp := testExpressionStatement(t, cons.Statements[0])
	testLiteralExpression(t, innerStmtExp.Expression, "x")

	if ifExp.Alternative != nil {
		t.Fatalf("ifExp.Alternative was not nil. got=%+v", ifExp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := parseInput(t, input, 1)
	stmtExp := testExpressionStatement(t, program.Statements[0])
	ifExp := testIfExpression(t, stmtExp.Expression)
	testInfixExpression(t, ifExp.Condition, "x", "<", "y")

	cons := testBlockStatement(t, ifExp.Consequence, 1)
	consStmtExp := testExpressionStatement(t, cons.Statements[0])
	testLiteralExpression(t, consStmtExp.Expression, "x")

	alt := testBlockStatement(t, ifExp.Alternative, 1)
	altStmtExp := testExpressionStatement(t, alt.Statements[0])
	testLiteralExpression(t, altStmtExp.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	program := parseInput(t, input, 1)
	stmtExp := testExpressionStatement(t, program.Statements[0])
	fn := testFunctionLiteral(t, stmtExp.Expression, 2)

	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	block := testBlockStatement(t, fn.Body, 1)
	blockStmt := testExpressionStatement(t, block.Statements[0])
	testInfixExpression(t, blockStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, 1)
		stmtExp := testExpressionStatement(t, program.Statements[0])
		lit := testFunctionLiteral(t, stmtExp.Expression, len(tt.expectedParams))
		for i, param := range tt.expectedParams {
			testLiteralExpression(t, lit.Parameters[i], param)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	program := parseInput(t, input, 1)
	expStmt := testExpressionStatement(t, program.Statements[0])
	call := testCallExpression(t, expStmt.Expression, 3)
	// function name
	testLiteralExpression(t, call.Function, "add")
	// arguments
	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 2, "*", 3)
	testInfixExpression(t, call.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		program := parseInput(t, tt.input, 1)
		expStmt := testExpressionStatement(t, program.Statements[0])
		call := testCallExpression(t, expStmt.Expression, len(tt.expectedArgs))
		testIdentifier(t, call.Function, tt.expectedIdent)

		for i, arg := range tt.expectedArgs {
			if call.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, call.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program := parseInput(t, input, 1)
	expStmt := testExpressionStatement(t, program.Statements[0])
	testStringLiteral(t, expStmt.Expression, "hello world")
}

// Test an individual let statement with a given name
func testLetStatement(t *testing.T, s ast.Statement, name string) *ast.LetStatement {
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return nil
	}

	if letStmt.TokenLiteral() != "let" {
		t.Errorf("letStmt.TokenLiteral not 'let'. got=%q", letStmt.TokenLiteral())
		return nil
	}

	testIdentifier(t, letStmt.Name, name)

	return letStmt
}

// Test an individual return statement
func testReturnStatement(t *testing.T, s ast.Statement) *ast.ReturnStatement {
	retStmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("s not *ast.LetStatement. got=%T", s)
	}

	if retStmt.TokenLiteral() != "return" {
		t.Fatalf("retStmt.TokenLiteral not 'return'. got=%q", retStmt.TokenLiteral())
	}

	return retStmt
}

// Test an expression statement
func testExpressionStatement(t *testing.T, s ast.Statement) *ast.ExpressionStatement {
	expStmt, ok := s.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("s not *ast.ExpressionStatement. got=%T", s)
	}

	return expStmt
}

// Test an expression statement
func testBlockStatement(t *testing.T, s ast.Statement, expectedStmts int) *ast.BlockStatement {
	block, ok := s.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("s not *ast.BlockStatement. got=%T", s)
	}

	if expectedStmts >= 0 && len(block.Statements) != expectedStmts {
		t.Fatalf("block.Statements does not have length %d. got=%d",
			expectedStmts, len(block.Statements))
	}

	return block
}

// Test a prefix expression
func testPrefixExpression(t *testing.T, e ast.Expression, operator string, right interface{}) *ast.PrefixExpression {
	exp, ok := e.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("e not *ast.PrefixExpression. got=%T", e)
	}

	if exp.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%s",
			operator, exp.Operator)
	}

	testLiteralExpression(t, exp.Right, right)

	return exp
}

// Test an infix expression
func testInfixExpression(t *testing.T, e ast.Expression, left interface{}, operator string, right interface{}) *ast.InfixExpression {
	exp, ok := e.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("e not *ast.InfixExpression. got=%T", e)
	}

	if exp.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%s",
			operator, exp.Operator)
	}

	testLiteralExpression(t, exp.Left, left)
	testLiteralExpression(t, exp.Right, right)

	return exp
}

// Test an in expression
func testIfExpression(t *testing.T, e ast.Expression) *ast.IfExpression {
	exp, ok := e.(*ast.IfExpression)
	if !ok {
		t.Fatalf("e not *ast.IfExpression. got=%T", e)
	}

	return exp
}

// Test a function literal
func testFunctionLiteral(t *testing.T, e ast.Expression, paramCount int) *ast.FunctionLiteral {
	lit, ok := e.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("e not *ast.FunctionLiteral. got=%T", e)
	}

	if paramCount >= 0 && len(lit.Parameters) != paramCount {
		t.Fatalf("lit.Parameters does not have length %d. got=%d",
			paramCount, len(lit.Parameters))
	}

	return lit
}

// Test a call expression
func testCallExpression(t *testing.T, e ast.Expression, argsCount int) *ast.CallExpression {
	exp, ok := e.(*ast.CallExpression)
	if !ok {
		t.Fatalf("e not *ast.CallExpression. got=%T", e)
	}

	if argsCount >= 0 && len(exp.Arguments) != argsCount {
		t.Fatalf("exp.Arguments length was not %d. got=%d",
			argsCount, len(exp.Arguments))
	}

	return exp
}

// Test a literal expression with an expected value
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) ast.Expression {
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
	t.Fatalf("type of exp not handled. got=%T", exp)
	return nil
}

// Test an identifier node with a specific name
func testIdentifier(t *testing.T, e ast.Expression, name string) *ast.Identifier {
	ident, ok := e.(*ast.Identifier)
	if !ok {
		t.Fatalf("e not *ast.Identifier. got=%T", e)
	}

	if ident.Value != name {
		t.Fatalf("ident.Value not '%s'. got=%s", name, ident.Value)
	}

	if ident.TokenLiteral() != name {
		t.Fatalf("ident.TokenLiteral() not '%s'. got=%s", name, ident.TokenLiteral())
	}

	return ident
}

// Test integer literal expression
func testIntegerLiteral(t *testing.T, e ast.Expression, value int64) *ast.IntegerLiteral {
	intLit, ok := e.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("e not *ast.IntegerLiteral. got=%T", e)
	}

	if intLit.Value != int64(value) {
		t.Fatalf("intLit.Value not '%d'. got=%d", value, intLit.Value)
	}

	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("intLit.TokenLiteral() not '%d'. got=%s", value, intLit.TokenLiteral())
		return nil
	}

	return intLit
}

// Test boolean literal expression
func testBoolean(t *testing.T, e ast.Expression, value bool) *ast.Boolean {
	boolean, ok := e.(*ast.Boolean)
	if !ok {
		t.Fatalf("e not *ast.Boolean. got=%T", e)
	}

	if boolean.Value != value {
		t.Fatalf("boolean.Value not '%t'. got=%t", value, boolean.Value)
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Fatalf("boolean.TokenLiteral() not '%t'. got=%s", value, boolean.TokenLiteral())
	}

	return boolean
}

// Test string literal expression
func testStringLiteral(t *testing.T, e ast.Expression, value string) *ast.StringLiteral {
	str, ok := e.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("e not *ast.StringLiteral. got=%T", e)
	}

	if str.Value != value {
		t.Fatalf("string.Value not '%s'. got=%s", value, str.Value)
	}

	if str.TokenLiteral() != value {
		t.Fatalf("str.TokenLiteral() not '%s'. got=%s", value, str.TokenLiteral())
	}

	return str
}

// Parse an input string, check for errors and statement lenth, and return program
func parseInput(t *testing.T, input string, expStmtLen int) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if expStmtLen >= 0 && len(program.Statements) != expStmtLen {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			expStmtLen, len(program.Statements))
	}
	return program
}

// Checks a parser for parsing errors
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d error(s)", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
