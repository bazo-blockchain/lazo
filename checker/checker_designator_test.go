package checker

import (
	"github.com/bazo-blockchain/lazo/parser/node"
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

	tester.assertDesignator(
		tester.syntaxTree.Contract.Variables[1].Expression,
		tester.globalScope.Contract.Fields[0],
		tester.globalScope.IntType)
}

func TestMixedDesignatorExpression(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		int x = 4
		int y = 2 * x
	`, true)

	binExpr := tester.syntaxTree.Contract.Variables[1].Expression.(*node.BinaryExpressionNode)
	tester.assertDesignator(
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
		string s

		function void test() {
			string t = s
		}
	`, true)

	tester.assertDesignator(
		tester.getFuncStatementNode(0, 0).(*node.VariableNode).Expression,
		tester.globalScope.Contract.Fields[0],
		tester.globalScope.StringType)
}

// Function Parameter Designators
// ------------------------------

func TestFuncParamDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test(bool a){
			bool b = a
		}
	`, true)

	tester.assertDesignator(
		tester.getFuncStatementNode(0,0).(*node.VariableNode).Expression,
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

	ifStmt := tester.getFuncStatementNode(0,0).(*node.IfStatementNode)
	tester.assertDesignator(
		ifStmt.Condition,
		tester.globalScope.Contract.Functions[0].Parameters[0],
		tester.globalScope.BoolType)

	tester.assertDesignator(
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

	tester.assertDesignator(
		tester.getFuncStatementNode(0,1).(*node.VariableNode).Expression,
		tester.getLocalVariableSymbol(0, 0),
		tester.globalScope.IntType)
}