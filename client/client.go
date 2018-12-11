package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		var source string
		fmt.Print("Enter command: ")
		myscanner := bufio.NewScanner(os.Stdin)
		myscanner.Scan()
		source = myscanner.Text()
		if len(source) == 0 {
			fmt.Println("Wrong input: no command")
			continue
		}
		// send message
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			fmt.Println(err)
			return
		}
		// get response
		fmt.Print("Server response:")
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		fmt.Print(string(buff[0:n]))
		fmt.Println()
	}
}
