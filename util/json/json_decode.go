package json

import (
	"NonRelDB/log"
	"encoding/json"
)

// UnpackFromJSON receives json bytes and returns map pointer.
func UnpackFromJSON(bytes []byte) *map[string]string {
	kvMap := make(map[string]string)
	err := json.Unmarshal(bytes, &kvMap)
	if err != nil {
		log.Warning.Println(err.Error())
		log.Warning.Println("No bytes. Will be returned zero map")
	}
	return &kvMap
}
