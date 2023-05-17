package request

import (
	"github.com/vault-thirteen/Fast-CGI/src/pkg/models/data"
)

/*
	typedef struct {
	    FCGI_Header header;
	    FCGI_EndRequestBody body;
	} FCGI_EndRequestRecord;
*/
type EndRequest struct {
	Header dm.Header         // 8 bytes.
	Body   dm.EndRequestBody // 8 bytes.
}

func NewEndRequest(requestId uint16, appStatus uint32, protocolStatus byte) (er *EndRequest) {
	return &EndRequest{
		Header: dm.Header{
			Version:       dm.FCGI_VERSION_1,
			Type:          dm.FCGI_END_REQUEST,
			RequestId:     requestId,
			ContentLength: 8, // Body is always 8 bytes long.
			PaddingLength: 0,
		},
		Body: dm.NewEndRequestBody(appStatus, protocolStatus),
	}
}

func (er *EndRequest) ToBytes() (ba []byte) {
	ba = make([]byte, 0, 16)
	ba = append(ba, er.Header.ToBytes()...)
	ba = append(ba, er.Body.ToBytes()...)
	return ba
}
