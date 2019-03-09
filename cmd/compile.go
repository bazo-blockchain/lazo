package cmd

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/lazo/checker"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(compileCommand)
}

var compileCommand = &cobra.Command{
	Use:   "compile",
	Short: "Compile the Lazo source code",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		} else {
			compile(args[0])
		}
	},
}

func compile(sourceFile string) {
	file, err := os.Open(sourceFile)
	check(err)

	lexer := lexer.New(bufio.NewReader(file))
	//for !lexer.IsEnd {
	//	tok := lexer.NextToken()
	//	fmt.Printf("%s \n", tok)
	//}

	parser := parser.New(lexer)
	syntaxTree, errors := parser.ParseProgram()
	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, errors)
	}

	checker := checker.New(syntaxTree)
	symbolTable, errors := checker.Run()
	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, errors)
	}
	fmt.Println(symbolTable)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}