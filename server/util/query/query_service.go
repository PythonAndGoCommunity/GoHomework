package query

import (
	"NonRelDB/server/util/inmemory"
)

// Get | Receives and key and returns value according its key.
func Get(key string) string {
	s := inmemory.Storage
	v := (*s)[key]
	if v != ""{
		return v;
	} else {
		return "Value with this key not found"
	}
}

// Set | Set value according to key.
func Set(key string, value string) string{
	s := inmemory.Storage
	(*s)[key] = value
	return "Value has changed"
}

// Del | Del value according to key.
func Del(key string) string{
	s := inmemory.Storage
	v := (*s)[key]
	if v != "" {
		delete(*s,key)
		return v;
	} else {
		return "Value with this key not found"
	}
}