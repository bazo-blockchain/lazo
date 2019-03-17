package checker

import (
	"testing"
)

func TestFunctionReturnBoolConstant(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function bool test() {
				return true
			}
		}
		`, true)
}

// TODO: fix
//func TestFunctionReturnBoolFail(t *testing.T) {
//	_ = newCheckerTestUtilWithRawInput(t, `
//		contract Test {
//			function bool test() {
//				bool b = 5
//				return b
//			}
//		}
//		`, false)
//}

func TestFunctionReturnInt(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function int test() {
				int i = 5
				return 5
			}
		}
		`, true)
}

func TestFunctionReturnString(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function string test() {
				string s = "test"
				return s
			}
		}
		`, true)
}

func TestFunctionReturnChar(t *testing.T) {
	_ = newCheckerTestUtilWithRawInput(t, `
		contract Test {
			function char test() {
				char c = 'c'
				return c
			}
		}
		`, true)
}
