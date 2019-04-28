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
	`)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(0))
}

func TestContractFieldExpression(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x = 4 * 12
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
	`)

	assert.Equal(t, tester.context.ContractVariables[1] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[1] == nil, false)
	tester.assertVariableInt(1, big.NewInt(36))
}

// Constructor
// -----------

func TestContractFieldAssignment(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x

		constructor(){
			x = 3
		}
	`)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(3))
}

func TestConstructorWithParam(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int a

		constructor(int x){
			int y = x + 1
			a = y
		}
	`, 2, 0, 4)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(5))
}

// CallFunc contract functions externally
// ----------------------------------

func TestFuncCallByHash(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int doNotCall() {
			return 4
		}

		function int doCall() {
			return 5
		}
	`, "(int)doCall()")

	tester.assertInt(big.NewInt(5))
}

func TestFuncCallByHashWithParams(t *testing.T) {
	funcData := []byte{
		2, 0, 2,
		2, 0, 4,
	}

	tester := newGeneratorTestUtilWithFunc(t, `
		function int doNotCall() {
			return 4
		}

		function int doCall(int x, int y) {
			return x * y
		}
	`, "(int)doCall(int,int)", funcData...)

	tester.assertInt(big.NewInt(8))
}

// Statements
// ----------

// Local Variables
// ---------------

func TestLocalVarIntDefaultValue(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x
			return x
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(0))
}

func TestLocalVarInt(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 3
			int y = 4
			return x
		}
	`, intTestSig)
	tester.assertInt(big.NewInt(3))

	tester = newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 3
			int y = 4
			return y
		}
	`, intTestSig)
	tester.assertInt(big.NewInt(4))
}

func TestLocVarBoolDefaultValue(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function bool test() {
			bool x
			return x
		}
	`, boolTestSig)

	tester.assertBool(false)
}

func TestLocVarBool(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function bool test() {
			bool x = true
			return x
		}
	`, boolTestSig)

	tester.assertBool(true)
}

func TestLocVarStringDefaultValue(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function String test() {
			String x
			return x
		}
	`, stringTestSig)

	tester.assertString("")
}

func TestLocVarString(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function String test() {
			String x = "hello"
			return x
		}
	`, stringTestSig)

	tester.assertString("hello")
}

func TestLocVarCharDefaultValue(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function char test() {
			char x
			return x
		}
	`, charTestSig)

	tester.assertChar('0')
}

func TestLocVarChar(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function char test() {
			char x = 'c'
			return x
		}
	`, charTestSig)

	tester.assertChar('c')
}

// Multi-Variables
// ---------------

func TestMultiVariables(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, bool) test() {
			int x, bool b = test2()
			return x, b
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, "(int,bool)test()")

	tester.assertIntAt(0, big.NewInt(1))
	tester.assertBoolAt(1, true)
}

// Assignments
// -----------

func TestAssignmentInt(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x
			int y
			x = 3
			return x
		}
	`, intTestSig)
	tester.assertInt(big.NewInt(3))

	tester = newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x
			int y
			x = 3
			return y
		}
	`, intTestSig)
	tester.assertInt(big.NewInt(0))
}

func TestReAssignmentInt(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 3
			int y = 4
			x = y
			return x
		}
	`, intTestSig)
	tester.assertInt(big.NewInt(4))
}

func TestReAssignmentBool(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function bool test() {
			bool x = true
			bool y = false
			x = y
			return x
		}
	`, boolTestSig)
	tester.assertBool(false)
}

func TestReAssignmentString(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function String test() {
			String x = "abc"
			String y = "def"
			x = y
			return x
		}
	`, stringTestSig)
	tester.assertString("def")
}

func TestReAssignmentChar(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function char test() {
			char x = 'c'
			char y = 'd'
			x = y
			return x
		}
	`, charTestSig)
	tester.assertChar('d')
}

// Multi-Assignments
// -----------------

func TestMultiAssignment(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, bool) test() {
			int x
			bool b
			x, b = test2()
			return x, b
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, "(int,bool)test()")

	tester.assertIntAt(0, big.NewInt(1))
	tester.assertBoolAt(1, true)
}

func TestMultiAssignmentWithField(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		int x
		function (int, bool) test() {
			bool b
			x, b = test2()
			return x, b
		}

		function (int, bool) test2() {
			return 1, true
		}
	`, "(int,bool)test()")

	tester.assertIntAt(0, big.NewInt(1))
	tester.assertBoolAt(1, true)
}

// Return statements
// -----------------

func TestSingleReturnValue(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			return 1
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1))
}

func TestTwoReturnValues(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, int) test() {
			return 1, 2
		}
	`, "(int,int)test()")

	tester.assertIntAt(0, big.NewInt(1))
	tester.assertIntAt(1, big.NewInt(2))
}

func TestThreeReturnValues(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, int, int) test() {
			return 1, 2, 3
		}
	`, "(int,int,int)test()")

	assert.Equal(t, len(tester.evalStack), 3)
	tester.assertIntAt(0, big.NewInt(1))
	tester.assertIntAt(1, big.NewInt(2))
	tester.assertIntAt(2, big.NewInt(3))
}

