package dm

import (
	"bytes"
	"encoding/binary"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	"io"

	"github.com/vault-thirteen/auxie/reader"
)

/*
	typedef struct {
		unsigned char version;
		unsigned char type;
		unsigned char requestIdB1;
		unsigned char requestIdB0;
		unsigned char contentLengthB1;
		unsigned char contentLengthB0;
		unsigned char paddingLength;
		unsigned char reserved;
		unsigned char contentData[contentLength];
		unsigned char paddingData[paddingLength];
	} FCGI_Record;
*/
type Record struct {
	Version       byte
	Type          RecordType
	RequestId     uint16
	ContentLength uint16
	PaddingLength byte
	Reserved      byte

	ContentData []byte
	PaddingData []byte
}

// Notes on RequestID.
//
// *	A Zero RequestId means a Management Record.
//		FCGI_NULL_REQUEST_ID
//
// *	A Non-Zero RequestId means an Application Record.

// Discrete Record is a single record.
// Stream Record is a single record of a sequence of stream records.
// Stream records are followed by an "End-of-Stream" record which has 0 length.

// Management Record Types:
//
// 1.	FCGI_GET_VALUES, FCGI_GET_VALUES_RESULT
//		ContentData of both records is a sequence of Name-Value Pairs.
//
// 2.	FCGI_UNKNOWN_TYPE

// Application Record Types:
//
// 1.	FCGI_BEGIN_REQUEST
//
// 2.	FCGI_PARAMS
//		ContentData of both records is a sequence of Name-Value Pairs.
//
// 3.	FCGI_STDIN, FCGI_DATA, FCGI_STDOUT, FCGI_STDERR
//		Byte Streams.
//
// 4.	FCGI_ABORT_REQUEST
//
// 5.	FCGI_END_REQUEST

func NewRecordFromStream(stream io.Reader) (rec *Record, err error) {
	rec = &Record{}
	rdr := reader.New(stream)

	rec.Version, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}

	rec.Type, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}

	rec.RequestId, err = rdr.ReadWord_BE()
	if err != nil {
		return nil, err
	}

	rec.ContentLength, err = rdr.ReadWord_BE()
	if err != nil {
		return nil, err
	}

	rec.PaddingLength, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}

	rec.Reserved, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}

	rec.ContentData, err = rdr.ReadBytes(int(rec.ContentLength))
	if err != nil {
		return nil, err
	}

	rec.PaddingData, err = rdr.ReadBytes(int(rec.PaddingLength))
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (r *Record) ToBytes() (ba []byte, err error) {
	var buf bytes.Buffer

	_, err = buf.Write([]byte{r.Version})
	if err != nil {
		return nil, err
	}

	_, err = buf.Write([]byte{r.Type})
	if err != nil {
		return nil, err
	}

	var buf2B = []byte{0, 0}
	binary.BigEndian.PutUint16(buf2B, r.RequestId)
	_, err = buf.Write(buf2B)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(buf2B, r.ContentLength)
	_, err = buf.Write(buf2B)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write([]byte{r.PaddingLength})
	if err != nil {
		return nil, err
	}

	_, err = buf.Write([]byte{r.Reserved})
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(r.ContentData)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(r.PaddingData)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Record) ParseContentAsNVPs() (nvps []*nvpair.nvpair, err error) {
	if r.ContentLength == 0 {
		return nil, nil
	}

	nvps = make([]*nvpair.NameValuePair, 0)
	bytesProcessed := 0
	rdr := bytes.NewReader(r.ContentData)

	var nvp *nvpair.NameValuePair
	for bytesProcessed < int(r.ContentLength) {
		nvp, err = nvpair.NewNameValuePairFromStream(rdr)
		if err != nil {
			return nil, err
		}

		nvps = append(nvps, nvp)
		bytesProcessed += nvp.Measure()
	}

	return nvps, nil
}
