package generator

import (
	"math/big"
	"testing"
)

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
			int x = 4
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
