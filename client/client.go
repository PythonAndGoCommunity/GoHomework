package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "9090", "listening port")
	flag.StringVar(&port, "p", "9090", "listening port")
	var host string
	flag.StringVar(&host, "h", "127.0.0.1", "listening IP")
	flag.StringVar(&host, "host", "127.0.0.1", "listening IP")

	flag.Parse()

	address := fmt.Sprintf("%v:%v", host, port)
	conn, err := net.Dial("tcp", address)
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
