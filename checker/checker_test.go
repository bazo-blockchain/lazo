package checker

import (
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"testing"
)

// Phase: Symbol Construction
// ==========================

func TestEmptyProgram(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, ``, false)
}

// Global Scope
// ------------

func TestGlobalScope(t *testing.T) {
	tester := newCheckerTestUtilWithRawInput(t, `
		contract Test {
		}
	`, true)

	gs := tester.symbolTable.GlobalScope
	assert.Check(t, gs.Contract != nil)
	assert.Equal(t, len(gs.Types), 4)
	assert.Equal(t, len(gs.BuiltInTypes), 4)
	assert.Equal(t, len(gs.BuiltInFunctions), 0)
	assert.Equal(t, len(gs.Constants), 3)

	// Built-in types
	assert.Equal(t, gs.NullType.GetIdentifier(), "@NULL")
	assert.Equal(t, gs.BoolType.GetIdentifier(), "bool")
	assert.Equal(t, gs.CharType.GetIdentifier(), "char")
	assert.Equal(t, gs.StringType.GetIdentifier(), "string")
	assert.Equal(t, gs.IntType.GetIdentifier(), "int")

	// Constants
	assert.Equal(t, gs.TrueConstant.GetIdentifier(), "true")
	assert.Equal(t, gs.FalseConstant.GetIdentifier(), "false")
	assert.Equal(t, gs.NullConstant.GetIdentifier(), "null")
}

// Contract Symbol
// ---------------

func TestValidContract(t *testing.T) {
	tester := newCheckerTestUtilWithRawInput(t, `
		contract Test {
			int x
			bool b = true
			char c = 'c'
			string s = "hello"

			function void test() {
				int x = 2
			}

			function void test2(bool x) {
				bool b = x
			}
		}
	`, true)

	tester.assertContract(4, 2)
}

func TestInvalidContractName(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract int {
		}
	`, false)
}

// Field Symbol
// ------------

func TestContractFields(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b
		int x = 2
		char c
		string s
	`,true)

	gs := tester.symbolTable.GlobalScope
	tester.assertField(0, gs.BoolType)
	tester.assertField(1, gs.IntType)
	tester.assertField(2, gs.CharType)
	tester.assertField(3, gs.StringType)
}

// Function Symbol with parameter and local variable symbols
//----------------------------------------------------------

func TestFunctionVoid(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
		}
	`, true)

	tester.assertFunction(0, 0, 0, 0)
}

func TestFunctionSingleReturnBool(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function bool test() {
			bool b = true
			return b
		}`, true)
	tester.assertFunction(0, 1, 0, 1)
}

func TestFunctionMultipleReturn(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, bool) test() {
			return 1, false
		}
	`, true)
	tester.assertFunction(0, 2, 0, 0)
}

// TODO Return type tests

func TestFunctionSingleParam(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int x) {
		}
	`, true)
	tester.assertFunction(0, 0, 1, 0)
	tester.assertFuncParam(0, 0, tester.symbolTable.GlobalScope.IntType)
}

func TestFunctionMultipleParams(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int x, bool y) {
		}
	`, true)
	tester.assertFunction(0, 0, 2, 0)
	tester.assertFuncParam(0, 0, tester.symbolTable.GlobalScope.IntType)
	tester.assertFuncParam(0, 1, tester.symbolTable.GlobalScope.BoolType)
}

func TestFunctionWithLocalVars(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int x
			char y
		}
	`, true)
	tester.assertFunction(0, 0, 0, 2)
	tester.assertLocalVariable(0, 0, tester.symbolTable.GlobalScope.IntType, 0)
	tester.assertLocalVariable(0, 1, tester.symbolTable.GlobalScope.CharType, 0)
}

func TestFunctionWithAssignment(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int x
			bool b
			string s
			x = 2
		}
	`, true)
	tester.assertFunction(0, 0, 0, 3)
	tester.assertLocalVariable(0, 0, tester.symbolTable.GlobalScope.IntType, 1)
	tester.assertLocalVariable(0, 1, tester.symbolTable.GlobalScope.BoolType, 1)
	tester.assertLocalVariable(0, 2, tester.symbolTable.GlobalScope.StringType, 1)

	varX := tester.symbolTable.GlobalScope.Contract.Functions[0].LocalVariables[0]
	assignX, ok := varX.VisibleIn[0].(*node.AssignmentStatementNode)
	assert.Assert(t, ok)
	assert.Equal(t, assignX.Left.Value, "x")
}

func TestFunctionWithIf(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int x
			
			if (true) {
				bool b
				int y
			}	
		}
	`, true)

	tester.assertFunction(0, 0, 0, 3)
	tester.assertLocalVariable(0, 0, tester.symbolTable.GlobalScope.IntType, 1)
	tester.assertLocalVariable(0, 1, tester.symbolTable.GlobalScope.BoolType, 0)
	tester.assertLocalVariable(0, 2, tester.symbolTable.GlobalScope.IntType, 0)

	varX := tester.symbolTable.GlobalScope.Contract.Functions[0].LocalVariables[0]
	_, ok := varX.VisibleIn[0].(*node.IfStatementNode)
	assert.Assert(t, ok)
}

func TestFunctionWithIfElse(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function int test() {
			int x
			x = 3
			int y = 4
			
			if (true) {
				char c
				c = 'c'
			} else {
				bool a
				a = true
			}
			
			string s
			return x
		}
	`, true)

	tester.assertFunction(0, 1, 0, 5)
	tester.assertLocalVariable(0, 0, tester.symbolTable.GlobalScope.IntType, 5)
	tester.assertLocalVariable(0, 1, tester.symbolTable.GlobalScope.IntType, 4)
	tester.assertLocalVariable(0, 2, tester.symbolTable.GlobalScope.CharType, 1)
	tester.assertLocalVariable(0, 3, tester.symbolTable.GlobalScope.BoolType, 1)
	tester.assertLocalVariable(0, 4, tester.symbolTable.GlobalScope.StringType, 1)
}

// Identifier Checks
// -----------------

func TestInvalidFieldName(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool void
		int int
		char this
		string null
	`,false)
	assert.Equal(t, len(tester.errors), 4)
}

func TestInvalidFunctionName(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void int() {
		}
	`,false)
	assert.Equal(t, len(tester.errors), 1)
}

func TestInvalidParamName(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int int) {
		}
	`,false)
	assert.Equal(t, len(tester.errors), 1)
}

func TestInvalidLocalVarNames(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			bool bool
			
			if(true) {
				char char
			} else {
				string string
			}
		}
	`, false)
	assert.Equal(t, len(tester.errors), 3)
}

func TestDuplicateFieldNames(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i
		bool i
	`,false)
	assert.Equal(t, len(tester.errors), 1)
}

func TestFieldVarShadowing(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		int i
		
		function void test(int i) {
		}

		function void test2() {
			int i
		}
	`,true)
}

func TestDuplicateLocalParamNames(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int i) {
			int i
		}
	`, false)
	assert.Equal(t, len(tester.errors), 1)
}

func TestDuplicateLocalVars(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int i
			bool i
			char i
		}
	`, false)
	assert.Equal(t, len(tester.errors), 2)
}

func TestLocalVarShadowing(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int i
			
			if(true) {
				bool i
			}
		}
	`, false)
	assert.Equal(t, len(tester.errors), 1)
}

func TestUniqueLocalVar(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {			
			if(true) {
				bool i
			} else {
				string i
			}
		}
	`, false)
	assert.Equal(t, len(tester.errors), 1)
}
