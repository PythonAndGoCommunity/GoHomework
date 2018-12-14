package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	var portFlag string
	flag.StringVar(&portFlag, "port", "9090", "Server port. Defalt: 9090")
	flag.StringVar(&portFlag, "p", "9090", "Server port. Defalt: 9090")
	flag.Parse()
	listener, err := net.Listen("tcp", ":"+portFlag)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	log.Printf("Server is running on %s\n", portFlag)
	log.Println("Ready to accept connections")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	dataBase := make(map[string]string)
	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)
	for {
		command, err := clientReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		commands := strings.Split(command[:len(command)-1], " ")
		switch commands[0] {
		case "set":
			if len(commands) == 3 {
				dataBase[commands[1]] = commands[2]
				clientWriter.WriteString("Ok")
			} else {
				clientWriter.WriteString("Wrong command. Right: set [key] [value]")
			}
		case "get":
			if len(commands) == 2 {
				clientWriter.WriteString("\"" + dataBase[commands[1]] + "\"")
			} else {
				clientWriter.WriteString("Wrong command. Right: get [key]")
			}
		case "del":
			if len(commands) == 2 {
				delete(dataBase, commands[1])
				clientWriter.WriteString("Ok")
			} else {
				clientWriter.WriteString("Wrong command. Right: det [key]")
			}
		default:
			clientWriter.WriteString("unknown command: " + commands[0])
		}
		clientWriter.WriteString("\n")
		clientWriter.Flush()

	}
}
