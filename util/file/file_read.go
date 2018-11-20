package file

import (
	"io/ioutil"
)

// OpenAndReadString | receives file name, reads this file and returns its string content.
func OpenAndReadString(name string) string{
	b, err := ioutil.ReadFile(name)
	if err != nil{
		panic(err)
	}
	return string(b)
}

// OpenAndRead | Receives file name, reads this file and returns byte array from it. 
func OpenAndRead(name string) []byte{
	b, err := ioutil.ReadFile(name)
	if err != nil{
		panic(err)
	}
	return b
}