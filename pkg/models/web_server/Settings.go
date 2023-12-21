package ws

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	ae "github.com/vault-thirteen/auxie/errors"
)

const (
	SymbolDot = '.'
)

type Settings struct {
	DocumentRootPath   string   `json:"documentRootPath"`   // Path to the 'www' folder.
	FolderDefaultFiles []string `json:"folderDefaultFiles"` // List of default file names for a folder.

	ServerProtocol   string `json:"serverProtocol"`   // HTTP.
	GatewayInterface string `json:"gatewayInterface"` // CGI/1.1.
	ServerSoftware   string `json:"serverSoftware"`   // DemoGoServer/0.0.0.
	ServerName       string `json:"serverName"`       // Domain name: localhost.
	ServerHost       string `json:"serverHost"`       // IP address or domain name: localhost.
	ServerPort       string `json:"serverPort"`       // 8000.

	PhpServerNetwork  string   `json:"phpServerNetwork"`  // tcp.
	PhpServerHost     string   `json:"phpServerHost"`     // 127.0.0.1.
	PhpServerPort     string   `json:"phpServerPort"`     // 9000.
	PhpFileExtensions []string `json:"phpFileExtensions"` // "php", "phtml", ...

	// PHP is known to use an old-school variant of the 'Location' HTTP header.
	// FixRelativeRedirects, when enabled, fixed outdated URLs.
	// This feature is experimental and not safe.
	FixRelativeRedirects bool `json:"fixRelativeRedirects"`

	// IsCgiExtraPathEnabled flag enables support for CGI feature called "Extra
	// Path". This feature allows to make crazy-looking URLs which are
	// impossible to be parsed, something like the following:
	// http://some.machine/cgi-bin/display.pl/cgi/cgi_doc.txt
	IsCgiExtraPathEnabled bool `json:"isCgiExtraPathEnabled"`

	IsCachingEnabled           bool `json:"isCachingEnabled"`
	FileServerCacheSizeLimit   int  `json:"fileServerCacheSizeLimit"`
	FileServerCacheVolumeLimit int  `json:"fileServerCacheVolumeLimit"`
	FileServerCacheRecordTtl   uint `json:"fileServerCacheRecordTtl"`

	// This setting forbids redirecting phpBB requests using `/installer/status`
	// CGI extra path. Such requests are done directly.
	PhpbbDoNotRedirectExtraPathInstallerStatus bool `json:"phpbbDoNotRedirectExtraPathInstallerStatus"`
}

func NewSettings(settingsFilePath string) (set *Settings, err error) {
	var file *os.File
	file, err = os.Open(settingsFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := file.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	var fileText []byte
	fileText, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	set = &Settings{}
	err = json.Unmarshal(fileText, set)
	if err != nil {
		return nil, err
	}

	set.PhpFileExtensions = convertFileExtensionsFromNormalToGolang(set.PhpFileExtensions)

	return set, nil
}

func convertFileExtensionsFromNormalToGolang(exts []string) (golangExts []string) {
	golangExts = make([]string, 0, len(exts))

	for _, ext := range exts {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if len(ext) == 0 {
			continue
		}
		ext = prependDot(ext)

		golangExts = append(golangExts, ext)
	}

	return golangExts
}

func prependDot(sIn string) (sOut string) {
	symbols := []rune(sIn)

	if len(symbols) == 0 {
		return sOut
	}

	if symbols[0] != SymbolDot {
		return string(SymbolDot) + sIn
	}

	return sIn
}
