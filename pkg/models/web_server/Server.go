package ws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	cl "github.com/vault-thirteen/Fast-CGI/pkg/Client"
	sr "github.com/vault-thirteen/Fast-CGI/pkg/models/ScriptRunner"
	pm "github.com/vault-thirteen/Fast-CGI/pkg/models/php"
	sfs "github.com/vault-thirteen/Simple-File-Server"
	mime "github.com/vault-thirteen/auxie/MIME"
	"github.com/vault-thirteen/auxie/file"
	"github.com/vault-thirteen/auxie/header"
)

const (
	ErrRemoteAddrParts       = "remote address error"
	ErrContentLengthMismatch = "content length mismatch"
)

const (
	MimeTypeDefault = "application/octet-stream"
)

type Server struct {
	settings     *Settings
	httpServer   *http.Server
	cgiClient    *cl.Client
	scriptRunner *sr.ScriptRunner
	fileServer   *sfs.SimpleFileServer

	// MIME types.
	// Key: file extension (dot-prefixed); Value: MIME type.
	mimeTypes map[string]string
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

	srv.scriptRunner = sr.New()

	srv.fileServer, err = sfs.NewSimpleFileServer(
		srv.settings.DocumentRootPath,
		srv.settings.FolderDefaultFiles,
		srv.settings.IsCachingEnabled,
		srv.settings.FileServerCacheSizeLimit,
		srv.settings.FileServerCacheVolumeLimit,
		srv.settings.FileServerCacheRecordTtl,
	)
	if err != nil {
		return nil, err
	}

	srv.mimeTypes = srv.getMimeTypes()

	return srv, nil
}

func (srv *Server) getMimeTypes() (mimeTypes map[string]string) {
	return map[string]string{
		".css":   mime.TypeTextCss,
		".htm":   mime.TypeTextHtml,
		".html":  mime.TypeTextHtml,
		".js":    mime.TypeApplicationJavascript,
		".json":  mime.TypeApplicationJson,
		".pdf":   mime.TypeApplicationPdf,
		".txt":   mime.TypeTextPlain,
		".xhtml": mime.TypeApplicationXhtmlXml,
		".xml":   mime.TypeApplicationXml,

		// Images.
		".avif": mime.TypeImageAvif,
		".bmp":  mime.TypeImageBmp,
		".gif":  mime.TypeImageGif,
		".ico":  mime.TypeImageVndMicrosoftIcon,
		".jpeg": mime.TypeImageJpeg,
		".jpg":  mime.TypeImageJpeg,
		".png":  mime.TypeImagePng,
		".svg":  mime.TypeImageSvgXml,
		".webp": mime.TypeImageWebp,

		// Audio.
		".aac":  mime.TypeAudioAac,
		".mp3":  mime.TypeAudioMpeg,
		".opus": mime.TypeAudioOpus,

		// Video.
		".mp4":  mime.TypeVideoMp4,
		".mpeg": mime.TypeVideoMpeg,

		// Archive.
		".rar": mime.TypeApplicationVndRar,
		".zip": mime.TypeApplicationZip,
	}
}

