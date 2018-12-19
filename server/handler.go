package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

func handle(commands chan command, conn net.Conn) {
	defer func() {
		conn.Close()
		log.Println("Connection closed")
	}()

	log.Println("Connection from", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fs := strings.Fields(ln)

		result := make(chan string)
		commands <- command{
			fields: fs,
			result: result,
		}

		io.WriteString(conn, <-result+"\n")
	}
}
