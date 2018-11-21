package handler

import (
	"net"
)

func HandleConnection(c net.Conn){
	defer c.Close()
	for {
		req := make([]byte,1024)
		lenReq, err := c.Read(req)
		if err != nil {
			panic(err)
		}
		query := string(req[:lenReq]) 
		resp_str := HandleQuery(query)
		resp := []byte(resp_str)
		c.Write(resp)
	}
}