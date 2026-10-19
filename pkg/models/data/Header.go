package dm

const (
	FCGI_HEADER_LEN      = 8
	FCGI_VERSION_1       = 1
	FCGI_NULL_REQUEST_ID = 0
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
	} FCGI_Header;
*/
type Header struct {
	Version       byte
	Type          RecordType
	RequestId     uint16
	ContentLength uint16
	PaddingLength byte
	Reserved      byte
}

func (h Header) ToBytes() (ba []byte) {
	ba = make([]byte, 0, 8)
	ba = append(ba, h.Version)
	ba = append(ba, h.Type)
	ba = append(ba, Uint16ToBytes(h.RequestId)...)
	ba = append(ba, Uint16ToBytes(h.ContentLength)...)
	ba = append(ba, h.PaddingLength)
	ba = append(ba, h.Reserved)
	return ba
}
