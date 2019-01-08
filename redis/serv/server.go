package serv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"redis/save"
	"redis/types"
	"strings"
)

//ServConnHandler - scanning input connection and send to ServCmndsHandler
func ServConnHandler(ServConnCh chan types.Server, conn net.Conn) {
	scnnr := bufio.NewScanner(conn)
	for scnnr.Scan() {
		line := scnnr.Text()
		inptFlds := strings.Fields(line)
		rslt := make(chan string)
		ServConnCh <- types.Server{
			HandFlds: inptFlds,
			Rslt:     rslt,
		}
		io.WriteString(conn, <-rslt)
	}
}

//ServCmndsHandler - containing GET, SET, DEL Commands.
func ServCmndsHandler(ServConnCh chan types.Server, gMemory string) {
	var memData = make(map[string]string)
	for cmnd := range ServConnCh {
		if len(cmnd.HandFlds) < 2 {
			cmnd.Rslt <- "SET 'key value', GET 'key', DEL 'key'.\n"
			continue
		}
		switch cmnd.HandFlds[0] {
		//GET <key>
		case "GET":
			if len(cmnd.HandFlds) != 2 {
				cmnd.Rslt <- "Get what?"
			}
			key := cmnd.HandFlds[1]
			value := memData[key]
			if len(memData) == 0 {
				cmnd.Rslt <- "Data is empty"
			} else {
				cmnd.Rslt <- value
			}
			//SET <key>
		case "SET":
			if len(cmnd.HandFlds) != 3 {
				cmnd.Rslt <- "Missing value\n"
			}
			key := cmnd.HandFlds[1]
			value := cmnd.HandFlds[2]
			memData[key] = value
			if gMemory == "disk" {
				memDisk, err := json.Marshal(memData)
				if err != nil {
					fmt.Println("JSON", string(memDisk), err)
				}
				info := string(memDisk)
				save.SaveOnDisk(info)
				cmnd.Rslt <- "JSON: KEY - VALUE SET\n"
			} else {
				cmnd.Rslt <- "KEY - VALUE SET\n"
			}
			//DEL <KEY>
		case "DEL":
			key := cmnd.HandFlds[1]
			value, ok := memData[key]
			if ok {
				delete(memData, key)
				cmnd.Rslt <- key + " - " + value + " DELETED\n"
			} else {
				cmnd.Rslt <- "KEY not found\n"
			}
		default:
			cmnd.Rslt <- "I don't know this command :" + cmnd.HandFlds[0] + "\n"
		}
	}
}
