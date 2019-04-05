package util

import "math/big"

// GetSignByte returns the sign byte of a value on the stack
func GetSignByte(value *big.Int) byte {
	var sign byte
	if value.Sign() == -1 {
		sign = 1
	}
	return sign
}
