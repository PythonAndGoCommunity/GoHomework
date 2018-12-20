package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var storage string
var port string
var verboseLog bool

func init() {
	flag.StringVar(&storage, "mode", "disk", "The possible storage option; default = disk")
	flag.StringVar(&storage, "m", "disk", "The possible storage option; default = disk")
	flag.StringVar(&port, "port", "9090", "The port for listening on; default = 9090")
	flag.StringVar(&port, "p", "9090", "The port for listening on; default = 9090")
	flag.BoolVar(&verboseLog, "verbose", false, "Full log of the client requests; default = false")
	flag.BoolVar(&verboseLog, "v", false, "Full log of the client requests; default = false")
}

const (
	COMMAND_SET       = "SET"
	COMMAND_GET       = "GET"
	COMMAND_DEL       = "DEL"
	COMMAND_KEYS      = "KEYS"
	COMMAND_PUBLISH   = "PUBLISH"
	COMMAND_SUBSCRIBE = "SUBSCRIBE"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Channel struct {
	Name    string
	Clients *[]net.Conn
}

type ChannelCollection struct {
	Channels []Channel
}

func (c *ChannelCollection) Subscribe(name string, conn net.Conn) (string, bool) {
	for _, channel := range c.Channels {
		if channel.Name == name {
			for _, cnn := range *channel.Clients {
				if cnn == conn {
					return fmt.Sprintf("You are already subscribed to [%s] channel", name), false
				}
			}
			*channel.Clients = append(*channel.Clients, conn)
			return "", false
		}
	}
	c.Channels = append(c.Channels, Channel{Name: name, Clients: &[]net.Conn{}})
	fmt.Printf("[%s] channel created\n", name)
	c.Subscribe(name, conn)
	return fmt.Sprintf("You are now subscribed to [%s] channel", name), true
}

func (c *ChannelCollection) Publish(message string, conn net.Conn) (string, bool) {
	channelIndex := 0
	connIndex := 0
	var channel Channel
	var cnn net.Conn
OUTER:
	for channelIndex, channel = range c.Channels {
		for connIndex, cnn = range *channel.Clients {
			if cnn == conn {
				break OUTER
			}
		}
	}
	for index, cnn := range *c.Channels[channelIndex].Clients {
		if connIndex == index {
			continue
		}
		_, err := cnn.Write([]byte(fmt.Sprintf("%s: %s\n", conn.RemoteAddr().String(), message)))
		if err != nil {
			return fmt.Sprint(err), false
		}
	}
	return "", true
}

func (c *ChannelCollection) Unsubscribe(newChannelName string, conn net.Conn) {
	indx := 0
	var channel Channel
OUTER:
	for indx, channel = range c.Channels {
		for index, cnn := range *channel.Clients {
			if cnn == conn && channel.Name != newChannelName {
				mySlice := make([]net.Conn, len(*channel.Clients))
				mySlice = append(mySlice[:index], mySlice[index+1:]...)
				*channel.Clients = mySlice
				break OUTER
			}
		}
	}
	if len(*c.Channels[indx].Clients) == 0 {
		fmt.Printf("[%s] channel removed\n", c.Channels[indx].Name)
		c.Channels = append(c.Channels[:indx], c.Channels[indx+1:]...)
	}
}

var channels = &ChannelCollection{}
var keyArray = []string{}
var keyValueMap = make(map[string]string)
var isDatabaseInUse = false

func main() {
	const SERVER_ADDRESS = "127.0.0.1"
	flag.Parse()
	fmt.Printf(".::SERVER OPTIONS::.\n")

	message, isCorrect := CheckPort()
	if !isCorrect {
		fmt.Print(message)
		os.Exit(0)
	}
	fmt.Print(message)

	message, isCorrect = CheckMode()
	if !isCorrect {
		fmt.Print(message)
		os.Exit(0)
	}
	fmt.Print(message)

	fmt.Printf("Verbose: %v [OK]\n", verboseLog)
	fmt.Println("\nThe server is starting...")

	fullServerAddress := fmt.Sprintf("%s:%s", SERVER_ADDRESS, port)
	listener, err := net.Listen("tcp", fullServerAddress)
	if err != nil {
		fmt.Println("Unable to listen on %s\n", fullServerAddress)
		os.Exit(0)
	}

	fmt.Printf("Listening on %s\n", fullServerAddress)

	defer listener.Close()

	go SetDatabaseRefresher()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Unable to accept client: %s\n", fmt.Sprint(err))
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected from %s\n", remoteAddr)

	disconnectTimer := time.NewTimer(60 * time.Second)
	defer disconnectTimer.Stop()

	go func() {
		<-disconnectTimer.C
		conn.Write([]byte("Disconnected due to inactivity\n"))
		conn.Close()
	}()

	connScanner := bufio.NewScanner(conn)

	for connScanner.Scan() {
		if !disconnectTimer.Stop() {
			disconnectTimer.Reset(60 * time.Second)
		}
		messageResponse, state := handleMessage(connScanner.Text(), conn)
		if !state {
			conn.Write([]byte(fmt.Sprintf("%s\n", messageResponse)))
		} else if state && messageResponse != "" {
			conn.Write([]byte(fmt.Sprintf("%s\n", messageResponse)))
		}
	}

	fmt.Printf("Client from %s disconnected\n", remoteAddr)
}

func handleMessage(message string, conn net.Conn) (string, bool) {

	if verboseLog {
		fmt.Printf("[V]%s: %s\n", conn.RemoteAddr().String(), message)
	}

	splittedMessage, isSplitted := CheckInputForWhitespaces(message)
	if !isSplitted {
		return "", true
	}

	switch splittedMessage[0] {
	case COMMAND_SET:
		if storage == "disk" {
			commandResponse, isSet := CommandSET(splittedMessage[1:])
			if !isSet {
				return fmt.Sprintf("%s", commandResponse), false
			}
			return fmt.Sprintf("%s", commandResponse), true
		}
		commandResponse, isSet := CommandSET_memory(splittedMessage[1:])
		if !isSet {
			return fmt.Sprintf("%s", commandResponse), false
		}
		return fmt.Sprintf("%s", commandResponse), true

	case COMMAND_GET:
		if storage == "disk" {
			commandResponse, isGotten := CommandGET(splittedMessage[1:])
			if !isGotten {
				return fmt.Sprintf("%s", commandResponse), false
			}
			return fmt.Sprintf("%s", commandResponse), true
		}
		commandResponse, isGotten := CommandGET_memory(splittedMessage[1:])
		if !isGotten {
			return fmt.Sprintf("%s", commandResponse), false
		}
		return fmt.Sprintf("%s", commandResponse), true

	case COMMAND_DEL:
		if storage == "disk" {
			commandResponse, isDisabled := CommandDEL(splittedMessage[1:])
			if !isDisabled {
				return fmt.Sprintf("%s", commandResponse), false
			}
			return fmt.Sprintf("%s", commandResponse), true
		}
		commandResponse, isDisabled := CommandDEL_memory(splittedMessage[1:])
		if !isDisabled {
			return fmt.Sprintf("%s", commandResponse), false
		}
		return fmt.Sprintf("%s", commandResponse), true

	case COMMAND_KEYS:
		if storage == "disk" {
			commandResponse, isFound := CommandKEYS(splittedMessage[1:], conn)
			if !isFound {
				return fmt.Sprintf("%s", commandResponse), false
			}
			return "", true
		}
		commandResponse, isFound := CommandKEYS_memory(splittedMessage[1:], conn)
		if !isFound {
			return fmt.Sprintf("%s", commandResponse), false
		}
		return "", true

	case COMMAND_SUBSCRIBE:
		commandResponse, isSubscribed := CommandSUBSCRIBE(splittedMessage[1:], conn)
		if !isSubscribed {
			return fmt.Sprintf("%s", commandResponse), false
		}
		return fmt.Sprintf("%s", commandResponse), true

	case COMMAND_PUBLISH:
		commandResponse, isSent := CommandPUBLISH(splittedMessage[1:], conn)
		if !isSent {
			return fmt.Sprintf("%s", commandResponse), false
		}
		return fmt.Sprintf("%s", commandResponse), true

	default:
		return fmt.Sprintf("Command %s is not found", splittedMessage[0]), false
	}
}

func WriteKeyValueToFile(key string, value string) bool {

	file, err := os.OpenFile("database.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Unable to 'open'/'create and open' file: %s\n", fmt.Sprint(err))
		return false
	}
	isDatabaseInUse = true
	defer func() {
		isDatabaseInUse = false
	}()
	defer file.Close()

	keyValue := KeyValue{
		Key:   key,
		Value: value,
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	err = encoder.Encode(&keyValue)
	if err != nil {
		fmt.Printf("Unable to encode [%s:%s] to JSON format: %s\n", key, value, fmt.Sprint(err))
		return false
	}

	_, err = io.Copy(file, buffer)
	if err != nil {
		fmt.Printf("Unable to copy [%s:%s] to the database: %s\n", key, value, fmt.Sprint(err))
		return false
	}

	return true
}

func CommandDEL(splittedCommandArgs []string) (string, bool) {

	runeCounter := 0
	runeIndexes := []int{}
	key := ""

	for index, value := range splittedCommandArgs {
		if value == "$" {
			runeCounter += 1
			runeIndexes = append(runeIndexes, index)
		}
	}

	if runeCounter%2 != 0 || runeCounter > 2 {
		return "Incorrect amount of $ runes", false
	}

	if runeCounter == 0 {
		if len(splittedCommandArgs) != 1 {
			return "Incorrect amount of DEL args", false
		}

		key = splittedCommandArgs[0]
		if key == "/$" {
			key = "$"
		}
	} else {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			return "KEY can not be empty", false
		}

		key = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(key); index += 1 {
			if key[index] == '$' && key[index-1] == '/' {
				key = key[:index-1] + key[index:]
			}
		}
	}

	file, err := os.OpenFile("database.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Unable to 'open'/'create and open' file: %s\n", fmt.Sprint(err))
		return "Unable to 'open'/'create and open' file", false
	}
	isDatabaseInUse = true
	defer func() {
		isDatabaseInUse = false
	}()
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	keyValue := &KeyValue{}

