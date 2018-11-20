package query

import (
	"NonRelDB-server/util/inmemory"
)

// Get | Receives and key and returns value according its key.
func Get(key string) interface{}{
	s := inmemory.Storage
	v := (*s)[key]
	if v != nil{
		return v;
	} else {
		return "Value with this key not found"
	}
}

// Set | Set value according to key.
func Set(key string, value interface{}) string{
	s := inmemory.Storage
	(*s)[key] = value
	return "Value has changed"
}

// Del | Del value according to key.
func Del(key string) interface{}{
	s := inmemory.Storage
	v := (*s)[key]
	if v != nil {
		delete(*s,key)
		return v;
	} else {
		return "Value with this key not found"
	}
}