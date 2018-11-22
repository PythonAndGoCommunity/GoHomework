package handler

import (
	"net"
	"NonRelDB/server/log"
)

func HandleListener(l net.Listener){
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			log.Warning.Printf("Failed connection from %s",c.RemoteAddr().String())
			c.Close()
		}
		log.Info.Printf("%s was connected to server",c.RemoteAddr().String())
		go HandleConnection(c)
	}	
}