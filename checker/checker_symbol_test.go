package checker

import (
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"testing"
)

// Phases 1 & 2: Symbol Construction & Type Resolution
// ===================================================

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

	gs := tester.globalScope
	assert.Check(t, gs.Contract != nil)
	assert.Equal(t, len(gs.Types), 4)
	assert.Equal(t, len(gs.BuiltInTypes), 4)
	assert.Equal(t, len(gs.BuiltInFunctions), 0)
	assert.Equal(t, len(gs.Constants), 3)

	// Built-in types
	assert.Equal(t, gs.NullType.Identifier(), "@NULL")
	assert.Equal(t, gs.BoolType.Identifier(), "bool")
	assert.Equal(t, gs.CharType.Identifier(), "char")
	assert.Equal(t, gs.StringType.Identifier(), "string")
	assert.Equal(t, gs.IntType.Identifier(), "int")

	// Constants
	assert.Equal(t, gs.TrueConstant.Identifier(), "true")
	assert.Equal(t, gs.FalseConstant.Identifier(), "false")
	assert.Equal(t, gs.NullConstant.Identifier(), "null")
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
	`, true)

	gs := tester.globalScope
	tester.assertField(0, gs.BoolType)
	tester.assertField(1, gs.IntType)
	tester.assertField(2, gs.CharType)
	tester.assertField(3, gs.StringType)
}

func TestUnknownFieldType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Integer l
	`, false)
	tester.assertTotalErrors(1)
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

func TestFunctionMultipleVoids(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (void, void) test() {
		}
	`, false)

	tester.assertTotalErrors(2)
	tester.assertFunction(0, 0, 0, 0)
}

func TestFunctionVoidInt(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (void, int) test() {
		}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertFunction(0, 1, 0, 0)
}

func TestFunctionIntVoid(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, void) test() {
		}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertFunction(0, 1, 0, 0)
}

func TestFunctionSingleReturnBool(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function bool test() {
			bool b = true
			return b
		}`, true)
	tester.assertFunction(0, 1, 0, 1)
	tester.assertReturnType(0, 0, tester.globalScope.BoolType)
}

func TestFunctionMultipleReturn(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, char, bool) test() {
			return 1, 'c', true
		}
	`, true)
	tester.assertFunction(0, 3, 0, 0)
	tester.assertReturnType(0, 0, tester.globalScope.IntType)
	tester.assertReturnType(0, 1, tester.globalScope.CharType)
	tester.assertReturnType(0, 2, tester.globalScope.BoolType)
}

func TestFunctionMaximumReturnTypes(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, char, bool, string) test() {
		}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertFunction(0, 4, 0, 0)
}

func TestFunctionSingleParam(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int x) {
		}
	`, true)
	tester.assertFunction(0, 0, 1, 0)
	tester.assertFuncParam(0, 0, tester.globalScope.IntType)
}

func TestFunctionMultipleParams(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int x, bool y) {
		}
	`, true)
	tester.assertFunction(0, 0, 2, 0)
	tester.assertFuncParam(0, 0, tester.globalScope.IntType)
	tester.assertFuncParam(0, 1, tester.globalScope.BoolType)
}

func TestFunctionWithLocalVars(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int x
			char y
		}
	`, true)
	tester.assertFunction(0, 0, 0, 2)
	tester.assertLocalVariable(0, 0, tester.globalScope.IntType, 1)
	tester.assertLocalVariable(0, 1, tester.globalScope.CharType, 0)
}

func TestFunctionWithMultipleVars(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			string s
			int x, bool y = test2()
			char z
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, true)

	tester.assertFunction(0, 0, 0, 4)
	tester.assertMultiLocalVariable(0, 1, 0, tester.globalScope.IntType, 1)
	tester.assertMultiLocalVariable(0, 2, 1, tester.globalScope.BoolType, 1)
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
	tester.assertLocalVariable(0, 0, tester.globalScope.IntType, 3)
	tester.assertLocalVariable(0, 1, tester.globalScope.BoolType, 2)
	tester.assertLocalVariable(0, 2, tester.globalScope.StringType, 1)

	varX := tester.globalScope.Contract.Functions[0].LocalVariables[0]
	assignX, ok := varX.VisibleIn[2].(*node.AssignmentStatementNode)
	assert.Assert(t, ok)
	assert.Equal(t, assignX.Left.Value, "x")
}

func TestFunctionWithMultiAssign(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			string a
			int x
			bool y
			x, y = test2()
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, true)

	varX := tester.globalScope.Contract.Functions[0].LocalVariables[1]
	varY := tester.globalScope.Contract.Functions[0].LocalVariables[2]

	multiAssign, ok := varX.VisibleIn[1].(*node.MultiAssignmentStatementNode)
	assert.Assert(t, ok)
	multiAssignY, ok := varY.VisibleIn[0].(*node.MultiAssignmentStatementNode)
	assert.Assert(t, ok)
	assert.Equal(t, multiAssign, multiAssignY)

	assert.Equal(t, multiAssign.Designators[0].Value, "x")
	assert.Equal(t, multiAssign.Designators[1].Value, "y")
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
	tester.assertLocalVariable(0, 0, tester.globalScope.IntType, 3)
	tester.assertLocalVariable(0, 1, tester.globalScope.BoolType, 1)
	tester.assertLocalVariable(0, 2, tester.globalScope.IntType, 0)

	varX := tester.globalScope.Contract.Functions[0].LocalVariables[0]
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
	tester.assertLocalVariable(0, 0, tester.globalScope.IntType, 9)
	tester.assertLocalVariable(0, 1, tester.globalScope.IntType, 7)
	tester.assertLocalVariable(0, 2, tester.globalScope.CharType, 1)
	tester.assertLocalVariable(0, 3, tester.globalScope.BoolType, 1)
	tester.assertLocalVariable(0, 4, tester.globalScope.StringType, 1)
}

// ID Checks
// -----------------

func TestInvalidFieldName(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool void
		int int
		char this
		string null
	`, false)
	tester.assertTotalErrors(4)
}

func TestInvalidFunctionName(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void int() {
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestInvalidParamName(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int int) {
		}
	`, false)
	tester.assertTotalErrors(1)
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
	tester.assertTotalErrors(3)
}

func TestInvalidMultiLocalVarNames(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int int, bool bool = test2()
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, false)

	tester.assertErrorAt(0, "Reserved keyword 'int' cannot be used")
	tester.assertErrorAt(1, "Reserved keyword 'bool' cannot be used")
}

func TestDuplicateFieldNames(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i
		bool i
	`, false)
	tester.assertTotalErrors(1)
}

func TestFieldVarShadowing(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		int i
		
		function void test(int i) {
		}

		function void test2() {
			int i
		}
	`, true)
}

func TestDuplicateLocalParamNames(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(int i) {
			int i
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestDuplicateLocalVars(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int i
			bool i
			char i
		}
	`, false)
	tester.assertTotalErrors(2)
}

func TestDuplicateMultiLocalVarNames(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int i, bool i = test2()
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Identifier 'i' is already declared")
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
	tester.assertTotalErrors(1)
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
	tester.assertTotalErrors(1)
}
