package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
	cv "project_1/commonVariables"
	"strings"
)




func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Println("ERROR: ", err)
		}
	}()

    ServerInitialization(ServerArgumentParsing(os.Args))
}


//ServerArgumentParsing - checks, if arguments passed with server are valid
func ServerArgumentParsing (args []string) (string, string) {
	var address = ":9090"
	var mode = "disk"
	var flagModeOrPort = cv.Nothing
	argsLen := len(args)
	if argsLen == 2 {
		panic(cv.UsageServerErrorMessage)
	}
	if argsLen >= 3 {
		if args[1] == "-m" || args[1] == "--mode"{
			if args[2] == "disk" {
				flagModeOrPort = cv.Mode
				mode = "disk"
			} else {
				panic(cv.UsageServerErrorMessage)
			}
		} else if args[1] == "-p" || args[1] == "--port"{
			address = args[2]
			flagModeOrPort = cv.Port
		} else {
			panic(cv.UsageServerErrorMessage)
		}
	}
	if argsLen == 4 {
		panic(cv.UsageServerErrorMessage)
	}
	if argsLen == 5 {
		if flagModeOrPort == cv.Mode {
			if args[3] == "-p" || args[3] == "--port"{
				address = args[4]
			} else {
				panic(cv.UsageServerErrorMessage)
			}
		} else if flagModeOrPort == cv.Port {
			if args[3] == "-m" || args[3] == "--mode" {
				if args[4] != "disk" {
					panic(cv.UsageServerErrorMessage)
				} else {
					mode = "disk"
				}
			}
		} else {
			panic(cv.UsageServerErrorMessage)
		}
	}
	if argsLen > 5 {
		panic(cv.UsageServerErrorMessage)
	}
	return address, mode
}

//ServerInitialization - starts the server
func ServerInitialization(address string, mode string) {
	channelClient, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	defer channelClient.Close()

	log.Printf("Server is ready to work on port %s", address)

	commands := make(chan []string)
	answers := make(chan cv.Answer)
	go Database(commands, answers)

	for {
		client, err2 := channelClient.Accept()
		if err2 != nil {
			log.Fatalln(err2)
		}

		go HandleClients(commands, client, answers)
	}
}

//HandleClients - handles communication with client
func HandleClients(commands chan<- []string, client net.Conn, answers <-chan cv.Answer) {
	defer func() {
		log.Println("Client ", client.RemoteAddr(), " disconnected.")
		client.Close()
	}()

	log.Println("Client ", client.RemoteAddr(), " connected.")

	scanner := bufio.NewScanner(client)
	for scanner.Scan() {
		convertedText := scanner.Text()
		log.Println("Got a message - ", convertedText)
		command := strings.Fields(convertedText)
		commands<-command
		result := <-answers
		cmdAnswer := strings.Replace(result.Answer, "_", " ", -1)
		if result.State == cv.ERROR {
			client.Write([]byte("ERROR: " + cmdAnswer + "\n"))
		} else {
			client.Write([]byte(cmdAnswer + "\n"))
		}
	}
}

//Database - stores the data and processes commands
func Database(commands <-chan []string, answers chan<- cv.Answer) {
	data := make(map[string]string)

	for command := range commands {
		if len(command) < 1 {
			answers <- cv.Answer{Answer: "", State: cv.ERROR}
			continue
		}
		if len(command) < 2 {
			answers <- cv.Answer{Answer: "Expected at least 1 argument.", State: cv.ERROR}
			continue
		}
		switch command[0] {

		case "GET":
			if len(command) > 2 {
				answers <- cv.Answer{Answer: cv.GetErrorMessage, State: cv.ERROR}
				continue
			} else {
				value, contains := data[command[1]]
				if contains {
					answers <- cv.Answer{Answer: value, State: cv.Present}
				} else {
					answers <- cv.Answer{Answer: "(nil)", State: cv.Absent}
				}
			}

		case "SET":
			if len(command) > 3 {
				answers <- cv.Answer{Answer: cv.SetErrorMessage, State: cv.ERROR}
				continue
			} else {
				prevValue, contains := data[command[1]]
				data[command[1]] = command[2]
				if contains {
					answers <- cv.Answer{
						Answer: "Replaced previous value - " + prevValue + ".", State: cv.Ignored}
				} else {
					answers <- cv.Answer{Answer: "OK.", State: cv.Ignored}
				}
			}

		case "DEL":
			numberOfKeys := len(command) - 1
			deleted := 0
			iterator := 1
			var deleteReports []string
			for iterator <= numberOfKeys {
				if _, contains := data[command[iterator]]; contains {
					delete(data, command[iterator])
					deleted++
					deleteReports = append(deleteReports, strconv.Itoa(iterator) + " - deleted.")
				} else {
					deleteReports = append(deleteReports, strconv.Itoa(iterator) + " - ignored.")
				}
				iterator++
			}
			answers <- cv.Answer{ Answer: strings.Join(deleteReports, " "), State: cv.Ignored}

		default:
			answers <- cv.Answer{ Answer: "Wrong command.", State: cv.ERROR}
			continue
		}
		log.Println(data)
	}
}
