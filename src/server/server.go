package server

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	cmdSet    = "SET"
	cmdGet    = "GET"
	cmdDel    = "DEL"
	cmdKeys   = "KEYS"
	publish   = "publish"
	subscribe = "subscribe"
	dump      = "dump"
	restore   = "restore"
	dataJSON  = "data.json"

	defaultPort = "9090"
	diskMode    = "disk"
	memoryMode  = "memory"

	fMode    = "--mode"
	mode     = "-m"
	fPort    = "--port"
	port     = "-p"
	fVerbose = "--verbose"
	verbose  = "-v"

	protocolTCP = "tcp"
)

var (
	_mainPort    = defaultPort
	_mainMode    = true  //mode=true - save in disk; mode=false - save in memory
	_mainVerbose = false //verbose=false - without logging; verbose=true - with logging

	dataMap = map[string]string{}

	clients       *list.List
	publishAddr   = list.New()
	subscribeAddr = list.New()
)

//for parsing CLI
func trimLastSymbol(s string) string {
	if last := len(s) - 1; last >= 0 && s[last] == ',' {
		s = s[:last]
		//Info.Printf("Trimmed arguments: %s", s)
		return s
	}
	return s
}

func parseArguments() {

	Info.Println("Parse arguments")
	var cmds = make(map[string]func(string))
	// заполняем команды
	cmds[port] = func(s string) {
		//номер порта
		s = trimLastSymbol(s)
		Info.Printf("Parse -p=%s\n", s)
		//update client/server port

		//проверка на число
		if port, err := strconv.Atoi(s); err == nil {
			fmt.Printf("%q looks like a number.\n", s)
			//Valid numbers for ports are: 0 to 2^16-1 = 0 to 65535
			//But user ports 1024 to 49151
			if port < 49152 && port > 1024 {
				_mainPort = s
			}
		} else {
			_mainPort = defaultPort
		}
		Info.Printf("New port: %s\r\n", _mainPort)
	}

	cmds[fPort] = func(s string) {

		s = trimLastSymbol(s)
		Info.Printf("Parse --port=%s\r\n", s)
		//update client/server port

		//проверка на число
		if port, err := strconv.Atoi(s); err == nil {
			//Valid numbers for ports are: 0 to 2^16-1 = 0 to 65535
			//But user ports 1024 to 49151
			if port < 49152 && port > 1024 {
				_mainPort = s
			}
		} else {
			_mainPort = defaultPort
		}
		Info.Printf("New port: %s\r\n", _mainPort)
	}
	cmds[mode] = func(s string) {
		s = trimLastSymbol(s)
		Info.Printf("Parse -m=%s\r\n", s)
		if s == memoryMode {
			_mainMode = false
		} else if s == diskMode {
			_mainMode = true
		} else {
			Warning.Println("Argument 'mode' error, will be used 'disk'")
			_mainMode = true
		}

		if _mainMode == true {
			Info.Printf("New mode work with data: %s", diskMode)
		} else {
			Info.Printf("New mode work with data: %s", memoryMode)
		}
	}
	cmds[fMode] = func(s string) {
		s = trimLastSymbol(s)
		Info.Printf("Parse --mode=%s\r\n", s)
		if s == memoryMode {
			_mainMode = false
		} else if s == diskMode {
			_mainMode = true
		} else {
			Warning.Println("Argument 'mode' error, will be used 'disk'")
			_mainMode = true
		}
		if _mainMode == true {
			Info.Printf("New mode work with data: %s", diskMode)
		} else {
			Info.Printf("New mode work with data: %s", memoryMode)
		}
	}
	cmds[verbose] = func(s string) {
		Info.Println("Parse -v")
		_mainVerbose = true
		Info.Println("Work with logs")
	}
	cmds[fVerbose] = func(s string) {
		Info.Println("Parse --verbose")
		_mainVerbose = true
		Info.Println("Work with all logs")
	}

	//анализ команд
	for _, arg := range os.Args[1:] {
		cmd := strings.Split(arg, "=")

		if len(cmd) > 2 {
			fmt.Println(arg, "don't know... ")
			continue
		}

		name := cmd[0]
		param := ""
		if trimLastSymbol(cmd[0]) == verbose || trimLastSymbol(cmd[0]) == fVerbose {
			name = trimLastSymbol(cmd[0])
			param = ""
		} else {
			param = cmd[1]
		}

		// ищем функцию
		fn := cmds[name]
		if fn == nil {
			fmt.Println(name, "i don't know... ")
			continue
		}

		// исполняем команду
		fn(param)
	}
}

// JSONWriter structure for writing to a file in JSON format
type JSONWriter struct {
	mutex    *sync.Mutex
	fileName string
}

func newJSONWriter(fileName string) *JSONWriter {
	name := fileName
	return &JSONWriter{mutex: &sync.Mutex{}, fileName: name}
}

