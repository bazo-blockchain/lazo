package cmd

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

	vm := vm.NewVM(context)
	isSuccess := vm.Exec(true)
	if !isSuccess {
		panic("Code execution failed")
	}
	result, _ := vm.PeekResult()
	fmt.Printf("%d", result) // [0, 7] => +7
}