OUTER:
	for fileScanner.Scan() {
		decoder := json.NewDecoder(strings.NewReader(fileScanner.Text()))
		err = decoder.Decode(&keyValue)
		if err != nil {
			fmt.Printf("Unable to decode JSON to [KEY:VALUE]: %s\n", fmt.Sprint(err))
			return "Unable to decode JSON to [KEY:VALUE]", false
		}

		for _, value := range keyArray {
			if keyValue.Key == value {
				continue OUTER
			}
		}

		if keyValue.Key == key {
			if len(keyArray) > 10 {
				if isDatabaseInUse {
					return "Database is in use. Try again later", false
				}
				RefreshDatabase()
			}
			keyArray = append(keyArray, key)

			return fmt.Sprintf("[%s:%s] pair was added for removing from the database", keyValue.Key, keyValue.Value), true
		}
	}

	return fmt.Sprintf("[%s] key is not found", key), false
}

func RefreshDatabase() bool {
	primaryFile, err := os.OpenFile("database.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Unable to 'open'/'create and open' the database: %s\n", fmt.Sprint(err))
		return false
	}
	isDatabaseInUse = true
	defer func() {
		isDatabaseInUse = false
	}()
	defer primaryFile.Close()

	conversionFile, err := os.OpenFile("tempdatabase.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("Unable to 'open'/'create and open' temporary storage: %s\n", fmt.Sprint(err))
		return false
	}
	defer conversionFile.Close()

	fileScanner := bufio.NewScanner(primaryFile)
	keyValue := &KeyValue{}

