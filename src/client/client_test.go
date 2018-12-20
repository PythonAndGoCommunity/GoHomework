package client

import (
	"os"
	"testing"
)

func TestTrimLastSymbol(t *testing.T) {
	t.Log("\tTrim last comma test")

	tests := []struct {
		input string
		want  string
	}{
		{"9090,", "9090"},
		{"9090", "9090"},
		{"-p=9090, -h=127.0.0.1", "-p=9090, -h=127.0.0.1"},
		{",", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if trimLastSymbol(tt.input) == tt.want {
				t.Log("\t[OK]\tShould get '" + tt.want + "'")
			} else {
				t.Error("\t[ERR]\tShould get '" + tt.want + "'")
			}
		})
	}
}

func Test_parseArguments(t *testing.T) {
	t.Log("Parse arguments test")

	tests := []struct {
		name    string
		Args    []string
		port    string
		host    string
		dump    bool
		restore bool
	}{
		{"no arguments", []string{"client.go"}, defaultPort, defaultHost, false, false},
		{"one wrong argument", []string{"client.go", "-pt=t"}, defaultPort, defaultHost, false, false},
		{"one argument", []string{"client.go", "--port=9091"}, "9091", defaultHost, false, false},
		{"port after 49152", []string{"client.go", "--port=50505"}, "9090", defaultHost, false, false},
		{"port after 49152", []string{"client.go", "-p=50505"}, "9090", defaultHost, false, false},
		{"port less 1024", []string{"client.go", "--port=1000"}, "9090", defaultHost, false, false},
		{"port less 1024", []string{"client.go", "-p=1000"}, "9090", defaultHost, false, false},
		{"two arguments", []string{"client.go", "-p=1000", "-h=127.0.0.2"}, defaultPort, "127.0.0.2", false, false},
		{"two arguments", []string{"client.go", "-p=9090", "--host=127.0.0.2"}, defaultPort, "127.0.0.2", false, false},
		{"three arguments", []string{"client.go", "-p=9091", "--host=127.0.0.3", "--dump"}, "9091", "127.0.0.3", true, false},
		{"four arguments", []string{"client.go", "-p=9092", "--host=127.0.0.4", "--dump", "--restore"}, "9092", "127.0.0.4", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_mainPort = defaultPort
			_mainHost = defaultHost
			_dump = false
			_restore = false

			os.Args = tt.Args
			parseArguments()
			if tt.host == _mainHost && tt.port == _mainPort && tt.dump == _dump && tt.restore == _restore {
				t.Log("\t[OK]\t parse with " + tt.name)
			} else {
				t.Error("\t[ERR]\t parse with " + tt.name + " port: " + _mainPort)
			}
		})
	}
}

