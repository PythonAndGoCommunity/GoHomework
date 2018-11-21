package handler

import (
	"net"
)

func HandleListener(l net.Listener){
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			// log.Print(err.Error())
			// errorLogger.Println(err.Error())
			c.Close()
		}
		// infoLogger.Println("Connection successfully accepted")
		go HandleConnection(c)
	}	
}