package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	cv "project_1/commonVariables"
	"strings"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Println("ERROR: ", err)
		}
	}()

	ClientInitialization(ClientArgumentParsing(os.Args))
}

//ClientArgumentParsing - checks, if arguments passed with client are valid
func ClientArgumentParsing(args []string) (string, string) {
	var address = ":9090"
	var host = "127.0.0.1"
	var flagPortOrHost = cv.Nothing
	argsLen := len(args)
	if argsLen == 2 {
		panic(cv.UsageClientErrorMessage)
	}
	if argsLen >= 3 {
		if args[1] == "-h" || args[1] == "--host" {
			host = args[2]
			flagPortOrHost = cv.Host
		} else if args[1] == "-p" || args[1] == "--port" {
			address = args[2]
			flagPortOrHost = cv.Port
		} else {
			panic(cv.UsageClientErrorMessage)
		}
	}
	if argsLen == 4 {
		panic(cv.UsageClientErrorMessage)
	}
	if argsLen == 5 {
		if flagPortOrHost == cv.Host {
			if args[3] == "-p" || args[3] == "--port" {
				address = args[4]
			} else {
				panic(cv.UsageClientErrorMessage)
			}
		} else if flagPortOrHost == cv.Port {
			if args[3] == "-h" || args[3] == "--host" {
				host = args[4]
			}
		} else {
			panic(cv.UsageClientErrorMessage)
		}
	}
	if argsLen > 5 {
		panic(cv.UsageClientErrorMessage)
	}
	return address, host
}

//ClientInitialization - starts the client and exchanges messages with server
func ClientInitialization(address string, host string) {
	serverConnection, err1 := net.Dial("tcp", host + address)
	if err1 != nil {
		fmt.Println(err1)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		commandIsWrong, messageCheckResult := CheckMessage(message)
		if commandIsWrong {
			fmt.Println(messageCheckResult)
		} else {
			serverConnection.Write([]byte(messageCheckResult))
			scanner := bufio.NewReader(serverConnection)
			answer, _ := scanner.ReadString('\n')
			fmt.Print(answer)
		}
	}
}

//CheckMessage - checks, if the command, that we are trying to use, is correct
func CheckMessage(message string) (bool, string) {
	var commandType cv.CommandFlag
	commandIsWrong := false
	switch strings.Split(message, " ")[0] {
	case "GET":
		commandType = cv.GET
	case "SET":
		commandType = cv.SET
	case "DEL":
		commandType = cv.DEL
	default:
		commandIsWrong = true
	}
	if commandIsWrong {
		return commandIsWrong, cv.UsageCommandsErrorMessage
	}
		messageLength := len(message)
		insideValue := false
		quotesCounter := 0
		for i := 4; i < messageLength; i++ {
			if message[i] == ' ' {
				if insideValue {
					part1 := message[:i]
					part2 := message[i+1:]
					array := []string{part1, part2}
					message = strings.Join(array, "_")
				}
			}
			if message[i] == '"' {
				if insideValue {
					insideValue = false
					if i + 1 < messageLength {
						if message[i+1] != ' ' && message[i+1] != '\n' && message[i+1] != '\t' {
							return true, cv.ArgumentsErrorMessage
						}
					}
				} else {
					insideValue = true
					if i - 1 >= 4 {
						if message[i-1] != ' ' && message[i-1] != '\t' {
							return true, cv.ArgumentsErrorMessage
						}
					}
				}
				quotesCounter++
				part1 := message[:i]
				part2 := message[i+1:]
				array := []string{part1, part2}
				message = strings.Join(array, "")
				messageLength--
				i--
			}
		}
		command := strings.Fields(message)
		switch commandType {
		case cv.GET:
			if (quotesCounter != 0 && quotesCounter != 2) || len(command) != 2 {
				commandIsWrong = true
			}
		case cv.SET:
			if (quotesCounter != 0 && quotesCounter != 2 && quotesCounter != 4) ||
				len(command) != 3 {
				commandIsWrong = true
			}
		case cv.DEL:
			if quotesCounter % 2 != 0 {
				commandIsWrong = true
			}
		}
		if commandIsWrong {
			return commandIsWrong, cv.ArgumentsErrorMessage
		}
	return commandIsWrong, message
}