package file

import (
	"NonRelDB/log"
	"os"
)

// CreateAndWriteString creates file and writes string to it.
func CreateAndWriteString(name string, value string) {
	file, err := os.Create(name)

	defer file.Close()

	if err != nil {
		log.Error.Panicln(err.Error())
	}

	_, err = file.WriteString(value)
	if err != nil {
		log.Error.Panicln(err.Error())
	}
}

// CreateAndWrite creates file and writes byte array to it.
func CreateAndWrite(name string, value []byte) {
	file, err := os.Create(name)

	defer file.Close()

	if err != nil {
		log.Error.Panicln(err.Error())
	}

	_, err = file.Write(value)
	if err != nil {
		log.Error.Panicln(err.Error())
	}
}
