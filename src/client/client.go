package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

const DEFAULT_PORT = "9090"
const DEFAULT_HOST = "127.0.0.1"

func main() {
	var connectPort string
	var connectHost string
	flag.StringVar(&connectPort, "port", DEFAULT_PORT, "Port to listen on.")
	flag.StringVar(&connectPort, "p", DEFAULT_PORT, "Port to listen on.")
	flag.StringVar(&connectHost, "host", DEFAULT_HOST, "Host to listen on.")
	flag.StringVar(&connectHost, "h", DEFAULT_HOST, "Host to listen on.")
	flag.Parse()
	fmt.Print("Sending commands on port:" + connectPort + "\n")
	fmt.Print("To host:" + connectHost + "\n")

	conn, conn_err := net.Dial("tcp", connectHost+":"+connectPort)
	if conn_err != nil {
		fmt.Println(conn_err)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte(text))
		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Response: " + response)
	}
}
