package handler

import (
	"net"
	"bufio"
	"os"
	"fmt"
	"NonRelDB/log"
)

func HandleConnection(c net.Conn){
	for {
		consoleReader := bufio.NewReader(os.Stdin)
		// netReader := bufio.NewReader(c)
		// netWriter := bufio.NewWriter(c)

		fmt.Print("\nNonRelDB> ")
		command, err := consoleReader.ReadString('\n')
		command = command[:len(command)-1]
		
		if err != nil {
			log.Error.Panicln(err.Error())
		}
	
		if command == "exit"{
			fmt.Println("Good bye")
			c.Write([]byte(command))
			return
		}

		req := []byte(command)
		c.Write(req)

		resp := make([]byte, 1024)
		lenResp, err := c.Read(resp)

		if err != nil {
			log.Error.Panicln(err.Error())
		}

		value := string(resp[:lenResp])

		fmt.Println(value)
	}
}