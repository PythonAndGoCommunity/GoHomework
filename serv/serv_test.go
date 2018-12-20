package main

import (
	//"errors"
	"testing"
)



func TestRedisRAM(t *testing.T) {
	clInput := make(chan string)
	clOutput := make(chan string)
	clErr := make(chan error)
	clientCh := clientChan{
		input:  clInput,
		output: clOutput,
		err:    clErr,
	}


	tests := []struct {
		cmd		string
		out		string
	}{
		{"", "",},
		{"qwe", "",},
		{"GET", "",},
		{"GET qwe", "",},
		{"GET qwe qwe", "",},
		{"GET qwe qwe qwe", "",},
		{"SET", "",},
		{"SET qwe", "",},
		{"SET qwe qwe qwe", "",},
		{"DEL", "",},
		{"DEL qwe", "",},
		{"DEL qwe qwe", "",},
		{"DEL qwe qwe qwe", "",},
		{"SET key1 val1", "val1",},
		{"SET key2 val2", "val2",},
		{"SET key3 val3", "val3",},
		{"SET key4 val4", "val4",},
		{"GET key4 ", "val4",},
		{"GET key2 ", "val2",},
		{"DEL key3 ", "val3",},
		{"DEL key2 ", "val2",},
		{"DEL key3 ", "",},
		{"GET key3 ", "",},
		{"GET key2 ", "",},
	}
	go redis(clientCh)
	for i, st := range tests {

		f := func(t *testing.T) {
			t.Logf("Run test #%d for %s", i, st.cmd)

			clientCh.input <- st.cmd
			out := <-clientCh.output
			if out != st.out {
				t.Errorf("[ERR]:\n\tOutput: want %s, - got %s", st.out, out)
			}
		}
		// Run sub test
		t.Run("RAM", f)
	}

}

func TestRedisDISK(t *testing.T) {

	clInput := make(chan string)
	clOutput := make(chan string)
	clErr := make(chan error)
	clientCh := clientChan{
		input:  clInput,
		output: clOutput,
		err:    clErr,
	}

	tests := []struct {
		cmd		string
		out		string
	}{
		{"", "",},
		{"qwe", "",},
		{"GET", "",},
		{"GET qwe", "",},
		{"GET qwe qwe", "",},
		{"GET qwe qwe qwe", "",},
		{"SET", "",},
		{"SET qwe", "",},
		{"SET qwe qwe qwe", "",},
		{"DEL", "",},
		{"DEL qwe", "",},
		{"DEL qwe qwe", "",},
		{"DEL qwe qwe qwe", "",},
		{"SET key1 val1", "val1",},
		{"SET key2 val2", "val2",},
		{"SET key3 val3", "val3",},
		{"SET key4 val4", "val4",},
		{"GET key4 ", "val4",},
		{"GET key2 ", "val2",},
		{"DEL key3 ", "val3",},
		{"DEL key2 ", "val2",},
		{"DEL key3 ", "",},
		{"GET key3 ", "",},
		{"GET key2 ", "",},
	}
	go saveData("redisDatabase", clientCh)
	for i, st := range tests {

		f := func(t *testing.T) {
			t.Logf("Run test #%d for %s", i, st.cmd)

			clientCh.input <- st.cmd
			out := <-clientCh.output
			if out != st.out {
				t.Errorf("[ERR]:\n\tOutput: want %s, - got %s", st.out, out)
			}
		}
		// Run sub test
		t.Run("DISK", f)
	}
}
