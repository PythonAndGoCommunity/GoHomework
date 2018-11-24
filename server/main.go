package main

import (
	"NonRelDB/server/handler"
	"NonRelDB/log"
	"NonRelDB/server/storage/inmemory"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var ip string
var port string
var mode string
var location string

func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "Defines host ip")
	flag.StringVar(&port, "port", "9090", "Defines host port")
	flag.StringVar(&port, "p", "9090", "Defines host port")
	flag.StringVar(&mode, "mode", "memory", "Defines storage location")
	flag.StringVar(&mode, "m", "memory", "Defines storage location")
	flag.StringVar(&location, "location", "storage.json", "Defines storage location")
	flag.StringVar(&location, "l", "storage.json", "Defines storage location")
	flag.Parse()
}

func storageInit() {
	if mode == "memory" {
		inmemory.InitDBInMemory()
	} else if mode == "disk" {
		inmemory.InitDBFromStorage(location)
	}
}

func cleanup() {
	sign := make(chan os.Signal, 2)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sign
		log.Info.Println("Ctrl+C pressed in Terminal")
		if (mode == "disk"){
			inmemory.SaveDBToStorage(location)
		}
		os.Exit(0)
	}()
}

func main() {
	storageInit()
	cleanup()

	l, err := net.Listen("tcp", ip+":"+port)
	
	if err != nil {
		log.Error.Panicln(err.Error())
	}

	log.Info.Printf("Server started listening on %s", l.Addr().String())
	handler.HandleListener(l)
}
