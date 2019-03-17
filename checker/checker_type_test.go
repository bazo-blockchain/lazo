package checker

import (
	"github.com/bazo-blockchain/lazo/parser/node"
	"testing"
)

// Phase 4: Type Checker
// =====================

// Field Types
// -----------

func TestFieldBuiltInType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = true
		int x = 2
		char c = 'c'
		string s = "test"
	`, true)

	gs := tester.globalScope
	tester.assertField(0, gs.BoolType)
	tester.assertField(1, gs.IntType)
	tester.assertField(2, gs.CharType)
	tester.assertField(3, gs.StringType)
}

func TestFieldTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = 2
		int x = 'c'
		char c = "test"
		string s = true
	`, false)
	tester.assertTotalErrors(4)
}

// Local Variable Types
// --------------------

func TestLocalVarBuiltInType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(bool b1, int x1, char c1, string s1) {
			bool b = b1
			int x = x1
			char c = c1
			string s = s1
		}`, true)

	gs := tester.globalScope
	tester.assertLocalVariable(0, 0, gs.BoolType, 3)
	tester.assertLocalVariable(0, 1, gs.IntType, 2)
	tester.assertLocalVariable(0, 2, gs.CharType, 1)
	tester.assertLocalVariable(0, 3, gs.StringType, 0)
}

func TestLocalVarTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(bool b1, int x1, char c1, string s1) {
			bool b = x1
			int x = c1
			char c = s1
			string s = b1
		}`, false)
	tester.assertTotalErrors(4)
}

// Return Types
// ------------

func TestFunctionReturnVoid(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function void test() {
			return
		}
	`, true)
}

func TestFunctionReturnIntForVoid(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function void test() {
			return 1
		}
	`, false)
}

func TestFunctionReturnBoolConstant(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function bool test() {
			return true
		}
	`, true)
}

func TestFunctionReturnBoolFail(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function bool test() {
			return 5
		}
	`, false)
}

func TestFunctionReturnInt(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function int test() {
			int i = 5
			return 5
		}`, true)
}

func TestFunctionReturnString(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function string test() {
			string s = "test"
			return s
		}`, true)
}

func TestFunctionReturnChar(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function char test() {
			char c = 'c'
			return c
		}`, true)
}

func TestFunctionMultipleReturnTypes(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, char, bool) test() {
			return 1, 'c', true
		}`, true)

	returnStmt := tester.getFuncStatementNode(0,0).(*node.ReturnStatementNode)
	tester.assertExpressionType(returnStmt.Expressions[0], tester.globalScope.IntType)
	tester.assertExpressionType(returnStmt.Expressions[1], tester.globalScope.CharType)
	tester.assertExpressionType(returnStmt.Expressions[2], tester.globalScope.BoolType)
}

func TestFunctionReturnTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, char, bool) test() {
			return 'c', true, 1
		}`, false)
	tester.assertTotalErrors(3)
}

func TestFunctionMissingReturnValue(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, char, bool) test() {
			return 'c', true
		}`, false)
	tester.assertTotalErrors(1)
}

func TestFunctionTooManyReturnValues(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function int test() {
			return 1, true
		}`, false)
	tester.assertTotalErrors(1)
}
