package file

import (
	"io/ioutil"
)

func OpenAndReadString(name string) string{
	b, err := ioutil.ReadFile(name)
	if err != nil{
		panic(err)
	}
	return string(b)
}

func OpenAndRead(name string) []byte{
	b, err := ioutil.ReadFile(name)
	if err != nil{
		panic(err)
	}
	return b
}