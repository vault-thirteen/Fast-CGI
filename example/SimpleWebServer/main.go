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

	"github.com/vault-thirteen/Fast-CGI/example/SimpleWebServer/classes"
)

const ErrNotEnoughArguments = "not enough arguments"

func main() {
	settingsFilePath, err := getSettingsFilePath()
	mustBeNoError(err)

	var settings *c.Settings
	settings, err = c.NewSettings(settingsFilePath)
	mustBeNoError(err)

	var srv *c.Server
	srv, err = c.NewServer(settings)
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

func waitForQuitSignalFromOS(srv *c.Server) {
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
