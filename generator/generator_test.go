package generator

import (
	"math/big"
	"testing"
)

// Statements
// ----------

// TODO: Test if, assignment, local variable and return statements

// Expressions
// -----------

func TestAddition(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
        	return 1 + 2
		}
	`)

	tester.assertInt(big.NewInt(3))
}

// TODO: Test all type of expressions