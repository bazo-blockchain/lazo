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
		String s = "test"
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
		String s = true
	`, false)
	tester.assertTotalErrors(4)
}

// Local Variable Types
// --------------------

func TestLocalVarBuiltInType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(bool b1, int x1, char c1, String s1) {
			bool b = b1
			int x = x1
			char c = c1
			String s = s1
		}`, true)

	gs := tester.globalScope
	tester.assertFuncLocalVariable(0, 0, gs.BoolType, 3)
	tester.assertFuncLocalVariable(0, 1, gs.IntType, 2)
	tester.assertFuncLocalVariable(0, 2, gs.CharType, 1)
	tester.assertFuncLocalVariable(0, 3, gs.StringType, 0)
}

func TestLocalVarTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(bool b1, int x1, char c1, String s1) {
			bool b = x1
			int x = c1
			char c = s1
			String s = b1
		}`, false)
	tester.assertTotalErrors(4)
}

func TestConstructorLocalVars(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor(int a) {
			int b = a
			char c = a
		}
	`, false)

	gs := tester.globalScope
	constructor := gs.Contract.Constructor
	tester.assertLocalVariable(constructor.LocalVariables[0], constructor, gs.IntType, 1)
	tester.assertErrorAt(0, "expected Type char, given int")
}

// Return Types
// ------------

func TestConstructorReturn(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			return
		}
	`, false)
	tester.assertErrorAt(0, "return is not allowed in constructor")
}

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
		function String test() {
			String s = "test"
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

	returnStmt := tester.getFuncStatementNode(0, 0).(*node.ReturnStatementNode)
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

// Assignment Types
// ----------------

func TestFieldAssignmentType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x
		
		function void test() {
			x = 3
		}
	`, true)

	assignStmt := tester.getFuncStatementNode(0, 0).(*node.AssignmentStatementNode)
	tester.assertAssignment(assignStmt, tester.globalScope.IntType)
}

func TestFieldAssignmentTypeMismatch(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		bool b
		
		function void test() {
			b = 3
		}
	`, false)
}

// If Statement Types
// ------------------

func TestIfConditionBoolType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			if (true) {
			}
		}
	`, true)

	tester.assertExpressionType(
		tester.getFuncStatementNode(0, 0).(*node.IfStatementNode).Condition,
		tester.globalScope.BoolType)
}

func TestIfConditionIntType(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function void test() {
			if (1) {
			}
		}
	`, false)
}

// Binary Expression Types
// -----------------------

func TestLogicAndType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = true && true
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

func TestLogicOrTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = true || 1
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

func TestAdditionType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int a
		int b = 3 + a
	`, true)

	tester.assertExpressionType(tester.getFieldNode(1).Expression, tester.globalScope.IntType)
}

func TestSubtractionTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int a = true - 1
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
}

func TestMixedArithmeticExpr(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int a = 4 * 5 + 8 / 2 % 6
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
}

func TestEqualityComparisonType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool a = true == false
		bool b = 4 != 5
		bool c = 'c' == 'a'
		bool d = "hello" != "world"
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(1).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(2).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(3).Expression, tester.globalScope.BoolType)
}

func TestEqualityComparisonTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool a = true == 4
		bool b = 4 != false
		bool c = 'c' == "world"
		bool d = "hello" != 5
	`, false)
	tester.assertTotalErrors(4)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(1).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(2).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(3).Expression, tester.globalScope.BoolType)
}

func TestRelationalComparisonType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool a = 1 < 3
		bool b = 4 <= 4
		bool c = 'c' > 'a'
		bool d = 'c' >= 'a'
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(1).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(2).Expression, tester.globalScope.BoolType)
	tester.assertExpressionType(tester.getFieldNode(3).Expression, tester.globalScope.BoolType)
}

func TestRelationalComparisonTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool a = 1 < false
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

func TestBoolRelationalComparison(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool a = true > false
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

func TestStringRelationalComparison(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool a = "hello" >= "world"
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

// Unary Expression Types
// -----------------------

func TestUnaryPlusType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = +4
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
}

func TestUnaryMinusType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = -15
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
}

func TestUnaryTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = -true
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
}

func TestUnaryNotType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = !true
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

func TestUnaryNotTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = !4
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

func TestMixedExpressionType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = 1 > -2 == !false 
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

func TestMixedExpressionTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool b = 1 > -2 < 4 
	`, false)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.BoolType)
}

// Function Calls
// --------------

func TestFuncNameAsDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int y = test

		function int test() {
			return 1
		}
	`, false)

	tester.assertErrorAt(0, "Type mismatch: expected Type int, given nil")
}

func TestFuncNameAsLocalVar(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function int test() {
			int test2
			test2 = test2()
			return 1
		}

		function int test2() {
			return 1
		}
	`, false)

	tester.assertErrorAt(0, "test2 is not a function")
}

func TestFuncCallType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int y = test()
		bool b = test2()

		function int test() {
			return 1
		}

		function bool test2() {
			return true
		}
	`, true)

	tester.assertField(0, tester.globalScope.IntType)
	tester.assertField(1, tester.globalScope.BoolType)
}

func TestFuncCallArgsType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool y = test(1, true, "string")

		function bool test(int i, bool b, String s) {
			return true
		}
	`, true)

	fc := tester.getFieldNode(0).Expression.(*node.FuncCallNode)
	gs := tester.globalScope
	tester.assertExpressionType(fc.Args[0], gs.IntType)
	tester.assertExpressionType(fc.Args[1], gs.BoolType)
	tester.assertExpressionType(fc.Args[2], gs.StringType)
}

func TestFuncCallArgsMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool y = test(2)

		function bool test() {
		}
	`, false)

	tester.assertErrorAt(0, "expected 0 args, got 1")
}

func TestFuncCallArgsTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		bool y = test(true)

		function bool test(char c) {
		}
	`, false)

	tester.assertErrorAt(0, "expected Type char, got Type bool")
}

func TestVoidFuncCall(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			test2()
		}

		function void test2() {
		}
	`, true)

	st := tester.getFuncStatementNode(0, 0).(*node.CallStatementNode)
	tester.assertExpressionType(st.Call, nil)
}

func TestVoidFuncCallTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int y = test()

		function void test() {
		}
	`, false)

	tester.assertErrorAt(0, "expected 1 return value(s), but function returns 0")
}

func TestIntFuncCallAsStatement(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			test2()
		}

		function int test2() {
			return 1
		}
	`, false)

	st := tester.getFuncStatementNode(0, 0).(*node.CallStatementNode)
	tester.assertExpressionType(st.Call, tester.globalScope.IntType)

	tester.assertErrorAt(0, "function call as statement should be void")
}

func TestFuncCallBinary(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int y = test2() + 1
		}

		function int test2() {
			return 1
		}
	`, true)

	st := tester.getFuncStatementNode(0, 0).(*node.VariableNode)
	tester.assertExpressionType(st.Expression, tester.globalScope.IntType)
}

// Function Calls with multiple returns
// ------------------------------------

func TestFieldWithMultipleReturnValues(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = test()

		function (int, int) test() {
			return 1, 2
		}
	`, false)

	tester.assertErrorAt(0, "expected 1 return value(s), but function returns 2")
}

func TestLocalVarWithMultipleReturnValues(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int x = test2()
		}

		function (int, int) test2() {
			return 1, 2
		}
	`, false)

	tester.assertErrorAt(0, "expected 1 return value(s), but function returns 2")
}

func TestVoidFuncCallWithMultiVar(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int y, bool b = test2()
		}

		function void test2() {
			return
		}
	`, false)

	tester.assertErrorAt(0, "expected 2 return value(s), but function returns 0")
}

func TestFuncCallWithMultiVarInvalid(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int a, bool b = test2()
		}

		function int test2() {
			return
		}
	`, false)

	tester.assertErrorAt(0, "expected 2 return value(s), but function returns 1")
}

func TestFuncCallMultiVarTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int y, bool b = test2()
		}

		function (int, int) test2() {
			return 1, 2
		}
	`, false)

	tester.assertErrorAt(0, "Return type mismatch: expected int, given bool")
}

func TestMultiFuncCallBinary(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int y = test2() + 1
		}

		function (int, int) test2() {
			return 1
		}
	`, false)

	tester.assertErrorAt(0, "Arithmetic operators can only be applied to int types")
}

func TestMultiFuncCallReturn(t *testing.T) {
	_ = newCheckerTestUtil(t, `
		function (int, int, bool) test() {
        	return test2()
		}

    	function (int, int, bool) test2() {
        	return 1, 1, true
    	}
	`, true)
}

func TestMultiFuncCallReturnTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, int, bool) test() {
        	return test2()
		}

    	function (int, int, int) test2() {
        	return 1, 1, 1
    	}
	`, false)

	tester.assertErrorAt(0, "Return type mismatch: expected int, given bool")
}

func TestMultiFuncCallReturnMixed(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function (int, int) test() {
        	return test2(), 1 
		}

    	function (int, int) test2() {
        	return 1, 1
    	}
	`, false)

	tester.assertErrorAt(0, "Return type mismatch: expected int, given nil")
}
