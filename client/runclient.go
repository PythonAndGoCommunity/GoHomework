/*
	Title: Goredis clientside
	Description: Goredis is a simple implementation of Redis on Golang
    by Gutyra A.
*/

package main

import (
	"GoHomework/client/client"
	"fmt"
	"net"
	"os"
)

const protocol = "tcp"

func main() {
	host, port, dump, filename := client.GetCommandLineParams()
	exitMsgChan := make(chan string)

	conn, err := net.Dial(protocol, host + ":" + port)
	if err != nil {
		fmt.Println("ERROR: cannot connect to " + host + ":" + port)
		os.Exit(1)
	}

	go client.HandleServerResponds(conn, host + ":" + port, exitMsgChan)
	go client.HandleUserRequests(conn, exitMsgChan)

	if filename != "" {
		client.SendRequest(conn, "RESTORE \"" + filename + "\"\n")
	}

	if dump {
		client.SendRequest(conn, "DUMP\n")
	}

	fmt.Println(<-exitMsgChan)
	close(exitMsgChan)
}
