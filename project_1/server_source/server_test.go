package main

import (
	"fmt"
	"strings"
	"testing"
	cv "project_1/commonVariables"
)

var argumentTests = []struct {
	arguments string
	expectedAddress string
	expectedMode string
} {
	{"server -m disk", ":9090", "disk"},
	{"server -p :9071", ":9071", "disk"},
	{"client -p :9071 -m disk", ":9071", "disk"},
	{"client -m disk -p :9071", ":9071", "disk"},
	{"client", ":9090", "disk"},
}

var databaseTests = []struct {
	arguments []string
	expectedAnswer string
} {
	{  strings.Fields("SET key4 value4"), "OK."},
	{  strings.Fields("SET key4 value5"), "Replaced previous value - value4."},
	{  strings.Fields("GET key4"), "value5"},
	{  strings.Fields("GET key5"), "(nil)"},
	{  strings.Fields("DEL key4 key5"), "1 - deleted. 2 - ignored."},
	{  strings.Fields("SET key5 asdasd sada sad"), cv.SetErrorMessage},
	{  strings.Fields("GET key4 sad dasd"), cv.GetErrorMessage},
	{  strings.Fields(""), ""},
	{  strings.Fields("PET asda asd"), "Wrong command."},
	{  strings.Fields("HEY!"), "Expected at least 1 argument."},
}

func TestServerArgumentParsing(t *testing.T) {
	for _, test := range argumentTests {
		resA, resM := ServerArgumentParsing(strings.Fields(test.arguments))
		if resA != test.expectedAddress || resM != test.expectedMode {
			t.Errorf("Expected address to be %s. Expected mode to be %s. Got %s and %s.",
				test.expectedAddress, test.expectedMode, resA, resM)
		}
	}
	t.Run("3 argument string.", TestSAP6)
	t.Run("1 argument string.", TestSAP7)
	t.Run("5 argument string.", TestSAP8)
	t.Run("Random string.", TestSAP9)
}

func TestSAP6(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ServerArgumentParsing(strings.Fields("client -p :9071 -m"))
	t.Error("Allowed 3 arguments to be passed.")
}

func TestSAP7(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ServerArgumentParsing(strings.Fields("client -p"))
	t.Error("Allowed 1 argument to be passed.")
}

func TestSAP8(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ServerArgumentParsing(strings.Fields("client -m disk -p :9071 dljfkad"))
	t.Error("Allowed 5 arguments to be passed.")
}

func TestSAP9(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ServerArgumentParsing(strings.Fields("client -asdf 1ags31.0.0.1 -sadg :9sgad071"))
	t.Error("Passed some random string and it worked.")
}

func TestDatabase(t *testing.T) {
	commands := make(chan []string)
	answers := make(chan cv.Answer)
	go Database(commands, answers)
	for _, test := range databaseTests {
		commands <- test.arguments
		result := <-answers
		if result.Answer != test.expectedAnswer {
			t.Errorf("Expected answer - %s. Got - %s.", test.expectedAnswer, result.Answer)
		}
	}
}
