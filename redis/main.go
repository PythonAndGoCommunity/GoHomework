package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"redis/argscheck"
	"redis/client"
	"redis/serv"
	"redis/types"
)

const (
	protocolTCP = "tcp"
	separator   = " "
	buff        = 128
)

var (
	gPort   = ":9090"
	gMemory = "..." //expected disk or blank
	gIP     = "127.0.0.1"
)

func main() {
	argscheck.Start(os.Args, gPort, gMemory, gIP)
	ServFunc(gPort)
}

//ServFunc - containing Server and Client
func ServFunc(gPort string) {
	fmt.Printf("Port->%s\n", gPort)
	li, err := net.Listen(protocolTCP, gPort)
	if err != nil {
		fmt.Println("Error: ", err)
		log.Fatal(err)
	}
	defer li.Close()
	fmt.Println("Entered with :" + gIP + gPort)
	go client.Commands(protocolTCP, gPort, gIP, buff)
	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		ServConnCh := make(chan types.Server)
		go serv.ServConnHandler(ServConnCh, conn)
		go serv.ServCmndsHandler(ServConnCh, gMemory)
	}
}
