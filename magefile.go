//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
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

	err = sh.RunV("protoc", "--go_out=./", "./proto/consensus/consensus.proto")
	err = sh.RunV("protoc", "--go_out=./", "./proto/smr/smr.proto")
	err = sh.RunV("protoc", "--go_out=./", "./proto/client/client.proto")

	if err != nil {
		return err
	}
	return nil
}

// Run tests.
func Test() error {
	mg.Deps(Generate)
	return sh.RunV("go", "test", "-v", "./...")
}

// Build binary executables.
func Build() error {
	err := sh.RunV("go", "build", "-v", "-o", "./proposer/bin/proposer", "./proposer/")
	err := sh.RunV("go", "build", "-v", "-o", "./recorder/bin/recorder", "./recorder/")
	err := sh.RunV("go", "build", "-v", "-o", "./client/traffic/bin/traffic", "./client/traffic/")
	err := sh.RunV("go", "build", "-v", "-o", "./client/attack/bin/attack", "./client/attack/")
	if err != nil {
		return err
	}

	return nil
}

// Remove binary executables.
func Clean() error {
	return os.RemoveAll("bin")
}
