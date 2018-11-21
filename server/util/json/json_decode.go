package json

import (
	"encoding/json"
)

// UnpackFromJSON | Receives json bytes and returns map pointer.
func UnpackFromJSON(b []byte) *map[string]string{
	var m map[string]string
	err := json.Unmarshal(b, &m)
	if err != nil{
		panic(err)
	}
	return &m
}