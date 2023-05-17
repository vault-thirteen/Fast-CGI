package dm

import (
	"bytes"
	"encoding/binary"

	nvpair "github.com/vault-thirteen/Fast-CGI/src/pkg/models/NameValuePair"
)

// Uint16ToBytes converts an uint16 into an array (slice) of bytes using the
// Big-Endian byte order.
func Uint16ToBytes(x uint16) (ba []byte) {
	ba = []byte{0, 0}
	binary.BigEndian.PutUint16(ba, x)
	return ba
}

// Uint32ToBytes converts an uint32 into an array (slice) of bytes using the
// Big-Endian byte order.
func Uint32ToBytes(x uint32) (ba []byte) {
	ba = []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(ba, x)
	return ba
}

// CalculatePadding calculates the required padding size.
// FastCGI Specification recommends aligning data by 8 bytes.
func CalculatePadding(dataLength int) (padding int) {
	mod := dataLength % 8
	if mod == 0 {
		return 0
	}
	return 8 - mod
}

// PaddingToBytes prepares the padding bytes.
func PaddingToBytes(paddingLength byte) (ba []byte, err error) {
	if paddingLength == 0 {
		return []byte{}, nil
	}

	var buf bytes.Buffer

	padding := make([]byte, paddingLength)
	for i := 0; i < len(padding); i++ {
		padding[i] = 0
	}

	_, err = buf.Write(padding)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// WriteParametersToBytesBuffer puts parameters into a buffer of bytes.
func WriteParametersToBytesBuffer(buf *bytes.Buffer, parameters []*nvpair.NameValuePair) (err error) {
	var ba []byte
	for _, p := range parameters {
		ba, err = p.ToBytes()
		if err != nil {
			return err
		}

		_, err = buf.Write(ba)
		if err != nil {
			return err
		}
	}

	return nil
}

// WritePaddingToBytesBuffer puts th padding into a buffer of bytes.
func WritePaddingToBytesBuffer(buf *bytes.Buffer, paddingLength byte) (err error) {
	var ba []byte
	ba, err = PaddingToBytes(paddingLength)
	if err != nil {
		return err
	}

	_, err = buf.Write(ba)
	if err != nil {
		return err
	}

	return nil
}

// FilterRecordsByRequestId returns only those records, that have the specified
// Request ID.
func FilterRecordsByRequestId(recsInput []*Record, requestId uint16) (recsOutput []*Record) {
	recsOutput = make([]*Record, 0, len(recsInput))

	for _, rec := range recsInput {
		if rec.RequestId == requestId {
			recsOutput = append(recsOutput, rec)
		}
	}

	return recsOutput
}

// GetContentsFromRecordsByType extracts contents from records having the
// specified type, concatenates it and returns it.
func GetContentsFromRecordsByType(recs []*Record, recordType RecordType) (contents []byte) {
	var buf bytes.Buffer

	for _, rec := range recs {
		if rec.Type == recordType {
			buf.Write(rec.ContentData)
		}
	}

	return buf.Bytes()
}

// GetStdErrFromRecords extracts stderr records from the list of records and
// concatenates their content.
func GetStdErrFromRecords(recs []*Record) (stderr []byte) {
	return GetContentsFromRecordsByType(recs, FCGI_STDERR)
}

// GetStdOutFromRecords extracts stdout records from the list of records and
// concatenates their content.
func GetStdOutFromRecords(recs []*Record) (stdout []byte) {
	return GetContentsFromRecordsByType(recs, FCGI_STDOUT)
}