func TestReturnMultipleValuesSameTypes(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, int) test() {
			int x = 1
			int y = 2
			return x, y
		}
	`, "(int,int)test()")

	tester.assertIntAt(0, big.NewInt(1))
	tester.assertIntAt(1, big.NewInt(2))
}

func TestReturnMultipleValuesDifferentTypes(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, bool) test() {
			int x = 1
			bool y = true
			return x,y
		}
	`, "(int,bool)test()")

	tester.assertIntAt(0, big.NewInt(1))
	tester.assertBoolAt(1, true)
}

// If statements
// ---------------

func TestIfStatement(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			if (true) {
				return 1
			}
			return 0
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1))
}

func TestSkipIfStatement(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			if (false) {
				return 1
			}
			return 0
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(0))
}

func TestIfElseStatement(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			if (true) {
				return 1
			} else {
				return 0
			}
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1))
}

func TestIfElseStatementAlternative(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			if (false) {
				return 1
			} else {
				return 0
			}
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(0))
}

func TestNestedIfStatement(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			if (true) {
				if (true) {
					return 1
				}
			} 
			return 0
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1))
}

func TestNestedIfStatementAlternative(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
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
	`, intTestSig)

	tester.assertInt(big.NewInt(2))
}

func TestSetter(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		int x = 4
	
		function void set() {
			x = 5
		}
	`, "()set()")

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(5))
}

// Function Calls
// --------------

func TestFuncCall(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			return add(10, 20)
		}

		function int add(int x, int y) {
			return x + y
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(30))
}

func TestFuncCallWithMultiReturn(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, int) test() {
			return calc(10, 20)
		}

		function (int, int) calc(int x, int y) {
			return x + y, x * y
		}
	`, "(int,int)test()")

	tester.assertIntAt(0, big.NewInt(30))
	tester.assertIntAt(1, big.NewInt(200))
}

func TestFuncCallVoid(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		int x = 4
		
		function void test() {
			set()
		}
	
		function void set() {
			x = 5
		}
	`, voidTestSig)

	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	assert.Equal(t, tester.context.ContractVariables[0] == nil, false)
	tester.assertVariableInt(0, big.NewInt(5))
}

// Struct
// ------

func TestEmptyStruct(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
		}

		function Person test() {
			Person p
			return p
		}
	`, "(Person)test()")

	expected := []byte{
		0x02,       // array type
		0x00, 0x00, // array length = 2
	}

	tester.assertBytes(expected...)
}

func TestDefaultStructValues(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
		}

		function Person test() {
			Person p
			return p
		}
	`, "(Person)test()")

	expected := []byte{
		0x02,       // array type
		0x00, 0x02, // array length = 2
		0x00, 0x00, // index 0: empty string size = 0
		0x00, 0x01, // index 1: int size = 1
		0x00, // int 0
	}

	tester.assertBytes(expected...)
}