OUTER:
	for fileScanner.Scan() {
		decoder := json.NewDecoder(strings.NewReader(fileScanner.Text()))

		err = decoder.Decode(&keyValue)
		if err != nil {
			fmt.Printf("Unable to decode %s to [KEY:VALUE]: %s\n", fileScanner.Text(), fmt.Sprint(err))
			return false
		}

		for _, value := range keyArray {
			if keyValue.Key == value {
				continue OUTER
			}
		}

		_, err = conversionFile.WriteString(fmt.Sprintf("%s\n", fileScanner.Text()))
		if err != nil {
			fmt.Printf("Unable to write data to conversion file: %s\n", fmt.Sprint(err))
			return false
		}
	}

	err = os.Remove("database.json")
	if err != nil {
		fmt.Printf("Unable to remove primary database: %s\n", fmt.Sprint(err))

		err = os.Remove("tempdatabase.json")
		if err != nil {
			fmt.Printf("Unable to remove conversion database: %s\n", fmt.Sprint(err))
		}
		return false
	}

	err = os.Rename("tempdatabase.json", "database.json")
	if err != nil {
		fmt.Printf("Unable to rename conversion file: %s\n", fmt.Sprint(err))
		return false
	}

	keyArray = nil
	return true
}

func CommandSET(splittedCommandArgs []string) (string, bool) {

	if len(splittedCommandArgs) < 2 {
		return "Incorrect amount of SET arguments", false
	}

	key := ""
	value := ""
	runeIndexes := []int{}

	runeCounter := 0
	for index, value := range splittedCommandArgs {
		if value == "$" {
			runeCounter += 1
			runeIndexes = append(runeIndexes, index)
		}
	}

	if runeCounter%2 != 0 || runeCounter > 4 {
		return "Incorrect amount of $ runes", false
	}

	if runeCounter == 0 {
		if len(splittedCommandArgs) == 2 {
			key = splittedCommandArgs[0]
			if key == "/$" {
				key = "$"
			}
			value = splittedCommandArgs[1]
			if value == "/$" {
				value = "$"
			}

			errMessage, isKeyFree := CheckKeyInFile(key)
			if !isKeyFree {
				return errMessage, false
			}

			checker := WriteKeyValueToFile(key, value)
			if !checker {
				return fmt.Sprintf("Unable to add [%s:%s] pair to the database", key, value), false
			}
			return fmt.Sprintf("[%s:%s] pair was added to the database", key, value), true
		}
		return "Incorrect amount of SET arguments", false
	}

	if runeCounter == 2 {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			if splittedCommandArgs[0] == "$" {
				return "KEY can not be empty", false
			}
			return "VALUE can not be empty", false
		}

		if splittedCommandArgs[0] == "$" {
			key = strings.Join(splittedCommandArgs[1:runeIndexes[1]], " ")
			for index := 1; index < len(key); index += 1 {
				if key[index] == '$' && key[index-1] == '/' {
					key = key[:index-1] + key[index:]
				}
			}
			if runeIndexes[1]+1 == len(splittedCommandArgs) {
				return "VALUE can not be empty", false
			}
			value = splittedCommandArgs[runeIndexes[1]+1]

			errMessage, isKeyFree := CheckKeyInFile(key)
			if !isKeyFree {
				return errMessage, false
			}

			checker := WriteKeyValueToFile(key, value)
			if !checker {
				return fmt.Sprintf("Unable to add [%s:%s] pair to the database", key, value), false
			}
			return fmt.Sprintf("[%s:%s] pair was added to the database", key, value), true
		}

		key = splittedCommandArgs[0]
		value = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(value); index += 1 {
			if value[index] == '$' && value[index-1] == '/' {
				value = value[:index-1] + value[index:]
			}
		}

		errMessage, isKeyFree := CheckKeyInFile(key)
		if !isKeyFree {
			return errMessage, false
		}

		checker := WriteKeyValueToFile(key, value)
		if !checker {
			return fmt.Sprintf("Unable to add [%s:%s] pair to the database", key, value), false
		}
		return fmt.Sprintf("[%s:%s] pair was added to the database", key, value), true
	}

	if runeCounter == 4 {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			return "KEY can not be empty", false
		}
		if runeIndexes[3]-runeIndexes[2] < 2 {
			return "VALUE can not be empty", false
		}

		key = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(key); index += 1 {
			if key[index] == '$' && key[index-1] == '/' {
				key = key[:index-1] + key[index:]
			}
		}
		value = strings.Join(splittedCommandArgs[runeIndexes[2]+1:runeIndexes[3]], " ")
		for index := 1; index < len(value); index += 1 {
			if value[index] == '$' && value[index-1] == '/' {
				value = value[:index-1] + value[index:]
			}
		}
	}

	errMessage, isKeyFree := CheckKeyInFile(key)
	if !isKeyFree {
		return errMessage, false
	}

	checker := WriteKeyValueToFile(key, value)
	if !checker {
		return fmt.Sprintf("Unable to add [%s:%s] pair to the database", key, value), false
	}
	return fmt.Sprintf("[%s:%s] pair was added to the database", key, value), true
}

