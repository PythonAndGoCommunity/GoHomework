package query

import (
	"regexp"
	"strings"
	"NonRelDB/server/storage/inmemory"
	"NonRelDB/log"
)

// Get | Receives and key and returns value according its key.
func Get(key string) string {
	s := inmemory.GetStorage()
	v := (*s)[key]
	if v != ""{
		return v;
	} else {
		return "Value with this key not found"
	}
}

// Set | Set value according to key.
func Set(key string, value string) string{
	s := inmemory.GetStorage()
	(*s)[key] = value
	return "Value has changed"
}

// Del | Del value according to key.
func Del(key string) string{
	s := inmemory.GetStorage()
	v := (*s)[key]
	if v != "" {
		delete(*s,key)
		return v;
	} else {
		return "Value with this key not found"
	}
}

// Keys | Returns keys which match to pattern.
func Keys(pattern string) string {
	s := inmemory.GetStorage()
	var keys []string

	re, err := regexp.Compile(pattern)

	if err != nil {
		log.Warning.Println(err.Error())
		return "Pattern is incorrect"
	}

	for k := range (*s) {
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