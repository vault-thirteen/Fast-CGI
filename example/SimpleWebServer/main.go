package main

// This example is a super simplified demonstration example of how an HTTP
// server in Golang can run PHP scripts. This package uses functions with
// unoptimized and super simple code for good readability. Please note that
// "production"-ready code would be different.

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	ws "github.com/vault-thirteen/Fast-CGI/pkg/models/web_server"
)

const ErrNotEnoughArguments = "not enough arguments"

func main() {
	settingsFilePath, err := getSettingsFilePath()
	mustBeNoError(err)

	var settings *ws.Settings
	settings, err = ws.NewSettings(settingsFilePath)
	mustBeNoError(err)

	var srv *ws.Server
	srv, err = ws.NewServer(settings)
	mustBeNoError(err)

	srv.Run()
	waitForQuitSignalFromOS(srv)
}

func getSettingsFilePath() (sfp string, err error) {
	if len(os.Args) < 2 {
		return sfp, errors.New(ErrNotEnoughArguments)
	}

	return os.Args[1], nil
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func waitForQuitSignalFromOS(srv *ws.Server) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	for sig := range osSignals {
		switch sig {
		case syscall.SIGINT,
			syscall.SIGTERM:
			log.Println("quit signal from OS has been received: ", sig)
			mustBeNoError(srv.Stop())
			close(osSignals)
		}
	}
}
