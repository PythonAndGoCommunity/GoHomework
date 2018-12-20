package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

// var TEST_PAIRS = map[string]string{
// 	"UNSUPPORTED" + "\n":   "Unsupported command: " + "\n",
// 	"SET" + "\n":           "Error in command syntax. Syntax: set [key] [value]" + "\n",
// 	"GET" + "\n":           "Error in command syntax. Syntax: get [key]" + "\n",
// 	"DEL" + "\n":           "Error in command syntax. Syntax: del [key]" + "\n",
// 	"SET key1 val1" + "\n": "SET successful" + "\n",
// 	"SET key2 val2" + "\n": "SET successful" + "\n",
// 	"GET key2" + "\n":      "val2" + "\n",
// 	"DEL key2" + "\n":      "DEL successful" + "\n",
// 	"KEYS" + "\n":          "key1 val1" + "\n",
// }

var test_pairs = []struct {
	test_case   string
	test_result string
}{

	{"UNSUPPORTED" + "\n", "Unsupported command: " + "\n"},
	{"SET" + "\n", "Error in command syntax. Syntax: set [key] [value]" + "\n"},
	{"GET" + "\n", "Error in command syntax. Syntax: get [key]" + "\n"},
	{"DEL" + "\n", "Error in command syntax. Syntax: del [key]" + "\n"},
	{"SET key1 val1" + "\n", "SET successful" + "\n"},
	{"SET key2 val2" + "\n", "SET successful" + "\n"},
	{"GET key2" + "\n", "val2" + "\n"},
	{"DEL key2" + "\n", "DEL successful" + "\n"},
	{"KEYS" + "\n", "key1" + "\n"},
}

func TestServer(t *testing.T) {
	go main()
	time.Sleep(100 * time.Millisecond)

	conn, conn_err := net.Dial("tcp", "127.0.0.1:9090")
	if conn_err != nil {
		t.Error(conn_err)
	}
	for _, test_pair := range test_pairs {
		test_case := test_pair.test_case
		test_result := test_pair.test_result

		conn.Write([]byte(test_case))

		response, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Response: " + response)
		//fmt.Print("Result: " + test_result)
		if response != test_result {
			t.Errorf("Test failed, expected: '%s', got:  '%s'", test_result, response)
		}
	}

}
