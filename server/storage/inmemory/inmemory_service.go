package inmemory

import (
	"os"
	"NonRelDB/server/util/json"
	"NonRelDB/server/util/file"
	"NonRelDB/server/log"
)

// Storage | Global variable for kv storage.
var Storage *map[string]string

// InitDBInMemory | Init kv db in memory.
func InitDBInMemory(){
	s := make(map[string]string)
	Storage = &s	
	log.Info.Println("DB successfully created in-memory")
}

// InitDBFromStorage | Receives filename and load its content to inmemory storage.
func InitDBFromStorage(filename string){
	log.Info.Println("[" + filename+"]")
	log.Info.Println(len(filename))

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Warning.Println(err.Error())
		log.Warning.Printf("Storage doesnt exist. Will be created new with name %s", filename)

		f, err := os.Create(filename)
		
		if err != nil {
			log.Error.Panicln(err.Error())
		}
		log.Info.Printf("DB successfully initialized from %s", filename)
		f.Close()
	}

	j := file.OpenAndReadString(filename)
	jb := []byte(j)
	Storage = json.UnpackFromJSON(jb)
}

// SaveDBToStorage | Receives file name and saves inmemory storage to it.
func SaveDBToStorage(filename string){
	jb := json.PackMapToJSON(*Storage)
	j := string(jb)

	file.CreateAndWriteString(filename, j)
	log.Info.Printf("DB successfully saved to %s", filename)
}
