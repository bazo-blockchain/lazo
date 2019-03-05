package main

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer"
	"github.com/bazo-blockchain/lazo/parser"
	"io"
	"os"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "help" {
		printUsage(os.Stderr)
		return
	}

	compile(os.Args[1])
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
	program, errors := parser.ParseProgram()
	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, errors)
	}
	fmt.Println(program)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Lazo is a smart contract language for the Bazo Blockchain")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: \"lazo [source file]\"")
	fmt.Fprintln(w, "Example: \"lazo program.lazo\"")
}
