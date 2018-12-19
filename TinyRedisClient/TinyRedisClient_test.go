package main

import "testing"

func TestGetAddress(t *testing.T) {
	address := getAddress()
	if address == "127.0.0.1:9090" {
		t.Log("\t[OK]")
	} else {
		t.Error("\t[ERR]\tError with address")
	}
}

func TestConnectionAttempt(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Error with connection")
		} else {
			t.Log("\t[OK]")
		}
	}()

	connectionAttempt("127:0:0:1:9090")
}
