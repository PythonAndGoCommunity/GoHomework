package db

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"sync"
)

const defaultPath = "go-kvdb.db"

// DataBase struct.
type DataBase struct {
	mode string
	sync.RWMutex
	m map[string]string
}

// newDB creates a new database object.
func newDB(mode string) *DataBase {
	return &DataBase{
		mode: mode,
		m:    map[string]string{},
	}
}

// InitDB init a database according to the given mode.
func InitDB(mode string) (*DataBase, error) {
	db := newDB(mode)
	log.Println("Database mode:", mode)

	switch mode {
	case "memory":
		return db, nil
	case "disk":
		// check if the database file exists
		if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
			log.Print("The database file not found. Creating a new one")
			// create a database file if not exist
			err := db.Save()
			if err != nil {
				return nil, err
			}
		}
		// load a database from a file
		err := db.Load()
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errors.New("Unknown database mode")
	}
}

// Load database from a file on disk.
func (db *DataBase) Load() error {
	f, errFile := os.Open(defaultPath)
	if errFile != nil {
		log.Println("Error loading database from", defaultPath)
		return errFile
	}

	// decode to a database struct map
	errDec := decodeGob(f, &db.m)
	if errDec != nil {
		log.Println("Error decoding the database file")
		return errDec
	}

	log.Println("Database loaded from", defaultPath)
	return nil
}

// Save database on disk into a file.
func (db *DataBase) Save() error {
	f, errFile := os.Create(defaultPath)
	if errFile != nil {
		return errFile
	}
	defer f.Close()

	errEnc := encodeGob(f, db.m)
	if errEnc != nil {
		log.Println("Error encoding the database for saving")
		return errEnc
	}
	return nil
}

// Get returns the value of a key and the key state.
func (db *DataBase) Get(key string) (string, bool) {
	db.RLock()
	defer db.RUnlock()

	value, ok := db.m[key]
	return value, ok
}

// Set a key with the given value.
func (db *DataBase) Set(key, value string) error {
	db.Lock()
	defer db.Unlock()

	oldValue := db.m[key]
	db.m[key] = value

	// TODO refactor code repetition in Set and Delete functions,
	// maybe we need to implement 'transactions':
	// func (db *DataBase) (operation, payload) error {
	// 		db.Lock()
	// 		operation(payload)
	// 		if mode==disk { db.Save() }
	// 		revert operation if saving failed
	// 		db.Unlock()
	// }
	if db.mode == "disk" {
		err := db.Save()
		if err != nil {
			// revert old value and return error
			db.m[key] = oldValue
			return err
		}
	}

	return nil
}

// Delete deletes a key from database.
func (db *DataBase) Delete(key string) error {
	db.Lock()
	defer db.Unlock()

	oldValue := db.m[key]
	delete(db.m, key)

	if db.mode == "disk" {
		err := db.Save()
		if err != nil {
			// revert old value and return error
			db.m[key] = oldValue
			return err
		}
	}

	return nil
}

// Keys returns all keys matching pattern.
// 
// Bug?: what if a key itself has '*'?
func (db *DataBase) Keys(pattern string) ([]string, error) {
	db.RLock()
	defer db.RUnlock()

	var result []string
	for k := range db.m {
		m, err := globMatch(pattern, k)
		if err != nil {
			return nil, err
		}
		if m {
			result = append(result, k)
		}
	}
	return result, nil
}

// Dump serializes database data into a json format.
func (db *DataBase) Dump() ([]byte, error) {
	return json.Marshal(db.m)
}

// encodeGob encodes using gob.
func encodeGob(r io.Writer, object interface{}) error {
	encoder := gob.NewEncoder(r)
	return encoder.Encode(object)
}

// decodeGob decodes using gob.
func decodeGob(r io.Reader, object interface{}) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(object)
}

var globsRegex = regexp.MustCompile(`\*+`)

// globMatch reports whether the string contains any match of the pattern.
// Pattern could include '*' symbol which matches zero or more characters.
func globMatch(pattern string, s string) (bool, error) {
	// normalize pattern by replacing all \*+ with .*
	p := globsRegex.ReplaceAllString(pattern, ".*")
	// compile a regex for matching
	r, err := regexp.Compile("^" + p + "$")
	return r.MatchString(s), err
}
