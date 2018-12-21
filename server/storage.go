package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/SiarheiKresik/go-kvdb/db"
)

func storage(cmd chan command, mode string) {
	// initialize database
	data, err := db.InitDB(mode)
	if err != nil {
		log.Fatalln("Error while initializing database:", err)
	}

	// TODO this huge code block needs to be refactored to separate functions
	// TODO storage function is not right place to parse commands from string
	// wait and handle commands from the commands channel
	for cmd := range cmd {
		if len(cmd.fields) < 1 {
			cmd.result <- ""
			continue
		}

		// Executing command
		switch cmd.fields[0] {

		// GET <KEY>
		case "GET":
			result, ok := data.Get(cmd.fields[1])
			if ok {
				cmd.result <- fmt.Sprintf("\"%s\"", result)
			} else {
				cmd.result <- "no such key"
			}

		// SET <KEY> <VALUE>
		case "SET":
			if len(cmd.fields) != 3 {
				cmd.result <- "EXPECTED VALUE"
				continue
			}
			err := data.Set(cmd.fields[1], cmd.fields[2])
			if err != nil {
				cmd.result <- "error"
			} else {
				cmd.result <- "Ok"
			}

		// DEL <KEY>
		case "DEL":
			err := data.Delete(cmd.fields[1])
			if err != nil {
				cmd.result <- "error"
			} else {
				cmd.result <- "Ok"
			}

		// KEYS
		case "KEYS":
			// set pattern to '*' if pattern is not setted by client
			var pattern string
			if len(cmd.fields) == 1 {
				pattern = "*"
			} else {
				pattern = cmd.fields[1]
			}
			// get result from database
			keys, err := data.Keys(pattern)
			// send response
			if err != nil {
				cmd.result <- "Invalid pattern"
			} else {
				l := len(keys)
				cmd.result <- strings.Join(keys, " ") + fmt.Sprintf(", number of keys: %d", l)
			}

		// for debugging only
		case "DUMP":
			json, err := data.Dump()
			if err != nil {
				cmd.result <- "Error while dumping"
			} else {
				cmd.result <- string(json)
			}

		// for debugging only
		case "SAVE":
			err := data.Save()
			if err != nil {
				cmd.result <- "Error while saving"
			} else {
				cmd.result <- "Saved"
			}

		// for debugging only
		case "LOAD":
			err := data.Load()
			if err != nil {
				cmd.result <- "Error while loading"
			} else {
				cmd.result <- "Loaded"
			}

		default:
			cmd.result <- "Invalid command " + cmd.fields[0]
		}
	}
}
