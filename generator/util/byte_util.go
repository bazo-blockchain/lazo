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

// GetBytesFromUInt16 returns the bytes for uint16 value
func GetBytesFromUInt16(value uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(value))
	return bytes
}
