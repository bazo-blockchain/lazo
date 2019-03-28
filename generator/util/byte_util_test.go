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
