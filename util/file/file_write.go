package file

import (
	"os"
	"NonRelDB/log"
)

// CreateAndWriteString | Creates file and writes string to it.
func CreateAndWriteString(name string, value string){
	f, err := os.Create(name)

	defer f.Close()

	if err != nil{
		log.Error.Panicln(err.Error())
	}

	_, err = f.WriteString(value)
	if err != nil{
		log.Error.Panicln(err.Error())
	}
}

// CreateAndWrite | Creates file and writes byte array to it.
func CreateAndWrite(name string, value []byte){
	f, err := os.Create(name)

	defer f.Close()

	if err != nil{
		log.Error.Panicln(err.Error())
	}

	_, err = f.Write(value)
	if err != nil{
		log.Error.Panicln(err.Error())
	}
}

