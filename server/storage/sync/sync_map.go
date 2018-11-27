package sync

import (
	"sync"
	"regexp"
	"strings"
	"NonRelDB/log"
)

// Map | Map synchronized with mutex. 
type Map struct {
	sync.Mutex
	m *map[string]string
}

// GetMap | Getter for map.
func (sm *Map) GetMap() *map[string]string {
	return sm.m
}

// SetMap | Setter for map.
func (sm *Map) SetMap(m *map[string]string){
	sm.m = m
}

// Get | Receives and key and returns value according its key.
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

// Set | Set value according to key.
func (sm *Map) Set(key string, value string) string{
	sm.Lock()
	defer sm.Unlock()
	(*sm.m)[key] = value
	return "Value has changed"
}

// Del | Del value according to key.
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

// Keys | Returns keys which match to pattern.
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