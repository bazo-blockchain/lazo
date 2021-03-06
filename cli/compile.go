package cli

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/lazo/checker"
	"github.com/bazo-blockchain/lazo/checker/symbol"
	"github.com/bazo-blockchain/lazo/generator"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/lexer/token"
	"github.com/bazo-blockchain/lazo/parser"
	"github.com/bazo-blockchain/lazo/parser/node"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var stage string

func init() {
	rootCmd.AddCommand(compileCommand)

	compileCommand.Flags().StringVarP(
		&stage,
		"stage",
		"s",
		"g",
		"Compilation stage. \nAvailable stages: l=lexer, p=parser, c=checker, g=generator")
}

var compileCommand = &cobra.Command{
	Use:     "compile [source file]",
	Short:   "Compile the Lazo source code",
	Example: "  lazo compile program.lazo --stage=l",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		} else {
			compile(args[0])
		}
	},
}

// compile compiles the given Lazo source code into Bazo byte code.
func compile(sourceFile string) ([]byte, [][]byte) {
	file, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}

	lexer := scan(file)
	syntaxTree := parse(lexer)
	symbolTable := check(syntaxTree)
	return generate(symbolTable)
}

func scan(file io.Reader) *lexer.Lexer {
	lexer := lexer.New(bufio.NewReader(file))
	if stage == "l" {
		tok := lexer.NextToken()
		for {
			if ftok, ok := tok.(*token.FixToken); ok && ftok.Value == token.EOF {
				break
			}
			fmt.Println(tok)
			tok = lexer.NextToken()
		}
		os.Exit(0)
	}
	return lexer
}

func parse(l *lexer.Lexer) *node.ProgramNode {
	parser := parser.New(l)
	syntaxTree, errors := parser.ParseProgram()

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, errors)
		fmt.Println(syntaxTree)
		os.Exit(1)
	}

	if stage == "p" {
		fmt.Println(syntaxTree)
		os.Exit(0)
	}

	return syntaxTree
}

func check(syntaxTree *node.ProgramNode) *symbol.SymbolTable {
	checker := checker.New(syntaxTree)
	symbolTable, errors := checker.Run()

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, errors)
		os.Exit(1)
	}

	if stage == "c" {
		fmt.Println(symbolTable)
		os.Exit(0)
	}

	return symbolTable
}

func generate(symbolTable *symbol.SymbolTable) ([]byte, [][]byte) {
	generator := generator.New(symbolTable)
	metadata, errors := generator.Run()

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, errors)
		os.Exit(1)
	}

	return metadata.CreateContract()
}
