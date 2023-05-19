package example

import (
	"bytes"
	"errors"

	"github.com/vault-thirteen/Fast-CGI/pkg/Client"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/http"
	"github.com/vault-thirteen/errorz"
)

// RunSimplePhpScript runs a simple PHP script and gets its output.
// Only the `SCRIPT_FILENAME` parameter is provided to the PHP script, that is
// why it is simple. The PHP-CGI server must be started manually before running
// this function.
func RunSimplePhpScript(
	serverNetwork string,
	serverAddress string,
	scriptFilePath string,
) (stdOut []byte, stdErr []byte, err error) {
	var c *cl.Client
	c, err = cl.New(serverNetwork, serverAddress)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		derr := c.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()

	requestId := uint16(1)
	params := []*nvpair.NameValuePair{
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ScriptFilename, scriptFilePath),
	}
	stdin := []byte{}

	var tcpData bytes.Buffer
	var ba []byte

	ba = c.CreateBeginRequest(requestId, dm.FCGI_RESPONDER, dm.FCGI_KEEP_CONN)
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	ba, err = c.CreateParamsRequest(requestId, params)
	if err != nil {
		return nil, nil, err
	}
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	ba, err = c.CreateParamsRequest(requestId, nil)
	if err != nil {
		return nil, nil, err
	}
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	ba, err = c.CreateStdInRequest(requestId, stdin)
	if err != nil {
		return nil, nil, err
	}
	_, err = tcpData.Write(ba)
	if err != nil {
		return nil, nil, err
	}

	err = c.SendRequest(tcpData.Bytes())
	if err != nil {
		return nil, nil, err
	}

	var recs []*dm.Record
	recs, err = c.ReadResponseUntilEnd()
	if err != nil {
		return nil, nil, err
	}

	recs = dm.FilterRecordsByRequestId(recs, requestId)

	stdOut = dm.GetStdOutFromRecords(recs)
	stdErr = dm.GetStdErrFromRecords(recs)

	return stdOut, stdErr, nil
}

// RunSimplePhpScriptAndGetHttpData runs a simple PHP script, gets its output,
// splits the outpud into HTTP headers and HTTP body. Only the `SCRIPT_FILENAME`
// parameter is provided to the PHP script, that is why it is simple. The
// PHP-CGI server must be started manually before running this function.
func RunSimplePhpScriptAndGetHttpData(
	serverNetwork string,
	serverAddress string,
	scriptFilePath string,
) (headers []*http.Header, body []byte, err error) {
	var stdOut []byte
	var stdErr []byte
	stdOut, stdErr, err = RunSimplePhpScript(serverNetwork, serverAddress, scriptFilePath)
	if err != nil {
		return nil, nil, err
	}

	if len(stdErr) > 0 {
		return nil, nil, errors.New(string(stdErr))
	}

	return http.SplitHttpHeadersFromStdout(stdOut)
}
