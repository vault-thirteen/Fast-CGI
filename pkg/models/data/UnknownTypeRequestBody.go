package dm

/*
	 typedef struct {
		unsigned char type;
		unsigned char reserved[7];
	} FCGI_UnknownTypeBody;
*/
type UnknownTypeRequestBody struct {
	Type     RecordType
	Reserved [7]byte
}

func NewUnknownTypeRequestBody(recordType RecordType) (utrb UnknownTypeRequestBody) {
	return UnknownTypeRequestBody{
		Type: recordType,
	}
}

func (utrb UnknownTypeRequestBody) ToBytes() (ba []byte) {
	ba = make([]byte, 0, 8)
	ba = append(ba, utrb.Type)
	ba = append(ba, utrb.Reserved[:]...)
	return ba
}