/**
if
f = True - SET
f = False - DEL
k - key in map
v - value in map
*/
func (w *JSONWriter) update(k string, v string, f bool) {

	w.mutex.Lock()

	if f == true {
		dataMap[k] = v
	} else {
		delete(dataMap, k)
	}

	Info.Println(dataMap)

	if _mainMode {
		//преобразуем в JSON формат
		jsonString, err := json.Marshal(dataMap)
		//Info.Println(jsonString)
		if err != nil {
			Error.Println(err)
		} else {
			//Info.Println(jsonString)
			err = ioutil.WriteFile(w.fileName, jsonString, 0644)
			if err != nil {
				Error.Println(err)
			}
		}
	}
	w.mutex.Unlock()
}

func (w *JSONWriter) get() map[string]string {

	w.mutex.Lock()

	// Open jsonFile
	jsonFile, err := os.Open(w.fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		Warning.Println(err)
	} else {
		Info.Println("Successfully Opened json file")
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	err = json.Unmarshal([]byte(byteValue), &result)

	w.mutex.Unlock()

	if err != nil {
		Warning.Println(err)
		result = map[string]string{}
	} else {
		Info.Printf("JSON file: %s", result)
	}
	return result

}

var (
	//Trace for minor doing
	Trace   *log.Logger
	//Info for user doing
	Info    *log.Logger
	//Warning for case minor error
	Warning *log.Logger
	//Error for errors
	Error   *log.Logger
)

func initLoggers(

	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE:  ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO:    ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR:   ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// Println - only print string
func Println(l *log.Logger, s string) {

	if _mainVerbose {
		l.Println(s)
	}
}

// Printf - print string with 1 argument
func Printf(l *log.Logger, s string, s2 interface{}) {
	if _mainVerbose {
		l.Printf(s, s2)
	}
}

func handleClient(socket net.Conn, writer JSONWriter) {
	for {
		buffer, err := bufio.NewReader(socket).ReadString('\n')
		if err != nil {
			Printf(Info, "User %s go away", socket.RemoteAddr())
			//fmt.Println("User go away")

			if subscribeAddr.Len() > 0 {
				for j := subscribeAddr.Front(); j != nil; j = j.Next() {
					if j.Value == socket.RemoteAddr().String() {
						subscribeAddr.Remove(j)
					}
				}
				//Printf(Info, "After s:%s\r\n", subscribeAddr)
				//fmt.Println(subscribeAddr.Len())
			}

			if publishAddr.Len() > 0 {
				for j := publishAddr.Front(); j != nil; j = j.Next() {
					if j.Value == socket.RemoteAddr().String() {
						publishAddr.Remove(j)
					}
				}
				//Printf(Info, "After p:%s\r\n", publishAddr)
				//fmt.Println(publishAddr.Len())
			}

			socket.Close()
			return

		}
		for i := clients.Front(); i != nil; i = i.Next() {
			//fmt.Fprint(i.Value.(net.Conn), buffer)

			//обработка пришедшей команды
			go parseInputMessage(writer, i.Value.(net.Conn), buffer)
		}
	}
}

//This will be fixed in Go 1.12.))))))
func parseInputMessage(jsonW JSONWriter, conn net.Conn, s string) {
	Printf(Info, "input message: %s", s)

	s = strings.Trim(s, "\r\n")
	line := strings.Split(s, " ")

	cmd := line[0]
	cmd = strings.ToUpper(cmd)
	switch cmd {
	case cmdSet:
		{
			var error = false
			Printf(Info, "command cmdSet: %s", cmd)
			var key, value = "", ""
			if len(line) > 1 {
				key = line[1]
				if len(line) > 2 {
					value = strings.Trim(line[2], "\r\n")
				}
				if len(line) > 3 {
					//TODO сообщить об ошибке
					error = true
				}

			} else {
				key, value = "", ""
			}
			Info.Printf("\nk:%s v:%s", key, value)

			if error {
				Println(Error, "ERROR: Syntax error\r\n")
				//отправить данные клиенту
				fmt.Fprint(conn, "ERROR: Syntax error\r\n")
			} else {
				//запись в JSON файл
				jsonW.update(key, value, true)
				Println(Trace, "Record was saved\r\n")
				//оправить сообщение пользователю
				fmt.Fprint(conn, "OK\r\n")
			}
		}
	case cmdGet:
		{
			Printf(Info, "command cmdGet:%s", cmd)
			if len(line) == 2 {
				key := strings.Trim(string(line[1]), "\r\n")
				Printf(Info, "key: %s", key)
				_, ok := dataMap[key]
				if ok {
					//проверка c чтением из файла
					if _mainMode {
						value := jsonW.get()[key]
						Printf(Trace, "Read value from file: %s\r\n", value)
						//отправить клиенту
						fmt.Fprint(conn, value+"\r\n")
					} else {
						//если ключ есть
						value := dataMap[key]
						Printf(Trace, "value:%s\r\n", value)
						//отправить клиенту
						fmt.Fprint(conn, value+"\r\n")
					}

				} else {
					//если ключа нет
					fmt.Printf("\nvalue:(nil)")
					//отправить клиенту
					fmt.Fprint(conn, "(nil)\r\n")
				}
			} else {
				//TODO сообщить об ошибке
				fmt.Fprint(conn, "ERROR: wrong number of arguments (given "+string(len(line)-1)+" expected 1)\r\n")
			}
		}
	case cmdDel:
		{
			//удаление по ключу
			Printf(Info, "command cmdDel:%s", cmd)
			if len(line) > 1 {
				key := strings.Trim(string(line[1]), "\r\n")
				_, ok := dataMap[key]
				if ok {

					//update MAP and write to JSON file, because update map
					jsonW.update(key, "nil", false)

					//TODO добавить удаление многих элементов
					//отправить клиенту сколько удалено
					fmt.Fprint(conn, "1\r\n")
				} else {
					//отправить клиенту сколько удалено
					fmt.Fprint(conn, "0\r\n")
				}
				Printf(Info, "key:%s", key)
				Printf(Info, "key in map:%t", ok)

			} else {
				//TODO сообщить об ошибке
				//отправить клиенту
				fmt.Fprint(conn, "ERROR: wrong number of arguments for 'del' command\r\n")
			}
		}
	case cmdKeys:
		{
			Printf(Info, "command cmdKeys:%s\r\n", cmd)

			if len(line) > 1 {
				pattern := strings.Trim(string(line[1]), "\r\n")
				Printf(Info, "pattern keys: %s", pattern)

				keys := make([]string, 0, len(dataMap))
				i := 0
				for key := range dataMap {
					matched, err := regexp.MatchString(pattern, key)
					if matched && err == nil {
						Info.Printf("key: %s, pattern: %s", key, pattern)
						keys = append(keys, key)
						i++
					}
				}

				Printf(Info, "all found keys: %s", keys)
				if i > 0 {
					fmt.Fprint(conn, strings.Join(keys, ",")+"\r\n")
				} else {
					fmt.Fprint(conn, "(nil)\r\n")
				}

			} else {
				//TODO сообщить об ошибке
				//send client
				fmt.Fprint(conn, "ERROR: wrong number of arguments for 'keys' command\r\n")
			}

		}
	case publish:
		{

			if len(line) > 1 {
				channel := strings.Trim(string(line[1]), "\r\n")
				Printf(Info, "Publish channel: %s", channel)
				Info.Printf("Command: %s, address of client: %s\r\n", cmd, conn.RemoteAddr().String())

				publishAddr.PushFront(conn.RemoteAddr())
				if subscribeAddr.Len() > 0 {
					for j := subscribeAddr.Front(); j != nil; j = j.Next() {
						if j.Value == conn.RemoteAddr().String() {
							subscribeAddr.Remove(j)
						}
					}
					//Info.Printf("Address remove s:%s, rest len: %s\r\n", subscribeAddr, subscribeAddr.Len())
				}
			} else {
				//TODO сообщить об ошибке
				//отправить клиенту
				fmt.Fprint(conn, "ERROR: wrong number of arguments for 'publish' command\r\n")
			}

		}
	case subscribe:
		{
			Info.Printf("Command: %s, address of client: %s\r\n", cmd, conn.RemoteAddr().String())
			subscribeAddr.PushFront(conn.RemoteAddr())
			if publishAddr.Len() > 0 {
				for j := publishAddr.Front(); j != nil; j = j.Next() {
					if j.Value == conn.RemoteAddr().String() {
						publishAddr.Remove(j)
					}
				}
				Printf(Info, "After p:%s\r\n", publishAddr)
				fmt.Println(publishAddr.Len())
			}
		}
	case dump:
		{
			Info.Printf("command cmdKeys:%s", cmd)
			file, errF := os.Open(dataJSON) // For read access.
			if errF != nil {
				Error.Println("Unable to open file")
			}
			defer file.Close() // make sure to close the file even if we panic.
			n, error := io.Copy(conn, file)
			if error != nil {
				Error.Printf("Send file error %s\r\n", error.Error())
			}
			Info.Println(n, "Bytes sent")
			file.Close()
			conn.Close()
		}
	case restore:
		{
			Info.Printf("command cmdKeys:%s", cmd)

			data := line[1]
			var result map[string]string
			errUnm := json.Unmarshal([]byte(data), &result)

			if errUnm != nil {
				Warning.Println(errUnm)
				result = map[string]string{}
			} else {
				Info.Printf("JSON file: %s", result)
			}
			for k, v := range result {
				jsonW.update(k, v, true)
			}
		}
	default:
		{
			//отправить клиенту
			fmt.Fprint(conn, "\r\n")
		}

	}
}

func main() {

	initLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr) //init loggers

	parseArguments() //parse Args

	Println(Info, "Server start")
	clients = list.New()

	jsonWriter := newJSONWriter(dataJSON)
	dataMap = jsonWriter.get()

	server, err := net.Listen(protocolTCP, ":"+_mainPort)
	if err != nil {
		Printf(Error, "Error: %s", err.Error())
		return
	}
	defer server.Close() //close connection
	for {
		client, err := server.Accept()
		if err != nil {
			Printf(Error, "Error: %s", err.Error())
			return
		}

		Printf(Info, "New user connected %s\n", client.RemoteAddr())
		clients.PushBack(client)
		go handleClient(client, *jsonWriter)
	}

}

