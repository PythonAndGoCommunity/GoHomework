package main

import (
	"NonRelDB/client/handler"
	"NonRelDB/log"
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

var host string
var port string
var dump bool
var restore bool
var location string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Defines host ip")
	flag.StringVar(&host, "h", "127.0.0.1", "Defines host ip")
	flag.StringVar(&port, "port", "9090", "Defines host port")
	flag.StringVar(&port, "p", "9090", "Defines host port")
	flag.BoolVar(&dump, "dump", false, "Requests db dump in json format from server")
	flag.BoolVar(&restore, "restore", false, "Restores db from dumped file")
	flag.Parse()
}

// main entry point for client.
func main() {
	c, err := net.Dial("tcp", host+":"+port)
	defer c.Close()

	if err != nil {
		log.Error.Panicln(err.Error())
	}

	if dump {
		fmt.Fprintf(c, "dump\n")
		dbDump, err := bufio.NewReader(c).ReadString('\n')

		if err != nil {
			log.Error.Panicln(err.Error())
		}

		fmt.Println(dbDump)
		return
	}

	if restore {
		fmt.Fprintf(c, "restore\n")
		dbRestore, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			log.Error.Panicln(err.Error())
		}

		fmt.Fprintf(c, dbRestore)
		return
	}

	handler.HandleConnection(c)
}
