package checker

import (
	"github.com/bazo-blockchain/lazo/parser/node"
	"gotest.tools/assert"
	"testing"
)

// Phase 3: Designator Resolution
// =============================

func TestUndefinedDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			x = 4
		}
	`, false)

	tester.assertTotalErrors(1)
}

// Field Designators
// -----------------

func TestFieldDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = 4
		int y = x
	`, true)

	tester.assertBasicDesignator(
		tester.syntaxTree.Contract.Fields[1].Expression,
		tester.globalScope.Contract.Fields[0],
		tester.globalScope.IntType)
}

func TestMixedDesignatorExpression(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = 4
		int y = 2 * x
	`, true)

	binExpr := tester.syntaxTree.Contract.Fields[1].Expression.(*node.BinaryExpressionNode)
	tester.assertBasicDesignator(
		binExpr.Right,
		tester.globalScope.Contract.Fields[0],
		tester.globalScope.IntType)
}

func TestUndefinedFieldDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int y = x
	`, false)
	tester.assertTotalErrors(1)
}

func TestFieldDesignatorInFunction(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		String s

		function void test() {
			String t = s
		}
	`, true)

	tester.assertBasicDesignator(
		tester.getFuncStatementNode(0, 0).(*node.VariableNode).Expression,
		tester.globalScope.Contract.Fields[0],
		tester.globalScope.StringType)
}

// Constructor Designators
// ------------------------

func TestConstructorDesignators(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x

		constructor(int a){
			int b = a
			int c = b
			int d = x
		}
	`, true)

	gs := tester.globalScope
	constructor := gs.Contract.Constructor

	tester.assertBasicDesignator(
		tester.getConstructorStatementNode(0).(*node.VariableNode).Expression,
		constructor.Parameters[0],
		gs.IntType)

	tester.assertBasicDesignator(
		tester.getConstructorStatementNode(1).(*node.VariableNode).Expression,
		constructor.LocalVariables[0],
		gs.IntType)

	tester.assertBasicDesignator(
		tester.getConstructorStatementNode(2).(*node.VariableNode).Expression,
		gs.Contract.Fields[0],
		gs.IntType)
}

func TestUndefinedConstructorDesignators(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		constructor(){
			int b = a
			c = 3
		}
	`, false)

	tester.assertErrorAt(0, "Designator a is undefined")
	tester.assertErrorAt(1, "Designator c is undefined")
}

// Function Parameter Designators
// ------------------------------

func TestFuncParamDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(bool a){
			bool b = a
		}
	`, true)

	tester.assertBasicDesignator(
		tester.getFuncStatementNode(0, 0).(*node.VariableNode).Expression,
		tester.globalScope.Contract.Functions[0].Parameters[0],
		tester.globalScope.BoolType)
}

func TestFuncParamInsideIf(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(bool a, char c){
			if (a) {
				char d = c
			}
		}
	`, true)

	ifStmt := tester.getFuncStatementNode(0, 0).(*node.IfStatementNode)
	tester.assertBasicDesignator(
		ifStmt.Condition,
		tester.globalScope.Contract.Functions[0].Parameters[0],
		tester.globalScope.BoolType)

	tester.assertBasicDesignator(
		ifStmt.Then[0].(*node.VariableNode).Expression,
		tester.globalScope.Contract.Functions[0].Parameters[1],
		tester.globalScope.CharType)
}

// Function Local Variable Designators
// -----------------------------------

func TestLocalVarDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			int x
			int y = x
		}
	`, true)

	tester.assertBasicDesignator(
		tester.getFuncStatementNode(0, 1).(*node.VariableNode).Expression,
		tester.getLocalVariableSymbol(0, 0),
		tester.globalScope.IntType)
}

