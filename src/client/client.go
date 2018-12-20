package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	publish   = "publish"
	subscribe = "subscribe"
	cmdExit   = "EXIT"
	dataJSON  = "data.json"

	defaultPort = "9090"
	defaultHost = "127.0.0.1"

	fPort = "--port"
	port  = "-p"
	fHost = "--host"
	host  = "-h"

	dump    = "--dump"
	restore = "--restore"

	protocolTCP = "tcp"
)

var (
	_mainPort = defaultPort
	_mainHost = defaultHost

	_dump    = false
	_restore = false

	_subscriber = false //подписчик
	_publisher  = false //слушатель
)

//для парсинга коммандной строки
func trimLastSymbol(s string) string {
	if last := len(s) - 1; last >= 0 && s[last] == ',' {
		s = s[:last]
		return s
	}
	return s
}

func parseArguments() {

	fmt.Println("Parse arguments")
	var cmds = make(map[string]func(string))
	// заполняем команды
	cmds[port] = func(s string) {
		//номер порта
		s = trimLastSymbol(s)
		fmt.Printf("Parse -p=%s\r\n", s)
		if port, err := strconv.Atoi(s); err == nil {
			//Valid numbers for ports are: 0 to 2^16-1 = 0 to 65535
			//But user ports 1024 to 49151
			if port < 49152 && port > 1024 {
				_mainPort = s
			}
		} else {
			_mainPort = defaultPort
		}
		fmt.Printf("New port: %s\r\n", _mainPort)
	}
	cmds[fPort] = func(s string) {
		s = trimLastSymbol(s)
		fmt.Printf("Parse --port=%s\r\n", s)
		if port, err := strconv.Atoi(s); err == nil {
			//Valid numbers for ports are: 0 to 2^16-1 = 0 to 65535
			//But user ports 1024 to 49151
			if port < 49152 && port > 1024 {
				_mainPort = s
			}
		} else {
			_mainPort = defaultPort
		}
		fmt.Printf("New port: %s\r\n", _mainPort)
	}
	cmds[host] = func(s string) {
		s = trimLastSymbol(s)
		//путь к файлу для сохранения
		fmt.Printf("Parse --h=%s\r\n", s)
		// TODO add validate ip address
		fmt.Printf("New IP address: %s\r\n", s)
		_mainHost = s // update ip address
	}
	cmds[fHost] = func(s string) {
		s = trimLastSymbol(s)
		//путь к файлу с сохранением
		fmt.Printf("--host=%s\r\n", s)
		// TODO add validate ip address
		fmt.Printf("New IP address: %s\r\n", s)
		_mainHost = s // update ip address
	}
	cmds[dump] = func(s string) {
		s = trimLastSymbol(s)
		//путь к файлу с сохранением
		fmt.Println("--dump")
		_dump = true
	}
	cmds[restore] = func(s string) {
		s = trimLastSymbol(s)
		//путь к файлу с сохранением
		fmt.Println("--restore")
		_restore = true
	}

	//анализ команд
	for _, arg := range os.Args[1:] {
		cmd := strings.Split(arg, "=")

		if len(cmd) > 2 {
			fmt.Println(arg, "don't know... ")
			continue
		}

		// для красноречия
		name := cmd[0]
		param := ""
		if trimLastSymbol(cmd[0]) == dump || trimLastSymbol(cmd[0]) == restore {
			name = trimLastSymbol(cmd[0])
			param = ""
		} else {
			param = cmd[1]
		}

		// ищем функцию
		fn := cmds[name]
		if fn == nil {
			fmt.Println(name, "don't know... ")
			continue
		}

		// исполняем команду
		fn(param)
	}
}

func main() {

	if len(os.Args) > 0 {
		parseArguments()
	}

	// connect to this socket
	address := _mainHost + ":" + _mainPort
	conn, err := net.Dial(protocolTCP, address)

	if err == nil && conn != nil {

		if _dump {
			// send to socket dump command
			fmt.Fprintf(conn, "dump\n")
			// get file
			file, _ := os.Create(dataJSON)
			defer file.Close()
			n, err := io.Copy(file, conn)
			if err == io.EOF {
				fmt.Println(err.Error())
			}
			fmt.Println("Bytes received", n)
		} else if _restore {
			// send to socket restore command
			fmt.Fprintf(conn, "restore ")

			file, errF := os.Open(dataJSON) // For read access.
			if errF != nil {
				fmt.Println("Unable to open file, " + errF.Error())
			}
			defer file.Close() // make sure to close the file even if we panic.
			n, error := io.Copy(conn, file)
			if error != nil {
				fmt.Printf("Send file error %s\r\n", error.Error())
			}
			fmt.Println(n, "Bytes sent")
			fmt.Fprintf(conn, "\r\n")
			file.Close()
		} else {

			for {
				// read in input from stdin
				if !_subscriber {
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Text to send: ")
					text, _ := reader.ReadString('\n')
					if strings.ToUpper(strings.TrimRight(text, "\r\n")) == cmdExit {
						//пользователь решил удалиться)))
						break
					}

					if strings.ToUpper(strings.TrimRight(text, "\r\n")) == publish {
						_publisher = true
						_subscriber = false
					}
					if strings.ToUpper(strings.TrimRight(text, "\r\n")) == subscribe {
						_subscriber = true
						_publisher = false
					}

					// send to socket
					fmt.Fprintf(conn, text+"\n")
				}

				// listen for reply
				if !_publisher {
					message, _ := bufio.NewReader(conn).ReadString('\n')
					fmt.Println("Message from client: " + message)
				}
			}
		}

	} else {
		fmt.Print("Didn't connect")
		if err != nil {
			fmt.Printf(", Error: %s\n\r", err.Error())
		}
	}
	defer conn.Close()
}

