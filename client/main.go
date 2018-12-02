package main

import (
	"NonRelDB/util/file"
	"fmt"
	"bufio"
	"net"
	"flag"
	"NonRelDB/log"
	"NonRelDB/client/handler"
)

var host string
var port string
var dump bool
var restore bool
var location string

func init(){
	flag.StringVar(&host, "host", "127.0.0.1", "Defines host ip")
	flag.StringVar(&host, "h", "127.0.0.1", "Defines host ip")
	flag.StringVar(&port, "port", "9090", "Defines host port")
	flag.StringVar(&port, "p", "9090", "Defines host port")
	flag.BoolVar(&dump, "dump", false, "Requests db dump in json format from server")
	flag.BoolVar(&restore, "restore", false, "Restores received dump to file")
	flag.StringVar(&location, "location", "dump.json", "Defines location of dump")
	flag.Parse()
}

func main(){
	c, err := net.Dial("tcp", host + ":" + port)
	defer c.Close()

	if err != nil {
		log.Error.Panicln(err.Error())
	}

	if dump {
		// c.Write([]byte("dump\n"))
		fmt.Fprintf(c, "dump\n")
		dbDump, err := bufio.NewReader(c).ReadString('\n')

		if err != nil {
			log.Error.Panicln(err.Error())
		}

		fmt.Println(dbDump)

		if restore {
			file.CreateAndWriteString(location, dbDump)
		}
		return
	}

	handler.HandleConnection(c)
}