package dm

// Record Types.
const (
	FCGI_BEGIN_REQUEST     = 1
	FCGI_ABORT_REQUEST     = 2
	FCGI_END_REQUEST       = 3
	FCGI_PARAMS            = 4
	FCGI_STDIN             = 5
	FCGI_STDOUT            = 6
	FCGI_STDERR            = 7
	FCGI_DATA              = 8
	FCGI_GET_VALUES        = 9
	FCGI_GET_VALUES_RESULT = 10
	FCGI_UNKNOWN_TYPE      = 11
	FCGI_MAXTYPE           = FCGI_UNKNOWN_TYPE
)

type RecordType = byte
