package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var host string
var port string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "The possible host option; default = 127.0.0.1")
	flag.StringVar(&host, "h", "127.0.0.1", "The possible host option; default = 127.0.0.1")
	flag.StringVar(&port, "port", "9090", "The port for listening on; default = 9090")
	flag.StringVar(&port, "p", "9090", "The port for listening on; default = 9090")
}

func main() {
	flag.Parse()

	errMessage, status := CheckHost()
	if !status {
		fmt.Println(errMessage)
	}
	errMessage, status = CheckPort()
	if !status {
		fmt.Println(errMessage)
	}

	dest := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)
	if err != nil {
		fmt.Printf("Can not connect to the server: %s\n", fmt.Sprint(err))
		os.Exit(0)
	}

	go readConnection(conn)

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		conn.Write([]byte(fmt.Sprintf("%s\n", inputScanner.Text())))
	}
}

func readConnection(conn net.Conn) {
	fmt.Println("Connected to the server")
	connScanner := bufio.NewScanner(conn)
	for connScanner.Scan() {
		fmt.Printf(">>> %s\n", connScanner.Text())
	}
	fmt.Println("Lost connection to the server")
	os.Exit(0)
}

func CheckPort() (string, bool) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Sprintf("%s is not valid port value\n", port), false
	}
	if len(port) < 4 || len(port) > 5 ||
		intPort < 0 || intPort > 65535 {
		return fmt.Sprintf("%s is not valid port value!\n", port), false
	}
	return "", true
}

func CheckHost() (string, bool) {
	splittedArgs := strings.Split(host, ".")
	for _, value := range splittedArgs {
		intArg, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Sprintf("%s is not valid host value\n", value), false
		}
		if intArg > 255 || intArg < 0 {
			return fmt.Sprintf("%s is not valid host value\n", value), false
		}
	}
	return "", true
}
