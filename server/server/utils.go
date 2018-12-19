package server

import (
	"fmt"
	"log"
)

func PrintlnAndLog(msg string) {
	fmt.Println(msg)
	if Verbose {
		log.Println(msg)
	}
}
