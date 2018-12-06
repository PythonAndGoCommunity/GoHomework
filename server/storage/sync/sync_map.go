package sync

import (
	"sync"
	"regexp"
	"strings"
	"NonRelDB/log"
)

// Map map synchronized with mutex. 
type Map struct {
	sync.Mutex
	m *map[string]string
}

// GetMap getter for map.
func (sm *Map) GetMap() *map[string]string {
	return sm.m
}

// SetMap setter for map.
func (sm *Map) SetMap(m *map[string]string){
	sm.m = m
}

// Get receives and key and returns value according its key.
func (sm *Map) Get(key string) string {
	sm.Lock()
	defer sm.Unlock()
	v := (*sm.m)[key]
	if v != ""{
		return v;
	} else {
		return "Value with this key not found"
	}
}

// Set set value according to key.
func (sm *Map) Set(key string, value string) string{
	sm.Lock()
	defer sm.Unlock()
	(*sm.m)[key] = value
	return "Value has changed"
}

// Del del value according to key.
func (sm *Map) Del(key string) string{
	sm.Lock()
	defer sm.Unlock()
	v := (*sm.m)[key]
	if v != "" {
		delete((*sm.m),key)
		return v;
	} else {
		return "Value with this key not found"
	}
}

// Keys returns keys which match to pattern.
func (sm *Map) Keys(pattern string) string {
	sm.Lock()
	defer sm.Unlock()

	var keys []string

	re, err := regexp.Compile(pattern)

	if err != nil {
		log.Warning.Println(err.Error())
		return "Pattern is incorrect"
	}

	for k := range ((*sm.m)) {
		if re.MatchString(k) {
			keys = append(keys, k)
		}
	}

	if keys != nil {
		return strings.Join(keys,",")
	} else {
		return "Keys with this pattern not found"
	}
}