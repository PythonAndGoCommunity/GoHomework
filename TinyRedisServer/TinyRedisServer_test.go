package main

import (
	"testing"
)

type cmd struct {
	fields []string
	result chan string
}

func TestContains(t *testing.T) {
	var stringsSet = []string{"one", "two"}
	var testString = "two"
	if contains(stringsSet, testString) {
		t.Log("\t[OK]")
	} else {
		t.Error("\t[ERR]\tShould return true")
	}

	testString = "three"
	if !contains(stringsSet, testString) {
		t.Log("\t[OK]")
	} else {
		t.Error("\t[ERR]\tShould return false")
	}

	if contains(nil, testString) {
		t.Log("\t[OK]")
	} else {
		t.Error("\t[ERR]\tShould return true")
	}
}

func TestCmdFlagParse(t *testing.T) {
	port, mode := cmdFlagParse()
	if port == "9090" && mode == "ram" {
		t.Log("\t[OK]")
	} else {
		t.Error("\t[ERR]\tError with port and mode")
	}
}

func TestWriteToFile(t *testing.T) {
	var mapData = make(map[string]string)
	mapData["0"] = "0"
	if err := writeToFile(mapData); err != nil {
		t.Error("\t[ERR]\tError with writing to file")
	} else {
		t.Log("\t[OK]")
	}
}

func TestReadFromFile(t *testing.T) {
	if _, err := readFromFile(); err != nil {
		t.Error("\t[ERR]\tError with reading from file")
	} else {
		t.Log("\t[OK]")
	}
}

func TestCreateStorage(t *testing.T) {
	var mapData = make(map[string]string)

	if line := createStorage(mapData, ""); line != "" {
		t.Error("\t[ERR]\tError with storage commands processing")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "1"); line != "Expected at least 2 arguments!" {
		t.Error("\t[ERR]\tError with storage commands processing")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "SET 0"); line != "Expected value!" {
		t.Error("\t[ERR]\tError with setting to storage")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "SET 1 1 1"); line != "Expected Key and Value!" {
		t.Error("\t[ERR]\tError with setting to storage")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "SET 1 2"); line != "done" {
		t.Error("\t[ERR]\tError with setting to storage")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "GET 1"); line != "2" {
		t.Error("\t[ERR]\tError with getting from storage")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "DEL 1"); line != "deleted" {
		t.Error("\t[ERR]\tError with deleting from storage")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "GET 1"); line != "nil" {
		t.Error("\t[ERR]\tError with getting from storage")
	} else {
		t.Log("\t[OK]")
	}

	if line := createStorage(mapData, "REQUEST 1"); line != "Invalid command \"REQUEST\"" {
		t.Error("\t[ERR]\tError with storage commands processing")
	} else {
		t.Log("\t[OK]")
	}
}
