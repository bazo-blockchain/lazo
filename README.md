# Lazo: A Smart Contract Language for the Bazo Blockchain

[![Build Status](https://travis-ci.org/bazo-blockchain/lazo.svg?branch=master)](https://travis-ci.org/bazo-blockchain/lazo)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bazo-blockchain_lazo&metric=alert_status)](https://sonarcloud.io/dashboard?id=bazo-blockchain_lazo)

Lazo is a statically typed, imperative and non-turing complete programming language.
Please refer to [lazo-specification](https://github.com/bazo-blockchain/lazo-specification) for the complete language features.

## Usage

The Lazo tool works with the CLI commands.
Run `lazo` to see all the available commands and their usages.

    $ lazo
    Lazo is a tool for managing Lazo source code on the Bazo Blockchain
    
    Usage:
      lazo [flags]
      lazo [command]
    
    Available Commands:
      compile     Compile the Lazo source code
      help        Help about any command
      version     Print the version number of Lazo
    
    Flags:
      -h, --help   help for lazo
    
    Use "lazo [command] --help" for more information about a command.

## Development
###  Dependency Management

Packages are managed by [dep](https://golang.github.io/dep/). Install dep and run `dep ensure` to install all the dependencies.

### Run Compiler from Source

    go run main.go compile program.lazo

It will compile the given source code file "*program.lazo*".

### Run Unit Tests

    go test ./... 

It will run all tests in the current directory and all of its subdirectories.

To see the test coverage, run `./scripts/test.sh` and then open the **coverage.html** file.

### Run Lints

TODO

### Build Compiler

    go build 

It will create an executable for the current operating system (e.g. `lazo.exe` in Windows).

### Install Compiler

    go install
    
It will build an executable and place it in the `$GOPATH/bin` directory.
Thus, `lazo` command will be available in the terminal from anywhere.
