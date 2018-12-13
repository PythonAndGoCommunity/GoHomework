package client

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func HandleServerResponds(conn net.Conn, invitationMessage string, done chan string) {
	scanner := bufio.NewScanner(conn)
	for {
		fmt.Print(invitationMessage + "> ")
		ok := scanner.Scan()
		if !ok {
			done <- "\rServer closed the connection."
			return
		}

		text := scanner.Text()
		if text != "" {
			fmt.Println(text)
		} else {
			fmt.Println()
		}
	}
}

func HandleUserRequests(conn net.Conn, done chan string) {
	var buff bytes.Buffer
	reader := bufio.NewReader(os.Stdin)
	for {
		for {
			b, _ := reader.ReadByte()

			if strings.Compare(string(b), "\t") == 0 {
				// commandEnd := completeCommand(buff.String())
				// fmt.Fprintln(os.Stdout, commandEnd)
			}

			buff.WriteString(string(b))
			if strings.Compare(string(b), "\n") == 0 {
				break
			}
		}

		text := buff.String()
		if strings.Compare(text, "EXIT\n") == 0 {
			done <- "\rYou have been disconnected."
		}

		SendRequest(conn, text)
		buff.Reset()
	}
}

func SendRequest(conn net.Conn, text string) {
	err := conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		fmt.Println("ERROR: Cannot configure your request.")
	}

	_, err = conn.Write([]byte(text))
	if err != nil {
		fmt.Println("ERROR: Cannot send your request.")
	}
}

func completeCommand(commandPart string) string {
	commands := [4]string{"SET", "GET", "PUBLISH", "SUBSCRIBE"}

	for i := range commands {
		if strings.HasPrefix(commands[i], commandPart) {
			return commands[i][len(commandPart):]
		}
	}

	return ""
}
