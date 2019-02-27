package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printUsage(os.Stderr)
		return
	}

	compile(os.Args[1])
}

// TODO: Use cli library (e.g. cobra) to show help (available commands, flags and usage)
func printUsage(w io.Writer){
	fmt.Fprintln(w, "Lazo is a smart contract language for the Bazo Blockchain")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: \"lazo [source file]\"")
	fmt.Fprintln(w, "Example: \"lazo program.lazo\"")
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
