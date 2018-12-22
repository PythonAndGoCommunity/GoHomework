package inmemory

import (
	"NonRelDB/log"
	"NonRelDB/util/file"
	"NonRelDB/util/json"
	"NonRelDB/util/sync"
	"os"
)

// Global variable for kv storage.
var storage *sync.Map

// GetStorage getter for storage.
func GetStorage() *sync.Map {
	return storage
}

// SetStorage setter for storage
func SetStorage(syncMap *sync.Map) {
	storage = syncMap
}

// InitDBInMemory init kv db in memory.
func InitDBInMemory() {
	storage = &sync.Map{}
	syncMap := make(map[string]string)
	storage.SetMap(&syncMap)
	log.Info.Println("DB successfully created in-memory")
}

// InitDBFromStorage receives filename and load its content to inmemory storage.
func InitDBFromStorage(filename string) {
	storage = &sync.Map{}
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

	jsonString := file.OpenAndReadString(filename)
	jsonBytes := []byte(jsonString)
	storage.SetMap(json.UnpackFromJSON(jsonBytes))
	log.Info.Printf("DB successfully initialized from %s", filename)
}

// RestoreDBFromDump restores db from received dump.
func RestoreDBFromDump(dump []byte) {
	storage.SetMap(json.UnpackFromJSON(dump))
}

// SaveDBToStorage receives file name and saves inmemory storage to it.
func SaveDBToStorage(filename string) {
	jsonBytes := json.PackMapToJSON((*storage.GetMap()))
	jsonString := string(jsonBytes)

	file.CreateAndWriteString(filename, jsonString)
	log.Info.Printf("DB successfully saved to %s", filename)
}
