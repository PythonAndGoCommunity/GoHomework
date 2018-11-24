package json

import (
	"encoding/json"
	"NonRelDB/log"
)

// UnpackFromJSON | Receives json bytes and returns map pointer.
func UnpackFromJSON(b []byte) *map[string]string{
	m := make(map[string]string)
	err := json.Unmarshal(b, &m)
	if err != nil{
		log.Warning.Println(err.Error())
		log.Warning.Println("No bytes. Will be returned zero map")
	}
	return &m
}