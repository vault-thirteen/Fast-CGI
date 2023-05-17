package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/vault-thirteen/Fast-CGI/example"
)

const (
	TestServerNetwork = "tcp"
	TestServerAddress = "127.0.0.1:9000"
)

const (
	ErrNotEnoughArguments = "not enough arguments"
)

// This example runs a simple PHP script and gets its output.
// Only the `SCRIPT_FILENAME` parameter is provided to the PHP script, that is
// why it is simple. The PHP-CGI server must be started manually before running
// this function.
func main() {
	var err error
	var scriptFilePath string
	scriptFilePath, err = getScriptFilePath()
	if err != nil {
		log.Println(err)
		showOutro()
	}

	err = runSimplePhpScript(scriptFilePath)
	mustBeNoError(err)
}

func getScriptFilePath() (scriptFilePath string, err error) {
	if len(os.Args) == 1 {
		return scriptFilePath, errors.New(ErrNotEnoughArguments)
	}

	scriptFilePath = os.Args[1]

	return scriptFilePath, nil
}

func showOutro() {
	fmt.Println("Usage: program.exe [Path-To-Script]")
	os.Exit(1)
}

func runSimplePhpScript(scriptFilePath string) (err error) {
	var stdOut, stdErr []byte
	stdOut, stdErr, err = example.RunSimplePhpScript(TestServerNetwork, TestServerAddress, scriptFilePath)
	if err != nil {
		return err
	}

	if len(stdErr) > 0 {
		_, err = fmt.Fprintln(os.Stderr, string(stdErr))
		if err != nil {
			return err
		}
	}

	if len(stdOut) > 0 {
		_, err = fmt.Fprintln(os.Stdout, string(stdOut))
		if err != nil {
			return err
		}
	}

	return nil
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
