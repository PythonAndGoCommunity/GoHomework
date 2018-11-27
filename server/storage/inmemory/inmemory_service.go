package inmemory

import (
	"os"
	"NonRelDB/util/json"
	"NonRelDB/util/file"
	"NonRelDB/log"
	"NonRelDB/server/storage/sync"
)

// Storage | Global variable for kv storage.

var storage sync.Map

// GetStorage | Getter for storage.
func GetStorage() *sync.Map {
	return &storage
}

// SetStorage | Setter for storage
func SetStorage(sm sync.Map){
	storage = sm
}

// InitDBInMemory | Init kv db in memory.
func InitDBInMemory(){
	storage := sync.Map{}
	m := make(map[string]string)
	storage.SetMap(&m)
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
	storage.SetMap(json.UnpackFromJSON(jb))
}

// SaveDBToStorage | Receives file name and saves inmemory storage to it.
func SaveDBToStorage(filename string){
	jb := json.PackMapToJSON((*storage.GetMap()))
	j := string(jb)

	file.CreateAndWriteString(filename, j)
	log.Info.Printf("DB successfully saved to %s", filename)
}
