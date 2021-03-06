package generator

import (
	"bytes"
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

	tester.assertVariableInt(0, big.NewInt(0))
	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[0], []byte{0})
}

func TestContractFieldExpression(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x = 4 * 12
	`)

	tester.assertVariableInt(0, big.NewInt(48))
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[0], []byte{0, 48})
}

func TestMultipleContractFields(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x = 4 * 12
		int y = 3 * 12
	`)

	tester.assertVariableInt(1, big.NewInt(36))
	assert.Equal(t, tester.context.ContractVariables[1] == nil, true)
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[1], []byte{0, 36})
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

	tester.assertVariableInt(0, big.NewInt(5))
	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[0], []byte{0, 5})
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

// Shorthand Assignments
// ---------------------

func TestPostfixIncrement(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			int x = 2
			x++
			return x
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(3))
}

func TestPostfixDecrement(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x
		
		constructor() {
			x--
		}
	`)

	tester.assertVariableInt(0, big.NewInt(-1))
}

func TestShorthandAssignments(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int add = 1
		int sub = 2
		int mul = -3
		int div = 4
		int mod = 5
		int exp = 2
		int shiftL = 1
		int shiftR = 8
		int bitAnd = 5
		int bitOr = 5
		int bitXor = 5

		constructor() {
			add += 1
			sub -= 1
			mul *= -1
			div /= 2
			mod %= 2
			exp **= 3
			shiftL <<= 3
			shiftR >>= 3
			bitAnd &= 3
			bitOr |= 3
			bitXor ^= 3	
		}
	`)

	expected := []int64{
		2, // 1 + 1
		1, // 2 - 1
		3, // -3 * -1
		2, // 4 / 2
		1, // 5 % 2
		8, // 2 ** 3
		8, // 1 << 3
		1, // 8 >> 3
		1, // 5 & 3
		7, // 5 | 3
		6, // 5 ^ 3
	}
	assert.Equal(t, len(tester.context.ContractVariables), len(expected))
	for i := 0; i < len(expected); i++ {
		tester.assertVariableInt(i, big.NewInt(expected[i]))
	}
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

	tester.assertVariableInt(0, big.NewInt(5))
	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[0], []byte{0, 5})
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

	tester.assertVariableInt(0, big.NewInt(5))
	assert.Equal(t, tester.context.ContractVariables[0] == nil, true)
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[0], []byte{0, 5})
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

func TestStructFieldAccess(t *testing.T) {
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

func TestStructFieldAssignment(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
		}

		function (String, int) test() {
			Person p
			p.name = "a"
			p.balance = 100
			return p.name, p.balance
		}
	`, "(String,int)test()")

	assert.Equal(t, len(tester.evalStack), 2)
	tester.assertBytesAt(0, 97)
	tester.assertIntAt(1, big.NewInt(100))
}

func TestStructNewCreation(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
		}

		function (String, int) test() {
			Person p = new Person()
			return p.name, p.balance
		}
	`, "(String,int)test()")

	assert.Equal(t, len(tester.evalStack), 2)
	tester.assertBytesAt(0)
	tester.assertIntAt(1, big.NewInt(0))
}

func TestStructNewCreationWithInitialValue(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
		}

		function (String, int) test() {
			Person p = new Person("a")
			return p.name, p.balance
		}
	`, "(String,int)test()")

	assert.Equal(t, len(tester.evalStack), 2)
	tester.assertBytesAt(0, 97)
	tester.assertIntAt(1, big.NewInt(0))
}

func TestStructNewCreationWithInitialValues(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
		}

		function (String, int) test() {
			Person p = new Person("a", 200)
			return p.name, p.balance
		}
	`, "(String,int)test()")

	assert.Equal(t, len(tester.evalStack), 2)
	tester.assertBytesAt(0, 97)
	tester.assertIntAt(1, big.NewInt(200))
}

func TestStructCreationWithNamedFieldValues(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
			bool valid
		}

		function (String, int, bool) test() {
			Person p = new Person(balance=300)
			return p.name, p.balance, p.valid
		}
	`, "(String,int,bool)test()")

	assert.Equal(t, len(tester.evalStack), 3)
	tester.assertBytesAt(0)
	tester.assertIntAt(1, big.NewInt(300))
	tester.assertBoolAt(2, false)
}

func TestStructCreationWithNamedFieldValues2(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			String name
			int balance
			bool valid
		}

		function (String, int, bool) test() {
			Person p = new Person(valid=true, name="a", balance=400)
			return p.name, p.balance, p.valid
		}
	`, "(String,int,bool)test()")

	assert.Equal(t, len(tester.evalStack), 3)
	tester.assertBytesAt(0, 97)
	tester.assertIntAt(1, big.NewInt(400))
	tester.assertBoolAt(2, true)
}

