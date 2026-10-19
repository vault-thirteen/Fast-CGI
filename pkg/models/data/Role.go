package dm

// Roles.
const (
	FCGI_RESPONDER  = 1
	FCGI_AUTHORIZER = 2
	FCGI_FILTER     = 3
)

type Role = uint16
