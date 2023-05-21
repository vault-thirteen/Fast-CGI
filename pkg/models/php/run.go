package pm

import (
	"bytes"
	"errors"

	"github.com/vault-thirteen/Fast-CGI/pkg/Client"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
	"github.com/vault-thirteen/errorz"
)

// RunOnceSimplePhpScript runs a simple PHP script once and gets its output.
// Only the `SCRIPT_FILENAME` parameter is provided to the PHP script, that is
// why it is simple. The PHP-CGI server must be started manually before running
// this function.
func RunOnceSimplePhpScript(serverNetwork string, serverAddress string, scriptFilePath string) (stdOut []byte, stdErr []byte, err error) {
	parameters := []*nvpair.NameValuePair{
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ScriptFilename, scriptFilePath),
	}

	return RunOncePhpScript(serverNetwork, serverAddress, 1, parameters, []byte{})
}

// RunOnceSimplePhpScriptAndGetHttpData runs a simple PHP script once, gets its
// output, splits the output into HTTP headers and HTTP body. Only the
// `SCRIPT_FILENAME` parameter is provided to the PHP script, that is why it is
// simple. The PHP-CGI server must be started manually before running this
// function.
func RunOnceSimplePhpScriptAndGetHttpData(serverNetwork string, serverAddress string, scriptFilePath string) (data *Data, err error) {
	var stdOut []byte
	var stdErr []byte
	stdOut, stdErr, err = RunOnceSimplePhpScript(serverNetwork, serverAddress, scriptFilePath)
	if err != nil {
		return nil, err
	}

	if len(stdErr) > 0 {
		return nil, errors.New(string(stdErr))
	}

	return SplitHeadersFromStdout(stdOut)
}

// RunOncePhpScript runs a PHP script once.
// Path to the script file must be set as a 'SCRIPT_FILENAME' parameter inside
// the 'parameters' argument.
func RunOncePhpScript(serverNetwork string, serverAddress string, requestId uint16, parameters []*nvpair.NameValuePair, stdin []byte) (stdOut []byte, stdErr []byte, err error) {
	var client *cl.Client
	client, err = cl.New(serverNetwork, serverAddress)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		derr := client.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()

	return ExecPhpScript(client, requestId, parameters, stdin)
}

// RunOncePhpScriptAndGetHttpData runs a PHP script once, gets its output,
// splits the output into HTTP headers and HTTP body. Path to the script file
// must be set as a 'SCRIPT_FILENAME' parameter inside the 'parameters'
// argument. The PHP-CGI server must be started manually before running this
// function.
func RunOncePhpScriptAndGetHttpData(serverNetwork string, serverAddress string, requestId uint16, parameters []*nvpair.NameValuePair, stdin []byte) (data *Data, err error) {
	var stdOut []byte
	var stdErr []byte
	stdOut, stdErr, err = RunOncePhpScript(serverNetwork, serverAddress, requestId, parameters, stdin)
	if err != nil {
		return nil, err
	}

	if len(stdErr) > 0 {
		return nil, errors.New(string(stdErr))
	}

	return SplitHeadersFromStdout(stdOut)
}

// ExecPhpScript executes a PHP script using the specified client.
// Path to the script file must be set as a 'SCRIPT_FILENAME' parameter inside
// the 'parameters' argument.
func ExecPhpScript(client *cl.Client, requestId uint16, parameters []*nvpair.NameValuePair, stdin []byte) (stdOut []byte, stdErr []byte, err error) {
	var tcpData bytes.Buffer
	var ba []byte

	ba = client.CreateBeginRequest(requestId, dm.FCGI_RESPONDER, dm.FCGI_KEEP_CONN)
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	ba, err = client.CreateParamsRequest(requestId, parameters)
	if err != nil {
		return nil, nil, err
	}
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	ba, err = client.CreateParamsRequest(requestId, nil)
	if err != nil {
		return nil, nil, err
	}
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	ba, err = client.CreateStdInRequest(requestId, stdin)
	if err != nil {
		return nil, nil, err
	}
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	err = client.SendRequest(tcpData.Bytes())
	if err != nil {
		return nil, nil, err
	}

	var recs []*dm.Record
	recs, err = client.ReadResponseUntilEnd()
	if err != nil {
		return nil, nil, err
	}

	recs = dm.FilterRecordsByRequestId(recs, requestId)

	stdOut = dm.GetStdOutFromRecords(recs)
	stdErr = dm.GetStdErrFromRecords(recs)

	return stdOut, stdErr, nil
}

// ExecPhpScriptAndGetHttpData executes a PHP script using the specified
// client, gets its output, splits the output into HTTP headers and HTTP body.
// Path to the script file must be set as a 'SCRIPT_FILENAME' parameter inside
// the 'parameters' argument. The PHP-CGI server must be started manually
// before running this function.
func ExecPhpScriptAndGetHttpData(client *cl.Client, requestId uint16, parameters []*nvpair.NameValuePair, stdin []byte) (data *Data, err error) {
	var stdOut []byte
	var stdErr []byte
	stdOut, stdErr, err = ExecPhpScript(client, requestId, parameters, stdin)
	if err != nil {
		return nil, err
	}

	if len(stdErr) > 0 {
		return nil, errors.New(string(stdErr))
	}

	return SplitHeadersFromStdout(stdOut)
}
