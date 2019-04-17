package parser

import (
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"testing"
)

// Program Nodes
// --------------

func TestEmptyProgram(t *testing.T) {
	p := newParserFromInput("")
	program, _ := p.ParseProgram()

	assertNoErrors(t, p)
	assertProgram(t, program, false)
}

func TestProgramWithNewlines(t *testing.T) {
	p := newParserFromInput("\n \n  \n \n")
	_, _ = p.ParseProgram()

	assertNoErrors(t, p)
}

func TestInvalidProgram(t *testing.T) {
	p := newParserFromInput("hello")
	_, _ = p.ParseProgram()

	assertHasError(t, p)
}

// Contract Nodes
// --------------

func TestEmptyContract(t *testing.T) {
	p := newParserFromInput(`
		contract Test {
		
		}
	`)
	program, _ := p.ParseProgram()

	assertNoErrors(t, p)
	assertProgram(t, program, true)
	assertContract(t, program.Contract, "Test", 0, 0)
}

func TestContractWithVariable(t *testing.T) {
	p := newParserFromInput(`contract Test {
		int x
		int y
	}`)
	c := p.parseContract()

	assertNoErrors(t, p)
	assertContract(t, c, "Test", 2, 0)
	assertField(t, c.Fields[0], "int", "x", "")
	assertField(t, c.Fields[1], "int", "y", "")

	// Positions
	assert.Equal(t, c.Pos().String(), "1:1")
	assert.Equal(t, c.Fields[0].Pos().String(), "2:3")
	assert.Equal(t, c.Fields[1].Pos().String(), "3:3")
}

func TestContractWithFunction(t *testing.T) {
	p := newParserFromInput(`contract Test {
		function void test() {

		}
	}`)
	c := p.parseContract()

	assertNoErrors(t, p)
	assertContract(t, c, "Test", 0, 1)
	assertFunction(t, c.Functions[0], "test", 1, 0, 0)
}

// Field Nodes
// --------------

func TestField(t *testing.T) {
	p := newParserFromInput("int x \n")
	f := p.parseField()

	assertField(t, f, "int", "x", "")
	assertNoErrors(t, p)
}

func TestFieldDeclarationWithoutNewLine(t *testing.T) {
	p := newParserFromInput("int x")
	_ = p.parseField()
	assertHasError(t, p)
}

func TestFieldAssignment(t *testing.T) {
	p := newParserFromInput("int x = 5\n")
	f := p.parseField()

	assertField(t, f, "int", "x", "5")
	assertNoErrors(t, p)
}

// Constructor Nodes
// -----------------

func TestEmptyConstructor(t *testing.T) {
	p := newParserFromInput("constructor() { \n } \n")
	c := p.parseConstructor()

	assertConstructor(t, c, 0, 0)
	assertNoErrors(t, p)
}

func TestConstructorWithParamsAndStatements(t *testing.T) {
	p := newParserFromInput(`constructor(int a, bool b) { 
		int x
		int y
		y = 2
	}
	`)
	c := p.parseConstructor()

	assertConstructor(t, c, 2, 3)
	assertNoErrors(t, p)
}

func TestMultipleConstructors(t *testing.T) {
	p := newParserFromInput(`contract Test {
		constructor(int a) {
		}
		
		constructor() {
		}
	}`)
	c := p.parseContract()

	assertConstructor(t, c.Constructor, 1, 0)
	assertErrorAt(t, p, 0, "Only one constructor is allowed")
}

// Function Nodes
// --------------

