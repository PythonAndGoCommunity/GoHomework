package handler

import (
	"strings"
	"net"
	"bufio"
	"os"
	"fmt"
	"NonRelDB/log"
	"NonRelDB/util/regex"
)

func SendRequest(req string, c net.Conn){
	fmt.Fprintf(c, req)
}

func HandleConnection(c net.Conn){
	consoleReader := bufio.NewReader(os.Stdin)
	netReader := bufio.NewReader(c)
	for {
		fmt.Print("NonRelDB> ")
		req, err := consoleReader.ReadString('\n')
		
		if err != nil {
			log.Error.Panicln(err.Error())
		}

		if regex.QueryReg.MatchString(req){
			SendRequest(req, c)
			resp, err := netReader.ReadString('\n')

			if err != nil {
				log.Error.Panicln(err.Error())
			}

			fmt.Println(resp)

		} else if regex.ExitReg.MatchString(req){
			fmt.Println("Good bye")
			SendRequest(req, c)
			return

		} else {
			reqParts := strings.Split(req, " ")[:2]

			switch reqCtx := strings.ToLower(reqParts[0]); reqCtx{
				case "subscribe":{
					SendRequest(req, c)
					HandleTopic(c, *netReader, reqParts[1])
				}
				case "publish":{
					SendRequest(req, c)
				}
		}
	}
}
}