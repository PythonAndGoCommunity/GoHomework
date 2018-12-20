package server

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"
)

func TestTrimLastSymbol(t *testing.T) {
	t.Log("\tTrim last comma test")
	if trimLastSymbol("9090,") == "9090" {
		t.Log("\t[OK]\tShould get '9090'")
	} else {
		t.Error("\t[ERR]\tShould get '9090'")
	}

	if trimLastSymbol("9090") == "9090" {
		t.Log("\t[OK]\tShould get '9090'")
	} else {
		t.Error("\t[ERR]\tShould get '9090'")
	}

	if trimLastSymbol("-p=9090") == "-p=9090" {
		t.Log("\t[OK]\tShould get '-p=9090'")
	} else {
		t.Error("\t[ERR]\tShould get '-p=9090'")
	}

	if trimLastSymbol(",") == "" {
		t.Log("\t[OK]\tShould get ''")
	} else {
		t.Error("\t[ERR]\tShould get ''")
	}

}

func TestPrintln(t *testing.T) {
	t.Log("\tLoggers test")
	initLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	_mainVerbose = true
	Println(Trace, "test println")
	Println(Info, "test println")
	Println(Warning, "test println")
	Println(Error, "test println")
	cmd := "test"
	Printf(Trace, "%s printf", cmd)
	Printf(Info, "%s printf", cmd)
	Printf(Warning, "%s printf", cmd)
	Printf(Error, "%s printf", cmd)
	_mainVerbose = false
}

func Test_newJsonWriter(t *testing.T) {
	t.Log("New JSON writer test")
	file, err := ioutil.TempFile(os.TempDir(), "temp.json")

	if err == nil {

		type args struct {
			fileName string
		}
		tests := []struct {
			name string
		}{
			{"temp file"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				newJSONWriter(file.Name())
				t.Log("\t[OK]\t temp file " + tt.name)
			})
		}
	} else {
		t.Error("\t[ERR]\tCan't create temporary file")
	}
	defer os.Remove(file.Name())
}

func Test_parseArguments(t *testing.T) {
	t.Log("Parse arguments test")
	tests := []struct {
		name    string
		Args    []string
		port    string
		mode    bool
		verbose bool
	}{
		{"no arguments", []string{"server.go"}, "9090", true, false},
		{"wrong arguments 1", []string{"server.go", " -g=t"}, "9090", true, false},
		{"wrong arguments 2", []string{"server.go", "--port=t"}, "9090", true, false},
		{"one argument", []string{"server.go", "-p=9091"}, "9091", true, false},
		{"port after 49152", []string{"server.go", "-p=55000"}, "9090", true, false},
		{"port less 1024", []string{"server.go", "-p=1020"}, "9090", true, false},
		{"port after 49152", []string{"server.go", "--port=55000"}, "9090", true, false},
		{"port less 1024", []string{"server.go", "--port=1020"}, "9090", true, false},
		{"two arguments", []string{"server.go", "--port=9091,", "-m=disk"}, "9091", true, false},
		{"mode broke", []string{"server.go", "--port=9091,", "-m=broke_disk"}, "9091", true, false},
		{"mode broke", []string{"server.go", "--port=9091,", "--mode=broke_disk"}, "9091", true, false},
		{"three arguments", []string{"server.go", "--port=9092,", "--mode=memory,", "-v"}, "9092", false, true},
		{"different arguments", []string{"server.go", "--port=9090,", "--mode=disk,", "--verbose"}, "9090", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_mainPort = defaultPort
			_mainMode = true     //mode=true - save in disk; mode=false - save in memory
			_mainVerbose = false //verbose=false - without logging; verbose=true - with logging

			os.Args = tt.Args
			parseArguments()
			if tt.mode == _mainMode && tt.port == _mainPort && tt.verbose == _mainVerbose {
				t.Log("\t[OK]\t parse with " + tt.name)
			} else {
				t.Error("\t[ERR]\t parse with " + tt.name + " port: " + _mainPort)
			}
		})
	}
}

