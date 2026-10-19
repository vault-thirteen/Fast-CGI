package main

import (
	"fmt"
	"os"

	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/Fast-CGI/pkg/models/php"
)

const (
	TestServerNetwork = "tcp"
	TestServerAddress = "127.0.0.1:9000"
)

const (
	Err_NotEnoughArguments = "not enough arguments"
)

const (
	Msg_Usage       = "Usage: program.exe [Path-To-Script]"
	Msg_HttpHeaders = "HTTP headers:"
	Msg_HttpBody    = "HTTP body:"
	Msg_HLine       = "--------------------------------------------------------------------------------"

	Msgf_HttpStatus = "HTTP status: %v (%v)."
	Msgf_NV         = "[%v] = [%v]"
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
		errors.LogError(err)
		showOutro()
	}

	//err = runSimplePhpScript(scriptFilePath)
	err = runSimplePhpScriptAndGetHttpData(scriptFilePath)
	mustBeNoError(err)
}

func getScriptFilePath() (scriptFilePath string, err error) {
	if len(os.Args) == 1 {
		return scriptFilePath, errors.ErrorFromText(Err_NotEnoughArguments)
	}

	scriptFilePath = os.Args[1]

	return scriptFilePath, nil
}

func showOutro() {
	cio.PrintMessage(Msg_Usage)
	os.Exit(1)
}

func runSimplePhpScript(scriptFilePath string) (err error) {
	var stdOut, stdErr []byte
	stdOut, stdErr, err = pm.RunOnceSimplePhpScript(TestServerNetwork, TestServerAddress, scriptFilePath)
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

func runSimplePhpScriptAndGetHttpData(scriptFilePath string) (err error) {
	var httpData *pm.Data
	httpData, err = pm.RunOnceSimplePhpScriptAndGetHttpData(TestServerNetwork, TestServerAddress, scriptFilePath)
	if err != nil {
		return err
	}

	cio.PrintMessage(cio.FormatMessage(Msgf_HttpStatus, httpData.StatusCode, httpData.StatusText))
	cio.PrintMessage()

	cio.PrintMessage(Msg_HttpHeaders)
	cio.PrintMessage(Msg_HLine)
	for _, hdr := range httpData.Headers {
		cio.PrintMessage(cio.FormatMessage(Msgf_NV, hdr.Name, hdr.Value))
	}
	cio.PrintMessage(Msg_HLine)
	cio.PrintMessage()

	cio.PrintMessage(Msg_HttpBody)
	cio.PrintMessage(Msg_HLine)
	cio.PrintMessage(string(httpData.Body))
	cio.PrintMessage(Msg_HLine)
	cio.PrintMessage()

	return nil
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
