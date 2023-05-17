package request

import (
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
)

/*
	typedef struct {
	    FCGI_Header header;
	    FCGI_UnknownTypeBody body;
	} FCGI_UnknownTypeRecord;
*/
type UnknownTypeRequest struct {
	Header dm.Header
	Body   dm.UnknownTypeRequestBody
}

func NewUnknownTypeRequest(recordType dm.RecordType) (utr *UnknownTypeRequest) {
	return &UnknownTypeRequest{
		Header: dm.Header{
			Version:       dm.FCGI_VERSION_1,
			Type:          dm.FCGI_UNKNOWN_TYPE,
			RequestId:     dm.FCGI_NULL_REQUEST_ID,
			ContentLength: 8, // Body is always 8 bytes long.
			PaddingLength: 0,
		},
		Body: dm.NewUnknownTypeRequestBody(recordType),
	}
}

func (utr *UnknownTypeRequest) ToBytes() (ba []byte) {
	ba = make([]byte, 0, 16)
	ba = append(ba, utr.Header.ToBytes()...)
	ba = append(ba, utr.Body.ToBytes()...)
	return ba
}
