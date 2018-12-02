package handler

import (
	"fmt"
	"strings"
	"bufio"
	"net"
	"NonRelDB/log"
	"NonRelDB/util/json"
	"NonRelDB/server/storage/inmemory"
)

func HandleConnection(c net.Conn){
	defer c.Close()

	netReader := bufio.NewReader(c)
	
	for {

		query, err := netReader.ReadString('\n')
		query = strings.TrimSuffix(query,"\n")

		fmt.Println("[" + query + "]")

		if err != nil {
			log.Error.Println(err.Error())
			return
		} 

		if query == "exit" {
			log.Info.Printf("%s disconnected from server",c.RemoteAddr().String())
			return
		} else if query == "dump" {
			log.Info.Printf("Sent db dump to %s", c.RemoteAddr().String())
			dbDump := string(json.PackMapToJSON((*inmemory.GetStorage().GetMap())))
			fmt.Fprintf(c, dbDump + "\n")
			return 
		}

		log.Info.Printf("Received request from %s -> %s", c.RemoteAddr().String(), query)

		resp := HandleQuery(strings.TrimSpace(query))
		log.Info.Printf("Sent response to %s -> %s", c.RemoteAddr().String(), resp)
		fmt.Fprintf(c, resp + "\n")
	}
}