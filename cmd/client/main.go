package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/ITandElectronics/GoHomework/protocol"
)

func main() {
	var (
		host string
		port int
	)

	flag.IntVar(&port, "port", 9090, "Remote server port")
	flag.IntVar(&port, "p", 9090, "Remote server port")
	flag.StringVar(&host, "host", "127.0.0.1", "Remote server address")
	flag.StringVar(&host, "h", "127.0.0.1", "Remote server address")
	flag.Usage = func() {
		usage := `Usage of %s:
  --host, -h string
	  Remote server address (default "127.0.0.1")
  --port, -p int
	  Remote server port (default 9090)
`
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.Parse()

	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("couldn't connect to server: %v\n", err)
	}

	input := bufio.NewReader(os.Stdin)
	response := bufio.NewReader(conn)
	for {
		fmt.Print("Enter command: ")
		line, _, err := input.ReadLine()
		if err != nil {
			fmt.Printf("[error] couldn't read line from stdin: %v\n", err)
			continue
		}
		msg, err := protocol.DecodeMessage(line)
		if err != nil {
			fmt.Printf("[error] invalid message format: %v\n", err)
			continue
		}
		if err := protocol.ValidateMessage(msg); err != nil {
			fmt.Printf("validation failed: %v\n", err)
			continue
		}
		if _, err = conn.Write(append(line, '\n')); err != nil {
			fmt.Printf("[error] coudn't write to connection: %v\n", err)
			continue
		}

		resp, _, err := response.ReadLine()
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection is closed")
				return
			}
			fmt.Printf("[error] couldn't read response: %v\n", err)
			continue
		}
		fmt.Printf("[server] %s\n", string(resp))
	}

}