func TestStructMultiFieldAssignment(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
			bool isValid
		}

		function (int, bool) test() {
			Person p
			p.balance, p.isValid = test2()
			
			return p.balance, p.isValid
		}

		function (int, bool) test2() {
			return 100, true
		}
	`, "(int,bool)test()")

	assert.Equal(t, len(tester.evalStack), 2)
	tester.assertIntAt(0, big.NewInt(100))
	tester.assertBoolAt(1, true)
}

func TestStructNestedFieldAccess(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
			Person friend
		}

		function int test() {
			Person p
			p.friend = new Person()
			return p.friend.balance
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(0))
}

func TestStructNestedFieldAssignment(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
			Person friend
		}

		function int test() {
			Person p = new Person(friend = new Person())
			p.friend.balance = 100
			return p.friend.balance
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(100))
}

func TestThisStructNestedFieldAssignment(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
			Person friend
		}

		Person p

		function int test() {
			this.p = new Person()
			this.p.friend = new Person()
			this.p.friend.balance = 100
			return this.p.friend.balance
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(100))
}

func TestStructNestedMultiFieldAssignment(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
			bool isValid
			Person friend
		}

		Person p

		function (int, bool) test() {
			this.p = new Person()
			this.p.friend = new Person()
			this.p.friend.balance, p.friend.isValid = test2()
			
			return this.p.friend.balance, p.friend.isValid
		}

		function (int, bool) test2() {
			return 100, true
		}
	`, "(int,bool)test()")

	assert.Equal(t, len(tester.evalStack), 2)
	tester.assertIntAt(0, big.NewInt(100))
	tester.assertBoolAt(1, true)
}

// Arrays
// ------

func TestUninitializedArray(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[] a
	`)

	variable, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(variable, []byte{})
}

func TestElementAssignment1(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x
		int[] a = new int[1]

		constructor() {
			x = a[0]
		}
	`)

	variable, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(variable, []byte{0})

}

func TestElementAssignment2(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[] a = new int[1]

		constructor() {
			a[0] = 1
		}
	`)

	expected := []byte{
		0x02,       // array type
		0x00, 0x01, // size = 1
		0x00, 0x02, // length of value
		0x00, 0x01, // int value
	}

	variable, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(variable, expected)
}

func TestElementAccessMultiAssignment(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[] a = new int[2]

		constructor() {
			a[0], a[1] = initArray()
		}

		function (int, int) initArray() {
			return 1, 2
		}
	`)

	expected := []byte{
		0x02,       // array type
		0x00, 0x02, // size = 2
		0x00, 0x02, // length of first element
		0x00, 0x01, // first element
		0x00, 0x02, // length of second element
		0x00, 0x02, // second element
	}

	variable, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(variable, expected)
}

func TestArrayLengthCreation(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[] a = new int[2]
	`)

	expected := []byte{
		0x02,       // array type
		0x00, 0x02, // size = 2
		0x00, 0x01, // length of first element
		0x00,       // first element
		0x00, 0x01, // length of second element
		0x00, // second element
	}

	variable, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(variable, expected)
}

func TestNestedArrayLengthCreation(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[][] a = new int[2][2]
	`)

	tester.assertErrorAt(0, "Generator currently does not support array nesting")
}

func TestArrayValueCreation(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[] a = new int[]{1, 2}
	`)

	expected := []byte{
		0x02,       // array type
		0x00, 0x02, // size = 2
		0x00, 0x02, // length of first element
		0x00, 0x01, // first element
		0x00, 0x02, // length of second element
		0x00, 0x02, // second element
	}

	variable, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(variable, expected)
}

func TestNestedArrayValueCreation(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[][] a = new int[][]{{1, 2}, {3}}
	`)

	tester.assertErrorAt(0, "Generator currently does not support array nesting")
}

