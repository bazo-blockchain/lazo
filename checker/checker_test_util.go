package checker

import (
	"bufio"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser"
	"gotest.tools/assert"
	"strings"
	"testing"
)

type CheckerTestUtil struct {
	t *testing.T
	symbolTable *symbol.SymbolTable
	errors      []error
}

func NewCheckerTestUtil(t *testing.T, code string, isValidCode bool) *CheckerTestUtil {
	p := parser.New(lexer.New(bufio.NewReader(strings.NewReader(code))))
	program, err := p.ParseProgram()
	assert.Equal(t, len(err), 0, "Program has syntax errors", err)

	tester := &CheckerTestUtil{
		t: t,
	}
	tester.symbolTable, tester.errors = New(program).Run()
	assert.Equal(t, len(tester.errors) == 0, isValidCode)

	return tester
}