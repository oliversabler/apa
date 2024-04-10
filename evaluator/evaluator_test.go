package evaluator

import (
	"testing"

	"github.com/oliversabler/egglang/lexer"
	"github.com/oliversabler/egglang/object"
	"github.com/oliversabler/egglang/parser"
)

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!sant", false},
		{"!falskt", true},
		{"!5", false},
		{"!!sant", true},
		{"!!falskt", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
        låt nyAddering = funktion(x) {
            funktion(y) { x + y };
        };

        låt adderaTvå = nyAddering(2);
        adderaTvå(2);`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + sant;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + sant; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-sant",
			"unknown operator: -BOOLEAN",
		},
		{
			"sant + falskt;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; sant + falskt; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"om (10 > 1) { sant + falskt; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
            om (10 > 1) {
                om (10 > 1) {
                    tillbaka sant + falskt;
                }

                tillbaka 1;
            }
            `,
			"unknown operator: BOOLEAN + BOOLEAN",
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
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"sant", true},
		{"falskt", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"sant == sant", true},
		{"falskt == falskt", true},
		{"sant == falskt", false},
		{"sant != falskt", true},
		{"falskt != sant", true},
		{"(1 < 2) == sant", true},
		{"(1 < 2) == falskt", false},
		{"(1 > 2) == sant", false},
		{"(1 > 2) == falskt", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

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
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"om (sant) { 10 }", 10},
		{"om (falskt) { 10 }", nil},
		{"om (1) { 10 }", 10},
		{"om (1 < 2) { 10 }", 10},
		{"om (1 > 2) { 10 }", nil},
		{"om (1 > 2) { 10 } annars { 20 }", 20},
		{"om (1 < 2) { 10 } annars { 20 }", 10},
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

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"låt a = 5; a;", 5},
		{"låt a = 5 * 5; a;", 25},
		{"låt a = 5; låt b = a; b;", 5},
		{"låt a = 5; låt b = a; låt c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"tillbaka 10;", 10},
		{"tillbaka 10; 9;", 10},
		{"tillbaka 2 * 5; 9;", 10},
		{"9; tillbaka 2 * 5; 9;", 10},
		{`om (10 > 1) {
              om (10 > 1) {
                  tillbaka 10;
              }

              tillbaka 1;
          }`, 10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"låt identifiera = funktion(x) { x; }; identifiera(5);", 5},
		{"låt identifiera = funktion(x) { tillbaka x; }; identifiera(5);", 5},
		{"låt dubbel = funktion(x) { x * 2; }; dubbel(5);", 10},
		{"låt addera = funktion(x, y) { x + y; }; addera(5, 5);", 10},
		{"låt addera = funktion(x, y) { x + y; }; addera(5 + 5, addera(5, 5));", 20},
		{"funktion(x) { x; }(5);", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "funktion(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("wrong number of parameters, want 1. parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameters is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%v+)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object was wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%v+)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}
