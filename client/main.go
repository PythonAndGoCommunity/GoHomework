package main

import (
	"net"
	"flag"
	"NonRelDB/log"
	"NonRelDB/client/handler"
)

var host string
var port string

func init(){
	flag.StringVar(&host, "host", "127.0.0.1", "Defines host ip")
	flag.StringVar(&host, "h", "127.0.0.1", "Defines host ip")
	flag.StringVar(&port, "port", "9090", "Defines host port")
	flag.StringVar(&port, "p", "9090", "Defines host port")
	flag.Parse()
}

func main(){
	c, err := net.Dial("tcp", host + ":" + port)
	defer c.Close()

	if err != nil {
		log.Error.Panicln(err.Error())
	}

	handler.HandleConnection(c)
}