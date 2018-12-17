package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestStoreMethods(t *testing.T) {
	ioutil.WriteFile("storage", []byte(""), 0755)
	file, _ := os.OpenFile("storage", os.O_CREATE|os.O_RDWR, 0775)
	file.WriteString("bird;Twit\n")
	file.Seek(0, 0)
	newStorage := storage{
		map[string]string{
			"cat":    "Tom",
			"dog":    "Bob",
			"duck":   "Donald",
			"mouse":  "Djery",
			"monkey": "Lisa",
			"snake":  "Ka",
			"lion":   "Symba",
		},
		"disk",
		file,
	}

	t.Log("GET method of storage")
	ch1 := make(chan string, 1)
	go newStorage.GET("cat", ch1)
	result := <-ch1
	if result == "cat:Tom" {
		t.Log("\t[OK]\tShould get Tom")
	} else {
		t.Error("\t[ERR]\tShould get Tom")
	}

	ch2 := make(chan string, 1)
	go newStorage.GET("wolf", ch2)
	if <-ch2 == "no pair contains key=wolf" {
		t.Log("\t[OK]\tStore no contain key dog")
	} else {
		t.Error("\t[ERR]\tShould Error wrong value return")
	}

	ch2s := make(chan string, 1)
	go newStorage.GET("bird", ch2s)
	if <-ch2s == "bird:Twit" {
		t.Log("\t[OK]\tStore no contain key dog")
	} else {
		t.Error("\t[ERR]\tShould Error wrong value return")
	}

	t.Log("SET method of storage")
	ch3 := make(chan string, 1)
	go newStorage.SET("panda", "Pow", ch3)
	if <-ch3 == "pair panda:Pow created" {
		t.Log("\t[OK]\tPair panda:Pow inserted in to storage")
	} else {
		t.Error("\t[ERR]\tStore contains key panda, do not may overwrite it")
	}

	ch4 := make(chan string, 1)
	go newStorage.SET("cat", "Poor", ch4)
	if <-ch4 == "store contains pair with key cat: Tom" {
		t.Log("\t[OK]\tStore contains key cat")
	} else {
		t.Error("\t[ERR]\tPair cat: Poor wrote to storage")
	}

	t.Log("DEL method of storage")
	ch5 := make(chan string, 1)
	go newStorage.DEL("dog", ch5)
	if <-ch5 == "pair deleted" {
		t.Log("\t[OK]\tKey dog deleted")
	} else {
		t.Error("\t[ERR]\tNo find key dog in to storage")
	}

	ch6 := make(chan string, 1)
	go newStorage.DEL("cow", ch6)
	if <-ch6 == "no pair for delete" {
		t.Log("\t[OK]\tStore don't contains key cow")
	} else {
		t.Error("\t[ERR]\tKey cow deleted")
	}

	t.Log("KEYS method of storage")
	ch7 := make(chan string, 1)
	go newStorage.KEYS("on", ch7)
	res := <-ch7
	//fmt.Println(res, "res")
	if reflect.DeepEqual("monkey:Lisa\nlion:Symba\n", res) {
		t.Log("\t[OK]\tMethod keys is work")
	} else if reflect.DeepEqual("lion:Symba\nmonkey:Lisa\n", res) {
		t.Log("\t[OK]\tMethod keys is work")
	} else {
		t.Error("\t[ERR]\tMethod keys not work")
	}

	ch8 := make(chan string, 1)
	go newStorage.KEYS("bear", ch8)
	if reflect.DeepEqual("None:None\n", <-ch8) {
		t.Log("\t[OK]\tMethod keys is work")
	} else {
		t.Error("\t[ERR]\tMethod keys not work")
	}

	t.Log("Write data to file method")
	newStorage.writeToFileAllData()
	readingData := newStorage.readFromFileAllData()
	if reflect.DeepEqual(readingData, newStorage.storageMap) {
		t.Log("\t[OK]\tWrite/read file")
	} else {
		t.Error("\t[ERR]\tKey cow deleted")
	}
}