func TestEmptyFunction(t *testing.T) {
	p := newParserFromInput("function void test(){\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 1, 0, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithParam(t *testing.T) {
	p := newParserFromInput("function void test(int a){\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 1, 1, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithParams(t *testing.T) {
	p := newParserFromInput("function void test(int a, int b){\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 1, 2, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithMultipleRTypes(t *testing.T) {
	p := newParserFromInput("function (int, int) test(){\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 2, 0, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithParamsAndRTypes(t *testing.T) {
	p := newParserFromInput("function (int, int) test(int a, int b){\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 2, 2, 0)
	assertNoErrors(t, p)
}

func TestFunctionWithStatement(t *testing.T) {
	p := newParserFromInput("function void test(){\nint a\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 1, 0, 1)
	assertNoErrors(t, p)
}

func TestFunctionWithMultipleStatements(t *testing.T) {
	p := newParserFromInput("function void test(){\nint a\nint b\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 1, 0, 2)
	assertNoErrors(t, p)
}

func TestFullFunction(t *testing.T) {
	p := newParserFromInput("function (int, int) test(int a, int b){\nint a\nint b\n}\n")
	f := p.parseFunction()
	assertFunction(t, f, "test", 2, 2, 2)
	assertNoErrors(t, p)
}

func TestFunctionWORType(t *testing.T) {
	p := newParserFromInput("function test(int a, int b){\n}\n")
	p.parseFunction()
	assertHasError(t, p)
}

func TestFunctionMissingNewline(t *testing.T) {
	p := newParserFromInput("function void test(int a, int b){\n}")
	p.parseFunction()
	assertHasError(t, p)
}

func TestFunctionMissingNewlineInBody(t *testing.T) {
	p := newParserFromInput("function void test(int a, int b){}\n")
	p.parseFunction()
	assertHasError(t, p)
}

func TestFunctionMissingParamComma(t *testing.T) {
	p := newParserFromInput("function void test(int a int b){}\n")
	p.parseFunction()
	assertHasError(t, p)
}

// Statement Nodes
// ---------------

func TestEmptyStatementBlock(t *testing.T) {
	p := newParserFromInput("{\n}\n")
	v := p.parseStatementBlock()

	assertStatementBlock(t, v, 0)
	assertNoErrors(t, p)
}

func TestStatementBlock(t *testing.T) {
	p := newParserFromInput("{\nint a = 5\n}\n")
	v := p.parseStatementBlock()

	assertStatementBlock(t, v, 1)
	assertNoErrors(t, p)
}

func TestMultipleStatementBlock(t *testing.T) {
	p := newParserFromInput("{\nint a = 5\nint b = 4\n}\n")
	v := p.parseStatementBlock()

	assertStatementBlock(t, v, 2)
	assertNoErrors(t, p)
}

// Local Variable Statements
// -------------------------

func TestCharVariableStatement(t *testing.T) {
	p := newParserFromInput("char a = 'c'\n")
	v := p.parseVariableStatement().(*node.VariableNode)

	assertVariableStatement(t, v, "char", "a", "c")
	assertNoErrors(t, p)
}

func TestIntVariableStatement(t *testing.T) {
	p := newParserFromInput("int a = 5\n")
	v := p.parseVariableStatement().(*node.VariableNode)

	assertVariableStatement(t, v, "int", "a", "5")
	assertNoErrors(t, p)
}

func TestVariableStatementWONewline(t *testing.T) {
	p := newParserFromInput("char a = 'c'")
	p.parseVariableStatement()

	assertHasError(t, p)
}

// Multi Local Variable Statements
// -------------------------------

func TestMultiVariableStatement(t *testing.T) {
	p := newParserFromInput("int x, bool b = call() \n")
	mv, ok := p.parseVariableStatement().(*node.MultiVariableNode)

	assert.Assert(t, ok)
	assert.Equal(t, mv.Types[0].Identifier, "int")
	assert.Equal(t, mv.Identifiers[0], "x")
	assert.Equal(t, mv.Types[1].Identifier, "bool")
	assert.Equal(t, mv.Identifiers[1], "b")
	assertFuncCall(t, mv.FuncCall, "call")
	assertNoErrors(t, p)
}

func TestMultiVariableStatementMissingType(t *testing.T) {
	p := newParserFromInput("int x, = call() \n")
	_, ok := p.parseVariableStatement().(*node.MultiVariableNode)

	assert.Assert(t, ok)
	assertErrorAt(t, p, 0, "Identifier expected")
}

func TestMultiVariableStatementMissingID(t *testing.T) {
	p := newParserFromInput("int x, bool = call() \n")
	_, ok := p.parseVariableStatement().(*node.MultiVariableNode)

	assert.Assert(t, ok)
	assertErrorAt(t, p, 0, "Identifier expected")
}

func TestMultiVariableStatementMissingFuncCall(t *testing.T) {
	p := newParserFromInput("int x, bool b = y, true \n")
	_, ok := p.parseVariableStatement().(*node.MultiVariableNode)

	assert.Assert(t, ok)
	assertErrorAt(t, p, 0, "Symbol ( expected, but got ,")
}

func TestMultiVariableStatementMissingNewLine(t *testing.T) {
	p := newParserFromInput("int x, bool b = call(1, 2)")
	_, ok := p.parseVariableStatement().(*node.MultiVariableNode)

	assert.Assert(t, ok)
	assertErrorAt(t, p, 0, "Symbol \\n expected, but got EOF")
}

// Return statements
// ------------------

func TestReturnStatementMissingNewline(t *testing.T) {
	p := newParserFromInput("return")
	p.parseReturnStatement()
	assertHasError(t, p)
}

func TestEmptyReturnStatement(t *testing.T) {
	p := newParserFromInput("return \n")
	v := p.parseReturnStatement()

	assertReturnStatement(t, v, 0)
	assertNoErrors(t, p)
}

func TestSingleReturnStatement(t *testing.T) {
	p := newParserFromInput("return 1\n")
	v := p.parseReturnStatement()

	assertReturnStatement(t, v, 1)
	assertNoErrors(t, p)
}

func TestMultipleReturnStatement(t *testing.T) {
	p := newParserFromInput("return 1, 2\n")
	v := p.parseReturnStatement()

	assertReturnStatement(t, v, 2)
	assertNoErrors(t, p)
}

// If Statement
// ------------

func TestIfStatement(t *testing.T) {
	p := newParserFromInput("if(true){\n} else{\n}\n")
	v := p.parseIfStatement()

	assertIfStatement(t, v, "true", 0, 0)
	assertNoErrors(t, p)
}

func TestIfStatementSingleStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a \n} else{\nint b\n}\n")
	v := p.parseIfStatement()

	assertIfStatement(t, v, "true", 1, 1)
	assertNoErrors(t, p)
}

func TestIfStatementSingleThenStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a \n} else{\n}\n")
	v := p.parseIfStatement()

	assertIfStatement(t, v, "true", 1, 0)
	assertNoErrors(t, p)
}

func TestIfStatementSingleElseStatement(t *testing.T) {
	p := newParserFromInput("if(true){\n} else{\n int a \n}\n")
	v := p.parseIfStatement()

	assertIfStatement(t, v, "true", 0, 1)
	assertNoErrors(t, p)
}

func TestIfStatementMultipleStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a\n int b\n} else{\nint c\n int d\n}\n")
	v := p.parseIfStatement()

	assertIfStatement(t, v, "true", 2, 2)
	assertNoErrors(t, p)
}

func TestIfStatementMultipleThenStatement(t *testing.T) {
	p := newParserFromInput("if(true){\nint a\n int b\n} else{\n}\n")
	v := p.parseIfStatement()

	assertIfStatement(t, v, "true", 2, 0)
	assertNoErrors(t, p)
}

func TestIfStatementMultipleElseStatement(t *testing.T) {
	p := newParserFromInput("if(true){\n} else{\nint c\n int d\n}\n")
	v := p.parseStatement()

	assertIfStatement(t, v.(*node.IfStatementNode), "true", 0, 2)
	assertNoErrors(t, p)
}

func TestIfStatementWOElse(t *testing.T) {
	p := newParserFromInput("if(true){\n}\n")
	v := p.parseStatementWithFixToken()

	assertIfStatement(t, v.(*node.IfStatementNode), "true", 0, 0)
	assertNoErrors(t, p)
}

func TestIfStatementWOElseWONewline(t *testing.T) {
	p := newParserFromInput("if(true){\n}")
	p.parseIfStatement()

	assertHasError(t, p)
}

// Function Call Statements
// ------------------------

func TestFuncCallStatement(t *testing.T) {
	p := newParserFromInput("call() \n")
	s, ok := p.parseStatement().(*node.CallStatementNode)

	assertNoErrors(t, p)
	assert.Assert(t, ok)
	assertFuncCall(t, s.Call, "call")
}

func TestFuncCallStatementWithoutNL(t *testing.T) {
	p := newParserFromInput("call()")
	_ = p.parseStatement()

	assertErrorAt(t, p, 0, "Symbol \\n expected, but got EOF")
}

// Assignment
// ----------

func TestAssignmentStatement(t *testing.T) {
	p := newParserFromInput("a = 5\n")
	d := p.parseDesignator()
	v := p.parseAssignmentStatement(d)

	assertAssignmentStatement(t, v, "a", "5")
	assertNoErrors(t, p)
}

func TestAssignmentStatementChar(t *testing.T) {
	p := newParserFromInput("a = 'c'\n")
	s := p.parseStatementWithIdentifier()

	assertAssignmentStatement(t, s.(*node.AssignmentStatementNode), "a", "c")
	assertNoErrors(t, p)
}

func TestAssignmentStatementWONewline(t *testing.T) {
	p := newParserFromInput("a = 'c'")
	d := p.parseDesignator()
	_ = p.parseAssignmentStatement(d)

	assertHasError(t, p)
}

func TestAssignmentWithFuncCall(t *testing.T) {
	p := newParserFromInput("x = call() \n")
	s := p.parseStatement()

	assertNoErrors(t, p)
	assertAssignmentStatement(t, s.(*node.AssignmentStatementNode), "x", "call([])")
}

// Multi Assignment
// ----------------

func TestMultiAssignmentStatement(t *testing.T) {
	p := newParserFromInput("a, b = call(test(1)) \n")
	ma, ok := p.parseStatement().(*node.MultiAssignmentStatementNode)

	assert.Assert(t, ok)
	assertNoErrors(t, p)
	assert.Equal(t, ma.Designators[0].String(), "a")
	assert.Equal(t, ma.Designators[1].String(), "b")
	assertFuncCall(t, ma.FuncCall, "call", "test([1])")
}

func TestMultiAssignmentStatementMissingFuncCall(t *testing.T) {
	p := newParserFromInput("a, b = 1, 2 \n")
	_, ok := p.parseStatement().(*node.MultiAssignmentStatementNode)

	assert.Assert(t, ok)
	assertErrorAt(t, p, 0, "Identifier expected")
}

func TestMultiAssignmentStatementInvalidDesignator(t *testing.T) {
	p := newParserFromInput("a, int x = call() \n")
	_, ok := p.parseStatement().(*node.MultiAssignmentStatementNode)

	assert.Assert(t, ok)
	assertErrorAt(t, p, 0, "Symbol = expected, but got x")
}

// Statement with Fix token
// ------------------------

func TestStatementWithFixTokenReturn(t *testing.T) {
	p := newParserFromInput("return\n")
	v := p.parseStatementWithFixToken()

	assertStatement(t, v, "\n [1:1] RETURNSTMT []")
	assertNoErrors(t, p)
}

func TestStatementWithFixTokenReturnValue(t *testing.T) {
	p := newParserFromInput("return 5\n")
	v := p.parseStatementWithFixToken()

	assertStatement(t, v, "\n [1:1] RETURNSTMT [5]")
	assertNoErrors(t, p)
}

func TestStatementWithFixTokenMultipleReturnValue(t *testing.T) {
	p := newParserFromInput("return 5, 4\n")
	v := p.parseStatementWithFixToken()

	assertStatement(t, v, "\n [1:1] RETURNSTMT [5 4]")
	assertNoErrors(t, p)
}

// Statement with ID
// -------------------------

func TestStatementWithIdentifier(t *testing.T) {
	p := newParserFromInput("int a\n")
	v := p.parseStatementWithIdentifier()

	assertStatement(t, v, "\n [1:1] VAR int a")
	assertNoErrors(t, p)
}

func TestStatementWithIdentifierAssignment(t *testing.T) {
	p := newParserFromInput("int a = 5\n")
	v := p.parseStatementWithIdentifier()

	assertStatement(t, v, "\n [1:1] VAR int a = 5")
	assertNoErrors(t, p)
}

// Type Nodes
//-----------

func TestTypeNode(t *testing.T) {
	p := newParserFromInput("int")
	v := p.parseType()
	assertType(t, v, "int")
}