func TestFuncNameAsLocalVarName(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			int test
			int y = test
		}
	`, true)

	tester.assertBasicDesignator(
		tester.getFuncStatementNode(0, 1).(*node.VariableNode).Expression,
		tester.getLocalVariableSymbol(0, 0),
		tester.globalScope.IntType)
}

func TestUndefinedLocalVarDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			int y = x
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestDesignatorWithAssignment(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			int x
			x = 3
		}
	`, true)

	tester.assertBasicDesignator(
		tester.getFuncStatementNode(0, 1).(*node.AssignmentStatementNode).Left,
		tester.getLocalVariableSymbol(0, 0),
		tester.globalScope.IntType)
}

func TestUndefinedLocalVarAssignment(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			x = 3
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestUndefinedMultiVarAssignment(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			x, y = test2()
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, false)

	tester.assertTotalErrors(2)
	tester.assertErrorAt(0, "Designator x is undefined")
	tester.assertErrorAt(1, "Designator y is undefined")
}

func TestUndefinedDesignatorAssignment(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			int x
			x = y
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestFuncNameAssignment(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			test = 3
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestContractNameAssignment(t *testing.T) {
	tester := newCheckerTestUtilWithRawInput(t, `
		contract Hello {
			function void test(){
				Hello = 3
			}
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestLocalVarAccessFromSubScope(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			bool b
			int x

			if (b) {
				x = 4
			}
		}
	`, true)

	ifStmt := tester.getFuncStatementNode(0, 2).(*node.IfStatementNode)
	tester.assertBasicDesignator(
		ifStmt.Condition,
		tester.getLocalVariableSymbol(0, 0),
		tester.globalScope.BoolType)

	tester.assertBasicDesignator(
		ifStmt.Then[0].(*node.AssignmentStatementNode).Left,
		tester.getLocalVariableSymbol(0, 1),
		tester.globalScope.IntType)
}

func TestLocalVarAccessOutOfScope(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			if (true) {
				int x
			}
			x = 4
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestLocalVarAccessIfElse(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(){
			if (true) {
				int x
			} else {
				x = 4
			}
		}
	`, false)
	tester.assertTotalErrors(1)
}

func TestLocalVarWithReturn(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function int test(){
			int x
			return x
		}
	`, true)

	returnStmt := tester.getFuncStatementNode(0, 1).(*node.ReturnStatementNode)
	tester.assertBasicDesignator(
		returnStmt.Expressions[0],
		tester.getLocalVariableSymbol(0, 0),
		tester.globalScope.IntType)
}

func TestUndefinedLocalVarReturn(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function int test(){
			return x
		}
	`, false)
	tester.assertTotalErrors(1)
}

// Function Call
// -------------

func TestUndefinedFuncCall(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int y = test()
	`, false)

	tester.assertErrorAt(0, "Designator test is undefined")
}

func TestDesignatorWithFuncCall(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x
		int y = test(x)

		function int test(int y) {
			return y
		}
	`, true)

	fc := tester.getFieldNode(1).Expression.(*node.FuncCallNode)
	assert.Equal(t, tester.symbolTable.GetDeclByDesignator(fc.Designator), tester.globalScope.Contract.Functions[0])
}

func TestUndefinedDesignatorWithFuncCall(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int y = test(x)

		function int test(int y) {
			return y
		}
	`, false)

	tester.assertErrorAt(0, "Designator x is undefined")
}

// Struct
// ------

func TestStructAssignment(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
		}

		Person p = new Person()
		Person p2 = p
	`, true)

	tester.assertBasicDesignator(
		tester.getFieldNode(1).Expression,
		tester.globalScope.Contract.Fields[0],
		tester.globalScope.Structs["Person"])
}

func TestStructFieldAccess(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		struct Person {
			int balance
		}

		Person p = new Person(1000)
		int i = p.balance
	`, true)

	memberAccess := tester.getFieldNode(1).Expression.(*node.MemberAccessNode)
	tester.assertMemberAccess(
		memberAccess,
		tester.globalScope.Structs["Person"].Fields[0],
		tester.globalScope.IntType)

	tester.assertBasicDesignator(
		memberAccess.Designator,
		tester.globalScope.Contract.Fields[0],
		tester.globalScope.Structs["Person"])
}
