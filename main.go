package main

import (
	"bufio"
	"fmt"
	"github.com/bazo-blockchain/lazo/lexer"
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
	for !lexer.EOF {
		tok := lexer.NextToken()
		fmt.Printf("%s \n", tok)
	}
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
