package handler

import (
	"strings"
	"net"
	"bufio"
	"os"
	"fmt"
	"NonRelDB/log"
)

func HandleConnection(c net.Conn){
	consoleReader := bufio.NewReader(os.Stdin)
	netReader := bufio.NewReader(c)
	for {
		fmt.Print("NonRelDB> ")
		command, err := consoleReader.ReadString('\n')
		
		if err != nil {
			log.Error.Panicln(err.Error())
		}
	
		if command == "exit\n"{
			fmt.Println("Good bye")
			fmt.Fprintf(c, command)
			return
		} else if strings.Contains(command, "subscribe"){
			fmt.Fprintf(c, command)
			HandleTopic(c, *netReader, strings.Split(command, " ")[1])
		} else if strings.Contains(command, "publish"){
			fmt.Fprintf(c, command)
			continue
		}
		
		fmt.Fprintf(c, command)

		resp, err := netReader.ReadString('\n')

		if err != nil {
			log.Error.Panicln(err.Error())
		}

		fmt.Println(resp)
	}
}