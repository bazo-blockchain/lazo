package util

import (
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func TestGetSignByte(t *testing.T) {
	testData := map[int64]byte{
		-1: 1,
		0:  0,
		1:  0,
	}

	for k, v := range testData {
		assert.Equal(t, GetSignByte(big.NewInt(k)), v)
	}
}

func TestGetBytesFromUInt16(t *testing.T) {
	bytes := GetBytesFromUInt16(uint16(0x1234))
	assert.Equal(t, bytes[0], byte(0x12))
	assert.Equal(t, bytes[1], byte(0x34))
}
