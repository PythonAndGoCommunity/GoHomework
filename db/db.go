package db

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const defaultPath = "./go-kvdb.db"

type mode int

const (
	disk mode = iota + 1
	memory
)

// DataBase ...
type DataBase struct {
	sync.RWMutex
	m    map[string]string
	mod mode
}

// NewDB ...
func NewDB(mod mode) *DataBase {
	return &DataBase{m: map[string]string{"a": "1", "b": "2"}}
}

// Open ...
func (db *DataBase) Open(path string) error {
	db.Lock()
	defer db.Unlock()

	err := readGob(path, &db.m)
	fmt.Print(err)
	return err
}

// Save ...
func (db *DataBase) Save() error {
	db.Lock()
	defer db.Unlock()

	// f, _ := db.Dump()
	// err := ioutil.WriteFile("./go-kvdb.db", f, 0644)
	// if err != nil {
	// 	return err
	// }
	// return nil

	// path := "./go-kvdb.db"
	err := writeGob(path, db.m)
	return err
}

// Get ...
func (db *DataBase) Get(key string) (string, bool) {
	db.RLock()
	defer db.RUnlock()

	value, ok := db.m[key]
	fmt.Println("getting", key, ": ", value, ok)
	return value, ok
}

// Set ...
func (db *DataBase) Set(key, value string) {
	db.Lock()
	defer db.Unlock()

	db.m[key] = value
}

// Delete ...
func (db *DataBase) Delete(key string) {
	db.Lock()
	defer db.Unlock()

	delete(db.m, key)
}

// TODO add wildcard support

// Keys ...
func (db *DataBase) Keys() []string {
	db.RLock()
	defer db.RUnlock()

	var result []string
	for k := range db.m {
		result = append(result, k)

	}
	return result
}

// Dump ...
func (db *DataBase) Dump() ([]byte, error) {
	result, err := json.Marshal(db.m)
	return result, err
}

// func (db *DataBase) Restore(path string) error {
// 	result, err := json.Marshal(db.m)
// 	return result, err
// }

// Serialize object using gob and save result to disk.
func writeGob(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Read object from disk and deserialize it using gob.
func readGob(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
