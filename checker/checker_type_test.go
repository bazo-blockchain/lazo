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
	tester.assertErrorAt(0, "expected char, given int")
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

// Shorthand Assignment Types
// ---------------------------

func TestPostfixIncAndDecrementType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			int x
			x++
			x--
		}
	`, true)

	assignment := tester.getConstructorStatementNode(1).(*node.ShorthandAssignmentStatementNode)
	tester.assertExpressionType(assignment.Designator, tester.globalScope.IntType)
	tester.assertExpressionType(assignment.Expression, tester.globalScope.IntType)

	assignment = tester.getConstructorStatementNode(2).(*node.ShorthandAssignmentStatementNode)
	tester.assertExpressionType(assignment.Designator, tester.globalScope.IntType)
	tester.assertExpressionType(assignment.Expression, tester.globalScope.IntType)
}

func TestPostfixIncAndDecrementTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			bool x
		}

		constructor() {
			Person p
			p.x++
			p.x--
		}
	`, false)

	tester.assertTotalErrors(2)
	tester.assertErrorAt(0, "expected Type int, got Type bool")
	tester.assertErrorAt(1, "expected Type int, got Type bool")
}

func TestShorthandAssignmentIntType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			int x
			x += 2
			x -= 2
			x *= 2
			x /= 2
			x **= 2
			x <<= 2
			x >>= 2
			x &= 2
			x |= 2
			x ^= 2
		}
	`, true)

	for i := 1; i <= 10; i++ {
		assignment := tester.getConstructorStatementNode(i).(*node.ShorthandAssignmentStatementNode)
		tester.assertExpressionType(assignment.Designator, tester.globalScope.IntType)
		tester.assertExpressionType(assignment.Expression, tester.globalScope.IntType)
	}
}

func TestShorthandAssignmentStringType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			String s = "hello"
			s += "World"
		}
	`, true)

	assignment := tester.getConstructorStatementNode(1).(*node.ShorthandAssignmentStatementNode)
	tester.assertExpressionType(assignment.Designator, tester.globalScope.StringType)
	tester.assertExpressionType(assignment.Expression, tester.globalScope.StringType)
}

func TestShorthandAssignmentTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			int x
			x += true

			String s
			s += x
		}
	`, false)

	tester.assertTotalErrors(2)
	tester.assertErrorAt(0, "expected Type int, got Type bool")
	tester.assertErrorAt(1, "expected Type String, got Type int")
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

// Ternary Expressions
// -------------------

func TestTernaryExpressionType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = 1 == 2 ? 1 : 2 
	`, true)

	ternaryExpr := tester.getFieldNode(0).Expression.(*node.TernaryExpression)
	tester.assertExpressionType(ternaryExpr, tester.globalScope.IntType)
	tester.assertExpressionType(ternaryExpr.Condition, tester.globalScope.BoolType)
	tester.assertExpressionType(ternaryExpr.True, tester.globalScope.IntType)
	tester.assertExpressionType(ternaryExpr.False, tester.globalScope.IntType)
}

func TestTernaryExpressionConditionTypeError(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = 1 + 2 ? 1 : 2 
	`, false)

	tester.assertErrorAt(0, "condition should be bool type")
}

func TestTernaryExpressionReturnTypeError(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = true ? 1 : false 
	`, false)

	tester.assertErrorAt(0, "ternary expression should return same type")
}

func TestTernaryExpressionAssignmentError(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = true ? true : false
	`, false)

	tester.assertErrorAt(0, "Type mismatch: expected int, given bool")
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
	tester.assertErrorAt(0, "Logic operators can only be applied to bool types")
}

func TestBitwiseLogicTypes(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = 1 & 2
		int x = i | 3
		int y = x ^ i
		int z = i & x | y ^ 4
	`, true)

	intType := tester.globalScope.IntType
	tester.assertExpressionType(tester.getFieldNode(0).Expression, intType)
	tester.assertExpressionType(tester.getFieldNode(1).Expression, intType)
	tester.assertExpressionType(tester.getFieldNode(2).Expression, intType)
	tester.assertExpressionType(tester.getFieldNode(3).Expression, intType)
}

func TestBitwiseLogicTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = true & 2
		int x = 'c' | false
		int y = x ^ "string"
	`, false)

	tester.assertTotalErrors(3)
	tester.assertErrorAt(0, "Bitwise logic operators can only be applied to int types")
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
	tester.assertErrorAt(0, "Arithmetic operators can only be applied to int types")
}

func TestMixedArithmeticExpr(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int a = 4 * 5 + 8 / 2 % 6
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
}

func TestStringConcatenation(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		String s = "hello" + "world"
		String s2 = "int " + (String) 1
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.StringType)
	tester.assertExpressionType(tester.getFieldNode(1).Expression, tester.globalScope.StringType)
}

func TestStringConcatenationError(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		String s = "hello" + 1
	`, false)

	tester.assertErrorAt(0, "+ operator can only be applied to int/string types")
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
	tester.assertErrorAt(0, "Equality comparison should have the same type")
	tester.assertErrorAt(1, "Equality comparison should have the same type")
	tester.assertErrorAt(2, "Equality comparison should have the same type")
	tester.assertErrorAt(3, "Equality comparison should have the same type")

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

func TestBitwiseShiftType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = 2 << 3
		int x = i >> 1
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
	tester.assertExpressionType(tester.getFieldNode(1).Expression, tester.globalScope.IntType)
}

func TestBitwiseShiftTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = 2 << true
		int x = 'c' >> 1
	`, false)

	tester.assertTotalErrors(2)
	tester.assertErrorAt(0, "Bitwise shift operators can only be applied to int types")
	tester.assertErrorAt(1, "Bitwise shift operators can only be applied to int types")
}

// Type cast expression
// --------------------

func TestTypeCastToString(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		String s = (String) 1
		String s2 = (String) true
		String s3 = (String) 'c'
		String s4 = (String) s
	`, true)

	typeCast := tester.getFieldNode(0).Expression.(*node.TypeCastNode)

	stringType := tester.globalScope.StringType
	tester.assertExpressionType(typeCast, stringType)
	tester.assertExpressionType(typeCast.Expression, tester.globalScope.IntType)
	tester.assertExpressionType(tester.getFieldNode(1).Expression, stringType)
	tester.assertExpressionType(tester.getFieldNode(2).Expression, stringType)
	tester.assertExpressionType(tester.getFieldNode(3).Expression, stringType)
}

func TestTypeCastToStringError(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] i

		String s = (String) i
	`, false)

	tester.assertTotalErrors(2)
	tester.assertErrorAt(0, "String type cast is not supported for Array of Type int")
	tester.assertExpressionType(tester.getFieldNode(1).Expression, nil)
}

func TestTypeCastUnsupportedType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = (int) true
	`, false)

	tester.assertTotalErrors(2)
	tester.assertErrorAt(0, "Unsupported type cast to Type int")
	tester.assertExpressionType(tester.getFieldNode(0).Expression, nil)
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

func TestBitwiseNotType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = ~4
	`, true)

	tester.assertExpressionType(tester.getFieldNode(0).Expression, tester.globalScope.IntType)
}

func TestBitwiseNotTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int i = ~true
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "~ unary operator can only be applied to int type")
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

	tester.assertErrorAt(0, "Type mismatch: expected int, given nil")
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

func TestVoidFuncCallArgTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			test2(true)
		}

		function void test2(int a) {
		}
	`, false)

	tester.assertErrorAt(0, "expected Type int, got Type bool")
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

func TestFuncCallMultiVarAssignmentTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			int y
			bool b
			y, b = test2()
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

	tester.assertErrorAt(0, "+ operator can only be applied to int/string types")
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

// Struct Types
// ------------

func TestStructCreationType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			String name
			int balance
		}

		Person p = new Person()
	`, true)

	gs := tester.globalScope
	tester.assertField(0, gs.Structs["Person"])
}

func TestStructCreationUndefinedType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Person p = new Person()
	`, false)

	tester.assertErrorAt(0, "Invalid type 'Person'")
}

func TestStructCreationUndefinedType2(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			String name
			int balance
		}

		Person p = new Person2()
	`, false)

	tester.assertErrorAt(0, "Struct Person2 is undefined")
}

func TestStructCreationUndefinedType3(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			String name
			int balance
		}

		Person p = new Person2(name="test")
	`, false)

	tester.assertErrorAt(0, "Struct Person2 is undefined")
}

func TestStructCreationTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
		}

		struct Person2 {
		}

		Person p = new Person2()
	`, false)

	tester.assertErrorAt(0, "expected Person, given Person2")
}

func TestStructCreationFieldType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			String name
			int balance
		}

		Person p = new Person("hello")
		Person p2 = new Person("hello", 120)
	`, true)

	gs := tester.globalScope
	tester.assertField(0, gs.Structs["Person"])

	sc := tester.getFieldNode(0).Expression.(*node.StructCreationNode)
	tester.assertExpressionType(sc.FieldValues[0], gs.StringType)

	sc = tester.getFieldNode(1).Expression.(*node.StructCreationNode)
	tester.assertExpressionType(sc.FieldValues[0], gs.StringType)
	tester.assertExpressionType(sc.FieldValues[1], gs.IntType)
}

func TestStructCreationFieldTypeMisMatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			String name
			int balance
		}

		Person p = new Person(120, "hello")
		Person p2 = new Person("hello", 120, true)
	`, false)

	tester.assertTotalErrors(3)
	tester.assertErrorAt(0, "expected Type String, got Type int")
	tester.assertErrorAt(1, "expected Type int, got Type String")
	tester.assertErrorAt(2, "Struct Person has only 3 field(s), got 2 value(s)")
}

func TestStructCreationWithNamedField(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			String name
			int balance
		}

		Person p = new Person(balance=120)
		Person p2 = new Person(balance=120, name="hello")
	`, true)

	gs := tester.globalScope
	tester.assertField(0, gs.Structs["Person"])

	sc := tester.getFieldNode(0).Expression.(*node.StructNamedCreationNode)
	tester.assertExpressionType(sc.FieldValues[0], gs.IntType)

	sc = tester.getFieldNode(1).Expression.(*node.StructNamedCreationNode)
	tester.assertExpressionType(sc.FieldValues[0], gs.IntType)
	tester.assertExpressionType(sc.FieldValues[1], gs.StringType)
}

func TestStructCreationWithNamedFieldTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			String name
			int balance
		}

		Person p = new Person(balance="120")
		Person p2 = new Person(age=18)
		Person p3 = new Person(balance=120, name="hello", age=18)
	`, false)

	tester.assertTotalErrors(3)
	tester.assertErrorAt(0, "expected Type int, got Type String")
	tester.assertErrorAt(1, "Field age not found")
	tester.assertErrorAt(2, "Struct Person has only 3 field(s), got 2 value(s)")
}

func TestStructFieldTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			int balance
		}
		
		constructor() {
			Person p = new Person()
			p.balance = true
			bool b = p.balance
		}
		
	`, false)

	tester.assertErrorAt(0, "assignment of bool is not compatible with target int")
	tester.assertErrorAt(1, "expected bool, given int")
}

func TestThisReturnStatement(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			int x = test()
		}

		function Test test() {
			return this
		}
	`, false)

	tester.assertErrorAt(0, "Invalid type 'Test'")
}

func TestAssignToThisMember(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = 0
		constructor() {
			this.x = 5
		}
	`, true)

	assignment := tester.getConstructorStatementNode(0).(*node.AssignmentStatementNode)

	tester.assertAssignment(assignment, tester.globalScope.IntType)

}

func TestAssignToThis(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			this = this
		}
	`, false)

	tester.assertErrorAt(0, "Assigning to 'this' is not allowed!")
}

func TestAssignThisMemberToVar(t *testing.T) {
	tester := newCheckerTestUtil(t, `
	int x = 5
	int y = 0
	constructor() {
		y = this.x
	}
`, true)

	assignment := tester.getConstructorStatementNode(0).(*node.AssignmentStatementNode)

	tester.assertAssignment(assignment, tester.globalScope.IntType)
}

func TestAssignThisToVar(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x
		constructor() {
			x = this
		}
	`, false)

	tester.assertErrorAt(0, "'this' cannot be assigned!")
}

func TestThisVariableDeclaration(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			int this = 0
		}
	`, false)
	tester.assertErrorAt(0, "Reserved keyword 'this' cannot be used as an identifier")
}

func TestThisFieldDeclaration(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int this
	`, false)
	tester.assertErrorAt(0, "Reserved keyword 'this' cannot be used as an identifier")
}