func TestStructLoadField(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
		}

		function (String, int) test() {
			Person p
			return p.name, p.balance
		}
	`, "(String,int)test()")

	assert.Equal(t, len(tester.evalStack), 2)
	tester.assertBytesAt(0)
	tester.assertIntAt(1, big.NewInt(0))
}

// Arithmetic Expressions
// ----------------------

func TestAddition(t *testing.T) {
	assertIntExpr(t, "1 + 2", 3)
}

func TestAdditionVar(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 1
			int y = 2
			return x + y
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(3))
}

func TestSubtraction(t *testing.T) {
	assertIntExpr(t, "2 - 1", 1)
	assertIntExpr(t, "1 - 2", -1)
}

func TestSubtractionVar(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 1
			return 2 - x
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1))
}

func TestMultiplication(t *testing.T) {
	assertIntExpr(t, "2 * 3", 6)
}

func TestMultiplicationVar(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 2
			return x * 3
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(6))
}

func TestSubMulOrder(t *testing.T) {
	assertIntExpr(t, "8 - 4 * 2", 0)
	assertIntExpr(t, "8 * 4 - 2", 30)
}

func TestDivision(t *testing.T) {
	assertIntExpr(t, "10 / 5", 2)
}

func TestDivisionVar(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 5
			return 10 / x
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(2))
}

func TestDivisionRound(t *testing.T) {
	assertIntExpr(t, "5 / 2", 2)
}

func TestModulo(t *testing.T) {
	assertIntExpr(t, "5 % 2", 1)
}

func TestExponent(t *testing.T) {
	assertIntExpr(t, "2 ** 3", 8)
	assertIntExpr(t, "2 ** 3 ** 2", 512) // 2^9
	assertIntExpr(t, "2 ** 3 ** 4 ** 0", 8)
}

func TestExponentVar(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 3
			return 2 ** x
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(8))
}

func TestExpWithMul(t *testing.T) {
	assertIntExpr(t, "2 * 3 ** 4", 162)
	assertIntExpr(t, "2 ** 3 * 4", 32)
}

func TestMixedOperators(t *testing.T) {
	assertIntExpr(t, "5 * 4 + 2 ** 3 - 1", 27)
}

// Logical Expressions
// -------------------

func TestLogicAnd(t *testing.T) {
	assertBoolExpr(t, "true && true", true)
	assertBoolExpr(t, "true && false", false)
}

func TestLogicAndShortCircuit(t *testing.T) {
	assertBoolExpr(t, "false && true", false)
	assertBoolExpr(t, "false && false", false)
}

func TestLogicOr(t *testing.T) {
	assertBoolExpr(t, "false || false", false)
	assertBoolExpr(t, "false || true", true)
}

func TestLogicOrShortCircuit(t *testing.T) {
	assertBoolExpr(t, "true || false", true)
	assertBoolExpr(t, "true || true", true)
}

func TestLogicNot(t *testing.T) {
	assertBoolExpr(t, "!true", false)
	assertBoolExpr(t, "!false", true)
	assertBoolExpr(t, "!!true", true)
}

// Equality Comparison
// --------------------

func TestIntEqual(t *testing.T) {
	assertBoolExpr(t, "4 == 4", true)
	assertBoolExpr(t, "-4 == -4", true)
	assertBoolExpr(t, "1 == 2", false)
}

func TestIntUnequal(t *testing.T) {
	assertBoolExpr(t, "4 != 4", false)
	assertBoolExpr(t, "-4 != -4", false)
	assertBoolExpr(t, "1 != 2", true)
}

func TestBoolEqual(t *testing.T) {
	assertBoolExpr(t, "true == true", true)
	assertBoolExpr(t, "false == false", true)
	assertBoolExpr(t, "true == false", false)
}

func TestBoolUnequal(t *testing.T) {
	assertBoolExpr(t, "true != true", false)
	assertBoolExpr(t, "false != false", false)
	assertBoolExpr(t, "true != false", true)
}

func TestCharEqual(t *testing.T) {
	assertBoolExpr(t, "'a' == 'a'", true)
	assertBoolExpr(t, "'a' == 'b'", false)
}

func TestCharUnequal(t *testing.T) {
	assertBoolExpr(t, "'a' != 'a'", false)
	assertBoolExpr(t, "'a' != 'b'", true)
}

func TestStringEqual(t *testing.T) {
	assertBoolExpr(t, " \"hello\" == \"hello\" ", true)
	assertBoolExpr(t, " \"hello\" == \"world\" ", false)
}

func TestStringUnequal(t *testing.T) {
	assertBoolExpr(t, " \"hello\" != \"hello\" ", false)
	assertBoolExpr(t, " \"hello\" != \"world\" ", true)
}

// Relational Comparison
// --------------------

func TestIntLess(t *testing.T) {
	assertBoolExpr(t, "1 < 3", true)
	assertBoolExpr(t, "1 < 1", false)
	assertBoolExpr(t, "3 < 1", false)
}

func TestIntLessEqual(t *testing.T) {
	assertBoolExpr(t, "1 <= 3", true)
	assertBoolExpr(t, "3 <= 3", true)
	assertBoolExpr(t, "3 <= 1", false)
}

func TestIntGreater(t *testing.T) {
	assertBoolExpr(t, "1 > 3", false)
	assertBoolExpr(t, "1 > 1", false)
	assertBoolExpr(t, "3 > 1", true)
}

func TestIntGreaterEqual(t *testing.T) {
	assertBoolExpr(t, "1 >= 3", false)
	assertBoolExpr(t, "3 >= 1", true)
}

func TestCharLess(t *testing.T) {
	assertBoolExpr(t, "'a' < 'b'", true)
	assertBoolExpr(t, "'a' < 'a'", false)
	assertBoolExpr(t, "'b' < 'a'", false)
}

func TestCharLessEqual(t *testing.T) {
	assertBoolExpr(t, "'a' <= 'b'", true)
	assertBoolExpr(t, "'b' <= 'b'", true)
	assertBoolExpr(t, "'b' <= 'a'", false)
}

func TestCharGreater(t *testing.T) {
	assertBoolExpr(t, "'a' > 'b'", false)
	assertBoolExpr(t, "'a' > 'a'", false)
	assertBoolExpr(t, "'b' > 'a'", true)
}

func TestCharGreaterEqual(t *testing.T) {
	assertBoolExpr(t, "'a' >= 'b'", false)
	assertBoolExpr(t, "'b' >= 'b'", true)
	assertBoolExpr(t, "'b' >= 'a'", true)
}
