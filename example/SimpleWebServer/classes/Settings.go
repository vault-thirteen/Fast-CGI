package c

import (
	"encoding/json"
	"io"
	"os"

	"github.com/vault-thirteen/errorz"
)

type Settings struct {
	DocumentRootPath string `json:"documentRootPath"` // Path to the 'www' folder.
	ServerProtocol   string `json:"serverProtocol"`   // HTTP.
	GatewayInterface string `json:"gatewayInterface"` // CGI/1.1.
	ServerSoftware   string `json:"serverSoftware"`   // DemoGoServer/0.0.0.
	ServerName       string `json:"serverName"`       // Domain name: localhost.
	ServerHost       string `json:"serverHost"`       // IP address or domain name: localhost.
	ServerPort       string `json:"serverPort"`       // 8000.
	PhpServerNetwork string `json:"phpServerNetwork"` // tcp.
	PhpServerHost    string `json:"phpServerHost"`    // 127.0.0.1.
	PhpServerPort    string `json:"phpServerPort"`    // 9000.
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
			err = errorz.Combine(err, derr)
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

	return set, nil
}
