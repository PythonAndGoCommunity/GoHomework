package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

const DataFile = "server/server/data/data.json"

var Channels = make(map[net.Conn][]string)
var Data DataType
var Verbose bool
var ExitChannel = make(chan bool)

func ConfigureLogging() {
	f, err := os.OpenFile("server/server.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ERROR: Cannot open file %v", err)
	}
	log.SetOutput(f)
}

func HandleConnection(conn net.Conn) {
	clientAddress := conn.RemoteAddr().String()
	PrintlnAndLog("New connection established (" + clientAddress + ")")

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		PrintlnAndLog("New message " + scanner.Text() + " from " + clientAddress)

		var response string
		if strings.Compare(scanner.Text(), "") != 0 {
			response = ApplyCommand(scanner.Text(), conn)
		} else {
			response = "ERR Empty message"
		}

		_, err := conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("ERROR: cannot respond the request.")
		}
	}

	PrintlnAndLog("Client with address " + clientAddress + " disconnected.")
}

func ApplyCommand(cmdString string, conn net.Conn) string {
	cmdName, cmdList := getCommandArguments(cmdString)
	response := "ERR Unknown command"

	if strings.Compare(cmdName, "STOP") == 0 {
		SaveData(DataFile)
		response = "You've stopped the server successfully."
		ExitChannel <- true

	} else if strings.Compare(cmdName, "HELP") == 0 {
		response = " Goredis is a simple implementation of Redis on Golang.\n"
		response += "\r Available commands: HELP SET GET DEL KEYS (UN)SUBSCRIBE PUBLISH EXIT STOP.\n"
		response += "\r To get more help on them, type any command with no arguments."

	} else if strings.Contains(cmdName, "SET") {
		if len(cmdList) > 1 {
			response = "OK"
			AddEntry(cmdList[0], cmdList[1])
			SaveData(DataFile)
		} else {
			response = "Usage: SET <key> <value>"
		}

	} else if strings.Contains(cmdName, "GET") {
		if len(cmdList) == 1 {
			response = GetEntry(cmdList[0])
		} else {
			response = "Usage: GET <key> or GET <key> <value>"
		}

	} else if strings.Contains(cmdName, "DEL") {
		if len(cmdList) > 0 {
			response = RemoveEntries(cmdList)
			SaveData(DataFile)
		} else {
			response = "Usage: DEL <key>"
		}

	} else if strings.Contains(cmdName, "KEYS") {
		if len(cmdList) == 0 || strings.Compare("", cmdList[0]) == 0 || strings.Compare("*", cmdList[0]) == 0 {
			response = ShowAllKeys()
		} else {
			response = FindKeys(cmdList[0])
		}

	} else if strings.Contains(cmdName, "UNSUBSCRIBE") {
		if len(cmdList) == 1 {
			Channels[conn] = removeFromSlice(cmdList[0], Channels[conn])
			response = "OK"
		} else {
			response = "Usage: UNSUBSCRIBE <channel>"
		}

	} else if strings.Contains(cmdName, "SUBSCRIBE") {
		if len(cmdList) == 1 {
			Channels[conn] = append(Channels[conn], cmdList[0])
			response = "OK"
		} else {
			response = "Usage: SUBSCRIBE <channel>"
		}

	} else if strings.Contains(cmdName, "PUBLISH") {
		if len(cmdList) == 2 {
			for k, v := range Channels {
				if stringInSlice(cmdList[0], v) && conn != k {
					_, err := k.Write([]byte("\rNew message in channel " + cmdList[0] + ": " + cmdList[1] + "\n"))
					if err != nil {
						fmt.Println("ERROR: Cannot send message in channel.")
					}
				}
			}
			response = "\rOK              "
		}

	} else if strings.Contains(cmdName, "DUMP") {
		jsonData, _ := LoadFromFile(DataFile)
		response = fmt.Sprintf("\rDumped database: %s", jsonData)

	} else if strings.Contains(cmdName, "RESTORE") {
		err := LoadData(cmdList[0])
		if err != nil {
			PrintlnAndLog("Data cannot be restored by client request.")
			response = fmt.Sprintf("\rCannot restore data from " + cmdList[0] + ".")
		} else {
			SaveData(DataFile)
			PrintlnAndLog("Data has been restored by client request.")
			response = fmt.Sprintf("\rYou've been successfully restored data from " + cmdList[0] + ".")
		}
	}

	return response
}

func getCommandArguments(cmdString string) (string, []string) {
	cmdString = getRidOfSpaces(cmdString)
	pattern := regexp.MustCompile(`^[A-Za-z]+`)
	cmdName := pattern.FindString(cmdString)
	cmdName = strings.ToUpper(cmdName)

	cmdList := splitArguments(cmdString)
	return cmdName, cmdList[1:]
}

func splitArguments(s string) []string {
	var args []string
	inQuotes := false
	lastSplitted := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			if inQuotes == false {
				inQuotes = true
			} else {
				inQuotes = false
			}
		}

		if !inQuotes {
			if s[i] == ' ' {
				arg := strings.Replace(s[lastSplitted:i], "\"", "", -1)
				args = append(args, arg)
				lastSplitted = i + 1
			} else if i == len(s) - 1 {
				arg := strings.Replace(s[lastSplitted:i + 1], "\"", "", -1)
				args = append(args, arg)
			}
		}
	}
	return args
}

func getRidOfSpaces(s string) string {
	pattern := regexp.MustCompile(`[\s]{2,}`)
	s = pattern.ReplaceAllString(s, " ")
	return s
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Compare(a, b) == 0 {
			return true
		}
	}
	return false
}

func removeFromSlice(a string, list []string) []string {
	for i, b := range list {
		if strings.Compare(a, b) == 0 {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}
