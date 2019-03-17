package checker

import "testing"

// Phase 3: Designator Resolution
// =============================

func TestNotDefinedDesignator(t *testing.T){
	tester := newCheckerTestUtil(t, `
		function void test() {
			x = 4
		}
	`, false)

	tester.assertTotalErrors(1)
}

// TODO add more tests
