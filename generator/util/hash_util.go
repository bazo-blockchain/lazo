package util

import "crypto/sha256"

// CreateFuncHash hashes the function signature and returns the first four bytes.
// Function signature includes return types, function name and parameter types, e.g. (int,bool)test(int)
func CreateFuncHash(funcSig string) [4]byte {
	h := sha256.Sum256([]byte(funcSig))
	var arr [4]byte
	for i := 0; i < 4; i++ {
		arr[i] = h[i]
	}
	return arr
}