func TestArrayInStructUpdateError(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		struct Person{
			int[] nums
		}
		constructor(){
			Person p = new Person(nums = new int[2])
			p.nums[0] = 1
		}
	`)

	tester.assertErrorAt(0, "Multiple dereferences on value types are not supported")
}

func TestArrayInStructUpdate(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person{
			int[] nums
		}
		
		function int test(){
			Person p = new Person(nums = new int[2])
			int[] copy = p.nums
			copy[0] = 1
			p.nums = copy
			
			return p.nums[0] 
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1))
}

func TestStructArrayUpdateError(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person{
			int balance
		}
		
		function int test(){
			Person[] p = new Person[2]
			p[0] = new Person(100)
			p[0].balance = 101

			return p[0].balance
		}
	`, intTestSig)

	tester.assertErrorAt(0, "Updating struct value type in array/map is not supported")
}

func TestStructArrayUpdate(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person{
			int balance
		}
		
		function int test(){
			Person[] p = new Person[2]
			p[0] = new Person(100)

			Person copy = p[0]
			copy.balance = 101
			p[0] = copy

			return p[0].balance
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(101))
}

// Map
// ---

func TestEmptyMap(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		Map<int, int> m
	`)

	expected := []byte{
		0x01,       // Map type
		0x00, 0x00, // Map size
	}

	field, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(field, expected)
}

func TestMapPush(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		Map<int, bool> m

		constructor() {
			m[3] = true 
		}
	`)

	expected := []byte{
		0x01,
		0x00, 0x01,
		0x00, 0x02, // length of key
		0x00, 0x03, // key 3
		0x00, 0x01, // length of value
		0x01, // value true
	}

	field, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(field, expected)
}

func TestMapPushMultiple(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		Map<String, int> m

		constructor() {
			m["a"] = 1
			m["ab"] = 2
		}
	`)

	expected := []byte{
		0x01,
		0x00, 0x02,
		0x00, 0x01, // length of key
		0x61,       // key "a"
		0x00, 0x02, // length of value
		0x00, 0x01, // value
		0x00, 0x02, // length of key
		0x61, 0x62, // key "ab"
		0x00, 0x02, // length of value
		0x00, 0x02, // value
	}

	field, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(field, expected)
}

func TestMapPushOverride(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		Map<String, int> m

		constructor() {
			m["a"] = 1
			m["a"] = 2
		}
	`)

	expected := []byte{
		0x01,
		0x00, 0x01,
		0x00, 0x01, // length of key
		0x61,       // key "a"
		0x00, 0x02, // length of value
		0x00, 0x02, // value
	}

	field, err := tester.context.GetContractVariable(0)
	assert.NilError(t, err)
	tester.compareBytes(field, expected)
}

func TestMapGetVal(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			Map<String, int> m
			m["a"] = 1234
			return m["a"]
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1234))
}

func TestMapDeleteKey(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function Map<String, int> test() {
			Map<String, int> m
			m["a"] = 1234
			delete m["a"]
			return m
		}
	`, "(Map<String,int>)test()")

	tester.assertBytes(1, 0, 0) // Empty map
}

func TestMapContainsKey(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function bool test() {
			Map<String, int> m
			m["a"] = 1234
			return m.contains("a")
		}
	`, boolTestSig)

	tester.assertBool(true)
}

func TestMapContainsKeyFalse(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function bool test() {
			Map<String, int> m
			return m.contains("a")
		}
	`, boolTestSig)

	tester.assertBool(false)
}

func TestMapStructValError(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
		}
	
		constructor() {
			Map<char, Person> m
			m['a'] = new Person(1000)
			m['a'].balance = 1001
		}
	`, intTestSig)

	tester.assertErrorAt(0, "Updating struct value type in array/map is not supported")
}

func TestMapStructVal(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
		}
	
		function int test() {
			Map<char, Person> m
			m['a'] = new Person(1000)

			Person p = m['a']
			p.balance = 1001
			return m['a'].balance
		}
	`, intTestSig)

	// Struct is a value type, so changing the copy will not affect the map's value
	tester.assertInt(big.NewInt(1000))
}

func TestMapStructValOverride(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		struct Person {
			int balance
		}
	
		function int test() {
			Map<char, Person> m
			m['a'] = new Person(1000)

			Person p = m['a']
			p.balance = 1001

			m['a'] = p
			return m['a'].balance
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(1001))
}

func TestMapArrayValueUpdateError(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		constructor(){
			Map<String, int[]> m
			m["a"] = new int[2]

			m["a"][0] = 2
		}
	`)

	tester.assertErrorAt(0, "Multiple dereferences on value types are not supported")
}

