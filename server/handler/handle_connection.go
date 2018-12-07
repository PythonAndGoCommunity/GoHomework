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

// HandleRequest handling request from client.
func HandleRequest(req string, c net.Conn) rune{
	if regex.QueryReg.MatchString(req){
		resp:= HandleQuery(req)
		SendResponse(resp, c)
		return 'c'

	} else if regex.ExitReg.MatchString(req){
		log.Info.Printf("%s disconnected from server", c.RemoteAddr().String())
		return 'r'

	} else if regex.DumpReg.MatchString(req){
		dbDump := string(json.PackMapToJSON((*inmemory.GetStorage().GetMap())))
		fmt.Fprintf(c, dbDump + "\n")
		log.Info.Printf("Sent db dump to %s", c.RemoteAddr().String())
		return 'r'

	} else if regex.TopicReg.MatchString(req){
		reqParts := strings.Split(req, " ")

		if len(reqParts) == 2 {
			if strings.ToLower(reqParts[0]) == "subscribe"{
				topic.Subscribe(reqParts[1], c)
				return 'c'

			} else if strings.ToLower(reqParts[0]) == "unsubscribe"{
				topic.Unsubscribe(reqParts[1], c)
				return 'r'

			}

		} else if len(reqParts) >= 3{
			if strings.ToLower(reqParts[0]) == "publish"{
				msg := regex.DoubleQuoteReg.FindString(req)
				topic.Publish(reqParts[1], msg)
				return 'c'
			}

		}

	} else {
		SendResponse("Bad request", c)
	}

	return 'c'
}

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

		switch handleCtx := HandleRequest(req, c); handleCtx{
			case 'r': return
			case 'c': continue
		}
	}
}