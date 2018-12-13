/*
	Title: Goredis serverside
	Description: Goredis is a simple implementation of Redis on Golang
    by Gutyra A.
*/

package main

import (
	"GoHomework/server/server"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	defaultProtocol = "tcp"
	defaultAddress  = "localhost"
)

func main() {
	storage, port, verbose := server.GetCommandLineParams()
	err := server.LoadData(server.DataFile)
	if err != nil {
		os.Exit(1)
	}
	server.ConfigureLogging()

	ln, err := net.Listen(defaultProtocol, defaultAddress + ":" + port)
	if err != nil {
		fmt.Println("ERROR: Cannot create listener.")
	}

	setStorage(storage)
	setVerbose(verbose)

	go acceptConnections(ln)

	<-server.ExitChannel
	server.SaveData(server.DataFile)
	close(server.ExitChannel)
	fmt.Println("Server stopped. Data saved.")
	ln.Close()
}

func setStorage(storage string) {
	if strings.Compare(storage, "memory") == 0 {
		server.Storage = false
	} else {
		server.Storage = true
	}
}

func setVerbose(verbose bool) {
	if verbose {
		server.PrintlnAndLog("Server launched in verbose mode.")
		server.Verbose = true
	} else {
		fmt.Println("Server launched.")
		server.Verbose = false
	}
}

func acceptConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			server.PrintlnAndLog("ERROR: Cannot accept new connection.")
		}
		server.Channels[conn] = make([]string, 0)
		go server.HandleConnection(conn)
	}
}
