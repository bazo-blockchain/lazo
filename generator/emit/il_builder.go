package emit

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator/data"
	"github.com/bazo-blockchain/lazo/generator/il"
	"strings"
)

/**
 *	IL Builder constructs Metadata
 */
type ILBuilder struct {
	symbolTable       *symbol.SymbolTable
	Metadata          *data.Metadata
	functionData      map[*symbol.FunctionSymbol]*data.FunctionData
	functionPositions map[*symbol.FunctionSymbol]uint16
	Errors            []error
}

func NewILBuilder(symbolTable *symbol.SymbolTable) *ILBuilder {
	builder := &ILBuilder{
		symbolTable:       symbolTable,
		Metadata:          &data.Metadata{},
		functionData:      map[*symbol.FunctionSymbol]*data.FunctionData{},
		functionPositions: map[*symbol.FunctionSymbol]uint16{},
	}
	builder.generateMetadata()
	return builder
}

func (b *ILBuilder) generateMetadata() {
	contract := b.symbolTable.GlobalScope.Contract
	b.registerContract(contract)
	b.fixContract(contract)
}

func (b *ILBuilder) Complete() *data.Metadata {
	b.fixOperands(b.Metadata.Contract.Instructions)
	for _, function := range b.Metadata.Contract.Functions {
		b.fixOperands(function.Instructions)
	}
	return b.Metadata
}

func (b *ILBuilder) GetFunctionData(function *symbol.FunctionSymbol) *data.FunctionData {
	return b.functionData[function]
}

func (b *ILBuilder) SetFunctionPos(symbol *symbol.FunctionSymbol, pos uint16) {
	b.functionPositions[symbol] = pos
}

func (b *ILBuilder) fixOperands(code []*il.Instruction) {
	for _, instruction := range code {
		if typeSymbol, ok := instruction.Operand.(*symbol.TypeSymbol); ok {
			instruction.Operand = b.getTypeRef(typeSymbol)
		} else if functionSymbol, ok := instruction.Operand.(*symbol.FunctionSymbol); ok {
			operand := make([]byte, 3)
			binary.BigEndian.PutUint16(operand, uint16(b.functionPositions[functionSymbol]))
			operand[2] = byte(len(functionSymbol.Parameters))
			instruction.Operand = operand
		}
	}
}

func (b *ILBuilder) registerContract(contract *symbol.ContractSymbol) {
	b.Metadata.Contract = &data.ContractData{
		Identifier: contract.GetIdentifier(),
	}
	for _, function := range contract.Functions {
		b.registerFunction(function)
	}
}

func (b *ILBuilder) registerFunction(function *symbol.FunctionSymbol) {
	functionData := &data.FunctionData{
		Identifier: function.GetIdentifier(),
		Hash:       createFuncHash(createFuncSignature(function)),
	}
	b.Metadata.Contract.Functions = append(b.Metadata.Contract.Functions, functionData)
	b.functionData[function] = functionData
}

func (b *ILBuilder) fixContract(contract *symbol.ContractSymbol) {
	contractData := b.Metadata.Contract

	for _, field := range contract.Fields {
		contractData.Fields = append(contractData.Fields, b.getTypeRef(field.Type))
	}

	for _, function := range contract.Functions {
		b.fixFunction(function)
	}
}

func (b *ILBuilder) fixFunction(function *symbol.FunctionSymbol) {
	functionData := b.GetFunctionData(function)

	for _, rtype := range function.ReturnTypes {
		functionData.ReturnTypes = append(functionData.ReturnTypes, b.getTypeRef(rtype))
	}

	for _, param := range function.Parameters {
		functionData.ParamTypes = append(functionData.ParamTypes, b.getTypeRef(param.Type))
	}

	for _, local := range function.LocalVariables {
		functionData.LocalTypes = append(functionData.LocalTypes, b.getTypeRef(local.Type))
	}
}

func (b *ILBuilder) getTypeRef(sym *symbol.TypeSymbol) data.TypeData {
	scope := b.symbolTable.GlobalScope
	if sym.GetIdentifier() == scope.BoolType.GetIdentifier() {
		return data.BoolType
	} else if sym.GetIdentifier() == scope.CharType.GetIdentifier() {
		return data.CharType
	} else if sym.GetIdentifier() == scope.StringType.GetIdentifier() {
		return data.StringType
	} else if sym.GetIdentifier() == scope.IntType.GetIdentifier() {
		return data.IntType
	} else {
		panic(fmt.Sprintf("Error: Unsupported Type %s", sym.GetIdentifier()))
	}
}

// Helper Functions
// ----------------

func createFuncHash(funcSig string) [4]byte {
	h := sha256.Sum256([]byte(funcSig))
	var arr [4]byte
	for i := 0; i < 4; i++ {
		arr[i] = h[i]
	}
	return arr
}

func createFuncSignature(function *symbol.FunctionSymbol) string {
	var sb strings.Builder

	sb.WriteRune('(')
	for i, r := range function.ReturnTypes {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(r.Identifier)
	}
	sb.WriteRune(')')

	sb.WriteString(function.Identifier)

	sb.WriteRune('(')
	for i, p := range function.Parameters {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(p.Type.Identifier)
	}
	sb.WriteRune(')')

	return sb.String()
}