func (srv *Server) router(rw http.ResponseWriter, req *http.Request) {
	var psi = &pm.PhpScriptInfo{
		OriginalUrlPath: req.URL.Path,
		UrlRelPath:      req.URL.Path,
	}

	psi.QueryParamExtraPath = req.URL.Query().Get(pm.QueryParamExtraPath)

	var isFolderCheckDisabled = false
	if srv.settings.IsCgiExtraPathEnabled {
		path, extraPath, err := srv.fileServer.FindExtraPath(psi.UrlRelPath)
		if err == nil {
			psi.UrlRelPath = path
			psi.UrlExtraPath = extraPath
		}

		// The installer of phpBB has a bug:
		// The installation process is started by a redirect to the page:
		// "http://localhost:8000/phpBB3/install/app.php", no trailing slash.
		// The link to itself from the page leads to another URL:
		// "http://localhost:8000/phpBB3/install/app.php/", with trailing slash.
		// What can I say ? You guys are very lame. Because of you I had to
		// write all this stupid shit. This video very vividly shows how phpBB
		// and CGI work together in a modern web server. Enjoy the movie:
		// https://www.youtube.com/watch?v=JHA6OxF3k0g
		if psi.UrlExtraPath == ExtraPathSingleSlash {
			psi.UrlExtraPath = ""
			isFolderCheckDisabled = true
		}
	}

	if srv.redirectToFriendlyUrlWithoutExtraPathIfNeeded(rw, req, psi) {
		return
	}

	//log.Println(fmt.Sprintf("path=[%v], extra-path=[%v].", psi.UrlRelPath, psi.UrlExtraPath)) //DEBUG.
	psi.FilePath = strings.ReplaceAll(psi.UrlRelPath, `/`, string(os.PathSeparator))

	// If a folder is requested, replace it with a default file.
	if (sfs.IsPathFolder(psi.OriginalUrlPath)) && (!isFolderCheckDisabled) {
		fileName, err := srv.fileServer.GetFolderDefaultFilename(psi.UrlRelPath)
		if err != nil {
			switch err.Error() {
			case sfs.Err_PathIsNotValid:
				srv.respondWithNotAllowed(rw)
				return

			default:
				srv.respondWithInternalServerError(rw, err)
				return
			}
		}
		if len(fileName) == 0 {
			srv.respondWithNotFound(rw)
			return
		}

		psi.FilePath = filepath.Join(psi.FilePath, fileName)
	}

	psi.FileName = filepath.Base(psi.FilePath)
	psi.FileExt = filepath.Ext(psi.FileName)

	if srv.isExtOfPhpScript(psi.FileExt) {
		psi.FileAbsPath = filepath.Join(srv.settings.DocumentRootPath, psi.FilePath)
		if len(psi.UrlExtraPath) > 0 {
			psi.FileAbsExtraPath = filepath.Join(srv.settings.DocumentRootPath, psi.UrlExtraPath)
		}
		srv.runPhpScript(rw, req, psi)
	} else {
		srv.serveOrdinaryFile(rw, psi.FilePath, psi.FileExt)
	}
}

func (srv *Server) getMimeTypeByExt(ext string) (mimeType string) {
	var ok bool
	mimeType, ok = srv.mimeTypes[ext]
	if !ok {
		log.Println(fmt.Sprintf("unknown file extension: %v.", ext))
		return MimeTypeDefault
	}

	return mimeType
}

func (srv *Server) serveOrdinaryFile(rw http.ResponseWriter, relFilePath string, fileExt string) {
	fileExists, err := srv.fileServer.FileExists(relFilePath)
	if err != nil {
		if err.Error() == file.ErrObjectIsNotFile {
			srv.respondWithNotAllowed(rw)
			return
		} else {
			srv.respondWithInternalServerError(rw, err)
			return
		}
	}
	if !fileExists {
		srv.respondWithNotFound(rw)
		return
	}

	var fileContents []byte
	fileContents, err = srv.fileServer.GetFile(relFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			srv.respondWithNotFound(rw)
			return
		} else {
			srv.respondWithInternalServerError(rw, err)
			return
		}
	}

	srv.respondWithData(rw, fileContents, fileExt)
	return
}

func (srv *Server) respondWithInternalServerError(rw http.ResponseWriter, err error) {
	log.Println(err)
	rw.WriteHeader(http.StatusInternalServerError)
}

func (srv *Server) respondWithNotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
}

func (srv *Server) respondWithNotAllowed(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusForbidden)
}

func (srv *Server) respondWithData(rw http.ResponseWriter, data []byte, fileExt string) {
	rw.Header().Set(header.HttpHeaderContentType, srv.getMimeTypeByExt(fileExt))
	rw.Header().Set(header.HttpHeaderServer, srv.settings.ServerSoftware)

	var err error
	_, err = rw.Write(data)
	if err != nil {
		log.Println(err)
	}
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
