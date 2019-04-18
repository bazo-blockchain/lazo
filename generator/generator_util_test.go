package generator

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/bazo-vm/vm"
	"github.com/bazo-blockchain/lazo/checker"
	"github.com/bazo-blockchain/lazo/generator/data"
	"github.com/bazo-blockchain/lazo/generator/util"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser"
	"gotest.tools/assert"
	"math/big"
	"strings"
	"testing"
)

const (
	VoidTestSig   = "()test()"
	IntTestSig    = "(int)test()"
	BoolTestSig   = "(bool)test()"
	CharTestSig   = "(char)test()"
	StringTestSig = "(string)test()"
)

type generatorTestUtil struct {
	t        *testing.T
	metadata *data.Metadata
	context  *vm.MockContext
	result   []byte
	errors   []error
}

func newGeneratorTestUtil(t *testing.T, contractCode string) *generatorTestUtil {
	txData := []byte{
		1, 0, // call constructor
	}

	return newGeneratorTestUtilWithRawInput(
		t,
		fmt.Sprintf("contract Test {\n %s \n }", contractCode),
		txData,
	)
}

func newGeneratorTestUtilWithFunc(t *testing.T, contractCode string,
	funcSignature string, funcData ...byte) *generatorTestUtil {
	funcHash := util.CreateFuncHash(funcSignature)
	txData := append(funcData, 4)
	txData = append(txData, funcHash[:]...)

	return newGeneratorTestUtilWithRawInput(
		t,
		fmt.Sprintf("contract Test {\n %s \n }", contractCode),
		txData,
	)
}

func newGeneratorTestUtilWithRawInput(t *testing.T, code string, txData []byte) *generatorTestUtil {
	p := parser.New(lexer.New(bufio.NewReader(strings.NewReader(code))))
	program, err := p.ParseProgram()
	assert.Equal(t, len(err), 0, "Program has syntax errors", err)

	symbolTable, err := checker.New(program).Run()
	assert.Equal(t, len(err), 0, "Program has semantic errors")

	tester := &generatorTestUtil{
		t: t,
	}

	tester.metadata, tester.errors = New(symbolTable).Run()
	assert.Equal(t, len(err), 0, "Error while generating byte code")

	byteCode, variables := tester.metadata.CreateContract()
	context := vm.NewMockContext(byteCode)
	context.ContractVariables = variables
	context.Data = txData
	context.Fee += (uint64(len(variables))) * 1000
	context.Fee += 10000 // To be able to calculate 2^16
	tester.context = context

	bazoVM := vm.NewVM(context)
	isSuccess := bazoVM.Exec(false)
	result, vmError := bazoVM.PeekResult()
	assert.Assert(t, isSuccess, string(result))

	tester.result = result
	tester.errors = append(tester.errors, vmError)
	return tester
}

func (gt *generatorTestUtil) assertInt(value *big.Int) {
	bytes := append([]byte{util.GetSignByte(value)}, value.Bytes()...)
	gt.assertBytes(bytes...)
}

// Can be deleted as soon as VM is fixed
func (gt *generatorTestUtil) assertBool(value bool) {
	if value {
		gt.assertBytes(1)
	} else {
		gt.assertBytes(0)
	}
}

func (gt *generatorTestUtil) assertString(value string) {
	gt.assertBytes([]byte(value)...)
}

func (gt *generatorTestUtil) assertChar(value rune) {
	gt.assertBytes(byte(value))
}

func (gt *generatorTestUtil) assertBytes(bytes ...byte) {
	gt.compareBytes(gt.result, bytes)
}

func (gt *generatorTestUtil) assertVariableInt(index int, value *big.Int) {
	bytes := gt.context.ContractVariables[index]
	expected := append([]byte{util.GetSignByte(value)}, value.Bytes()...)
	gt.compareBytes(bytes, expected)
}

func (gt *generatorTestUtil) compareBytes(actual []byte, expected []byte) {
	assert.Equal(gt.t, len(actual), len(expected))

	for i, b := range actual {
		assert.Equal(gt.t, b, expected[i])
	}
}

func assertBoolExpr(t *testing.T, expr string, expected bool) {
	code := fmt.Sprintf("function bool test() {\n return %s \n }", expr)
	tester := newGeneratorTestUtilWithFunc(t, code, BoolTestSig)
	tester.assertBool(expected)
}

func assertIntExpr(t *testing.T, expr string, expected int64) {
	code := fmt.Sprintf("function int test() {\n return %s \n }", expr)
	tester := newGeneratorTestUtilWithFunc(t, code, IntTestSig)
	tester.assertInt(big.NewInt(expected))
}
