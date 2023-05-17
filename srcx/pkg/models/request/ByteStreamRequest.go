package request

import (
	"bytes"
	"errors"
	"math"

	"github.com/vault-thirteen/Fast-CGI/srcx/pkg/models/common"
	dm "github.com/vault-thirteen/Fast-CGI/srcx/pkg/models/data"
)

// ByteStreamRequest is a generic request for requests that use byte stream as
// a content, such as StdInRequest, DataRequest, StdOutRequest and
// StdErrRequest.
type ByteStreamRequest struct {
	Header dm.Header
	Bytes  []byte
}

func NewByteStreamRequest(requestType byte, requestId uint16, bytes []byte) (bsr *ByteStreamRequest, err error) {
	contentLength := len(bytes)
	if contentLength > math.MaxUint16 {
		return nil, errors.New(common.ErrContentIsTooLong)
	}

	bsr = &ByteStreamRequest{
		Header: dm.Header{
			Version:       dm.FCGI_VERSION_1,
			Type:          requestType,
			RequestId:     requestId,
			ContentLength: uint16(contentLength),
			PaddingLength: byte(dm.CalculatePadding(contentLength)),
		},
		Bytes: bytes,
	}

	return bsr, nil
}

func NewStdInRequest(requestId uint16, stdin []byte) (sir *ByteStreamRequest, err error) {
	return NewByteStreamRequest(dm.FCGI_STDIN, requestId, stdin)
}

func NewDataRequest(requestId uint16, data []byte) (dr *ByteStreamRequest, err error) {
	return NewByteStreamRequest(dm.FCGI_DATA, requestId, data)
}

func NewStdOutRequest(requestId uint16, stdout []byte) (sor *ByteStreamRequest, err error) {
	return NewByteStreamRequest(dm.FCGI_STDOUT, requestId, stdout)
}

func NewStdErrRequest(requestId uint16, stderr []byte) (ser *ByteStreamRequest, err error) {
	return NewByteStreamRequest(dm.FCGI_STDERR, requestId, stderr)
}

func (bsr *ByteStreamRequest) ToBytes() (ba []byte, err error) {
	var buf bytes.Buffer

	_, err = buf.Write(bsr.Header.ToBytes())
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(bsr.Bytes)
	if err != nil {
		return nil, err
	}

	err = dm.WritePaddingToBytesBuffer(&buf, bsr.Header.PaddingLength)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
