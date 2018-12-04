package handler

import (
	"regexp"
	"fmt"
	"strings"
	"bufio"
	"net"
	"NonRelDB/log"
	"NonRelDB/util/json"
	"NonRelDB/server/storage/inmemory"
	"NonRelDB/server/topic"
)

var queryReg *regexp.Regexp
var exitReg *regexp.Regexp
var dumpReg *regexp.Regexp
var subscribeReg *regexp.Regexp
var unsubscribeReg *regexp.Regexp
var topicNameReg *regexp.Regexp
var publishReg *regexp.Regexp
var publishTopicNameReg *regexp.Regexp
var msgReg *regexp.Regexp

func init(){
	queryReg = regexp.MustCompile("^(get\\s(.*))|(set\\s(.*)\\s\"(.*)\")|(del\\s(.*))|(keys\\s\"(.*?)\")$")
	exitReg = regexp.MustCompile("^exit$")
	dumpReg = regexp.MustCompile("^dump$")
	subscribeReg = regexp.MustCompile("^subscribe\\s(.*)$")
	unsubscribeReg = regexp.MustCompile("^unsubscrube\\s(.*)$")
	topicNameReg = regexp.MustCompile("\\s(.*)$")
	publishReg = regexp.MustCompile("^publish\\s(.*)\\s\"(.*)\"$")
	publishTopicNameReg = regexp.MustCompile("\\s(.*)\\s")
	msgReg = regexp.MustCompile("\"(.*)\"$")
}

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

		if queryReg.MatchString(req){
			resp := HandleQuery(req)
			fmt.Fprintf(c, resp + "\n")
			log.Info.Printf("Sent response to %s -> %s", c.RemoteAddr().String(), resp)

		} else if exitReg.MatchString(req){
			log.Info.Printf("%s disconnected from server", c.RemoteAddr().String())

		} else if dumpReg.MatchString(req){
			dbDump := string(json.PackMapToJSON((*inmemory.GetStorage().GetMap())))
			fmt.Fprintf(c, dbDump + "\n")
			log.Info.Printf("Sent db dump to %s", c.RemoteAddr().String())
		
		} else if subscribeReg.MatchString(req){
			topicName := strings.Trim(topicNameReg.FindString(req), " ")
			topic.Subscribe(topicName, c)
			continue
		
		} else if unsubscribeReg.MatchString(req){
			topicName := strings.Trim(topicNameReg.FindString(req), " ")
			topic.Unsubscribe(topicName, c)
			return
		
		} else if publishReg.MatchString(req){
			topicName := strings.Trim(publishTopicNameReg.FindString(req), " ")
			msg := strings.Trim(msgReg.FindString(req), "\"")
			topic.Publish(topicName, msg)
			continue
			
		}
	}
}