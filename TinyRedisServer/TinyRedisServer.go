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

var address = "127.0.0.1"
var port, mode string
var dataMutex = sync.RWMutex{}
var mapData = make(map[string]string)
var diskDir = "data.txt"

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

	commands := make(chan command)
	go storage(commands) // performing client's commands

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(commands, conn)
	}
}

func contains(strs []string, str string) bool {
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
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		if msg == "" {
			continue
		}
		flds := strings.Fields(msg)

		result := make(chan string)
		commands <- command{
			fields: flds,
			result: result,
		}
		conn.Write([]byte(<-result + "\n>>> "))
	}
}

func storage(cmd chan command) {
	for cmd := range cmd {
		if len(cmd.fields) < 1 {
			cmd.result <- ""
			continue
		}
		if len(cmd.fields) < 2 {
			cmd.result <- "Expected at least 2 arguments"
			continue
		}

		fmt.Println("Command:", cmd.fields)
		switch cmd.fields[0] {
		case "GET":
			dataMutex.RLock()
			if mode == "disk" {
				if err := readFromFile(); err != nil {
					panic(err)
				}
			}
			if val, ok := mapData[cmd.fields[1]]; !ok {
				cmd.result <- "nil"
			} else {
				cmd.result <- val
			}
			dataMutex.RUnlock()

		case "SET":
			dataMutex.Lock()
			if len(cmd.fields) < 3 {
				cmd.result <- "Warning: EXPECTED VALUE"
				continue
			} else if len(cmd.fields) > 3 {
				cmd.result <- "Expected Key and Value"
				continue
			}
			if mode == "disk" {
				if err := readFromFile(); err != nil {
					panic(err)
				}
			}
			mapData[cmd.fields[1]] = cmd.fields[2]
			if mode == "disk" {
				if err := writeToFile(); err != nil {
					panic(err)
				}
			}
			cmd.result <- "done"
			dataMutex.Unlock()

		case "DEL":
			dataMutex.Lock()
			if mode == "disk" {
				if err := readFromFile(); err != nil {
					panic(err)
				}
			}
			delete(mapData, cmd.fields[1])
			if mode == "disk" {
				if err := writeToFile(); err != nil {
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
func readFromFile() error {
	file, err := os.Open(diskDir)
	if err != nil {
		return err
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
	return scanner.Err()
}

func writeToFile() error {
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