func CheckKeyInFile(key string) (string, bool) {

	file, err := os.OpenFile("database.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Unable to 'open'/'create and open' the database: %s\n", fmt.Sprint(err))
		return "Unable to 'open'/'create and open' the database", false
	}
	isDatabaseInUse = true
	defer func() {
		isDatabaseInUse = false
	}()
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	keyValue := &KeyValue{}

OUTER:
	for fileScanner.Scan() {
		decoder := json.NewDecoder(strings.NewReader(fileScanner.Text()))

		err = decoder.Decode(&keyValue)
		if err != nil {
			fmt.Printf("Unable to decode JSON to [KEY:VALUE]: %s\n", fmt.Sprint(err))
			return "Unable to decode JSON to [KEY:VALUE]", false
		}

		for _, value := range keyArray {
			if value == keyValue.Key {
				continue OUTER
			}
		}

		if keyValue.Key == key {
			return fmt.Sprintf("[%s] key is already in the database", key), false
		}
	}
	return "", true
}

func CommandGET(splittedCommandArgs []string) (string, bool) {

	runeCounter := 0
	runeIndexes := []int{}
	key := ""

	for index, value := range splittedCommandArgs {
		if value == "$" {
			runeCounter += 1
			runeIndexes = append(runeIndexes, index)
		}
	}

	if runeCounter%2 != 0 || runeCounter > 2 {
		return "Incorrect amount of $ runes", false
	}

	if runeCounter == 0 {
		if len(splittedCommandArgs) != 1 {
			return "Incorrect amount of GET args", false
		}

		key = splittedCommandArgs[0]
		if key == "/$" {
			key = "$"
		}
	} else {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			return "KEY can not be empty", false
		}

		key = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(key); index += 1 {
			if key[index] == '$' && key[index-1] == '/' {
				key = key[:index-1] + key[index:]
			}
		}
	}

	file, err := os.OpenFile("database.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Unable to 'open'/'create and open' file: %s\n", fmt.Sprint(err))
		return "Unable to 'open'/'create and open' file", false
	}
	isDatabaseInUse = true
	defer func() {
		isDatabaseInUse = false
	}()
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	keyValue := &KeyValue{}

