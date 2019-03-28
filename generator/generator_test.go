package generator

import (
	"gotest.tools/assert"
	"math/big"
	"testing"
)

// Contract Fields
// ---------------

func TestContractFieldDefault(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x

		function int test() {
			return x
		}
	`)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(0))
}

func TestContractFieldExpression(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x = 4 * 12

		function int test() {
			return x
		}
	`)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(48))
}

func TestMultipleContractFields(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x = 4 * 12
		int y = 3 * 12

		function int test() {
			return y
		}
	`)

	assert.Equal(t, tester.context.ContractVariables[1] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[1] == nil, false)
	tester.assertVariableInt(1, big.NewInt(36))
}

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

// Call contract functions externally
// ----------------------------------

func TestFuncCallByHash(t *testing.T) {
	txData := []byte{
		3,
		0x51, 0xA3, 0x52, 0xE1,
	}

	tester := newGeneratorTextUtilWithTx(t, `
		function int doNotCall() {
			return 4
		}

		function int doCall() {
			return 5
		}
	`, txData)

	tester.assertInt(big.NewInt(5))
}

func TestFuncCallByHashWithParams(t *testing.T) {
	txData := []byte{
		1, 0, 2,
		1, 0, 4,
		3, 0x35, 0x2E, 0x00, 0x80,
	}

	tester := newGeneratorTextUtilWithTx(t, `
		function int doNotCall() {
			return 4
		}

		function int doCall(int x, int y) {
			return x * y
		}
	`, txData)

	tester.assertInt(big.NewInt(8))
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

func TestReturnMultipleValuesSameTypes(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function (int, int) test() {
			int x = 1
			int y = 2
			return x, y
		}
	`)
	tester.assertInt(big.NewInt(2))
}

func TestReturnMultipleValuesDifferentTypes(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function (int, bool) test() {
			int x = 1
			bool y = true
			return x,y
		}
	`)
	tester.assertBool(true)
}

func TestIfStatement(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			if (true) {
				return 1
			}
			return 0
		}
	`)

	tester.assertInt(big.NewInt(1))
}

func TestSkipIfStatement(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			if (false) {
				return 1
			}
			return 0
		}
	`)

	tester.assertInt(big.NewInt(0))
}

func TestIfElseStatement(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			if (true) {
				return 1
			} else {
				return 0
			}
		}
	`)

	tester.assertInt(big.NewInt(1))
}

func TestIfElseStatementAlternative(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			if (false) {
				return 1
			} else {
				return 0
			}
		}
	`)

	tester.assertInt(big.NewInt(0))
}

func TestNestedIfStatement(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			if (true) {
				if (true) {
					return 1
				}
			} 
			return 0
		}
	`)

	tester.assertInt(big.NewInt(1))
}

func TestNestedIfStatementAlternative(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			if (true) {
				if (false) {
					return 1
				} else {
					return 2
				}
			} 
			return 0
		}
	`)

	tester.assertInt(big.NewInt(2))
}

func TestSingleReturnValue(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			return 1
		}
	`)

	tester.assertInt(big.NewInt(1))
}

func TestTwoReturnValues(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function (int, int) test() {
			return 1, 2
		}
	`)

	tester.assertInt(big.NewInt(2))
}

func TestThreeReturnValues(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function (int, int, int) test() {
			return 1, 2, 3
		}
	`)

	tester.assertInt(big.NewInt(3))
}

func TestSetter(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x = 4
	
		function void set() {
			x = 5
		}
	`)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(5))
}

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

func TestAdditionVar(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x = 1
			int y = 2
			return x + y
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

func TestSubtractionVar(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x = 1
			return 2 - x
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

func TestMultiplicationVar(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x = 2
			return x * 3
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

func TestDivisionVar(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
			int x = 5
			return 10 / x
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

//func TestExponent(t *testing.T) {
//	tester := newGeneratorTestUtil(t, `
//		function int test() {
//			return 2 ** 3
//		}
//	`)
//
//	tester.assertInt(big.NewInt(8))
//}
//
//// TODO: Fix exponent
//func TestExponentVar(t *testing.T) {
//	tester := newGeneratorTestUtil(t, `
//		function int test() {
//			int x = 3
//			return 2 ** x
//		}
//	`)
//
//	tester.assertInt(big.NewInt(8))
//}
//
//// TODO: Fix exponent (right associativity 2^9)
//func TestNestedExponents(t *testing.T) {
//	tester := newGeneratorTestUtil(t, `
//		function int test() {
//			return 2 ** 3 ** 2
//		}
//	`)
//
//	tester.assertInt(big.NewInt(512))
//}
//
//// TODO: Test with different basis
//func TestMultipleExponent(t *testing.T) {
//	tester := newGeneratorTestUtil(t, `
//		function int test() {
//			return 2 ** 2 ** 2 ** 2
//		}
//	`)
//
//	tester.assertInt(big.NewInt(65536))
//}

func TestPointBeforeLine(t *testing.T) {
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

	tester.assertInternalBool(false)
}

func TestLogicNotNot(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function bool test() {
			return !!true
		}
	`)

	tester.assertInternalBool(true)
}