func TestThisParameter(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			int x = test(this)
		}

		function int test(int x) {
			return x
		}
	`, false)
	tester.assertErrorAt(0, "'this' cannot be used as an argument")
}

// Arrays
// ------

func TestArrayInitialization(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new int[1]
	`, true)

	tester.assertField(0, tester.globalScope.Types["int[]"])
	tester.assertArrayLengthCreation(tester.getFieldNode(0).Expression, tester.globalScope.Types["int[]"])
}

func TestArrayInitializationWithValues(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new int[]{1, 2, 3}
	`, true)

	tester.assertField(0, tester.globalScope.Types["int[]"])
	tester.assertArrayValueCreation(tester.getFieldNode(0).Expression, tester.globalScope.Types["int[]"], tester.globalScope.IntType)
}

func TestArrayInitializationWithVariableLength(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = 2
		int[] a = new int[x]
	`, true)

	creation := tester.getFieldNode(1).Expression.(*node.ArrayLengthCreationNode)
	tester.assertArrayLengthCreation(creation, tester.globalScope.Types["int[]"])
}

func TestInvalidArrayInitialization(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new char[1]
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Type mismatch: expected int[], given char[]")
}

func TestInvalidArrayInitializationWithValues(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new char[]{'a', 'b', 'c'}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Type mismatch: expected int[], given char[]")
}

func TestArrayInitializationWithValuesOfDifferentType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new int[]{'a', 2, 3}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "expected Type int, got Type char")
}

func TestArrayNestedLengthInitialization(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[][] a = new int[1][2]
	`, true)

	tester.assertTotalErrors(0)
}

func TestArrayNestedValueInitialization(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[][] a = new int[][]{{1, 2}, {3, 4}}
	`, true)
	tester.assertTotalErrors(0)
}

func TestArrayNestedValueInitializationDifferentLength(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[][] a = new int[][]{{1, 2}, {3}}
	`, true)

	tester.assertTotalErrors(0)
}

func TestInvalidNestedArrayAssignment1(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new int[2][2]
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Type mismatch: expected int[], given int[][]")
}

func TestInvalidNestedArrayAssignment2(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new int[][]{{1, 2}, {3}}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Type mismatch: expected int[], given int[][]")
}

func TestInvalidArrayElementAccess(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int[] a = new int[1]
		constructor() {
			a[true] = 1
		}
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Array index must be of type int")
}

// Maps
// ----

func TestMapTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<String, int> m
		Map<int, int> m2 = m
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Type mismatch: expected Map<int,int>, given Map<String,int>")
}

func TestMapElementAccessType(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<String, int> m
		int i = m["key"]
	`, true)

	tester.assertField(1, tester.globalScope.IntType)
}

func TestMapElementAccessTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<String, int> m
		String s = m["key"]
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "Type mismatch: expected String, given int")
}

func TestMapKeyTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<String, int> m
		int i = m[1]
	`, false)

	tester.assertTotalErrors(1)
	tester.assertErrorAt(0, "expected Type String, got Type int")
}

func TestMapElementAssignmentTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor(){
			Map<bool, int> m
			m[true] = "string"
		}
	`, false)

	tester.assertErrorAt(0, "assignment of String is not compatible with target int")
}

func TestMapNestedElementAccess(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<int, int[]> m
		int i = m[1][0]
	`, true)

	tester.assertField(1, tester.globalScope.IntType)
}

func TestMapNestedElementAccessError(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<int, int> m
		int i = m[1][0]
	`, false)

	tester.assertErrorAt(0, "Designator m[1][0] does not refer to an array/map type")
}

func TestDeleteStatement(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			Map<char, int> m
			delete m['c']
		}
	`, true)

	deleteStmt := tester.getConstructorStatementNode(1).(*node.DeleteStatementNode)
	tester.assertExpressionType(deleteStmt.Element, tester.globalScope.IntType)
}

func TestInvalidDeleteStatement(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor() {
			int[] i = new int[2]
			delete i[0]
		}
	`, false)

	tester.assertErrorAt(0, "delete requires map type")
}

func TestMapContainsKey(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<char, int> m
		bool b = m.contains('c')
	`, true)

	tester.assertExpressionType(tester.getFieldNode(1).Expression, tester.globalScope.BoolType)
}

func TestMapContainsKeyTypeMismatch(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		Map<char, int> m
		bool b = m.contains("c")
	`, false)

	tester.assertErrorAt(0, "expected Type char, got Type String")
}
