package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

//Reading commands from terminal
func startup() (addr string, err error) {
	//default values
    addr = ""
	host := "127.0.0.1"
	port := ":9090"
	err = nil

	args := os.Args
	for i := 1; i < len(args); i += 2 {
		switch args[i] {
		// -p, --port	: set the port for listening on
		case "-p", "--port":
			if args[i+1][0] == ':' {
				port = args[i+1]
			} else {
				port = ":" + args[i+1]
			}
		// -m, --mode	: enable mirroring data to the drive
		case "-h", "--host":
			host = args[i+1]
		default:
			err = errors.New("ERROR: Unknown command\n")            
			return
		}
	}
    addr = host + port
	return
}


func main() {
	addr, cmdErr := startup()
	if cmdErr == nil {
		conn, netErr := net.Dial("tcp", addr)
		if netErr != nil {
			log.Fatalln(netErr)
		} else {
			defer conn.Close()
			fmt.Println("Connected to", addr)
			for {
				// read in input from stdin
				reader := bufio.NewReader(os.Stdin)
				command, _ := reader.ReadString('\n')
				fmt.Fprintf(conn, command)
				message, _ := bufio.NewReader(conn).ReadString('\n')
				fmt.Print(message)
			}
		}

	} else {
		//Some CMD error
		fmt.Println(cmdErr)
	}
}
