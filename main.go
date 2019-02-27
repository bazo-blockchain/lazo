package main

import (
	"fmt"
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
	_, err := os.Open(sourceFile)
	check(err)
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