OUTER:
	for fileScanner.Scan() {
		decoder := json.NewDecoder(strings.NewReader(fileScanner.Text()))
		err = decoder.Decode(&keyValue)
		if err != nil {
			fmt.Printf("Unable to decode JSON to [KEY:VALUE]: %s\n", fmt.Sprint(err))
			return "Unable to decode JSON to [KEY:VALUE]", false
		}

		for _, value := range keyArray {
			if value == keyValue.Key {
				continue OUTER
			}
		}

		if keyValue.Key == key {
			return fmt.Sprintf("[KEY:%s][VALUE:%s]", keyValue.Key, keyValue.Value), true
		}
	}

	return fmt.Sprintf("[%s] key is not found", key), false
}

func CommandKEYS(splittedCommandArgs []string, conn net.Conn) (string, bool) {
	key := ""
	if splittedCommandArgs[0] == "*" {
		if len(splittedCommandArgs) == 1 {
			return "Invalid amount of KEYS arguments", false
		}
		key = strings.Join(splittedCommandArgs[1:], " ")
	} else {
		key = strings.Join(splittedCommandArgs, " ")
	}

	searchTrigger := false

	file, err := os.OpenFile("database.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("Unable to 'open'/'create and open' file: %s\n", fmt.Sprint(err))
		return "Unable to 'open'/'create and open' file", false
	}
	isDatabaseInUse = true
	defer func() {
		isDatabaseInUse = false
	}()
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	keyValue := &KeyValue{}

OUTER:
	for fileScanner.Scan() {
		decoder := json.NewDecoder(strings.NewReader(fileScanner.Text()))

		err = decoder.Decode(&keyValue)
		if err != nil {
			fmt.Printf("Unable to decode JSON to [KEY:VALUE]: %s\n", fmt.Sprint(err))
			return "Unable to decode JSON to [KEY:VALUE]", false
		}

		for _, value := range keyArray {
			if value == keyValue.Key {
				continue OUTER
			}
		}

		if strings.Contains(keyValue.Key, key) && splittedCommandArgs[0] != "*" {
			conn.Write([]byte(fmt.Sprintf("[%s]\n", keyValue.Key)))
			searchTrigger = true
		} else if strings.ContainsAny(keyValue.Key, key) && splittedCommandArgs[0] == "*" {
			conn.Write([]byte(fmt.Sprintf("[%s]\n", keyValue.Key)))
			searchTrigger = true
		}
	}

	if !searchTrigger {
		return fmt.Sprintf("KEYS[%s] are not found", key), false
	}
	return "", true
}

