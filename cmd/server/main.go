package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ITandElectronics/GoHomework/disk"
	"github.com/ITandElectronics/GoHomework/server"
)

func main() {
	const storagePath = "./db.json"

	var (
		port int
		mode string
	)

	flag.IntVar(&port, "port", 9090, "Port to listen on")
	flag.IntVar(&port, "p", 9090, "Port to listen on")
	flag.StringVar(&mode, "mode", "disk", "Storage options. One of [disk]")
	flag.StringVar(&mode, "m", "disk", "Storage options. One of [disk]")
	flag.Usage = func() {
		usage := `Usage of %s:
  --mode, -m string
	  Storage options. One of [disk] (default "disk")
  --port, -p int
	  Port to listen on (default 9090)
	`
		fmt.Fprintf(os.Stderr, usage, os.Args[0])

	}
	flag.Parse()

	fmt.Printf("server is going to start on '%d' and work in '%s' mode\n", port, mode)
	storage, err := disk.New(storagePath)
	if err != nil {
		log.Fatalf("coudn't create stroage: %v\n", err)
	}
	s, err := server.New(port, storage)
	if err != nil {
		log.Fatalf("coudln't create server: %v\n", err)
	}
	log.Fatalf("unable to start server: %v\n", s.Run())
}
