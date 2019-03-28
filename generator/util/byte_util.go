package util

import "math/big"

func GetSignByte(value *big.Int) byte {
	var sign byte
	if value.Sign() == -1 {
		sign = 1
	}
	return sign
}