func CommandSUBSCRIBE(cleanSplittedMessage []string, conn net.Conn) (string, bool) {
	if len(cleanSplittedMessage) == 0 {
		return "Channel name can not be empty", false
	}

	channelName := strings.Join(cleanSplittedMessage, " ")

	commandResponse, isSubscribed := channels.Subscribe(channelName, conn)
	if !isSubscribed {
		channels.Unsubscribe(channelName, conn)
		return commandResponse, false
	}

	channels.Unsubscribe(channelName, conn)
	return commandResponse, true
}

func CommandPUBLISH(cleanSplittedMessage []string, conn net.Conn) (string, bool) {
	if len(cleanSplittedMessage) == 0 {
		return "Chat message can not be empty", false
	}

	message := strings.Join(cleanSplittedMessage, " ")
	messageResponse, isSent := channels.Publish(message, conn)
	if !isSent {
		return fmt.Sprint(messageResponse), false
	}
	return "", true
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

	portCheck := exec.Command("lsof", "-i", fmt.Sprintf(":%s", port))

	portCheckOutput, err := portCheck.CombinedOutput()
	if err != nil {
		if fmt.Sprint(err) != "exit status 1" {
			return fmt.Sprintf("%s %s\n", fmt.Sprint(err), string(portCheckOutput)), false
		}
	} else {
		if string(portCheckOutput) != "" {
			return fmt.Sprintf("Port %s is not available\n", port), false
		}
	}
	return fmt.Sprintf("Port %s [OK]\n", port), true
}

func CheckMode() (string, bool) {
	if storage != "disk" && storage != "memory" {
		return fmt.Sprintf("%s is not valid mode parameter!\n", storage), false
	} else {
		return fmt.Sprintf("Mode: '%s' [OK]\n", storage), true
	}
}

