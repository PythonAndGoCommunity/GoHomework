package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//Commands - read stdin, input connection, write to connection
func Commands(protocolTCP string, gPort string, gIP string, buff int) (string, error) {
	var (
		buffer    = make([]byte, buff)
		IPHandler = []string{gIP + gPort}
		answer    string
	)

	Dconn, err := net.Dial(protocolTCP, gPort)
	if err != nil {
		log.Fatal(err)
	}

	for {
		first, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatal(err)

		}
		mssg := strings.Replace(first, "\n", " ", -1)
		Splitmssg := strings.Fields(mssg)

		switch Splitmssg[0] {
		case "--exit", "-q":
			fmt.Println("Exit. See ya.")
			os.Exit(1)
			//return mssg, nil // for test
		case "--connect":
			if len(Splitmssg) == 2 {
				var toggle = false
				for _, addr := range IPHandler {
					if addr == Splitmssg[1] {
						toggle = true
						break
					}
				}
				if !toggle {
					conn, err := net.Dial(protocolTCP, Splitmssg[1])
					if err == nil {
						IPHandler = append(IPHandler, Splitmssg[1])
						conn.Write([]byte(fmt.Sprintf("User '%s:%s join'\n", gIP, gPort)))
						conn.Close()
					}
				}
			}
		default:
			//write to socket
			_, werr := Dconn.Write([]byte(mssg))
			if werr != nil {
				log.Fatal(werr)
			}
			//read from socket
			read, err := Dconn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
			answer = string(buffer[:read])
			fmt.Printf("Answer: %s.", answer)
		}
	}
	return answer, nil
}
