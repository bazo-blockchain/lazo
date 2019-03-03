# Lazo: A Smart Contract Language for the Bazo Blockchain

[![Build Status](https://travis-ci.org/bazo-blockchain/lazo.svg?branch=master)](https://travis-ci.org/bazo-blockchain/lazo)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bazo-blockchain_lazo&metric=alert_status)](https://sonarcloud.io/dashboard?id=bazo-blockchain_lazo)

Please refer to [lazo-specification](https://github.com/bazo-blockchain/lazo-specification) for the language features.

## Dependency Management

Packages are managed by [dep](https://golang.github.io/dep/). Install dep and run `dep ensure` to install all the dependencies.

## Run Compiler

`go run main.go program.lazo`

It will compile the given source code file "*program.lazo*".

## Build Compiler

`go build` 

It will create an executable for the current operating system (e.g. `lazo.exe` in Windows).

## Run Tests

`go test ./...` 

It will run all tests in the current directory and all of its subdirectories.

To see the test coverage, run `./scripts/test.sh` and then open the **coverage.html** file.  

