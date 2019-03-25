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

func TestLocVarBoolDefautValue(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			bool x
			return x
		}
	`)

	tester.assertBool(false)
}

func TestLocVarBool(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			bool x = true
			return x
		}
	`)

	tester.assertBool(true)
}

func TestLocVarStringDefautValue(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function string test() {
			string x
			return x
		}
	`)

	tester.assertString("")
}

func TestLocVarString(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function string test() {
			string x = "hello"
			return x
		}
	`)

	tester.assertString("hello")
}

func TestLocVarCharDefautValue(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function char test() {
			char x
			return x
		}
	`)

	tester.assertChar('0')
}

func TestLocVarChar(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function char test() {
			char x = 'c'
			return x
		}
	`)

	tester.assertChar('c')
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

func TestReAssignmentBool(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			bool x = true
			bool y = false
			x = y
			return x
		}
	`)
	tester.assertBool(false)
}

func TestReAssignmentString(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function string test() {
			string x = "abc"
			string y = "def"
			x = y
			return x
		}
	`)
	tester.assertString("def")
}

func TestReAssignmentChar(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function char test() {
			char x = 'c'
			char y = 'd'
			x = y
			return x
		}
	`)
	tester.assertChar('d')
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

func TestSubtraction(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 2 - 1
		}
	`)

	tester.assertInt(big.NewInt(1))
}

func TestMultiplication(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 2 * 3
		}
	`)

	tester.assertInt(big.NewInt(6))
}

func TestDivision(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 10 / 5
		}
	`)

	tester.assertInt(big.NewInt(2))
}

func TestDivisionRound(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 5 / 2
		}
	`)

	tester.assertInt(big.NewInt(2))
}

func TestModulo(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 5 % 2
		}
	`)

	tester.assertInt(big.NewInt(1))
}

func TestExponent(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 2 ** 3
		}
	`)

	tester.assertInt(big.NewInt(8))
}

func TestNestedExponents(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 2 ** 2 ** 2
		}
	`)

	tester.assertInt(big.NewInt(16))
}

func TestMultipleExponent(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 2 ** 2 ** 2 ** 2
		}
	`)

	tester.assertInt(big.NewInt(65536))
}

func TestPointBeforeLine(t *testing.T){
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 8 - 4 * 2
		}
	`)

	tester.assertInt(big.NewInt(0))
}

func TestNegativeResult(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 1 - 2
		}
	`)

	tester.assertInt(big.NewInt(-1))
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

func TestLogicNot(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return !true
		}
	`)

	tester.assertBoolAfterNot(false)
}

func TestLogicNotNot(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return !!true
		}
	`)

	tester.assertBoolAfterNot(true)
}

// TODO: Test all type of expressions
