package generator

import (
	"gotest.tools/assert"
	"math/big"
	"testing"
)

// Contract Fields
// ---------------

func TestContractFieldAssignment(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x

		function void test() {
			x = 3
		}
	`)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(3))
}

// Statements
// ----------

func TestLocalVarIntDefaultValue(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x
			return x
		}
	`)

	tester.assertInt(big.NewInt(0))
}

func TestLocalVarInt(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x = 3
			int y = 4
			return x
		}
	`)
	tester.assertInt(big.NewInt(3))

	tester = newGeneratorTestUtil(t, `
		function int test() {
			int x = 3
			int y = 4
			return y
		}
	`)
	tester.assertInt(big.NewInt(4))
}

func TestAssignmentInt(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x
			int y
			x = 3
			return x
		}
	`)
	tester.assertInt(big.NewInt(3))

	tester = newGeneratorTestUtil(t, `
		function int test() {
			int x
			int y
			x = 3
			return y
		}
	`)
	tester.assertInt(big.NewInt(0))
}

func TestReAssignmentInt(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x = 3
			int y = 4
			x = y
			return x
		}
	`)
	tester.assertInt(big.NewInt(4))
}

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

func TestLogicAndTrue(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return true && true
		}
	`)

	tester.assertBool(true)
}

func TestLogicAndFalse(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return true && false
		}
	`)

	tester.assertBool(false)
}

func TestLogicAndFalseShortCircuit(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return false && true
		}
	`)

	tester.assertBool(false)
}

func TestLogicAndFalseShortCircuit2(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return false && false
		}
	`)

	tester.assertBool(false)
}

func TestLogicOrFalse(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return false || false
		}
	`)

	tester.assertBool(false)
}

func TestLogicOrTrue(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return false || true
		}
	`)

	tester.assertBool(true)
}

func TestLogicOrTrueShortCircuit(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return true || false
		}
	`)

	tester.assertBool(true)
}

func TestLogicOrTrueShortCircuit2(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return true || true
		}
	`)

	tester.assertBool(true)
}

// TODO: Test all type of expressions
