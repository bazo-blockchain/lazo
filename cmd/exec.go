package cmd

import (
	"fmt"
	"github.com/bazo-blockchain/bazo-vm/vm"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func init() {
	rootCmd.AddCommand(execCommand)
}

var execCommand = &cobra.Command{
	Use:     "exec [il file]",
	Short:   "Execute the enhanced Bazo byte code",
	Example: "lazo exec program.bc",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		} else {
			execute(args[0])
		}
	},
}

func execute(ilSourceFile string) {
	byteCode, err := ioutil.ReadFile(ilSourceFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(byteCode)
	//ilCode = []byte{
	//	vm.PUSH, 1, 0, 4, // push 1 argument, 0=positive, int 4
	//	vm.PUSH, 1, 0, 3,
	//	vm.ADD,			 // add 4 + 3 = 7
	//	vm.HALT,
	//}

	vm := vm.NewVM(vm.NewMockContext(byteCode))
	isSuccess := vm.Exec(true)
	if !isSuccess {
		panic("Code execution failed")
	}
	result, _ := vm.PeekResult()
	fmt.Printf("%d", result) // [0, 7] => +7
}
