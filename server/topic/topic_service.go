package topic

import (
	"NonRelDB/log"
	"fmt"
	"net"
	"NonRelDB/util/collection"
)

var topics map[string][]net.Conn

func init(){
	topics = make(map[string][]net.Conn)
}

func Subscribe(name string, c net.Conn){
	if topics[name] == nil {
		topics[name] = make([]net.Conn, 10)
	} 
	topics[name] = append(topics[name], c)
	log.Info.Printf("%s just subscribed %s", c.RemoteAddr().String(), name)
}

func Unsubscribe(name string, c net.Conn) {
	if topics[name] != nil || len(topics) != 0{
		index := collection.ConnIndex(topics[name], c)
		if index != -1 {
			topics[name][index] = nil
			log.Info.Printf("%s just unsubscribed %s", c.RemoteAddr().String(), name)
		}
	}
}

func Publish(name string, msg string){
	if topics[name] != nil || len(topics) != 0 {
		log.Info.Printf("%s just published in %s", msg , name)
		for _, listener := range topics[name] {
			str := fmt.Sprintf("[%s]: %s", name, msg)
			if (listener != nil){
				fmt.Fprintf(listener, str + "\n")
			}
		}
	}
}