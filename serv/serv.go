package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)


//MAX key and value len
const keyBuffLen = 63 //10
const valueBuffLen = 63 //12

//Reading commands from terminal
func startup() (port string, disk bool, filePath string, err error) {
	//default values
	port = ":9090"
	disk = false
	//filePath = "/data/redisDatabase"
	filePath = "redisDatabase"
	//filePath = "text"
	err = nil

	args := os.Args
	for i := 1; i < len(args); i += 2 {
		switch args[i] {
		// -p, --port	: set the port for listening on
		case "-p", "--port":
			if args[i+1][0] == ':' {
				port = args[i+1]
			} else {
				port = ":" + args[i+1]
			}
		// -m, --mode	: enable mirroring data to the drive
		case "-m", "--mode":
			if args[i+1] == "disk" {
				disk = true
			} else {
				if args[i+1] == "memory" {
					disk = false
				}
			}
		default:
			err = errors.New("ERROR: Unknown command")
			return
		}
	}
	return
}

//Using RAM as storage
func redis(clientCh clientChan) {
	fmt.Println("Redis started\nRAM mode")
	//Redis RAM storage
	var data = make(map[string]string)
Loop:
	for cmd := range clientCh.input {
		cmdList := strings.Fields(cmd)
		//check client's input
		if !checkInput(&clientCh, &cmdList){
			continue
		}

		switch cmdList[0] {
		// SET <KEY> <VALUE>
		case "SET":
			value, err := setRAM(&cmdList, &data)
			go func() {
				clientCh.output <- value
				clientCh.err <- err
			}()
		// GET <KEY>
		case "GET":
			value, err := getRAM(&cmdList, &data)
			go func() {
				clientCh.output <- value
				clientCh.err <- err
			}()
		// DEL <KEY>
		case "DEL":
			value, err := delRAM(&cmdList, &data)
			go func() {
				clientCh.output <- value
				clientCh.err <- err
			}()
		//exit
		case "exit":
			go func() {
				clientCh.output <- ""
				err := errors.New("ERROR")
				clientCh.err <- err
			}()
			break Loop
		default:
			// sending error
			go func() {
				clientCh.output <- ""
				err := errors.New("ERROR: Wrong command\nUse:\n SET <KEY> <VALUE>\n GET <KEY>\n DEL <KEY>")
				clientCh.err <- err
			}()
		}

	}
	fmt.Println("STOP")
}

// SET <KEY> <VALUE>
func setRAM(cmdList *[]string, data *(map[string]string))(value string, err error){

	if len(*cmdList) != 3 {
		err := errors.New("ERROR: Wrong command\nSET <KEY> <VALUE>")
		return "", err
	}
	(*data)[(*cmdList)[1]] = (*cmdList)[2]
	return (*cmdList)[2], nil
}

// GET <KEY>
func getRAM(cmdList *[]string, data *(map[string]string))(value string, err error){
	if len(*cmdList) != 2 {
		err := errors.New("ERROR: Wrong command\nGET <KEY>")
		return "", err
	}
	value, ok := (*data)[(*cmdList)[1]]
	if ok {
		return value, nil
	}
	err = errors.New("ERROR: Unknown key \"" + (*cmdList)[1] + "\"")
	return "", err
}

// DEL <KEY>
func delRAM(cmdList *[]string, data *(map[string]string))(value string, err error){
	if len(*cmdList) != 2 {
		err := errors.New("ERROR: Wrong command\nGET <KEY>")
		return "", err
	}
	value, ok := (*data)[(*cmdList)[1]]
	if ok {
		delete(*data, (*cmdList)[1])
		return value, nil
	}
	err = errors.New("ERROR: Unknown key \"" + (*cmdList)[1] + "\"")
	return "", err
}

//Check client's input
func checkInput(clientCh *clientChan, cmdList *[]string)(bool){
	if len(*cmdList) < 1 || len(*cmdList) > 3{
		go func() {
			clientCh.output <- ""
			err := errors.New("ERROR: Wrong command\nUse:\n SET <KEY> <VALUE>\n GET <KEY>\n DEL <KEY>")
			clientCh.err <- err
		}()
		return false
	}
	if len(*cmdList) == 2 && len((*cmdList)[1]) > keyBuffLen {
		go func() {
			clientCh.output <- ""
			err := errors.New("ERROR: The key is too long")
			clientCh.err <- err
		}()
		return false
	}
	if len(*cmdList) == 3 && (len((*cmdList)[1]) > keyBuffLen || len((*cmdList)[2]) > valueBuffLen){
		go func() {
			clientCh.output <- ""
			err := errors.New("ERROR: Key or value is too long")
			clientCh.err <- err
		}()
		return false
	}
	return true
}


//Connection to the client
func handle(conn net.Conn, clientCh clientChan) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
Loop:
	for scanner.Scan() {
		command := scanner.Text()
		clientCh.input <- command

		answer := <-clientCh.output
		err := <-clientCh.err
		switch err{
		case nil:
			answer += "\r\n"
		case errors.New("exit"):
			break Loop
		default:
			answer = err.Error() + "\r\n"
		}
		io.WriteString(conn, answer)
	}
}

//Clients <-> Redis channels
type clientChan struct {
	input  chan string
	output chan string
	err    chan error
}

func main() {
	port, disk, filePath, cmdErr := startup()
	if cmdErr == nil {
		li, netErr := net.Listen("tcp", port)
		if netErr != nil {
			log.Fatalln(netErr)
		} else {
			defer li.Close()

			//Prepare channels for client commands and saving data
			clInput := make(chan string)
			clOutput := make(chan string)
			clErr := make(chan error)
			clientCh := clientChan {
				input:	clInput,
				output:	clOutput,
				err:	clErr,
			}

			if disk {
				//Storage - disk
				go saveData(filePath, clientCh)
			} else {
				//Storage - RAM
				go redis(clientCh)
			}

			for {
				conn, err := li.Accept()
				if err != nil {
					log.Fatalln(err)
				}
				go handle(conn, clientCh)
			}
		}
	} else {
		//Some CMD error
		fmt.Println(cmdErr)
	}
	fmt.Println("stop")
}
