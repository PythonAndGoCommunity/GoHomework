package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

func (object storage) GET(key string, response chan string) {
	value, ok := object.storageMap[key]
	if !ok && object.storageMode == "disk" {
		valueInFile := object.readFromFile(key)
		if valueInFile != "false" {
			object.storageMap[key] = valueInFile
			response <- key + ":" + valueInFile //send data to client
		} else {
			response <- "no pair contains key=" + key //send data to client
		}
	} else {
		response <- key + ":" + value //send data to client
	}
}

func (object storage) SET(key string, value string, response chan string) {
	val, ok := object.storageMap[key]
	if ok == true {
		response <- fmt.Sprintf("store contains pair with key %v: %v", key, val) //send data to client
	} else {
		inFileValue := object.readFromFile(key)
		if inFileValue != "false" {
			object.storageMap[key] = inFileValue
			response <- "pair with key " + key + "exist. " + key + ": " + inFileValue //send data to client
			return
		} else {
			object.storageMap[key] = value
			response <- "pair " + key + ":" + value + " created" //send data to client
			return
		}
	}
}

func (object storage) DEL(key string, response chan string) {
	_, ok := object.storageMap[key]
	if ok == true {
		delete(object.storageMap, key)
		object.deleteFromFile(key)
		response <- "pair deleted" //send data to client
	} else {
		deleteFromFileResult := object.deleteFromFile(key)
		if deleteFromFileResult == true {
			response <- "pair deleted" //send data to client
		} else {
			response <- "no pair for delete" //send data to client
		}
	}
}

func (object storage) KEYS(pattern string, response chan string) {
	result := make(map[string]string)
	for key, value := range object.storageMap {
		keyExist, err := regexp.MatchString(pattern, key)
		if err != nil {
			result["error"] = "true"
		} else if keyExist {
			result[key] = value
		}
	}
	var resultToString string
	if len(result) == 0 {
		result["None"] = "None"
	}
	for key, value := range result {
		resultToString += fmt.Sprintf("%v:%v\n", key, value)
	}
	response <- resultToString
}

func (object storage) WRITE() {
	object.writeToFileAllData()
}

func (object storage) readFromFile(key string) string {
	object.storageFile.Seek(0, 0)
	scaner := bufio.NewScanner(object.storageFile)
	for scaner.Scan() {
		line := scaner.Text()
		splitLine := strings.Split(line, ";")
		//splitLine[0] is key; splitLine[1] is val
		if splitLine[0] == key {
			return splitLine[1]
		}
	}
	return "false"
}

func (object storage) readFromFileAllData() map[string]string {
	object.storageFile.Seek(0, 0)
	scaner := bufio.NewScanner(object.storageFile)
	result := make(map[string]string)
	for scaner.Scan() {
		line := scaner.Text()
		splitLine := strings.Split(line, ";")
		//splitLine[0] is key; splitLine[1] is val
		result[splitLine[0]] = splitLine[1]
	}
	return result
}

func (object storage) deleteFromFile(key string) bool {
	object.storageFile.Seek(0, 0)
	scaner := bufio.NewScanner(object.storageFile)
	var stringResult string
	var deleteBool = false
	for scaner.Scan() {
		line := strings.Split(scaner.Text(), ";")
		if line[0] == key {
			stringResult += ""
			deleteBool = true
		} else {
			stringResult += line[0] + ";" + line[1] + "\n"
		}
	}
	if deleteBool {
		object.storageFile.Seek(0, 0)
		ioutil.WriteFile(string(object.storageFile.Name()), []byte(stringResult), 0755)
		return deleteBool
	} else {
		return deleteBool
	}

}

func (object storage) writeToFileAllData() {
	//object.storageFile.Seek(0,0)
	allDataInFile := object.readFromFileAllData()
	var resultString string
	resultMap := object.storageMap
	for key, fileValue := range allDataInFile {
		_, ok := resultMap[key]
		if !ok {
			resultMap[key] = fileValue
		}
	}
	for key, value := range resultMap {
		resultString += fmt.Sprintf("%v;%v\n", key, value)
	}
	ioutil.WriteFile(string(object.storageFile.Name()), []byte(resultString), 0755)
	//object.storageFile.Sync()
}

type storage struct {
	storageMap  map[string]string
	storageMode string
	storageFile *os.File
}

func main() {
	var port string
	flag.StringVar(&port, "port", "9090", "listening port")
	flag.StringVar(&port, "p", "9090", "listening port")
	var host string
	flag.StringVar(&host, "h", "127.0.0.1", "listening IP")
	flag.StringVar(&host, "host", "127.0.0.1", "listening IP")
	var mode string
	flag.StringVar(&mode, "m", "disk", "storage mode disk(default) or memory")
	flag.StringVar(&mode, "mode", "disk", "storage mode disk(default) or memory")

	flag.Parse()
	fmt.Println(port, host)

	logFile, _ := os.OpenFile("server.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	logWriter := log.New(logFile, "", log.Ldate|log.Ltime)

	listener, err := net.Listen("tcp", host+":"+port)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	logWriter.Println("start server at address: " + host + ":" + port)

	request := make(chan string)
	response := make(chan string)
	go goStorage(request, response, mode)

	for {
		conn, err := listener.Accept()
		logWriter.Printf("received connection from %v", conn.RemoteAddr())
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn, request, response, logWriter) // запускаем горутину для обработки запроса

	}

}

func handleConnection(conn net.Conn, request chan string, response chan string, logWriter *log.Logger) {
	defer conn.Close()
	for {
		// read data from request
		input := make([]byte, 1024*32)
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			log.Printf("received message read error")
			break
		}
		log.Printf("received message from client(%v) %v", conn.RemoteAddr(), string(input[0:n]))
		logWriter.Printf("received message from client(%v) '%v'", conn.RemoteAddr(), string(input[0:n]))
		request <- string(input[0:n])
		resp := []byte(<-response)
		log.Printf("send response to client(%v) %v", conn.RemoteAddr(), string(resp))
		logWriter.Printf("send response to client(%v) %v", conn.RemoteAddr(), string(resp))
		conn.Write(resp)
	}
}

func goStorage(request chan string, response chan string, mode string) {
	file, _ := os.OpenFile("storage", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0775)
	defer file.Close()
	newStorage := storage{
		make(map[string]string),
		mode,
		file,
	}
	requestsCount := 0
	for entry := range request {
		entryList := strings.Split(entry, " ")
		keyword := strings.ToUpper(entryList[0])
		switch keyword {
		case "GET":
			if len(entryList) == 2 {
				key := string(entryList[1])
				newStorage.GET(key, response)
			} else {
				response <- "Wrong operator count" //send data to client
			}
		case "SET":
			if len(entryList) == 3 {
				key := string(entryList[1])
				value := string(entryList[2])
				newStorage.SET(key, value, response)
			} else {
				response <- "Wrong operator count" //send data to client
			}
		case "DEL":
			if len(entryList) == 2 {
				key := string(entryList[1])
				newStorage.DEL(key, response)
			} else {
				response <- "Wrong operator count" //send data to client
			}
		case "KEYS":
			if len(entryList) == 2 {
				pattern := string(entryList[1])
				pattern = strings.Replace(pattern, "*", "", -1)
				newStorage.KEYS(pattern, response)
			} else if len(entryList) == 1 {
				pattern := ".*"
				newStorage.KEYS(pattern, response)
			} else {
				response <- "keys are not exist" //send data to client
			}
		default:
			response <- "Wrong expression:(" //send data to client
		}
		requestsCount += 1
		if requestsCount >= 3 {
			requestsCount = 0
			newStorage.WRITE()
		}
	}
}
