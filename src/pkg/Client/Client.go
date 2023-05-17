package cl

import (
	"net"

	"github.com/vault-thirteen/Fast-CGI/src/pkg/models/NameValuePair"
	"github.com/vault-thirteen/Fast-CGI/src/pkg/models/data"
	"github.com/vault-thirteen/Fast-CGI/src/pkg/models/request"
)

type Client struct {
	serverAddress *net.TCPAddr
	conn          *net.TCPConn
}

func New(network string, address string) (c *Client, err error) {
	c = &Client{}

	c.serverAddress, err = net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	c.conn, err = net.DialTCP(network, nil, c.serverAddress)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Close() (err error) {
	return c.conn.Close()
}

// 4.1.1. FCGI_GET_VALUES.
func (c *Client) CreateGetValuesRequest(params []*nvpair.NameValuePair) (ba []byte, err error) {
	var r *request.ValuesRequest
	r, err = request.NewGetValuesRequest(params)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

// 4.1.2. FCGI_GET_VALUES_RESULT.
func (c *Client) CreateGetValuesResultRequest(params []*nvpair.NameValuePair) (ba []byte, err error) {
	var r *request.ValuesRequest
	r, err = request.NewGetValuesResultRequest(params)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

// 4.2. FCGI_UNKNOWN_TYPE.
func (c *Client) CreateUnknownTypeRequest(recordType dm.RecordType) (ba []byte) {
	return request.NewUnknownTypeRequest(recordType).ToBytes()
}

// 5.1. FCGI_BEGIN_REQUEST.
func (c *Client) CreateBeginRequest(requestId uint16, role dm.Role, flags byte) (ba []byte) {
	return request.NewBeginRequest(requestId, role, flags).ToBytes()
}

// 5.2. FCGI_PARAMS.
func (c *Client) CreateParamsRequest(requestId uint16, params []*nvpair.NameValuePair) (ba []byte, err error) {
	var r *request.ValuesRequest
	r, err = request.NewParamsRequest(requestId, params)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

// 5.3.1. FCGI_STDIN.
func (c *Client) CreateStdInRequest(requestId uint16, stdin []byte) (ba []byte, err error) {
	var r *request.ByteStreamRequest
	r, err = request.NewStdInRequest(requestId, stdin)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

// 5.3.2. FCGI_DATA.
func (c *Client) CreateDataRequest(requestId uint16, data []byte) (ba []byte, err error) {
	var r *request.ByteStreamRequest
	r, err = request.NewDataRequest(requestId, data)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

// 5.3.3. FCGI_STDOUT.
func (c *Client) CreateStdOutRequest(requestId uint16, stdout []byte) (ba []byte, err error) {
	var r *request.ByteStreamRequest
	r, err = request.NewStdOutRequest(requestId, stdout)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

// 5.3.4. FCGI_STDERR.
func (c *Client) CreateStdErrRequest(requestId uint16, stderr []byte) (ba []byte, err error) {
	var r *request.ByteStreamRequest
	r, err = request.NewStdErrRequest(requestId, stderr)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

// 5.4. FCGI_ABORT_REQUEST.
func (c *Client) CreateAbortRequest(requestId uint16) (ba []byte) {
	return request.NewAbortRequest(requestId).ToBytes()
}

// 5.5. FCGI_END_REQUEST.
func (c *Client) CreateEndRequest(requestId uint16, appStatus uint32, protocolStatus byte) (ba []byte) {
	return request.NewEndRequest(requestId, appStatus, protocolStatus).ToBytes()
}

func (c *Client) SendRequest(data []byte) (err error) {
	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ReadRawRecord() (r *dm.Record, err error) {
	return dm.NewRecordFromStream(c.conn)
}

func (c *Client) ReadResponseUntilEnd() (recs []*dm.Record, err error) {
	recs = make([]*dm.Record, 0)
	var rec *dm.Record

	for {
		rec, err = dm.NewRecordFromStream(c.conn)
		if err != nil {
			return nil, err
		}

		recs = append(recs, rec)

		if rec.Type == dm.FCGI_END_REQUEST {
			break
		}
	}

	return recs, nil
}
