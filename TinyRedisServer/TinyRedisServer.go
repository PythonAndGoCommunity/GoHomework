package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type cmdFlagHandler struct { // different flags handler
	options       []string
	defaultVal    string
	warning       string
	possibleRange []string
}

type command struct {
	fields []string
	result chan string
}

var address = "127.0.0.1"
var port, mode string
var dataMutex = sync.RWMutex{}
var diskDir = "ServerData.txt" // defaulf file for "-m=disk"option
var commands = make(chan command)

func main() {
	port, mode = cmdFlagParse() // parsing results from cmd
	fullAddress := address + ":" + port

	l, err := net.Listen("tcp", fullAddress)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.Printf("Server is running on %s\n", fullAddress)
	log.Println("Ready to accept connections")

	var mapData = make(map[string]string) // creation of map for server data
	createStorage(mapData)                // creation of storage

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(commands, conn)
	}
}

func contains(strs []string, str string) bool { // function for flag parsing
	if strs == nil {
		return true
	}
	for _, val := range strs {
		if str == val {
			return true
		}
	}
	return false
}

func cmdFlagParse() (string, string) {
	portHandler := cmdFlagHandler{ // different flag options
		options:    []string{"p", "port"},
		defaultVal: "9090",
		warning:    "specify port to use",
	}

	modeHandler := cmdFlagHandler{
		options:       []string{"m", "mode"},
		defaultVal:    "ram",
		warning:       "specify memory mode",
		possibleRange: []string{"ram", "disk"},
	}

	// long and short forms are for long and short flags, e.g. -p and --port
	portShortResult := flag.String(portHandler.options[0], portHandler.defaultVal, portHandler.warning)
	portLongResult := flag.String(portHandler.options[1], portHandler.defaultVal, portHandler.warning)

	modeShortResult := flag.String(modeHandler.options[0], modeHandler.defaultVal, modeHandler.warning)
	modeLongResult := flag.String(modeHandler.options[1], modeHandler.defaultVal, modeHandler.warning)

	flag.Parse()

	var portResult, modeResult string
	if *portLongResult != portHandler.defaultVal {
		portResult = *portLongResult
	} else {
		portResult = *portShortResult
	}

	if *modeLongResult != modeHandler.defaultVal && contains(modeHandler.possibleRange, *modeLongResult) {
		modeResult = *modeLongResult
	} else if contains(modeHandler.possibleRange, *modeShortResult) {
		modeResult = *modeShortResult
	} else {
		modeResult = modeHandler.defaultVal
	}

	return portResult, modeResult
}

func handle(commands chan command, conn net.Conn) { // client input processing
	defer func() {
		conn.Close()
		log.Println("Connection closed")
	}()

	log.Println("Connection from", conn.RemoteAddr())
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		// if client was unexpectadly unconnected, server is able to detect it
		// and close connection via following lines:
		if err != nil {
			return
		}

		flds := strings.Fields(msg) // result formatting

		if len(flds) > 0 {
			if flds[0] == "EXIT" { // decent exit processing
				return
			}
		}

		result := make(chan string)
		commands <- command{
			fields: flds,
			result: result,
		}
		conn.Write([]byte(<-result + "\n>>> "))
	}
}

func createStorage(mapData map[string]string, testStrs ...string) string {
	go storage(commands, mapData) // performing client's commands

	if testStrs != nil { // it is here just for test purposes
		result := make(chan string)
		flds := strings.Fields(testStrs[0])
		commands <- command{
			fields: flds,
			result: result,
		}
		return <-result
	}
	return ""
}

func storage(cmd chan command, mapData map[string]string) {
	for cmd := range cmd {
		if len(cmd.fields) < 1 {
			cmd.result <- ""
			continue
		}
		if len(cmd.fields) < 2 {
			cmd.result <- "Expected at least 2 arguments!"
			continue
		}

		switch cmd.fields[0] {
		case "GET":
			dataMutex.RLock()
			if mode == "disk" {
				var err error
				if mapData, err = readFromFile(); err != nil {
					panic(err)
				}
			}
			val, ok := mapData[cmd.fields[1]]
			if !ok {
				cmd.result <- "nil"
			} else {
				cmd.result <- val
			}
			dataMutex.RUnlock()

		case "SET":
			if len(cmd.fields) < 3 {
				cmd.result <- "Expected value!"
				continue
			} else if len(cmd.fields) > 3 {
				cmd.result <- "Expected Key and Value!"
				continue
			}
			dataMutex.Lock()
			if mode == "disk" {
				var err error
				if mapData, err = readFromFile(); err != nil {
					panic(err)
				}
			}
			mapData[cmd.fields[1]] = cmd.fields[2]
			if mode == "disk" {
				if err := writeToFile(mapData); err != nil {
					panic(err)
				}
			}
			cmd.result <- "done"
			dataMutex.Unlock()

		case "DEL":
			dataMutex.Lock()
			if mode == "disk" {
				var err error
				if mapData, err = readFromFile(); err != nil {
					panic(err)
				}
			}
			delete(mapData, cmd.fields[1])
			if mode == "disk" {
				if err := writeToFile(mapData); err != nil {
					panic(err)
				}
			}
			cmd.result <- "deleted"
			dataMutex.Unlock()

		default:
			cmd.result <- fmt.Sprintf("Invalid command \"%s\"", cmd.fields[0])
		}
	}
}

// reading from and writing to disk via following functions:
func readFromFile() (map[string]string, error) {
	var mapData = make(map[string]string)
	file, err := os.Open(diskDir)
	if err != nil {
		return mapData, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, line := range lines { // fillling out map variable
		splitted := strings.Fields(string(line))
		mapData[splitted[0]] = splitted[1]
	}
	return mapData, scanner.Err()
}

func writeToFile(mapData map[string]string) error {
	file, err := os.Create(diskDir)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	for key, val := range mapData { // conversion from map to lines in file
		line := key + " " + val
		lines = append(lines, line)
	}

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
