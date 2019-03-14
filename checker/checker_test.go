package checker

import (
	"testing"
)

func TestValidProgram(t *testing.T) {
	_ = NewCheckerTestUtil(t, `
		contract Test {
			int x
			bool b = true
			char c = 'c'
			string s = "hello"

			function void test() {
				int x = 2
			}
		}
	`, true)
}

func TestInvalidProgram(t *testing.T) {
	_ = NewCheckerTestUtil(t, `
		contract int {
		}
	`, false)
}
