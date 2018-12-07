package file

import (
	"io/ioutil"
	"NonRelDB/log"
)

// OpenAndReadString receives file name, reads this file and returns its string content.
func OpenAndReadString(name string) string{
	bytes, err := ioutil.ReadFile(name)
	
	if err != nil{
		log.Error.Panicln(err.Error())
	}
	
	return string(bytes)
}

// OpenAndRead receives file name, reads this file and returns byte array from it. 
func OpenAndRead(name string) []byte{
	bytes, err := ioutil.ReadFile(name)
	
	if err != nil{
		log.Error.Panicln(err.Error())
	}
	
	return bytes
}