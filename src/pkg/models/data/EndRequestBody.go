package dm

// Protocol Status.
const (
	FCGI_REQUEST_COMPLETE = 0
	FCGI_CANT_MPX_CONN    = 1
	FCGI_OVERLOADED       = 2
	FCGI_UNKNOWN_ROLE     = 3
)

/*
	typedef struct {
		unsigned char appStatusB3;
		unsigned char appStatusB2;
		unsigned char appStatusB1;
		unsigned char appStatusB0;
		unsigned char protocolStatus;
		unsigned char reserved[3];
	} FCGI_EndRequestBody;
*/
type EndRequestBody struct {
	AppStatus      uint32
	ProtocolStatus byte
	Reserved       [3]byte
}

func NewEndRequestBody(appStatus uint32, protocolStatus byte) (erb EndRequestBody) {
	return EndRequestBody{
		AppStatus:      appStatus,
		ProtocolStatus: protocolStatus,
	}
}

func (erb EndRequestBody) ToBytes() (ba []byte) {
	ba = make([]byte, 0, 8)
	ba = append(ba, Uint32ToBytes(erb.AppStatus)...)
	ba = append(ba, erb.ProtocolStatus)
	ba = append(ba, erb.Reserved[:]...)
	return ba
}