func SetDatabaseRefresher() {
	DatabaseRefresher := time.NewTicker(5 * time.Second)
	defer DatabaseRefresher.Stop()

	for _ = range DatabaseRefresher.C {
		if len(keyArray) > 0 {
			if isDatabaseInUse {
				fmt.Printf("Can not refresh databse: Database is in use\n")
			} else {
				RefreshDatabase()
			}
		}
	}
}

func CheckInputForWhitespaces(message string) ([]string, bool) {
	splittedMessage := strings.Split(message, " ")
	if (len(strings.TrimSpace(message))) == 0 {
		return nil, false
	}

	cleanSplittedMessage := []string{}
	for _, value := range splittedMessage {
		if len(strings.TrimSpace(value)) != 0 {
			cleanSplittedMessage = append(cleanSplittedMessage, value)
		}
	}
	return cleanSplittedMessage, true
}

func CommandSET_memory(splittedCommandArgs []string) (string, bool) {
	if len(splittedCommandArgs) < 2 {
		return "Incorrect amount of SET arguments", false
	}

	key := ""
	value := ""
	runeIndexes := []int{}

	runeCounter := 0
	for index, value := range splittedCommandArgs {
		if value == "$" {
			runeCounter += 1
			runeIndexes = append(runeIndexes, index)
		}
	}

	if runeCounter%2 != 0 || runeCounter > 4 {
		return "Incorrect amount of $ runes", false
	}

	if runeCounter == 0 {
		if len(splittedCommandArgs) == 2 {
			key = splittedCommandArgs[0]
			if key == "/$" {
				key = "$"
			}
			value = splittedCommandArgs[1]
			if value == "/$" {
				value = "$"
			}

			_, keyExist := keyValueMap[key]
			if keyExist {
				return fmt.Sprintf("[%s] key already exists", key), false
			}

			keyValueMap[key] = value

			return fmt.Sprintf("[%s:%s] pair was added to the storage", key, value), true
		}
		return "Incorrect amount of SET arguments", false
	}

	if runeCounter == 2 {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			if splittedCommandArgs[0] == "$" {
				return "KEY can not be empty", false
			}
			return "VALUE can not be empty", false
		}

		if splittedCommandArgs[0] == "$" {
			key = strings.Join(splittedCommandArgs[1:runeIndexes[1]], " ")
			for index := 1; index < len(key); index += 1 {
				if key[index] == '$' && key[index-1] == '/' {
					key = key[:index-1] + key[index:]
				}
			}
			if runeIndexes[1]+1 == len(splittedCommandArgs) {
				return "VALUE can not be empty", false
			}
			value = splittedCommandArgs[runeIndexes[1]+1]

			_, keyExist := keyValueMap[key]
			if keyExist {
				return fmt.Sprintf("[%s] key already exists"), false
			}

			keyValueMap[key] = value

			return fmt.Sprintf("[%s:%s] pair was added to the storage", key, value), true
		}

		key = splittedCommandArgs[0]
		value = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(value); index += 1 {
			if value[index] == '$' && value[index-1] == '/' {
				value = value[:index-1] + value[index:]
			}
		}

		_, keyExist := keyValueMap[key]
		if keyExist {
			return fmt.Sprintf("[%s] key already exists"), false
		}

		keyValueMap[key] = value

		return fmt.Sprintf("[%s:%s] pair was added to the storage", key, value), true
	}

	if runeCounter == 4 {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			return "KEY can not be empty", false
		}
		if runeIndexes[3]-runeIndexes[2] < 2 {
			return "VALUE can not be empty", false
		}

		key = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(key); index += 1 {
			if key[index] == '$' && key[index-1] == '/' {
				key = key[:index-1] + key[index:]
			}
		}
		value = strings.Join(splittedCommandArgs[runeIndexes[2]+1:runeIndexes[3]], " ")
		for index := 1; index < len(value); index += 1 {
			if value[index] == '$' && value[index-1] == '/' {
				value = value[:index-1] + value[index:]
			}
		}
	}

	_, keyExist := keyValueMap[key]
	if keyExist {
		return fmt.Sprintf("[%s] key already exists"), false
	}

	keyValueMap[key] = value

	return fmt.Sprintf("[%s:%s] pair was added to the storage", key, value), true
}

