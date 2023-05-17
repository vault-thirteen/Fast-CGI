package dm

// Flags.
const (
	FCGI_KEEP_CONN = 1
)

/*
	typedef struct {
		unsigned char roleB1;
		unsigned char roleB0;
		unsigned char flags;
		unsigned char reserved[5];
	} FCGI_BeginRequestBody;
*/
type BeginRequestBody struct {
	Role     Role
	Flags    byte
	Reserved [5]byte
}

func NewBeginRequestBody(role Role, flags byte) (brb BeginRequestBody) {
	return BeginRequestBody{
		Role:  role,
		Flags: flags,
	}
}

func (brb BeginRequestBody) ToBytes() (ba []byte) {
	ba = make([]byte, 0, 8)
	ba = append(ba, Uint16ToBytes(brb.Role)...)
	ba = append(ba, brb.Flags)
	ba = append(ba, brb.Reserved[:]...)
	return ba
}
