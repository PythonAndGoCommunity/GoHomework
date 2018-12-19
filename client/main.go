package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

const protocol = "tcp"

func main() {
	flag.Parse()

	fmt.Println("To exit type 'exit'")

	// validate host value
	host := net.ParseIP(config.host)
	if host == nil {
		log.Fatalln("invalid ip adress:", config.host)
	}

	// network address
	addr := net.JoinHostPort(config.host, strconv.Itoa(config.port))

	// create a client and connect
	client := NewClient()

	// connect to a server
	log.Printf("Connecting to server %s...\n", addr)

	remoteAddr, err := client.Connect(addr)
	if err != nil {
		log.Fatalf("Error connecting to %s", addr)
	}

	defer func() {
		log.Println("Closing connection...")
		client.Close()
		log.Println("Exit")

	}()

	log.Println("Connected")

	// create cli
	cli := NewCli()

	for {
		cli.Prompt(remoteAddr.String())

		// read user input
		text, err := cli.Read()
		if err != nil {
			log.Fatalln("text error:", err)
		}

		if text == "exit\n" {
			break
		}

		// send to socket
		err = client.Send(text)
		if err != nil {
			log.Print("error sending message:", err)
			continue
		}

		// listen for reply
		message, readerr := bufio.NewReader(client.conn).ReadString('\n')
		if readerr != nil {
			log.Fatalln(readerr)
			break
		}

		cli.Write(message)
	}
}
