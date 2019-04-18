package util

import (
	"encoding/binary"
	"math/big"
)

// GetSignByte returns the sign byte of a value on the stack
func GetSignByte(value *big.Int) byte {
	var sign byte
	if value.Sign() == -1 {
		sign = 1
	}
	return sign
}

func GetBytesFromUInt16(element uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(element))
	return bytes
}
