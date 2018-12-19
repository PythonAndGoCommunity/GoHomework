package test

import (
	"GoHomework/client/client"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	commandEnd := client.CompleteCommand("SUB")
	if strings.Compare(commandEnd, "SCRIBE") != 0 {
		t.Error("[ERR] CompleteCommand failed. Expected SUB -> SUBSCRIBE, got SUB -> " + commandEnd)
	}

	commandEnd = client.CompleteCommand("UNKNOWN")
	if strings.Compare(commandEnd, "") != 0 {
		t.Error("[ERR] CompleteCommand failed. Expected empty, got " + commandEnd)
	}
}
