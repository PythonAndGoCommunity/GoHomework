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
	expectedHost string
} {
	{"client -h 131.0.0.1", ":9090", "131.0.0.1"},
	{"client -p :9071", ":9071", "127.0.0.1"},
	{"client -p :9071 -h 131.0.0.1", ":9071", "131.0.0.1"},
	{"client -h 131.0.0.1 -p :9071", ":9071", "131.0.0.1"},
	{"client", ":9090", "127.0.0.1"},
}

var commandsTests = []struct {
	arguments string
	expectedFlag bool
	expectedMessage string
} {
	{"SET key1 value1", false, "SET key1 value1"},
	{"GET key1", false, "GET key1"},
	{"DEL key1", false, "DEL key1"},
	{"SET \"abc dg\" \"petro\"", false, "SET abc_dg petro"},
	{"GET \"abc dg\"", false, "GET abc_dg"},
	{"DEL \"abc dg\"", false, "DEL abc_dg"},
	{"SeT key1 value1", true, cv.UsageCommandsErrorMessage},
	{"SET \"abc dg \"petro\"", true, cv.ArgumentsErrorMessage},
	{"GET \"abc dg", true, cv.ArgumentsErrorMessage},
	{"DEL abc dg\"", true, cv.ArgumentsErrorMessage},
	{"DEL \"\"abc dg", true, cv.ArgumentsErrorMessage},

}

func TestClientArgumentParsing(t *testing.T) {
	for _, test := range argumentTests {
		resA, resH := ClientArgumentParsing(strings.Fields(test.arguments))
		if resA != test.expectedAddress || resH != test.expectedHost {
			t.Errorf("Expected address to be %s. Expected host to be %s. Got %s and %s.",
				test.expectedAddress, test.expectedHost, resA, resH)
		}
	}
	t.Run("3 argument string.", TestCAP6)
	t.Run("1 argument string.", TestCAP7)
	t.Run("5 argument string.", TestCAP8)
	t.Run("Random string.", TestCAP9)
}

func TestCAP6(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ClientArgumentParsing(strings.Fields("client -h 131.0.0.1 -p"))
	t.Error("Allowed 3 arguments to be passed.")
}

func TestCAP7(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ClientArgumentParsing(strings.Fields("client -h"))
	t.Error("Allowed 1 argument to be passed.")
}

func TestCAP8(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ClientArgumentParsing(strings.Fields("client -h 131.0.0.1 -p :9071 dljfkad"))
	t.Error("Allowed 5 arguments to be passed.")
}

func TestCAP9(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered.", r)
		}
	}()
	ClientArgumentParsing(strings.Fields("client -asdf 1ags31.0.0.1 -sadg :9sgad071"))
	t.Error("Passed some random string and it worked.")
}

func TestCheckMessage(t *testing.T) {
	for _, test := range commandsTests {
		resF, resM := CheckMessage(test.arguments)
		if resF == test.expectedFlag {
			if !resF {
				if resM != test.expectedMessage {
					t.Errorf("Expected corrected message to be %s. Got %s.",
						test.expectedMessage, resM)
				}
			}
		} else {
			if !resF {
				t.Errorf("Expected command to be wrong. Message - %s.", resM)
			} else {
				t.Errorf("Expected command to be correct. Message - %s.", resM)
			}
		}
	}
}