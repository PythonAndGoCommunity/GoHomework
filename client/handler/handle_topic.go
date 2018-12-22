package handler

import (
	"NonRelDB/log"
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// quit Handling Ctrl + C press.
func quit(c net.Conn, topic string) {
	sign := make(chan os.Signal, 2)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sign
		fmt.Println("Ctrl+C pressed in Terminal")
		fmt.Fprintf(c, "unsubscribe "+topic+"\n")
		os.Exit(0)
	}()
}

// HandleTopic listening messages from specified topic.
func HandleTopic(c net.Conn, r bufio.Reader, topic string) {
	quit(c, topic)
	fmt.Printf("Reading messages from %s... (press Ctrl + C to stop)\n", topic)
	for {
		msg, err := r.ReadString('\n')

		if err != nil {
			log.Error.Panicln(err.Error())
		}

		fmt.Print(msg)
	}
}
