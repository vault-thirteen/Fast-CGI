package dm

type ParameterName = string
type ParameterValue = []byte

const (
	Parameter_AuthType         = "AUTH_TYPE"
	Parameter_ContentLength    = "CONTENT_LENGTH"
	Parameter_ContentType      = "CONTENT_TYPE"
	Parameter_DocumentRoot     = "DOCUMENT_ROOT"
	Parameter_DocumentUri      = "DOCUMENT_URI"
	Parameter_GatewayInterface = "GATEWAY_INTERFACE"
	Parameter_PathInfo         = "PATH_INFO"
	Parameter_PathTranslated   = "PATH_TRANSLATED"
	Parameter_QueryString      = "QUERY_STRING"
	Parameter_RemoteAddr       = "REMOTE_ADDR"
	Parameter_RemoteHost       = "REMOTE_HOST"
	Parameter_RemoteIdent      = "REMOTE_IDENT"
	Parameter_RemotePort       = "REMOTE_PORT"
	Parameter_RemoteUser       = "REMOTE_USER"
	Parameter_RequestMethod    = "REQUEST_METHOD"
	Parameter_RequestUri       = "REQUEST_URI"
	Parameter_ScriptFilename   = "SCRIPT_FILENAME"
	Parameter_ScriptName       = "SCRIPT_NAME"
	Parameter_ServerAddr       = "SERVER_ADDR"
	Parameter_ServerName       = "SERVER_NAME"
	Parameter_ServerPort       = "SERVER_PORT"
	Parameter_ServerProtocol   = "SERVER_PROTOCOL"
	Parameter_ServerSoftware   = "SERVER_SOFTWARE"
)

// PHP Parameters.
const (
	Parameter_Https              = "HTTPS"
	Parameter_OrigPathInfo       = "ORIG_PATH_INFO"
	Parameter_PhpAuthDigest      = "PHP_AUTH_DIGEST"
	Parameter_PhpAuthPw          = "PHP_AUTH_PW"
	Parameter_PhpAuthUser        = "PHP_AUTH_USER"
	Parameter_PhpSelf            = "PHP_SELF"
	Parameter_RedirectRemoteUser = "REDIRECT_REMOTE_USER"
	Parameter_RedirectStatus     = "REDIRECT_STATUS"
	Parameter_RequestScheme      = "REQUEST_SCHEME"
	Parameter_RequestTimeFloat   = "REQUEST_TIME_FLOAT"
	Parameter_ServerAdmin        = "SERVER_ADMIN"
	Parameter_ServerSignature    = "SERVER_SIGNATURE"
)

const (
	ParameterPrefix_Http = "HTTP_"
)
