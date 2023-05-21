package rm

import (
	"bytes"
	"errors"
	"math"

	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/common"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
)

// ValuesRequest is a generic request for requests that use values, a.k.a.
// name-value pairs, as a content, such as GetValuesRequest, ParamsRequest and
// GetValuesResultRequest.
type ValuesRequest struct {
	Header dm.Header
	Values []*nvpair.NameValuePair
}

func NewValuesRequest(requestType byte, requestId uint16, values []*nvpair.NameValuePair) (vr *ValuesRequest, err error) {
	contentLength := nvpair.MeasureNameValuePairs(values)
	if contentLength > math.MaxUint16 {
		return nil, errors.New(cm.ErrContentIsTooLong)
	}

	vr = &ValuesRequest{
		Header: dm.Header{
			Version:       dm.FCGI_VERSION_1,
			Type:          requestType,
			RequestId:     requestId,
			ContentLength: uint16(contentLength),
			PaddingLength: byte(dm.CalculatePadding(contentLength)),
		},
		Values: values,
	}

	return vr, nil
}

func NewGetValuesRequest(parameters []*nvpair.NameValuePair) (gvr *ValuesRequest, err error) {
	return NewValuesRequest(dm.FCGI_GET_VALUES, dm.FCGI_NULL_REQUEST_ID, parameters)
}

func NewGetValuesResultRequest(parameters []*nvpair.NameValuePair) (gvrr *ValuesRequest, err error) {
	return NewValuesRequest(dm.FCGI_GET_VALUES_RESULT, dm.FCGI_NULL_REQUEST_ID, parameters)
}

func NewParamsRequest(requestId uint16, params []*nvpair.NameValuePair) (pr *ValuesRequest, err error) {
	return NewValuesRequest(dm.FCGI_PARAMS, requestId, params)
}

func (vr *ValuesRequest) ToBytes() (ba []byte, err error) {
	var buf bytes.Buffer

	_, err = buf.Write(vr.Header.ToBytes())
	if err != nil {
		return nil, err
	}

	err = dm.WriteParametersToBytesBuffer(&buf, vr.Values)
	if err != nil {
		return nil, err
	}

	err = dm.WritePaddingToBytesBuffer(&buf, vr.Header.PaddingLength)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
