package inmemory

import (
	"NonRelDB/server/util/json"
	"NonRelDB/server/util/file"
)

// Storage | Global variable for kv storage.
var Storage *map[string]string

// InitDBInMemory | Init kv db in memory.
func InitDBInMemory(){
	s := make(map[string]string)
	Storage = &s
}

// InitDBFromStorage | Receives filename and load its content to inmemory storage.
func InitDBFromStorage(filename string){
	j := file.OpenAndReadString(filename)
	jb := []byte(j)
	Storage = json.UnpackFromJSON(jb)
}

// SaveDBToStorage | Receives file name and saves inmemory storage to it.
func SaveDBToStorage(filename string){
	jb := json.PackMapToJSON(*Storage)
	j := string(jb)
	file.CreateAndWriteString(filename, j)
}
