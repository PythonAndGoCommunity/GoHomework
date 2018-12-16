package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type DataType map[string]string
var MemoryMode bool

func LoadFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR: cannot load data from " + filename)
		return data, err
	}
	return data, nil
}

func LoadData(filename string) error {
	dumped, err := LoadFromFile(filename)
	if err != nil {
		return err
	}
	Data = unDumpData(dumped)
	return nil
}

func SaveData(filename string) {
	if !MemoryMode {
		dumped := dumpData(Data)
		err := ioutil.WriteFile(filename, dumped, 0644)
		if err != nil {
			fmt.Println("ERROR: cannot save data to filename")
		}
	}
}

func AddEntry(key string, value string) {
	Data[key] = value
}

func CheckEntry(key string, value string) string {
	if v, ok := Data[key]; ok && strings.Compare(v, value) == 0 {
		return "OK"
	} else {
		return "(nil)"
	}
}

func GetEntry(key string) string {
	if value, ok := Data[key]; ok {
		return value
	}
	return "(nil)"
}

func RemoveEntries(keys []string) string {
	deleted := 0
	for i := range keys {
		_, ok := Data[keys[i]]
		if ok {
			deleted += 1
			delete(Data, keys[i])
		}
	}
	response := "(integer) " + strconv.Itoa(deleted)
	return response
}

func ShowAllKeys() string {
	response := ""
	for k := range Data {
		response = "\"" + k + "\" " + response
	}
	return response
}

func FindKeys(pattern string) string {
	response := ""
	for k := range Data {
		matched, err := regexp.MatchString(pattern, k)
		if err != nil {
			response = "ERROR: wrong pattern (" + pattern + ")."
			return response
		}
		if matched {
			response = "\"" + k + "\" " + response
		}
	}
	return response
}

func unDumpData(dumped []byte) DataType {
	var data DataType
	err := json.Unmarshal(dumped, &data)
	if err != nil {
		fmt.Println("ERROR: cannot undump data.")
	}
	return data
}

func dumpData(data DataType) []byte {
	dumped, err := json.Marshal(data)
	if err != nil {
		fmt.Println("ERROR: cannot dump data.")
	}
	return dumped
}
