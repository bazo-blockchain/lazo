package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	ilParser "github.com/tk-codes/bazo-smartcontract/src/parser"
	bazoVM "github.com/tk-codes/bazo-smartcontract/src/vm"
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
	contract, err := ioutil.ReadFile(ilSourceFile)
	if err != nil {
		panic(err)
	}

	ilCode := ilParser.Parse(string(contract))
	vm := bazoVM.NewVM()
	vm.SetContractCode(ilCode)
	vm.Exec(true)
	fmt.Println(vm.GetResult())
}
