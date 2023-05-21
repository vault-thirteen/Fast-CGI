package rm

import "github.com/vault-thirteen/Fast-CGI/pkg/models/data"

/* {FCGI_ABORT_REQUEST, R} */
type AbortRequest struct {
	Header dm.Header
}

func NewAbortRequest(requestId uint16) (br *AbortRequest) {
	return &AbortRequest{
		Header: dm.Header{
			Version:       dm.FCGI_VERSION_1,
			Type:          dm.FCGI_ABORT_REQUEST,
			RequestId:     requestId,
			ContentLength: 0,
			PaddingLength: 0,
		},
	}
}

func (ar *AbortRequest) ToBytes() (ba []byte) {
	return ar.Header.ToBytes()
}
