package cli

import (
	"fmt"
	"github.com/bazo-blockchain/bazo-vm/vm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCommand)
}

var runCommand = &cobra.Command{
	Use:     "run [source file]",
	Short:   "Compile and run the lazo source code on Bazo VM",
	Example: "lazo run program.lazo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		} else {
			execute(args[0])
		}
	},
}

func execute(sourceFile string) {
	code, variables := Compile(sourceFile)
	context := vm.NewMockContext(code)
	context.ContractVariables = variables
	context.Fee += (uint64(len(variables))) * 1000
	context.Data = []byte{
		1, // total bytes
		0, // Contract Init Flag
		//3,
		//19, 70, 101, 78, // Function hash
	}

	vm := vm.NewVM(context)
	isSuccess := vm.Exec(true)
	result, _ := vm.PeekResult()
	if !isSuccess {
		panic(fmt.Sprintf("Runtime Error: %s", result))
	}

	fmt.Printf("%d", result) // [0, 7] => +7
}
