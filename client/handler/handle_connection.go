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
	
		if command == "exit"{
			fmt.Println("Good bye")
			fmt.Print(c, command)
			return
		}

		fmt.Fprintf(c, strings.Trim(command," "))

		resp, err := netReader.ReadString('\n')

		if err != nil {
			log.Error.Panicln(err.Error())
		}

		fmt.Println(resp)
	}
}