package evaluator

import (
	"dara/lexer"
	"dara/parser"
	"testing"
)

func TestEvalNumberExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10.5", 10.5},
		{"-5", -5},
		{"-10.5", -10.5},
		{"6 % 2", 0},
		{"6 % 4", 2},
		{"6.2 % 4", 2.2},
		{"5.5 + 5.5", 11},
		{"5.5 * 2", 11},
		{"5.5 + 5.4", 10.9},
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
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"nil == nil", true},
		{`"a" == "a"`, true},
		{`"a" == "b"`, false},
		{`"a" != "b"`, true},
		{`"a" < "b"`, true},
		{`"a" > "b"`, false},
		{`"a" <= "b"`, true},
		{`"a" >= "b"`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"a"`, "a"},
		{`"a" + "b"`, "ab"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if true { 10 } else { 5 }", 10},
		{"if false { 10 } else { 5 }", 5},
		{"if true { 10 } else if false { 5 }", 10},
		{"if false { 10 } else if true { 5 }", 5},
		{"if false { 10 } else if false { 5 }", nil},
		{"if true { 10 } else if false { 5 } else { 3 }", 10},
		{"if false { 10 } else if true { 5 } else { 3 }", 5},
		{"if false { 10 } else if false { 5 } else { 3 }", 3},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestDefineExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"a := 5; a;", 5},
		{"a := 5 * 5; a;", 25},
		{"a := 5; b := a; b;", 5},
		{"a := 5; b := a; c := a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestAssignExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"a := 5; a := 25; a;", 25},
		{"a := 5; b := a; a = a + b + 5; a;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if 10 > 1 {
			if 10 > 1 {
				return 10;
			}
			return 1;
		}`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"identity := fn(x) { x; }; identity(5);", 5},
		{"identity := fn(x) { return x; }; identity(5);", 5},
		{"double := fn(x) { x * 2; }; double(5);", 10},
		{"add := fn(x, y) { x + y; }; add(5, 5);", 10},
		{"add := fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testNumberObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
   newAdder := fn(x) {
     fn(y) { x + y };
};
addTwo := newAdder(2); addTwo(2);`
	testNumberObject(t, testEval(input), 4)
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: number + boolean",
		},
		{
			`5 + "a";`,
			"type mismatch: number + string",
		},
		{
			"5 + true; 5;",
			"type mismatch: number + boolean",
		},
		{
			"-true",
			"invalid operation: operator - is not defined for true (boolean)",
		},
		{
			"!10",
			"invalid operation: operator ! is not defined for 10 (number)",
		},
		{
			`"a" - "b"`,
			"invalid operation: operator - is not defined for \"a\" (string)",
		},
		{
			"true + false;",
			"invalid operation: operator + is not defined for true (boolean)",
		},
		{
			"5; true + false; 5",
			"invalid operation: operator + is not defined for true (boolean)",
		},
		{
			"if 10 > 1 { true + false; }",
			"invalid operation: operator + is not defined for true (boolean)",
		},
		{
			"if 10 { true + false; }",
			"type mismatch: non-boolean condition 10 (number) in if statement",
		},
		{
			`if 10 > 1 {
				if 10 > 1 {
					return true + false;
				}
				return 1;
			}`,
			"invalid operation: operator + is not defined for true (boolean)",
		},
		{
			"foobar",
			"undeclared name: foobar",
		},
		{
			"foobar = 5",
			"undeclared name: foobar",
		},
		{
			"a := 5; 5()",
			"invalid operation: can not call non-function (number)",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func testNumberObject(t *testing.T, obj Object, expected float64) bool {
	result, ok := obj.(*Number)
	if !ok {
		t.Errorf("object is not Number. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%v, want=%v", result.Value, expected)
		return false
	}
	return true
}

func testNilObject(t *testing.T, obj Object) bool {
	if obj != NIL {
		t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj Object, expected string) bool {
	result, ok := obj.(*String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
		return false
	}
	return true
}

func testEval(input string) Object {
	var (
		l       = lexer.New(input)
		p       = parser.New(l)
		program = p.ParseProgram()
		env     = NewEnvironment()
	)
	return Eval(program, env)
}
