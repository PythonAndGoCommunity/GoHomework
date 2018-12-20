// server_main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
)

const DEFAULT_PORT = "9090"
const DEFAULT_HOST = "127.0.0.1"

var database = make(map[string]string)

func main() {
	var listenPort string
	var listenHost string
	flag.StringVar(&listenPort, "port", DEFAULT_PORT, "Port to listen on.")
	flag.StringVar(&listenPort, "p", DEFAULT_PORT, "Port to listen on.")
	flag.StringVar(&listenHost, "host", DEFAULT_HOST, "Host to listen on.")
	flag.StringVar(&listenHost, "h", DEFAULT_HOST, "Host to listen on.")
	flag.Parse()
	fmt.Print("Listening on port:" + listenPort + "\n")
	fmt.Print("Listening on host:" + listenHost + "\n")

	//fmt.Println(listenHost)
	//fmt.Println(listenPort)
	listener, err := net.Listen("tcp", listenHost+":"+listenPort)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	connReader := bufio.NewReader(conn)
	for {
		request, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		command := strings.Split(request[:len(request)-1], " ")
		fmt.Println(command[0])
		switch command[0] {
		case "SET":
			if len(command) == 3 {
				database[command[1]] = command[2]
				conn.Write([]byte("SET successful" + "\n"))
			} else {
				conn.Write([]byte("Error in command syntax. Syntax: set [key] [value]" + "\n"))
			}
		case "GET":
			if len(command) == 2 {
				conn.Write([]byte(database[command[1]] + "\n"))
			} else {
				conn.Write([]byte("Error in command syntax. Syntax: get [key]" + "\n"))
			}
		case "DEL":
			if len(command) == 2 {
				delete(database, command[1])
				conn.Write([]byte("DEL successful" + "\n"))
			} else {
				conn.Write([]byte("Error in command syntax. Syntax: del [key]" + "\n"))
			}
		case "KEYS":
			all_key := []string{}
			for key, _ := range database {
				all_key = append(all_key, key)
			}
			conn.Write([]byte(strings.Join(all_key, " ") + "\n"))

		default:
			conn.Write([]byte("Unsupported command: " + "\n"))
		}
	}

}
