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

type GeneratorTestUtil struct {
	t        *testing.T
	metadata *data.Metadata
	context  *vm.MockContext
	result   []byte
	errors   []error
}

func newGeneratorTestUtil(t *testing.T, contractCode string) *GeneratorTestUtil {
	return newGeneratorTestUtilWithRawInput(
		t,
		fmt.Sprintf("contract Test {\n %s \n }", contractCode),
	)
}

func newGeneratorTestUtilWithRawInput(t *testing.T, code string) *GeneratorTestUtil {
	p := parser.New(lexer.New(bufio.NewReader(strings.NewReader(code))))
	program, err := p.ParseProgram()
	assert.Equal(t, len(err), 0, "Program has syntax errors", err)

	symbolTable, err := checker.New(program).Run()
	assert.Equal(t, len(err), 0, "Program has semantic errors")

	tester := &GeneratorTestUtil{
		t: t,
	}

	tester.metadata, tester.errors = New(symbolTable).Run()
	assert.Equal(t, len(err), 0, "Error while generating byte code")

	byteCode, variables := tester.metadata.CreateContract()
	context := vm.NewMockContext(byteCode)
	context.ContractVariables = variables
	context.Fee += (uint64(len(variables))) * 1000
	context.Fee += 10000 // To be able to calculate 2^16
	tester.context = context

	bazoVM := vm.NewVM(context)
	isSuccess := bazoVM.Exec(true)
	result, vmError := bazoVM.PeekResult()
	assert.Assert(t, isSuccess, string(result))

	tester.result = result
	tester.errors = append(tester.errors, vmError)
	return tester
}

func (gt *GeneratorTestUtil) assertInt(value *big.Int) {
	bytes := append([]byte{util.GetSignByte(value)}, value.Bytes()...)
	gt.assertBytes(bytes...)
}

func (gt *GeneratorTestUtil) assertBool(value bool) {
	if value {
		gt.assertBytes(0, 1)
	} else {
		gt.assertBytes(0)
	}
}

func (gt *GeneratorTestUtil) assertString(value string) {
	bytes := append([]byte{0}, []byte(value)...)
	gt.assertBytes(bytes...)
}

func (gt *GeneratorTestUtil) assertChar(value rune) {
	bytes := []byte(string(value))
	bytes = append([]byte{0}, bytes...)
	gt.assertBytes(bytes...)
}

func (gt *GeneratorTestUtil) assertBytes(bytes ...byte) {
	gt.compareBytes(gt.result, bytes)
}

func (gt *GeneratorTestUtil) assertVariableInt(index int, value *big.Int) {
	bytes := gt.context.ContractVariables[index]
	expected := append([]byte{util.GetSignByte(value)}, value.Bytes()...)
	gt.compareBytes(bytes, expected)
}

func (gt *GeneratorTestUtil) compareBytes(actual []byte, expected []byte) {
	assert.Equal(gt.t, len(actual), len(expected))

	for i, b := range actual {
		assert.Equal(gt.t, b, expected[i])
	}
}
