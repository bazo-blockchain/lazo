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


