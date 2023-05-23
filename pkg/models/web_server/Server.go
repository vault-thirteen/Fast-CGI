package ws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vault-thirteen/Fast-CGI/pkg/Client"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/data"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/http"
	"github.com/vault-thirteen/Fast-CGI/pkg/models/php"
	"github.com/vault-thirteen/header"
)

const (
	ErrRemoteAddrParts       = "remote address error"
	ErrContentLengthMismatch = "content length mismatch"
)

const (
	HostPortDelimiter  = ":"
	GolangNetNetworkIP = "ip" // These constants should be exported by Golang ! Source: net/iprawsock.go.
)

type Server struct {
	settings   *Settings
	httpServer *http.Server
	cgiClient  *cl.Client
}

func NewServer(settings *Settings) (srv *Server, err error) {
	srv = &Server{
		settings: settings,
	}

	srv.httpServer = &http.Server{
		Addr:    net.JoinHostPort(srv.settings.ServerHost, srv.settings.ServerPort),
		Handler: http.Handler(http.HandlerFunc(srv.router)),
	}

	phpServerAddress := net.JoinHostPort(srv.settings.PhpServerHost, srv.settings.PhpServerPort)
	srv.cgiClient, err = cl.New(srv.settings.PhpServerNetwork, phpServerAddress)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Server) router(rw http.ResponseWriter, req *http.Request) {
	scriptFilePath := filepath.Join(srv.settings.DocumentRootPath, req.URL.Path)
	scriptFileName := filepath.Base(scriptFilePath)
	scriptFileExt := filepath.Ext(scriptFileName)

	if srv.isExtOfPhpScript(scriptFileExt) {
		srv.runPhpScript(rw, req, scriptFilePath, scriptFileName)
	} else {
		srv.respondWithNotAllowed(rw)
	}
}

func (srv *Server) isExtOfPhpScript(ext string) bool {
	for _, phpExt := range srv.settings.PhpFileExtensions {
		if ext == phpExt {
			return true
		}
	}

	return false
}

func (srv *Server) runPhpScript(rw http.ResponseWriter, req *http.Request, scriptFilePath string, scriptFileName string) {
	var requestId uint16 = 1
	var parameters []*nvpair.NameValuePair
	var stdin []byte
	var err error
	stdin, parameters, err = srv.prepareInputDataToRunPhpScript(req, scriptFilePath, scriptFileName)
	if err != nil {
		srv.respondWithInternalServerError(rw, err)
		return
	}

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

	// Headers.
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

func (srv *Server) prepareInputDataToRunPhpScript(req *http.Request, scriptFilePath string, scriptFileName string) (stdin []byte, parameters []*nvpair.NameValuePair, err error) {
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

	parameters = []*nvpair.NameValuePair{
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_AuthType, authScheme),                                      // 4.1.1.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ContentLength, strconv.Itoa(len(stdin))),                   // 4.1.2.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ContentType, req.Header.Get(header.HttpHeaderContentType)), // 4.1.3.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_DocumentRoot, srv.settings.DocumentRootPath),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_DocumentUri, req.URL.Path),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_GatewayInterface, srv.settings.GatewayInterface), // 4.1.4.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_PathInfo, ""), // 4.1.5.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_PathTranslated, ""), // 4.1.6.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_QueryString, req.URL.RawQuery), // 4.1.7.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RedirectStatus, strconv.Itoa(http.StatusOK)),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteAddr, remoteAddrParts[0]), // Host. 4.1.8.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteHost, ""),                 // FQDN. 4.1.9.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteIdent, ""),                // 4.1.10.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemotePort, remoteAddrParts[1]), // Port.
		//nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RemoteUser, ""),                 // 4.1.11.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RequestMethod, req.Method),     // 4.1.12.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RequestScheme, req.URL.Scheme), // Apache HTTP Server Header.
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_RequestUri, req.RequestURI),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ScriptFilename, scriptFilePath),
		nvpair.NewNameValuePairWithTextValueU(dm.Parameter_ScriptName, scriptFileName), // 4.1.13.
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

func (srv *Server) respondWithInternalServerError(rw http.ResponseWriter, err error) {
	log.Println(err)
	rw.WriteHeader(http.StatusInternalServerError)
}

func (srv *Server) respondWithNotAllowed(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusForbidden)
}

func (srv *Server) Run() {
	go srv.run()
}

func (srv *Server) run() {
	var err = srv.httpServer.ListenAndServe()
	if (err != nil) && (err != http.ErrServerClosed) {
		log.Println(err)
		mustBeNoError(srv.Stop())
	}
}

func (srv *Server) Stop() (err error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	fmt.Print("HTTP Server Shutdown ... ")
	err = srv.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Done")

	fmt.Print("FastCGI Client Shutdown ... ")
	err = srv.cgiClient.Close()
	if err != nil {
		return err
	}
	fmt.Println("Done")

	return nil
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
