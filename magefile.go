//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

var (
	Default = Build
)

// Install build dependencies.
func BuildDeps() error {
	err := sh.RunV("protoc", "--version")
	if err != nil {
		return err
	}
	err = sh.RunV("go", "get", "-u", "github.com/golang/protobuf/protoc-gen-go")
	if err != nil {
		return err
	}
	err = sh.RunV("go", "get", "-u", "google.golang.org/grpc")
	if err != nil {
		return err
	}

	return nil
}

// Install dependencies.
func Deps() error {
	err := sh.RunV("go", "mod", "vendor")
	if err != nil {
		return err
	}

	return nil
}

// Generate code.
func Generate() error {

	err := sh.RunV("protoc", "--go_out=./", "--go_opt=paths=source_relative", "--go-grpc_out=.", "--go-grpc_opt=paths=source_relative", "./replica/src/consensus.proto")
	err = sh.RunV("protoc", "--go_out=./", "./proto/client/client.proto")

	if err != nil {
		return err
	}
	return nil
}

// Build binary executables.
func Build() error {
	err := sh.RunV("go", "build", "-v", "-o", "./replica/bin/replica", "./replica/")
	err = sh.RunV("go", "build", "-v", "-o", "./client/bin/client", "./client/")
	if err != nil {
		return err
	}

	return nil
}
