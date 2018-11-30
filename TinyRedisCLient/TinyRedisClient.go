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
	port, host = cmdFlagParse()
	address := host + ":" + port

	conn, err := net.DialTimeout("tcp", address, 5*time.Second) // connection attempt
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to %s\n", address)
	defer fmt.Println("Connection closed")

	for {
		fmt.Print(">>> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n') // got client's command
		fmt.Fprintf(conn, text+"\n")       // which is being sent via port
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
	}
}

func contains(strs []string, str string) bool {
	if strs == nil {
		return true
	}
	for _, val := range strs {
		if str == val {
			return true
		}
	}
	return false
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

	portShortResult := flag.String(portHandler.options[0], portHandler.defaultVal, portHandler.warning)
	portLongResult := flag.String(portHandler.options[1], portHandler.defaultVal, portHandler.warning)

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
