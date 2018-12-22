package json

import (
	"NonRelDB/log"
	"encoding/json"
)

// PackToJSON receives string key and interface value and returns json bytes.
func PackToJSON(key string, value string) []byte {
	kvMap := map[string]string{key: value}
	bytes, err := json.Marshal(kvMap)
	if err != nil {
		log.Error.Println(err.Error())
	}
	return bytes
}

// PackToJSONIndent receives string key and interface value and returns json bytes with indent.
func PackToJSONIndent(key string, value string) []byte {
	kvMap := map[string]string{key: value}
	bytes, err := json.MarshalIndent(kvMap, "", " ")
	if err != nil {
		log.Error.Println(err.Error())
	}
	return bytes
}

// PackMapToJSON receives map and returns json bytes.
func PackMapToJSON(kvMap map[string]string) []byte {
	bytes, err := json.Marshal(kvMap)
	if err != nil {
		log.Error.Println(err.Error())
	}
	return bytes
}

// PackMapToJSONIndent receives map and returns json bytes with indent.
func PackMapToJSONIndent(kvMap map[string]string) []byte {
	bytes, err := json.MarshalIndent(kvMap, "", " ")
	if err != nil {
		log.Error.Println(err.Error())
	}
	return bytes
}
