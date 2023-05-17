package example

import (
	"bytes"
	"github.com/vault-thirteen/Fast-CGI/pkg/Client"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	dm2 "github.com/vault-thirteen/Fast-CGI/pkg/models/data"

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
		nvpair.NewNameValuePairWithTextValueU(dm2.Parameter_ScriptFilename, scriptFilePath),
	}
	stdin := []byte{}

	var tcpData bytes.Buffer
	var ba []byte

	ba = c.CreateBeginRequest(requestId, dm2.FCGI_RESPONDER, dm2.FCGI_KEEP_CONN)
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

	var recs []*dm2.Record
	recs, err = c.ReadResponseUntilEnd()
	if err != nil {
		return nil, nil, err
	}

	recs = dm2.FilterRecordsByRequestId(recs, requestId)

	stdOut = dm2.GetStdOutFromRecords(recs)
	stdErr = dm2.GetStdErrFromRecords(recs)

	return stdOut, stdErr, nil
}
