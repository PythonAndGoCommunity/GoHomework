package handler

import (
	"fmt"
	"strings"
	"bufio"
	"net"
	"NonRelDB/log"
	"NonRelDB/util/json"
	"NonRelDB/server/storage/inmemory"
	"NonRelDB/server/topic"
	"NonRelDB/util/regex"
)

// SendResponse sends response to specified connection.
func SendResponse(resp string, c net.Conn){
	fmt.Fprintf(c, resp + "\n")
	log.Info.Printf("Sent response to %s -> %s", c.RemoteAddr().String(), resp)
}

// HandleConnection handling communication with client.
func HandleConnection(c net.Conn){
	defer c.Close()

	netReader := bufio.NewReader(c)
	
	for {
		req, err := netReader.ReadString('\n')
		req = strings.TrimSuffix(req,"\n")

		if err != nil {
			log.Error.Println(err.Error())
			return
		} 

		log.Info.Printf("Received request from %s -> %s", c.RemoteAddr().String(), req)

		if regex.QueryReg.MatchString(req){
			resp:= HandleQuery(req)
			SendResponse(resp, c)
	
		} else if regex.ExitReg.MatchString(req){
			log.Info.Printf("%s disconnected from server", c.RemoteAddr().String())
			return
	
		} else if regex.DumpReg.MatchString(req){
			dbDump := string(json.PackMapToJSON((*inmemory.GetStorage().GetMap())))
			fmt.Fprintf(c, dbDump + "\n")
			log.Info.Printf("Sent db dump to %s", c.RemoteAddr().String())
			return
	
		} else if regex.TopicReg.MatchString(req){
			reqParts := strings.Split(req, " ")
	
			if len(reqParts) == 2 {
				if strings.ToLower(reqParts[0]) == "subscribe"{
					topic.Subscribe(reqParts[1], c)
	
				} else if strings.ToLower(reqParts[0]) == "unsubscribe"{
					topic.Unsubscribe(reqParts[1], c)
					return
	
				}
	
			} else if len(reqParts) >= 3{
				if strings.ToLower(reqParts[0]) == "publish"{
					msg := regex.DoubleQuoteReg.FindString(req)
					topic.Publish(reqParts[1], msg)
				}
	
			}
	
		} else {
			SendResponse("Bad request", c)
		}
	}
}