func TestMapArrayValueUpdate(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function int test(){
			Map<String, int[]> m
			m["a"] = new int[2]
			
			int[] copy = m["a"] 
			copy[0] = 2
			m["a"] = copy

			return m["a"][0]
		}
	`, intTestSig)

	tester.assertInt(big.NewInt(2))
}

func TestMapElementMultiAssignment(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function (int, int) test(){
			Map<String, int> m
			m["a"], m["b"] = test2()
			
			return m["a"], m["b"]
		}

		function (int, int) test2() {
			return 10, 20
		}
	`, "(int,int)test()")

	tester.assertIntAt(0, big.NewInt(10))
	tester.assertIntAt(1, big.NewInt(20))
}

// Ternary Expression
// ------------------

func TestTernaryExpression(t *testing.T) {
	assertTernaryExpr(t, "true ? 1 + 2 : 3 + 4", "int", "3")
}

func TestTernaryExpressionFalse(t *testing.T) {
	assertTernaryExpr(t, "false ? 1 : 2", "int", "2")
}

// Binary Arithmetic Expressions
// -----------------------------

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
	assertIntExpr(t, "-2 * 3", -6)
	assertIntExpr(t, "-2 * -3", 6)
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

func TestParentheses(t *testing.T) {
	assertIntExpr(t, "(2 + 3) * 4", 20)
}

func TestShiftLeft(t *testing.T) {
	// 1 --> 1000
	assertIntExpr(t, "1 << 3", 8)
	assertIntExpr(t, "1 << 0", 1)
}

func TestShiftRight(t *testing.T) {
	// 1000 --> 1
	assertIntExpr(t, "8 >> 3", 1)
	assertIntExpr(t, "-8 >> 3", -1)
	assertIntExpr(t, "1 >> 3", 0)
}

// Make sure the VM trace is set to false.
// Otherwise printing 536870913 bytes on console slows down the test extremely.
func TestShiftLeft_32bit_Max(t *testing.T) {
	expected := big.NewInt(1)
	expected.Lsh(expected, 0xffffffff)

	tester := newGeneratorTestUtilWithFunc(t, `
		function int test() {
			return 1 << 0xffffffff
		}
	`, intTestSig)
	actual := tester.result[1:] // first byte is sign. It can be ignored, since it is 0.

	// DO NOT use assertBytes helper function, because iterating 536870913 bytes
	// and asserting every byte takes extremely long.
	// Using built-in compare is much faster for much larger slice.
	result := bytes.Compare(actual, expected.Bytes())
	assert.Equal(t, result, 0)
}

// String Concatenation
// --------------------

func TestStringConcatenation(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function String test() {
			return "Hello" + "World"
		}
	`, stringTestSig)

	tester.assertErrorAt(0, "String concatenation is not supported")
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

// Bitwise Logical Expressions
// ---------------------------

func TestBitwiseNot(t *testing.T) {
	// Use http://bitwisecmd.com/ to verify the results visually
	assertIntExpr(t, "~5", -6)
	assertIntExpr(t, "~-6", 5)
}

func TestBitwiseAnd(t *testing.T) {
	// 0101 -> 5
	// 0011 -> 3
	// ---- &
	// 0001 -> 1
	assertIntExpr(t, "5 & 3", 1)
}

func TestBitwiseOr(t *testing.T) {
	// 0101 -> 5
	// 0011 -> 3
	// ---- |
	// 0111 -> 7
	assertIntExpr(t, "5 | 3", 7)
}

func TestBitwiseXOr(t *testing.T) {
	// 0101 -> 5
	// 0011 -> 3
	// ---- ^
	// 0110 -> 6
	assertIntExpr(t, "5 ^ 3", 6)
}

// Type Cast
// ---------

func TestIntToStringTypeCast(t *testing.T) {
	tester := newGeneratorTestUtilWithFunc(t, `
		function String test() {
			return (String) 1
		}
	`, stringTestSig)

	tester.assertErrorAt(0, "VM currently does not support types")
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

// This Keyword
// ------------

func TestThisMemberAccess(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int x

		constructor(){
			this.x = 3
		}
	`)

	tester.assertVariableInt(0, big.NewInt(3))
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[0], []byte{0, 3})
}

// Length Member
// -------------

func TestArrayLengthMemberAccess(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		int[] x = new int[2]
		int y

		constructor(){
			y = x.length
		}
	`)

	tester.assertVariableInt(1, big.NewInt(2))
	tester.context.PersistChanges()
	tester.compareBytes(tester.context.ContractVariables[1], []byte{0, 2})
}
