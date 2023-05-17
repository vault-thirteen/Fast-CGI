package request

import (
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
)

/*
	typedef struct {
	    FCGI_Header header;
	    FCGI_BeginRequestBody body;
	} FCGI_BeginRequestRecord;
*/
type BeginRequest struct {
	Header dm.Header           // 8 bytes.
	Body   dm.BeginRequestBody // 8 bytes.
}

func NewBeginRequest(requestId uint16, role dm.Role, flags byte) (br *BeginRequest) {
	return &BeginRequest{
		Header: dm.Header{
			Version:       dm.FCGI_VERSION_1,
			Type:          dm.FCGI_BEGIN_REQUEST,
			RequestId:     requestId,
			ContentLength: 8, // Body is always 8 bytes long.
			PaddingLength: 0,
		},
		Body: dm.NewBeginRequestBody(role, flags),
	}
}

func (br *BeginRequest) ToBytes() (ba []byte) {
	ba = make([]byte, 0, 16)
	ba = append(ba, br.Header.ToBytes()...)
	ba = append(ba, br.Body.ToBytes()...)
	return ba
}