func CommandGET_memory(splittedCommandArgs []string) (string, bool) {
	runeCounter := 0
	runeIndexes := []int{}
	key := ""

	for index, value := range splittedCommandArgs {
		if value == "$" {
			runeCounter += 1
			runeIndexes = append(runeIndexes, index)
		}
	}

	if runeCounter%2 != 0 || runeCounter > 2 {
		return "Incorrect amount of $ runes", false
	}

	if runeCounter == 0 {
		if len(splittedCommandArgs) != 1 {
			return "Incorrect amount of GET args", false
		}

		key = splittedCommandArgs[0]
		if key == "/$" {
			key = "$"
		}
	} else {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			return "KEY can not be empty", false
		}

		key = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(key); index += 1 {
			if key[index] == '$' && key[index-1] == '/' {
				key = key[:index-1] + key[index:]
			}
		}
	}

	value, keyExist := keyValueMap[key]
	if keyExist {
		return fmt.Sprintf("[%s:%s]", key, value), true
	}

	return fmt.Sprintf("[%s] key is not found", key), false
}

func CommandDEL_memory(splittedCommandArgs []string) (string, bool) {
	runeCounter := 0
	runeIndexes := []int{}
	key := ""

	for index, value := range splittedCommandArgs {
		if value == "$" {
			runeCounter += 1
			runeIndexes = append(runeIndexes, index)
		}
	}

	if runeCounter%2 != 0 || runeCounter > 2 {
		return "Incorrect amount of $ runes", false
	}

	if runeCounter == 0 {
		if len(splittedCommandArgs) != 1 {
			return "Incorrect amount of DEL args", false
		}

		key = splittedCommandArgs[0]
		if key == "/$" {
			key = "$"
		}
	} else {
		if runeIndexes[1]-runeIndexes[0] < 2 {
			return "KEY can not be empty", false
		}

		key = strings.Join(splittedCommandArgs[runeIndexes[0]+1:runeIndexes[1]], " ")
		for index := 1; index < len(key); index += 1 {
			if key[index] == '$' && key[index-1] == '/' {
				key = key[:index-1] + key[index:]
			}
		}
	}

	_, keyExist := keyValueMap[key]
	if keyExist {
		delete(keyValueMap, key)
		return fmt.Sprintf("[%s] key was deleted from the storage", key), true
	}

	return fmt.Sprintf("[%s] key is not found", key), false
}

func CommandKEYS_memory(splittedCommandArgs []string, conn net.Conn) (string, bool) {
	key := ""
	searchTrigger := false
	if splittedCommandArgs[0] == "*" {
		if len(splittedCommandArgs) == 1 {
			return "Invalid amount of KEYS arguments", false
		}
		key = strings.Join(splittedCommandArgs[1:], " ")
	} else {
		key = strings.Join(splittedCommandArgs, " ")
	}

	for iKey, _ := range keyValueMap {
		if strings.Contains(iKey, key) && splittedCommandArgs[0] != "*" {
			conn.Write([]byte(fmt.Sprintf("[%s]\n", iKey)))
			searchTrigger = true
		} else if strings.ContainsAny(iKey, key) && splittedCommandArgs[0] == "*" {
			conn.Write([]byte(fmt.Sprintf("[%s]\n", iKey)))
			searchTrigger = true
		}
	}
	if searchTrigger {
		return "", true
	}
	return fmt.Sprintf("[%s] key is not found", key), false
}
