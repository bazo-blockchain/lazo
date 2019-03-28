package emit

import (
	"github.com/bazo-blockchain/lazo/checker/symbol"
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

func TestFuncHashes(t *testing.T) {
	// sha256 values are generated using this tool: https://passwordsgenerator.net/sha256-hash-generator/
	testData := map[string][]byte{
		"()test()":                 {0xD1, 0xFC, 0x69, 0xEB},
		"(int)test()":              {0x13, 0x46, 0x65, 0x4E},
		"(int,bool)test()":         {0xC0, 0x26, 0x4C, 0xF0},
		"()test(int)":              {0x21, 0xB0, 0x59, 0xA3},
		"()test(int,int)":          {0x10, 0xD3, 0x28, 0x4F},
		"(int,bool)test(int,bool)": {0x77, 0xA1, 0x94, 0x20},
	}

	for k, v := range testData {
		hash := createFuncHash(k)
		for i, h := range hash {
			assert.Equal(t, h, v[i], k)
		}
	}
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

	actual := createFuncHash(createFuncSignature(f))
	expected := []byte{0x77, 0xA1, 0x94, 0x20}
	for i, h := range actual {
		assert.Equal(t, h, expected[i])
	}
}
