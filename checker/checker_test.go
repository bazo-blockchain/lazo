package checker

import (
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

// Function Symbol
//----------------

func TestVoidFunction(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function void test() {
		}
	`, true)

}

func TestFunctionReturnBoolConstant(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function bool test() {
				return true
			}
		}
		`, true)
}

func TestFunctionReturnBool(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function bool test() {
				bool b = true
				return b
			}
		}
		`, true)
}

// TODO: fix
//func TestFunctionReturnBoolFail(t *testing.T) {
//	_ = newCheckerTestUtilWithRawInput(t, `
//		contract Test {
//			function bool test() {
//				bool b = 5
//				return b
//			}
//		}
//		`, false)
//}

func TestFunctionReturnInt(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function int test() {
				int i = 5
				return 5
			}
		}
		`, true)
}

func TestFunctionReturnString(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function string test() {
				string s = "test"
				return s
			}
		}
		`, true)
}

func TestFunctionReturnChar(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function char test() {
				char c = 'c'
				return c
			}
		}
		`, true)
}



