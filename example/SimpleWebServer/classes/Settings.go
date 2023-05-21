package c

import (
	"os"
	"path/filepath"
)

const (
	PhpServerNetwork   = "tcp"
	PhpServerAddress   = "127.0.0.1:9000"
	WebServerName      = "localhost" // Domain name.
	WebServerHost      = "localhost" // IP address or domain name.
	WebServerPort      = "8000"
	ServerProtocol     = "HTTP"
	GatewayInterface   = "CGI/1.1"
	ServerSoftwareName = "DemoGoServer/0.0.0"
)

type Settings struct {
	documentRootPath string
	serverProtocol   string
	gatewayInterface string
	serverSoftware   string
	serverName       string
	serverHost       string
	serverPort       string
}

func NewSettings() (set *Settings, err error) {
	set = &Settings{
		//documentRootPath: "",
		serverProtocol:   ServerProtocol,
		gatewayInterface: GatewayInterface,
		serverSoftware:   ServerSoftwareName,
		serverName:       WebServerName,
		serverHost:       WebServerHost,
		serverPort:       WebServerPort,
	}

	set.documentRootPath, err = getDocumentRootPath()
	if err != nil {
		return nil, err
	}

	return set, nil
}

// getDocumentRootPath returns the path to the 'www' folder.
// In this demonstration example we assume that server's executable file is
// located in the 'www' folder for simplicity.
func getDocumentRootPath() (drp string, err error) {
	var exePath string
	exePath, err = os.Executable()
	if err != nil {
		return drp, err
	}

	drp = filepath.Dir(exePath)

	return drp, nil
}
