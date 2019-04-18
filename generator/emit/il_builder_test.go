package emit

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/util"
	"gotest.tools/assert"
	"testing"
)

func TestVoidFuncSignature(t *testing.T) {
	f := symbol.NewFunctionSymbol(nil, "test")

	assert.Equal(t, createFuncSignature(f), "()test()")
}

func TestIntFuncSignature(t *testing.T) {
	f := symbol.NewFunctionSymbol(nil, "test")
	f.ReturnTypes = append(f.ReturnTypes, symbol.NewTypeSymbol(f, "int"))

	assert.Equal(t, createFuncSignature(f), "(int)test()")
}

func TestMultipleReturnFuncSignature(t *testing.T) {
	f := symbol.NewFunctionSymbol(nil, "test")
	f.ReturnTypes = append(f.ReturnTypes, symbol.NewTypeSymbol(f, "int"))
	f.ReturnTypes = append(f.ReturnTypes, symbol.NewTypeSymbol(f, "bool"))

	assert.Equal(t, createFuncSignature(f), "(int,bool)test()")
}

func TestFuncSigSingleParam(t *testing.T) {
	f := symbol.NewFunctionSymbol(nil, "test")
	intType := symbol.NewTypeSymbol(f, "int")
	p1 := symbol.NewParameterSymbol(nil, "x")
	p1.Type = intType
	f.Parameters = append(f.Parameters, p1)

	assert.Equal(t, createFuncSignature(f), "()test(int)")
}

func TestFuncSigMultipleParams(t *testing.T) {
	f := symbol.NewFunctionSymbol(nil, "test")
	intType := symbol.NewTypeSymbol(f, "int")
	p1 := symbol.NewParameterSymbol(nil, "x")
	p2 := symbol.NewParameterSymbol(nil, "y")
	p1.Type = intType
	p2.Type = intType
	f.Parameters = append(f.Parameters, p1, p2)

	assert.Equal(t, createFuncSignature(f), "()test(int,int)")
}

func TestFuncSig(t *testing.T) {
	f := symbol.NewFunctionSymbol(nil, "test")
	intType := symbol.NewTypeSymbol(f, "int")
	boolType := symbol.NewTypeSymbol(f, "bool")
	f.ReturnTypes = append(f.ReturnTypes, intType)
	f.ReturnTypes = append(f.ReturnTypes, boolType)

	p1 := symbol.NewParameterSymbol(nil, "x")
	p2 := symbol.NewParameterSymbol(nil, "y")
	p1.Type = intType
	p2.Type = boolType
	f.Parameters = append(f.Parameters, p1, p2)

	assert.Equal(t, createFuncSignature(f), "(int,bool)test(int,bool)")
}

func TestFuncHashWithSymbol(t *testing.T) {
	f := symbol.NewFunctionSymbol(nil, "test")
	intType := symbol.NewTypeSymbol(f, "int")
	boolType := symbol.NewTypeSymbol(f, "bool")
	f.ReturnTypes = append(f.ReturnTypes, intType)
	f.ReturnTypes = append(f.ReturnTypes, boolType)

	p1 := symbol.NewParameterSymbol(nil, "x")
	p2 := symbol.NewParameterSymbol(nil, "y")
	p1.Type = intType
	p2.Type = boolType
	f.Parameters = append(f.Parameters, p1, p2)

	actual := util.CreateFuncHash(createFuncSignature(f))
	expected := []byte{0x77, 0xA1, 0x94, 0x20}
	for i, h := range actual {
		assert.Equal(t, h, expected[i])
	}
}
