package handler

import (
	"strings"
	"bufio"
	"net"
	"NonRelDB/log"
)

func HandleConnection(c net.Conn){
	defer c.Close()
	netReader := bufio.NewReader(c)
	for {
		query, err := netReader.ReadString('\n')

		if err != nil {
			log.Error.Println(err.Error())
			return
		} 

		if query == "exit" {
			log.Info.Printf("%s disconnected from server",c.RemoteAddr().String())
			return
		}

		log.Info.Printf("Received request from %s -> %s", c.RemoteAddr().String(), query)

		resp := HandleQuery(strings.TrimSpace(query))
		log.Info.Printf("Sent response to %s -> %s", c.RemoteAddr().String(), resp)
		c.Write([]byte(resp + "\n"))
	}
}