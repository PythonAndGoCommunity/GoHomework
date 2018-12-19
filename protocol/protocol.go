package protocol

import (
	"bufio"
	"bytes"
	"fmt"
)

// Msg represent message from a client
type Msg struct {
	Cmd string
	Key string
	Val string
}

// DecodeMessage parse message from bytes to Msg struct
func DecodeMessage(line []byte) (Msg, error) {
	scanner := bufio.NewScanner(bytes.NewReader(line))
	scanner.Split(bufio.ScanWords)
	data := make([]string, 3)
	for i := 0; i < 3 && scanner.Scan(); i++ {
		data[i] = scanner.Text()
	}
	return Msg{Cmd: data[0], Key: data[1], Val: data[2]}, scanner.Err()
}

// ValidateMessage check weither message is valid or not
func ValidateMessage(msg Msg) error {
	switch msg.Cmd {
	case "GET", "DEL":
		if err := checkField("key", msg.Key); err != nil {
			return err
		}
	case "SET":
		if err := checkField("key", msg.Key); err != nil {
			return err
		}
		if err := checkField("value", msg.Val); err != nil {
			return err
		}
	default:
		return fmt.Errorf("command '%s' is not supported", msg.Cmd)
	}

	return nil
}

func checkField(name, value string) error {
	if value == "" {
		return fmt.Errorf("invalid format of a command: %s is not provided", name)
	}
	return nil
}
