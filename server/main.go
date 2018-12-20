package main

import (
	"flag"
	"log"
	"net"
	"strconv"
)

const protocol = "tcp"

type command struct {
	fields []string
	result chan string
}

func main() {

	flag.Parse()

	// listen on the given address
	addr := net.JoinHostPort("", strconv.Itoa(config.port))
	li, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatalln(err)
	}

	defer li.Close()

	// welcome messages
	log.Printf("Server is listening %s", addr)
	log.Println("Ready to accept connections")
	if config.verbose {
		log.Println("Server runs in verbose mode")
	}

	commands := make(chan command)

	// initialize storage
	go storage(commands, config.mode)

	// initial client id
	id := 1

	// wait and handle connections
	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(commands, conn, config.verbose, id)
		id++
	}
}
