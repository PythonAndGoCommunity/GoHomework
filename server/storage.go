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

	// wait and handle commands from the commands channel
	for cmd := range cmd {
		if len(cmd.fields) < 1 {
			cmd.result <- ""
			continue
		}

		log.Println("Command:", cmd.fields)

		// Executing command
		switch cmd.fields[0] {

		// GET <KEY>
		case "GET":
			result, ok := data.Get(cmd.fields[1])
			fmt.Println("sending result:", result)
			if ok {
				cmd.result <- "\"" + result + "\""
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

			keys := data.Keys()
			l := len(keys)
			cmd.result <- strings.Join(keys, " ") + fmt.Sprintln(", number of keys:", l)

		// DUMP
		case "DUMP":
			json, err := data.Dump()
			if err != nil {
				cmd.result <- "Error while "
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
