package protocol

import (
	"testing"
)

func TestDecodeMessage(t *testing.T) {
	cases := []struct {
		testName        string
		inputLine       []byte
		expectedMessage Msg
	}{
		{
			testName:        "Valid 'GET' command",
			inputLine:       []byte("GET cats"),
			expectedMessage: Msg{Cmd: "GET", Key: "cats"},
		},
		{
			testName:        "Valid 'DEL' command",
			inputLine:       []byte("DEL dogs"),
			expectedMessage: Msg{Cmd: "DEL", Key: "dogs"},
		},
		{
			testName:        "Valid 'SET' command",
			inputLine:       []byte("SET name XXX"),
			expectedMessage: Msg{Cmd: "SET", Key: "name", Val: "XXX"},
		},
		{
			testName:        "Unsupported command",
			inputLine:       []byte("HGET region us-east-1"),
			expectedMessage: Msg{Cmd: "HGET", Key: "region", Val: "us-east-1"},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.testName, func(t *testing.T) {
			actualMessage, err := DecodeMessage(testCase.inputLine)
			if err != nil {
				t.Fatal(err)
			}
			compareMessages(t, actualMessage, testCase.expectedMessage)
		})
	}

}

func compareMessages(t *testing.T, actual, expected Msg) {
	if actual.Cmd != expected.Cmd {
		t.Fatalf("actual msg.Cmd(%s) != expected msg.Cmd(%s)", actual.Cmd, expected.Cmd)
	}
	if actual.Key != expected.Key {
		t.Fatalf("actual msg.Key(%s) != expected msg.Key(%s)", actual.Key, expected.Key)
	}
	if actual.Val != expected.Val {
		t.Fatalf("actual msg.Val(%s) !+ expected msg.Val(%s)", actual.Val, expected.Val)
	}
}

func TestValidateMessage(t *testing.T) {
	cases := []struct {
		testName      string
		msg           Msg
		expectedError string
	}{
		{
			testName:      "Valid 'GET' command",
			msg:           Msg{Cmd: "GET", Key: "cats"},
			expectedError: "",
		},
		{
			testName:      "'GET' without key",
			msg:           Msg{Cmd: "GET"},
			expectedError: "invalid format of a command: key is not provided",
		},
		{
			testName:      "Valid 'SET' command",
			msg:           Msg{Cmd: "SET", Key: "cats", Val: "---"},
			expectedError: "",
		},
		{
			testName:      "'SET' without key",
			msg:           Msg{Cmd: "SET"},
			expectedError: "invalid format of a command: key is not provided",
		},
		{
			testName:      "'SET' without value",
			msg:           Msg{Cmd: "SET", Key: "cats"},
			expectedError: "invalid format of a command: value is not provided",
		},
		{
			testName:      "Valid 'DEL'",
			msg:           Msg{Cmd: "DEL", Key: "cats"},
			expectedError: "",
		},
		{
			testName:      "'DEL' without key",
			msg:           Msg{Cmd: "DEL"},
			expectedError: "invalid format of a command: key is not provided",
		},
		{
			testName:      "Unsupported command should raise an error",
			msg:           Msg{Cmd: "HGET"},
			expectedError: "command 'HGET' is not supported",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.testName, func(t *testing.T) {
			err := ValidateMessage(testCase.msg)
			if err != nil && err.Error() != testCase.expectedError {
				t.Fatalf("actual error(%s) != expected error(%s)", err.Error(), testCase.expectedError)
			}
		})
	}

}
