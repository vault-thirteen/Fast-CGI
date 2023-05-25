package ws

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/http"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/php"
	"github.com/vault-thirteen/header"
)

const (
	HostPortDelimiter  = ":"
	GolangNetNetworkIP = "ip" // These constants should be exported by Golang ! Source: net/iprawsock.go.
)

const (
	ExtraPathSingleSlash     = `/`
	ExtraPathInstallerStatus = `/installer/status`
)

func (srv *Server) isExtOfPhpScript(ext string) bool {
	for _, phpExt := range srv.settings.PhpFileExtensions {
		if ext == phpExt {
			return true
		}
	}

	return false
}

func (srv *Server) prepareInputDataToRunPhpScript(req *http.Request, psi *pm.PhpScriptInfo) (stdin []byte, parameters []*nvpair.NameValuePair, err error) {
	stdin, err = io.ReadAll(req.Body)
	if err != nil {
		return nil, nil, err
	}
	if len(stdin) != int(req.ContentLength) {
		return nil, nil, errors.New(ErrContentLengthMismatch)
	}

	var authScheme string
	authScheme, _, err = hm.ParseAuthorizationHeader(req.Header.Get(header.HttpHeaderAuthorization))
	if err != nil {
		return nil, nil, err
	}

	remoteAddrParts := strings.Split(req.RemoteAddr, HostPortDelimiter)
	if len(remoteAddrParts) != 2 {
		return nil, nil, errors.New(ErrRemoteAddrParts)
	}

	var serverIPAddr *net.IPAddr
	serverIPAddr, err = net.ResolveIPAddr(GolangNetNetworkIP, srv.settings.ServerHost)
	if err != nil {
		return nil, nil, err
	}

	var ossd = &pm.OldSchoolStyleData{}
	if len(psi.QueryParamExtraPath) > 0 {
		// If extra path is set as a query parameter, we are in a compatibility
		// mode. We need to emulate the old-school-style CGI request.
		ossd.DocumentUri = req.URL.Path + psi.QueryParamExtraPath
		ossd.CgiExtraPath = psi.QueryParamExtraPath
		ossd.QueryString = ""
		ossd.RequestUri = req.RequestURI[:strings.Index(req.RequestURI, "?")] + psi.QueryParamExtraPath
	} else {
		ossd.DocumentUri = req.URL.Path
		ossd.CgiExtraPath = psi.UrlExtraPath
		ossd.QueryString = req.URL.RawQuery
		ossd.RequestUri = req.RequestURI

	}

	parameters = []*nvpair.NameValuePair{
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_AuthType, authScheme),                                      // 4.1.1.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ContentLength, strconv.Itoa(len(stdin))),                   // 4.1.2.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ContentType, req.Header.Get(header.HttpHeaderContentType)), // 4.1.3.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_DocumentRoot, srv.settings.DocumentRootPath),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_DocumentUri, ossd.DocumentUri),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_GatewayInterface, srv.settings.GatewayInterface), // 4.1.4.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_PathInfo, ossd.CgiExtraPath),                     // 4.1.5.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_PathTranslated, psi.FileAbsExtraPath),            // 4.1.6.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_QueryString, ossd.QueryString), // 4.1.7.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RedirectStatus, strconv.Itoa(http.StatusOK)),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteAddr, remoteAddrParts[0]), // Host. 4.1.8.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteHost, ""),                 // FQDN. 4.1.9.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteIdent, ""),                // 4.1.10.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemotePort, remoteAddrParts[1]), // Port.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteUser, ""),                 // 4.1.11.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RequestMethod, req.Method),     // 4.1.12.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RequestScheme, req.URL.Scheme), // Apache HTTP Server Header.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RequestUri, ossd.RequestUri),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ScriptFilename, psi.FileAbsPath),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ScriptName, psi.FileName), // 4.1.13.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ServerAddr, serverIPAddr.String()),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ServerName, srv.settings.ServerName),         // 4.1.14.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ServerPort, srv.settings.ServerPort),         // 4.1.15.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ServerProtocol, req.Proto),                   // 4.1.16.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ServerSoftware, srv.settings.ServerSoftware), // 4.1.17.
		// HTTP_XXX // 4.1.18.  Protocol-Specific Meta-Variables
	}

	// Add Client's HTTP Headers.
	hm.AddHttpHeadersToParameters(&parameters, req.Header)

	// There is a known bug or vulnerability with 'HTTP_HOST' header. It is
	// recommended to ignore this header in PHP while a client is able to
	// change this header manually. In any case, the server always sends its
	// address using the following variables: 'SERVER_ADDR', 'SERVER_NAME' and
	// 'SERVER_PORT', which should be used instead of client's 'HTTP_HOST'
	// header.

	return stdin, parameters, nil
}

