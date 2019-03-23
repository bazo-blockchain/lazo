package generator

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/bazo-vm/vm"
	"github.com/bazo-blockchain/lazo/checker"
	"github.com/bazo-blockchain/lazo/generator/data"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser"
	"gotest.tools/assert"
	"strings"
	"testing"
)

type GeneratorTestUtil struct {
	t        *testing.T
	metaData *data.Metadata
	code     []byte
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

	tester.metaData, tester.errors = New(symbolTable).Run()
	assert.Equal(t, len(err), 0, "Error while generating byte code")

	tester.code = tester.metaData.GetByteCode()
	bazoVM := vm.NewVM(vm.NewMockContext(tester.code))
	isSuccess := bazoVM.Exec(true)
	assert.Assert(t, isSuccess, "Code execution failed")

	result, vmError := bazoVM.PeekResult()
	tester.result = result
	tester.errors = append(tester.errors, vmError)
	return tester
}

func (gt *GeneratorTestUtil) assertBytes(bytes ...byte) {
	assert.Equal(gt.t, len(gt.result), len(bytes))

	for i, b := range bytes {
		assert.Equal(gt.t, gt.result[i], b)
	}
}
