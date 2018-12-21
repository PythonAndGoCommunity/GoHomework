package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

func handle(commands chan command, conn net.Conn, verbose bool, id int) {
	defer func() {
		conn.Close()
		log.Printf("Connection closed, client: id %d", id)
	}()

	log.Printf("Connection from %s, client id: %d", conn.RemoteAddr(), id)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fs := strings.Fields(ln)

		if verbose {
			log.Printf("Request: [%s], client id: %d", ln, id)
		}

		result := make(chan string)
		commands <- command{
			fields: fs,
			result: result,
		}

		io.WriteString(conn, <-result+"\n")
	}
}