func (srv *Server) runPhpScript(rw http.ResponseWriter, req *http.Request, psi *pm.PhpScriptInfo) {
	var requestId uint16 = 1
	var parameters []*nvpair.NameValuePair
	var stdin []byte
	var err error
	stdin, parameters, err = srv.prepareInputDataToRunPhpScript(req, psi)
	if err != nil {
		srv.respondWithInternalServerError(rw, err)
		return
	}

	//nvpair.PrintParameters(parameters) // DEBUG.

	var phpScriptOutput *pm.Data
	var phpErr error
	phpScriptOutput, phpErr = pm.ExecPhpScriptAndGetHttpData(srv.cgiClient, requestId, parameters, stdin)
	if phpErr != nil {
		// Headers.
		rw.Header().Set(header.HttpHeaderServer, srv.settings.ServerSoftware)

		// Status.
		rw.WriteHeader(http.StatusInternalServerError)

		// Body.
		_, err = rw.Write([]byte(phpErr.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	if srv.settings.FixRelativeRedirects {
		err = phpScriptOutput.FixLocationHeader(req.URL.Path)
		if err != nil {
			srv.respondWithInternalServerError(rw, err)
			return
		}
	}

	// Headers.
	if (srv.settings.IsCgiExtraPathEnabled) && (len(psi.UrlExtraPath) > 0) {
		rw.Header().Set(header.HttpHeaderContentLocation, srv.composeFriendlyUrlWithoutExtraPath(req.RequestURI, psi.UrlExtraPath))
	}
	for _, phpHdr := range phpScriptOutput.Headers {
		rw.Header().Set(phpHdr.Name, phpHdr.Value)
	}
	rw.Header().Set(header.HttpHeaderServer, srv.settings.ServerSoftware)

	// Status.
	if phpScriptOutput.StatusCode == 0 {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(int(phpScriptOutput.StatusCode))
	}

	// Body.
	_, err = rw.Write(phpScriptOutput.Body)
	if err != nil {
		log.Println(err)
	}
}

// composeFriendlyUrlWithoutExtraPath composes an adequate URL having no extra
// path instead of legacy CGI extra path.
// 'requestURI' is the original request URI with extra path.
func (srv *Server) composeFriendlyUrlWithoutExtraPath(requestURI, extraPath string) (friendlyUrl string) {
	return strings.TrimSuffix(requestURI, extraPath) + "?" + pm.QueryParamExtraPath + "=" + extraPath
}

// redirectToFriendlyUrlWithoutExtraPathIfNeeded redirects to a friendly URL without CGI extra path if needed.
// If a redirect has happened, 'true' is returned, otherwise – false.
func (srv *Server) redirectToFriendlyUrlWithoutExtraPathIfNeeded(rw http.ResponseWriter, req *http.Request, psi *pm.PhpScriptInfo) (isRedirectDone bool) {
	if len(psi.UrlExtraPath) == 0 {
		return false
	}

	if (srv.settings.PhpbbDoNotRedirectExtraPathInstallerStatus) && (psi.UrlExtraPath == ExtraPathInstallerStatus) {
		return false
	}

	srv.redirectToFriendlyUrlWithoutExtraPath(rw, req, psi)
	return true
}

func (srv *Server) redirectToFriendlyUrlWithoutExtraPath(rw http.ResponseWriter, req *http.Request, psi *pm.PhpScriptInfo) {
	rw.Header().Set(header.HttpHeaderLocation, srv.composeFriendlyUrlWithoutExtraPath(req.RequestURI, psi.UrlExtraPath))

	httpBody, err := io.ReadAll(req.Body)
	if err != nil {
		srv.respondWithInternalServerError(rw, err)
		return
	}

	if len(httpBody) > 0 {
		// 307 Redirect preserves original HTTP body.
		rw.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		// 302 Redirect does not preserve original HTTP body.
		rw.WriteHeader(http.StatusFound)
	}
}
