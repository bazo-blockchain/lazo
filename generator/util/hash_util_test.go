package util

import (
	"gotest.tools/assert"
	"testing"
)

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
		hash := CreateFuncHash(k)
		for i, h := range hash {
			assert.Equal(t, h, v[i], k)
		}
	}
}
