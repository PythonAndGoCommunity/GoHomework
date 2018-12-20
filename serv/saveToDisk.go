package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)


//Using disk as storage
func saveData(filePath string, clientCh clientChan) {
	//Redis disk storage
	file, fileErr := os.Create(filePath)
	//file, fileErr := os.OpenFile(filePath, os.O_RDWR, 0600)
	if fileErr != nil{
		panic(fileErr)
	}
	defer file.Close()
	format := "%" + strconv.Itoa(keyBuffLen) + "s %" + strconv.Itoa(valueBuffLen) + "s\n"

	fmt.Println("Redis started\nDisk mode")
	for cmd := range clientCh.input {
		cmdList := strings.Fields(cmd)
		//check client's input
		if !checkInput(&clientCh, &cmdList){
			continue
		}

		switch cmdList[0] {
		// SET <KEY> <VALUE>
		case "SET":
			value, err := setDisk(file, &cmdList, &format)
			go func() {
				clientCh.output <- value
				clientCh.err <- err
			}()
		// GET <KEY>
		case "GET":
			value, err := getDisk(file, &cmdList)
			go func() {
				clientCh.output <- value
				clientCh.err <- err
			}()
		// DEL <KEY>
		case "DEL":
			value, err := delDisk(file, &cmdList, &format)
			go func() {
				clientCh.output <- value
				clientCh.err <- err
			}()
		//exit
		case "stop":
			if cmdList[1] == "redis" {
				break
			}
		default:
			go func() {
				clientCh.output <- ""
				err := errors.New("ERROR: Unknown command")
				clientCh.err <- err
			}()
		}
	}
}


func setDisk(file *os.File, cmdList *[]string, format *string)(value string, err error){
	if len(*cmdList) != 3 {
		err := errors.New("ERROR: Wrong command\nSET <KEY> <VALUE>")
		return "", err
	}
	key := (*cmdList)[1]
	value = (*cmdList)[2]
	var offset int64
	_, _, offset, err = findKey(file, &key)
	if err != nil {
		return "", err
	}

	str := fmt.Sprintf(*format, key, value)
	buff := []byte(str)
	_, err = file.WriteAt(buff, offset)

	if err != nil {
		return "", err
	}
	return value, nil
}

func getDisk(file *os.File, cmdList *[]string)(value string, err error){
	if len(*cmdList) != 2 {
		err := errors.New("ERROR: Wrong command\nGET <KEY>")
		return "", err
	}
	key := (*cmdList)[1]
	ok, value, _, err := findKey(file, &key)
	if err != nil {
		return "", err
	}
	if ok {
		//Returning the value
		return value, nil
	}
	//There are no such key in db
	err = errors.New("ERROR: Unknown key \"" + key + "\"")
	return "", err
}

func delDisk(file *os.File, cmdList *[]string, format *string)(value string, err error){
	if len(*cmdList) != 2 {
		err := errors.New("ERROR: Wrong command\nGET <KEY>")
		return "", err
	}
	key := (*cmdList)[1]
	ok, value, offset, err := findKey(file, &key)
	if err != nil {
		return"", err
	}
	if ok {
		//Replacing the key and the value
		str := fmt.Sprintf(*format, "", "")
		buff := []byte(str)
		_, err := file.WriteAt(buff, offset)
		return value, err
	}
	//There are no such key in db
	err = errors.New("ERROR: Unknown key \"" + key + "\"")
	return"", err

}

func findKey(file *os.File, key *string)(ok bool, value string, offset int64, err error){
	reader := bufio.NewReader(file)
	for {
		file.Seek(offset,0)
		str, err := reader.ReadString('\n')
		if err != nil{
			if err == io.EOF || str == "\n"{
				err = nil
				return ok, value, offset, err
			}
			return ok, value, offset, err
		}
		words := strings.Fields(str)
		//Skipping empty strings
		if len(words) != 2{
			offset += int64(len(str))
			continue
		}
		if words[0] == *key {
			value = words[1]
			return true, value, offset, err
		}
		offset += int64(len(str))
	}
}
