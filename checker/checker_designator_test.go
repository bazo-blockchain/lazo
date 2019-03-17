package checker

import "testing"

// Phase 3: Designator Resolution
// =============================

func TestNotDefinedDesignator(t *testing.T) {
	tester := newCheckerTestUtil(t, `
		function void test() {
			x = 4
		}
	`, false)

	tester.assertTotalErrors(1)
}

// Field Designators

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
