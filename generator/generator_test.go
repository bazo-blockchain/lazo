package generator

import (
	"testing"
)

func TestCodeExecution(t *testing.T) {
	tester := newGeneratorTestUtil(t, `
		function int test() {
        	return 1 + 1
		}
	`)

	tester.assertBytes(0, 2)
}
