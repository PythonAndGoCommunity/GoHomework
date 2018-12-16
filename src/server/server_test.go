package main

import (
	"testing"
	"bufio"
	"net"
	"time"
)


var tests = map[string] string {
	"set a b\n" : "Ok\n",
	"get a\n" : "\"b\"\n",
	"del a\n" : "Ok\n",
}

func TestHandleConnection(t *testing.T){
	go main()
	time.Sleep(time.Millisecond)
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
			t.Error(err)
			
	}
	defer conn.Close()
	serverWriter := bufio.NewWriter(conn)
	serverReader := bufio.NewReader(conn)
	
	for reqst, resp := range tests{
		n, err := serverWriter.WriteString(reqst)
		if n == 0 || err != nil {
			t.Error(err)
		}
		serverWriter.Flush()
		reply, err := serverReader.ReadString('\n')
		if err != nil {
			t.Error(err)
		}
		if reply != resp {
			t.Error("Excepted ", resp, " got", reply)
		}
	}
}