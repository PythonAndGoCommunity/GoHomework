package file

import (
	"os"
)

func CreateAndWriteString(name string, value string){
	f, err := os.Create(name)
	if err != nil{
		panic(err)
	}
	_, err = f.WriteString(value)
	if err != nil{
		panic(err)
	}
	f.Close()
}

func CreateAndWrite(name string, value []byte){
	f, err := os.Create(name)
	if err != nil{
		panic(err)
	}

	_, err = f.Write(value)
	if err != nil{
		panic(err)
	}
	f.Close()
}