func mockClientConnectSet(t *testing.T, msg string) string {
	conn, _ := net.Dial("tcp", "127.0.0.1:9090")
	// send to socket
	fmt.Fprintf(conn, msg)
	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	defer conn.Close()
	return message
}

func mockClientConnectGet(t *testing.T, msg string) string {
	conn, _ := net.Dial("tcp", "127.0.0.1:9090")
	// send to socket
	fmt.Fprintf(conn, msg)
	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	time.Sleep(5 * time.Millisecond)
	defer conn.Close()
	return message
}

func mockClientConnectDel(t *testing.T, msg string) string {
	conn, _ := net.Dial("tcp", "127.0.0.1:9090")
	// send to socket
	fmt.Fprintf(conn, msg)
	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	defer conn.Close()
	return message
}

func mockClientConnectKeys(t *testing.T, msg string) string {
	conn, _ := net.Dial("tcp", "127.0.0.1:9090")
	fmt.Fprintf(conn, msg)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	defer conn.Close()
	return message
}

func Test_main(t *testing.T) {
	t.Log("CI test")

	set := "set kT0 v0\r\n"
	set1 := "set kT1 v1\r\n"
	set2 := "set kT2\r\n"

	os.Args = []string{"server.go", "-m=memory"}
	go main()
	time.Sleep(50 * time.Millisecond)

	if mockClientConnectSet(t, set) == "OK\r\n" {
		t.Log("\t[OK]\t good response from 'SET kT0 v0'")
	} else {
		t.Error("\t[ERR]\t with response from 'SET kT0 v0'")
	}
	if mockClientConnectSet(t, set1) == "OK\r\n" {
		t.Log("\t[OK]\t good response from 'SET kT1 v1'")
	} else {
		t.Error("\t[ERR]\t with response from 'SET kT1 v1'")
	}
	if mockClientConnectSet(t, set2) == "OK\r\n" {
		t.Log("\t[OK]\t good response from 'SET kT2 v2'")
	} else {
		t.Error("\t[ERR]\t with response from 'SET kT2 v2'")
	}
	time.Sleep(50 * time.Millisecond)

	get := "get kT1 kt2 kt3 kT4\r\n"
	get1 := "GET kT2\r\n"
	get2 := "Get kT3\r\n"

	if mockClientConnectGet(t, get) != "v1\r\n" {
		t.Log("\t[OK]\t good response from 'GET kT1'")
	} else {
		t.Error("\t[ERR]\t with response from 'GET kT1'")
	}
	if mockClientConnectGet(t, get1) == "\r\n" {
		t.Log("\t[OK]\t good response from 'GET kT2'")
	} else {
		t.Error("\t[ERR]\t with response from 'GET kT2'")
	}
	if mockClientConnectGet(t, get2) == "(nil)\r\n" {
		t.Log("\t[OK]\t good response from 'GET kT3'")
	} else {
		t.Error("\t[ERR]\t with response from 'GET kT3'")
	}
	time.Sleep(50 * time.Millisecond)

	del1 := "del kT3\r\n"
	del2 := "del\r\n"

	if mockClientConnectDel(t, del1) == "0\r\n" {
		t.Log("\t[OK]\t good response from 'del kT3'")
	} else {
		t.Error("\t[ERR]\t with response from 'del kT3'")
	}
	if mockClientConnectDel(t, del2) == "ERROR: wrong number of arguments for 'del' command\r\n" {
		t.Log("\t[OK]\t good response from 'del'")
	} else {
		t.Error("\t[ERR]\t with response from 'del'")
	}

	mockClientConnectDel(t, "DEL kT0\r\n")
	if mockClientConnectKeys(t, "KEYS k*\r\n") == "kT1,kT2\r\n" {
		t.Log("\t[OK]\t good response from 'DEL kT0' 'KEYS k*'")
	} else {
		t.Error("\t[ERR]\t with response from 'DEL kT0' 'KEYS k*'")
	}

}

