package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	_ "github.com/SiarheiKresik/go-kvdb/common"
)

func main() {

	const protocol = "tcp"

	addr := net.JoinHostPort(config.host, strconv.Itoa(config.port))
	// TODO check if host and port are valid

	log.Println("Connecting to server ", addr, "...")

	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	remoteAddr := conn.RemoteAddr()
	reader := bufio.NewReader(os.Stdin)

	for {
		// command line prompt
		fmt.Print("server: ", remoteAddr, " > ")

		// read in input from stdin
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("text error:", err)
		}

		// send to socket
		_, prnterr := fmt.Fprintf(conn, text)
		if prnterr != nil {
			log.Fatalln(prnterr)
		}

		// listen for reply
		message, readerr := bufio.NewReader(conn).ReadString('\n')
		if readerr != nil {
			log.Fatalln(readerr)
			break
		}
		fmt.Print(message)

	}
}
