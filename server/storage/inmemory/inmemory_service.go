package inmemory

import (
	"os"
	"NonRelDB/util/json"
	"NonRelDB/util/file"
	"NonRelDB/log"
)

// Storage | Global variable for kv storage.
var storage *map[string]string

// GetStorage | Getter for storage.
func GetStorage() *map[string]string {
	return storage
}

// InitDBInMemory | Init kv db in memory.
func InitDBInMemory(){
	s := make(map[string]string)
	storage = &s	
	log.Info.Println("DB successfully created in-memory")
}

// InitDBFromStorage | Receives filename and load its content to inmemory storage.
func InitDBFromStorage(filename string){
	defer log.Info.Printf("DB successfully initialized from %s", filename)

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Warning.Println(err.Error())
		log.Warning.Printf("Storage doesnt exist. Will be created new with name %s", filename)

		f, err := os.Create(filename)
		
		if err != nil {
			log.Error.Panicln(err.Error())
		}
		f.Close()
	}

	j := file.OpenAndReadString(filename)
	jb := []byte(j)
	storage = json.UnpackFromJSON(jb)
}

// SaveDBToStorage | Receives file name and saves inmemory storage to it.
func SaveDBToStorage(filename string){
	jb := json.PackMapToJSON(*storage)
	j := string(jb)

	file.CreateAndWriteString(filename, j)
	log.Info.Printf("DB successfully saved to %s", filename)
}
