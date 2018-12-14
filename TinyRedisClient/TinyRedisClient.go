package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var port, host string

type cmdFlagHandler struct { // command line flag handler
	options       []string
	defaultVal    string
	warning       string
	possibleRange []string
}

func main() {
	address := getAddress() // address retrieving

	conn := connectionAttempt(address) // connection handler
	fmt.Printf("Connected to %s\n", address)
	defer fmt.Println("Connection closed")

	for {
		fmt.Print(">>> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n') // got client's command
		fmt.Fprintf(conn, text+"\n")       // which is being sent via port

		if len(text) > 4 { // closing connection with server
			if text[0:4] == "EXIT" {
				return
			}
		}

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
	}
}

func getAddress() string { // returns complete server address
	port, host = cmdFlagParse()
	address := host + ":" + port
	return address
}

func cmdFlagParse() (string, string) { // cmd flag parser
	portHandler := cmdFlagHandler{ // default port and host data
		options:    []string{"p", "port"},
		defaultVal: "9090",
		warning:    "specify port to use",
	}
	hostHandler := cmdFlagHandler{
		options:    []string{"h", "host"},
		defaultVal: "127.0.0.1",
		warning:    "specify port to use",
	}
	// long and short forms are for long and short flags, e.g. -p and --port
	portShortResult := flag.String(portHandler.options[0],
		portHandler.defaultVal,
		portHandler.warning)

	portLongResult := flag.String(portHandler.options[1],
		portHandler.defaultVal,
		portHandler.warning)

	hostShortResult := flag.String(hostHandler.options[0], hostHandler.defaultVal, hostHandler.warning)
	hostLongResult := flag.String(hostHandler.options[1], hostHandler.defaultVal, hostHandler.warning)

	flag.Parse()

	var portResult, hostResult string
	if *portLongResult != portHandler.defaultVal {
		portResult = *portLongResult
	} else {
		portResult = *portShortResult
	}

	if *hostLongResult != hostHandler.defaultVal {
		hostResult = *hostLongResult
	} else {
		hostResult = *hostShortResult
	}
	return portResult, hostResult
}

func connectionAttempt(addr string) net.Conn { // five seconds for dialing server
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		panic(err)
	}
	return conn
}
