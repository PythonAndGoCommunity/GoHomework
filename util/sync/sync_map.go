package sync

import (
	"NonRelDB/log"
	"regexp"
	"strings"
	"sync"
)

// Map map synchronized with mutex.
type Map struct {
	sync.Mutex
	storage *map[string]string
}

// GetMap getter for map.
func (syncMap *Map) GetMap() *map[string]string {
	return syncMap.storage
}

// SetMap setter for map.
func (syncMap *Map) SetMap(storage *map[string]string) {
	syncMap.storage = storage
}

// Get receives and key and returns value according its key.
func (syncMap *Map) Get(key string) string {
	syncMap.Lock()
	defer syncMap.Unlock()
	v := (*syncMap.storage)[key]
	if v != "" {
		return v
	}
	return "Value with this key not found"
}

// Set set value according to key.
func (syncMap *Map) Set(key string, value string) string {
	syncMap.Lock()
	defer syncMap.Unlock()
	(*syncMap.storage)[key] = value
	return "Value has changed"
}

// Del del value according to key.
func (syncMap *Map) Del(key string) string {
	syncMap.Lock()
	defer syncMap.Unlock()
	v := (*syncMap.storage)[key]
	if v != "" {
		delete((*syncMap.storage), key)
		return v
	}
	return "Value with this key not found"
}

// Keys returns keys which match to pattern.
func (syncMap *Map) Keys(pattern string) string {
	syncMap.Lock()
	defer syncMap.Unlock()

	var keys []string

	regex, err := regexp.Compile(pattern)

	if err != nil {
		log.Warning.Println(err.Error())
		return "Pattern is incorrect"
	}

	for key := range *syncMap.storage {
		if regex.MatchString(key) {
			keys = append(keys, key)
		}
	}

	if keys != nil {
		return strings.Join(keys, ",")
	}
	return "Keys with this pattern not found"
}
